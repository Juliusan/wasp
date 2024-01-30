package nodeconn2

import (
	"context"
	"errors"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/api"
	"github.com/iotaledger/wasp/packages/isc"
)

type (
	SubmitIotaBlockCallback func(iotago.BlockID, error, api.BlockState, api.BlockFailureReason)
	ChainOutputHandler      func(output *isc.AnchorAccountOutput)
	OtherOutputHandler      func(output isc.OutputWithID)
)

type ChainNodeConnection interface {
	GetChainID() isc.ChainID

	// Publishing can be canceled via the context.
	// The result must be returned via the callback, unless ctx is canceled first.
	// PublishTX handles promoting and reattachments until the tx is confirmed or the context is canceled.
	SubmitIotaBlock(
		ctx context.Context,
		iotaBlock *iotago.Block,
		callback SubmitIotaBlockCallback,
	)

	// outputsReceived receives only unspent (in that slot) outputs
	outputsReceived(
		iotago.SlotIndex,
		*isc.AnchorOutputWithID,
		*isc.AccountOutputWithID,
		[]*isc.BasicOutputWithID,
		[]*isc.FoundryOutputWithID,
		[]*isc.NFTOutputWithID,
	)
}

type NodeConnection interface {
	AttachChain(
		ctx context.Context,
		chainID isc.ChainID,
		chainOutputHandler ChainOutputHandler,
		otherOutputHandler OtherOutputHandler,
		//recvMilestone MilestoneHandler,
		//onChainConnect func(),
		//onChainDisconnect func(),
	) ChainNodeConnection
}

var (
	ErrOperationAborted      = errors.New("operation was aborted")
	ErrSubmitIotaBlockFailed = errors.New("submitting block failed")
)
