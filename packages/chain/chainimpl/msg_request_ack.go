// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainimpl

import (
	"bytes"
	"io/ioutil"

	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/peering"
	"golang.org/x/xerrors"
)

type requestAckMsg struct {
	ReqID *iscp.RequestID
}

type requestAckMsgIn struct {
	requestAckMsg
	SenderNetID string
}

var _ peering.Serializable = &requestAckMsg{}

func (msg *requestAckMsg) GetMsgType() byte {
	return byte(chainMessageTypeOffledgerRequest)
}

func (msg *requestAckMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyChain
}

func (msg *requestAckMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if _, err = buf.Write(msg.ReqID.Bytes()); err != nil {
		return nil, xerrors.Errorf("failed to write requestIDs: %w", err)
	}
	return buf.Bytes(), nil
}

func (msg *requestAckMsg) Deserialize(buf []byte) error {
	b, err := ioutil.ReadAll(bytes.NewReader(buf))
	if err != nil {
		return xerrors.Errorf("failed to read requestIDs: %w", err)
	}
	reqID, err := iscp.RequestIDFromBytes(b)
	if err != nil {
		return xerrors.Errorf("failed to read requestIDs: %w", err)
	}
	msg.ReqID = &reqID
	return nil
}
