package statemanager

import (
	"context"

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
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
)

type stateManager struct {
	log             *logger.Logger
	chainID         *isc.ChainID
	stateManagerGPA gpa.GPA
	nodeRandomiser  *smUtils.NodeRandomiser
	nodeIDToPubKey  map[gpa.NodeID]*cryptolib.PublicKey
}

var _ consGR.StateMgr = &stateManager{}
var _ node.ChainStateMgr = &stateManager{}

func New(
	chainID *isc.ChainID,
	me *cryptolib.PublicKey,
	peerPubKeys []*cryptolib.PublicKey,
	store kvstore.KVStore,
	log *logger.Logger,
) node.ChainStateMgr { // TODO: consGR.StateMgr {
	smLog := log.Named("sm")
	nr := smUtils.NewNodeRandomiserNoInit(pubKeyAsNodeID(me))
	result := &stateManager{
		log:             smLog,
		chainID:         chainID,
		stateManagerGPA: smGPA.New(chainID, nr, store, smLog),
		nodeRandomiser:  nr,
	}
	result.updatePublicKeys(peerPubKeys)
	return result
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
	// TODO
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

func pubKeyAsNodeID(pubKey *cryptolib.PublicKey) gpa.NodeID {
	return gpa.NodeID(pubKey.String())
}
