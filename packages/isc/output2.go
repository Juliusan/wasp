// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc

import (
	iotago "github.com/iotaledger/iota.go/v4"
)

type OutputWithID interface {
	OutputID() iotago.OutputID
	Output() iotago.Output
	OutputType() iotago.OutputType
	TransactionID() iotago.TransactionID
}

type outputWithIDImpl struct {
	outputID iotago.OutputID
	output   iotago.Output
}

func newOutputWithID(output iotago.Output, outputID iotago.OutputID) *outputWithIDImpl {
	return &outputWithIDImpl{
		outputID: outputID,
		output:   output,
	}
}

func (o *outputWithIDImpl) OutputID() iotago.OutputID {
	return o.outputID
}

func (o *outputWithIDImpl) Output() iotago.Output {
	return o.output
}

func (o *outputWithIDImpl) OutputType() iotago.OutputType {
	return o.Output().Type()
}

func (o *outputWithIDImpl) TransactionID() iotago.TransactionID {
	return o.OutputID().TransactionID()
}

type AnchorOutputWithID struct {
	*outputWithIDImpl
}

var _ OutputWithID = &AnchorOutputWithID{}

func NewAnchorOutputWithID(output *iotago.AnchorOutput, outputID iotago.OutputID) *AnchorOutputWithID {
	return &AnchorOutputWithID{
		outputWithIDImpl: newOutputWithID(output, outputID),
	}
}

func (o *AnchorOutputWithID) AnchorOutput() *iotago.AnchorOutput {
	return o.Output().(*iotago.AnchorOutput)
}

type AccountOutputWithID struct {
	*outputWithIDImpl
}

var _ OutputWithID = &AccountOutputWithID{}

func NewAccountOutputWithID(output *iotago.AccountOutput, outputID iotago.OutputID) *AccountOutputWithID {
	return &AccountOutputWithID{
		outputWithIDImpl: newOutputWithID(output, outputID),
	}
}

func (o *AccountOutputWithID) AccountOutput() *iotago.AccountOutput {
	return o.Output().(*iotago.AccountOutput)
}

type BasicOutputWithID struct {
	*outputWithIDImpl
}

var _ OutputWithID = &BasicOutputWithID{}

func NewBasicOutputWithID(output *iotago.BasicOutput, outputID iotago.OutputID) *BasicOutputWithID {
	return &BasicOutputWithID{
		outputWithIDImpl: newOutputWithID(output, outputID),
	}
}

func (o *BasicOutputWithID) BasicOutput() *iotago.BasicOutput {
	return o.Output().(*iotago.BasicOutput)
}

type FoundryOutputWithID struct {
	*outputWithIDImpl
}

var _ OutputWithID = &FoundryOutputWithID{}

func NewFoundryOutputWithID(output *iotago.FoundryOutput, outputID iotago.OutputID) *FoundryOutputWithID {
	return &FoundryOutputWithID{
		outputWithIDImpl: newOutputWithID(output, outputID),
	}
}

func (o *FoundryOutputWithID) FoundryOutput() *iotago.FoundryOutput {
	return o.Output().(*iotago.FoundryOutput)
}

type NFTOutputWithID struct {
	*outputWithIDImpl
}

var _ OutputWithID = &NFTOutputWithID{}

func NewNFTOutputWithID(output *iotago.NFTOutput, outputID iotago.OutputID) *NFTOutputWithID {
	return &NFTOutputWithID{
		outputWithIDImpl: newOutputWithID(output, outputID),
	}
}

func (o *NFTOutputWithID) NFTOutput() *iotago.NFTOutput {
	return o.Output().(*iotago.NFTOutput)
}

type AnchorAccountOutput struct {
	anchorOutput  *AnchorOutputWithID
	accountOutput *AccountOutputWithID
}

func NewAnchorAccountOutput(anchorOutput *AnchorOutputWithID, accountOutput *AccountOutputWithID) *AnchorAccountOutput {
	return &AnchorAccountOutput{
		anchorOutput:  anchorOutput,
		accountOutput: accountOutput,
	}
}

func NewAnchorAccountOutputAnchor(anchorOutput *AnchorOutputWithID) *AnchorAccountOutput {
	return NewAnchorAccountOutput(anchorOutput, nil)
}

func NewAnchorAccountOutputAccount(accountOutput *AccountOutputWithID) *AnchorAccountOutput {
	return NewAnchorAccountOutput(nil, accountOutput)
}

func (o *AnchorAccountOutput) AnchorOutputWithID() *AnchorOutputWithID {
	return o.anchorOutput
}

func (o *AnchorAccountOutput) AnchorOutput() *iotago.AnchorOutput {
	aoID := o.AnchorOutputWithID()
	if aoID == nil {
		return nil
	}
	return aoID.AnchorOutput()
}

func (o *AnchorAccountOutput) AnchorOutputID() iotago.OutputID {
	aoID := o.AnchorOutputWithID()
	if aoID == nil {
		return iotago.OutputID{}
	}
	return aoID.OutputID()
}

func (o *AnchorAccountOutput) AccountOutputWithID() *AccountOutputWithID {
	return o.accountOutput
}

func (o *AnchorAccountOutput) AccountOutput() *iotago.AccountOutput {
	aoID := o.AccountOutputWithID()
	if aoID == nil {
		return nil
	}
	return aoID.AccountOutput()
}

func (o *AnchorAccountOutput) AccountOutputID() iotago.OutputID {
	aoID := o.AccountOutputWithID()
	if aoID == nil {
		return iotago.OutputID{}
	}
	return aoID.OutputID()
}
