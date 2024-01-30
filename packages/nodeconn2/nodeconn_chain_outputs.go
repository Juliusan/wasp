package nodeconn2

import (
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/isc"
)

type chainNodeConnOutputsReceivedData struct {
	slot      iotago.SlotIndex
	anchor    *isc.AnchorOutputWithID
	account   *isc.AccountOutputWithID
	basics    []*isc.BasicOutputWithID
	foundries []*isc.FoundryOutputWithID
	nfts      []*isc.NFTOutputWithID
}

func newChainNodeConnOutputsReceivedData(
	slot iotago.SlotIndex,
	anchor *isc.AnchorOutputWithID,
	account *isc.AccountOutputWithID,
	basics []*isc.BasicOutputWithID,
	foundries []*isc.FoundryOutputWithID,
	nfts []*isc.NFTOutputWithID,
) *chainNodeConnOutputsReceivedData {
	return &chainNodeConnOutputsReceivedData{
		slot:      slot,
		anchor:    anchor,
		account:   account,
		basics:    basics,
		foundries: foundries,
		nfts:      nfts,
	}
}

func (cncord *chainNodeConnOutputsReceivedData) slotIndex() iotago.SlotIndex {
	return cncord.slot
}

func (cncord *chainNodeConnOutputsReceivedData) anchorOutput() *isc.AnchorOutputWithID {
	return cncord.anchor
}

func (cncord *chainNodeConnOutputsReceivedData) accountOutput() *isc.AccountOutputWithID {
	return cncord.account
}

func (cncord *chainNodeConnOutputsReceivedData) basicOutputs() []*isc.BasicOutputWithID {
	return cncord.basics
}

func (cncord *chainNodeConnOutputsReceivedData) foundryOutputs() []*isc.FoundryOutputWithID {
	return cncord.foundries
}

func (cncord *chainNodeConnOutputsReceivedData) nftOutputs() []*isc.NFTOutputWithID {
	return cncord.nfts
}
