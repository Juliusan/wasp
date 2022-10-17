// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package smUtils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/core/kvstore/mapdb"
	"github.com/iotaledger/trie.go/trie"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
)

func TestBlockCacheSimple(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	blocks := getBlocks(t, 4)
	blockCache := NewBlockCache(log)
	blockCache.AddBlock(blocks[0])
	blockCache.AddBlock(blocks[1])
	blockCache.AddBlock(blocks[2])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.Nil(t, blockCache.GetBlock(blocks[3].GetHash()))
}

func TestBlockCacheCleaning(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	blocks := getBlocks(t, 6)
	blockCache := NewBlockCache(log)
	beforeTime := time.Now()
	blockCache.AddBlock(blocks[0])
	blockCache.AddBlock(blocks[1])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	blockCache.cleanOlderThan(beforeTime)
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	blockCache.cleanOlderThan(time.Now())
	require.Nil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.Nil(t, blockCache.GetBlock(blocks[1].GetHash()))
	blockCache.AddBlock(blocks[2])
	blockCache.AddBlock(blocks[3])
	inTheMiddleTime := time.Now()
	blockCache.AddBlock(blocks[4])
	blockCache.AddBlock(blocks[5])
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[5].GetHash()))
	blockCache.cleanOlderThan(inTheMiddleTime)
	require.Nil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.Nil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[5].GetHash()))
}

func TestBlockCacheCleaningAuto(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	blocks := getBlocks(t, 6)
	blockCacheTimers := NewBlockCacheTimers()
	blockCacheTimers.BlocksInCacheDuration = 300 * time.Millisecond
	blockCacheTimers.BlockCleaningPeriod = 90 * time.Millisecond
	blockCache := NewBlockCache(log, blockCacheTimers)
	blockCache.AddBlock(blocks[0])
	blockCache.AddBlock(blocks[1])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	time.Sleep(100 * time.Millisecond)
	blockCache.AddBlock(blocks[2])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	time.Sleep(100 * time.Millisecond)
	blockCache.AddBlock(blocks[3])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	time.Sleep(170 * time.Millisecond)
	blockCache.AddBlock(blocks[4])
	require.Nil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.Nil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	time.Sleep(100 * time.Millisecond)
	require.Nil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	time.Sleep(100 * time.Millisecond)
	require.Nil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	time.Sleep(200 * time.Millisecond)
	require.Nil(t, blockCache.GetBlock(blocks[4].GetHash()))
}

func getBlocks(t *testing.T, count int) []state.Block {
	store := mapdb.NewMapDB()
	chainID := isc.RandomChainID()
	originVS, err := state.CreateOriginState(store, chainID)
	require.NoError(t, err)
	result := make([]state.Block, count+1)
	vStates := make([]state.VirtualStateAccess, len(result))
	vsCommitments := make([]*state.L1Commitment, len(result))
	vStates[0] = originVS
	vsCommitments[0] = state.OriginL1Commitment()
	for i := 1; i < len(result); i++ {
		baseIndex := i / 2
		increment := uint64(1 + i%2)
		result[i], vStates[i], vsCommitments[i] = nextState(t, vStates[baseIndex], vsCommitments[baseIndex], increment)
	}
	return result[1:]
}

func nextState(
	t *testing.T,
	vs state.VirtualStateAccess,
	vsCommitment *state.L1Commitment,
	incrementOpt ...uint64,
) (state.Block, state.VirtualStateAccess, *state.L1Commitment) {
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

	suBlockIndex := state.NewStateUpdateWithBlockLogValues(prevBlockIndex+1, time.Now(), vsCommitment)

	suCounter := state.NewStateUpdate()
	counterBin = codec.EncodeUint64(counter + increment)
	suCounter.Mutations().Set(counterKey, counterBin)

	nextVS.ApplyStateUpdate(suBlockIndex)
	nextVS.ApplyStateUpdate(suCounter)
	nextVS.Commit()
	require.EqualValues(t, prevBlockIndex+1, nextVS.BlockIndex())

	block, err := nextVS.ExtractBlock()
	require.NoError(t, err)
	nextVSCommitment := state.NewL1Commitment(trie.RootCommitment(nextVS.TrieNodeStore()), block.GetHash())

	return block, nextVS, nextVSCommitment
}
