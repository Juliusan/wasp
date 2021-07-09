package chainimpl

import (
	"bytes"
	"fmt"
	/*	"math/rand"
		"sync"*/
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/goshimmer/packages/ledgerstate/utxodb"
	"github.com/iotaledger/goshimmer/packages/ledgerstate/utxoutil"
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/util/ready"
	"github.com/iotaledger/wasp/packages/vm/processors"
	//	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/chain"
	/*	"github.com/iotaledger/wasp/packages/chain/committee"
		"github.com/iotaledger/wasp/packages/chain/mempool"*/
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/coretypes/chainid"
	/*		"github.com/iotaledger/wasp/packages/coretypes/coreutil"
			"github.com/iotaledger/wasp/packages/hashing"
			"github.com/iotaledger/wasp/packages/kv"*/
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/registry"
	/*		"github.com/iotaledger/wasp/packages/registry/committee_record"
			"github.com/iotaledger/wasp/packages/solo"*/
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/testutil/testchain"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/testutil/testpeers"
	//		"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type MockedEnv struct {
	T         *testing.T
	ChainImpl *chainObj
	//		Quorum            uint16
	Log               *logger.Logger
	Ledger            *utxodb.UtxoDB
	OriginatorKeyPair *ed25519.KeyPair
	OriginatorAddress ledgerstate.Address
	NodeIDs           []string
	NetworkProviders  []peering.NetworkProvider
	NetworkBehaviour  *testutil.PeeringNetDynamic
	ChainID           chainid.ChainID
	NodeConn          *testchain.MockedNodeConn
	SolidState        state.VirtualState       // State manager mock
	StateOutput       *ledgerstate.AliasOutput // State manager mock
	StateManagerMock  *chainimplStateManagerMock
	VMRunner          *MockedVMRunner
	Registry          *registry.Impl
	//StateTransition   *testchain.MockedStateTransition
	/*		MockedACS         chain.AsynchronousCommonSubsetRunner
			InitStateOutput   *ledgerstate.AliasOutput
			mutex             sync.Mutex
			Nodes             []*mockedNode*/
}

/*type mockedNode struct {
	NodeID      string
	Env         *MockedEnv
	ChainCore   *testchain.MockedChainCore // Chain mock
	stateSync   coreutil.ChainStateSync    // Chain mock
	Mempool     chain.Mempool              // Consensus needs
	Consensus   chain.Consensus            // Consensus needs
	store       kvstore.KVStore            // State manager mock
	SolidState  state.VirtualState         // State manager mock
	StateOutput *ledgerstate.AliasOutput   // State manager mock
	Log         *logger.Logger
	mutex       sync.Mutex
}*/

func NewMockedEnv(t *testing.T /*n, quorum uint16,*/, debug bool) (*MockedEnv, *ledgerstate.Transaction) {
	/*	return newMockedEnv(t, n, quorum, debug, false)
		}

		func NewMockedEnvWithMockedACS(t *testing.T, n, quorum uint16, debug bool) (*MockedEnv, *ledgerstate.Transaction) {
			return newMockedEnv(t, n, quorum, debug, true)
		}

		func newMockedEnv(t *testing.T, n, quorum uint16, debug, mockACS bool) (*MockedEnv, *ledgerstate.Transaction) {*/
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}
	log := testlogger.WithLevel(testlogger.NewLogger(t, "04:05.000"), level, false)
	var err error

	//log.Infof("creating test environment with N = %d, T = %d", n, quorum)

	ret := &MockedEnv{
		T: t,
		//		Quorum: quorum,
		Log:      log,
		Ledger:   utxodb.New(),
		NodeConn: testchain.NewMockedNodeConnection("Node"),
		//		Nodes:  make([]*mockedNode, n),*/
	}

	/*	if mockACS {
			ret.MockedACS = testchain.NewMockedACSRunner(quorum, log)
			log.Infof("running MOCKED ACS consensus")
		} else {
			log.Infof("running REAL ACS consensus")
		}*/

	ret.NetworkBehaviour = testutil.NewPeeringNetDynamic(log)

	log.Infof("running DKG and setting up mocked network..")
	nodeIDs, identities := testpeers.SetupKeys(1)
	ret.NodeIDs = nodeIDs
	ret.NetworkProviders = testpeers.SetupNet(ret.NodeIDs, identities, ret.NetworkBehaviour, log)

	cfg := &chainimplTestConfigProvider{
		ownNetID:  ret.NodeIDs[0],
		neighbors: ret.NodeIDs,
	}

	ret.OriginatorKeyPair, ret.OriginatorAddress = ret.Ledger.NewKeyPairByIndex(0)
	_, err = ret.Ledger.RequestFunds(ret.OriginatorAddress)
	require.NoError(t, err)

	outputs := ret.Ledger.GetAddressOutputs(ret.OriginatorAddress)
	require.True(t, len(outputs) == 1)

	bals := map[ledgerstate.Color]uint64{ledgerstate.ColorIOTA: 100}
	txBuilder := utxoutil.NewBuilder(outputs...)
	err = txBuilder.AddNewAliasMint(bals, ret.OriginatorAddress, state.OriginStateHash().Bytes())
	require.NoError(t, err)
	err = txBuilder.AddRemainderOutputIfNeeded(ret.OriginatorAddress, nil)
	require.NoError(t, err)
	originTx, err := txBuilder.BuildWithED25519(ret.OriginatorKeyPair)
	require.NoError(t, err)
	err = ret.Ledger.AddTransaction(originTx)
	require.NoError(t, err)

	ret.StateOutput, err = utxoutil.GetSingleChainedAliasOutput(originTx)
	require.NoError(t, err)

	ret.ChainID = *chainid.NewChainID(ret.StateOutput.GetAliasAddress())

	var makeNodeConnFun MakeNodeConnFun = func(*logger.Logger) chain.NodeConnection {
		return ret.NodeConn
	}

	ret.StateManagerMock = &chainimplStateManagerMock{
		t:   t,
		log: log,
		env: ret,
	}
	var makeStateManagerFun MakeStateManagerFun = func(store kvstore.KVStore, c chain.ChainCore, peers peering.PeerDomainProvider, nodeconn chain.NodeConnection) chain.StateManager {
		return ret.StateManagerMock
	}

	store := mapdb.NewMapDB()
	ret.Registry = registry.NewRegistry(log, store)

	ret.SolidState, err = state.CreateOriginState(store, &ret.ChainID)
	require.NoError(t, err)

	ret.ChainImpl = NewChain(&ret.ChainID, log, makeNodeConnFun, makeStateManagerFun, cfg, store, ret.NetworkProviders[0], ret.Registry, ret.Registry, ret.Registry, processors.NewConfig()).(*chainObj)

	ret.VMRunner = NewMockedVMRunner(ret)
	/*ret.StateTransition = testchain.NewMockedStateTransition(t, ret.OriginatorKeyPair)
	ret.StateTransition.OnNextState(func(vstate state.VirtualState, tx *ledgerstate.Transaction) {
		log.Debugf("MockedEnv.onNextState: state index %d", vstate.BlockIndex())
		ret.SolidState = vstate
		stateOutput, err := utxoutil.GetSingleChainedAliasOutput(tx)
		require.NoError(t, err)
		ret.StateOutput = stateOutput
		go ret.ChainImpl.ReceiveTransaction(tx)
		go ret.NodeConn.PostTransaction(tx)
	})*/
	ret.NodeConn.OnPostTransaction(func(tx *ledgerstate.Transaction) {
		log.Debugf("Mocked node conn: OnPostTransaction %v", tx.ID().Base58())
		_, exists := ret.Ledger.GetTransaction(tx.ID())
		if exists {
			log.Debugf("Mocked node conn: OnPostTransaction: posted repeating originTx: %s", tx.ID().Base58())
			return
		}
		if err := ret.Ledger.AddTransaction(tx); err != nil {
			log.Errorf("Mocked node conn: OnPostTransaction: error adding transaction: %v", err)
			return
		}
		ret.Log.Infof("Mocked node conn: OnPostTransaction: posted transaction to ledger: %s", tx.ID().Base58())
	})
	/*	ret.NodeConn.OnPullState(func(addr *ledgerstate.AliasAddress) {
		log.Debugf("Mocked node conn: PullState request, address %v", addr.Base58())
	})*/

	return ret, originTx
}

func (env *MockedEnv) NextState(tx *ledgerstate.Transaction, targetAddr ledgerstate.Address) {
	env.SolidState = env.VMRunner.Run(env.SolidState, targetAddr, tx)
	stateOutput, err := utxoutil.GetSingleChainedAliasOutput(tx)
	require.NoError(env.T, err)
	env.StateOutput = stateOutput
}

/*func (env *MockedEnv) NextState(reqs ...coretypes.Request) {
	env.StateTransition.NextState(env.SolidState, env.StateOutput, time.Now(), reqs...)
}*/

/*func (c *MockedStateTransition) NextState(tx *ledgerstate.Transaction, reqs ...coretypes.Request) {
	if c.chainKey != nil {
		require.True(c.t, chainOutput.GetStateAddress().Equals(ledgerstate.NewED25519Address(c.chainKey.PublicKey)))
	}

	nextvs := vs.Clone()
	prevBlockIndex := vs.BlockIndex()
	counterKey := kv.Key(coreutil.StateVarBlockIndex + "counter")

	counterBin, err := nextvs.KVStore().Get(counterKey)
	require.NoError(c.t, err)

	counter, _, err := codec.DecodeUint64(counterBin)
	require.NoError(c.t, err)

	suBlockIndex := state.NewStateUpdateWithBlockIndexMutation(prevBlockIndex + 1)

	suCounter := state.NewStateUpdate()
	counterBin = codec.EncodeUint64(counter + 1)
	suCounter.Mutations().Set(counterKey, counterBin)

	utxo := utxodb.New()
	_, addr0 := utxo.NewKeyPairByIndex(0)
	_, addr1 := utxo.NewKeyPairByIndex(1)
	_, addr2 := utxo.NewKeyPairByIndex(2)

	suReqs := state.NewStateUpdate()
	for i, req := range reqs {
		fmt.Printf("XXX doing req %v\n", i)
		key := kv.Key(blocklog.NewRequestLookupKey(vs.BlockIndex()+1, uint16(i)).Bytes())
		suReqs.Mutations().Set(key, req.ID().Bytes())
	}

	nextvs.ApplyStateUpdates(suBlockIndex, suCounter, suReqs)
	require.EqualValues(c.t, prevBlockIndex+1, nextvs.BlockIndex())

	nextStateHash := nextvs.Hash()

	txBuilder := utxoutil.NewBuilder(chainOutput).WithTimestamp(ts)
	err = txBuilder.AddAliasOutputAsRemainder(chainOutput.GetAliasAddress(), nextStateHash[:])
	require.NoError(c.t, err)

	if c.chainKey != nil {
		tx, err := txBuilder.BuildWithED25519(c.chainKey)
		require.NoError(c.t, err)
		reqs, err := request.RequestsOnLedgerFromTransaction(tx, chainOutput.GetStateAddress())
		fmt.Printf("XXX REQUESTs %v, err %v\n", len(reqs), err)
		reqs, err = request.RequestsOnLedgerFromTransaction(tx, addr0)
		fmt.Printf("XXX REQUESTs0 %v, err %v\n", len(reqs), err)
		reqs, err = request.RequestsOnLedgerFromTransaction(tx, addr1)
		fmt.Printf("XXX REQUESTs1 %v, err %v\n", len(reqs), err)
		reqs, err = request.RequestsOnLedgerFromTransaction(tx, addr2)
		fmt.Printf("XXX REQUESTs2 %v, err %v\n", len(reqs), err)
		reqs, err = request.RequestsOnLedgerFromTransaction(tx, chainOutput.GetAliasAddress())
		fmt.Printf("XXX REQUESTsa %v, err %v\n", len(reqs), err)
		c.onNextState(nextvs, tx)
	} else {
		tx, _, err := txBuilder.BuildEssence()
		require.NoError(c.t, err)
		c.onVMResult(nextvs, tx)
	}
}*/

//func (env *MockedEnv) GetRequestsOnLedger(t *testing.T, amount int) []*request.RequestOnLedger {
func (env *MockedEnv) GetTxWithRequestsOnLedger(t *testing.T, amount int) (*ledgerstate.Transaction, ledgerstate.Address) {
	/*utxo := utxodb.New()
	keyPair, addr := utxo.NewKeyPairByIndex(0)
	_, err := utxo.RequestFunds(addr)
	require.NoError(t, err)*/

	//outputs := env.Ledger.GetAddressOutputs(env.OriginatorAddress)
	outputs := env.Ledger.GetAddressOutputs(env.ChainID.AsAddress())
	require.True(t, len(outputs) == 1)

	//_, targetAddr := env.Ledger.NewKeyPairByIndex(1)
	targetAddr := env.ChainID.AsAddress()
	txBuilder := utxoutil.NewBuilder(outputs...)
	//txBuilder := utxoutil.NewBuilder(env.StateOutput)
	var err error
	var i uint64
	for i = 0; int(i) < amount; i++ {
		err = txBuilder.AddExtendedOutputConsume(targetAddr, util.Uint64To8Bytes(i), map[ledgerstate.Color]uint64{ledgerstate.ColorIOTA: 1})
		require.NoError(t, err)
	}
	err = txBuilder.AddRemainderOutputIfNeeded(env.OriginatorAddress, nil)
	require.NoError(t, err)
	tx, err := txBuilder.BuildWithED25519(env.OriginatorKeyPair)
	require.NoError(t, err)
	require.NotNil(t, tx)

	return tx, targetAddr

	/*	requests, err := request.RequestsOnLedgerFromTransaction(tx, targetAddr)
		require.NoError(t, err)
		require.True(t, amount == len(requests))
		return requests*/
}

/*func (env *MockedEnv) CreateNodes(timers ConsensusTimers) {
	for i := range env.Nodes {
		env.Nodes[i] = env.NewNode(uint16(i), timers)
	}
}

func (env *MockedEnv) NewNode(nodeIndex uint16, timers ConsensusTimers) *mockedNode { //nolint:revive
	nodeID := env.NodeIDs[nodeIndex]
	log := env.Log.Named(nodeID)
	ret := &mockedNode{
		NodeID:    nodeID,
		Env:       env,
		NodeConn:  testchain.NewMockedNodeConnection("Node_" + nodeID),
		store:     mapdb.NewMapDB(),
		ChainCore: testchain.NewMockedChainCore(env.T, env.ChainID, log),
		stateSync: coreutil.NewChainStateSync(),
		Log:       log,
	}
	ret.ChainCore.OnGlobalStateSync(func() coreutil.ChainStateSync {
		return ret.stateSync
	})
	ret.ChainCore.OnGetStateReader(func() state.OptimisticStateReader {
		return state.NewOptimisticStateReader(ret.store, ret.stateSync)
	})
	ret.NodeConn.OnPostTransaction(func(tx *ledgerstate.Transaction) {
		env.mutex.Lock()
		defer env.mutex.Unlock()

		if _, already := env.Ledger.GetTransaction(tx.ID()); !already {
			if err := env.Ledger.AddTransaction(tx); err != nil {
				ret.Log.Error(err)
				return
			}
			stateOutput := transaction.GetAliasOutput(tx, env.ChainID.AsAddress())
			require.NotNil(env.T, stateOutput)

			ret.Log.Infof("stored transaction to the ledger: %s", tx.ID().Base58())
			for _, node := range env.Nodes {
				go func(n *mockedNode) {
					ret.mutex.Lock()
					defer ret.mutex.Unlock()
					n.StateOutput = stateOutput
					n.checkStateApproval()
				}(node)
			}
		} else {
			ret.Log.Infof("transaction already in the ledger: %s", tx.ID().Base58())
		}
	})
	ret.NodeConn.OnPullTransactionInclusionState(func(addr ledgerstate.Address, txid ledgerstate.TransactionID) {
		if _, already := env.Ledger.GetTransaction(txid); already {
			go ret.ChainCore.ReceiveMessage(&chain.InclusionStateMsg{
				TxID:  txid,
				State: ledgerstate.Confirmed,
			})
		}
	})
	ret.Mempool = mempool.New(ret.ChainCore.GetStateReader(), coretypes.NewInMemoryBlobCache(), log)

	cfg := &consensusTestConfigProvider{
		ownNetID:  nodeID,
		neighbors: env.NodeIDs,
	}
	//
	// Pass the ACS mock, if it was set in env.MockedACS.
	acs := make([]chain.AsynchronousCommonSubsetRunner, 0, 1)
	if env.MockedACS != nil {
		acs = append(acs, env.MockedACS)
	}
	cmtRec := &committee_record.CommitteeRecord{
		Address: env.StateAddress,
		Nodes:   env.NodeIDs,
	}
	cmt, err := committee.New(
		cmtRec,
		&env.ChainID,
		env.NetworkProviders[nodeIndex],
		cfg,
		env.DKSRegistries[nodeIndex],
		log,
		acs...,
	)
	require.NoError(env.T, err)
	cmt.Attach(ret.ChainCore)

	ret.StateOutput = env.InitStateOutput
	ret.SolidState, err = state.CreateOriginState(ret.store, &env.ChainID)
	ret.stateSync.SetSolidIndex(0)
	require.NoError(env.T, err)

	cons := New(ret.ChainCore, ret.Mempool, cmt, ret.NodeConn, timers)
	cons.vmRunner = testchain.NewMockedVMRunner(env.T, log)
	ret.Consensus = cons

	ret.ChainCore.OnReceiveAsynchronousCommonSubsetMsg(func(msg *chain.AsynchronousCommonSubsetMsg) {
		ret.Consensus.EventAsynchronousCommonSubsetMsg(msg)
	})
	ret.ChainCore.OnReceiveVMResultMsg(func(msg *chain.VMResultMsg) {
		ret.Consensus.EventVMResultMsg(msg)
	})
	ret.ChainCore.OnReceiveInclusionStateMsg(func(msg *chain.InclusionStateMsg) {
		ret.Consensus.EventInclusionsStateMsg(msg)
	})
	ret.ChainCore.OnReceiveStateCandidateMsg(func(msg *chain.StateCandidateMsg) {
		ret.mutex.Lock()
		defer ret.mutex.Unlock()
		newState := msg.State
		ret.Log.Infof("chainCore.StateCandidateMsg: state hash: %s, approving output: %s",
			msg.State.Hash(), coretypes.OID(msg.ApprovingOutputID))

		if ret.SolidState != nil && ret.SolidState.BlockIndex() == newState.BlockIndex() {
			ret.Log.Debugf("new state already committed for index %d", newState.BlockIndex())
			return
		}
		err := newState.Commit()
		require.NoError(env.T, err)

		ret.SolidState = newState
		ret.Log.Debugf("committed new state for index %d", newState.BlockIndex())

		ret.checkStateApproval()
	})
	ret.ChainCore.OnReceivePeerMessage(func(msg *peering.PeerMessage) {
		var err error
		if msg.MsgType == chain.MsgSignedResult {
			decoded := chain.SignedResultMsg{}
			if err = decoded.Read(bytes.NewReader(msg.MsgData)); err == nil {
				decoded.SenderIndex = msg.SenderIndex
				ret.Consensus.EventSignedResultMsg(&decoded)
			}
		}
		if err != nil {
			ret.Log.Errorf("unexpected peer message type = %d", msg.MsgType)
		}
	})
	return ret
}

func (env *MockedEnv) nodeCount() int {
	return len(env.NodeIDs)
}

func (env *MockedEnv) SetInitialConsensusState() {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	for _, node := range env.Nodes {
		go func(n *mockedNode) {
			if n.SolidState != nil && n.SolidState.BlockIndex() == 0 {
				n.EventStateTransition()
			}
		}(node)
	}
}

func (n *mockedNode) checkStateApproval() {
	if n.SolidState == nil || n.StateOutput == nil {
		return
	}
	if n.SolidState.BlockIndex() != n.StateOutput.GetStateIndex() {
		return
	}
	stateHash, err := hashing.HashValueFromBytes(n.StateOutput.GetStateData())
	require.NoError(n.Env.T, err)
	require.EqualValues(n.Env.T, stateHash, n.SolidState.Hash())

	reqIDsForLastState := make([]coretypes.RequestID, 0)
	prefix := kv.Key(util.Uint32To4Bytes(n.SolidState.BlockIndex()))
	err = n.SolidState.KVStoreReader().Iterate(prefix, func(key kv.Key, value []byte) bool {
		reqid, err := coretypes.RequestIDFromBytes(value)
		require.NoError(n.Env.T, err)
		reqIDsForLastState = append(reqIDsForLastState, reqid)
		return true
	})
	require.NoError(n.Env.T, err)
	n.Mempool.RemoveRequests(reqIDsForLastState...)

	n.Log.Infof("STATE APPROVED (%d reqs). Index: %d, State output: %s",
		len(reqIDsForLastState), n.SolidState.BlockIndex(), coretypes.OID(n.StateOutput.ID()))

	n.EventStateTransition()
}

func (n *mockedNode) EventStateTransition() {
	n.Log.Debugf("EventStateTransition")

	n.ChainCore.GlobalStateSync().SetSolidIndex(n.SolidState.BlockIndex())

	n.Consensus.EventStateTransitionMsg(&chain.StateTransitionMsg{
		State:          n.SolidState.Clone(),
		StateOutput:    n.StateOutput,
		StateTimestamp: time.Now(),
	})
}

func (env *MockedEnv) StartTimers() {
	for _, n := range env.Nodes {
		n.StartTimer()
	}
}

func (n *mockedNode) StartTimer() {
	n.Log.Debugf("started timer..")
	go func() {
		counter := 0
		for {
			n.Consensus.EventTimerMsg(chain.TimerTick(counter))
			counter++
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (env *MockedEnv) WaitTimerTick(until int) error {
	checkTimerTickFun := func(node *mockedNode) bool {
		snap := node.Consensus.GetStatusSnapshot()
		if snap != nil && snap.TimerTick >= until {
			return true
		}
		return false
	}
	return env.WaitForEventFromNodes("TimerTick", checkTimerTickFun)
}

func (env *MockedEnv) WaitStateIndex(quorum int, stateIndex uint32, timeout ...time.Duration) error {
	checkStateIndexFun := func(node *mockedNode) bool {
		snap := node.Consensus.GetStatusSnapshot()
		if snap != nil && snap.StateIndex >= stateIndex {
			return true
		}
		return false
	}
	return env.WaitForEventFromNodesQuorum("stateIndex", quorum, checkStateIndexFun, timeout...)
}

func (env *MockedEnv) WaitMempool(numRequests int, quorum int, timeout ...time.Duration) error { //nolint:gocritic
	checkMempoolFun := func(node *mockedNode) bool {
		snap := node.Consensus.GetStatusSnapshot()
		if snap != nil && snap.Mempool.InPoolCounter >= numRequests && snap.Mempool.OutPoolCounter >= numRequests {
			return true
		}
		return false
	}
	return env.WaitForEventFromNodesQuorum("mempool", quorum, checkMempoolFun, timeout...)
}

func (env *MockedEnv) WaitForEventFromNodes(waitName string, nodeConditionFun func(node *mockedNode) bool, timeout ...time.Duration) error {
	return env.WaitForEventFromNodesQuorum(waitName, env.nodeCount(), nodeConditionFun, timeout...)
}

func (env *MockedEnv) WaitForEventFromNodesQuorum(waitName string, quorum int, isEventOccuredFun func(node *mockedNode) bool, timeout ...time.Duration) error {
	to := 10 * time.Second
	if len(timeout) > 0 {
		to = timeout[0]
	}
	ch := make(chan int)
	nodeCount := env.nodeCount()
	deadline := time.Now().Add(to)
	for _, n := range env.Nodes {
		go func(node *mockedNode) {
			for time.Now().Before(deadline) {
				if isEventOccuredFun(node) {
					ch <- 1
				}
				time.Sleep(10 * time.Millisecond)
			}
			ch <- 0
		}(n)
	}
	var sum, total int
	for n := range ch {
		sum += n
		total++
		if sum >= quorum {
			return nil
		}
		if total >= nodeCount {
			return fmt.Errorf("Wait for %s: test timeouted", waitName)
		}
	}
	return fmt.Errorf("WaitMempool: timeout expired %v", to)
}

func (env *MockedEnv) PostDummyRequests(n int, randomize ...bool) {
	reqs := make([]coretypes.Request, n)
	for i := 0; i < n; i++ {
		reqs[i] = solo.NewCallParams("dummy", "dummy", "c", i).
			NewRequestOffLedger(env.OriginatorKeyPair)
	}
	rnd := len(randomize) > 0 && randomize[0]
	for _, n := range env.Nodes {
		if rnd {
			for _, req := range reqs {
				go func(node *mockedNode, r coretypes.Request) {
					time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
					node.Mempool.ReceiveRequests(r)
				}(n, req)
			}
		} else {
			n.Mempool.ReceiveRequests(reqs...)
		}
	}
}*/

var _ chain.StateManager = &chainimplStateManagerMock{}

type chainimplStateManagerMock struct {
	t   *testing.T
	log *logger.Logger
	env *MockedEnv
}

func (cT *chainimplStateManagerMock) Ready() *ready.Ready {
	cT.log.Debugf("mocked state manager: Ready called")
	r := ready.New(fmt.Sprintf("mocked state manager %s", cT.env.ChainImpl.ID().Base58()[:6]+".."))
	r.SetReady()
	return r
}

func (cT *chainimplStateManagerMock) EventGetBlockMsg(msg *chain.GetBlockMsg) {
	cT.log.Debugf("mocked state manager: EventGetBlockMsg called")
}

func (cT *chainimplStateManagerMock) EventBlockMsg(msg *chain.BlockMsg) {
	cT.log.Debugf("mocked state manager: EventBlockMsg called")
}

func (cT *chainimplStateManagerMock) EventStateMsg(msg *chain.StateMsg) {
	cT.log.Debugf("mocked state manager: EventStateMsg called")
}

func (cT *chainimplStateManagerMock) EventOutputMsg(msg ledgerstate.Output) {
	cT.log.Debugf("mocked state manager: EventOutputMsg called")
}

func (cT *chainimplStateManagerMock) EventStateCandidateMsg(msg *chain.StateCandidateMsg) {
	cT.log.Debugf("mocked state manager: EventStateCandidateMsg called")
}

func (cT *chainimplStateManagerMock) EventTimerMsg(msg chain.TimerTick) {
	cT.log.Debugf("mocked state manager: EventTimerMsg called")
	stats := cT.env.ChainImpl.mempool.Stats()
	cT.log.Debugf("XXX MEMPOOL: IN %v OUT %v", stats.InPoolCounter, stats.OutPoolCounter)
}

func (cT *chainimplStateManagerMock) GetStatusSnapshot() *chain.SyncInfo {
	cT.log.Debugf("mocked state manager: GetStatusSnapshot called")
	outputStateHash, _ := hashing.HashValueFromBytes(cT.env.StateOutput.GetStateData())
	return &chain.SyncInfo{
		Synced:                bytes.Equal(cT.env.SolidState.Hash().Bytes(), cT.env.StateOutput.GetStateData()),
		SyncedBlockIndex:      cT.env.SolidState.BlockIndex(),
		SyncedStateHash:       cT.env.SolidState.Hash(),
		SyncedStateTimestamp:  cT.env.SolidState.Timestamp(),
		StateOutputBlockIndex: cT.env.StateOutput.GetStateIndex(),
		StateOutputID:         cT.env.StateOutput.ID(),
		StateOutputHash:       outputStateHash,
		StateOutputTimestamp:  time.Now(), // not very correct, probably
	}
}

func (cT *chainimplStateManagerMock) Close() {
	cT.log.Debugf("mocked state manager: Close called")
}

// TODO: should this object be obtained from peering.NetworkProvider?
// Or should coretypes.PeerNetworkConfigProvider methods methods be part of
// peering.NetworkProvider interface
var _ coretypes.PeerNetworkConfigProvider = &chainimplTestConfigProvider{}

type chainimplTestConfigProvider struct {
	ownNetID  string
	neighbors []string
}

func (p *chainimplTestConfigProvider) OwnNetID() string {
	return p.ownNetID
}

func (p *chainimplTestConfigProvider) PeeringPort() int {
	return 0 // Anything
}

func (p *chainimplTestConfigProvider) Neighbors() []string {
	return p.neighbors
}

func (p *chainimplTestConfigProvider) String() string {
	return fmt.Sprintf("chainimplTestConfigProvider( ownNetID: %s, neighbors: %+v )", p.OwnNetID(), p.Neighbors())
}
