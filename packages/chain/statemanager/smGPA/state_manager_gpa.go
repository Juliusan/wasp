//
//
//
//
//
//

package smGPA

import (
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
	log            *logger.Logger
	chainID        *isc.ChainID
	store          kvstore.KVStore
	blockCache     *smUtils.BlockCache
	blockRequests  map[state.BlockHash]([]blockRequest) //nolint:gocritic // removing brackets doesn't make code simpler or clearer
	nodeRandomiser *smUtils.NodeRandomiser
}

var _ gpa.GPA = &stateManagerGPA{}

const (
	numberOfNodesToRequestBlockFromConst = 5
)

func New(chainID *isc.ChainID, nr *smUtils.NodeRandomiser, store kvstore.KVStore, log *logger.Logger) gpa.GPA {
	smLog := log.Named("sm")
	return &stateManagerGPA{
		log:            smLog,
		chainID:        chainID,
		store:          store,
		blockCache:     smUtils.NewBlockCache(smLog),
		blockRequests:  make(map[state.BlockHash][]blockRequest),
		nodeRandomiser: nr,
	}
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
	return nil, nil
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
	// TODO: aliasOutput!
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
	smT.log.Debugf("Input received: chain confirmed alias output  %s", isc.OID(aliasOutput.ID()))
	// TODO: aliasOutput!
	return gpa.NoMessages()
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
	return smT.traceBlockChainByRequest(newConsensusDecidedStateBlockRequest(cds, smT.createOriginState))
}

func (smT *stateManagerGPA) traceBlockChain(block state.Block, requests ...blockRequest) gpa.OutMessages {
	for _, request := range requests {
		request.blockAvailable(block)
	}
	nextBlockHash := block.PreviousL1Commitment().BlockHash
	currentBlock := block
	for !nextBlockHash.Equals(state.OriginBlockHash()) {
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
	smT.log.Debugf("Tracing the chain of blocks: the chain is complete, marking all the requests as completed")
	for _, request := range requests {
		request.markCompleted()
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
		// Make `numberOfNodesToRequestBlockFromConst` messages to random peers
		nodeIDs := smT.nodeRandomiser.GetRandomOtherNodeIDs(numberOfNodesToRequestBlockFromConst)
		smT.log.Debugf("Requesting block %s from %v random nodes %v", blockHash, numberOfNodesToRequestBlockFromConst, nodeIDs)
		response := gpa.NoMessages()
		for _, nodeID := range nodeIDs {
			response.Add(smMessages.NewGetBlockMessage(blockHash, nodeID))
		}
		return nil, response
	}
	smT.log.Debugf("Block %s is available", blockHash)
	return block, nil // Second parameter should not be used then
}

func (smT *stateManagerGPA) createGetBlockMessages(blockHash state.BlockHash) gpa.OutMessages {
	// pick some number of messages to request block `blockHash` from random peers
	return gpa.NoMessages()
}

func (smT *stateManagerGPA) createOriginState() (state.VirtualStateAccess, error) {
	return state.CreateOriginState(smT.store, smT.chainID)
}
