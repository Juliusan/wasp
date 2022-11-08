package smUtils

import (
	"testing"
	//	"time"

	"github.com/stretchr/testify/require"
	//		"golang.org/x/exp/maps"
	"pgregory.net/rapid"
	/*	"github.com/iotaledger/hive.go/core/kvstore"
		"github.com/iotaledger/hive.go/core/kvstore/mapdb"
		"github.com/iotaledger/hive.go/core/logger"
		"github.com/iotaledger/wasp/packages/database/dbkeys"
		"github.com/iotaledger/wasp/packages/isc"*/
	"github.com/iotaledger/wasp/packages/state"
	/*"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/util"*/)

//const constTestFolder = "basicWALTest"

type blockCacheTestSM struct { // State machine for block cache property based Rapid tests
	*blockCacheNoWALTestSM
	blocksInWAL []state.BlockHash
}

func (bctsmT *blockCacheTestSM) Init(t *rapid.T) {
	bctsmT.blockCacheNoWALTestSM = &blockCacheNoWALTestSM{}
	bctsmT.blockCacheNoWALTestSM.initWAL(t, NewMockedBlockWAL(), bctsmT.onAddBlock)
	bctsmT.blocksInWAL = make([]state.BlockHash, 0)
}

// Cleanup() // inheritted from blockCacheNoWALTestSM

func (bctsmT *blockCacheTestSM) Check(t *rapid.T) {
	bctsmT.blockCacheNoWALTestSM.Check(t)
	bctsmT.invariantAllBlocksInWALDifferent(t)
}

// AddNewBlock(t *rapid.T) // inheritted from blockCacheNoWALTestSM
// AddExistingBlock(t *rapid.T) // inheritted from blockCacheNoWALTestSM
// WriteBlockToDb(t *rapid.T) // inheritted from blockCacheNoWALTestSM
// CleanCache(t *rapid.T) // inheritted from blockCacheNoWALTestSM

// Maybe some files in WAL got corrupted
func (bctsmT *blockCacheTestSM) DeleteFromWAL(t *rapid.T) {
	if len(bctsmT.blocksInWAL) == 0 {
		t.Skip()
	}
	newWAL := NewMockedBlockWAL()
	newBlocksInWAL := make([]state.BlockHash, 0)
	gen := rapid.Bool()
	for i := range bctsmT.blocksInWAL {
		blockHash := bctsmT.blocksInWAL[i]
		if gen.Example(i) {
			t.Logf("Block %s was deleted from WAL", blockHash)
		} else {
			block, ok := bctsmT.blocks[blockHash]
			require.True(t, ok)
			err := newWAL.Write(block)
			require.NoError(t, err)
			newBlocksInWAL = append(newBlocksInWAL, blockHash)
		}
	}
	bctsmT.blocksInWAL = newBlocksInWAL
	bctsmT.bc.(*blockCache).wal = newWAL
	t.Logf("Delete some blocks from WAL completed")
}

func (bctsmT *blockCacheTestSM) GetBlockFromCache(t *rapid.T) {
	if len(bctsmT.blocksInCache) == 0 {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksInCache).Example()
	if ContainsBlock(blockHash, bctsmT.blocksInDB) || ContainsBlock(blockHash, bctsmT.blocksInWAL) {
		t.Skip()
	}
	bctsmT.tstGetBlockFromCacheAndDB(t, blockHash)
}

func (bctsmT *blockCacheTestSM) GetBlockFromWAL(t *rapid.T) {
	if len(bctsmT.blocksInWAL) == 0 {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksInWAL).Example()
	if ContainsBlock(blockHash, bctsmT.blocksInCache) || ContainsBlock(blockHash, bctsmT.blocksInDB) {
		t.Skip()
	}
	bctsmT.tstGetBlockNoCache(t, blockHash)
	t.Logf("Block from WAL %s obtained", blockHash)
}

func (bctsmT *blockCacheTestSM) GetBlockFromDB(t *rapid.T) {
	if len(bctsmT.blocksInDB) == 0 {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksInDB).Example()
	if ContainsBlock(blockHash, bctsmT.blocksInCache) || ContainsBlock(blockHash, bctsmT.blocksInWAL) {
		t.Skip()
	}
	bctsmT.tstGetBlockFromDB(t, blockHash)
}

func (bctsmT *blockCacheTestSM) GetBlockFromCacheAndDB(t *rapid.T) {
	if (len(bctsmT.blocksInCache) == 0) || len(bctsmT.blocksInDB) == 0 {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksInDB).Example()
	if !ContainsBlock(blockHash, bctsmT.blocksInCache) || ContainsBlock(blockHash, bctsmT.blocksInWAL) {
		t.Skip()
	}
	bctsmT.tstGetBlockFromCacheAndDB(t, blockHash)
}

func (bctsmT *blockCacheTestSM) GetBlockFromCacheAndWAL(t *rapid.T) {
	if len(bctsmT.blocksInCache) == 0 {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksInCache).Example()
	if !ContainsBlock(blockHash, bctsmT.blocksInWAL) || ContainsBlock(blockHash, bctsmT.blocksInDB) {
		t.Skip()
	}
	bctsmT.getAndCheckBlock(t, blockHash)
	t.Logf("Block from cache and WAL %s obtained", blockHash)
}

func (bctsmT *blockCacheTestSM) GetBlockFromWALAndDB(t *rapid.T) {
	if len(bctsmT.blocksInDB) == 0 {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksInDB).Example()
	if !ContainsBlock(blockHash, bctsmT.blocksInWAL) || ContainsBlock(blockHash, bctsmT.blocksInCache) {
		t.Skip()
	}
	bctsmT.tstGetBlockNoCache(t, blockHash)
	t.Logf("Block from WAL and DB %s obtained", blockHash)
}

func (bctsmT *blockCacheTestSM) GetBlockFromAll(t *rapid.T) {
	if len(bctsmT.blocksInDB) == 0 {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksInDB).Example()
	if !ContainsBlock(blockHash, bctsmT.blocksInWAL) || !ContainsBlock(blockHash, bctsmT.blocksInCache) {
		t.Skip()
	}
	bctsmT.getAndCheckBlock(t, blockHash)
	t.Logf("Block from cache, WAL and DB %s obtained", blockHash)
}

func (bctsmT *blockCacheTestSM) GetLostBlock(t *rapid.T) {
	if len(bctsmT.blocksInCache) == len(bctsmT.blocks) {
		t.Skip()
	}
	blockHash := rapid.SampledFrom(bctsmT.blocksNotInCache(t)).Example()
	if ContainsBlock(blockHash, bctsmT.blocksInDB) || ContainsBlock(blockHash, bctsmT.blocksInWAL) {
		t.Skip()
	}
	bctsmT.tstGetLostBlock(t, blockHash)
}

// Restart(t *rapid.T) // inheritted from blockCacheNoWALTestSM

func (bctsmT *blockCacheTestSM) onAddBlock(t *rapid.T, blockHash state.BlockHash) {
	if !ContainsBlock(blockHash, bctsmT.blocksInWAL) {
		bctsmT.blocksInWAL = append(bctsmT.blocksInWAL, blockHash)
	}
}

func (bctsmT *blockCacheTestSM) invariantAllBlocksInWALDifferent(t *rapid.T) {
	require.True(t, AllDifferent(bctsmT.blocksInWAL))
}

func TestBlockCacheFullPropBased(t *testing.T) {
	rapid.Check(t, rapid.Run[*blockCacheTestSM]())
}
