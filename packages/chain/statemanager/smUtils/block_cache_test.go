// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package smUtils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/testutil/testlogger"
)

func TestBlockCacheSimple(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	_, blocks, _ := GetBlocks(t, 4, 1)
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

	_, blocks, _ := GetBlocks(t, 6, 2)
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

	_, blocks, _ := GetBlocks(t, 6, 2)
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
