// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package testchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/contracts/native/inccounter"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/testutil/utxodb"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

////////////////////////////////////////////////////////////////////////////////
// TestChainLedger

type TestChainLedger struct {
	t           *testing.T
	utxoDB      *utxodb.UtxoDB
	governor    *cryptolib.KeyPair
	chainID     isc.ChainID
	fetchedReqs map[iotago.Address]map[iotago.OutputID]bool
}

func NewTestChainLedger(t *testing.T, utxoDB *utxodb.UtxoDB, originator *cryptolib.KeyPair) *TestChainLedger {
	return &TestChainLedger{
		t:           t,
		utxoDB:      utxoDB,
		governor:    originator,
		fetchedReqs: map[iotago.Address]map[iotago.OutputID]bool{},
	}
}

// Only set after MakeTxChainOrigin.
func (tcl *TestChainLedger) ChainID() isc.ChainID {
	return tcl.chainID
}

func (tcl *TestChainLedger) MakeTxChainOrigin(committeeAddress iotago.Address) (*iotago.SignedTransaction, *isc.ChainOutputs, isc.ChainID) {
	outs := tcl.utxoDB.GetUnspentOutputs(tcl.governor.Address())
	originTX, _, chainID, err := origin.NewChainOriginTransaction(
		tcl.governor,
		committeeAddress,
		tcl.governor.Address(),
		100*isc.Million,
		0,
		nil,
		outs,
		testutil.L1API.TimeProvider().SlotFromTime(time.Now()),
		allmigrations.DefaultScheme.LatestSchemaVersion(),
		testutil.L1API,
	)
	require.NoError(tcl.t, err)
	stateAnchor, anchorOutput, err := transaction.GetAnchorFromTransaction(originTX.Transaction)
	require.NoError(tcl.t, err)
	require.NotNil(tcl.t, stateAnchor)
	require.NotNil(tcl.t, anchorOutput)
	chainOutputs := isc.NewChainOutputs(anchorOutput, stateAnchor.OutputID, nil, iotago.OutputID{})
	require.NoError(tcl.t, tcl.utxoDB.AddToLedger(originTX))
	tcl.chainID = chainID
	return originTX, chainOutputs, chainID
}

func (tcl *TestChainLedger) MakeTxAccountsDeposit(account *cryptolib.KeyPair) []isc.Request {
	outs := tcl.utxoDB.GetUnspentOutputs(account.Address())
	tx, err := transaction.NewRequestTransaction(
		account,
		account.Address(),
		outs,
		&isc.RequestParameters{
			TargetAddress:                 tcl.chainID.AsAddress(),
			Assets:                        isc.NewAssetsBaseTokens(100_000_000),
			AdjustToMinimumStorageDeposit: false,
			Metadata: &isc.SendMetadata{
				TargetContract: accounts.Contract.Hname(),
				EntryPoint:     accounts.FuncDeposit.Hname(),
				GasBudget:      2 * gas.LimitsDefault.MinGasPerRequest,
			},
		},
		nil,
		testutil.L1API.TimeProvider().SlotFromTime(time.Now()),
		false,
		testutil.L1API,
	)
	require.NoError(tcl.t, err)
	require.NoError(tcl.t, tcl.utxoDB.AddToLedger(tx))
	return tcl.findChainRequests(tx.Transaction)
}

func (tcl *TestChainLedger) MakeTxDeployIncCounterContract() []isc.Request {
	sender := tcl.governor
	outs := tcl.utxoDB.GetUnspentOutputs(sender.Address())
	tx, err := transaction.NewRequestTransaction(
		sender,
		sender.Address(),
		outs,
		&isc.RequestParameters{
			TargetAddress:                 tcl.chainID.AsAddress(),
			Assets:                        isc.NewAssetsBaseTokens(2_000_000),
			AdjustToMinimumStorageDeposit: false,
			Metadata: &isc.SendMetadata{
				TargetContract: root.Contract.Hname(),
				EntryPoint:     root.FuncDeployContract.Hname(),
				Params: codec.MakeDict(map[string]interface{}{
					root.ParamProgramHash: inccounter.Contract.ProgramHash,
					root.ParamName:        inccounter.Contract.Name,
					inccounter.VarCounter: 0,
				}),
				GasBudget: 2 * gas.LimitsDefault.MinGasPerRequest,
			},
		},
		nil,
		testutil.L1API.TimeProvider().SlotFromTime(time.Now()),
		false,
		testutil.L1API,
	)
	require.NoError(tcl.t, err)
	require.NoError(tcl.t, tcl.utxoDB.AddToLedger(tx))
	return tcl.findChainRequests(tx.Transaction)
}

func (tcl *TestChainLedger) FakeStateTransition(chainOuts *isc.ChainOutputs, stateCommitment *state.L1Commitment) *isc.ChainOutputs {
	stateMetadata := transaction.NewStateMetadata(
		stateCommitment,
		gas.DefaultFeePolicy(),
		0,
		"",
	)
	anchorOutput := &iotago.AnchorOutput{
		Amount:     chainOuts.AnchorOutput.Amount,
		AnchorID:   tcl.chainID.AsAnchorID(),
		StateIndex: chainOuts.GetStateIndex() + 1,
		UnlockConditions: iotago.AnchorOutputUnlockConditions{
			&iotago.StateControllerAddressUnlockCondition{Address: tcl.governor.Address()},
			&iotago.GovernorAddressUnlockCondition{Address: tcl.governor.Address()},
		},
		Features: iotago.AnchorOutputFeatures{
			&iotago.SenderFeature{
				Address: tcl.chainID.AsAddress(),
			},
			&iotago.StateMetadataFeature{
				Entries: map[iotago.StateMetadataFeatureEntriesKey]iotago.StateMetadataFeatureEntriesValue{
					"": stateMetadata.Bytes(),
				},
			},
		},
	}

	// TODO how to transition the account output?
	accountID, accountOut, _ := chainOuts.AccountOutput()

	return isc.NewChainOutputs(
		anchorOutput,
		iotago.OutputID{byte(anchorOutput.StateIndex)},
		accountOut,
		accountID,
	)
}

func (tcl *TestChainLedger) FakeRotationTX(chainOuts *isc.ChainOutputs, nextCommitteeAddr iotago.Address) (*isc.ChainOutputs, *iotago.SignedTransaction) {
	tx, err := transaction.NewRotateChainStateControllerTx(
		tcl.chainID.AsAnchorID(),
		nextCommitteeAddr,
		chainOuts.AnchorOutputID,
		chainOuts.AnchorOutput,
		testutil.L1API.TimeProvider().SlotFromTime(time.Now()),
		testutil.L1API,
		tcl.governor,
	)
	if err != nil {
		panic(err)
	}
	outputs, err := tx.Transaction.OutputsSet()
	if err != nil {
		panic(err)
	}
	for outputID, output := range outputs {
		if output.Type() == iotago.OutputAnchor {
			ao := output.(*iotago.AnchorOutput)
			// TODO I'm not sure if the state index should be updated here
			ao.StateIndex = chainOuts.GetStateIndex() + 1 // Fake next state index, just for tests.
			accountOutputID, accountOutput, _ := chainOuts.AccountOutput()
			return isc.NewChainOutputs(ao, outputID, accountOutput, accountOutputID), tx
		}
	}
	panic("anchor output not found")
}

func (tcl *TestChainLedger) findChainRequests(tx *iotago.Transaction) []isc.Request {
	reqs := []isc.Request{}
	outputs, err := tx.OutputsSet()
	require.NoError(tcl.t, err)
	for outputID, output := range outputs {
		// If that's anchor output of the chain, then it is not a request.
		if output.Type() == iotago.OutputAnchor {
			anchorOut := output.(*iotago.AnchorOutput)
			if anchorOut.AnchorID == tcl.chainID.AsAnchorID() {
				continue // That's our anchor output, not the request, skip it here.
			}
			if anchorOut.AnchorID.Empty() {
				implicitAnchorID := iotago.AnchorIDFromOutputID(outputID)
				if implicitAnchorID == tcl.chainID.AsAnchorID() {
					continue // That's our origin anchor output, not the request, skip it here.
				}
			}
		}
		//
		// Otherwise check the receiving address.
		outAddr := output.UnlockConditionSet().Address()
		if outAddr == nil {
			continue
		}
		if !outAddr.Address.Equal(tcl.chainID.AsAddress()) {
			continue
		}
		req, err := isc.OnLedgerFromUTXO(output, outputID)
		if err != nil {
			continue
		}
		reqs = append(reqs, req)
	}
	return reqs
}
