// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package statemgr

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/chain/messages"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/state"
)

// EventGetBlockMsg is a request for a block while syncing
func (sm *stateManager) EnqueueGetBlock(index uint32, senderNetID string) {
	sm.eventGetBlockMsgCh <- &getBlockMsgIn{
		getBlockMsg: getBlockMsg{BlockIndex: index},
		SenderNetID: senderNetID,
	}
}

func (sm *stateManager) handleGetBlock(msg *getBlockMsgIn) {
	sm.log.Debugw("handleGetBlock: ",
		"sender", msg.SenderNetID,
		"block index", msg.BlockIndex,
	)
	if sm.stateOutput == nil { // Not a necessary check, only for optimization.
		sm.log.Debugf("handleGetBlock message ignored: stateOutput is nil")
		return
	}
	if msg.BlockIndex > sm.stateOutput.GetStateIndex() { // Not a necessary check, only for optimization.
		sm.log.Debugf("handleGetBlock message ignored 1: block #%d not found. Current state index: #%d",
			msg.BlockIndex, sm.stateOutput.GetStateIndex())
		return
	}
	blockBytes, err := state.LoadBlockBytes(sm.store, msg.BlockIndex)
	if err != nil {
		sm.log.Errorf("handleGetBlock: LoadBlockBytes: %v", err)
		return
	}
	if blockBytes == nil {
		sm.log.Debugf("handleGetBlock message ignored 2: block #%d not found. Current state index: #%d",
			msg.BlockIndex, sm.stateOutput.GetStateIndex())
		return
	}

	sm.log.Debugf("handleGetBlock for state index #%d --> responding to peer %s", msg.BlockIndex, msg.SenderNetID)

	sm.chain.SendMsgByNetID(msg.SenderNetID, peering.PeerMessagePartyStateManager, &blockMsg{
		BlockBytes: blockBytes,
	})
}

// EventBlockMsg
func (sm *stateManager) EnqueueBlock(block []byte, senderNetID string) {
	sm.eventBlockMsgCh <- &blockMsgIn{
		blockMsg:    blockMsg{BlockBytes: block},
		SenderNetID: senderNetID,
	}
}

func (sm *stateManager) handleBlock(msg *blockMsgIn) {
	sm.log.Debugf("handleBlock message received from %v", msg.SenderNetID)
	if sm.stateOutput == nil {
		sm.log.Debugf("handleBlock message ignored: stateOutput is nil")
		return
	}
	block, err := state.BlockFromBytes(msg.BlockBytes)
	if err != nil {
		sm.log.Warnf("handleBlock message ignored: wrong block received from peer %s. Err: %v", msg.SenderNetID, err)
		return
	}
	sm.log.Debugw("handleBlock from ",
		"sender", msg.SenderNetID,
		"block index", block.BlockIndex(),
		"approving output", iscp.OID(block.ApprovingOutputID()),
	)
	if sm.addBlockFromPeer(block) {
		sm.takeAction()
	}
}

func (sm *stateManager) EventOutputMsg(msg ledgerstate.Output) {
	sm.eventOutputMsgCh <- msg
}

func (sm *stateManager) eventOutputMsg(msg ledgerstate.Output) {
	sm.log.Debugf("EventOutputMsg received: %s", iscp.OID(msg.ID()))
	chainOutput, ok := msg.(*ledgerstate.AliasOutput)
	if !ok {
		sm.log.Debugf("EventOutputMsg ignored: output is of type %t, expecting *ledgerstate.AliasOutput", msg)
		return
	}
	if sm.outputPulled(chainOutput) {
		sm.takeAction()
	}
}

// EventStateTransactionMsg triggered whenever new state transaction arrives
// the state transaction may be confirmed or not
func (sm *stateManager) EventStateMsg(msg *messages.StateMsg) {
	sm.eventStateOutputMsgCh <- msg
}

func (sm *stateManager) eventStateMsg(msg *messages.StateMsg) {
	sm.log.Debugw("EventStateMsg received: ",
		"state index", msg.ChainOutput.GetStateIndex(),
		"chainOutput", iscp.OID(msg.ChainOutput.ID()),
	)
	stateHash, err := hashing.HashValueFromBytes(msg.ChainOutput.GetStateData())
	if err != nil {
		sm.log.Errorf("EventStateMsg ignored: failed to parse state hash: %v", err)
		return
	}
	sm.log.Debugf("EventStateMsg state hash is %v", stateHash.String())
	if sm.stateOutputReceived(msg.ChainOutput, msg.Timestamp) {
		sm.takeAction()
	}
}

func (sm *stateManager) EventStateCandidateMsg(state state.VirtualStateAccess, outputID ledgerstate.OutputID) {
	sm.eventStateCandidateMsgCh <- &messages.StateCandidateMsg{
		State:             state,
		ApprovingOutputID: outputID,
	}
}

func (sm *stateManager) eventStateCandidateMsg(msg *messages.StateCandidateMsg) {
	sm.log.Debugf("EventStateCandidateMsg received: state index: %d, timestamp: %v",
		msg.State.BlockIndex(), msg.State.Timestamp(),
	)
	if sm.stateOutput == nil {
		sm.log.Debugf("EventStateCandidateMsg ignored: stateOutput is nil")
		return
	}
	if sm.addStateCandidateFromConsensus(msg.State, msg.ApprovingOutputID) {
		sm.takeAction()
	}
}

func (sm *stateManager) EventTimerMsg(msg messages.TimerTick) {
	if msg%2 == 0 {
		sm.eventTimerMsgCh <- msg
	}
}

func (sm *stateManager) eventTimerMsg() {
	sm.log.Debugf("EventTimerMsg received")
	sm.takeAction()
}
