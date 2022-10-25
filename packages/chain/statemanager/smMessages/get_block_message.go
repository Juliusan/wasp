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
	if data[0] != MsgTypeGetBlockMessage {
		return nil, fmt.Errorf("Error creating get block message from bytes: wrong message type %v", data[0])
	}
	actualData := data[1:]
	if len(actualData) != state.BlockHashSize {
		return nil, fmt.Errorf("Error creating get block message from bytes: wrong size %v, expecting %v", len(actualData), state.BlockHashSize)
	}
	var blockHash state.BlockHash
	copy(blockHash[:], actualData)
	return NewGetBlockMessage(blockHash, "UNKNOWN"), nil
}

func (gbmT *GetBlockMessage) MarshalBinary() (data []byte, err error) {
	return append([]byte{MsgTypeGetBlockMessage}, gbmT.blockHash[:]...), nil
}

func (gbmT *GetBlockMessage) GetBlockHash() state.BlockHash {
	return gbmT.blockHash
}
