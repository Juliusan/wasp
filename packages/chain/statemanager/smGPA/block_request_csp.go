package smGPA

import (
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
	"github.com/iotaledger/wasp/packages/state"
)

type consensusStateProposalBlockRequest struct {
	consensusStateProposal *smInputs.ConsensusStateProposal
	done                   bool
	lastBlockHash          state.BlockHash
}

var _ blockRequest = &consensusStateProposalBlockRequest{}

func newConsensusStateProposalBlockRequest(input *smInputs.ConsensusStateProposal) (blockRequest, error) {
	stateCommitment, err := state.L1CommitmentFromAliasOutput(input.GetAliasOutputWithID().GetAliasOutput())
	if err != nil {
		return nil, err
	}
	return &consensusStateProposalBlockRequest{
		consensusStateProposal: input,
		done:                   false,
		lastBlockHash:          stateCommitment.BlockHash,
	}, nil
}

func (cspbrT *consensusStateProposalBlockRequest) getLastBlockHash() state.BlockHash {
	return cspbrT.lastBlockHash
}

func (cspbrT *consensusStateProposalBlockRequest) isValid() bool {
	return cspbrT.consensusStateProposal.IsValid()
}

func (cspbrT *consensusStateProposalBlockRequest) blockAvailable(block state.Block) {}

func (cspbrT *consensusStateProposalBlockRequest) markCompleted() {
	if !cspbrT.done {
		cspbrT.done = true
		cspbrT.consensusStateProposal.Respond()
	}
}
