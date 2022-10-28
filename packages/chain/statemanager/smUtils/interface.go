package smUtils

import (
	"github.com/iotaledger/wasp/packages/state"
)

type BlockWAL interface {
	Write(state.Block) error
	Contains(state.BlockHash) bool
	Read(state.BlockHash) (state.Block, error)
}
