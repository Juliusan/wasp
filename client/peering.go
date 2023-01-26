// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/http"
	"strings"

	"github.com/iotaledger/wasp/packages/webapi/v1/model"
	"github.com/iotaledger/wasp/packages/webapi/v1/routes"
	"github.com/iotaledger/wasp/packages/webapi/v2/models"
)

func (c *WaspClient) GetPeeringSelf() (*model.PeeringTrustedNode, error) {
	var response *model.PeeringTrustedNode
	err := c.do(http.MethodGet, routes.PeeringSelfGet(), nil, &response)
	return response, err
}

func (c *WaspClient) GetPeeringTrustedList() ([]*model.PeeringTrustedNode, error) {
	var response []*model.PeeringTrustedNode
	err := c.do(http.MethodGet, routes.PeeringTrustedList(), nil, &response)
	return response, err
}

func (c *WaspClient) GetPeeringTrusted(pubKey string) (*model.PeeringTrustedNode, error) {
	var response *model.PeeringTrustedNode
	err := c.do(http.MethodGet, routes.PeeringTrustedGet(pubKey), nil, &response)
	return response, err
}

func (c *WaspClient) CheckConnectedPeers(pubKeys []string) (*models.PeeringConnectedResponse, error) {
	var response *models.PeeringConnectedResponse
	err := c.do(http.MethodGet, "v2/node/peers/connected?publicKeys="+strings.Join(pubKeys, ","), nil, &response)
	return response, err
}

func (c *WaspClient) PutPeeringTrusted(pubKey, netID string) (*model.PeeringTrustedNode, error) {
	request := model.PeeringTrustedNode{
		PubKey: pubKey,
		NetID:  netID,
	}
	var response model.PeeringTrustedNode
	err := c.do(http.MethodPut, routes.PeeringTrustedPut(pubKey), request, &response)
	return &response, err
}

func (c *WaspClient) PostPeeringTrusted(pubKey, netID string) (*model.PeeringTrustedNode, error) {
	request := model.PeeringTrustedNode{
		PubKey: pubKey,
		NetID:  netID,
	}
	var response model.PeeringTrustedNode
	err := c.do(http.MethodPost, routes.PeeringTrustedPost(), request, &response)
	return &response, err
}

func (c *WaspClient) DeletePeeringTrusted(pubKey string) error {
	err := c.do(http.MethodDelete, routes.PeeringTrustedDelete(pubKey), nil, nil)
	return err
}
