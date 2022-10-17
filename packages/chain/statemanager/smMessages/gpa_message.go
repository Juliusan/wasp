package smMessages

import (
	"github.com/iotaledger/wasp/packages/gpa"
)

// Partly implements gpa.Message
type GpaMessage struct {
	recipient gpa.NodeID
	sender    gpa.NodeID
}

func newGPAMessage(recipient gpa.NodeID) *GpaMessage {
	return &GpaMessage{recipient: recipient}
}

// Implementation for gpa.Message
func (bmT *GpaMessage) Recipient() gpa.NodeID {
	return bmT.recipient
}

// Implementation for gpa.Message
func (bmT *GpaMessage) SetSender(sender gpa.NodeID) {
	bmT.sender = sender
}

func (bmT *GpaMessage) GetSender() gpa.NodeID {
	return bmT.sender
}
