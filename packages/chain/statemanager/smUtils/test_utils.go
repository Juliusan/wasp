// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package smUtils

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
)

func MakeNodeID(index int) gpa.NodeID {
	return gpa.NodeID(fmt.Sprintf("Node%v", index))
}

func MakeNodeIDs(indexes []int) []gpa.NodeID {
	result := make([]gpa.NodeID, len(indexes))
	for i := range indexes {
		result[i] = MakeNodeID(indexes[i])
	}
	return result
}
