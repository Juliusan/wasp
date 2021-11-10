// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dkg

import (
	"github.com/iotaledger/wasp/packages/peering"
)

func (n *Node) GetPartyType() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyDkg
}

func (n *Node) HandleMessage(senderNetID string, msg peering.Serializable) {
	//TODO
}

func (n *Node) SendMessageToPeer(peer peering.PeerSender, peeringID peering.PeeringID, step byte, msg msgByteCoder) {
	msg.SetStep(step)
	peer.SendMsg(peeringID, peering.PeerMessagePartyDkg, msg)
}
