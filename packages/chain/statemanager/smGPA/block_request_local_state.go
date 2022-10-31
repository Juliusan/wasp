package smGPA

import (
	"github.com/iotaledger/wasp/packages/state"
)

// TODO: blocks parameter should probably be removed after DB refactoring
type handleVStateFun func(blocks []state.Block, vState state.VirtualStateAccess)

type localStateBlockRequest struct {
	*stateBlockRequest
	lastBlockHash  state.BlockHash
	lastBlockIndex uint32 // TODO: temporar field. Remove it after DB refactoring.
	priority       uint32
	respondFun     handleVStateFun
}

var (
	_ blockRequest                    = &localStateBlockRequest{} // Is a blockRequest
	_ stateBlockRequestImplementation = &localStateBlockRequest{} // Implements abstract methods of stateBlockRequest
)

// TODO: `bi` is a temporar parameter. Remove it after DB refactoring.
func newLocalStateBlockRequest(bi uint32, bh state.BlockHash, priority uint32, respondFun handleVStateFun) blockRequest {
	result := &localStateBlockRequest{
		lastBlockHash:  bh,
		lastBlockIndex: bi,
		priority:       priority,
		respondFun:     respondFun,
	}
	result.stateBlockRequest = newStateBlockRequest(result)
	return result
}

func (lsbrT *localStateBlockRequest) getLastBlockHash() state.BlockHash {
	return lsbrT.lastBlockHash
}

func (lsbrT *localStateBlockRequest) getLastBlockIndex() uint32 { // TODO: temporar function. Remove it after DB refactoring.
	return lsbrT.lastBlockIndex
}

func (lsbrT *localStateBlockRequest) isImplementationValid() bool {
	return true
}

func (lsbrT *localStateBlockRequest) getPriority() uint32 {
	return lsbrT.priority
}

// TODO: blocks parameter should probably be removed after DB refactoring
func (lsbrT *localStateBlockRequest) respond(blocks []state.Block, vState state.VirtualStateAccess) {
	lsbrT.respondFun(blocks, vState)
}
