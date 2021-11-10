// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// Package commonsubset wraps the ACS (Asynchronous Common Subset) part from
// the HoneyBadgerBFT with the Peering for communication layer. The main objective
// here is to ensure eventual delivery of all the messages, as that is the assumption
// of the ACS algorithm.
//
// To ensure the eventual delivery of messages, we are resending them until they
// are acknowledged by a peer. The acknowledgements are carried along with other
// messages to decrease communication overhead.
//
// We are using a forked version of https://github.com/anthdm/hbbft. The fork was
// made to gave misc mistakes fixed that made apparent in the case of out-of-order
// message delivery.
//
// The main references:
//
//  - https://eprint.iacr.org/2016/199.pdf
//    The original HBBFT paper.
//
//  - https://dl.acm.org/doi/10.1145/2611462.2611468
//    Here the BBA part is originally presented.
//    Most of the mistakes in the lib was in this part.
//
//  - https://dl.acm.org/doi/10.1145/197917.198088
//    The definition of ACS is by Ben-Or. At least looks so.
//

//nolint:dupl // TODO there is a bunch of duplicated code in this file, should be refactored to reusable funcs
package commonsubset

import (
	"encoding/binary"
	"time"

	"github.com/anthdm/hbbft"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/chain/consensus/commoncoin"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/tcrypto"
	"golang.org/x/xerrors"
)

const (
	acsMsgType   = iota
	resendPeriod = 500 * time.Millisecond
)

// CommonSubject is responsible for executing a single instance of agreement.
// The main task for this object is to manage all the redeliveries and acknowledgements.
type CommonSubset struct {
	impl *hbbft.ACS // The ACS implementation.

	batchCounter   uint32               // Used to assign unique numbers to the messages.
	missingAcks    map[uint32]time.Time // Not yet acknowledged batches sent by us and the time of the last resend.
	pendingAcks    [][]uint32           // pendingAcks[PeerID] -- acks to be sent to the peer.
	recvMsgBatches []map[uint32]uint32  // [PeerId]ReceivedMbID -> AckedInMb (0 -- unacked). It remains 0, if acked in non-data message.
	sentMsgBatches map[uint32]*msgBatch // BatchID -> MsgBatch

	sessionID   uint64 // Unique identifier for this consensus instance. Used to route messages.
	stateIndex  uint32 // Sequence number of the CS transaction we are agreeing on.
	committee   chain.Committee
	netOwnIndex uint16 // Our index in the group.

	inputCh  chan []byte            // For our input to the consensus.
	recvCh   chan *msgBatch         // For incoming messages.
	closeCh  chan bool              // To implement closing of the object.
	outputCh chan map[uint16][]byte // The caller will receive its result via this channel.
	done     bool                   // Indicates, if the decision is already made.
	log      *logger.Logger         // Logger, of course.
}

func NewCommonSubset(
	sessionID uint64,
	stateIndex uint32,
	committee chain.Committee,
	netGroup peering.GroupProvider,
	dkShare *tcrypto.DKShare,
	allRandom bool, // Set to true to have real CC rounds for each epoch. That's for testing mostly.
	outputCh chan map[uint16][]byte,
	log *logger.Logger,
) (*CommonSubset, error) {
	ownIndex := netGroup.SelfIndex()
	allNodes := netGroup.AllNodes()
	nodeCount := len(allNodes)
	nodes := make([]uint64, nodeCount)
	nodePos := 0
	for ni := range allNodes {
		nodes[nodePos] = uint64(ni)
		nodePos++
	}
	if outputCh == nil {
		outputCh = make(chan map[uint16][]byte, 1)
	}
	var salt [8]byte
	binary.BigEndian.PutUint64(salt[:], sessionID)
	acsCfg := hbbft.Config{
		N:          nodeCount,
		F:          nodeCount - int(dkShare.T),
		ID:         uint64(ownIndex),
		Nodes:      nodes,
		BatchSize:  0, // Unused in ACS.
		CommonCoin: commoncoin.NewBlsCommonCoin(dkShare, salt[:], allRandom),
	}
	cs := CommonSubset{
		impl:           hbbft.NewACS(acsCfg),
		batchCounter:   0,
		missingAcks:    make(map[uint32]time.Time),
		pendingAcks:    make([][]uint32, nodeCount),
		recvMsgBatches: make([]map[uint32]uint32, nodeCount),
		sentMsgBatches: make(map[uint32]*msgBatch),
		sessionID:      sessionID,
		stateIndex:     stateIndex,
		committee:      committee,
		netOwnIndex:    ownIndex,
		inputCh:        make(chan []byte, 1),
		recvCh:         make(chan *msgBatch, 1),
		closeCh:        make(chan bool),
		outputCh:       outputCh,
		log:            log,
	}
	for i := range cs.recvMsgBatches {
		cs.recvMsgBatches[i] = make(map[uint32]uint32)
	}
	for i := range cs.pendingAcks {
		cs.pendingAcks[i] = make([]uint32, 0)
	}
	go cs.run()
	return &cs, nil
}

// OutputCh returns the output channel, so that the called don't need to track it.
func (cs *CommonSubset) OutputCh() chan map[uint16][]byte {
	return cs.outputCh
}

// Input accepts the current node's proposal for the consensus.
func (cs *CommonSubset) Input(input []byte) {
	cs.inputCh <- input
}

// HandleMsgBatch accepts a parsed msgBatch as an input from other node.
// This function is used in the CommonSubsetCoordinator to avoid parsing
// the received message multiple times.
func (cs *CommonSubset) HandleMsgBatch(msg peering.Serializable) {
	defer func() {
		if err := recover(); err != nil {
			// Just to avoid panics on writing to a closed channel.
			// This can happen on the ACS termination.
			cs.log.Warnf("Dropping msgBatch reason=%v", err)
		}
	}()
	cs.recvCh <- msg.(*msgBatch)
}

func (cs *CommonSubset) Close() {
	cs.impl.Stop()
	close(cs.closeCh)
	close(cs.recvCh)
	close(cs.inputCh)
}

func (cs *CommonSubset) run() {
	retry := time.After(resendPeriod)
	for {
		select {
		case <-retry:
			cs.timeTick()
			if !cs.done || len(cs.pendingAcks) != 0 {
				// The condition for stopping the resend is a bit tricky, because it is
				// not enough for this node to complete with the decision. This node
				// must help others to decide as well.
				retry = time.After(resendPeriod)
			}
		case input, ok := <-cs.inputCh:
			if !ok {
				return
			}
			if err := cs.handleInput(input); err != nil {
				cs.log.Errorf("Failed to handle input, reason: %v", err)
			}
		case mb, ok := <-cs.recvCh:
			if !ok {
				return
			}
			cs.handleMsgBatch(mb)
		case <-cs.closeCh:
			return
		}
	}
}

func (cs *CommonSubset) timeTick() {
	if cs.log.Desugar().Core().Enabled(logger.LevelDebug) { // Debug info for the ACS.
		rbcInstances, bbaInstances, rbcResults, bbaResults, msgQueue := cs.impl.DebugInfo()
		cs.log.Debugf("ACS::timeTick[%v], sessionId=%v, done=%v queue.len=%v, bba.res=%v", cs.impl.ID, cs.sessionID, cs.done, len(msgQueue), bbaResults)
		for i, bba := range bbaInstances {
			if !bba.Done() {
				cs.log.Debugf("ACS::timeTick[%v], sessionId=%v, impl.bba[%v]=%+v", cs.impl.ID, cs.sessionID, i, bba)
			}
		}
		for i, rbc := range rbcInstances {
			if _, ok := rbcResults[i]; !ok {
				cs.log.Debugf("ACS::timeTick[%v], sessionId=%v, impl.rbc[%v]=%+v", cs.impl.ID, cs.sessionID, i, rbc)
			}
		}
	}

	now := time.Now()
	resentBefore := now.Add(resendPeriod * (-2))
	for missingAck, lastSentTime := range cs.missingAcks {
		if lastSentTime.Before(resentBefore) {
			if mb, ok := cs.sentMsgBatches[missingAck]; ok {
				cs.send(mb)
				cs.missingAcks[missingAck] = now
			} else {
				// Batch is already cleaned up, so we don't need the ack.
				delete(cs.missingAcks, missingAck)
			}
		}
	}
}

func (cs *CommonSubset) handleInput(input []byte) error {
	var err error
	if err = cs.impl.InputValue(input); err != nil {
		return xerrors.Errorf("Failed to process ACS.InputValue: %w", err)
	}
	cs.sendPendingMessages()
	return nil
}

func (cs *CommonSubset) handleMsgBatch(recvBatch *msgBatch) {
	//
	// Cleanup all the acknowledged messages from the missing acks list.
	for _, receivedAck := range recvBatch.acks {
		delete(cs.missingAcks, receivedAck)
	}
	//
	// Resend the old responses, if ack was not received
	// and the message is already processed.
	if ackIn, ok := cs.recvMsgBatches[recvBatch.src][recvBatch.id]; ok {
		if ackIn == 0 {
			// Make a fake message just to acknowledge the message, because sender has
			// retried it already and we haven't sent the ack yet. We will not expect
			// an acknowledgement for this ack.
			cs.send(cs.makeAckOnlyBatch(recvBatch.src, recvBatch.id))
		} else {
			cs.send(cs.sentMsgBatches[ackIn])
		}
		return
	}
	//
	// If we have completed with the decision, we just acknowledging other's messages.
	// That is needed to help other members to decide on termination.
	if cs.done {
		if recvBatch.NeedsAck() {
			cs.send(cs.makeAckOnlyBatch(recvBatch.src, recvBatch.id))
		}
		return
	}
	if recvBatch.NeedsAck() {
		// Batches with id=0 are for acks only, they contain no data,
		// therefore we will not track them, just respond on the fly.
		// Otherwise we store the batch ID to be acknowledged with the
		// next outgoing message batch to that peer.
		cs.recvMsgBatches[recvBatch.src][recvBatch.id] = 0 // Received, not acknowledged yet.
		cs.pendingAcks[recvBatch.src] = append(cs.pendingAcks[recvBatch.src], recvBatch.id)
	}
	//
	// Process the messages.
	for _, m := range recvBatch.msgs {
		if err := cs.impl.HandleMessage(uint64(recvBatch.src), m); err != nil {
			cs.log.Errorf("Failed to handle message: %v, message=%+v", err, m)
		}
	}
	//
	// Send the outgoing messages.
	cs.sendPendingMessages()
	//
	// Check, maybe we are done.
	if cs.impl.Done() {
		var output map[uint64][]byte = cs.impl.Output()
		out16 := make(map[uint16][]byte)
		for index, share := range output {
			out16[uint16(index)] = share
		}
		cs.done = true
		cs.outputCh <- out16
	}
}

func (cs *CommonSubset) sendPendingMessages() {
	var outBatches []*msgBatch
	var err error
	if outBatches, err = cs.makeBatches(cs.impl.Messages()); err != nil {
		cs.log.Errorf("Failed to make out batch: %v", err)
	}
	now := time.Now()
	for _, b := range outBatches {
		b.acks = cs.pendingAcks[b.dst]
		cs.pendingAcks[b.dst] = make([]uint32, 0)
		for i := range b.acks {
			// Update the reverse index, i.e. for each message,
			// specify message which have acknowledged it.
			cs.recvMsgBatches[b.dst][b.acks[i]] = b.id
		}
		cs.sentMsgBatches[b.id] = b
		cs.missingAcks[b.id] = now
		cs.send(b)
	}
}

func (cs *CommonSubset) makeBatches(msgs []hbbft.MessageTuple) ([]*msgBatch, error) {
	batchMsgs := make([][]*hbbft.ACSMessage, cs.impl.N)
	for i := range batchMsgs {
		batchMsgs[i] = make([]*hbbft.ACSMessage, 0)
	}
	for _, m := range msgs {
		if acsMsg, ok := m.Payload.(*hbbft.ACSMessage); ok {
			batchMsgs[m.To] = append(batchMsgs[m.To], acsMsg)
		} else {
			return nil, xerrors.Errorf("unexpected message payload type: %T", m.Payload)
		}
	}
	batches := make([]*msgBatch, 0)
	for i := range batchMsgs {
		if len(batchMsgs[i]) == 0 {
			continue
		}
		cs.batchCounter++
		batches = append(batches, &msgBatch{
			sessionID:  cs.sessionID,
			stateIndex: cs.stateIndex,
			id:         cs.batchCounter,
			src:        uint16(cs.impl.ID),
			dst:        uint16(i),
			msgs:       batchMsgs[i],
			acks:       make([]uint32, 0), // Filled later.
		})
	}
	return batches, nil
}

func (cs *CommonSubset) makeAckOnlyBatch(peerID uint16, ackMB uint32) *msgBatch {
	acks := []uint32{}
	if ackMB != 0 {
		acks = append(acks, ackMB)
	}
	if len(cs.pendingAcks[peerID]) != 0 {
		acks = append(acks, cs.pendingAcks[peerID]...)
		cs.pendingAcks[peerID] = make([]uint32, 0)
	}
	if len(acks) == 0 {
		return nil
	}
	return &msgBatch{
		sessionID:  cs.sessionID,
		stateIndex: cs.stateIndex,
		id:         0, // Do not require an acknowledgement.
		src:        uint16(cs.impl.ID),
		dst:        peerID,
		msgs:       []*hbbft.ACSMessage{},
		acks:       acks,
	}
}

func (cs *CommonSubset) send(msgBatch *msgBatch) {
	if msgBatch == nil {
		// makeAckOnlyBatch can produce nil batches, if there is nothing
		// to acknowledge. We handle that here, to avoid IFs in multiple places.
		return
	}
	cs.log.Debugf("ACS::IO - Sending a msgBatch=%+v", msgBatch)
	cs.committee.SendMsgByIndex(msgBatch.dst, peering.PeerMessagePartyAcs, msgBatch)
}
