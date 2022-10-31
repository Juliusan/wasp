package smGPA

import (
	"github.com/iotaledger/wasp/packages/state"
)

type handleVStateFun func(vState state.VirtualStateAccess)

type localStateBlockRequest struct {
	*stateBlockRequest
	lastBlockHash state.BlockHash
	respondFun    handleVStateFun
}

var (
	_ blockRequest                    = &localStateBlockRequest{} // Is a blockRequest
	_ stateBlockRequestImplementation = &localStateBlockRequest{} // Implements abstract methods of stateBlockRequest
)

func newLocalStateBlockRequest(bh state.BlockHash, respondFun handleVStateFun) blockRequest {
	result := &localStateBlockRequest{
		lastBlockHash: bh,
		respondFun:    respondFun,
	}
	result.stateBlockRequest = newStateBlockRequest(result)
	return result
}

func (lsbrT *localStateBlockRequest) getLastBlockHash() state.BlockHash {
	return lsbrT.lastBlockHash
}

func (lsbrT *localStateBlockRequest) isImplementationValid() bool {
	return true
}

func (lsbrT *localStateBlockRequest) respond(vState state.VirtualStateAccess) {
	lsbrT.respondFun(vState)
}
