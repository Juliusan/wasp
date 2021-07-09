package chainimpl

import (
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/coretypes/coreutil"
	"github.com/iotaledger/wasp/packages/coretypes/request"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/processors"
	"github.com/iotaledger/wasp/packages/vm/runvm"
	"github.com/iotaledger/wasp/packages/vm/vmtypes"
	"github.com/iotaledger/wasp/packages/vm/wasmproc"
	"github.com/stretchr/testify/require"
)

type MockedVMRunner struct {
	env                *MockedEnv
	runner             vm.VMRunner
	proc               *processors.Cache
	validatorFeeTarget coretypes.AgentID
	globalSync         coreutil.ChainStateSync
}

func NewMockedVMRunner(env *MockedEnv) *MockedVMRunner {
	processorConfig := processors.NewConfig()
	err := processorConfig.RegisterVMType(vmtypes.WasmTime, func(binary []byte) (coretypes.VMProcessor, error) {
		return wasmproc.GetProcessor(binary, env.Log)
	})
	require.NoError(env.T, err)

	agentID := coretypes.NewAgentID(env.ChainID.AsAddress(), 0)

	return &MockedVMRunner{
		env:                env,
		runner:             runvm.NewVMRunner(),
		proc:               processors.MustNew(processorConfig),
		validatorFeeTarget: *agentID,
		globalSync:         coreutil.NewChainStateSync().SetSolidIndex(0),
	}
}

func (mvrT *MockedVMRunner) Run(vState state.VirtualState, targetAddr ledgerstate.Address, tx *ledgerstate.Transaction) state.VirtualState {
	requestsOnLedger, err := request.RequestsOnLedgerFromTransaction(tx, targetAddr)
	require.NoError(mvrT.env.T, err)
	requests := make([]coretypes.Request, len(requestsOnLedger))
	for i := range requestsOnLedger {
		//		requestsOnLedger[i].WithFallbackOptions(mvrT.env.StateOutput.GetAliasAddress(), time.Time{})
		requests[i] = requestsOnLedger[i]
		requests[i].SolidifyArgs(mvrT.env.Registry)
	}
	mvrT.env.Log.Infof("XXX REqUESTS %v", len(requests))

	//outs := mvrT.env.Ledger.GetAliasOutputs(mvrT.env.ChainID.AsAddress())
	//require.EqualValues(mvrT.env.T, 1, len(outs))

	task := &vm.VMTask{
		Processors:         mvrT.proc,
		ChainInput:         mvrT.env.StateOutput,
		Requests:           requests,
		Timestamp:          time.Now(),
		VirtualState:       vState.Clone(),
		Entropy:            hashing.RandomHash(nil),
		ValidatorFeeTarget: mvrT.validatorFeeTarget,
		SolidStateBaseline: mvrT.globalSync.GetSolidIndexBaseline(),
		Log:                mvrT.env.Log,
	}
	task.OnFinish = func(_ dict.Dict, err error, vmError error) {
		require.NoError(mvrT.env.T, vmError)
		mvrT.env.Log.Debugf("runVM OnFinish callback: responding by state index: %d state hash: %s",
			task.VirtualState.BlockIndex(), task.VirtualState.Hash())
	}
	mvrT.runner.Run(task)
	return task.VirtualState
}
