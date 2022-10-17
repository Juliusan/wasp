package smGPA

import (
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
	"github.com/iotaledger/wasp/packages/state"
)

type createOriginStateFun func() (state.VirtualStateAccess, error)

type consensusDecidedStateBlockRequest struct {
	consensusDecidedState *smInputs.ConsensusDecidedState
	done                  bool
	lastBlockHash         state.BlockHash
	blocks                []state.Block
	createOriginStateFun  createOriginStateFun
}

var _ blockRequest = &consensusStateProposalBlockRequest{}

func newConsensusDecidedStateBlockRequest(input *smInputs.ConsensusDecidedState, createOriginStateFun createOriginStateFun) blockRequest {
	return &consensusDecidedStateBlockRequest{
		consensusDecidedState: input,
		done:                  false,
		lastBlockHash:         input.GetStateCommitment().BlockHash,
		blocks:                make([]state.Block, 0),
		createOriginStateFun:  createOriginStateFun,
	}
}

func (cspbrT *consensusDecidedStateBlockRequest) getLastBlockHash() state.BlockHash {
	return cspbrT.lastBlockHash
}

func (cspbrT *consensusDecidedStateBlockRequest) blockAvailable(block state.Block) {
	cspbrT.blocks = append(cspbrT.blocks, block)
}

func (cspbrT *consensusDecidedStateBlockRequest) markCompleted() {
	if !cspbrT.done {
		cspbrT.done = true
		vState, err := cspbrT.createOriginStateFun()
		if err != nil {
			return
		}
		for i := len(cspbrT.blocks) - 1; i >= 0; i-- {
			calculatedStateCommitment := state.RootCommitment(vState.TrieNodeStore())
			if !state.EqualCommitments(calculatedStateCommitment, cspbrT.blocks[i].PreviousL1Commitment().StateCommitment) {
				return
			}
			err := vState.ApplyBlock(cspbrT.blocks[i])
			if err != nil {
				return
			}
			vState.Commit() // TODO: is it needed
		}
		cspbrT.consensusDecidedState.Respond(nil, nil, vState) // TODO: return alias output and state baseline
	}
}
