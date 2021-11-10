// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainimpl

import (
	"github.com/iotaledger/hive.go/marshalutil"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/request"
	"github.com/iotaledger/wasp/packages/peering"
)

type missingRequestMsg struct {
	Request iscp.Request
}

var _ peering.Serializable = &missingRequestMsg{}

func (msg *missingRequestMsg) GetMsgType() byte {
	return byte(chainMessageTypeMissingRequest)
}

func (msg *missingRequestMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyChain
}

func (msg *missingRequestMsg) Serialize() ([]byte, error) {
	return msg.Request.Bytes(), nil
}

func (msg *missingRequestMsg) Deserialize(buf []byte) error {
	var err error
	msg.Request, err = request.FromMarshalUtil(marshalutil.New(buf))
	if err != nil {
		return err
	}
	return nil
}
