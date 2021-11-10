// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"fmt"
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/hive.go/events"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/chain/messages"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/coreutil"
	"github.com/iotaledger/wasp/packages/iscp/request"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/tcrypto"
	"github.com/iotaledger/wasp/packages/util/ready"
	"github.com/iotaledger/wasp/packages/vm/processors"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

type ChainCore interface {
	ID() *iscp.ChainID
	GetCommitteeInfo() *CommitteeInfo
	RegisterPeerMessageParty(party peering.PeerMessageSimpleParty) error
	UnregisterPeerMessageParty(partyType peering.PeerMessagePartyType) error
	SendMsgByNetID(netID string, destinationParty peering.PeerMessagePartyType, msg peering.Serializable)
	SendMsgToRandomPeers(upToNumPeers uint16, destinationParty peering.PeerMessagePartyType, msg peering.Serializable)
	StateCandidateToStateManager(state.VirtualStateAccess, ledgerstate.OutputID)
	Events() ChainEvents
	Processors() *processors.Cache
	GlobalStateSync() coreutil.ChainStateSync
	GetStateReader() state.OptimisticStateReader
	Log() *logger.Logger

	// Most of these methods are made publick for mocking in tests
	EnqueueDismissChain(reason string) // This one should really be public
	EnqueueLedgerState(chainOutput *ledgerstate.AliasOutput, timestamp time.Time)
	EnqueueOffLedgerRequestPeerMsg(req *request.OffLedger, senderNetID string)
}

// ChainEntry interface to access chain from the chain registry side
type ChainEntry interface {
	ReceiveTransaction(*ledgerstate.Transaction)
	ReceiveInclusionState(ledgerstate.TransactionID, ledgerstate.InclusionState)
	ReceiveState(stateOutput *ledgerstate.AliasOutput, timestamp time.Time)
	ReceiveOutput(output ledgerstate.Output)

	Dismiss(reason string)
	IsDismissed() bool
}

// ChainRequests is an interface to query status of the request
type ChainRequests interface {
	GetRequestProcessingStatus(id iscp.RequestID) RequestProcessingStatus
	EventRequestProcessed() *events.Event
}

type ChainEvents interface {
	RequestProcessed() *events.Event
	ChainTransition() *events.Event
}

type Chain interface {
	ChainCore
	ChainRequests
	ChainEntry
}

// Committee is ordered (indexed 0..size-1) list of peers which run the consensus
type Committee interface {
	Address() ledgerstate.Address
	Size() uint16
	Quorum() uint16
	OwnPeerIndex() uint16
	DKShare() *tcrypto.DKShare
	RegisterPeerMessageParty(party peering.PeerMessageGroupParty) error
	UnregisterPeerMessageParty(partyType peering.PeerMessagePartyType) error
	SendMsgByIndex(targetPeerIndex uint16, destinationParty peering.PeerMessagePartyType, msg peering.Serializable) error
	SendMsgToPeers(destinationParty peering.PeerMessagePartyType, msg peering.Serializable, except ...uint16)
	IsAlivePeer(peerIndex uint16) bool
	QuorumIsAlive(quorum ...uint16) bool
	PeerStatus() []*PeerStatus
	IsReady() bool
	Close()
	RunACSConsensus(value []byte, sessionID uint64, stateIndex uint32, callback func(sessionID uint64, acs [][]byte))
	GetOtherValidatorsPeerIDs() []string
	GetRandomValidators(upToN int) []string
}

type NodeConnection interface {
	PullBacklog(addr *ledgerstate.AliasAddress)
	PullState(addr *ledgerstate.AliasAddress)
	PullConfirmedTransaction(addr ledgerstate.Address, txid ledgerstate.TransactionID)
	PullTransactionInclusionState(addr ledgerstate.Address, txid ledgerstate.TransactionID)
	PullConfirmedOutput(addr ledgerstate.Address, outputID ledgerstate.OutputID)
	PostTransaction(tx *ledgerstate.Transaction)
}

type StateManager interface {
	Ready() *ready.Ready
	EnqueueGetBlock(index uint32, senderNetID string)
	EnqueueBlock(block []byte, senderNetID string)
	EventStateMsg(msg *messages.StateMsg)
	EventOutputMsg(msg ledgerstate.Output)
	EventStateCandidateMsg(state.VirtualStateAccess, ledgerstate.OutputID)
	EventTimerMsg(msg messages.TimerTick)
	GetStatusSnapshot() *SyncInfo
	Close()
}

type Consensus interface {
	EventStateTransitionMsg(state.VirtualStateAccess, *ledgerstate.AliasOutput, time.Time)
	EnqueueSignedResult(chainInputID ledgerstate.OutputID, essenceHash hashing.HashValue, sigShare tbls.SigShare, senderIndex uint16)
	EnqueueSignedResultAck(chainInputID ledgerstate.OutputID, essenceHash hashing.HashValue, senderIndex uint16)
	EventInclusionsStateMsg(ledgerstate.TransactionID, ledgerstate.InclusionState)
	EventAsynchronousCommonSubsetMsg(msg *messages.AsynchronousCommonSubsetMsg)
	EventVMResultMsg(msg *messages.VMResultMsg)
	EventTimerMsg(messages.TimerTick)
	IsReady() bool
	Close()
	GetStatusSnapshot() *ConsensusInfo
	ShouldReceiveMissingRequest(req iscp.Request) bool
}

type Mempool interface {
	ReceiveRequests(reqs ...iscp.Request)
	ReceiveRequest(req iscp.Request) bool
	RemoveRequests(reqs ...iscp.RequestID)
	ReadyNow(nowis ...time.Time) []iscp.Request
	ReadyFromIDs(nowis time.Time, reqIDs ...iscp.RequestID) ([]iscp.Request, []int, bool)
	HasRequest(id iscp.RequestID) bool
	GetRequest(id iscp.RequestID) iscp.Request
	Info() MempoolInfo
	WaitRequestInPool(reqid iscp.RequestID, timeout ...time.Duration) bool // for testing
	WaitInBufferEmpty(timeout ...time.Duration) bool                       // for testing
	Close()
}

type AsynchronousCommonSubsetRunner interface {
	RunACSConsensus(value []byte, sessionID uint64, stateIndex uint32, callback func(sessionID uint64, acs [][]byte))
	Close()
}

type MempoolInfo struct {
	TotalPool      int
	ReadyCounter   int
	InBufCounter   int
	OutBufCounter  int
	InPoolCounter  int
	OutPoolCounter int
}

type SyncInfo struct {
	Synced                bool
	SyncedBlockIndex      uint32
	SyncedStateHash       hashing.HashValue
	SyncedStateTimestamp  time.Time
	StateOutputBlockIndex uint32
	StateOutputID         ledgerstate.OutputID
	StateOutputHash       hashing.HashValue
	StateOutputTimestamp  time.Time
}

type ConsensusInfo struct {
	StateIndex uint32
	Mempool    MempoolInfo
	TimerTick  int
}

type ReadyListRecord struct {
	Request iscp.Request
	Seen    map[uint16]bool
}

type CommitteeInfo struct {
	Address       ledgerstate.Address
	Size          uint16
	Quorum        uint16
	QuorumIsAlive bool
	PeerStatus    []*PeerStatus
}

type PeerStatus struct {
	Index     int
	PeeringID string
	IsSelf    bool
	Connected bool
}

type ChainTransitionEventData struct {
	VirtualState    state.VirtualStateAccess
	ChainOutput     *ledgerstate.AliasOutput
	OutputTimestamp time.Time
}

func (p *PeerStatus) String() string {
	return fmt.Sprintf("%+v", *p)
}

type RequestProcessingStatus int

const (
	RequestProcessingStatusUnknown = RequestProcessingStatus(iota)
	RequestProcessingStatusBacklog
	RequestProcessingStatusCompleted
)

const (
	// TimerTickPeriod time tick for consensus and state manager objects
	TimerTickPeriod = 100 * time.Millisecond
)
