// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc

import (
	"fmt"
	"io"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

const ChainIDLength = iotago.AccountIDLength

var emptyChainID = ChainID{}

// ChainID represents the global identifier of the chain
// It is wrapped AccountAddress, an address without a private key behind
type (
	ChainID    iotago.AccountID
	ChainIDKey string
)

// EmptyChainID returns an empty ChainID.
func EmptyChainID() ChainID {
	return emptyChainID
}

func ChainIDFromAddress(addr *iotago.AccountAddress) ChainID {
	return ChainIDFromAccountID(addr.AccountID())
}

// ChainIDFromAccountID creates new chain ID from account address
func ChainIDFromAccountID(aliasID iotago.AccountID) ChainID {
	return ChainID(aliasID)
}

// ChainIDFromBytes reconstructs a ChainID from its binary representation.
func ChainIDFromBytes(data []byte) (ret ChainID, err error) {
	_, err = rwutil.ReadFromBytes(data, &ret)
	return ret, err
}

func ChainIDFromString(bech32 string) (ChainID, error) {
	_, addr, err := iotago.ParseBech32(bech32)
	if err != nil {
		return ChainID{}, err
	}
	if addr.Type() != iotago.AddressAccount {
		return ChainID{}, fmt.Errorf("chainID must be an account address (%s)", bech32)
	}
	return ChainIDFromAddress(addr.(*iotago.AccountAddress)), nil
}

func ChainIDFromKey(key ChainIDKey) ChainID {
	chainID, err := ChainIDFromString(string(key))
	if err != nil {
		panic(err)
	}
	return chainID
}

// RandomChainID creates a random chain ID. Used for testing only
func RandomChainID(seed ...[]byte) ChainID {
	var h hashing.HashValue
	if len(seed) > 0 {
		h = hashing.HashData(seed[0])
	} else {
		h = hashing.PseudoRandomHash(nil)
	}
	chainID, err := ChainIDFromBytes(h[:ChainIDLength])
	if err != nil {
		panic(err)
	}
	return chainID
}

func (id ChainID) AsAddress() iotago.Address {
	addr := iotago.AccountAddress(id)
	return &addr
}

func (id ChainID) AsAccountAddress() iotago.AccountAddress {
	return iotago.AccountAddress(id)
}

func (id ChainID) AsAccountID() iotago.AccountID {
	return iotago.AccountID(id)
}

func (id ChainID) Bytes() []byte {
	return id[:]
}

func (id ChainID) Empty() bool {
	return id == emptyChainID
}

func (id ChainID) Equals(other ChainID) bool {
	return id == other
}

func (id ChainID) Key() ChainIDKey {
	return ChainIDKey(id.AsAccountID().String())
}

func (id ChainID) IsSameChain(agentID AgentID) bool {
	contract, ok := agentID.(*ContractAgentID)
	if !ok {
		return false
	}
	return id.Equals(contract.ChainID())
}

func (id ChainID) ShortString() string {
	return id.AsAddress().String()[0:10]
}

// String human-readable form (bech32)
func (id ChainID) String() string {
	return id.AsAddress().Bech32(parameters.NetworkPrefix())
}

func (id *ChainID) Read(r io.Reader) error {
	return rwutil.ReadN(r, id[:])
}

func (id *ChainID) Write(w io.Writer) error {
	return rwutil.WriteN(w, id[:])
}
