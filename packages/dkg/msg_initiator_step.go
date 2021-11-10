// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dkg

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

//
// This is a message used to synchronize the DKG procedure by
// ensuring the lock-step, as required by the DKG algorithm
// assumptions (Rabin as well as Pedersen).
//
type initiatorStepMsg struct {
	step byte
}

func (m *initiatorStepMsg) GetMsgType() byte {
	return initiatorStepMsgType
}

func (m *initiatorStepMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyDkg
}

func (m *initiatorStepMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteByte(&buf, m.step); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *initiatorStepMsg) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if m.step, err = util.ReadByte(r); err != nil {
		return err
	}
	return nil
}

/*func (m *initiatorInitMsg) Step() byte {
	return m.step
}*/

func (m *initiatorStepMsg) SetStep(step byte) {
	m.step = step
}

/*func (m *initiatorStepMsg) Error() error {
	return nil
}

func (m *initiatorStepMsg) IsResponse() bool {
	return false
}*/
