package smMessages

import (
	"fmt"

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
	if data[0] != MsgTypeBlockMessage {
		return nil, fmt.Errorf("Error creating block message from bytes: wrong message type %v", data[0])
	}
	block, err := state.BlockFromBytes(data[1:])
	if err != nil {
		return nil, fmt.Errorf("Error creating block message from bytes: %v", err)
	}
	return NewBlockMessage(block, "UNKNOWN"), nil
}

func (bmT *BlockMessage) MarshalBinary() (data []byte, err error) {
	return append([]byte{MsgTypeBlockMessage}, bmT.block.Bytes()...), nil
}

func (bmT *BlockMessage) GetBlock() state.Block {
	return bmT.block
}
