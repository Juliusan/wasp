package smGPA

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/core/kvstore/mapdb"
	"github.com/iotaledger/hive.go/core/logger"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/trie.go/trie"
	"github.com/iotaledger/wasp/packages/chain/aaa2/cons/gr"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smUtils"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
)

func TestBasic(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	chainID, blocks, stateOutputs := smUtils.GetBlocks(t, 8, 1)
	nodeID := smUtils.MakeNodeID(0)
	_, sm := createStateManagerGpa(t, chainID, nodeID, []gpa.NodeID{nodeID}, log)
	tc := gpa.NewTestContext(map[gpa.NodeID]gpa.GPA{nodeID: sm})

	for i := range blocks {
		cbpInput, cbpRespChan := smInputs.NewChainBlockProduced(context.Background(), stateOutputs[i], blocks[i])
		t.Logf("Supplying block %s to node %s", blocks[i].GetHash(), nodeID)
		tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeID: cbpInput}).RunAll()
		requireReceiveNoError(t, cbpRespChan, 5*time.Second)
	}
	cspInput, cspRespChan := smInputs.NewConsensusStateProposal(context.Background(), stateOutputs[7])
	tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeID: cspInput}).RunAll()
	requireReceiveAnything(t, cspRespChan, 5*time.Second)
	commitment, err := state.L1CommitmentFromBytes(stateOutputs[7].GetAliasOutput().StateMetadata)
	require.NoError(t, err)
	cdsInput, cdsRespChan := smInputs.NewConsensusDecidedState(context.Background(), stateOutputs[7].OutputID(), &commitment)
	tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeID: cdsInput}).RunAll()
	requireReceiveVState(t, cdsRespChan, 8, stateOutputs[7].OutputID(), &commitment, 5*time.Second)
	//smInputs.NewConsensusDecidedState()
}

func TestBlockCacheCleaningAuto(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	tp := smUtils.NewArtifficialTimeProvider()
	smTimers := NewStateManagerTimers(tp)
	smTimers.BlockCacheBlocksInCacheDuration = 300 * time.Millisecond
	smTimers.BlockCacheBlockCleaningPeriod = 70 * time.Millisecond

	chainID, blocks, _ := smUtils.GetBlocks(t, 6, 2)
	nodeID := smUtils.MakeNodeID(0)
	_, sm := createStateManagerGpa(t, chainID, nodeID, []gpa.NodeID{nodeID}, log, smTimers)
	tc := gpa.NewTestContext(map[gpa.NodeID]gpa.GPA{nodeID: sm})

	advanceTimeAndTimerTickFun := func(advance time.Duration) {
		tp.SetNow(tp.GetNow().Add(advance))
		tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeID: smInputs.NewStateManagerTimerTick(tp.GetNow())}).RunAll()
	}

	blockCache := sm.(*stateManagerGPA).blockCache
	blockCache.AddBlock(blocks[0])
	blockCache.AddBlock(blocks[1])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	advanceTimeAndTimerTickFun(100 * time.Millisecond)
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	blockCache.AddBlock(blocks[2])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	advanceTimeAndTimerTickFun(100 * time.Millisecond)
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	blockCache.AddBlock(blocks[3])
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	advanceTimeAndTimerTickFun(80 * time.Millisecond)
	tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeID: smInputs.NewStateManagerTimerTick(tp.GetNow())}).RunAll()
	require.NotNil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	advanceTimeAndTimerTickFun(100 * time.Millisecond)
	blockCache.AddBlock(blocks[4])
	require.Nil(t, blockCache.GetBlock(blocks[0].GetHash()))
	require.Nil(t, blockCache.GetBlock(blocks[1].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	advanceTimeAndTimerTickFun(100 * time.Millisecond)
	require.Nil(t, blockCache.GetBlock(blocks[2].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	advanceTimeAndTimerTickFun(100 * time.Millisecond)
	require.Nil(t, blockCache.GetBlock(blocks[3].GetHash()))
	require.NotNil(t, blockCache.GetBlock(blocks[4].GetHash()))
	advanceTimeAndTimerTickFun(200 * time.Millisecond)
	require.Nil(t, blockCache.GetBlock(blocks[4].GetHash()))
}

func createStateManagerGpa(t *testing.T, chainID *isc.ChainID, me gpa.NodeID, nodeIDs []gpa.NodeID, log *logger.Logger, timers ...StateManagerTimers) (*smUtils.NodeRandomiser, gpa.GPA) {
	nr := smUtils.NewNodeRandomiser(me, nodeIDs)
	store := mapdb.NewMapDB()
	log = log.Named(me.String()).Named("c-" + chainID.ShortString())
	sm := New(chainID, nr, store, log, timers...)
	return nr, sm
}

func requireReceiveNoError(t *testing.T, errChan <-chan (error), timeout time.Duration) { //nolint:gocritic
	select {
	case err := <-errChan:
		require.NoError(t, err)
	case <-time.After(timeout):
		t.Logf("Waiting to receive no error timeouted")
		t.FailNow()
	}
}

func requireReceiveAnything(t *testing.T, anyChan <-chan (interface{}), timeout time.Duration) { //nolint:gocritic
	select {
	case <-anyChan:
		return
	case <-time.After(timeout):
		t.Logf("Waiting to receive anything timeouted")
		t.FailNow()
	}
}

func requireReceiveVState(t *testing.T, respChan <-chan (*consGR.StateMgrDecidedState), index uint32, aliasOutputID iotago.OutputID, l1c *state.L1Commitment, timeout time.Duration) { //nolint:gocritic
	select {
	case smds := <-respChan:
		require.Equal(t, smds.VirtualStateAccess.BlockIndex(), index)
		//require.Equal(t, smds.AliasOutput.GetStateIndex(), index)   // TODO
		//require.Equal(t, smds.AliasOutput.OutputID(), aliasOutputID) // TODO
		require.True(t, state.EqualCommitments(trie.RootCommitment(smds.VirtualStateAccess.TrieNodeStore()), l1c.StateCommitment))
		return
	case <-time.After(timeout):
		t.Logf("Waiting to receive state timeouted")
		t.FailNow()
	}
}
