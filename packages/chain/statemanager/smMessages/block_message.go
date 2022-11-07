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
	result := NewBlockMessage(nil, "UNKNOWN") // NOTE: `block` will be set in `UnmarshalBinary` method
	err := result.UnmarshalBinary(data)
	return result, err
}

func (bmT *BlockMessage) MarshalBinary() (data []byte, err error) {
	return append([]byte{MsgTypeBlockMessage}, bmT.block.Bytes()...), nil
}

func (bmT *BlockMessage) UnmarshalBinary(data []byte) error {
	if data[0] != MsgTypeBlockMessage {
		return fmt.Errorf("Error creating block message from bytes: wrong message type %v", data[0])
	}
	var err error
	bmT.block, err = state.BlockFromBytes(data[1:])
	return err
}

func (bmT *BlockMessage) GetBlock() state.Block {
	return bmT.block
}
