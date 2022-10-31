package smGPA

import (
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
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

func (cdsbrT *consensusDecidedStateBlockRequest) isImplementationValid() bool {
	return cdsbrT.consensusDecidedState.IsValid()
}

func (cdsbrT *consensusDecidedStateBlockRequest) getPriority() uint32 {
	return topPriority
}

func (cdsbrT *consensusDecidedStateBlockRequest) respond(vState state.VirtualStateAccess) {
	cdsbrT.consensusDecidedState.Respond(vState)
}
