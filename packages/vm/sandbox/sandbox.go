// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package sandbox

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/sandbox/sandbox_utils"
	"github.com/iotaledger/wasp/packages/vm/vmcontext"
)

type sandbox struct {
	vmctx *vmcontext.VMContext
}

func init() {
	vmcontext.NewSandbox = func(vmctx *vmcontext.VMContext) coretypes.Sandbox {
		return &sandbox{vmctx: vmctx}
	}
}

var _ coretypes.Sandbox = &sandbox{}

func (s *sandbox) Utils() coretypes.Utils {
	return sandbox_utils.NewUtils()
}

func (s *sandbox) ChainOwnerID() *coretypes.AgentID {
	return s.vmctx.ChainOwnerID()
}

func (s *sandbox) ContractCreator() *coretypes.AgentID {
	return s.vmctx.ContractCreator()
}

func (s *sandbox) ContractID() *coretypes.ContractID {
	return s.vmctx.CurrentContractID()
}

func (s *sandbox) GetTimestamp() int64 {
	return s.vmctx.Timestamp()
}

func (s *sandbox) Params() dict.Dict {
	return s.vmctx.Params()
}

func (s *sandbox) State() kv.KVStore {
	return s.vmctx.State()
}

func (s *sandbox) Caller() *coretypes.AgentID {
	return s.vmctx.Caller()
}

// DeployContract deploys contract by the binary hash
// and calls "init" endpoint (constructor) with provided parameters
func (s *sandbox) DeployContract(programHash hashing.HashValue, name string, description string, initParams dict.Dict) error {
	return s.vmctx.DeployContract(programHash, name, description, initParams)
}

// Call calls an entry point of contact, passes parameters and funds
func (s *sandbox) Call(contractHname coretypes.Hname, entryPoint coretypes.Hname, params dict.Dict, transfer *coretypes.ColoredBalances) (dict.Dict, error) {
	return s.vmctx.Call(contractHname, entryPoint, params, transfer)
}

func (s *sandbox) RequestID() ledgerstate.OutputID {
	return s.vmctx.RequestID()
}

// Minted
func (s *sandbox) Minted() (ledgerstate.Color, uint64) {
	panic("Minted: implement me")
}

func (s *sandbox) GetEntropy() hashing.HashValue {
	return s.vmctx.Entropy()
}

func (s *sandbox) TransferToAddress(targetAddr ledgerstate.Address, transfer *coretypes.ColoredBalances) bool {
	panic("TransferToAddress: deprecated")

	//return s.vmctx.TransferToAddress(targetAddr, transfer)
}

func (s *sandbox) PostRequest(par coretypes.PostRequestParams) bool {
	panic("PostRequest: deprecated")
	//return s.vmctx.PostRequest(par)
}

func (s *sandbox) Send(target ledgerstate.Address, tokens *coretypes.ColoredBalances, metadata *coretypes.SendMetadata, options ...coretypes.SendOptions) bool {
	panic("SendTransfer: not implemented")
}

func (s *sandbox) Log() coretypes.LogInterface {
	return s.vmctx
}

func (s *sandbox) Event(msg string) {
	s.Log().Infof("eventlog::%s -> '%s'", s.vmctx.CurrentContractHname(), msg)
	s.vmctx.StoreToEventLog(s.vmctx.CurrentContractHname(), []byte(msg))
	s.vmctx.EventPublisher().Publish(msg)
}

func (s *sandbox) IncomingTransfer() *coretypes.ColoredBalances {
	return s.vmctx.GetIncoming()
}

func (s *sandbox) Balance(col ledgerstate.Color) uint64 {
	return s.vmctx.GetBalance(col)
}

func (s *sandbox) Balances() *coretypes.ColoredBalances {
	return s.vmctx.GetMyBalances()
}
