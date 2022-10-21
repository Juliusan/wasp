package smInputs

import (
	"context"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/chain/aaa2/cons/gr"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/state"
)

type ConsensusDecidedState struct {
	context         context.Context
	aliasOutputID   iotago.OutputID
	stateCommitment *state.L1Commitment
	resultCh        chan<- *consGR.StateMgrDecidedState
}

var _ gpa.Input = &ConsensusDecidedState{}

func NewConsensusDecidedState(ctx context.Context, aliasOutputID iotago.OutputID, sc *state.L1Commitment) (*ConsensusDecidedState, <-chan *consGR.StateMgrDecidedState) {
	resultChannel := make(chan *consGR.StateMgrDecidedState, 1)
	return &ConsensusDecidedState{
		context:         ctx,
		aliasOutputID:   aliasOutputID,
		stateCommitment: sc,
		resultCh:        resultChannel,
	}, resultChannel
}

func (cdsT *ConsensusDecidedState) GetAliasOutputID() iotago.OutputID {
	return cdsT.aliasOutputID
}

func (cdsT *ConsensusDecidedState) GetStateCommitment() *state.L1Commitment {
	return cdsT.stateCommitment
}

func (cdsT *ConsensusDecidedState) IsValid() bool {
	return cdsT.context.Err() == nil
}

func (cdsT *ConsensusDecidedState) Respond(aliasOutput *isc.AliasOutputWithID, stateBaseline coreutil.StateBaseline, virtualStateAccess state.VirtualStateAccess) {
	if cdsT.IsValid() && !cdsT.IsResultChClosed() {
		cdsT.resultCh <- &consGR.StateMgrDecidedState{
			AliasOutput:        aliasOutput,
			StateBaseline:      stateBaseline,
			VirtualStateAccess: virtualStateAccess,
		}
		cdsT.closeResultCh()
	}
}

func (cdsT *ConsensusDecidedState) IsResultChClosed() bool {
	return cdsT.resultCh == nil
}

func (cdsT *ConsensusDecidedState) closeResultCh() {
	close(cdsT.resultCh)
	cdsT.resultCh = nil
}
