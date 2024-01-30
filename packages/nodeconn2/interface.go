package nodeconn2

import (
	"context"
	"errors"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/isc"
)

type (
	PublishTXCallback  func(tx *iotago.Transaction, confirmed bool)
	ChainOutputHandler func(output *isc.AnchorAccountOutput)
	OtherOutputHandler func(output isc.OutputWithID)
)

type ChainNodeConnection interface {
	GetChainID() isc.ChainID

	// Publishing can be canceled via the context.
	// The result must be returned via the callback, unless ctx is canceled first.
	// PublishTX handles promoting and reattachments until the tx is confirmed or the context is canceled.
	PublishTX(
		ctx context.Context,
		tx *iotago.SignedTransaction,
		callback PublishTXCallback,
	) error

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
	// Alias outputs are expected to be returned in order. Considering the Hornet node, the rules are:
	//   - Upon Attach -- existing unspent anchor output is returned FIRST.
	//   - Upon receiving a spent/unspent AO from L1 they are returned in
	//     the same order, as the milestones are issued.
	//   - If a single milestone has several anchor outputs, they have to be ordered
	//     according to the chain of TXes.
	//
	// NOTE: Any out-of-order AO will be considered as a rollback or AO by the chain impl.
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

var ErrOperationAborted = errors.New("operation was aborted")
