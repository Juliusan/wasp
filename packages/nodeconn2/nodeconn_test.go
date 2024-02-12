// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package nodeconn2

import (
	"context"
	"flag"
	"fmt"
	"testing"
	//	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/app/daemon"
	"github.com/iotaledger/hive.go/app/shutdown"
	"github.com/iotaledger/inx-app/pkg/nodebridge"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/util/l1starter"
)

// NOTE:
// * To run these tests, the user must:
//		1) git clone git@github.com:iotaledger/iota-core.git
//		2) cd iota-core/tools/docker-network
//		3) docker compose pull && docker compose build # this failed for me... Not needed?
//		4) bash run.sh # and cancel it after it builds
// * To stop the computer from dying after these tests are run, do this:
//		docker kill $(docker ps -q)
// * If the test fails to find IPs, do this:
//		docker network prune

func TestSmth(t *testing.T) {
	logger := testlogger.NewLogger(t)
	starter := l1starter.New(flag.CommandLine, flag.CommandLine)
	starter.StartPrivtangleIfNecessary(func(format string, args ...any) { logger.LogDebugf(format, args...) })
	t.Cleanup(func() { starter.Cleanup(t.Logf) })
	nodeBridge := nodebridge.New(logger)
	ctx := context.Background()
	fmt.Printf("XXX KUKU1 9059\n")
	//go func() {
	err := nodeBridge.Connect(ctx, "localhost:9059", 10)
	require.NoError(t, err)
	fmt.Printf("XXX KUKU2\n")
	//}()
	//time.Sleep(1 * time.Second)
	_, err = NewNodeConn(ctx, logger, nodeBridge, shutdown.NewShutdownHandler(logger, daemon.New()))
	require.NoError(t, err)
	fmt.Printf("XXX KUKU3\n")
}
