package smGPA

import (
	"github.com/iotaledger/wasp/packages/state"
)

type stateBlockRequest struct { // Abstract struct for requests to obtain certain virtual state;
	implementation stateBlockRequestImplementation
	done           bool
	blocks         []state.Block
}

type stateBlockRequestImplementation interface { // Abstract methods of struct stateBlockRequest
	isImplementationValid() bool
	respond(state.VirtualStateAccess)
}

var _ blockRequest = &stateBlockRequest{}

func newStateBlockRequest(sbri stateBlockRequestImplementation) *stateBlockRequest {
	return &stateBlockRequest{
		implementation: sbri,
		done:           false,
		blocks:         make([]state.Block, 0),
	}
}

func (sbrT *stateBlockRequest) getLastBlockHash() state.BlockHash {
	panic("Abstract method, should be overridden")
}

func (sbrT *stateBlockRequest) isValid() bool {
	if sbrT.done {
		return false
	}
	return sbrT.implementation.isImplementationValid()
}

func (sbrT *stateBlockRequest) blockAvailable(block state.Block) {
	sbrT.blocks = append(sbrT.blocks, block)
}

func (sbrT *stateBlockRequest) markCompleted(createBaseStateFun createStateFun) {
	if sbrT.isValid() {
		sbrT.done = true
		baseState, err := createBaseStateFun()
		if err != nil {
			// Something failed in creating the base state. Just forget the request.
			return
		}
		vState := baseState
		for i := len(sbrT.blocks) - 1; i >= 0; i-- {
			calculatedStateCommitment := state.RootCommitment(vState.TrieNodeStore())
			if !state.EqualCommitments(calculatedStateCommitment, sbrT.blocks[i].PreviousL1Commitment().StateCommitment) {
				return
			}
			err := vState.ApplyBlock(sbrT.blocks[i])
			if err != nil {
				return
			}
			vState.Commit() // TODO: is it needed
		}
		sbrT.implementation.respond(vState)
	}
}
