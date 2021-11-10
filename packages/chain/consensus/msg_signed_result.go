// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package consensus

import (
	"bytes"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

type signedResultMsg struct {
	ChainInputID ledgerstate.OutputID
	EssenceHash  hashing.HashValue
	SigShare     tbls.SigShare
}

type signedResultMsgIn struct {
	signedResultMsg
	SenderIndex uint16
}

var _ peering.Serializable = &signedResultMsg{}

func (msg *signedResultMsg) GetMsgType() byte {
	return byte(consensusMessageTypeSignerResult)
}

func (msg *signedResultMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyConsensus
}

func (msg *signedResultMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteHashValue(&buf, msg.EssenceHash); err != nil {
		return nil, err
	}
	if err = util.WriteBytes16(&buf, msg.SigShare); err != nil {
		return nil, err
	}
	if err = util.WriteOutputID(&buf, msg.ChainInputID); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *signedResultMsg) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if err = util.ReadHashValue(r, &msg.EssenceHash); err != nil {
		return err
	}
	if msg.SigShare, err = util.ReadBytes16(r); err != nil {
		return err
	}
	if err := util.ReadOutputID(r, &msg.ChainInputID); /* nolint:revive */ err != nil {
		return err
	}
	return nil

}
