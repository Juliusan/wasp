package smGPA

import (
	"github.com/iotaledger/wasp/packages/state"
)

type createOriginStateFun func() (state.VirtualStateAccess, error)

type stateBlockRequest struct { // Abstract struct for requests to obtain certain virtual state;
	implementation       stateBlockRequestImplementation
	done                 bool
	blocks               []state.Block
	createOriginStateFun createOriginStateFun
}

type stateBlockRequestImplementation interface { // Abstract methods of struct stateBlockRequest
	isImplementationValid() bool
	respond(state.VirtualStateAccess)
}

var _ blockRequest = &stateBlockRequest{}

func newStateBlockRequest(sbri stateBlockRequestImplementation, createOriginStateFun createOriginStateFun) *stateBlockRequest {
	return &stateBlockRequest{
		implementation:       sbri,
		done:                 false,
		blocks:               make([]state.Block, 0),
		createOriginStateFun: createOriginStateFun,
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

func (sbrT *stateBlockRequest) markCompleted() {
	if sbrT.isValid() {
		vState, err := sbrT.createOriginStateFun()
		sbrT.done = true
		if err != nil {
			return
		}
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
