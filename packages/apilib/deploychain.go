// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package apilib

import (
	"fmt"
	"io"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/clients/multiclient"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/l1connection"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/util"

	"github.com/iotaledger/wasp/packages/registry"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
)

// TODO DeployChain on peering domain, not on committee

type CreateChainParams struct {
	Layer1Client         l1connection.Client
	CommitteeAPIHosts    []string
	N                    uint16
	T                    uint16
	OriginatorKeyPair    *cryptolib.KeyPair
	Textout              io.Writer
	Prefix               string
	InitParams           dict.Dict
	GovernanceController iotago.Address
}

// DeployChain creates a new chain on specified committee address
func DeployChain(
	par CreateChainParams,
	stateControllerAddr iotago.Address,
	govControllerAddr iotago.Address,
	deposit iotago.BaseToken,
	depositMana iotago.Mana,
	creationSlot iotago.SlotIndex,
) (isc.ChainID, error) {
	var err error
	textout := io.Discard
	if par.Textout != nil {
		textout = par.Textout
	}
	originatorAddr := par.OriginatorKeyPair.GetPublicKey().AsEd25519Address()

	fmt.Fprint(textout, par.Prefix)
	fmt.Fprintf(textout, "Creating new chain\n* Owner address:    %s\n* State controller: %s\n* committee size = %d\n* quorum = %d\n",
		originatorAddr, stateControllerAddr, par.N, par.T)
	fmt.Fprint(textout, par.Prefix)

	chainID, err := CreateChainOrigin(
		par.Layer1Client,
		par.OriginatorKeyPair,
		stateControllerAddr,
		govControllerAddr,
		deposit,
		depositMana,
		par.InitParams,
		creationSlot,
	)
	fmt.Fprint(textout, par.Prefix)
	if err != nil {
		fmt.Fprintf(textout, "Creating chain origin and init transaction.. FAILED: %v\n", err)
		return isc.ChainID{}, fmt.Errorf("DeployChain: %w", err)
	}
	fmt.Fprint(textout, par.Prefix)
	fmt.Fprintf(textout, "Chain has been created successfully on the Tangle.\n* ChainID: %s\n* State address: %s\n* committee size = %d\n* quorum = %d\n",
		chainID.Bech32(par.Layer1Client.Bech32HRP()), stateControllerAddr.Bech32(par.Layer1Client.Bech32HRP()), par.N, par.T)

	fmt.Fprintf(textout, "Make sure to activate the chain on all committee nodes\n")

	return chainID, err
}

func utxoIDsFromUtxoMap(utxoMap iotago.OutputSet) iotago.OutputIDs {
	var utxoIDs iotago.OutputIDs
	for id := range utxoMap {
		utxoIDs = append(utxoIDs, id)
	}
	return utxoIDs
}

// CreateChainOrigin creates and confirms origin transaction of the chain and init request transaction to initialize state of it
func CreateChainOrigin(
	layer1Client l1connection.Client,
	originator *cryptolib.KeyPair,
	stateController iotago.Address,
	governanceController iotago.Address,
	deposit iotago.BaseToken,
	depositMana iotago.Mana,
	initParams dict.Dict,
	creationSlot iotago.SlotIndex,
) (isc.ChainID, error) {
	originatorAddr := originator.GetPublicKey().AsEd25519Address()
	// ----------- request owner address' outputs from the ledger
	utxoMap, err := layer1Client.OutputMap(originatorAddr)
	if err != nil {
		return isc.ChainID{}, fmt.Errorf("CreateChainOrigin: %w", err)
	}

	tokenInfo, err := layer1Client.TokenInfo()
	if err != nil {
		return isc.ChainID{}, fmt.Errorf("CreateChainOrigin: %w", err)
	}

	// ----------- create origin transaction
	originTx, _, chainID, err := origin.NewChainOriginTransaction(
		originator,
		stateController,
		governanceController,
		deposit,
		depositMana,
		initParams,
		utxoMap,
		creationSlot,
		allmigrations.DefaultScheme.LatestSchemaVersion(),
		layer1Client.APIProvider(),
		tokenInfo,
	)
	if err != nil {
		return isc.ChainID{}, fmt.Errorf("CreateChainOrigin: %w", err)
	}

	blockIssuerID, err := util.BlockIssuerFromOutputs(utxoMap)
	if err != nil {
		return isc.ChainID{}, err
	}

	// ------------- post origin transaction and wait for confirmation
	_, err = layer1Client.PostTxAndWaitUntilConfirmation(originTx, blockIssuerID, originator)
	if err != nil {
		return isc.ChainID{}, fmt.Errorf("CreateChainOrigin: %w", err)
	}

	return chainID, nil
}

// ActivateChainOnNodes puts chain records into nodes and activates its
func ActivateChainOnNodes(clientResolver multiclient.ClientResolver, apiHosts []string, chainID isc.ChainID, l1API iotago.API) error {
	nodes := multiclient.New(clientResolver, apiHosts, l1API)
	// ------------ put chain records to hosts
	return nodes.PutChainRecord(registry.NewChainRecord(chainID, true, []*cryptolib.PublicKey{}))
}
