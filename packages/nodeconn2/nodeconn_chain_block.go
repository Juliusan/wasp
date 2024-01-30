package nodeconn2

import (
	"context"
	"time"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/api"
)

type chainNodeConnSubmitIotaBlockData struct {
	ctx       context.Context
	timestamp time.Time
	block     *iotago.Block
	callback  SubmitIotaBlockCallback
}

var submitBlockValidityConst = 1 * time.Hour

func newChainNodeConnSubmitIotaBlockData(
	ctx context.Context,
	block *iotago.Block,
	callback SubmitIotaBlockCallback,
) *chainNodeConnSubmitIotaBlockData {
	return &chainNodeConnSubmitIotaBlockData{
		ctx:       ctx,
		timestamp: time.Now(),
		block:     block,
		callback:  callback,
	}
}

func (cncsbd *chainNodeConnSubmitIotaBlockData) context() context.Context {
	return cncsbd.ctx
}

func (cncsbd *chainNodeConnSubmitIotaBlockData) isValid() bool {
	return cncsbd.timestamp.Add(submitBlockValidityConst).After(time.Now())
}

func (cncsbd *chainNodeConnSubmitIotaBlockData) iotaBlock() *iotago.Block {
	return cncsbd.block
}

func (cncsbd *chainNodeConnSubmitIotaBlockData) callBack(blockID iotago.BlockID, err error, state api.BlockState, reason api.BlockFailureReason) {
	cncsbd.callback(blockID, err, state, reason)
}
