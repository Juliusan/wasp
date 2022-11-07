// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package smUtils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/core/kvstore/mapdb"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/tpkg"
	"github.com/iotaledger/trie.go/trie"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/state"
)

func MakeNodeID(index int) gpa.NodeID {
	return gpa.NodeID(fmt.Sprintf("Node%v", index))
}

func MakeNodeIDs(indexes []int) []gpa.NodeID {
	result := make([]gpa.NodeID, len(indexes))
	for i := range indexes {
		result[i] = MakeNodeID(indexes[i])
	}
	return result
}

func GetOriginState(t require.TestingT) (*isc.ChainID, *isc.AliasOutputWithID, state.VirtualStateAccess) {
	store := mapdb.NewMapDB()
	aliasOutput0ID := iotago.OutputIDFromTransactionIDAndIndex(getRandomTxID(), 0).UTXOInput()
	chainID := isc.ChainIDFromAliasID(iotago.AliasIDFromOutputID(aliasOutput0ID.ID()))
	originVS, err := state.CreateOriginState(store, &chainID)
	require.NoError(t, err)
	stateAddress := cryptolib.NewKeyPair().GetPublicKey().AsEd25519Address()
	aliasOutput0 := &iotago.AliasOutput{
		Amount:        tpkg.TestTokenSupply,
		AliasID:       *chainID.AsAliasID(), // NOTE: not very correct: origin output's AliasID should be empty; left here to make mocking transitions easier
		StateMetadata: state.OriginL1Commitment().Bytes(),
		Conditions: iotago.UnlockConditions{
			&iotago.StateControllerAddressUnlockCondition{Address: stateAddress},
			&iotago.GovernorAddressUnlockCondition{Address: stateAddress},
		},
		Features: iotago.Features{
			&iotago.SenderFeature{
				Address: stateAddress,
			},
		},
	}
	return &chainID, isc.NewAliasOutputWithID(aliasOutput0, aliasOutput0ID), originVS
}

func GetBlocks(t require.TestingT, count, branchingFactor int) (*isc.ChainID, []state.Block, []*isc.AliasOutputWithID) {
	chainID, aliasOutput0, originVS := GetOriginState(t)
	result := make([]state.Block, count+1)
	vStates := make([]state.VirtualStateAccess, len(result))
	vStates[0] = originVS
	aliasOutputs := make([]*isc.AliasOutputWithID, len(result))
	aliasOutputs[0] = aliasOutput0
	for i := 1; i < len(result); i++ {
		baseIndex := (i + branchingFactor - 2) / branchingFactor
		increment := uint64(1 + i%branchingFactor)
		result[i], aliasOutputs[i], vStates[i] = GetNextState(t, vStates[baseIndex], aliasOutputs[baseIndex], increment)
	}
	return chainID, result[1:], aliasOutputs[1:]
}

func GetNextState(
	t require.TestingT,
	vs state.VirtualStateAccess,
	consumedAliasOutputWithID *isc.AliasOutputWithID,
	incrementOpt ...uint64,
) (state.Block, *isc.AliasOutputWithID, state.VirtualStateAccess) {
	nextVS := vs.Copy()
	prevBlockIndex := vs.BlockIndex()
	counterKey := kv.Key(coreutil.StateVarBlockIndex + "counter")
	counterBin, err := nextVS.KVStore().Get(counterKey)
	require.NoError(t, err)
	counter, err := codec.DecodeUint64(counterBin, 0)
	require.NoError(t, err)
	var increment uint64
	if len(incrementOpt) > 0 {
		increment = incrementOpt[0]
	} else {
		increment = 1
	}

	consumedAliasOutput := consumedAliasOutputWithID.GetAliasOutput()
	vsCommitment, err := state.L1CommitmentFromBytes(consumedAliasOutput.StateMetadata)
	require.NoError(t, err)
	suBlockIndex := state.NewStateUpdateWithBlockLogValues(prevBlockIndex+1, time.Now(), &vsCommitment)

	suCounter := state.NewStateUpdate()
	counterBin = codec.EncodeUint64(counter + increment)
	suCounter.Mutations().Set(counterKey, counterBin)

	nextVS.ApplyStateUpdate(suBlockIndex)
	nextVS.ApplyStateUpdate(suCounter)
	nextVS.Commit()
	require.EqualValues(t, prevBlockIndex+1, nextVS.BlockIndex())

	block, err := nextVS.ExtractBlock()
	require.NoError(t, err)

	aliasOutput := &iotago.AliasOutput{
		Amount:         consumedAliasOutput.Amount,
		NativeTokens:   consumedAliasOutput.NativeTokens,
		AliasID:        consumedAliasOutput.AliasID,
		StateIndex:     consumedAliasOutput.StateIndex + 1,
		StateMetadata:  state.NewL1Commitment(trie.RootCommitment(nextVS.TrieNodeStore()), block.GetHash()).Bytes(),
		FoundryCounter: consumedAliasOutput.FoundryCounter,
		Conditions:     consumedAliasOutput.Conditions,
		Features:       consumedAliasOutput.Features,
	}
	aliasOutputID := iotago.OutputIDFromTransactionIDAndIndex(getRandomTxID(), 0).UTXOInput()
	aliasOutputWithID := isc.NewAliasOutputWithID(aliasOutput, aliasOutputID)

	return block, aliasOutputWithID, nextVS
}

func getRandomTxID() iotago.TransactionID {
	var result iotago.TransactionID
	rand.Read(result[:])
	return result
}
