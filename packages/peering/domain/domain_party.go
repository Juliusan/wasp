// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0
package domain

import (
	"github.com/iotaledger/wasp/packages/peering"
)

type domainParty struct {
	party  peering.PeerMessageSimpleParty
	domain *domainImpl
}

var _ peering.PeerMessageSimpleParty = &domainParty{}

func NewPeerMessageDomainSimpleParty(p peering.PeerMessageSimpleParty, d *domainImpl) peering.PeerMessageSimpleParty {
	return &domainParty{
		party:  p,
		domain: d,
	}
}

func (dpT *domainParty) GetPartyType() peering.PeerMessagePartyType {
	return dpT.party.GetPartyType()
}

func (dpT *domainParty) HandleMessage(senderNetID string, msg peering.Serializable) {
	_, ok := dpT.domain.nodes[senderNetID]
	if !ok {
		dpT.domain.log.Warnf("dropping message %T from %v, because it does not belong to the peer domain.", msg, senderNetID)
		return
	}
	// TODO: is it ok to assume that senderNetID == dpT.domain.nodes[senderNetID].NetID() ???
	if senderNetID == dpT.domain.netProvider.Self().NetID() {
		dpT.domain.log.Warnf("dropping message %T from self (%v).", msg, senderNetID)
		return
	}
	dpT.HandleMessage(senderNetID, msg)
}
