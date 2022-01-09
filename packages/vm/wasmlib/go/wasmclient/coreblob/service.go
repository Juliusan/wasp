// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

package coreblobclient

import "github.com/iotaledger/wasp/packages/vm/wasmlib/go/wasmclient"

const (
	ArgBlobs = "this"
	ArgField = "field"
	ArgHash  = "hash"

	ResBlobSizes = "this"
	ResBytes     = "bytes"
	ResHash      = "hash"
)

///////////////////////////// storeBlob /////////////////////////////

type StoreBlobFunc struct {
	wasmclient.ClientFunc
	args wasmclient.Arguments
}

func (f *StoreBlobFunc) Blobs(v []byte) {
	f.args.SetBytes(ArgBlobs, v)
}

func (f *StoreBlobFunc) Post() wasmclient.Request {
	return f.ClientFunc.Post(0xddd4c281, &f.args)
}

///////////////////////////// getBlobField /////////////////////////////

type GetBlobFieldView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetBlobFieldView) Field(v string) {
	f.args.SetString(ArgField, v)
}

func (f *GetBlobFieldView) Hash(v wasmclient.Hash) {
	f.args.SetHash(ArgHash, v)
}

func (f *GetBlobFieldView) Call() GetBlobFieldResults {
	f.args.Mandatory(ArgField)
	f.args.Mandatory(ArgHash)
	f.ClientView.Call("getBlobField", &f.args)
	return GetBlobFieldResults{res: f.Results()}
}

type GetBlobFieldResults struct {
	res wasmclient.Results
}

func (r *GetBlobFieldResults) Bytes() []byte {
	return r.res.GetBytes(ResBytes)
}

///////////////////////////// getBlobInfo /////////////////////////////

type GetBlobInfoView struct {
	wasmclient.ClientView
	args wasmclient.Arguments
}

func (f *GetBlobInfoView) Hash(v wasmclient.Hash) {
	f.args.SetHash(ArgHash, v)
}

func (f *GetBlobInfoView) Call() GetBlobInfoResults {
	f.args.Mandatory(ArgHash)
	f.ClientView.Call("getBlobInfo", &f.args)
	return GetBlobInfoResults{res: f.Results()}
}

type GetBlobInfoResults struct {
	res wasmclient.Results
}

func (r *GetBlobInfoResults) BlobSizes() map[string]int32 {
	res := make(map[string]int32)
	r.res.ForEach(func(key string, val string) {
		res[string(key)] = r.res.GetInt32(val)
	})
	return res
}

///////////////////////////// listBlobs /////////////////////////////

type ListBlobsView struct {
	wasmclient.ClientView
}

func (f *ListBlobsView) Call() ListBlobsResults {
	f.ClientView.Call("listBlobs", nil)
	return ListBlobsResults{res: f.Results()}
}

type ListBlobsResults struct {
	res wasmclient.Results
}

func (r *ListBlobsResults) BlobSizes() map[wasmclient.Hash]int32 {
	res := make(map[wasmclient.Hash]int32)
	r.res.ForEach(func(key string, val string) {
		res[wasmclient.Hash(key)] = r.res.GetInt32(val)
	})
	return res
}

///////////////////////////// CoreBlobService /////////////////////////////

type CoreBlobService struct {
	wasmclient.Service
}

func NewCoreBlobService(cl *wasmclient.ServiceClient, chainID string) (*CoreBlobService, error) {
	s := &CoreBlobService{}
	err := s.Service.Init(cl, chainID, 0xfd91bc63, nil)
	return s, err
}

func (s *CoreBlobService) StoreBlob() StoreBlobFunc {
	return StoreBlobFunc{ClientFunc: s.AsClientFunc()}
}

func (s *CoreBlobService) GetBlobField() GetBlobFieldView {
	return GetBlobFieldView{ClientView: s.AsClientView()}
}

func (s *CoreBlobService) GetBlobInfo() GetBlobInfoView {
	return GetBlobInfoView{ClientView: s.AsClientView()}
}

func (s *CoreBlobService) ListBlobs() ListBlobsView {
	return ListBlobsView{ClientView: s.AsClientView()}
}
