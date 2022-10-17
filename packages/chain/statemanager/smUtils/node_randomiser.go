//
//
//
//
//
//

package smUtils

import (
	"sync"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

type NodeRandomiser struct {
	//log           *logger.Logger
	me          gpa.NodeID
	nodeIDs     []gpa.NodeID
	permutation *util.Permutation16
	mutex       sync.RWMutex
}

func NewNodeRandomiser(me gpa.NodeID, nodeIDs []gpa.NodeID) *NodeRandomiser {
	result := NewNodeRandomiserNoInit(me)
	result.init(nodeIDs)
	return result
}

// Before using the returned NodeRandomiser, it must be initted: UpdateNodeIDs
// method must be called.
func NewNodeRandomiserNoInit(me gpa.NodeID) *NodeRandomiser {
	return &NodeRandomiser{
		//log:           smLog,
		me:          me,
		nodeIDs:     nil, // Will be set in result.UpdateNodeIDs([]gpa.NodeID).
		permutation: nil, // Will be set in result.UpdateNodeIDs([]gpa.NodeID).
	}
}

func (nrT *NodeRandomiser) init(allNodeIDs []gpa.NodeID) {
	nrT.nodeIDs = make([]gpa.NodeID, 0, len(allNodeIDs)-1)
	for _, nodeID := range allNodeIDs {
		if nodeID != nrT.me { // Do not include self to the permutation.
			nrT.nodeIDs = append(nrT.nodeIDs, nodeID)
		}
	}
	var err error
	nrT.permutation, err = util.NewPermutation16(uint16(len(nrT.nodeIDs)))
	if err != nil {
		// TODO
		//d.log.Warnf("Error generating cryptographically secure random domains permutation: %v", err)
		return
	}
}

func (nrT *NodeRandomiser) UpdateNodeIDs(nodeIDs []gpa.NodeID) {
	nrT.mutex.Lock()
	defer nrT.mutex.Unlock()

	nrT.init(nodeIDs)
}

func (nrT *NodeRandomiser) IsInitted() bool {
	return nrT.permutation != nil
}

func (nrT *NodeRandomiser) GetRandomOtherNodeIDs(upToNumPeers int) []gpa.NodeID {
	nrT.mutex.RLock()
	defer nrT.mutex.RUnlock()

	if upToNumPeers > len(nrT.nodeIDs) {
		upToNumPeers = len(nrT.nodeIDs)
	}
	ret := make([]gpa.NodeID, upToNumPeers)
	for i := 0; i < upToNumPeers; {
		ret[i] = nrT.nodeIDs[nrT.permutation.NextNoCycles()]
		distinct := true
		for j := 0; j < i && distinct; j++ {
			if ret[i].Equals(ret[j]) {
				distinct = false
			}
		}
		if distinct {
			i++
		}
	}
	return ret
}
