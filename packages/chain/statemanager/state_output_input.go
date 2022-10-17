package statemanager

import (
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
)

type StateOutputInput struct {
	stateOutput *isc.AliasOutputWithID
}

var _ gpa.Input = &StateOutputInput{}

func NewStateOutputInput(output *isc.AliasOutputWithID) *StateOutputInput {
	return &StateOutputInput{stateOutput: output}
}

func (soiT *StateOutputInput) GetStateOutput() *isc.AliasOutputWithID {
	return soiT.stateOutput
}
