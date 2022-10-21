// TODO: WAL

package smUtils

import (
	"time"

	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/state"
)

type blockTime struct {
	time      time.Time
	blockHash state.BlockHash
}

type BlockCache struct {
	log          *logger.Logger
	blocks       map[state.BlockHash]state.Block
	times        []*blockTime
	timeProvider TimeProvider
}

func NewBlockCache(tp TimeProvider, log *logger.Logger) *BlockCache {
	result := &BlockCache{
		log:          log.Named("bc"),
		blocks:       make(map[state.BlockHash]state.Block),
		times:        make([]*blockTime, 0),
		timeProvider: tp,
	}
	return result
}

func (bcT *BlockCache) AddBlock(block state.Block) error {
	blockHash := block.GetHash()
	//TODO: store to DB
	bcT.addBlockToCache(blockHash, block)
	return nil
}

func (bcT *BlockCache) addBlockToCache(blockHash state.BlockHash, block state.Block) {
	bcT.blocks[blockHash] = block
	bcT.times = append(bcT.times, &blockTime{
		time:      bcT.timeProvider.GetNow(),
		blockHash: blockHash,
	})
	bcT.log.Debugf("Block %v added to DB and cache", blockHash)
}

func (bcT *BlockCache) GetBlock(blockHash state.BlockHash) state.Block {
	block, ok := bcT.blocks[blockHash]

	if ok {
		return block
	}
	//TODO: search the DB

	// if rasta –> bcT.addBlockToCache(blockHash, block)
	return nil
}

func (bcT *BlockCache) CleanOlderThan(limit time.Time) {
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
