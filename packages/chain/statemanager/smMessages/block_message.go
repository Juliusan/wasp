package smMessages

import (
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
)

type BlockMessage struct {
	*GpaMessage
	block state.Block
}

var _ gpa.Message = &BlockMessage{}

func NewBlockMessage(block state.Block, to gpa.NodeID) *BlockMessage {
	return &BlockMessage{
		GpaMessage: newGPAMessage(to),
		block:      block,
	}
}

func NewBlockMessageFromBytes(data []byte) (*BlockMessage, error) {
	block, err := state.BlockFromBytes(data)
	if err != nil {
		return nil, err
	}
	return NewBlockMessage(block, "UNKNOWN"), nil
}

func (bmT *BlockMessage) MarshalBinary() (data []byte, err error) {
	return bmT.block.Bytes(), nil
}

func (bmT *BlockMessage) GetBlock() state.Block {
	return bmT.block
}
