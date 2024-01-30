package isc

import (
	iotago "github.com/iotaledger/iota.go/v4"
)

func TransactionIDsEqual(id1, id2 iotago.TransactionID) bool {
	return id1 == id2
}
