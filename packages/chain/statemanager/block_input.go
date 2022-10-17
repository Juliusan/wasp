package statemanager

import (
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
)

type BlockInput struct {
	block state.Block
}

var _ gpa.Input = &BlockInput{}

func NewBlockInput(block state.Block) *BlockInput {
	return &BlockInput{block: block}
}

func (biT *BlockInput) GetBlock() state.Block {
	return biT.block
}
