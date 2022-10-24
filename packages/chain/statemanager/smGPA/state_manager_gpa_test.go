package smGPA

import (
	"context"
	"fmt"
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

// Single node network. 8 blocks are sent to state manager. The result is checked
// by sending consensus requests, which force the access of the blocks.
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
		require.NoError(t, requireReceiveNoError(t, cbpRespChan, 5*time.Second))
	}
	cspInput, cspRespChan := smInputs.NewConsensusStateProposal(context.Background(), stateOutputs[7])
	tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeID: cspInput}).RunAll()
	require.NoError(t, requireReceiveAnything(cspRespChan, 5*time.Second))
	commitment, err := state.L1CommitmentFromBytes(stateOutputs[7].GetAliasOutput().StateMetadata)
	require.NoError(t, err)
	cdsInput, cdsRespChan := smInputs.NewConsensusDecidedState(context.Background(), stateOutputs[7].OutputID(), &commitment)
	tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeID: cdsInput}).RunAll()
	require.NoError(t, requireReceiveVState(t, cdsRespChan, 8, stateOutputs[7].OutputID(), &commitment, 5*time.Second))
}

// 10 nodes in a network. 8 blocks are sent to state manager of the first node.
// The result is checked by sending consensus requests to all the other 9 nodes,
// which force the access (and retrieval) of the blocks. For successful retrieval,
// several timer events are required for nodes to try to request blocks from peers.
func TestManyNodes(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	smTimers := NewStateManagerTimers()
	smTimers.StateManagerGetBlockRetry = 100 * time.Millisecond

	chainID, blocks, stateOutputs := smUtils.GetBlocks(t, 16, 1)
	nodeIDs := smUtils.MakeNodeIDs([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	sms := make(map[gpa.NodeID]gpa.GPA)
	for _, nodeID := range nodeIDs {
		_, sm := createStateManagerGpa(t, chainID, nodeID, nodeIDs, log, smTimers)
		sms[nodeID] = sm
	}
	tc := gpa.NewTestContext(sms)

	for i := range blocks {
		cbpInput, cbpRespChan := smInputs.NewChainBlockProduced(context.Background(), stateOutputs[i], blocks[i])
		t.Logf("Supplying block %s to node %s", blocks[i].GetHash(), nodeIDs[0])
		tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeIDs[0]: cbpInput}).RunAll()
		require.NoError(t, requireReceiveNoError(t, cbpRespChan, 5*time.Second))
	}
	//Nodes are checked sequentially
	now := time.Now()
	for i := 1; i < len(nodeIDs); i++ {
		cspInput, cspRespChan := smInputs.NewConsensusStateProposal(context.Background(), stateOutputs[7])
		tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeIDs[i]: cspInput}).RunAll()
		var j int
		received := false
		for j = 0; j < 10 && !received; {
			t.Logf("Waiting block %s on node %s, iteration %v", blocks[7].GetHash(), nodeIDs[i], j)
			if requireReceiveAnything(cspRespChan, 0*time.Second) != nil {
				now = now.Add(200 * time.Millisecond)
				sendTimerTickToNodes(tc, nodeIDs, now)
				j++
			} else {
				received = true
			}
		}
		require.Less(t, j, 10)
		commitment, err := state.L1CommitmentFromBytes(stateOutputs[7].GetAliasOutput().StateMetadata)
		require.NoError(t, err)
		cdsInput, cdsRespChan := smInputs.NewConsensusDecidedState(context.Background(), stateOutputs[7].OutputID(), &commitment)
		tc.WithInputs(map[gpa.NodeID]gpa.Input{nodeIDs[i]: cdsInput}).RunAll()
		require.NoError(t, requireReceiveVState(t, cdsRespChan, 8, stateOutputs[7].OutputID(), &commitment, 5*time.Second))
	}
	//Nodes are checked in parallel
	cspInputs := make(map[gpa.NodeID]gpa.Input)
	cspRespChans := make(map[gpa.NodeID]<-chan interface{})
	for i := 1; i < len(nodeIDs); i++ {
		nodeID := nodeIDs[i]
		cspInputs[nodeID], cspRespChans[nodeID] = smInputs.NewConsensusStateProposal(context.Background(), stateOutputs[15])
	}
	tc.WithInputs(cspInputs).RunAll()
	for nodeID, cspRespChan := range cspRespChans {
		var j int
		received := false
		for j = 0; j < 10 && !received; {
			t.Logf("Waiting block %s on node %s, iteration %v", blocks[15].GetHash(), nodeID, j)
			if requireReceiveAnything(cspRespChan, 0*time.Second) != nil {
				now = now.Add(200 * time.Millisecond)
				sendTimerTickToNodes(tc, nodeIDs, now)
				j++
			} else {
				received = true
			}
		}
		require.Less(t, j, 10)
	}
	commitment, err := state.L1CommitmentFromBytes(stateOutputs[15].GetAliasOutput().StateMetadata)
	require.NoError(t, err)
	cdsInputs := make(map[gpa.NodeID]gpa.Input)
	cdsRespChans := make(map[gpa.NodeID]<-chan *consGR.StateMgrDecidedState)
	for i := 1; i < len(nodeIDs); i++ {
		nodeID := nodeIDs[i]
		cdsInputs[nodeID], cdsRespChans[nodeID] = smInputs.NewConsensusDecidedState(context.Background(), stateOutputs[15].OutputID(), &commitment)
	}
	tc.WithInputs(cdsInputs).RunAll()
	for nodeID, cdsRespChan := range cdsRespChans {
		t.Logf("Waiting block %s on node %s", blocks[15].GetHash(), nodeID)
		require.NoError(t, requireReceiveVState(t, cdsRespChan, 16, stateOutputs[15].OutputID(), &commitment, 5*time.Second))
	}
}

// Single node network. Checks if block cache is cleaned via state manager
// timer events.
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

func requireReceiveNoError(t *testing.T, errChan <-chan (error), timeout time.Duration) error { //nolint:gocritic
	select {
	case err := <-errChan:
		require.NoError(t, err)
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("Waiting to receive no error timeouted")
	}
}

func requireReceiveAnything(anyChan <-chan (interface{}), timeout time.Duration) error { //nolint:gocritic
	select {
	case <-anyChan:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("Waiting to receive anything timeouted")
	}
}

func requireReceiveVState(t *testing.T, respChan <-chan (*consGR.StateMgrDecidedState), index uint32, aliasOutputID iotago.OutputID, l1c *state.L1Commitment, timeout time.Duration) error { //nolint:gocritic
	select {
	case smds := <-respChan:
		require.Equal(t, smds.VirtualStateAccess.BlockIndex(), index)
		//require.Equal(t, smds.AliasOutput.GetStateIndex(), index)   // TODO
		//require.Equal(t, smds.AliasOutput.OutputID(), aliasOutputID) // TODO
		require.True(t, state.EqualCommitments(trie.RootCommitment(smds.VirtualStateAccess.TrieNodeStore()), l1c.StateCommitment))
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("Waiting to receive state timeouted")
	}
}

func sendTimerTickToNodes(tc *gpa.TestContext, nodeIDs []gpa.NodeID, now time.Time) {
	inputs := make(map[gpa.NodeID]gpa.Input)
	for i := range nodeIDs {
		inputs[nodeIDs[i]] = smInputs.NewStateManagerTimerTick(now)
	}
	tc.WithInputs(inputs).RunAll()
}
