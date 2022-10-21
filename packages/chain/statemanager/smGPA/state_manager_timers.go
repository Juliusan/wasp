// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package smGPA

import (
	"time"

	"github.com/iotaledger/wasp/packages/chain/statemanager/smUtils"
)

type StateManagerTimers struct {
	// How long should the block stay in block cache before being deleted
	BlockCacheBlocksInCacheDuration time.Duration
	// How often should the block cache be cleaned
	BlockCacheBlockCleaningPeriod time.Duration
	// How often get block requests should be repeated
	StateManagerGetBlockRetry time.Duration

	TimeProvider smUtils.TimeProvider
}

func NewStateManagerTimers(tpOpt ...smUtils.TimeProvider) StateManagerTimers {
	var tp smUtils.TimeProvider
	if len(tpOpt) > 0 {
		tp = tpOpt[0]
	} else {
		tp = smUtils.NewDefaultTimeProvider()
	}
	return StateManagerTimers{
		BlockCacheBlocksInCacheDuration: 1 * time.Hour,
		BlockCacheBlockCleaningPeriod:   1 * time.Minute,
		StateManagerGetBlockRetry:       3 * time.Second,
		TimeProvider:                    tp,
	}
}
