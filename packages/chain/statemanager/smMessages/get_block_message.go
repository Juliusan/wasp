package smMessages

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
)

type GetBlockMessage struct {
	*GpaMessage
	blockHash state.BlockHash
}

var _ gpa.Message = &GetBlockMessage{}

func NewGetBlockMessage(blockHash state.BlockHash, to gpa.NodeID) *GetBlockMessage {
	return &GetBlockMessage{
		GpaMessage: newGPAMessage(to),
		blockHash:  blockHash,
	}
}

func NewGetBlockMessageFromBytes(data []byte) (*GetBlockMessage, error) {
	if len(data) != state.BlockHashSize {
		return nil, fmt.Errorf("Wrong size of get block message: %v, expecting %v", len(data), state.BlockHashSize)
	}
	var blockHash state.BlockHash
	copy(blockHash[:], data)
	return NewGetBlockMessage(blockHash, "UNKNOWN"), nil
}

func (gbmT *GetBlockMessage) MarshalBinary() (data []byte, err error) {
	return gbmT.blockHash[:], nil
}

func (gbmT *GetBlockMessage) GetBlockHash() state.BlockHash {
	return gbmT.blockHash
}
