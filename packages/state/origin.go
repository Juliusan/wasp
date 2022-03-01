package state

import (
	"encoding/hex"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/kv/trie"
)

const OriginStateCommitmentHex = "5924dc2f04542fc93b02fa5c8b230f62110a9fbda78fca024cf58087bd32204f"

func OriginStateCommitment() trie.VCommitment {
	retBin, err := hex.DecodeString(OriginStateCommitmentHex)
	if err != nil {
		panic(err)
	}
	ret, err := CommitmentModel.VectorCommitmentFromBytes(retBin)
	if err != nil {
		panic(err)
	}
	return ret
}

func OriginStateData() *iscp.StateData {
	return iscp.NewStateData(OriginStateCommitment())
}