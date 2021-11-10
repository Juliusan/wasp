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
	"io"

	"github.com/iotaledger/wasp/packages/util"
)

// Fills only data fields. Other fields must be filled by the caller.
//nolint:gocritic
func NewPeerMessageIn(buf []byte) (*PeerMessageIn, error) {
	var err error
	r := bytes.NewBuffer(buf)
	pmd := PeerMessageData{}
	if err = pmd.PeeringID.Read(r); err != nil {
		return nil, err
	}
	if err = util.ReadInt64(r, &pmd.Timestamp); err != nil {
		return nil, err
	}
	if pmd.DestinationParty, err = readPeerMessagePartyType(r); err != nil {
		return nil, err
	}
	if pmd.MsgDeserializer, err = readPeerMessagePartyType(r); err != nil {
		return nil, err
	}
	if pmd.MsgType, err = util.ReadByte(r); err != nil {
		return nil, err
	}
	if pmd.MsgData, err = util.ReadBytes32(r); err != nil {
		return nil, err
	}
	return &PeerMessageIn{PeerMessageData: pmd}, nil
}

func readPeerMessagePartyType(r io.Reader) (PeerMessagePartyType, error) {
	byteValue, err := util.ReadByte(r)
	if err != nil {
		return 0, err
	}
	return NewPeerMessagePartyType(byteValue)
}
