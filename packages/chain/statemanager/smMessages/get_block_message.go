package smMessages

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util"
)

type GetBlockMessage struct {
	*GpaMessage
	blockHash  state.BlockHash
	blockIndex uint32 // TODO: temporar field. Remove it after DB is refactored.
}

var _ gpa.Message = &GetBlockMessage{}

// TODO: `blockIndex` is a temporar parameter. Remove it after DB is refactored.
func NewGetBlockMessage(blockIndex uint32, blockHash state.BlockHash, to gpa.NodeID) *GetBlockMessage {
	return &GetBlockMessage{
		GpaMessage: newGPAMessage(to),
		blockHash:  blockHash,
		blockIndex: blockIndex,
	}
}

func NewGetBlockMessageFromBytes(data []byte) (*GetBlockMessage, error) {
	if data[0] != MsgTypeGetBlockMessage {
		return nil, fmt.Errorf("Error creating get block message from bytes: wrong message type %v", data[0])
	}
	// TODO: temporar code. Remove it after DB is refactored.
	if len(data) < 5 {
		return nil, fmt.Errorf("Error creating get block message from bytes: wrong size %v, expecting 5 or more", len(data))
	}
	blockIndex, err := util.Uint32From4Bytes(data[1:5])
	if err != nil {
		return nil, err
	}
	// End of temporar code
	actualData := data[5:] //data[1:]
	if len(actualData) != state.BlockHashSize {
		return nil, fmt.Errorf("Error creating get block message from bytes: wrong size %v, expecting %v", len(actualData), state.BlockHashSize)
	}
	var blockHash state.BlockHash
	copy(blockHash[:], actualData)
	return NewGetBlockMessage(blockIndex, blockHash, "UNKNOWN"), nil
}

func (gbmT *GetBlockMessage) MarshalBinary() (data []byte, err error) {
	result := append([]byte{MsgTypeGetBlockMessage}, util.Uint32To4Bytes(gbmT.blockIndex)...)
	return append(result, gbmT.blockHash[:]...), nil
}

func (gbmT *GetBlockMessage) GetBlockHash() state.BlockHash {
	return gbmT.blockHash
}

func (gbmT *GetBlockMessage) GetBlockIndex() uint32 { // TODO: temporar function. Remove it after DB is refactored.
	return gbmT.blockIndex
}