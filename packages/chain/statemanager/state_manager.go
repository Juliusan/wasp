package statemanager

import (
	"context"
	"time"

	"github.com/iotaledger/hive.go/core/kvstore"
	"github.com/iotaledger/hive.go/core/logger"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/chain/aaa2/cons/gr"
	"github.com/iotaledger/wasp/packages/chain/aaa2/node"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smGPA"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smUtils"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util/pipe"
)

type stateManager struct {
	log             *logger.Logger
	chainID         *isc.ChainID
	stateManagerGPA gpa.GPA
	nodeRandomiser  *smUtils.NodeRandomiser
	nodeIDToPubKey  map[gpa.NodeID]*cryptolib.PublicKey
	inputPipe       pipe.Pipe
	messagePipe     pipe.Pipe
	net             peering.NetworkProvider
	netPeeringID    peering.PeeringID
	ctx             context.Context
	cleanupFun      func()
}

var (
	_ consGR.StateMgr    = &stateManager{}
	_ node.ChainStateMgr = &stateManager{}
)

const (
	constMsgTypeStm    byte = iota
	constTimerTickTime      = 1 * time.Second
)

func New(
	ctx context.Context,
	chainID *isc.ChainID,
	me *cryptolib.PublicKey,
	peerPubKeys []*cryptolib.PublicKey,
	net peering.NetworkProvider,
	walFolder string,
	store kvstore.KVStore,
	log *logger.Logger,
) (node.ChainStateMgr, error) { // TODO: consGR.StateMgr {
	smLog := log.Named("sm")
	nr := smUtils.NewNodeRandomiserNoInit(pubKeyAsNodeID(me))
	stateManagerGPA, err := smGPA.New(chainID, nr, walFolder, store, smLog)
	if err != nil {
		smLog.Errorf("Failed to create state manager GPA: %v", err)
		return nil, err
	}
	result := &stateManager{
		log:             smLog,
		chainID:         chainID,
		stateManagerGPA: stateManagerGPA,
		nodeRandomiser:  nr,
		inputPipe:       pipe.NewDefaultInfinitePipe(),
		messagePipe:     pipe.NewDefaultInfinitePipe(),
		net:             net,
		netPeeringID:    peering.PeeringIDFromBytes(hashing.HashDataBlake2b(chainID.Bytes(), []byte("STM")).Bytes()),
		ctx:             ctx,
	}
	result.updatePublicKeys(peerPubKeys)

	attachID := result.net.Attach(&result.netPeeringID, peering.PeerMessageReceiverStateManager, func(recv *peering.PeerMessageIn) {
		if recv.MsgType != constMsgTypeStm {
			result.log.Warnf("Unexpected message, type=%v", recv.MsgType)
			return
		}
		result.messagePipe.In() <- recv
	})

	result.cleanupFun = func() {
		result.inputPipe.Close()
		result.messagePipe.Close()
		result.net.Detach(attachID)
	}

	go result.run()
	return result, nil
}

// -------------------------------------
// Implementations of node.ChainStateMgr
// -------------------------------------

func (smT *stateManager) BlockProduced(ctx context.Context, aliasOutput *isc.AliasOutputWithID, block state.Block) <-chan error {
	input, resultCh := smInputs.NewChainBlockProduced(ctx, aliasOutput, block)
	smT.addInput(input)
	return resultCh
}

func (smT *stateManager) ReceiveConfirmedAliasOutput(aliasOutput *isc.AliasOutputWithID) {
	smT.addInput(smInputs.NewChainReceiveConfirmedAliasOutput(aliasOutput))
}

func (smT *stateManager) AccessNodesUpdated(accessNodePubKeys []*cryptolib.PublicKey) {
	// TODO: do it in one state manager thread
	smT.updatePublicKeys(accessNodePubKeys)
}

// -------------------------------------
// Implementations of consGR.StateMgr (TODO: or whatever - for consensus)
// -------------------------------------

// ConsensusStateProposal asks State manager to ensure that all the blocks for aliasOutput are available.
// `true` is sent via the returned channel uppon successful retrieval of every block for aliasOutput.
func (smT *stateManager) ConsensusStateProposal(ctx context.Context, aliasOutput *isc.AliasOutputWithID) <-chan interface{} {
	input, resultCh := smInputs.NewConsensusStateProposal(ctx, aliasOutput)
	smT.addInput(input)
	return resultCh
}

// ConsensusDecidedState asks State manager to return a virtual state vith stateCommitment as its state commitment
func (smT *stateManager) ConsensusDecidedState(ctx context.Context, aliasOutputID iotago.OutputID, stateCommitment *state.L1Commitment) <-chan *consGR.StateMgrDecidedState {
	input, resultCh := smInputs.NewConsensusDecidedState(ctx, aliasOutputID, stateCommitment)
	smT.addInput(input)
	return resultCh
}

// -------------------------------------
// Internal functions
// -------------------------------------

func (smT *stateManager) addInput(input gpa.Input) {
	smT.inputPipe.In() <- input
}

func (smT *stateManager) updatePublicKeys(peerPubKeys []*cryptolib.PublicKey) {
	smT.nodeIDToPubKey = make(map[gpa.NodeID]*cryptolib.PublicKey)
	peerNodeIDs := make([]gpa.NodeID, len(peerPubKeys))
	for i := range peerPubKeys {
		peerNodeIDs[i] = pubKeyAsNodeID(peerPubKeys[i])
		smT.nodeIDToPubKey[peerNodeIDs[i]] = peerPubKeys[i]
	}
	smT.nodeRandomiser.UpdateNodeIDs(peerNodeIDs)
}

func (smT *stateManager) run() {
	defer smT.cleanupFun()
	ctxCloseCh := smT.ctx.Done()
	inputPipeCh := smT.inputPipe.Out()
	messagePipeCh := smT.messagePipe.Out()
	timerTickCh := time.After(constTimerTickTime)
	for {
		select {
		case input, ok := <-inputPipeCh:
			if ok {
				smT.handleInput(input.(gpa.Input))
			} else {
				inputPipeCh = nil
			}
		case msg, ok := <-messagePipeCh:
			if ok {
				smT.handleMessage(msg.(*peering.PeerMessageIn))
			} else {
				messagePipeCh = nil
			}
		case now, ok := <-timerTickCh:
			if ok {
				smT.handleTimerTick(now)
				timerTickCh = time.After(constTimerTickTime)
			} else {
				timerTickCh = nil
			}
		case <-ctxCloseCh:
			smT.log.Debugf("Stopping state manager, because context was closed")
			return
		}
	}
}

func (smT *stateManager) handleInput(input gpa.Input) {
	outMsgs := smT.stateManagerGPA.Input(input)
	smT.sendMessages(outMsgs)
}

func (smT *stateManager) handleMessage(peerMsg *peering.PeerMessageIn) {
	msg, err := smT.stateManagerGPA.UnmarshalMessage(peerMsg.MsgData)
	if err != nil {
		smT.log.Warnf("Parsing message failed: %v", err)
		return
	}
	msg.SetSender(pubKeyAsNodeID(peerMsg.SenderPubKey))
	outMsgs := smT.stateManagerGPA.Message(msg)
	smT.sendMessages(outMsgs)
}

func (smT *stateManager) handleTimerTick(now time.Time) {
	smT.handleInput(smInputs.NewStateManagerTimerTick(now))
}

func (smT *stateManager) sendMessages(msgs gpa.OutMessages) {
	msgs.MustIterate(func(msg gpa.Message) {
		msgData, err := msg.MarshalBinary()
		if err != nil {
			smT.log.Warnf("Failed to marshal message for sending: %v", err)
			return
		}
		pm := &peering.PeerMessageData{
			PeeringID:   smT.netPeeringID,
			MsgReceiver: peering.PeerMessageReceiverStateManager,
			MsgType:     constMsgTypeStm,
			MsgData:     msgData,
		}
		smT.net.SendMsgByPubKey(smT.nodeIDToPubKey[msg.Recipient()], pm)
	})
}

func pubKeyAsNodeID(pubKey *cryptolib.PublicKey) gpa.NodeID {
	return gpa.NodeID(pubKey.String())
}
