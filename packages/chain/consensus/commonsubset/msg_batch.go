// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // TODO there is a bunch of duplicated code in this file, should be refactored to reusable funcs
package commonsubset

import (
	"bytes"

	"github.com/anthdm/hbbft"
	"github.com/iotaledger/wasp/packages/chain/consensus/commoncoin"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
	"golang.org/x/xerrors"
)

/** msgBatch groups messages generated at one step in the protocol for a single recipient. */
type msgBatch struct {
	sessionID  uint64
	stateIndex uint32
	id         uint32              // ID of the batch for the acks, ID=0 => Acks are not needed.
	src        uint16              // Sender of the batch.
	dst        uint16              // Recipient of the batch.
	msgs       []*hbbft.ACSMessage // New messages to send.
	acks       []uint32            // Acknowledgements.
}

const (
	acsMsgTypeRbcProofRequest byte = 1 << 4
	acsMsgTypeRbcEchoRequest  byte = 2 << 4
	acsMsgTypeRbcReadyRequest byte = 3 << 4
	acsMsgTypeAbaBvalRequest  byte = 4 << 4
	acsMsgTypeAbaAuxRequest   byte = 5 << 4
	acsMsgTypeAbaCCRequest    byte = 6 << 4
	acsMsgTypeAbaDoneRequest  byte = 7 << 4
)

var _ peering.Serializable = &msgBatch{}

func (msg *msgBatch) GetMsgType() byte {
	return byte(acsMsgType)
}

func (msg *msgBatch) GetDeserializer() peering.PeerMessagePartyType {
	return peering.PeerMessagePartyAcs
}

//nolint:funlen, gocyclo // TODO this function is too long and has a high cyclomatic complexity, should be refactored
func (msg *msgBatch) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = util.WriteUint64(&buf, msg.sessionID); err != nil {
		return nil, xerrors.Errorf("failed to write msgBatch.sessionID: %w", err)
	}
	if err = util.WriteUint32(&buf, msg.stateIndex); err != nil {
		return nil, xerrors.Errorf("failed to write msgBatch.stateIndex: %w", err)
	}
	if err = util.WriteUint32(&buf, msg.id); err != nil {
		return nil, xerrors.Errorf("failed to write msgBatch.id: %w", err)
	}
	if err = util.WriteUint16(&buf, msg.src); err != nil {
		return nil, xerrors.Errorf("failed to write msgBatch.src: %w", err)
	}
	if err = util.WriteUint16(&buf, msg.dst); err != nil {
		return nil, xerrors.Errorf("failed to write msgBatch.dst: %w", err)
	}
	if err = util.WriteUint16(&buf, uint16(len(msg.msgs))); err != nil {
		return nil, xerrors.Errorf("failed to write msgBatch.msgs.len: %w", err)
	}
	for i, acsMsg := range msg.msgs {
		if err = util.WriteUint16(&buf, uint16(acsMsg.ProposerID)); err != nil {
			return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].ProposerID: %w", i, err)
		}
		switch acsMsgPayload := acsMsg.Payload.(type) {
		case *hbbft.BroadcastMessage:
			switch rbcMsgPayload := acsMsgPayload.Payload.(type) {
			case *hbbft.ProofRequest:
				if err = util.WriteByte(&buf, acsMsgTypeRbcProofRequest); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].type: %w", i, err)
				}
				if err = util.WriteBytes16(&buf, rbcMsgPayload.RootHash); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].RootHash: %w", i, err)
				}
				if err = util.WriteUint32(&buf, uint32(len(rbcMsgPayload.Proof))); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Proof.len: %w", i, err)
				}
				for pi, p := range rbcMsgPayload.Proof {
					if err = util.WriteBytes32(&buf, p); err != nil {
						return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Proof[%v]: %w", i, pi, err)
					}
				}
				if err = util.WriteUint16(&buf, uint16(rbcMsgPayload.Index)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Index: %w", i, err)
				}
				if err = util.WriteUint16(&buf, uint16(rbcMsgPayload.Leaves)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Leaves: %w", i, err)
				}
			case *hbbft.EchoRequest:
				if err = util.WriteByte(&buf, acsMsgTypeRbcEchoRequest); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].type: %w", i, err)
				}
				if err = util.WriteBytes16(&buf, rbcMsgPayload.RootHash); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].RootHash: %w", i, err)
				}
				if err = util.WriteUint32(&buf, uint32(len(rbcMsgPayload.Proof))); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Proof.len: %w", i, err)
				}
				for pi, p := range rbcMsgPayload.Proof {
					if err = util.WriteBytes32(&buf, p); err != nil {
						return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Proof[%v]: %w", i, pi, err)
					}
				}
				if err = util.WriteUint16(&buf, uint16(rbcMsgPayload.Index)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Index: %w", i, err)
				}
				if err = util.WriteUint16(&buf, uint16(rbcMsgPayload.Leaves)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Leaves: %w", i, err)
				}
			case *hbbft.ReadyRequest:
				if err = util.WriteByte(&buf, acsMsgTypeRbcReadyRequest); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].type: %w", i, err)
				}
				if err = util.WriteBytes16(&buf, rbcMsgPayload.RootHash); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].RootHash: %w", i, err)
				}
			default:
				return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v]: unexpected broadcast message type", i)
			}
		case *hbbft.AgreementMessage:
			switch abaMsgPayload := acsMsgPayload.Message.(type) {
			case *hbbft.BvalRequest:
				encoded := acsMsgTypeAbaBvalRequest
				if abaMsgPayload.Value {
					encoded |= 1
				}
				if err = util.WriteByte(&buf, encoded); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].value: %w", i, err)
				}
				if err = util.WriteUint16(&buf, uint16(acsMsgPayload.Epoch)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].epoch: %w", i, err)
				}
			case *hbbft.AuxRequest:
				encoded := acsMsgTypeAbaAuxRequest
				if abaMsgPayload.Value {
					encoded |= 1
				}
				if err = util.WriteByte(&buf, encoded); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].value: %w", i, err)
				}
				if err = util.WriteUint16(&buf, uint16(acsMsgPayload.Epoch)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].epoch: %w", i, err)
				}
			case *hbbft.CCRequest:
				coinMsg := abaMsgPayload.Payload.(*commoncoin.BlsCommonCoinMsg)
				if err = util.WriteByte(&buf, acsMsgTypeAbaCCRequest); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].type: %w", i, err)
				}
				if err = util.WriteUint16(&buf, uint16(acsMsgPayload.Epoch)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].epoch: %w", i, err)
				}
				if err = coinMsg.Write(&buf); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].Payload: %w", i, err)
				}
			case *hbbft.DoneRequest:
				encoded := acsMsgTypeAbaDoneRequest
				if err = util.WriteByte(&buf, encoded); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].type: %w", i, err)
				}
				if err = util.WriteUint16(&buf, uint16(acsMsgPayload.Epoch)); err != nil {
					return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v].epoch: %w", i, err)
				}
			default:
				return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v]: unexpected agreemet message type", i)
			}
		default:
			return nil, xerrors.Errorf("failed to write msgBatch.msgs[%v]: unexpected acs message type", i)
		}
	}
	if err = util.WriteUint16(&buf, uint16(len(msg.acks))); err != nil {
		return nil, xerrors.Errorf("failed to write msgBatch.acks.len: %w", err)
	}
	for i, ack := range msg.acks {
		if err = util.WriteUint32(&buf, ack); err != nil {
			return nil, xerrors.Errorf("failed to write msgBatch.acks[%v]: %w", i, err)
		}
	}
	return buf.Bytes(), nil
}

//nolint:funlen, gocyclo // TODO this function is too long and has a high cyclomatic complexity, should be refactored
func (msg *msgBatch) Deserialize(buf []byte) error {
	var err error
	r := bytes.NewBuffer(buf)
	if err = util.ReadUint64(r, &msg.sessionID); err != nil {
		return xerrors.Errorf("failed to read msgBatch.sessionID: %w", err)
	}
	if err = util.ReadUint32(r, &msg.stateIndex); err != nil {
		return xerrors.Errorf("failed to read msgBatch.stateIndex: %w", err)
	}
	if err = util.ReadUint32(r, &msg.id); err != nil {
		return xerrors.Errorf("failed to read msgBatch.id: %w", err)
	}
	if err = util.ReadUint16(r, &msg.src); err != nil {
		return xerrors.Errorf("failed to read msgBatch.src: %w", err)
	}
	if err = util.ReadUint16(r, &msg.dst); err != nil {
		return xerrors.Errorf("failed to read msgBatch.dst: %w", err)
	}
	//
	// Msgs.
	var msgsLen uint16
	if err = util.ReadUint16(r, &msgsLen); err != nil {
		return xerrors.Errorf("failed to read msgBatch.msgs.len: %w", err)
	}
	msg.msgs = make([]*hbbft.ACSMessage, msgsLen)
	for mi := range msg.msgs {
		acsMsg := hbbft.ACSMessage{}
		msg.msgs[mi] = &acsMsg
		// ProposerID.
		var proposerID uint16
		if err = util.ReadUint16(r, &proposerID); err != nil {
			return xerrors.Errorf("failed to read msgBatch.msgs[%v].ProposerID: %w", mi, err)
		}
		acsMsg.ProposerID = uint64(proposerID)
		// By Type.
		var msgType byte
		if msgType, err = util.ReadByte(r); err != nil {
			return xerrors.Errorf("failed to read msgBatch.msgs[%v].type: %w", mi, err)
		}
		switch msgType & 0xF0 {
		case acsMsgTypeRbcProofRequest:
			proofRequest := hbbft.ProofRequest{}
			if proofRequest.RootHash, err = util.ReadBytes16(r); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].RootHash: %w", mi, err)
			}
			var proofLen uint32
			if err = util.ReadUint32(r, &proofLen); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].Proof.len: %w", mi, err)
			}
			proofRequest.Proof = make([][]byte, proofLen)
			for pi := range proofRequest.Proof {
				if proofRequest.Proof[pi], err = util.ReadBytes32(r); err != nil {
					return xerrors.Errorf("failed to read msgBatch.msgs[%v].Proof[%v]: %w", mi, pi, err)
				}
			}
			var proofRequestIndex uint16
			if err = util.ReadUint16(r, &proofRequestIndex); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].Index: %w", mi, err)
			}
			proofRequest.Index = int(proofRequestIndex)
			var proofRequestLeaves uint16
			if err = util.ReadUint16(r, &proofRequestLeaves); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].Leaves: %w", mi, err)
			}
			proofRequest.Leaves = int(proofRequestLeaves)
			acsMsg.Payload = &hbbft.BroadcastMessage{
				Payload: &proofRequest,
			}
		case acsMsgTypeRbcEchoRequest:
			echoRequest := hbbft.EchoRequest{}
			if echoRequest.RootHash, err = util.ReadBytes16(r); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].RootHash: %w", mi, err)
			}
			var proofLen uint32
			if err = util.ReadUint32(r, &proofLen); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].Proof.len: %w", mi, err)
			}
			echoRequest.Proof = make([][]byte, proofLen)
			for pi := range echoRequest.Proof {
				if echoRequest.Proof[pi], err = util.ReadBytes32(r); err != nil {
					return xerrors.Errorf("failed to read msgBatch.msgs[%v].Proof[%v]: %w", mi, pi, err)
				}
			}
			var echoRequestIndex uint16
			if err = util.ReadUint16(r, &echoRequestIndex); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].Index: %w", mi, err)
			}
			echoRequest.Index = int(echoRequestIndex)
			var echoRequestLeaves uint16
			if err = util.ReadUint16(r, &echoRequestLeaves); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].Leaves: %w", mi, err)
			}
			echoRequest.Leaves = int(echoRequestLeaves)
			acsMsg.Payload = &hbbft.BroadcastMessage{
				Payload: &echoRequest,
			}
		case acsMsgTypeRbcReadyRequest:
			readyRequest := hbbft.ReadyRequest{}
			if readyRequest.RootHash, err = util.ReadBytes16(r); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].RootHash: %w", mi, err)
			}
			acsMsg.Payload = &hbbft.BroadcastMessage{
				Payload: &readyRequest,
			}
		case acsMsgTypeAbaBvalRequest:
			var epoch uint16
			if err = util.ReadUint16(r, &epoch); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].epoch: %w", mi, err)
			}
			acsMsg.Payload = &hbbft.AgreementMessage{
				Epoch:   int(epoch),
				Message: &hbbft.BvalRequest{Value: msgType&0x01 == 1},
			}
		case acsMsgTypeAbaAuxRequest:
			var epoch uint16
			if err = util.ReadUint16(r, &epoch); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].epoch: %w", mi, err)
			}
			acsMsg.Payload = &hbbft.AgreementMessage{
				Epoch:   int(epoch),
				Message: &hbbft.AuxRequest{Value: msgType&0x01 == 1},
			}
		case acsMsgTypeAbaCCRequest:
			var epoch uint16
			if err = util.ReadUint16(r, &epoch); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].epoch: %w", mi, err)
			}
			var ccReq commoncoin.BlsCommonCoinMsg
			if err = ccReq.Read(r); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].Payload: %w", mi, err)
			}
			acsMsg.Payload = &hbbft.AgreementMessage{
				Epoch: int(epoch),
				Message: &hbbft.CCRequest{
					Payload: &ccReq,
				},
			}
		case acsMsgTypeAbaDoneRequest:
			var epoch uint16
			if err = util.ReadUint16(r, &epoch); err != nil {
				return xerrors.Errorf("failed to read msgBatch.msgs[%v].epoch: %w", mi, err)
			}
			acsMsg.Payload = &hbbft.AgreementMessage{
				Epoch:   int(epoch),
				Message: &hbbft.DoneRequest{},
			}
		default:
			return xerrors.Errorf("failed to read msgBatch.msgs[%v]: unexpected message type %v, b=%+v", mi, msgType, msg)
		}
	}
	//
	// Acks.
	var acksLen uint16
	if err = util.ReadUint16(r, &acksLen); err != nil {
		return xerrors.Errorf("failed to read msgBatch.acks.len: %w", err)
	}
	msg.acks = make([]uint32, acksLen)
	for ai := range msg.acks {
		if err = util.ReadUint32(r, &msg.acks[ai]); err != nil {
			return xerrors.Errorf("failed to read msgBatch.acks[%v]: %w", ai, err)
		}
	}
	return nil
}

func (msg *msgBatch) NeedsAck() bool {
	return msg.id != 0
}
