// TODO: WAL

package smUtils

import (
	"sync"
	"time"

	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/state"
)

type blockTime struct {
	time      time.Time
	blockHash state.BlockHash
}

type BlockCache struct {
	log    *logger.Logger
	blocks map[state.BlockHash]state.Block
	times  []*blockTime
	timers BlockCacheTimers
	mutex  sync.Mutex
}

func NewBlockCache(log *logger.Logger, timersOpt ...BlockCacheTimers) *BlockCache {
	var timers BlockCacheTimers
	if len(timersOpt) > 0 {
		timers = timersOpt[0]
	} else {
		timers = NewBlockCacheTimers()
	}
	result := &BlockCache{
		log:    log.Named("bc"),
		blocks: make(map[state.BlockHash]state.Block),
		times:  make([]*blockTime, 0),
		timers: timers,
	}
	result.startCleaningLoop()
	return result
}

func (bcT *BlockCache) AddBlock(block state.Block) {
	blockHash := block.GetHash()
	//TODO: store to DB
	bcT.addBlockToCache(blockHash, block)
}

func (bcT *BlockCache) addBlockToCache(blockHash state.BlockHash, block state.Block) {
	bcT.mutex.Lock()
	defer bcT.mutex.Unlock()

	bcT.blocks[blockHash] = block
	bcT.times = append(bcT.times, &blockTime{
		time:      time.Now(),
		blockHash: blockHash,
	})
	bcT.log.Debugf("Block %v added to DB and cache", blockHash)
}

func (bcT *BlockCache) GetBlock(blockHash state.BlockHash) state.Block {
	bcT.mutex.Lock()
	block, ok := bcT.blocks[blockHash]
	bcT.mutex.Unlock()

	if ok {
		return block
	}
	//TODO: search the DB

	bcT.mutex.Lock()
	defer bcT.mutex.Unlock()
	// if rasta –> bcT.addBlockToCache(blockHash, block)
	return nil
}

func (bcT *BlockCache) cleanOlderThan(limit time.Time) {
	bcT.mutex.Lock()
	defer bcT.mutex.Unlock()

	for i, bt := range bcT.times {
		if bt.time.After(limit) {
			bcT.times = bcT.times[i:]
			return
		}
		delete(bcT.blocks, bt.blockHash)
		bcT.log.Debugf("Block %v deleted from cache", bt.blockHash)
	}
	bcT.times = make([]*blockTime, 0) // All the blocks were deleted
}

func (bcT *BlockCache) startCleaningLoop() {
	go func() {
		for {
			bcT.cleanOlderThan(time.Now().Add(-bcT.timers.BlocksInCacheDuration))
			time.Sleep(bcT.timers.BlockCleaningPeriod)
		}
	}()
}
