package smMessages

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/chain/statemanager/smUtils"
)

func TestMarchalUnmarshalBlockMessage(t *testing.T) {
	_, blocks, _ := smUtils.GetBlocks(t, 4, 1)
	for i := range blocks {
		t.Logf("Checking block %v: %v", i, blocks[i].GetHash())
		marshaled, err := NewBlockMessage(blocks[i], "SOMETHING").MarshalBinary()
		require.NoError(t, err)
		unmarshaled, err := NewBlockMessageFromBytes(marshaled)
		require.NoError(t, err)
		require.True(t, blocks[i].Equals(unmarshaled.GetBlock()))
	}
}
