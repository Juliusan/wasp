// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainimpl

import (
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/chain/consensus"
	"github.com/iotaledger/wasp/packages/chain/messages"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/request"
	"github.com/iotaledger/wasp/packages/peering"
	"golang.org/x/xerrors"
	//	"github.com/iotaledger/wasp/packages/state"
	//	"github.com/iotaledger/wasp/packages/util"
)

func (c *chainObj) handleMessagesLoop() {
	nDismissChainMsgChannel := c.dismissChainMsgChannel.Out()
	nStateMsgChannel := c.stateMsgChannel.Out()
	nOffLedgerRequestPeerMsgChannel := c.offLedgerRequestPeerMsgChannel.Out()
	nRequestAckPeerMsgChannel := c.requestAckPeerMsgChannel.Out()
	nMissingRequestIDsPeerMsgChannel := c.missingRequestIDsPeerMsgChannel.Out()
	nMissingRequestPeerMsgChannel := c.missingRequestPeerMsgChannel.Out()
	nTimerTickMsgChannel := c.timerTickMsgChannel.Out()
	for {
		select {
		case msg, ok := <-nDismissChainMsgChannel:
			if ok {
				c.handleDismissChain(msg.(DismissChainMsg))
			} else {
				nDismissChainMsgChannel = nil
			}
		case msg, ok := <-nStateMsgChannel:
			if ok {
				c.handleLedgerState(msg.(*messages.StateMsg))
			} else {
				nStateMsgChannel = nil
			}
		case msg, ok := <-nOffLedgerRequestPeerMsgChannel:
			if ok {
				c.handleOffLedgerRequestPeerMsg(msg.(*offLedgerRequestMsgIn))
			} else {
				nOffLedgerRequestPeerMsgChannel = nil
			}
		case msg, ok := <-nRequestAckPeerMsgChannel:
			if ok {
				c.handleRequestAckPeerMsg(msg.(*requestAckMsgIn))
			} else {
				nRequestAckPeerMsgChannel = nil
			}
		case msg, ok := <-nMissingRequestIDsPeerMsgChannel:
			if ok {
				c.handleMissingRequestIDsPeerMsg(msg.(*consensus.MissingRequestIDsMsgIn))
			} else {
				nMissingRequestIDsPeerMsgChannel = nil
			}
		case msg, ok := <-nMissingRequestPeerMsgChannel:
			if ok {
				c.handleMissingRequestPeerMsg(msg.(*missingRequestMsg))
			} else {
				nMissingRequestPeerMsgChannel = nil
			}
		case msg, ok := <-nTimerTickMsgChannel:
			if ok {
				c.handleTimerTick(msg.(messages.TimerTick))
			} else {
				nTimerTickMsgChannel = nil
			}
		}
		if nDismissChainMsgChannel == nil &&
			nStateMsgChannel == nil &&
			nOffLedgerRequestPeerMsgChannel == nil &&
			nRequestAckPeerMsgChannel == nil &&
			nMissingRequestIDsPeerMsgChannel == nil &&
			nMissingRequestPeerMsgChannel == nil &&
			nTimerTickMsgChannel == nil {
			return
		}
	}
}

func (c *chainObj) EnqueueDismissChain(reason string) {
	c.dismissChainMsgChannel.In() <- DismissChainMsg{Reason: reason}
}

func (c *chainObj) handleDismissChain(msg DismissChainMsg) {
	c.Dismiss(msg.Reason)
}

func (c *chainObj) EnqueueLedgerState(chainOutput *ledgerstate.AliasOutput, timestamp time.Time) {
	c.stateMsgChannel.In() <- &messages.StateMsg{
		ChainOutput: chainOutput,
		Timestamp:   timestamp,
	}
}

// handleLedgerState processes the only chain output which exists on the chain's address
// If necessary, it creates/changes/rotates committee object
func (c *chainObj) handleLedgerState(msg *messages.StateMsg) {
	sh, err := hashing.HashValueFromBytes(msg.ChainOutput.GetStateData())
	if err != nil {
		c.log.Error(xerrors.Errorf("parsing state hash: %w", err))
		return
	}
	c.log.Debugf("processStateMessage. stateIndex: %d, stateHash: %s, stateAddr: %s, state transition: %v",
		msg.ChainOutput.GetStateIndex(), sh.String(),
		msg.ChainOutput.GetStateAddress().Base58(), !msg.ChainOutput.GetIsGovernanceUpdated(),
	)
	cmt := c.getCommittee()

	if cmt != nil {
		err = c.rotateCommitteeIfNeeded(msg.ChainOutput, cmt)
	} else {
		err = c.createCommitteeIfNeeded(msg.ChainOutput)
	}
	if err != nil {
		c.log.Errorf("processStateMessage: %v", err)
		return
	}
	c.stateMgr.EventStateMsg(msg)
}

func (c *chainObj) EnqueueOffLedgerRequestPeerMsg(req *request.OffLedger, senderNetID string) {
	c.offLedgerRequestPeerMsgChannel.In() <- &offLedgerRequestMsgIn{
		offLedgerRequestMsg: offLedgerRequestMsg{Req: req},
		SenderNetID:         senderNetID,
	}
}

func (c *chainObj) handleOffLedgerRequestPeerMsg(msg *offLedgerRequestMsgIn) {
	req := msg.Req
	senderNetID := msg.SenderNetID
	c.log.Debugf("handleOffLedgerRequestPeerMsg: reqID: %s, peerID: %s", req.ID().Base58(), senderNetID)
	c.sendRequestAcknowledgementMsg(req.ID(), senderNetID)
	if !c.mempool.ReceiveRequest(req) {
		return
	}
	c.log.Debugf("handleOffLedgerRequestPeerMsg - added to mempool: reqID: %s, peerID: %s", req.ID().Base58(), senderNetID)
	c.broadcastOffLedgerRequest(req)
}

func (c *chainObj) enqueueRequestAckPeerMsg(requestID *iscp.RequestID, senderNetID string) {
	c.requestAckPeerMsgChannel.In() <- &requestAckMsgIn{
		requestAckMsg: requestAckMsg{ReqID: requestID},
		SenderNetID:   senderNetID,
	}
}

func (c *chainObj) handleRequestAckPeerMsg(msg *requestAckMsgIn) {
	c.ReceiveRequestAckMessage(msg.ReqID, msg.SenderNetID)
}

func (c *chainObj) enqueueMissingRequestIDsPeerMsg(ids []iscp.RequestID, senderNetID string) {
	c.missingRequestIDsPeerMsgChannel.In() <- &consensus.MissingRequestIDsMsgIn{
		MissingRequestIDsMsg: consensus.MissingRequestIDsMsg{IDs: ids},
		SenderNetID:          senderNetID,
	}
}

func (c *chainObj) handleMissingRequestIDsPeerMsg(msg *consensus.MissingRequestIDsMsgIn) {
	if !c.pullMissingRequestsFromCommittee {
		return
	}
	peerID := msg.SenderNetID
	for _, reqID := range msg.IDs {
		c.log.Debugf("Sending MissingRequestsToPeer: reqID: %s, peerID: %s", reqID.Base58(), peerID)
		if req := c.mempool.GetRequest(reqID); req != nil {
			c.SendMsgByNetID(peerID, peering.PeerMessagePartyChain, &missingRequestMsg{
				Request: req,
			})
		}
	}
}

func (c *chainObj) enqueueMissingRequestPeerMsg(request iscp.Request) {
	c.missingRequestPeerMsgChannel.In() <- &missingRequestMsg{Request: request}
}

func (c *chainObj) handleMissingRequestPeerMsg(msg *missingRequestMsg) {
	if !c.pullMissingRequestsFromCommittee {
		return
	}
	if c.consensus.ShouldReceiveMissingRequest(msg.Request) {
		c.mempool.ReceiveRequest(msg.Request)
	}
}

func (c *chainObj) enqueueTimerTick(tick int) {
	c.timerTickMsgChannel.In() <- messages.TimerTick(tick)
}

func (c *chainObj) handleTimerTick(msg messages.TimerTick) {
	if msg%2 == 0 {
		c.stateMgr.EventTimerMsg(msg / 2)
	} else if c.consensus != nil {
		c.consensus.EventTimerMsg(msg / 2)
	}
	if msg%40 == 0 {
		stats := c.mempool.Info()
		c.log.Debugf("mempool total = %d, ready = %d, in = %d, out = %d", stats.TotalPool, stats.ReadyCounter, stats.InPoolCounter, stats.OutPoolCounter)
	}
}

/*// EventGetBlockMsg is a request for a block while syncing
func (sm *stateManager) EventGetBlockMsg(msg *messages.GetBlockMsg) {
	sm.eventGetBlockMsgCh <- msg
}

func (sm *stateManager) eventGetBlockMsg(msg *messages.GetBlockMsg) {
	sm.log.Debugw("EventGetBlockMsg received: ",
		"sender", msg.SenderNetID,
		"block index", msg.BlockIndex,
	)
	if sm.stateOutput == nil { // Not a necessary check, only for optimization.
		sm.log.Debugf("EventGetBlockMsg ignored: stateOutput is nil")
		return
	}
	if msg.BlockIndex > sm.stateOutput.GetStateIndex() { // Not a necessary check, only for optimization.
		sm.log.Debugf("EventGetBlockMsg ignored 1: block #%d not found. Current state index: #%d",
			msg.BlockIndex, sm.stateOutput.GetStateIndex())
		return
	}
	blockBytes, err := state.LoadBlockBytes(sm.store, msg.BlockIndex)
	if err != nil {
		sm.log.Errorf("EventGetBlockMsg: LoadBlockBytes: %v", err)
		return
	}
	if blockBytes == nil {
		sm.log.Debugf("EventGetBlockMsg ignored 2: block #%d not found. Current state index: #%d",
			msg.BlockIndex, sm.stateOutput.GetStateIndex())
		return
	}

	sm.log.Debugf("EventGetBlockMsg for state index #%d --> responding to peer %s", msg.BlockIndex, msg.SenderNetID)

	sm.peers.SendSimple(msg.SenderNetID, messages.MsgBlock, util.MustBytes(&messages.BlockMsg{
		BlockBytes: blockBytes,
	}))
}

// EventBlockMsg
func (sm *stateManager) EventBlockMsg(msg *messages.BlockMsg) {
	sm.eventBlockMsgCh <- msg
}

func (sm *stateManager) eventBlockMsg(msg *messages.BlockMsg) {
	sm.log.Debugf("EventBlockMsg received from %v", msg.SenderNetID)
	if sm.stateOutput == nil {
		sm.log.Debugf("EventBlockMsg ignored: stateOutput is nil")
		return
	}
	block, err := state.BlockFromBytes(msg.BlockBytes)
	if err != nil {
		sm.log.Warnf("EventBlockMsg ignored: wrong block received from peer %s. Err: %v", msg.SenderNetID, err)
		return
	}
	sm.log.Debugw("EventBlockMsg from ",
		"sender", msg.SenderNetID,
		"block index", block.BlockIndex(),
		"approving output", iscp.OID(block.ApprovingOutputID()),
	)
	if sm.addBlockFromPeer(block) {
		sm.takeAction()
	}
}

func (sm *stateManager) EventOutputMsg(msg ledgerstate.Output) {
	sm.eventOutputMsgCh <- msg
}

func (sm *stateManager) eventOutputMsg(msg ledgerstate.Output) {
	sm.log.Debugf("EventOutputMsg received: %s", iscp.OID(msg.ID()))
	chainOutput, ok := msg.(*ledgerstate.AliasOutput)
	if !ok {
		sm.log.Debugf("EventOutputMsg ignored: output is of type %t, expecting *ledgerstate.AliasOutput", msg)
		return
	}
	if sm.outputPulled(chainOutput) {
		sm.takeAction()
	}
}

// EventStateTransactionMsg triggered whenever new state transaction arrives
// the state transaction may be confirmed or not
func (sm *stateManager) EventStateMsg(msg *messages.StateMsg) {
	sm.eventStateOutputMsgCh <- msg
}

func (sm *stateManager) eventStateMsg(msg *messages.StateMsg) {
	sm.log.Debugw("EventStateMsg received: ",
		"state index", msg.ChainOutput.GetStateIndex(),
		"chainOutput", iscp.OID(msg.ChainOutput.ID()),
	)
	stateHash, err := hashing.HashValueFromBytes(msg.ChainOutput.GetStateData())
	if err != nil {
		sm.log.Errorf("EventStateMsg ignored: failed to parse state hash: %v", err)
		return
	}
	sm.log.Debugf("EventStateMsg state hash is %v", stateHash.String())
	if sm.stateOutputReceived(msg.ChainOutput, msg.Timestamp) {
		sm.takeAction()
	}
}

func (sm *stateManager) EventStateCandidateMsg(state state.VirtualStateAccess, outputID ledgerstate.OutputID) {
	sm.eventStateCandidateMsgCh <- &messages.StateCandidateMsg{
		State:             state,
		ApprovingOutputID: outputID,
	}
}

func (sm *stateManager) eventStateCandidateMsg(msg *messages.StateCandidateMsg) {
	sm.log.Debugf("EventStateCandidateMsg received: state index: %d, timestamp: %v",
		msg.State.BlockIndex(), msg.State.Timestamp(),
	)
	if sm.stateOutput == nil {
		sm.log.Debugf("EventStateCandidateMsg ignored: stateOutput is nil")
		return
	}
	if sm.addStateCandidateFromConsensus(msg.State, msg.ApprovingOutputID) {
		sm.takeAction()
	}
}

func (sm *stateManager) EventTimerMsg(msg messages.TimerTick) {
	if msg%2 == 0 {
		sm.eventTimerMsgCh <- msg
	}
}

func (sm *stateManager) eventTimerMsg() {
	sm.log.Debugf("EventTimerMsg received")
	sm.takeAction()
}
*/
