// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0
package group

import (
	"github.com/iotaledger/wasp/packages/peering"
)

type groupParty struct {
	party peering.PeerMessageGroupParty
	group *groupImpl
}

var _ peering.PeerMessageSimpleParty = &groupParty{}

func NewPeerMessageGroupSimpleParty(p peering.PeerMessageGroupParty, g *groupImpl) peering.PeerMessageSimpleParty {
	return &groupParty{
		party: p,
		group: g,
	}
}

func (gpT *groupParty) GetPartyType() peering.PeerMessagePartyType {
	return gpT.party.GetPartyType()
}

func (gpT *groupParty) HandleMessage(senderNetID string, msg peering.Serializable) {
	senderIndex, err := gpT.group.PeerIndexByNetID(senderNetID)
	if err != nil || senderIndex == NotInGroup {
		gpT.group.log.Warnf("Dropping message %T from %v, because it does not belong to the group.", msg, senderNetID)
	}
	gpT.party.HandleGroupMessage(senderNetID, senderIndex, msg)
}
