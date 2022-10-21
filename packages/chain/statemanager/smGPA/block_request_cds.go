package smGPA

import (
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
	"github.com/iotaledger/wasp/packages/state"
)

type createOriginStateFun func() (state.VirtualStateAccess, error)

type consensusDecidedStateBlockRequest struct {
	consensusDecidedState *smInputs.ConsensusDecidedState
	done                  bool
	blocks                []state.Block
	createOriginStateFun  createOriginStateFun
}

var _ blockRequest = &consensusDecidedStateBlockRequest{}

func newConsensusDecidedStateBlockRequest(input *smInputs.ConsensusDecidedState, createOriginStateFun createOriginStateFun) blockRequest {
	return &consensusDecidedStateBlockRequest{
		consensusDecidedState: input,
		done:                  false,
		blocks:                make([]state.Block, 0),
		createOriginStateFun:  createOriginStateFun,
	}
}

func (cdsbrT *consensusDecidedStateBlockRequest) getLastBlockHash() state.BlockHash {
	return cdsbrT.consensusDecidedState.GetStateCommitment().BlockHash
}

func (cdsbrT *consensusDecidedStateBlockRequest) isValid() bool {
	return cdsbrT.consensusDecidedState.IsValid()
}

func (cdsbrT *consensusDecidedStateBlockRequest) blockAvailable(block state.Block) {
	cdsbrT.blocks = append(cdsbrT.blocks, block)
}

func (cdsbrT *consensusDecidedStateBlockRequest) markCompleted() {
	if !cdsbrT.done {
		cdsbrT.done = true
		vState, err := cdsbrT.createOriginStateFun()
		if err != nil {
			return
		}
		for i := len(cdsbrT.blocks) - 1; i >= 0; i-- {
			calculatedStateCommitment := state.RootCommitment(vState.TrieNodeStore())
			if !state.EqualCommitments(calculatedStateCommitment, cdsbrT.blocks[i].PreviousL1Commitment().StateCommitment) {
				return
			}
			err := vState.ApplyBlock(cdsbrT.blocks[i])
			if err != nil {
				return
			}
			vState.Commit() // TODO: is it needed
		}
		cdsbrT.consensusDecidedState.Respond(nil, nil, vState) // TODO: return alias output and state baseline
	}
}
