// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

package coreroot

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"

type DeployContractCall struct {
	Func    *wasmlib.ScFunc
	Params  MutableDeployContractParams
}

type GrantDeployPermissionCall struct {
	Func    *wasmlib.ScFunc
	Params  MutableGrantDeployPermissionParams
}

type RequireDeployPermissionsCall struct {
	Func    *wasmlib.ScFunc
	Params  MutableRequireDeployPermissionsParams
}

type RevokeDeployPermissionCall struct {
	Func    *wasmlib.ScFunc
	Params  MutableRevokeDeployPermissionParams
}

type SubscribeBlockContextCall struct {
	Func    *wasmlib.ScFunc
	Params  MutableSubscribeBlockContextParams
}

type FindContractCall struct {
	Func    *wasmlib.ScView
	Params  MutableFindContractParams
	Results ImmutableFindContractResults
}

type GetContractRecordsCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetContractRecordsResults
}

type Funcs struct{}

var ScFuncs Funcs

func (sc Funcs) DeployContract(ctx wasmlib.ScFuncCallContext) *DeployContractCall {
	f := &DeployContractCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncDeployContract)}
	f.Params.proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

func (sc Funcs) GrantDeployPermission(ctx wasmlib.ScFuncCallContext) *GrantDeployPermissionCall {
	f := &GrantDeployPermissionCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncGrantDeployPermission)}
	f.Params.proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

func (sc Funcs) RequireDeployPermissions(ctx wasmlib.ScFuncCallContext) *RequireDeployPermissionsCall {
	f := &RequireDeployPermissionsCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncRequireDeployPermissions)}
	f.Params.proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

func (sc Funcs) RevokeDeployPermission(ctx wasmlib.ScFuncCallContext) *RevokeDeployPermissionCall {
	f := &RevokeDeployPermissionCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncRevokeDeployPermission)}
	f.Params.proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

func (sc Funcs) SubscribeBlockContext(ctx wasmlib.ScFuncCallContext) *SubscribeBlockContextCall {
	f := &SubscribeBlockContextCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncSubscribeBlockContext)}
	f.Params.proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

func (sc Funcs) FindContract(ctx wasmlib.ScViewCallContext) *FindContractCall {
	f := &FindContractCall{Func: wasmlib.NewScView(ctx, HScName, HViewFindContract)}
	f.Params.proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.proxy)
	return f
}

func (sc Funcs) GetContractRecords(ctx wasmlib.ScViewCallContext) *GetContractRecordsCall {
	f := &GetContractRecordsCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetContractRecords)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.proxy)
	return f
}

var exportMap = wasmlib.ScExportMap{
	Names: []string{
		FuncDeployContract,
		FuncGrantDeployPermission,
		FuncRequireDeployPermissions,
		FuncRevokeDeployPermission,
		FuncSubscribeBlockContext,
		ViewFindContract,
		ViewGetContractRecords,
	},
	Funcs: []wasmlib.ScFuncContextFunction{
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
	},
	Views: []wasmlib.ScViewContextFunction{
		wasmlib.ViewError,
		wasmlib.ViewError,
	},
}

func OnLoad(index int32) {
	if index >= 0 {
		panic("Calling core contract?")
	}

	exportMap.Export()
}
