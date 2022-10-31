//
//
//
//
//
//

package smGPA

import (
	"fmt"
	"time"

	"github.com/iotaledger/hive.go/core/kvstore"
	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smInputs"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smMessages"
	"github.com/iotaledger/wasp/packages/chain/statemanager/smUtils"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
)

type stateManagerGPA struct {
	log                     *logger.Logger
	chainID                 *isc.ChainID
	store                   kvstore.KVStore
	blockCache              smUtils.BlockCache
	blockRequests           map[state.BlockHash]([]blockRequest) //nolint:gocritic // removing brackets doesn't make code simpler or clearer
	nodeRandomiser          smUtils.NodeRandomiser
	solidState              state.VirtualStateAccess
	solidStateBlockHash     state.BlockHash
	solidStateOutputSeq     uint32
	stateOutputSeqLastUsed  uint32
	timers                  StateManagerTimers
	lastGetBlocksTime       time.Time
	lastCleanBlockCacheTime time.Time
	lastCleanRequestsTime   time.Time
}

var _ gpa.GPA = &stateManagerGPA{}

const (
	numberOfNodesToRequestBlockFromConst = 5
)

func New(chainID *isc.ChainID, nr smUtils.NodeRandomiser, walFolder string, store kvstore.KVStore, log *logger.Logger, timersOpt ...StateManagerTimers) (gpa.GPA, error) {
	var err error
	var timers StateManagerTimers
	smLog := log.Named("gpa")
	if len(timersOpt) > 0 {
		timers = timersOpt[0]
	} else {
		timers = NewStateManagerTimers()
	}
	var wal smUtils.BlockWAL
	if walFolder == "" {
		wal = smUtils.NewMockedBlockWAL()
	} else {
		wal, err = smUtils.NewBlockWAL(walFolder, chainID, log)
		if err != nil {
			smLog.Debugf("Error creating block WAL: %v", err)
			return nil, err
		}
	}
	blockCache, err := smUtils.NewBlockCache(timers.TimeProvider, wal, smLog)
	if err != nil {
		smLog.Debugf("Error creating block cache: %v", err)
		return nil, err
	}
	result := &stateManagerGPA{
		log:                     smLog,
		chainID:                 chainID,
		store:                   store,
		blockCache:              blockCache,
		blockRequests:           make(map[state.BlockHash][]blockRequest),
		nodeRandomiser:          nr,
		solidStateOutputSeq:     0,
		solidStateBlockHash:     state.OriginBlockHash(),
		stateOutputSeqLastUsed:  0,
		timers:                  timers,
		lastGetBlocksTime:       time.Time{},
		lastCleanBlockCacheTime: time.Time{},
	}
	result.solidState, err = result.createOriginState()
	if err != nil {
		result.log.Errorf("Error creating origin state: %v", err)
		return nil, err
	}
	return result, nil
}

// -------------------------------------
// Implementation for gpa.GPA interface
// -------------------------------------

func (smT *stateManagerGPA) Input(input gpa.Input) gpa.OutMessages {
	switch inputCasted := input.(type) {
	case *smInputs.ChainBlockProduced: // From chain
		return smT.handleChainBlockProduced(inputCasted)
	case *smInputs.ChainReceiveConfirmedAliasOutput: // From chain
		return smT.handleChainReceiveConfirmedAliasOutput(inputCasted.GetStateOutput())
	case *smInputs.ConsensusStateProposal: // From consensus
		return smT.handleConsensusStateProposal(inputCasted)
	case *smInputs.ConsensusDecidedState: // From consensus
		return smT.handleConsensusDecidedState(inputCasted)
	case *smInputs.StateManagerTimerTick: // From state manager go routine
		return smT.handleStateManagerTimerTick(inputCasted.GetTime())
	default:
		smT.log.Warnf("Unknown input received, ignoring it: type=%T, message=%v", input, input)
		return gpa.NoMessages()
	}
}

func (smT *stateManagerGPA) Message(msg gpa.Message) gpa.OutMessages {
	switch msgCasted := msg.(type) {
	case *smMessages.GetBlockMessage:
		return smT.handlePeerGetBlock(msgCasted.GetSender(), msgCasted.GetBlockHash())
	case *smMessages.BlockMessage:
		return smT.handlePeerBlock(msgCasted.GetSender(), msgCasted.GetBlock())
	default:
		smT.log.Warnf("Unknown message received, ignoring it: type=%T, message=%v", msg, msg)
		return gpa.NoMessages()
	}
}

func (smT *stateManagerGPA) Output() gpa.Output {
	return nil
}

func (smT *stateManagerGPA) StatusString() string {
	return ""
}

func (smT *stateManagerGPA) UnmarshalMessage(data []byte) (gpa.Message, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("Error unmarshalling message: slice of length %v is too short", data)
	}
	switch data[0] {
	case smMessages.MsgTypeBlockMessage:
		return smMessages.NewBlockMessageFromBytes(data)
	case smMessages.MsgTypeGetBlockMessage:
		return smMessages.NewGetBlockMessageFromBytes(data)
	default:
		return nil, fmt.Errorf("Error unmarshalling message: message type %v unknown", data[0])
	}
}

// -------------------------------------
// Internal functions
// -------------------------------------

func (smT *stateManagerGPA) handlePeerGetBlock(from gpa.NodeID, blockHash state.BlockHash) gpa.OutMessages {
	smT.log.Debugf("Message received from peer %s: request to get block %s", from, blockHash)
	block := smT.blockCache.GetBlock(blockHash)
	if block == nil {
		return gpa.NoMessages()
	}
	return gpa.NoMessages().Add(smMessages.NewBlockMessage(block, from))
}

func (smT *stateManagerGPA) handlePeerBlock(from gpa.NodeID, block state.Block) gpa.OutMessages {
	smT.log.Debugf("Message received from peer %s: block %s", from, block.GetHash())
	_, ok := smT.blockRequests[block.GetHash()]
	if !ok {
		return gpa.NoMessages()
	}
	messages, _ := smT.handleGeneralBlock(block)
	return messages
}

func (smT *stateManagerGPA) handleChainBlockProduced(input *smInputs.ChainBlockProduced) gpa.OutMessages {
	smT.log.Debugf("Input received: chain block produced: block %s, alias output %s",
		input.GetBlock().GetHash(), isc.OID(input.GetAliasOutputWithID().ID()))
	messages, err := smT.handleGeneralBlock(input.GetBlock())
	input.Respond(err)
	return messages
}

func (smT *stateManagerGPA) handleGeneralBlock(block state.Block) (gpa.OutMessages, error) {
	err := smT.blockCache.AddBlock(block)
	if err != nil {
		return gpa.NoMessages(), err
	}
	blockHash := block.GetHash()
	requests, ok := smT.blockRequests[blockHash]
	if !ok {
		return gpa.NoMessages(), nil
	}
	delete(smT.blockRequests, blockHash)
	return smT.traceBlockChain(block, requests...), nil
}

func (smT *stateManagerGPA) handleChainReceiveConfirmedAliasOutput(aliasOutput *isc.AliasOutputWithID) gpa.OutMessages {
	aliasOutputID := isc.OID(aliasOutput.ID())
	smT.log.Debugf("Input received: chain confirmed alias output %s", aliasOutputID)
	smT.stateOutputSeqLastUsed++
	seq := smT.stateOutputSeqLastUsed
	smT.log.Debugf("Alias output %s is %v-th in state manager", aliasOutputID, seq)
	stateCommitment, err := state.L1CommitmentFromAliasOutput(aliasOutput.GetAliasOutput())
	if err != nil {
		smT.log.Errorf("Error retrieving state commitment from alias output %s: %v", aliasOutputID, err)
		return gpa.NoMessages()
	}
	request := newLocalStateBlockRequest(stateCommitment.BlockHash, seq, func(vs state.VirtualStateAccess) {
		smT.log.Debugf("Virtual state for alias output %s (%v-th in state manager) is ready", aliasOutputID, seq)
		if seq <= smT.solidStateOutputSeq {
			smT.log.Warnf("State for output %s (%v-th in state manager) is not needed: it is already overwritten by %v-th output",
				aliasOutputID, seq, smT.solidStateOutputSeq)
			return
		}
		//TODO: store blocks to DB
		smT.solidState = vs
		smT.solidStateBlockHash = stateCommitment.BlockHash
		// NOTE: A situation in which a request is waiting for block, which is
		// between the old solid state and the new one, should not be possible:
		// While the blocks for new solid state are being collected to handle this
		// `localStateBlockRequest`, all of them are guaranteed to be in the WAL
		// at least. Thus if some request is already waiting for the block when
		// it arrives for this `localStateBlockRequest`, both requests are handled
		// and no longer wait for this block. Otherwise, if a request for the block
		// arrives after this `localStateBlockRequest` obtains it, it is retrieved
		// from the WAL and the request is handled imediatelly without a need to
		// wait for input from other nodes.
		// If such sittuation appears possible after all, it should be handled here.
	})
	return smT.traceBlockChainByRequest(request)
}

func (smT *stateManagerGPA) handleConsensusStateProposal(csp *smInputs.ConsensusStateProposal) gpa.OutMessages {
	smT.log.Debugf("Input received: consensus state proposal for output %s", isc.OID(csp.GetAliasOutputWithID().ID()))
	request, err := newConsensusStateProposalBlockRequest(csp)
	if err != nil {
		smT.log.Errorf("Error creating consensus state proposal block request: %v", err)
		return gpa.NoMessages()
	}
	return smT.traceBlockChainByRequest(request)
}

func (smT *stateManagerGPA) handleConsensusDecidedState(cds *smInputs.ConsensusDecidedState) gpa.OutMessages {
	smT.log.Debugf("Input received: consensus request for decided state for output %s and commitment %s",
		isc.OID(cds.GetAliasOutputID().UTXOInput()), cds.GetStateCommitment())
	return smT.traceBlockChainByRequest(newConsensusDecidedStateBlockRequest(cds))
}

func (smT *stateManagerGPA) traceBlockChain(block state.Block, requests ...blockRequest) gpa.OutMessages {
	for _, request := range requests {
		request.blockAvailable(block)
	}
	nextBlockHash := block.PreviousL1Commitment().BlockHash
	currentBlock := block
	for !nextBlockHash.Equals(smT.solidStateBlockHash) && !nextBlockHash.Equals(state.OriginBlockHash()) {
		smT.log.Debugf("Tracing the chain of blocks: %s is not the first one, looking for its parent %s", currentBlock.GetHash(), nextBlockHash)
		var response gpa.OutMessages
		currentBlock, response = smT.getBlockOrRequestMessages(nextBlockHash, requests...)
		if currentBlock == nil {
			return response
		}
		for _, request := range requests {
			request.blockAvailable(currentBlock)
		}
		nextBlockHash = currentBlock.PreviousL1Commitment().BlockHash
	}
	var stateType string
	var createBaseStateFun createStateFun
	if !nextBlockHash.Equals(state.OriginBlockHash()) {
		stateType = "solid"
		createBaseStateFun = func() (state.VirtualStateAccess, error) {
			// NOTE: an assumption is that `createBaseStateFun` is called in
			// the same function call as it is created in. This is needed in order
			// for `smT.solidState` to be exactly the same as it was during
			// the creation of this fun.
			return smT.solidState.Copy(), nil
		}
	} else {
		stateType = "origin"
		createBaseStateFun = smT.createOriginState
	}
	smT.log.Debugf("Tracing the chain of blocks: the chain is complete, marking all the requests as completed based on %s state", stateType)
	// Competing all the top priority requests (the ones that do not chagne the
	// state of manager state) and the one with the largest priority among the others.
	// Other requests are just marked as completed without any call to respond function.
	// The idea is that all the consensus requests must be completed and only
	// a single request which was created after receiving alias output from L1:
	// the one for the newest alias output. Moreover, consensus requests must be
	// processed before the alias output one to use the old solid state, which is
	// modified by alias output request.
	maxPriority := uint32(0)
	delayedRequests := make([]blockRequest, 0)
	for _, request := range requests {
		priority := request.getPriority()
		if priority == topPriority {
			request.markCompleted(createBaseStateFun)
		} else {
			if maxPriority == 0 {
				maxPriority = priority
				delayedRequests = append(delayedRequests, request)
			} else if priority > maxPriority {
				maxPriority = priority
				delayedRequests = append(delayedRequests, delayedRequests[0])
				delayedRequests[0] = request
			} else {
				delayedRequests = append(delayedRequests, request)
			}
		}
	}
	if len(delayedRequests) > 0 {
		delayedRequests[0].markCompleted(createBaseStateFun)
		for _, request := range delayedRequests[1:] {
			request.markCompleted(func() (state.VirtualStateAccess, error) { return nil, nil })
		}
	}
	return gpa.NoMessages()
}

func (smT *stateManagerGPA) traceBlockChainByRequest(request blockRequest) gpa.OutMessages {
	lastBlockHash := request.getLastBlockHash()
	smT.log.Debugf("Tracing the chain of blocks ending with block %s", lastBlockHash)
	block, response := smT.getBlockOrRequestMessages(lastBlockHash, request)
	if block == nil {
		return response
	}
	return smT.traceBlockChain(block, request)
}

func (smT *stateManagerGPA) getBlockOrRequestMessages(blockHash state.BlockHash, requests ...blockRequest) (state.Block, gpa.OutMessages) {
	block := smT.blockCache.GetBlock(blockHash)
	if block == nil {
		// Mark that the requests are waiting for `blockHash` block
		currrentRequests, ok := smT.blockRequests[blockHash]
		if !ok {
			smT.log.Debugf("Block %s is missing, it is the first request waiting for it", blockHash)
			smT.blockRequests[blockHash] = requests
		} else {
			smT.log.Debugf("Block %s is missing, %v requests are waiting for it in addition to this one", len(currrentRequests))
			smT.blockRequests[blockHash] = append(currrentRequests, requests...)
		}
		return nil, smT.makeGetBlockRequestMessages(blockHash)
	}
	smT.log.Debugf("Block %s is available", blockHash)
	return block, nil // Second parameter should not be used then
}

// Make `numberOfNodesToRequestBlockFromConst` messages to random peers
func (smT *stateManagerGPA) makeGetBlockRequestMessages(blockHash state.BlockHash) gpa.OutMessages {
	nodeIDs := smT.nodeRandomiser.GetRandomOtherNodeIDs(numberOfNodesToRequestBlockFromConst)
	smT.log.Debugf("Requesting block %s from %v random nodes %v", blockHash, numberOfNodesToRequestBlockFromConst, nodeIDs)
	response := gpa.NoMessages()
	for _, nodeID := range nodeIDs {
		response.Add(smMessages.NewGetBlockMessage(blockHash, nodeID))
	}
	return response
}

func (smT *stateManagerGPA) handleStateManagerTimerTick(now time.Time) gpa.OutMessages {
	result := gpa.NoMessages()
	smT.log.Debugf("Input received: timer tick %v", now)
	if now.After(smT.lastGetBlocksTime.Add(smT.timers.StateManagerGetBlockRetry)) {
		smT.log.Debugf("Timer tick: resending get block messages...")
		for blockHash, _ := range smT.blockRequests { //nolint:gofumpt,gofmt,revive,gosimple
			result.AddAll(smT.makeGetBlockRequestMessages(blockHash))
		}
		smT.lastGetBlocksTime = now
	}
	if now.After(smT.lastCleanBlockCacheTime.Add(smT.timers.BlockCacheBlockCleaningPeriod)) {
		smT.log.Debugf("Timer tick: cleaning block cache...")
		smT.blockCache.CleanOlderThan(now.Add(-smT.timers.BlockCacheBlocksInCacheDuration))
		smT.lastCleanBlockCacheTime = now
	}
	if now.After(smT.lastCleanRequestsTime.Add(smT.timers.StateManagerRequestCleaningPeriod)) {
		smT.log.Debugf("Timer tick: cleaning requests...")
		newBlockRequestsMap := make(map[state.BlockHash]([]blockRequest)) //nolint:gocritic
		for blockHash, blockRequests := range smT.blockRequests {
			outI := 0
			for _, blockRequest := range blockRequests {
				if blockRequest.isValid() { // Request is valid - keeping it
					blockRequests[outI] = blockRequest
					outI++
				}
			}
			for i := outI; i < len(blockRequests); i++ {
				blockRequests[i] = nil // Not needed requests at the end - freeing memory
			}
			blockRequests = blockRequests[:outI]
			if len(blockRequests) > 0 {
				newBlockRequestsMap[blockHash] = blockRequests
			}
		}
		smT.blockRequests = newBlockRequestsMap
		smT.lastCleanRequestsTime = now
	}
	smT.log.Debugf("Timer tick %v handled", now)
	return result
}

func (smT *stateManagerGPA) createOriginState() (state.VirtualStateAccess, error) {
	return state.CreateOriginState(smT.store, smT.chainID)
}
