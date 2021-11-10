// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dkg

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

//
// initiatorStatusMsg
//
type initiatorStatusMsg struct {
	step  byte
	error error
}

func (m *initiatorStatusMsg) GetMsgType() byte {
	return initiatorStatusMsgType
}

func (m *initiatorStatusMsg) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyDkg
}

func (m *initiatorStatusMsg) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err := util.WriteByte(&buf, m.step); err != nil {
		return err
	}
	var errMsg string
	if m.error != nil {
		errMsg = m.error.Error()
	}
	if err := util.WriteString16(&buf, errMsg); err != nil {
		return err
	}
	return buf.Bytes(), nil
}

func (m *initiatorStatusMsg) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if m.step, err = util.ReadByte(r); err != nil {
		return err
	}
	var errMsg string
	if errMsg, err = util.ReadString16(r); err != nil {
		return err
	}
	if errMsg != "" {
		m.error = errors.New(errMsg)
	} else {
		m.error = nil
	}
	return nil
}

/*func (m *initiatorStatusMsg) Step() byte {
	return m.step
}*/

func (m *initiatorStatusMsg) SetStep(step byte) {
	m.step = step
}

/*func (m *initiatorStatusMsg) Error() error {
	return m.error
}

func (m *initiatorStatusMsg) IsResponse() bool {
	return true
}*/
