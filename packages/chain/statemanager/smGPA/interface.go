package smGPA

import (
	"github.com/iotaledger/wasp/packages/state"
)

type createStateFun func() (state.VirtualStateAccess, error)

type blockRequest interface {
	getLastBlockHash() state.BlockHash
	blockAvailable(state.Block)
	isValid() bool
	markCompleted(createStateFun) // NOTE: not all the requests need the base state, so a function to create one is passed rather than the created state
}
