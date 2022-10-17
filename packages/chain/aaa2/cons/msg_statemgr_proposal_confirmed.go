// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"github.com/iotaledger/wasp/packages/gpa"
)

type msgStateMgrProposalConfirmed struct {
	node gpa.NodeID
}

var _ gpa.Message = &msgStateMgrProposalConfirmed{}

func NewMsgStateMgrProposalConfirmed(node gpa.NodeID) gpa.Message {
	return &msgStateMgrProposalConfirmed{node: node}
}

func (m *msgStateMgrProposalConfirmed) Recipient() gpa.NodeID {
	return m.node
}

func (m *msgStateMgrProposalConfirmed) SetSender(sender gpa.NodeID) {
	if sender != m.node {
		panic("wrong sender/receiver for a local message")
	}
}

func (m *msgStateMgrProposalConfirmed) MarshalBinary() ([]byte, error) {
	panic("trying to marshal a local message")
}