// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dkg

import (
	"bytes"
	"time"

	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

//
// This is a message sent by the initiator to all the peers to
// initiate the DKG process.
//
type initiatorInitMsg struct {
	step         byte
	dkgRef       string // Some unique string to identify duplicate initialization.
	peerNetIDs   []string
	peerPubs     []ed25519.PublicKey
	initiatorPub ed25519.PublicKey
	threshold    uint16
	timeout      time.Duration
	roundRetry   time.Duration
}

func (m *initiatorInitMsg) GetMsgType() byte {
	return initiatorInitMsgType
}

func (m *initiatorInitMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyDkg
}

func (m *initiatorInitMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteByte(&buf, m.step); err != nil {
		return nil, err
	}
	if err = util.WriteString16(&buf, m.dkgRef); err != nil {
		return nil, err
	}
	if err = util.WriteStrings16(&buf, m.peerNetIDs); err != nil {
		return nil, err
	}
	if err = util.WriteUint16(&buf, uint16(len(m.peerPubs))); err != nil {
		return nil, err
	}
	for i := range m.peerPubs {
		if err = util.WriteBytes16(&buf, m.peerPubs[i].Bytes()); err != nil {
			return nil, err
		}
	}
	if err = util.WriteBytes16(&buf, m.initiatorPub.Bytes()); err != nil {
		return nil, err
	}
	if err = util.WriteUint16(&buf, m.threshold); err != nil {
		return nil, err
	}
	if err = util.WriteInt64(&buf, m.timeout.Milliseconds()); err != nil {
		return nil, err
	}
	if err = util.WriteInt64(&buf, m.roundRetry.Milliseconds()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *initiatorInitMsg) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if m.step, err = util.ReadByte(r); err != nil {
		return err
	}
	if m.dkgRef, err = util.ReadString16(r); err != nil {
		return err
	}
	if m.peerNetIDs, err = util.ReadStrings16(r); err != nil {
		return err
	}
	var arrLen uint16
	if err = util.ReadUint16(r, &arrLen); err != nil {
		return err
	}
	m.peerPubs = make([]ed25519.PublicKey, arrLen)
	for i := range m.peerPubs {
		var peerPubBytes []byte
		if peerPubBytes, err = util.ReadBytes16(r); err != nil {
			return err
		}
		if m.peerPubs[i], _, err = ed25519.PublicKeyFromBytes(peerPubBytes); err != nil {
			return err
		}
	}
	var initiatorPubBytes []byte
	if initiatorPubBytes, err = util.ReadBytes16(r); err != nil {
		return err
	}
	if m.initiatorPub, _, err = ed25519.PublicKeyFromBytes(initiatorPubBytes); err != nil {
		return err
	}
	if err = util.ReadUint16(r, &m.threshold); err != nil {
		return err
	}
	var timeoutMS int64
	if err = util.ReadInt64(r, &timeoutMS); err != nil {
		return err
	}
	m.timeout = time.Duration(timeoutMS) * time.Millisecond
	var roundRetryMS int64
	if err = util.ReadInt64(r, &roundRetryMS); err != nil {
		return err
	}
	m.roundRetry = time.Duration(roundRetryMS) * time.Millisecond
	return nil
}

/*func (m *initiatorInitMsg) Step() byte {
	return m.step
}*/

func (m *initiatorInitMsg) SetStep(step byte) {
	m.step = step
}

/*func (m *initiatorInitMsg) Error() error {
	return nil
}

func (m *initiatorInitMsg) IsResponse() bool {
	return false
}*/
