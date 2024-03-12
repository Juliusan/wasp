// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// TODO: Test connect/reconnect - start node conn, and later the hornet.
// TODO: Test connect/reconnect - on a running node stop and later restart hornet.

package tests

import (
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/l1connection"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
)

func createChain(t *testing.T) isc.ChainID {
	originator := cryptolib.NewKeyPair()
	layer1Client := l1connection.NewClient(l1.Config, testlogger.NewLogger(t))
	err := layer1Client.RequestFunds(originator)
	require.NoError(t, err)

	utxoMap, err := layer1Client.OutputMap(originator.Address())
	require.NoError(t, err)

	originTx, _, chainID, err := origin.NewChainOriginTransaction(
		originator,
		originator.Address(),
		originator.Address(),
		0,
		iotago.Mana(0),
		nil,
		utxoMap,
		layer1Client.APIProvider().LatestAPI().TimeProvider().SlotFromTime(time.Now()),
		allmigrations.DefaultScheme.LatestSchemaVersion(),
		layer1Client.APIProvider(),
		lo.Must(layer1Client.TokenInfo()),
	)
	require.NoError(t, err)
	blockIssuerID, err := util.BlockIssuerFromOutputs(utxoMap)
	require.NoError(t, err)
	_, err = layer1Client.PostTxAndWaitUntilConfirmation(originTx, blockIssuerID, originator)
	require.NoError(t, err)

	return chainID
}

func TestNodeConn(t *testing.T) {
	panic("TODO this needs to be re-written")
	// if testing.Short() {
	// 	t.Skip("Skipping nodeconn test in short mode")
	// }

	// l1.StartPrivtangleIfNecessary(t.Logf)

	// log := testlogger.NewLogger(t)
	// peerCount := 1

	// //
	// // Start a peering network.
	// // peeringID := peering.RandomPeeringID()
	// peeringURLs, peerIdentities := testpeers.SetupKeys(uint16(peerCount))
	// networkLog := testlogger.WithLevel(log.NewChildLogger("Network"), 0)
	// _, networkCloser := testpeers.SetupNet(
	// 	peeringURLs,
	// 	peerIdentities,
	// 	testutil.NewPeeringNetReliable(networkLog),
	// 	networkLog,
	// )
	// t.Log("Peering network created.")

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// ctxInit, cancelInit := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancelInit()

	// nodeBridge := nodebridge.New(log.NewChildLogger("NodeBridge"))
	// err := nodeBridge.Connect(ctxInit, l1.Config.INXAddress, 10)
	// require.NoError(t, err)

	// go nodeBridge.Run(ctx)

	// nc, err := nodeconn.New(ctxInit, log, nodeBridge, nil)
	// require.NoError(t, err)

	// // run the node connection
	// go nc.Run(ctx)

	// nc.WaitUntilInitiallySynced(ctxInit)

	// //
	// // Check the chain operations.
	// chainID := createChain(t)
	// chainOuts := make(map[iotago.OutputID]iotago.Output)
	// chainOICh := make(chan iotago.OutputID, 100)
	// chainStateOuts := make(map[iotago.OutputID]iotago.Output)
	// chainStateOutsICh := make(chan iotago.OutputID, 100)

	// drainChannel := func(channel chan iotago.OutputID) {
	// 	for {
	// 		select {
	// 		case <-channel:
	// 		default:
	// 			return
	// 		}
	// 	}
	// }

	// drainChannels := func() {
	// 	drainChannel(chainOICh)
	// 	drainChannel(chainStateOutsICh)
	// }

	// nc.AttachChain(
	// 	context.Background(),
	// 	chainID,
	// 	func(outputInfo *isc.OutputInfo) {
	// 		chainOuts[outputInfo.OutputID] = outputInfo.Output
	// 		chainOICh <- outputInfo.OutputID
	// 	},
	// 	func(outputInfo *isc.OutputInfo) {
	// 		chainStateOuts[outputInfo.OutputID] = outputInfo.Output
	// 		chainStateOutsICh <- outputInfo.OutputID
	// 	},
	// 	func(timestamp time.Time) {},
	// 	nil,
	// 	nil,
	// )

	// client := l1connection.NewClient(l1.Config, log)

	// drainChannels()

	// // // Post a TX directly, and wait for it in the message stream (e.g. a request).
	// err = client.RequestFunds(chainID.AsAddress())
	// require.NoError(t, err)

	// t.Log("Waiting for outputs posted via tangle...")
	// oid := <-chainOICh
	// t.Logf("Waiting for outputs posted via tangle... Done, have %v=%v", oid.ToHex(), chainOuts[oid])

	// drainChannels()

	// wallet := cryptolib.NewKeyPair()
	// client.RequestFunds(wallet)
	// tx, err := client.MakeSimpleValueTX(wallet, chainID.AsAddress(), 1*isc.Million)
	// require.NoError(t, err)

	// ctxPublish, cancelPublish := context.WithCancel(context.Background())
	// nc.PublishTX(ctxPublish, chainID, tx, func(tx *iotago.Transaction, confirmed bool) {
	// 	require.True(t, confirmed)
	// 	cancelPublish()
	// })

	// t.Log("Waiting for outputs posted via nodeConn...")
	// oid := <-chainOICh
	// t.Logf("Waiting for outputs posted via nodeConn... Done, have %v=%v", oid.ToHex(), chainOuts[oid])

	// //
	// // Cleanup.
	// require.NoError(t, networkCloser.Close())
}
