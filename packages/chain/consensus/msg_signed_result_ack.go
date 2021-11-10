// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package consensus

import (
	"bytes"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

type signedResultAckMsg struct {
	ChainInputID ledgerstate.OutputID
	EssenceHash  hashing.HashValue
}

type signedResultAckMsgIn struct {
	signedResultAckMsg
	SenderIndex uint16
}

var _ peering.Serializable = &signedResultAckMsg{}

func (msg *signedResultAckMsg) GetMsgType() byte {
	return byte(consensusMessageTypeSignerResultAck)
}

func (msg *signedResultAckMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyConsensus
}

func (msg *signedResultAckMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteHashValue(&buf, msg.EssenceHash); err != nil {
		return nil, err
	}
	if err = util.WriteOutputID(&buf, msg.ChainInputID); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *signedResultAckMsg) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if err = util.ReadHashValue(r, &msg.EssenceHash); err != nil {
		return err
	}
	if err = util.ReadOutputID(r, &msg.ChainInputID); err != nil {
		return err
	}
	return nil

}
