// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package peering

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

func initConnectedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-connected <pubKey|netID>?",
		Short: "Check if provided peers are inter-connected.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pubKeys := make([]string, len(args))
			for i := range pubKeys {
				if peering.CheckNetID(args[i]) == nil {
					// TODO: get public key from net ID
				}
				pubKeys[i] = args[i]
			}
			waspClient := config.WaspClient(config.MustWaspAPI())
			connected, err := waspClient.CheckConnectedPeers(pubKeys)
			log.Check(err)

			// Find out how many different public keys are there in query results
			resPubKeys := make([]string, len(connected.Sources))
			keyToIndexMap := make(map[string]int)
			for i := range connected.Sources {
				sPubKey := connected.Sources[i].PublicKey
				resPubKeys[i] = sPubKey
				keyToIndexMap[sPubKey] = i
			}
			for i := range connected.Sources {
				for j := range connected.Sources[i].Destinations {
					dPubKey := connected.Sources[i].Destinations[j].PublicKey
					_, ok := keyToIndexMap[dPubKey]
					if !ok {
						resPubKeys = append(resPubKeys, dPubKey)
						keyToIndexMap[dPubKey] = len(resPubKeys) - 1
					}
				}
			}

			// Make table header
			header := make([]string, len(resPubKeys)+2)
			header[0] = "Nr"
			header[1] = "PublicKey"

			// Fill table header column and table header rows
			result := make([][]string, len(resPubKeys))
			for i := range resPubKeys {
				header[i+2] = fmt.Sprint(i)
				result[i] = make([]string, len(resPubKeys)+2)
				result[i][0] = fmt.Sprint(i)
				result[i][1] = resPubKeys[i]
			}
			// Fill table data
			for i := range connected.Sources {
				for j := range connected.Sources[i].Destinations {
					dIndex := keyToIndexMap[connected.Sources[i].Destinations[j].PublicKey]
					failReason := connected.Sources[i].Destinations[j].FailReason
					if failReason == "" {
						failReason = "OK"
					}
					result[i][dIndex+2] = failReason
				}
			}
			log.PrintTable(header, result)
		},
	}
}
