// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// Package peering provides an overlay network for communicating
// between nodes in a peer-to-peer style with low overhead
// encoding and persistent connections. The network provides only
// the asynchronous communication.
//
// It is intended to use for the committee consensus protocol.
//
package peering

import (
	"strconv"
	"strings"
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"golang.org/x/xerrors"
)

const (
	MsgTypeReserved  = byte(0)
	MsgTypeHandshake = byte(1)
	MsgTypeMsgChunk  = byte(2)

	// FirstUserMsgCode is the first committee message type.
	// All the equal and larger msg types are committee messages.
	// those with smaller are reserved by the package for heartbeat and handshake messages
	FirstUserMsgCode = byte(0x10)
)

// PeeringID is relates peers in different nodes for a particular
// communication group. E.g. PeeringID identifies a committee in
// the consensus, etc.
type PeeringID [ledgerstate.AddressLength]byte

// NetworkProvider stands for the peer-to-peer network, as seen
// from the viewpoint of a single participant.
type NetworkProvider interface {
	Run(stopCh <-chan struct{})
	Self() PeerSender
	PeerGroup(peerAddrs []string) (GroupProvider, error)
	PeerDomain(peerAddrs []string) (PeerDomainProvider, error)
	PeerByNetID(peerNetID string) (PeerSender, error)
	PeerByPubKey(peerPub *ed25519.PublicKey) (PeerSender, error)
	PeerStatus() []PeerStatusProvider
	PeerMessagePartyCollection
}

// TrustedNetworkManager is used maintain a configuration which peers are trusted.
// In a typical implementation, this interface should be implemented by the same
// struct, that implements the NetworkProvider. These implementations should interact,
// e.g. when we distrust some peer, all the connections to it should be cut immediately.
type TrustedNetworkManager interface {
	IsTrustedPeer(pubKey ed25519.PublicKey) error
	TrustPeer(pubKey ed25519.PublicKey, netID string) (*TrustedPeer, error)
	DistrustPeer(pubKey ed25519.PublicKey) (*TrustedPeer, error)
	TrustedPeers() ([]*TrustedPeer, error)
}

// TrustedPeer carries a peer information we use to trust it.
type TrustedPeer struct {
	PubKey ed25519.PublicKey
	NetID  string
}

/*type PeerCollection interface {
	RegisterPeerMessageParty(peeringID PeeringID, party PeerMessageParty) error
	UnregisterPeerMessageParty(peeringID PeeringID, partyType PeerMessagePartyType) error
}*/

type PeerMessagePartyCollection interface {
	RegisterPeerMessageParty(peeringID PeeringID, party PeerMessageSimpleParty) error
	UnregisterPeerMessageParty(peeringID PeeringID, partyType PeerMessagePartyType) error
}

type PeerMessagePartyGroup interface {
	RegisterPeerMessageParty(peeringID PeeringID, party PeerMessageGroupParty) error
	UnregisterPeerMessageParty(peeringID PeeringID, partyType PeerMessagePartyType) error
}

// GroupProvider stands for a subset of a peer-to-peer network
// that is responsible for achieving some common goal, eg,
// consensus committee, DKG group, etc.
//
// Indexes are only meaningful in the groups, not in the
// network or a particular peers.
type GroupProvider interface {
	SelfIndex() uint16
	PeerIndex(peer PeerSender) (uint16, error)
	PeerIndexByNetID(peerNetID string) (uint16, error)
	SendMsgByIndex(peerIdx uint16, peeringID PeeringID, destinationParty PeerMessagePartyType, msg Serializable)
	Broadcast(peeringID PeeringID, destinationParty PeerMessagePartyType, msg Serializable, includingSelf bool, except ...uint16)
	/* TODO	ExchangeRound(
		peers map[uint16]PeerSender,
		recvCh chan *RecvEvent,
		retryTimeout time.Duration,
		giveUpTimeout time.Duration,
		sendCB func(peerIdx uint16, peer PeerSender),
		recvCB func(recv *RecvEvent) (bool, error),
	) error */
	AllNodes(except ...uint16) map[uint16]PeerSender   // Returns all the nodes in the group except specified.
	OtherNodes(except ...uint16) map[uint16]PeerSender // Returns other nodes in the group (excluding Self and specified).
	Close()
	PeerMessagePartyGroup
}

// PeerDomainProvider implements unordered set of peers which can dynamically change
// All peers in the domain shares same peeringID. Each peer within domain is identified via its netID
type PeerDomainProvider interface {
	SendMsgByNetID(netID string, peeringID PeeringID, destinationParty PeerMessagePartyType, msg Serializable)
	SendMsgToRandomPeers(upToNumPeers uint16, peeringID PeeringID, destinationParty PeerMessagePartyType, msg Serializable)
	ReshufflePeers(seedBytes ...[]byte)
	GetRandomPeers(upToNumPeers int) []string
	Close()
	PeerMessagePartyCollection
}

// PeerSender represents an interface to some remote peer.
type PeerSender interface {

	// NetID identifies the peer.
	NetID() string

	// PubKey of the peer is only available, when it is
	// authenticated, therefore it can return nil, if pub
	// key is not known yet. You can call await before calling
	// this function to ensure the public key is already resolved.
	PubKey() *ed25519.PublicKey

	// SendMsg works in an asynchronous way, and therefore the
	// errors are not returned here.
	SendMsg(peeringID PeeringID, destinationParty PeerMessagePartyType, msg Serializable)

	// IsAlive indicates, if there is a working connection with the peer.
	// It is always an approximate state.
	IsAlive() bool

	// Await for the connection to be established, handshaked, and the
	// public key resolved.
	Await(timeout time.Duration) error

	// Close releases the reference to the peer, this informs the network
	// implementation, that it can disconnect, cleanup resources, etc.
	// You need to get new reference to the peer (PeerSender) to use it again.
	Close()
}

// PeerStatusProvider is used to access the current state of the network peer
// without allocating it (increading usage counters, etc). This interface
// overlaps with the PeerSender, and most probably they both will be implemented
// by the same object.
type PeerStatusProvider interface {
	NetID() string
	PubKey() *ed25519.PublicKey
	IsAlive() bool
	NumUsers() int
}

// RecvEvent stands for a received message along with
// the reference to its sender peer.
/*type RecvEvent struct {
	From PeerSender
	Msg  *PeerMessage
}*/

// PeerMessage is an envelope for all the messages exchanged via
// the peering module.
type PeerMessageData struct {
	PeeringID        PeeringID
	Timestamp        int64
	DestinationParty PeerMessagePartyType
	MsgDeserializer  PeerMessagePartyType
	MsgType          byte
	MsgData          []byte
}

type PeerMessageOut struct {
	PeerMessageData
	serialized *[]byte
}

type PeerMessageIn struct {
	PeerMessageData
	SenderIndex uint16 // TODO: Only meaningful in a group, and when calculated by the client.
	SenderNetID string // TODO: Non persistent. Only used by PeeringDomain, filled by the receiver
}

type PeerMessagePartyType byte

const (
	peerMessagePartyStart PeerMessagePartyType = iota
	PeerMessagePartyAcs
	PeerMessagePartyChain
	PeerMessagePartyCommittee
	PeerMessagePartyConsensus
	PeerMessagePartyDkg
	PeerMessagePartyStateManager
	peerMessagePartyEnd
)

func NewPeerMessagePartyType(b byte) (PeerMessagePartyType, error) {
	if (byte(peerMessagePartyStart) < b) && (b < byte(peerMessagePartyEnd)) {
		return PeerMessagePartyType(b), nil
	}
	return peerMessagePartyStart, xerrors.Errorf("byte %v is not of PeerMessagePartyType", b)
}

type Serializable interface {
	GetMsgType() byte // Must be unique to each deserialiser PeerMessageParty
	GetDeserializer() PeerMessagePartyType
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

type PeerMessageParty interface {
	// Each PeerMessageParty should statically implement a following method:
	// GetEmptyMessage(msgType byte) (Serializable, error)
	GetPartyType() PeerMessagePartyType
}

type PeerMessageSimpleParty interface {
	PeerMessageParty
	HandleMessage(senderNetID string, msg Serializable)
}

type PeerMessageGroupParty interface {
	PeerMessageParty
	HandleGroupMessage(senderNetID string, SenderIndex uint16, msg Serializable)
}

// ParseNetID parses the NetID and returns the corresponding host and port.
func ParseNetID(netID string) (string, int, error) {
	parts := strings.Split(netID, ":")
	if len(parts) != 2 {
		return "", 0, xerrors.Errorf("invalid NetID: %v", netID)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, xerrors.Errorf("invalid port in NetID: %v", netID)
	}
	return parts[0], port, nil
}
