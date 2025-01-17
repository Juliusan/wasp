// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

package coreblocklogclient

import "github.com/iotaledger/wasp/packages/vm/wasmlib/go/wasmclient"

const (
	ArgBlockIndex    = "n"
	ArgContractHname = "h"
	ArgFromBlock     = "f"
	ArgRequestID     = "u"
	ArgToBlock       = "t"

	ResBlockIndex             = "n"
	ResBlockInfo              = "i"
	ResEvent                  = "e"
	ResGoverningAddress       = "g"
	ResRequestID              = "u"
	ResRequestIndex           = "r"
	ResRequestProcessed       = "p"
	ResRequestRecord          = "d"
	ResStateControllerAddress = "s"
)

///////////////////////////// controlAddresses /////////////////////////////

type ControlAddressesView struct {
	wasmclient.ClientView
}

func (f *ControlAddressesView) Call() ControlAddressesResults {
	f.ClientView.Call("controlAddresses", nil)
	return ControlAddressesResults{res: f.Results()}
}

type ControlAddressesResults struct {
	res wasmclient.Results
}

func (r *ControlAddressesResults) BlockIndex() int32 {
	return r.res.ToInt32(r.res.Get(ResBlockIndex))
}

func (r *ControlAddressesResults) GoverningAddress() wasmclient.Address {
	return r.res.ToAddress(r.res.Get(ResGoverningAddress))
}

func (r *ControlAddressesResults) StateControllerAddress() wasmclient.Address {
	return r.res.ToAddress(r.res.Get(ResStateControllerAddress))
}

///////////////////////////// getBlockInfo /////////////////////////////

type GetBlockInfoView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetBlockInfoView) BlockIndex(v int32) {
	f.args.Set(ArgBlockIndex, f.args.FromInt32(v))
}

func (f *GetBlockInfoView) Call() GetBlockInfoResults {
	f.args.Mandatory(ArgBlockIndex)
	f.ClientView.Call("getBlockInfo", &f.args)
	return GetBlockInfoResults{res: f.Results()}
}

type GetBlockInfoResults struct {
	res wasmclient.Results
}

func (r *GetBlockInfoResults) BlockInfo() []byte {
	return r.res.ToBytes(r.res.Get(ResBlockInfo))
}

///////////////////////////// getEventsForBlock /////////////////////////////

type GetEventsForBlockView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetEventsForBlockView) BlockIndex(v int32) {
	f.args.Set(ArgBlockIndex, f.args.FromInt32(v))
}

func (f *GetEventsForBlockView) Call() GetEventsForBlockResults {
	f.args.Mandatory(ArgBlockIndex)
	f.ClientView.Call("getEventsForBlock", &f.args)
	return GetEventsForBlockResults{res: f.Results()}
}

type GetEventsForBlockResults struct {
	res wasmclient.Results
}

func (r *GetEventsForBlockResults) Event() []byte {
	return r.res.ToBytes(r.res.Get(ResEvent))
}

///////////////////////////// getEventsForContract /////////////////////////////

type GetEventsForContractView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetEventsForContractView) ContractHname(v wasmclient.Hname) {
	f.args.Set(ArgContractHname, f.args.FromHname(v))
}

func (f *GetEventsForContractView) FromBlock(v int32) {
	f.args.Set(ArgFromBlock, f.args.FromInt32(v))
}

func (f *GetEventsForContractView) ToBlock(v int32) {
	f.args.Set(ArgToBlock, f.args.FromInt32(v))
}

func (f *GetEventsForContractView) Call() GetEventsForContractResults {
	f.args.Mandatory(ArgContractHname)
	f.ClientView.Call("getEventsForContract", &f.args)
	return GetEventsForContractResults{res: f.Results()}
}

type GetEventsForContractResults struct {
	res wasmclient.Results
}

func (r *GetEventsForContractResults) Event() []byte {
	return r.res.ToBytes(r.res.Get(ResEvent))
}

///////////////////////////// getEventsForRequest /////////////////////////////

type GetEventsForRequestView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetEventsForRequestView) RequestID(v wasmclient.RequestID) {
	f.args.Set(ArgRequestID, f.args.FromRequestID(v))
}

func (f *GetEventsForRequestView) Call() GetEventsForRequestResults {
	f.args.Mandatory(ArgRequestID)
	f.ClientView.Call("getEventsForRequest", &f.args)
	return GetEventsForRequestResults{res: f.Results()}
}

type GetEventsForRequestResults struct {
	res wasmclient.Results
}

func (r *GetEventsForRequestResults) Event() []byte {
	return r.res.ToBytes(r.res.Get(ResEvent))
}

///////////////////////////// getLatestBlockInfo /////////////////////////////

type GetLatestBlockInfoView struct {
	wasmclient.ClientView
}

func (f *GetLatestBlockInfoView) Call() GetLatestBlockInfoResults {
	f.ClientView.Call("getLatestBlockInfo", nil)
	return GetLatestBlockInfoResults{res: f.Results()}
}

type GetLatestBlockInfoResults struct {
	res wasmclient.Results
}

func (r *GetLatestBlockInfoResults) BlockIndex() int32 {
	return r.res.ToInt32(r.res.Get(ResBlockIndex))
}

func (r *GetLatestBlockInfoResults) BlockInfo() []byte {
	return r.res.ToBytes(r.res.Get(ResBlockInfo))
}

///////////////////////////// getRequestIDsForBlock /////////////////////////////

type GetRequestIDsForBlockView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetRequestIDsForBlockView) BlockIndex(v int32) {
	f.args.Set(ArgBlockIndex, f.args.FromInt32(v))
}

func (f *GetRequestIDsForBlockView) Call() GetRequestIDsForBlockResults {
	f.args.Mandatory(ArgBlockIndex)
	f.ClientView.Call("getRequestIDsForBlock", &f.args)
	return GetRequestIDsForBlockResults{res: f.Results()}
}

type GetRequestIDsForBlockResults struct {
	res wasmclient.Results
}

func (r *GetRequestIDsForBlockResults) RequestID() wasmclient.RequestID {
	return r.res.ToRequestID(r.res.Get(ResRequestID))
}

///////////////////////////// getRequestReceipt /////////////////////////////

type GetRequestReceiptView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetRequestReceiptView) RequestID(v wasmclient.RequestID) {
	f.args.Set(ArgRequestID, f.args.FromRequestID(v))
}

func (f *GetRequestReceiptView) Call() GetRequestReceiptResults {
	f.args.Mandatory(ArgRequestID)
	f.ClientView.Call("getRequestReceipt", &f.args)
	return GetRequestReceiptResults{res: f.Results()}
}

type GetRequestReceiptResults struct {
	res wasmclient.Results
}

func (r *GetRequestReceiptResults) BlockIndex() int32 {
	return r.res.ToInt32(r.res.Get(ResBlockIndex))
}

func (r *GetRequestReceiptResults) RequestIndex() int16 {
	return r.res.ToInt16(r.res.Get(ResRequestIndex))
}

func (r *GetRequestReceiptResults) RequestRecord() []byte {
	return r.res.ToBytes(r.res.Get(ResRequestRecord))
}

///////////////////////////// getRequestReceiptsForBlock /////////////////////////////

type GetRequestReceiptsForBlockView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetRequestReceiptsForBlockView) BlockIndex(v int32) {
	f.args.Set(ArgBlockIndex, f.args.FromInt32(v))
}

func (f *GetRequestReceiptsForBlockView) Call() GetRequestReceiptsForBlockResults {
	f.args.Mandatory(ArgBlockIndex)
	f.ClientView.Call("getRequestReceiptsForBlock", &f.args)
	return GetRequestReceiptsForBlockResults{res: f.Results()}
}

type GetRequestReceiptsForBlockResults struct {
	res wasmclient.Results
}

func (r *GetRequestReceiptsForBlockResults) RequestRecord() []byte {
	return r.res.ToBytes(r.res.Get(ResRequestRecord))
}

///////////////////////////// isRequestProcessed /////////////////////////////

type IsRequestProcessedView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *IsRequestProcessedView) RequestID(v wasmclient.RequestID) {
	f.args.Set(ArgRequestID, f.args.FromRequestID(v))
}

func (f *IsRequestProcessedView) Call() IsRequestProcessedResults {
	f.args.Mandatory(ArgRequestID)
	f.ClientView.Call("isRequestProcessed", &f.args)
	return IsRequestProcessedResults{res: f.Results()}
}

type IsRequestProcessedResults struct {
	res wasmclient.Results
}

func (r *IsRequestProcessedResults) RequestProcessed() string {
	return r.res.ToString(r.res.Get(ResRequestProcessed))
}

///////////////////////////// CoreBlockLogService /////////////////////////////

type CoreBlockLogService struct {
	wasmclient.Service
}

func NewCoreBlockLogService(cl *wasmclient.ServiceClient, chainID string) (*CoreBlockLogService, error) {
	s := &CoreBlockLogService{}
	err := s.Service.Init(cl, chainID, 0xf538ef2b, nil)
	return s, err
}

func (s *CoreBlockLogService) ControlAddresses() ControlAddressesView {
	return ControlAddressesView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetBlockInfo() GetBlockInfoView {
	return GetBlockInfoView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetEventsForBlock() GetEventsForBlockView {
	return GetEventsForBlockView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetEventsForContract() GetEventsForContractView {
	return GetEventsForContractView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetEventsForRequest() GetEventsForRequestView {
	return GetEventsForRequestView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetLatestBlockInfo() GetLatestBlockInfoView {
	return GetLatestBlockInfoView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetRequestIDsForBlock() GetRequestIDsForBlockView {
	return GetRequestIDsForBlockView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetRequestReceipt() GetRequestReceiptView {
	return GetRequestReceiptView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) GetRequestReceiptsForBlock() GetRequestReceiptsForBlockView {
	return GetRequestReceiptsForBlockView{ClientView: s.AsClientView()}
}

func (s *CoreBlockLogService) IsRequestProcessed() IsRequestProcessedView {
	return IsRequestProcessedView{ClientView: s.AsClientView()}
}
