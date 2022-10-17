package smInputs

import (
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
)

type ChainBlockProduced struct {
	aliasOutput *isc.AliasOutputWithID
	block       state.Block
}

var _ gpa.Input = &ConsensusStateProposal{}

func NewChainBlockProduced(aliasOutput *isc.AliasOutputWithID, block state.Block) *ChainBlockProduced {
	return &ChainBlockProduced{
		aliasOutput: aliasOutput,
		block:       block,
	}
}

func (capbT *ChainBlockProduced) GetAliasOutputWithID() *isc.AliasOutputWithID {
	return capbT.aliasOutput
}

func (capbT *ChainBlockProduced) GetBlock() state.Block {
	return capbT.block
}
