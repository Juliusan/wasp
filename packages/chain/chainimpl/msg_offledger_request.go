// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainimpl

import (
	"github.com/iotaledger/hive.go/marshalutil"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/request"
	"github.com/iotaledger/wasp/packages/peering"
	"golang.org/x/xerrors"
)

type offLedgerRequestMsg struct {
	ChainID *iscp.ChainID
	Req     *request.OffLedger
}

type offLedgerRequestMsgIn struct {
	offLedgerRequestMsg
	SenderNetID string
}

var _ peering.Serializable = &offLedgerRequestMsg{}

func (msg *offLedgerRequestMsg) GetMsgType() byte {
	return byte(chainMessageTypeOffledgerRequest)
}

func (msg *offLedgerRequestMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyChain
}

func (msg *offLedgerRequestMsg) Serialize() ([]byte, error) {
	return marshalutil.New().
		Write(msg.ChainID).
		Write(msg.Req).
		Bytes(), nil
}

func (msg *offLedgerRequestMsg) Deserialize(buf []byte) error {
	var err error
	var ok bool
	mu := marshalutil.New(buf)
	if msg.ChainID, err = iscp.ChainIDFromMarshalUtil(mu); err != nil {
		return err
	}
	req, err := request.FromMarshalUtil(mu)
	if err != nil {
		return err
	}
	if msg.Req, ok = req.(*request.OffLedger); !ok {
		return xerrors.New("OffLedgerRequestMsgFromBytes: wrong type of request data")
	}
	return nil
}
