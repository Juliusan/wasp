// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package statemgr

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

// BlockMsg StateManager in response to GetBlockMsg sends block data to the querying node's StateManager
type blockMsg struct {
	BlockBytes []byte
}

type blockMsgIn struct {
	blockMsg
	SenderNetID string
}

var _ peering.Serializable = &blockMsg{}

func (msg *blockMsg) GetMsgType() byte {
	return byte(stateManagerMessageTypeBlock)
}

func (msg *blockMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyStateManager
}

func (msg *blockMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteBytes32(&buf, msg.BlockBytes); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *blockMsg) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if msg.BlockBytes, err = util.ReadBytes32(r); err != nil {
		return err
	}
	return nil

}
