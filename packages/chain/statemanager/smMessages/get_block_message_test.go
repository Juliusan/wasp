package smMessages

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/chain/statemanager/smUtils"
)

func TestMarchalUnmarshalGetBlockMessage(t *testing.T) {
	_, blocks, _ := smUtils.GetBlocks(t, 4, 1)
	for i := range blocks {
		blockHash := blocks[i].GetHash()
		t.Logf("Checking block %v: %v", i, blockHash)
		marshaled, err := NewGetBlockMessage(blockHash, "SOMETHING").MarshalBinary()
		require.NoError(t, err)
		unmarshaled, err := NewGetBlockMessageFromBytes(marshaled)
		require.NoError(t, err)
		require.True(t, blockHash.Equals(unmarshaled.GetBlockHash()))
	}
}
