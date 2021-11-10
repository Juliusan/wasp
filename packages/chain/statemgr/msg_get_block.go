// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package statemgr

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

// GetBlockMsg StateManager queries specific block data from another peer (access node)
type getBlockMsg struct {
	BlockIndex uint32
}

type getBlockMsgIn struct {
	getBlockMsg
	SenderNetID string
}

var _ peering.Serializable = &getBlockMsg{}

func (msg *getBlockMsg) GetMsgType() byte {
	return byte(stateManagerMessageTypeBlock)
}

func (msg *getBlockMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyStateManager
}

func (msg *getBlockMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteUint32(&buf, msg.BlockIndex); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *getBlockMsg) Deserialize(buf []byte) error {
	r := bytes.NewBuffer(buf)
	return util.ReadUint32(r, &msg.BlockIndex)
}
