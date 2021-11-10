package domain

import (
	"crypto/rand"
	"sort"
	"sync"

	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/util"
)

type domainImpl struct {
	attachedTo  peering.PeeringID
	netProvider peering.NetworkProvider
	nodes       map[string]peering.PeerSender
	permutation *util.Permutation16
	netIDs      []string
	log         *logger.Logger
	mutex       *sync.RWMutex
}

var _ peering.PeerDomainProvider = &domainImpl{}

// NewPeerDomain creates a collection. Ignores self
func NewPeerDomain(netProvider peering.NetworkProvider, initialNodes []peering.PeerSender, log *logger.Logger) *domainImpl {
	ret := &domainImpl{
		netProvider: netProvider,
		nodes:       make(map[string]peering.PeerSender),
		permutation: util.NewPermutation16(uint16(len(initialNodes)), nil),
		netIDs:      make([]string, 0, len(initialNodes)),
		log:         log,
		mutex:       &sync.RWMutex{},
	}
	for _, sender := range initialNodes {
		ret.nodes[sender.NetID()] = sender
	}
	ret.reshufflePeers()
	return ret
}

func NewPeerDomainByNetIDs(netProvider peering.NetworkProvider, peerNetIDs []string, log *logger.Logger) (*domainImpl, error) {
	peers := make([]peering.PeerSender, 0, len(peerNetIDs))
	for _, nid := range peerNetIDs {
		if nid == netProvider.Self().NetID() {
			continue
		}
		peer, err := netProvider.PeerByNetID(nid)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}
	return NewPeerDomain(netProvider, peers, log), nil
}

func (d *domainImpl) SendMsgByNetID(netID string, peeringID peering.PeeringID, destinationParty peering.PeerMessagePartyType, msg peering.Serializable) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	d.sendMsgByNetIDNoLock(netID, peeringID, destinationParty, msg)
}

func (d *domainImpl) sendMsgByNetIDNoLock(netID string, peeringID peering.PeeringID, destinationParty peering.PeerMessagePartyType, msg peering.Serializable) {
	peer, ok := d.nodes[netID]
	if !ok {
		d.log.Warnf("sendMsgByNetID: NetID %v is not in the domain", netID)
		return
	}
	peer.SendMsg(peeringID, destinationParty, msg)
}

func (d *domainImpl) SendMsgToRandomPeers(upToNumPeers uint16, peeringID peering.PeeringID, destinationParty peering.PeerMessagePartyType, msg peering.Serializable) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if int(upToNumPeers) > len(d.nodes) {
		upToNumPeers = uint16(len(d.nodes))
	}
	for i := uint16(0); i < upToNumPeers; i++ {
		d.sendMsgByNetIDNoLock(d.netIDs[d.permutation.Next()], peeringID, destinationParty, msg)
	}
}

func (d *domainImpl) GetRandomPeers(upToNumPeers int) []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if upToNumPeers > len(d.netIDs) {
		upToNumPeers = len(d.netIDs)
	}
	ret := make([]string, upToNumPeers)
	for i := range ret {
		ret[i] = d.netIDs[d.permutation.Next()]
	}
	return ret
}

func (d *domainImpl) AddPeer(netID string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, ok := d.nodes[netID]; ok {
		return nil
	}
	if netID == d.netProvider.Self().NetID() {
		return nil
	}
	peer, err := d.netProvider.PeerByNetID(netID)
	if err != nil {
		return err
	}
	d.nodes[netID] = peer
	d.reshufflePeers()

	return nil
}

func (d *domainImpl) RemovePeer(netID string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.nodes, netID)
	d.reshufflePeers()
}

func (d *domainImpl) ReshufflePeers(seedBytes ...[]byte) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.reshufflePeers(seedBytes...)
}

func (d *domainImpl) reshufflePeers(seedBytes ...[]byte) {
	d.netIDs = make([]string, 0, len(d.nodes))
	for netID := range d.nodes {
		d.netIDs = append(d.netIDs, netID)
	}
	sort.Strings(d.netIDs)
	var seedB []byte
	if len(seedBytes) == 0 {
		var b [8]byte
		seedB = b[:]
		_, _ = rand.Read(seedB)
	} else {
		seedB = seedBytes[0]
	}
	d.permutation.Shuffle(seedB)
}

func (d *domainImpl) RegisterPeerMessageParty(peeringID peering.PeeringID, party peering.PeerMessageSimpleParty) error {
	return d.netProvider.RegisterPeerMessageParty(peeringID, NewPeerMessageDomainSimpleParty(party, d))
}

func (d *domainImpl) UnregisterPeerMessageParty(peeringID peering.PeeringID, partyType peering.PeerMessagePartyType) error {
	return d.netProvider.UnregisterPeerMessageParty(peeringID, partyType)
}

//TODO
func (d *domainImpl) Close() {
	for i := range d.nodes {
		d.nodes[i].Close()
	}
}
