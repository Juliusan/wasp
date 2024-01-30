package nodeconn2

import (
	"context"

	"github.com/iotaledger/hive.go/log"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/nodeclient"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util/pipe"
)

type chainNodeConn struct {
	ctx                context.Context
	log                log.Logger
	indexerClient      nodeclient.IndexerClient
	chainID            isc.ChainID
	chainOutputHandler ChainOutputHandler
	otherOutputHandler OtherOutputHandler

	outputsReceivedPipe pipe.Pipe[*chainNodeConnOutputsReceivedData]
}

var _ ChainNodeConnection = &chainNodeConn{}

func newChainNodeConn(
	ctx context.Context,
	log log.Logger,
	chainID isc.ChainID,
	indexerClient nodeclient.IndexerClient,
	chainOutputHandler ChainOutputHandler,
	otherOutputHandler OtherOutputHandler,
) ChainNodeConnection {
	// TODO
	result := &chainNodeConn{
		ctx:                 ctx,
		log:                 log.NewChildLogger(chainID.ShortString()),
		chainID:             chainID,
		indexerClient:       indexerClient,
		chainOutputHandler:  chainOutputHandler,
		otherOutputHandler:  otherOutputHandler,
		outputsReceivedPipe: pipe.NewInfinitePipe[*chainNodeConnOutputsReceivedData](),
	}
	go result.initAndRun()
	return result
}

func (cnc *chainNodeConn) initAndRun() {
	err := cnc.initialise()
	if err == nil {
		cnc.run()
	}
}

func (cnc *chainNodeConn) initialise() error {
	cnc.log.LogDebugf("Initialising...")
	anchorAddress := cnc.chainID.AsAnchorID().ToAddress()
	outputID, anchorOutput, slotAnchor, err := cnc.indexerClient.Anchor(cnc.ctx, anchorAddress.(*iotago.AnchorAddress))
	if err != nil {
		cnc.log.LogErrorf("Initialising: failed to obtain latest anchor output: %v", err)
		return err
	}
	cnc.log.LogDebugf("Initialising: latest anchor output index %v %s obtained", anchorOutput.StateIndex, outputID)
	anchorOutputWithID := isc.NewAnchorOutputWithID(anchorOutput, *outputID)

	outputID, accountOutput, slotAccount, err := cnc.indexerClient.Account(cnc.ctx, anchorAddress.(*iotago.AccountAddress))
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
	cnc.log.LogDebugf("Initialising completed")
	return nil
}

func (cnc *chainNodeConn) GetChainID() isc.ChainID {
	return cnc.chainID
}

func (*chainNodeConn) PublishTX(
	ctx context.Context,
	tx *iotago.SignedTransaction,
	callback PublishTXCallback,
) error {
	// TODO
	return nil
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
	//defer smT.cleanupFun()
	outputsReceivedPipeCh := cnc.outputsReceivedPipe.Out()
	//timerTickCh := smT.parameters.TimeProvider.After(smT.parameters.StateManagerTimerTickPeriod)
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
			/*case now, ok := <-timerTickCh:
			  	if ok {
			  		smT.handleTimerTick(now)
			  		timerTickCh = smT.parameters.TimeProvider.After(smT.parameters.StateManagerTimerTickPeriod)
			  	} else {
			  		timerTickCh = nil
			  	}
			  case <-statusTimerCh:
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
