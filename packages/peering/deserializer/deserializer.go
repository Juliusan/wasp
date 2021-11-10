// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package deserializer

import (
	"github.com/iotaledger/wasp/packages/chain/chainimpl"
	"github.com/iotaledger/wasp/packages/chain/committee"
	"github.com/iotaledger/wasp/packages/chain/consensus"
	"github.com/iotaledger/wasp/packages/chain/statemgr"
	"github.com/iotaledger/wasp/packages/dkg"
	"github.com/iotaledger/wasp/packages/peering"
	"golang.org/x/xerrors"
)

func GetEmptyMessage(deserializer peering.PeerMessagePartyType, msgType byte) (peering.Serializable, error) {
	switch deserializer {
	case peering.PeerMessagePartyChain:
		return chainimpl.GetEmptyMessage(msgType)
	case peering.PeerMessagePartyCommittee:
		return committee.GetEmptyMessage(msgType)
	case peering.PeerMessagePartyConsensus:
		return consensus.GetEmptyMessage(msgType)
	case peering.PeerMessagePartyDkg:
		return dkg.GetEmptyMessage(msgType)
	case peering.PeerMessagePartyStateManager:
		return statemgr.GetEmptyMessage(msgType)
	default:
		return nil, xerrors.Errorf("unknown deserializer: %v", deserializer)
	}

}
