package nodeconn2

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/iotaledger/hive.go/ds/shrinkingmap"
	"github.com/iotaledger/hive.go/log"
	"github.com/iotaledger/inx-app/pkg/nodebridge"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/api"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util/pipe"
)

type chainNodeConn struct {
	ctx                context.Context
	log                log.Logger
	nodeBridge         nodebridge.NodeBridge
	chainID            isc.ChainID
	chainOutputHandler ChainOutputHandler
	otherOutputHandler OtherOutputHandler
	chainStopHandler   chainStopHandler

	outputsReceivedPipe    pipe.Pipe[*chainNodeConnOutputsReceivedData]
	submitIotaBlockPipe    pipe.Pipe[*chainNodeConnSubmitIotaBlockData]
	confirmedIotaBlockPipe pipe.Pipe[*api.BlockMetadataResponse]

	submittedIotaBlocks *shrinkingmap.ShrinkingMap[iotago.BlockID, *chainNodeConnSubmitIotaBlockData]
}

type (
	chainStopHandler func()
)

var _ ChainNodeConnection = &chainNodeConn{}

func newChainNodeConn(
	ctx context.Context,
	log log.Logger,
	chainID isc.ChainID,
	nodeBridge nodebridge.NodeBridge,
	chainOutputHandler ChainOutputHandler,
	otherOutputHandler OtherOutputHandler,
	chainStopHandler chainStopHandler,
) ChainNodeConnection {
	// TODO
	result := &chainNodeConn{
		ctx:                    ctx,
		log:                    log.NewChildLogger(chainID.ShortString()),
		chainID:                chainID,
		nodeBridge:             nodeBridge,
		chainOutputHandler:     chainOutputHandler,
		otherOutputHandler:     otherOutputHandler,
		chainStopHandler:       chainStopHandler,
		outputsReceivedPipe:    pipe.NewInfinitePipe[*chainNodeConnOutputsReceivedData](),
		submitIotaBlockPipe:    pipe.NewInfinitePipe[*chainNodeConnSubmitIotaBlockData](),
		confirmedIotaBlockPipe: pipe.NewInfinitePipe[*api.BlockMetadataResponse](),
		submittedIotaBlocks:    shrinkingmap.New[iotago.BlockID, *chainNodeConnSubmitIotaBlockData](),
	}
	go result.initAndRun()
	return result
}

func (cnc *chainNodeConn) initAndRun() {
	err := cnc.initialise()
	if err == nil {
		cnc.run()
	} else {
		// TODO: possible race condition: if nodeconn is not fast enough to add chain nodeconn to its map,
		// this function (which aims at deleting it) might be called earlier, causing this chain nodeconn
		// to be included in nodeconn's map forever. See, if this happens.
		cnc.chainStopHandler()
	}
}

func (cnc *chainNodeConn) initialise() error {
	cnc.log.LogDebugf("Initialising...")
	indexerClient, err := cnc.nodeBridge.Indexer(cnc.ctx)
	if err != nil {
		cnc.log.LogErrorf("Initialising: failed to get indexer client: %v", err)
		return err
	}
	anchorAddress := cnc.chainID.AsAnchorID().ToAddress()
	outputID, anchorOutput, slotAnchor, err := indexerClient.Anchor(cnc.ctx, anchorAddress.(*iotago.AnchorAddress))
	if err != nil {
		cnc.log.LogErrorf("Initialising: failed to obtain latest anchor output: %v", err)
		return err
	}
	cnc.log.LogDebugf("Initialising: latest anchor output index %v %s obtained", anchorOutput.StateIndex, outputID)
	anchorOutputWithID := isc.NewAnchorOutputWithID(anchorOutput, *outputID)

	outputID, accountOutput, slotAccount, err := indexerClient.Account(cnc.ctx, anchorAddress.(*iotago.AccountAddress))
	if err != nil {
		cnc.log.LogDebugf("Initialising: failed to obtain latest account output: %v", err)
		return err
	}
	cnc.log.LogDebugf("Initialising: latest account output %s obtained", outputID)
	accountOutputWithID := isc.NewAccountOutputWithID(accountOutput, *outputID)

	var otherOutputs []isc.OutputWithID
	for empty := false; !empty; {
		select {
		case outputsReceivedData := <-cnc.outputsReceivedPipe.Out():
			cnc.log.LogDebugf("Initialising: checking received outputs of %v slot", outputsReceivedData.slotIndex())
			// TODO: is slot always a good way to sort
			if outputsReceivedData.slotIndex() >= slotAnchor {
				cnc.log.LogDebugf("Initialising: anchor output %s of slot %v is newer", outputsReceivedData.anchor.OutputID(), outputsReceivedData.slotIndex())
				anchorOutputWithID = outputsReceivedData.anchor
			}
			if outputsReceivedData.slotIndex() >= slotAccount {
				cnc.log.LogDebugf("Initialising: account output %s of slot %v is newer", outputsReceivedData.account.OutputID(), outputsReceivedData.slotIndex())
				accountOutputWithID = outputsReceivedData.account
			}
			otherOutputs = append(otherOutputs, joinOtherOutputs(outputsReceivedData)...)
		default:
			empty = true
		}
	}
	cnc.log.LogDebugf("Initialising: anchor output %s, account output %s and %v other outputs will be used initially",
		anchorOutputWithID.OutputID(), accountOutputWithID.OutputID(), len(otherOutputs))
	cnc.sendOutputs(anchorOutputWithID, accountOutputWithID, otherOutputs)

	go func() {
		cnc.log.LogDebugf("ListenToConfirmedBlocks thread: subscribing to blocks...")
		err := cnc.nodeBridge.ListenToConfirmedBlocks(cnc.ctx, func(blockMetadata *api.BlockMetadataResponse) error {
			cnc.confirmedIotaBlockPipe.In() <- blockMetadata
			return nil
		})
		if err != nil && !errors.Is(err, io.EOF) {
			cnc.log.LogErrorf("ListenToConfirmedBlocks thread: subscribing failed: %v", err)
			//nc.shutdownHandler.SelfShutdown("Subscribing to LedgerUpdates failed", true)
		}
		cnc.log.LogDebugf("ListenToConfirmedBlocks thread: completed")
		/*if nc.ctx.Err() == nil {
			// shutdown in case there isn't a shutdown already in progress
			nc.shutdownHandler.SelfShutdown("INX connection closed", true)
		}*/
	}()
	cnc.log.LogDebugf("Initialising completed")
	return nil
}

func (cnc *chainNodeConn) GetChainID() isc.ChainID {
	return cnc.chainID
}

func (cnc *chainNodeConn) SubmitIotaBlock(
	ctx context.Context,
	iotaBlock *iotago.Block,
	callback SubmitIotaBlockCallback,
) {
	cnc.submitIotaBlockPipe.In() <- newChainNodeConnSubmitIotaBlockData(ctx, iotaBlock, callback)
}

func (cnc *chainNodeConn) outputsReceived(
	slot iotago.SlotIndex,
	anchor *isc.AnchorOutputWithID,
	account *isc.AccountOutputWithID,
	basics []*isc.BasicOutputWithID,
	foundries []*isc.FoundryOutputWithID,
	nfts []*isc.NFTOutputWithID,
) {
	cnc.outputsReceivedPipe.In() <- newChainNodeConnOutputsReceivedData(slot, anchor, account, basics, foundries, nfts)
}

func (cnc *chainNodeConn) run() {
	defer cnc.chainStopHandler()
	outputsReceivedPipeCh := cnc.outputsReceivedPipe.Out()
	submitIotaBlockPipeCh := cnc.submitIotaBlockPipe.Out()
	confirmedIotaBlockPipeCh := cnc.confirmedIotaBlockPipe.Out()
	timerTickCh := time.After(1 * time.Minute) // TODO: parametrise smT.parameters.TimeProvider.After(smT.parameters.StateManagerTimerTickPeriod)
	//statusTimerCh := smT.parameters.TimeProvider.After(constStatusTimerTime)
	for {
		if cnc.ctx.Err() != nil {
			/*if smT.shutdownCoordinator == nil {
				return
			}*/
			// TODO what should the statemgr wait for?
			//smT.shutdownCoordinator.WaitNestedWithLogging(1 * time.Second)
			cnc.log.LogDebugf("Stopping chain nodeconn, because context was closed")
			//smT.shutdownCoordinator.Done()
			return
		}
		select {
		case outputsReceivedData, ok := <-outputsReceivedPipeCh:
			if ok {
				cnc.handleOutputsReceived(outputsReceivedData)
			} else {
				outputsReceivedPipeCh = nil
			}
		case submitIotaBlockData, ok := <-submitIotaBlockPipeCh:
			if ok {
				cnc.handleSubmitIotaBlock(submitIotaBlockData)
			} else {
				submitIotaBlockPipeCh = nil
			}
		case confirmedIotaBlockData, ok := <-confirmedIotaBlockPipeCh:
			if ok {
				cnc.handleConfirmedIotaBlock(confirmedIotaBlockData)
			} else {
				confirmedIotaBlockPipeCh = nil
			}
		case now, ok := <-timerTickCh:
			if ok {
				cnc.handleTimerTick(now)
				timerTickCh = time.After(1 * time.Minute) // TODO: parametrise smT.parameters.TimeProvider.After(smT.parameters.StateManagerTimerTickPeriod)
			} else {
				timerTickCh = nil
			}
			/*case <-statusTimerCh:
			statusTimerCh = smT.parameters.TimeProvider.After(constStatusTimerTime)
			smT.log.Debugf("State manager loop iteration; there are %v inputs, %v messages, %v public key changes waiting to be handled",
				smT.inputPipe.Len(), smT.messagePipe.Len(), smT.nodePubKeysPipe.Len())*/
		case <-cnc.ctx.Done():
			continue
		}
	}
}

func (cnc *chainNodeConn) handleOutputsReceived(data *chainNodeConnOutputsReceivedData) {
	otherOutputs := joinOtherOutputs(data)
	cnc.sendOutputs(data.anchor, data.account, otherOutputs)
}

func joinOtherOutputs(data *chainNodeConnOutputsReceivedData) []isc.OutputWithID {
	basicOutputs := data.basicOutputs()
	foundryOutputs := data.foundryOutputs()
	nftOutputs := data.nftOutputs()
	otherOutputs := make([]isc.OutputWithID, len(basicOutputs)+len(foundryOutputs)+len(nftOutputs))
	i := 0
	for j := 0; j < len(basicOutputs); j++ {
		otherOutputs[i] = basicOutputs[j]
		i++
	}
	for j := 0; j < len(foundryOutputs); j++ {
		otherOutputs[i] = foundryOutputs[j]
		i++
	}
	for j := 0; j < len(nftOutputs); j++ {
		otherOutputs[i] = nftOutputs[j]
		i++
	}
	return otherOutputs
}

func (cnc *chainNodeConn) sendOutputs(anchorOutput *isc.AnchorOutputWithID, accountOutput *isc.AccountOutputWithID, otherOutputs []isc.OutputWithID) {
	if anchorOutput != nil || accountOutput != nil {
		cnc.chainOutputHandler(isc.NewAnchorAccountOutput(anchorOutput, accountOutput))
	}
	for _, output := range otherOutputs {
		cnc.otherOutputHandler(output)
	}
}

func (cnc *chainNodeConn) handleSubmitIotaBlock(data *chainNodeConnSubmitIotaBlockData) {
	blockTime := data.iotaBlock().Header.IssuingTime
	cnc.log.LogDebugf("Submitting Iota block issued at %v...", blockTime)
	blockID, err := cnc.nodeBridge.SubmitBlock(data.context(), data.iotaBlock())
	if err != nil {
		data.callBack(blockID, err, api.BlockStateUnknown, api.BlockFailureInvalid)
		cnc.log.LogErrorf("Submitting Iota block issued at %v failed: %v", blockTime, err)
	}
	// NOTE: not checking for block ID duplication here
	cnc.submittedIotaBlocks.Set(blockID, data)
	cnc.log.LogErrorf("Submitting Iota block issued at %v succeeded, block ID: %s", blockTime, blockID)
}

func (cnc *chainNodeConn) handleConfirmedIotaBlock(response *api.BlockMetadataResponse) {
	data, exists := cnc.submittedIotaBlocks.Get(response.BlockID)
	if exists {
		cnc.log.LogDebugf("Iota block %s confirmed...", response.BlockID)
		if data.context().Err() == nil {
			data.callBack(response.BlockID, nil, response.BlockState, response.BlockFailureReason)
		} else {
			cnc.log.LogDebugf("Iota block %s confirmed: request context is cancelled, skipping callback: %v",
				response.BlockID, data.context().Err())
		}
		cnc.submittedIotaBlocks.Delete(response.BlockID)
		cnc.log.LogDebugf("Iota block %s confirmed handled", response.BlockID)
	}
}

func (cnc *chainNodeConn) handleTimerTick(now time.Time) {
	cnc.log.LogDebugf("Timer tick...")
	cnc.submittedIotaBlocks.ForEach(func(blockID iotago.BlockID, data *chainNodeConnSubmitIotaBlockData) bool {
		if data.context().Err() != nil {
			cnc.log.LogDebugf("Timer tick: submit Iota block %s request's context canceled: %v, removing it", blockID, data.context().Err())
			cnc.submittedIotaBlocks.Delete(blockID)
			return true
		}
		if !data.isValid() {
			cnc.log.LogDebugf("Timer tick: submit Iota block %s request's timeouted, removing it", blockID)
			data.callBack(blockID, errors.New("timeouted"), api.BlockStateUnknown, api.BlockFailureInvalid)
			cnc.submittedIotaBlocks.Delete(blockID)
			return true
		}
		cnc.log.LogDebugf("Timer tick: requesting for block's %s metadata...", blockID)
		response, err := cnc.nodeBridge.BlockMetadata(cnc.ctx, blockID)
		if err != nil {
			cnc.log.LogErrorf("Timer tick: failed to request for block's %s metadata", blockID)
			return true
		}
		if response.BlockState == api.BlockStateFinalized ||
			response.BlockState == api.BlockStateRejected ||
			response.BlockState == api.BlockStateFailed {
			cnc.log.LogDebugf("Timer tick: block %s is finalized/rejected/failed (%v), reason: %v, responding",
				blockID, response.BlockState, response.BlockFailureReason)
			data.callBack(blockID, ErrSubmitIotaBlockFailed, response.BlockState, response.BlockFailureReason)
			cnc.submittedIotaBlocks.Delete(blockID)
			return true
		}
		cnc.log.LogDebugf("Timer tick: block %s is still being handled", blockID)
		return true
	})
	cnc.log.LogDebugf("Timer tick handled, %v submit block requests waiting for response left", cnc.submittedIotaBlocks.Size())
}
