// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package smUtils

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
)

func TestGetRandomOtherNodeIDs(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	meIndex := 3
	nodeIDsToGet := 5
	iterationCount := 13

	nodeIDs := MakeNodeIDs([]int{0, 1, 2, 3, 4, 5, 6, 7}) // 7 nodes excluding self
	me := nodeIDs[meIndex]
	randomiser := NewNodeRandomiser(me, nodeIDs, log)
	testGetRandomOtherNodeIDs(t, randomiser, nodeIDsToGet, nodeIDsToGet, iterationCount, nodeIDs, me)
}

func TestGetRandomOtherNodeIDsToFew(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	meIndex := 3
	nodeIDsToGet := 5
	iterationCount := 1

	nodeIDs := MakeNodeIDs([]int{0, 1, 2, 3}) // 3 nodes excluding self
	me := nodeIDs[meIndex]
	randomiser := NewNodeRandomiser(me, nodeIDs, log)
	testGetRandomOtherNodeIDs(t, randomiser, nodeIDsToGet, 3, iterationCount, nodeIDs, me)
}

func TestGetRandomOtherNodeIDsAfterChanges(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	nodeIDsToGet := 5
	iterationCount := 7

	nodeIDs0 := MakeNodeIDs([]int{0, 1, 2, 3, 4, 5, 6, 7})
	nodeIDs1 := MakeNodeIDs([]int{0, 2, 3, 5, 6, 7})
	nodeIDs2 := MakeNodeIDs([]int{0, 2, 3, 5, 6, 7, 8})
	nodeIDs3 := MakeNodeIDs([]int{0, 2, 3, 5, 6, 7})
	nodeIDs4 := MakeNodeIDs([]int{0, 2, 3, 4, 5, 6, 7, 9})
	me := nodeIDs0[0]
	randomiser := NewNodeRandomiser(me, nodeIDs0, log)
	testGetRandomOtherNodeIDs(t, randomiser, nodeIDsToGet, nodeIDsToGet, iterationCount, nodeIDs0, me)
	randomiser.UpdateNodeIDs(nodeIDs1)
	testGetRandomOtherNodeIDs(t, randomiser, nodeIDsToGet, nodeIDsToGet, iterationCount, nodeIDs1, me)
	randomiser.UpdateNodeIDs(nodeIDs2)
	testGetRandomOtherNodeIDs(t, randomiser, nodeIDsToGet, nodeIDsToGet, iterationCount, nodeIDs2, me)
	randomiser.UpdateNodeIDs(nodeIDs3)
	testGetRandomOtherNodeIDs(t, randomiser, nodeIDsToGet, nodeIDsToGet, iterationCount, nodeIDs3, me)
	randomiser.UpdateNodeIDs(nodeIDs4)
	testGetRandomOtherNodeIDs(t, randomiser, nodeIDsToGet, nodeIDsToGet, iterationCount, nodeIDs4, me)
}

func testGetRandomOtherNodeIDs(t *testing.T, randomiser NodeRandomiser, nodeIDsToGet, nodeIDsGot, iterationCount int, nodeIDs []gpa.NodeID, me gpa.NodeID) {
	nodeIDFounds := make(map[gpa.NodeID]bool)
	for i := 0; i < iterationCount; i++ {
		t.Logf("Iteration %v...", i)
		randomNodeIDs := randomiser.GetRandomOtherNodeIDs(nodeIDsToGet)
		require.Equal(t, nodeIDsGot, len(randomNodeIDs))
		for j := range randomNodeIDs {
			nodeIDFounds[randomNodeIDs[j]] = true
			t.Logf("\tComparing nodeID %v, id %v with me...", j, randomNodeIDs[j])
			require.False(t, randomNodeIDs[j].Equals(me))
			t.Logf("\tComparing nodeIDs %v, id %v...", j, randomNodeIDs[j])
			for k := range randomNodeIDs[j+1:] {
				kk := k + j + 1
				t.Logf("\t\t and %v, id %v", kk, randomNodeIDs[kk])
				require.False(t, randomNodeIDs[j].Equals(randomNodeIDs[kk]))
			}
		}
	}
	t.Logf("Checking if all nodeIDs were returned...")
	for i := range nodeIDs {
		_, ok := nodeIDFounds[nodeIDs[i]]
		if nodeIDs[i].Equals(me) {
			t.Logf("\tMe nodeID %v should not be returned", nodeIDs[i])
			require.False(t, ok)
		} else {
			t.Logf("\tNodeID %v should be at least once", nodeIDs[i])
			require.True(t, ok)
		}
	}
	t.Logf("Checking if all returned nodeIDs are correct...")
	containsFun := func(ni gpa.NodeID) bool {
		for i := range nodeIDs {
			if ni.Equals(nodeIDs[i]) {
				return true
			}
		}
		return false
	}
	for nodeID := range nodeIDFounds {
		t.Logf("\tNodeID %v", nodeID)
		require.True(t, containsFun(nodeID))
	}
}
