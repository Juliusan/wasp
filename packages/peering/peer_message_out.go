// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// Package peering provides an overlay network for communicating
// between nodes in a peer-to-peer style with low overhead
// encoding and persistent connections. The network provides only
// the asynchronous communication.
//
// It is intended to use for the committee consensus protocol.
//
package peering

import (
	"bytes"
	"time"

	"github.com/iotaledger/wasp/packages/util"
)

func NewPeerMessageOut(peeringID PeeringID, destinationParty PeerMessagePartyType, msg Serializable) (*PeerMessageOut, error) {
	msgData, err := msg.Serialize()
	if err != nil {
		return nil, err
	}
	return &PeerMessageOut{
		PeerMessageData: PeerMessageData{
			PeeringID:        peeringID,
			Timestamp:        time.Now().UnixNano(),
			DestinationParty: destinationParty,
			MsgDeserializer:  msg.GetDeserializer(),
			MsgType:          msg.GetMsgType(),
			MsgData:          msgData,
		},
		serialized: nil,
	}, nil
}

func (pmoT *PeerMessageOut) Bytes() ([]byte, error) {
	if pmoT.serialized == nil {
		serialized, err := pmoT.bytes()
		if err != nil {
			return nil, err
		}
		pmoT.serialized = &serialized
	}
	return *(pmoT.serialized), nil
}

func (pmoT *PeerMessageOut) bytes() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = pmoT.PeeringID.Write(&buf); err != nil {
		return nil, err
	}
	if err = util.WriteInt64(&buf, pmoT.Timestamp); err != nil {
		return nil, err
	}
	if err = util.WriteByte(&buf, byte(pmoT.DestinationParty)); err != nil {
		return nil, err
	}
	if err = util.WriteByte(&buf, byte(pmoT.MsgDeserializer)); err != nil {
		return nil, err
	}
	if err = util.WriteByte(&buf, pmoT.MsgType); err != nil {
		return nil, err
	}
	if err = util.WriteBytes32(&buf, pmoT.MsgData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TODO: what????
func (m *PeerMessageOut) IsUserMessage() bool {
	return m.MsgType >= FirstUserMsgCode
}
