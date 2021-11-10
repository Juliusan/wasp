// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dkg

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
	"go.dedis.ch/kyber/v3"
	"golang.org/x/xerrors"
)

type initiatorDoneMsg struct {
	step      byte
	pubShares []kyber.Point
	blsSuite  kyber.Group // Transient, for un-marshaling only.
}

func (m *initiatorDoneMsg) GetMsgType() byte {
	return initiatorDoneMsgType
}

func (m *initiatorDoneMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyDkg
}

func (m *initiatorDoneMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteByte(&buf, m.step); err != nil {
		return nil, err
	}
	if err = util.WriteUint16(&buf, uint16(len(m.pubShares))); err != nil {
		return nil, err
	}
	for i := range m.pubShares {
		if err = util.WriteMarshaled(&buf, m.pubShares[i]); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (m *initiatorDoneMsg) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if m.step, err = util.ReadByte(r); err != nil {
		return err
	}
	var arrLen uint16
	if err = util.ReadUint16(r, &arrLen); err != nil {
		return err
	}
	m.pubShares = make([]kyber.Point, arrLen)
	for i := range m.pubShares {
		m.pubShares[i] = m.blsSuite.Point()
		if err = util.ReadMarshaled(r, m.pubShares[i]); err != nil {
			return xerrors.Errorf("failed to unmarshal initiatorDoneMsg.pubShares: %w", err)
		}
	}
	return nil
}

/*func (m *initiatorDoneMsg) Step() byte {
	return m.step
}*/

func (m *initiatorDoneMsg) SetStep(step byte) {
	m.step = step
}

/*func (m *initiatorDoneMsg) Error() error {
	return nil
}

func (m *initiatorDoneMsg) IsResponse() bool {
	return false
}*/
