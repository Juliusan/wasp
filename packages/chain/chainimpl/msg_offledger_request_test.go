// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainimpl

import (
	"bytes"
	"testing"

	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/request"
	"github.com/iotaledger/wasp/packages/iscp/requestargs"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/stretchr/testify/require"
)

const foo = "foo"

func TestMarshalling(t *testing.T) {
	// construct a dummy offledger request
	contract := iscp.Hn("somecontract")
	entrypoint := iscp.Hn("someentrypoint")
	args := requestargs.New(
		dict.Dict{foo: []byte("bar")},
	)

	msg := &offLedgerRequestMsg{
		ChainID: iscp.RandomChainID(),
		Req:     request.NewOffLedger(contract, entrypoint, args),
	}

	// marshall the msg
	msgBytes, err := msg.Serialize()
	require.NoError(t, err)

	// unmashal the message from bytes and ensure everything checks out
	unmarshalledMsg := &offLedgerRequestMsg{}
	err = unmarshalledMsg.Deserialize(msgBytes)
	require.NoError(t, err)

	require.True(t, unmarshalledMsg.ChainID.AliasAddress.Equals(msg.ChainID.AliasAddress))
	require.True(t, bytes.Equal(unmarshalledMsg.Req.Bytes(), msg.Req.Bytes()))
}
