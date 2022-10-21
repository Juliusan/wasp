package smGPA

import (
	"github.com/iotaledger/wasp/packages/state"
)

type blockRequest interface {
	getLastBlockHash() state.BlockHash
	blockAvailable(state.Block)
	isValid() bool
	markCompleted()
}
