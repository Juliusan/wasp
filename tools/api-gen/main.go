package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"

	"github.com/iotaledger/iota.go/v4/api"
	"github.com/iotaledger/wasp/components/app"
	"github.com/iotaledger/wasp/packages/authentication"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/evm/jsonrpc"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/webapi"
)

type NodeIdentityProviderMock struct{}

func (n *NodeIdentityProviderMock) NodeIdentity() *cryptolib.KeyPair {
	return cryptolib.NewKeyPair()
}

func (n *NodeIdentityProviderMock) NodePublicKey() *cryptolib.PublicKey {
	return cryptolib.NewEmptyPublicKey()
}

func main() {
	mockLog := testlogger.NewSilentLogger(false, "api-gen")
	e := echo.New()

	if app.Version == "" {
		app.Version = "0"
	}

	swagger := webapi.CreateEchoSwagger(e, app.Version)
	webapi.Init(mockLog, swagger, app.Version, nil, nil, nil, nil, nil, nil, &NodeIdentityProviderMock{}, nil, nil, nil, nil, authentication.AuthConfiguration{Scheme: authentication.AuthJWT}, time.Second, nil, nil, nil, jsonrpc.ParametersDefault(), nil, &api.InfoResBaseToken{})

	root, ok := swagger.(*echoswagger.Root)
	if !ok {
		panic("failed to get swagger root")
	}

	schema, err := root.GetSpec(nil, "/doc")
	if err != nil {
		panic(err.Error())
	}

	jsonSchema, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	fmt.Print(string(jsonSchema))
}
