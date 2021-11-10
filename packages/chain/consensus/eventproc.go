// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package consensus

import (
	"fmt"
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/chain/messages"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/state"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

func (c *Consensus) EventStateTransitionMsg(state state.VirtualStateAccess, stateOutput *ledgerstate.AliasOutput, stateTimestamp time.Time) {
	c.eventStateTransitionMsgCh <- &messages.StateTransitionMsg{
		State:          state,
		StateOutput:    stateOutput,
		StateTimestamp: stateTimestamp,
	}
}

func (c *Consensus) eventStateTransitionMsg(msg *messages.StateTransitionMsg) {
	c.log.Debugf("StateTransitionMsg received: state index: %d, state output: %s, timestamp: %v",
		msg.State.BlockIndex(), iscp.OID(msg.StateOutput.ID()), msg.StateTimestamp)
	c.setNewState(msg)
	c.takeAction()
}

func (c *Consensus) EnqueueSignedResult(chainInputID ledgerstate.OutputID, essenceHash hashing.HashValue, sigShare tbls.SigShare, senderIndex uint16) {
	c.eventSignedResultMsgCh <- &signedResultMsgIn{
		signedResultMsg: signedResultMsg{
			ChainInputID: chainInputID,
			EssenceHash:  essenceHash,
			SigShare:     sigShare,
		},
		SenderIndex: senderIndex,
	}
}

func (c *Consensus) handleSignedResult(msg *signedResultMsgIn) {
	c.log.Debugf("handleSignedResult message received: from sender %d, hash=%s, chain input id=%v",
		msg.SenderIndex, msg.EssenceHash, iscp.OID(msg.ChainInputID))
	c.receiveSignedResult(msg)
	c.takeAction()
}

func (c *Consensus) EnqueueSignedResultAck(chainInputID ledgerstate.OutputID, essenceHash hashing.HashValue, senderIndex uint16) {
	c.eventSignedResultAckMsgCh <- &signedResultAckMsgIn{
		signedResultAckMsg: signedResultAckMsg{
			ChainInputID: chainInputID,
			EssenceHash:  essenceHash,
		},
		SenderIndex: senderIndex,
	}
}

func (c *Consensus) handleSignedResultAck(msg *signedResultAckMsgIn) {
	c.log.Debugf("handleSignedResultAck message received: from sender %d, hash=%s, chain input id=%v",
		msg.SenderIndex, msg.EssenceHash, iscp.OID(msg.ChainInputID))
	c.receiveSignedResultAck(msg.ChainInputID, msg.EssenceHash, msg.SenderIndex)
	c.takeAction()
}

func (c *Consensus) EventInclusionsStateMsg(txID ledgerstate.TransactionID, state ledgerstate.InclusionState) {
	c.eventInclusionStateMsgCh <- &messages.InclusionStateMsg{
		TxID:  txID,
		State: state,
	}
}

func (c *Consensus) eventInclusionState(msg *messages.InclusionStateMsg) {
	c.log.Debugf("InclusionStateMsg received:  %s: '%s'", msg.TxID.Base58(), msg.State.String())
	c.processInclusionState(msg)

	c.takeAction()
}

func (c *Consensus) EventAsynchronousCommonSubsetMsg(msg *messages.AsynchronousCommonSubsetMsg) {
	c.eventACSMsgCh <- msg
}

func (c *Consensus) eventAsynchronousCommonSubset(msg *messages.AsynchronousCommonSubsetMsg) {
	c.log.Debugf("AsynchronousCommonSubsetMsg received for session %v: len = %d", msg.SessionID, len(msg.ProposedBatchesBin))
	c.receiveACS(msg.ProposedBatchesBin, msg.SessionID)

	c.takeAction()
}

func (c *Consensus) EventVMResultMsg(msg *messages.VMResultMsg) {
	c.eventVMResultMsgCh <- msg
}

func (c *Consensus) eventVMResultMsg(msg *messages.VMResultMsg) {
	var essenceString string
	if msg.Task.ResultTransactionEssence == nil {
		essenceString = "essence is nil"
	} else {
		essenceString = fmt.Sprintf("essence hash: %s", hashing.HashData(msg.Task.ResultTransactionEssence.Bytes()))
	}
	c.log.Debugf("VMResultMsg received: state index: %d state hash: %s %s",
		msg.Task.VirtualStateAccess.BlockIndex(), msg.Task.VirtualStateAccess.StateCommitment(), essenceString)
	c.processVMResult(msg.Task)
	c.takeAction()
}

func (c *Consensus) EventTimerMsg(msg messages.TimerTick) {
	c.eventTimerMsgCh <- msg
}

func (c *Consensus) eventTimerMsg(msg messages.TimerTick) {
	c.lastTimerTick.Store(int64(msg))
	c.refreshConsensusInfo()
	if msg%40 == 0 {
		if snap := c.GetStatusSnapshot(); snap != nil {
			c.log.Infof("timer tick #%d: state index: %d, mempool = (%d, %d, %d)",
				snap.TimerTick, snap.StateIndex, snap.Mempool.InPoolCounter, snap.Mempool.OutPoolCounter, snap.Mempool.TotalPool)
		}
	}
	c.takeAction()
	if c.stateOutput != nil {
		c.log.Debugf("Consensus::eventTimerMsg: stateIndex=%v, workflow=%+v",
			c.stateOutput.GetStateIndex(),
			c.workflow,
		)
	} else {
		c.log.Debugf("Consensus::eventTimerMsg: stateIndex=nil, workflow=%+v",
			c.workflow,
		)
	}
}
