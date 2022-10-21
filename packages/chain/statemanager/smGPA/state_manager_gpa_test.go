package smGPA

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/core/kvstore/mapdb"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/trie.go/trie"
	"github.com/iotaledger/wasp/packages/chain/aaa2/cons/gr"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smUtils"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
)

func TestBasic(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	chainID, blocks, stateOutputs := smUtils.GetBlocks(t, 8, 1)

	nodeID := smUtils.MakeNodeID(0)
	nr := smUtils.NewNodeRandomiser(nodeID, []gpa.NodeID{nodeID})
	store := mapdb.NewMapDB()
	log = log.Named(nodeID.String()).Named("c-" + chainID.ShortString())

	sm := New(chainID, nr, store, log)
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
