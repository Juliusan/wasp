// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainimpl

import (
	"sync"
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	txstream "github.com/iotaledger/goshimmer/packages/txstream/client"
	"github.com/iotaledger/hive.go/events"
	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/chain/committee"
	"github.com/iotaledger/wasp/packages/chain/consensus"
	"github.com/iotaledger/wasp/packages/chain/mempool"
	"github.com/iotaledger/wasp/packages/chain/nodeconnimpl"
	"github.com/iotaledger/wasp/packages/chain/statemgr"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/subrealm"
	"github.com/iotaledger/wasp/packages/metrics"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/publisher"
	"github.com/iotaledger/wasp/packages/registry"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/processors"
	"go.uber.org/atomic"
	"golang.org/x/xerrors"
	"gopkg.in/eapache/channels.v1"
)

var (
	_ chain.Chain                    = &chainObj{}
	_ chain.ChainCore                = &chainObj{}
	_ chain.ChainEntry               = &chainObj{}
	_ chain.ChainRequests            = &chainObj{}
	_ chain.ChainEvents              = &chainObj{}
	_ peering.PeerMessageSimpleParty = &chainObj{}
	_ peering.PeerMessageGroupParty  = &chainObj{}
)

type chainMessageType byte

const (
	chainMessageTypeStart chainMessageType = iota
	chainMessageTypeMissingRequest
	chainMessageTypeOffledgerRequest
	chainMessageTypeRequestAck
	chainMessageTypeEnd
)

type chainObj struct {
	committee                        atomic.Value
	mempool                          chain.Mempool
	mempoolLastCleanedIndex          uint32
	dismissed                        atomic.Bool
	dismissOnce                      sync.Once
	chainID                          *iscp.ChainID
	chainStateSync                   coreutil.ChainStateSync
	stateReader                      state.OptimisticStateReader
	procset                          *processors.Cache
	stateMgr                         chain.StateManager
	consensus                        chain.Consensus
	log                              *logger.Logger
	nodeConn                         chain.NodeConnection
	db                               kvstore.KVStore
	peerNetworkConfig                registry.PeerNetworkConfigProvider
	netProvider                      peering.NetworkProvider
	dksProvider                      registry.DKShareRegistryProvider
	committeeRegistry                registry.CommitteeRegistryProvider
	blobProvider                     registry.BlobCache
	eventRequestProcessed            *events.Event
	eventChainTransition             *events.Event
	peers                            *peering.PeerDomainProvider
	offLedgerReqsAcksMutex           sync.RWMutex
	offLedgerReqsAcks                map[iscp.RequestID][]string
	offledgerBroadcastUpToNPeers     int
	offledgerBroadcastInterval       time.Duration
	pullMissingRequestsFromCommittee bool
	peeringID                        peering.PeeringID
	attachIDs                        []interface{}
	chainMetrics                     metrics.ChainMetrics
	dismissChainMsgChannel           *channels.InfiniteChannel
	stateMsgChannel                  *channels.InfiniteChannel
	offLedgerRequestPeerMsgChannel   *channels.InfiniteChannel
	requestAckPeerMsgChannel         *channels.InfiniteChannel
	missingRequestIDsPeerMsgChannel  *channels.InfiniteChannel
	missingRequestPeerMsgChannel     *channels.InfiniteChannel
	timerTickMsgChannel              *channels.InfiniteChannel
}

type committeeStruct struct {
	valid bool
	cmt   chain.Committee
}

func NewChain(
	chainID *iscp.ChainID,
	log *logger.Logger,
	txstreamClient *txstream.Client,
	peerNetConfig registry.PeerNetworkConfigProvider,
	db kvstore.KVStore,
	netProvider peering.NetworkProvider,
	dksProvider registry.DKShareRegistryProvider,
	committeeRegistry registry.CommitteeRegistryProvider,
	blobProvider registry.BlobCache,
	processorConfig *processors.Config,
	offledgerBroadcastUpToNPeers int,
	offledgerBroadcastInterval time.Duration,
	pullMissingRequestsFromCommittee bool,
	chainMetrics metrics.ChainMetrics,
) chain.Chain {
	log.Debugf("creating chain object for %s", chainID.String())

	chainLog := log.Named(chainID.Base58()[:6] + ".")
	chainStateSync := coreutil.NewChainStateSync()
	ret := &chainObj{
		mempool:           mempool.New(state.NewOptimisticStateReader(db, chainStateSync), blobProvider, chainLog, chainMetrics),
		procset:           processors.MustNew(processorConfig),
		chainID:           chainID,
		log:               chainLog,
		nodeConn:          nodeconnimpl.New(txstreamClient, chainLog),
		db:                db,
		chainStateSync:    chainStateSync,
		stateReader:       state.NewOptimisticStateReader(db, chainStateSync),
		peerNetworkConfig: peerNetConfig,
		netProvider:       netProvider,
		dksProvider:       dksProvider,
		committeeRegistry: committeeRegistry,
		blobProvider:      blobProvider,
		eventRequestProcessed: events.NewEvent(func(handler interface{}, params ...interface{}) {
			handler.(func(_ iscp.RequestID))(params[0].(iscp.RequestID))
		}),
		eventChainTransition: events.NewEvent(func(handler interface{}, params ...interface{}) {
			handler.(func(_ *chain.ChainTransitionEventData))(params[0].(*chain.ChainTransitionEventData))
		}),
		offLedgerReqsAcks:                make(map[iscp.RequestID][]string),
		offledgerBroadcastUpToNPeers:     offledgerBroadcastUpToNPeers,
		offledgerBroadcastInterval:       offledgerBroadcastInterval,
		pullMissingRequestsFromCommittee: pullMissingRequestsFromCommittee,
		peeringID:                        chainID.Array(),
		attachIDs:                        make([]interface{}, 0),
		chainMetrics:                     chainMetrics,
		dismissChainMsgChannel:           channels.NewInfiniteChannel(),
		stateMsgChannel:                  channels.NewInfiniteChannel(),
		offLedgerRequestPeerMsgChannel:   channels.NewInfiniteChannel(),
		requestAckPeerMsgChannel:         channels.NewInfiniteChannel(),
		missingRequestIDsPeerMsgChannel:  channels.NewInfiniteChannel(),
		missingRequestPeerMsgChannel:     channels.NewInfiniteChannel(),
		timerTickMsgChannel:              channels.NewInfiniteChannel(),
	}
	ret.committee.Store(&committeeStruct{})
	ret.eventChainTransition.Attach(events.NewClosure(ret.processChainTransition))

	peers, err := netProvider.PeerDomain(peerNetConfig.Neighbors())
	if err != nil {
		log.Errorf("NewChain: %v", err)
		return nil
	}
	ret.stateMgr = statemgr.New(db, ret, peers, ret.nodeConn)
	ret.peers = &peers
	ret.RegisterPeerMessageParty(ret)
	go ret.handleMessagesLoop()
	ret.startTimer()
	return ret
}

func (c *chainObj) GetPartyType() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyChain
}

func (c *chainObj) HandleMessage(senderNetID string, msg peering.Serializable) {
	switch msgt := msg.(type) {
	case *offLedgerRequestMsg:
		c.EnqueueOffLedgerRequestPeerMsg(msgt.Req, senderNetID)
	case *requestAckMsg:
		c.enqueueRequestAckPeerMsg(msgt.ReqID, senderNetID)
	case *missingRequestMsg:
		c.enqueueMissingRequestPeerMsg(msgt.Request)
	}
}

func (c *chainObj) HandleGroupMessage(senderNetID string, senderIndex uint16, msg peering.Serializable) {
	msgt := msg.(*consensus.MissingRequestIDsMsg)
	c.enqueueMissingRequestIDsPeerMsg(msgt.IDs, senderNetID)
}

func (c *chainObj) SendMsgByNetID(netID string, destinationParty peering.PeerMessagePartyType, msg peering.Serializable) {
	(*c.peers).SendMsgByNetID(netID, c.peeringID, destinationParty, msg)
}

func (c *chainObj) SendMsgToRandomPeers(upToNumPeers uint16, destinationParty peering.PeerMessagePartyType, msg peering.Serializable) {
	(*c.peers).SendMsgToRandomPeers(upToNumPeers, c.peeringID, destinationParty, msg)
}

// processChainTransition processes the unique chain output which exists on the chain's address
// If necessary, it creates/changes/rotates committee object
func (c *chainObj) processChainTransition(msg *chain.ChainTransitionEventData) {
	stateIndex := msg.VirtualState.BlockIndex()
	c.log.Debugf("processChainTransition: processing state %d", stateIndex)
	if !msg.ChainOutput.GetIsGovernanceUpdated() {
		c.log.Debugf("processChainTransition state %d: output %s is not governance updated; state hash %s; last cleaned state is %d",
			stateIndex, iscp.OID(msg.ChainOutput.ID()), msg.VirtualState.StateCommitment().String(), c.mempoolLastCleanedIndex)
		// normal state update:
		c.stateReader.SetBaseline()
		chainID := iscp.NewChainID(msg.ChainOutput.GetAliasAddress())
		var reqids []iscp.RequestID
		for i := c.mempoolLastCleanedIndex + 1; i <= msg.VirtualState.BlockIndex(); i++ {
			c.log.Debugf("processChainTransition state %d: cleaning state %d", stateIndex, i)
			var err error
			reqids, err = blocklog.GetRequestIDsForBlock(c.stateReader, i)
			if reqids == nil {
				// The error means a database error. The optimistic state read failure can't occur here
				// because the state transition message is only sent only after state is committed and before consensus
				// start new round
				c.log.Panicf("processChainTransition. unexpected error: %v", err)
				return // to avoid "possible nil pointer dereference" in later use of `reqids`
			}
			// remove processed requests from the mempool
			c.log.Debugf("processChainTransition state %d cleaning state %d: removing %d requests", stateIndex, i, len(reqids))
			c.mempool.RemoveRequests(reqids...)
			chain.PublishRequestsSettled(chainID, i, reqids)
			// publish events
			for _, reqid := range reqids {
				c.eventRequestProcessed.Trigger(reqid)
			}
			c.publishNewBlockEvents(stateIndex)

			c.log.Debugf("processChainTransition state %d: state %d cleaned, deleted requests: %+v",
				stateIndex, i, iscp.ShortRequestIDs(reqids))
		}
		chain.PublishStateTransition(chainID, msg.ChainOutput, len(reqids))
		chain.LogStateTransition(msg, reqids, c.log)

		c.mempoolLastCleanedIndex = stateIndex
	} else {
		c.log.Debugf("processChainTransition state %d: output %s is governance updated; state hash %s",
			stateIndex, iscp.OID(msg.ChainOutput.ID()), msg.VirtualState.StateCommitment().String())
		chain.LogGovernanceTransition(msg, c.log)
		chain.PublishGovernanceTransition(msg.ChainOutput)
	}
	c.consensus.EventStateTransitionMsg(msg.VirtualState, msg.ChainOutput, msg.OutputTimestamp)
	c.log.Debugf("processChainTransition completed: state index: %d, state hash: %s", stateIndex, msg.VirtualState.StateCommitment().String())
}

func (c *chainObj) publishNewBlockEvents(blockIndex uint32) {
	if blockIndex == 0 {
		// don't run on state #0, root contracts not initialized yet.
		return
	}

	kvPartition := subrealm.NewReadOnly(c.stateReader.KVStoreReader(), kv.Key(blocklog.Contract.Hname().Bytes()))

	evts, err := blocklog.GetBlockEventsInternal(kvPartition, blockIndex)
	if err != nil {
		c.log.Panicf("publishNewBlockEvents - something went wrong getting events for block. %v", err)
	}

	go func() {
		for _, msg := range evts {
			c.log.Debugf("publishNewBlockEvents: '%s'", msg)
			publisher.Publish("vmmsg", c.chainID.Base58(), msg)
		}
	}()
}

func (c *chainObj) rotateCommitteeIfNeeded(anchorOutput *ledgerstate.AliasOutput, currentCmt chain.Committee) error {
	if currentCmt.Address().Equals(anchorOutput.GetStateAddress()) {
		// nothing changed. no rotation
		return nil
	}
	// address changed
	if !anchorOutput.GetIsGovernanceUpdated() {
		return xerrors.Errorf("rotateCommitteeIfNeeded: inconsistency. Governance transition expected... New output: %s", anchorOutput.String())
	}
	rec, err := c.getOwnCommitteeRecord(anchorOutput.GetStateAddress())
	if err != nil {
		return xerrors.Errorf("rotateCommitteeIfNeeded: %w", err)
	}
	// rotation needed
	// close current in any case
	c.log.Infof("CLOSING COMMITTEE for %s", currentCmt.Address().Base58())

	currentCmt.Close()
	c.consensus.Close()
	c.setCommittee(nil)
	c.consensus = nil
	if rec != nil {
		// create new if committee record is available
		if err = c.createNewCommitteeAndConsensus(rec); err != nil {
			return xerrors.Errorf("rotateCommitteeIfNeeded: creating committee and consensus: %v", err)
		}
	}
	return nil
}

func (c *chainObj) createCommitteeIfNeeded(anchorOutput *ledgerstate.AliasOutput) error {
	// check if I am in the committee
	rec, err := c.getOwnCommitteeRecord(anchorOutput.GetStateAddress())
	if err != nil {
		return xerrors.Errorf("rotateCommitteeIfNeeded: %w", err)
	}
	if rec != nil {
		// create if record is present
		if err = c.createNewCommitteeAndConsensus(rec); err != nil {
			return xerrors.Errorf("rotateCommitteeIfNeeded: creating committee and consensus: %v", err)
		}
	}
	return nil
}

func (c *chainObj) getOwnCommitteeRecord(addr ledgerstate.Address) (*registry.CommitteeRecord, error) {
	rec, err := c.committeeRegistry.GetCommitteeRecord(addr)
	if err != nil {
		return nil, xerrors.Errorf("createCommitteeIfNeeded: reading committee record: %v", err)
	}
	if rec == nil {
		// committee record wasn't found in th registry, I am not the part of the committee
		return nil, nil
	}
	// just in case check if I am among committee nodes
	// should not happen
	if !util.StringInList(c.peerNetworkConfig.OwnNetID(), rec.Nodes) {
		return nil, xerrors.Errorf("createCommitteeIfNeeded: I am not among nodes of the committee record. Inconsistency")
	}
	return rec, nil
}

func (c *chainObj) createNewCommitteeAndConsensus(cmtRec *registry.CommitteeRecord) error {
	c.log.Debugf("createNewCommitteeAndConsensus: creating a new committee...")
	cmt, err := committee.New(
		cmtRec,
		c.chainID,
		c.netProvider,
		c.peerNetworkConfig,
		c.dksProvider,
		c.log,
	)
	if err != nil {
		c.setCommittee(nil)
		return xerrors.Errorf("createNewCommitteeAndConsensus: failed to create committee object for state address %s: %w",
			cmtRec.Address.Base58(), err)
	}
	cmt.RegisterPeerMessageParty(c)
	c.log.Debugf("creating new consensus object...")
	c.consensus = consensus.New(c, c.mempool, cmt, c.nodeConn, c.pullMissingRequestsFromCommittee, c.chainMetrics)
	c.setCommittee(cmt)

	c.log.Infof("NEW COMMITTEE OF VALIDATORS has been initialized for the state address %s", cmtRec.Address.Base58())
	return nil
}

func (c *chainObj) getCommittee() chain.Committee {
	ret := c.committee.Load().(*committeeStruct)
	if !ret.valid {
		return nil
	}
	return ret.cmt
}

func (c *chainObj) setCommittee(cmt chain.Committee) {
	if cmt == nil {
		c.committee.Store(&committeeStruct{})
	} else {
		c.committee.Store(&committeeStruct{
			valid: true,
			cmt:   cmt,
		})
	}
}
