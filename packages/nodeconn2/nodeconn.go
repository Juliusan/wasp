package nodeconn2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/iotaledger/hive.go/app/shutdown"
	"github.com/iotaledger/hive.go/log"
	"github.com/iotaledger/hive.go/runtime/timeutil"
	"github.com/iotaledger/inx-app/pkg/nodebridge"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/nodeclient"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util/pipe"
)

type nodeConn struct {
	ctx             context.Context
	log             log.Logger
	nodeBridge      nodebridge.NodeBridge
	indexerClient   nodeclient.IndexerClient
	inxNodeClient   *nodeclient.Client
	chainNodeConns  map[isc.ChainID]ChainNodeConnection
	shutdownHandler *shutdown.ShutdownHandler

	newChainPipe     pipe.Pipe[ChainNodeConnection]
	stoppedChainPipe pipe.Pipe[isc.ChainID]
	ledgerUpdatePipe pipe.Pipe[*nodebridge.LedgerUpdate]
}

var _ NodeConnection = &nodeConn{}

func NewNodeConn(
	ctx context.Context,
	log log.Logger,
	nodeBridge nodebridge.NodeBridge,
	shutdownHandler *shutdown.ShutdownHandler,
) (NodeConnection, error) {
	ncLog := log.NewChildLogger("nc")
	indexerClient, err := nodeBridge.Indexer(ctx)
	if err != nil {
		ncLog.LogErrorf("Failed to get indexer client: %v", err)
		return nil, err
	}
	inxNodeClient, err := nodeBridge.INXNodeClient()
	if err != nil {
		ncLog.LogErrorf("Failed to get inx node client: %v", err)
		return nil, err
	}
	result := &nodeConn{
		ctx:              ctx,
		log:              ncLog,
		nodeBridge:       nodeBridge,
		indexerClient:    indexerClient,
		inxNodeClient:    inxNodeClient,
		chainNodeConns:   make(map[isc.ChainID]ChainNodeConnection),
		shutdownHandler:  shutdownHandler,
		newChainPipe:     pipe.NewInfinitePipe[ChainNodeConnection](),
		stoppedChainPipe: pipe.NewInfinitePipe[isc.ChainID](),
		ledgerUpdatePipe: pipe.NewInfinitePipe[*nodebridge.LedgerUpdate](),
	}
	err = result.initialise()
	if err != nil {
		return nil, err
	}
	result.run()
	return result, nil
}

func (nc *nodeConn) initialise() error {
	// the node bridge needs to be started before waiting for L1 to become synced,
	// otherwise the NodeStatus would never be updated and "syncAndSetProtocolParameters" would be stuck
	// in an infinite loop
	go func() {
		nc.log.LogDebugf("NodeBridge thread: started")
		nc.nodeBridge.Run(nc.ctx)
		nc.log.LogDebugf("NodeBridge thread: node bridge stopped")

		// if the Run function returns before the context was actually canceled,
		// it means that the connection to L1 node must have failed.
		if !errors.Is(nc.ctx.Err(), context.Canceled) {
			nc.log.LogErrorf("NodeBridge thread: INX connection to node dropped: %v", nc.ctx.Err())
			nc.shutdownHandler.SelfShutdown("INX connection to node dropped", true)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer timeutil.CleanupTicker(ticker)

	// make sure the node is connected to at least one other peer
	// otherwise the node status may not reflect the network status
	for connected := false; !connected; {
		nc.log.LogDebug("Initialising: waiting for L1 to be connected to other peers...")
		select {
		case <-ticker.C:
			mgmClient, err := nc.nodeBridge.Management(nc.ctx)
			if err != nil {
				nc.log.LogErrorf("Initialising: failed to obtain management client: %v", err)
				return err
			}
			peersResp, err := mgmClient.Peers(nc.ctx)
			if err != nil {
				nc.log.LogErrorf("Initialising: failed to obtain peers: %v", err)
				return err
			}
			peers := peersResp.Peers

			// check for at least one connected peer
			for i := 0; i < len(peers) && !connected; i++ {
				if peers[i].Connected {
					connected = true
				}
			}
		case <-nc.ctx.Done():
			// context was canceled
			nc.log.LogErrorf("Initialising: context was cancelled while waiting for L1 to be connected to other peers: %v", nc.ctx.Err())
			return ErrOperationAborted
		}
	}

	// Waiting for L1 to be healthy
	for healthy := false; !healthy; {
		nc.log.LogDebug("Initialising: waiting for L1 to become healthy...")
		select {
		case <-ticker.C:
			healthy = nc.nodeBridge.IsNodeHealthy()
		case <-nc.ctx.Done():
			// context was canceled
			nc.log.LogErrorf("Initialising: context was cancelled while waiting for L1 to become healthy: %v", nc.ctx.Err())
			return ErrOperationAborted
		}
	}

	// Subscribing to ledger updates
	commitment := nc.nodeBridge.LatestFinalizedCommitment()
	currentSlot := commitment.Commitment.Slot
	go func() {
		nc.log.LogDebugf("ListenToLedgerUpdates thread: subscribing to commitments since %v...", currentSlot)
		err := nc.nodeBridge.ListenToLedgerUpdates(nc.ctx, currentSlot, iotago.MaxSlotIndex, func(update *nodebridge.LedgerUpdate) error {
			nc.ledgerUpdatePipe.In() <- update
			return nil
		})
		if err != nil && !errors.Is(err, io.EOF) {
			nc.log.LogErrorf("ListenToLedgerUpdates thread: subscribing failed: %v", err)
			nc.shutdownHandler.SelfShutdown("Subscribing to LedgerUpdates failed", true)
		}
		nc.log.LogDebugf("ListenToLedgerUpdates thread: completed")
		/*if nc.ctx.Err() == nil {
			// shutdown in case there isn't a shutdown already in progress
			nc.shutdownHandler.SelfShutdown("INX connection closed", true)
		}*/
	}()

	return nil
}

func (nc *nodeConn) AttachChain(
	ctx context.Context,
	chainID isc.ChainID,
	chainOutputHandler ChainOutputHandler,
	otherOutputHandler OtherOutputHandler,
) ChainNodeConnection {
	cnc := newChainNodeConn(ctx, nc.log, chainID, nc.indexerClient, chainOutputHandler, otherOutputHandler)
	nc.newChainPipe.In() <- cnc
	return cnc
}

func (nc *nodeConn) chainNodeConnStopped(chainID isc.ChainID) {
	nc.stoppedChainPipe.In() <- chainID
}

func (nc *nodeConn) run() {
	//defer smT.cleanupFun()
	newChainPipeCh := nc.newChainPipe.Out()
	stoppedChainPipeCh := nc.stoppedChainPipe.Out()
	ledgerUpdatePipeCh := nc.ledgerUpdatePipe.Out()
	//timerTickCh := smT.parameters.TimeProvider.After(smT.parameters.StateManagerTimerTickPeriod)
	//statusTimerCh := smT.parameters.TimeProvider.After(constStatusTimerTime)
	for {
		if nc.ctx.Err() != nil {
			/*if smT.shutdownCoordinator == nil {
				return
			}*/
			// TODO what should the statemgr wait for?
			//smT.shutdownCoordinator.WaitNestedWithLogging(1 * time.Second)
			nc.log.LogDebugf("Stopping nodeconn, because context was closed")
			//smT.shutdownCoordinator.Done()
			return
		}
		select {
		case chainNodeConn, ok := <-newChainPipeCh:
			if ok {
				nc.handleNewChain(chainNodeConn)
			} else {
				newChainPipeCh = nil
			}
		case chainID, ok := <-stoppedChainPipeCh:
			if ok {
				nc.handleStoppedChain(chainID)
			} else {
				stoppedChainPipeCh = nil
			}
		case ledgerUpdate, ok := <-ledgerUpdatePipeCh:
			if ok {
				nc.handleLedgerUpdate(ledgerUpdate)
			} else {
				ledgerUpdatePipeCh = nil
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
		case <-nc.ctx.Done():
			continue
		}
	}
}

func (nc *nodeConn) handleNewChain(chainNodeConn ChainNodeConnection) {
	nc.chainNodeConns[chainNodeConn.GetChainID()] = chainNodeConn
}

func (nc *nodeConn) handleStoppedChain(chainID isc.ChainID) {
	delete(nc.chainNodeConns, chainID)
}

func (nc *nodeConn) handleLedgerUpdate(ledgerUpdate *nodebridge.LedgerUpdate) {
	slotIndex := ledgerUpdate.CommitmentID.Slot()
	nc.log.LogDebugf("Ledger update for slot index %v: received %v created outputs and %v consumed outputs",
		slotIndex, len(ledgerUpdate.Created), len(ledgerUpdate.Consumed))
	type outputsStruct struct {
		anchor    *isc.AnchorOutputWithID
		account   *isc.AccountOutputWithID
		basics    []*isc.BasicOutputWithID
		foundries []*isc.FoundryOutputWithID
		nfts      []*isc.NFTOutputWithID
	}
	spentOutputs := make(map[iotago.OutputID]bool)
	for _, output := range ledgerUpdate.Consumed {
		spentOutputs[output.OutputID] = true
	}

	outputMap := make(map[isc.ChainID]*outputsStruct)
	addToOutputMapFun := func(anchorID iotago.AnchorID, outputID iotago.OutputID, outputType string, addToOutputsStructFun func(*outputsStruct)) {
		chainID := isc.ChainIDFromAnchorID(anchorID)
		_, ok := nc.chainNodeConns[chainID]
		if ok {
			outputsOfChain, ok := outputMap[chainID]
			if !ok {
				outputsOfChain = &outputsStruct{
					anchor:    nil,
					account:   nil,
					basics:    make([]*isc.BasicOutputWithID, 0),
					foundries: make([]*isc.FoundryOutputWithID, 0),
					nfts:      make([]*isc.NFTOutputWithID, 0),
				}
				outputMap[chainID] = outputsOfChain
			}
			_, ok = spentOutputs[outputID]
			if ok {
				nc.log.LogDebugf("Ledger update for slot index %v: %s output %s for chain %s is skipped, because it is already consumed in this update",
					slotIndex, outputType, outputID, chainID)
			} else {
				addToOutputsStructFun(outputsOfChain)
				nc.log.LogDebugf("Ledger update for slot index %v: %s output %s will be passed to chain %s", slotIndex, outputType, outputID, chainID)
			}
		} // else { } // NOTE: outputs skipped, because they belong to untracked chain, are not logged
	}
	for _, output := range ledgerUpdate.Created {
		outputID := output.OutputID
		outputOutput := output.Output
		switch outputOutput.Type() {
		case iotago.OutputAnchor:
			var anchorID iotago.AnchorID
			anchorOutput := outputOutput.(*iotago.AnchorOutput)
			if anchorOutput.AnchorID.Empty() {
				anchorID = iotago.AnchorIDFromOutputID(outputID)
			} else {
				anchorID = anchorOutput.AnchorID
			}
			addToOutputMapFun(anchorID, outputID, "anchor", func(chainOutputs *outputsStruct) {
				if chainOutputs.anchor != nil {
					nc.log.LogPanicf("Ledger update for slot index %v: more than one unspent anchor output found: %s and %s",
						outputID, chainOutputs.anchor.OutputID())
				}
				chainOutputs.anchor = isc.NewAnchorOutputWithID(anchorOutput, outputID)
			})
		case iotago.OutputAccount:
			accountOutput := outputOutput.(*iotago.AccountOutput)
			anchorID, err := getAnchorIDFromAddress(accountOutput.Ident())
			if err != nil {
				nc.log.LogErrorf("Ledger update for slot index %v: cannot obtain anchor ID of account output %s: %v", outputID, err)
			} else {
				addToOutputMapFun(anchorID, outputID, "account", func(chainOutputs *outputsStruct) {
					if chainOutputs.account != nil {
						nc.log.LogPanicf("Ledger update for slot index %v: more than one unspent account output found: %s and %s",
							outputID, chainOutputs.account.OutputID())
					}
					chainOutputs.account = isc.NewAccountOutputWithID(accountOutput, outputID)
				})
			}
		case iotago.OutputBasic:
			basicOutput := outputOutput.(*iotago.BasicOutput)
			anchorID, err := getAnchorIDFromAddress(basicOutput.Ident())
			if err != nil {
				nc.log.LogErrorf("Ledger update for slot index %v: cannot obtain anchor ID of basic output %s: %v", outputID, err)
			} else {
				addToOutputMapFun(anchorID, outputID, "basic", func(chainOutputs *outputsStruct) {
					chainOutputs.basics = append(chainOutputs.basics, isc.NewBasicOutputWithID(basicOutput, outputID))
				})
			}
		case iotago.OutputFoundry:
			foundryOutput := outputOutput.(*iotago.FoundryOutput)
			anchorID, err := getAnchorIDFromAddress(foundryOutput.Ident())
			if err != nil {
				nc.log.LogErrorf("Ledger update for slot index %v: cannot obtain anchor ID of foundry output %s: %v", outputID, err)
			} else {
				addToOutputMapFun(anchorID, outputID, "foundry", func(chainOutputs *outputsStruct) {
					chainOutputs.foundries = append(chainOutputs.foundries, isc.NewFoundryOutputWithID(foundryOutput, outputID))
				})
			}
		case iotago.OutputNFT:
			nftOutput := outputOutput.(*iotago.NFTOutput)
			anchorID, err := getAnchorIDFromAddress(nftOutput.Ident())
			if err != nil {
				nc.log.LogErrorf("Ledger update for slot index %v: cannot obtain anchor ID of NFT output %s: %v", outputID, err)
			} else {
				addToOutputMapFun(anchorID, outputID, "NFT", func(chainOutputs *outputsStruct) {
					chainOutputs.nfts = append(chainOutputs.nfts, isc.NewNFTOutputWithID(nftOutput, outputID))
				})
			}
		default:
			nc.log.LogPanicf("Ledger update for slot index %v: unknown output type: %T", output)
		}
	}

	for chainID, outputs := range outputMap {
		chainNodeConn := nc.chainNodeConns[chainID]
		chainNodeConn.outputsReceived(
			slotIndex,
			outputs.anchor,
			outputs.account,
			outputs.basics,
			outputs.foundries,
			outputs.nfts,
		)
	}
	nc.log.LogDebugf("Ledger update for slot index %v handled", slotIndex)
}

func getAnchorIDFromAddress(address iotago.Address) (iotago.AnchorID, error) {
	anchorAddress, ok := address.(*iotago.AnchorAddress)
	if !ok {
		return iotago.AnchorID{}, fmt.Errorf("address %s is not an AnchorAddress", address)
	}
	return anchorAddress.AnchorID(), nil
}
