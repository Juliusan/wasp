// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package consensus

import (
	"github.com/iotaledger/hive.go/marshalutil"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/peering"
)

type MissingRequestIDsMsg struct {
	IDs []iscp.RequestID
}

type MissingRequestIDsMsgIn struct {
	MissingRequestIDsMsg
	SenderNetID string
}

var _ peering.Serializable = &MissingRequestIDsMsg{}

func (msg *MissingRequestIDsMsg) GetMsgType() byte {
	return byte(consensusMessageTypeMissingRequestIDs)
}

func (msg *MissingRequestIDsMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyConsensus
}

//TODO is marshal util needed?
func (msg *MissingRequestIDsMsg) Serialize() ([]byte, error) {
	mu := marshalutil.New()
	mu.WriteUint16(uint16(len(msg.IDs)))
	for i := range msg.IDs {
		mu.WriteBytes(msg.IDs[i].Bytes())
	}
	return mu.Bytes(), nil
}

//TODO is marshal util needed?
func (msg *MissingRequestIDsMsg) Deserialize(buf []byte) error {
	mu := marshalutil.New(buf)
	num, err := mu.ReadUint16()
	if err != nil {
		return err
	}
	msg.IDs = make([]iscp.RequestID, num)
	for i := range msg.IDs {
		if msg.IDs[i], err = iscp.RequestIDFromMarshalUtil(mu); err != nil {
			return err
		}
	}
	return nil
}
