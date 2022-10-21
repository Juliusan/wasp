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
	blockCache := NewBlockCache(NewDefaultTimeProvider(), log)
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
	blockCache := NewBlockCache(NewDefaultTimeProvider(), log)
	beforeTime := time.Now()
	blockCache.AddBlock(blocks[0])
	blockCache.AddBlock(blocks[1])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	blockCache.CleanOlderThan(beforeTime)
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	blockCache.CleanOlderThan(time.Now())
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
	blockCache.CleanOlderThan(inTheMiddleTime)
	require.Nil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.Nil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[5].GetHash()))
}
