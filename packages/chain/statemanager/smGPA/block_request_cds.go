package smGPA

import (
	"github.com/iotaledger/wasp/packages/chain/statemanager/smGPA/smInputs"
	"github.com/iotaledger/wasp/packages/state"
)

type consensusDecidedStateBlockRequest struct {
	*stateBlockRequest
	consensusDecidedState *smInputs.ConsensusDecidedState
}

var (
	_ blockRequest                    = &consensusDecidedStateBlockRequest{} // Is a blockRequest
	_ stateBlockRequestImplementation = &consensusDecidedStateBlockRequest{} // Implements abstract methods of stateBlockRequest
)

func newConsensusDecidedStateBlockRequest(input *smInputs.ConsensusDecidedState) blockRequest {
	result := &consensusDecidedStateBlockRequest{
		consensusDecidedState: input,
	}
	result.stateBlockRequest = newStateBlockRequest(result)
	return result
}

func (cdsbrT *consensusDecidedStateBlockRequest) getLastBlockHash() state.BlockHash {
	return cdsbrT.consensusDecidedState.GetStateCommitment().BlockHash
}

func (cdsbrT *consensusDecidedStateBlockRequest) getLastBlockIndex() uint32 { // TODO: temporar function.  Remove it after DB refactoring.
	return cdsbrT.consensusDecidedState.GetBlockIndex()
}

func (cdsbrT *consensusDecidedStateBlockRequest) isImplementationValid() bool {
	return cdsbrT.consensusDecidedState.IsValid()
}

func (cdsbrT *consensusDecidedStateBlockRequest) getPriority() uint32 {
	return topPriority
}

// TODO: first parameter should probably be removed after DB refactoring
func (cdsbrT *consensusDecidedStateBlockRequest) respond(_ []state.Block, vState state.VirtualStateAccess) {
	cdsbrT.consensusDecidedState.Respond(vState)
}
