package execution

import (
	"math/big"
	"time"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/api"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/vm/processors"
)

// The following interfaces define the common functionality for SC execution (VM/external view calls)

type WaspContext interface {
	LocateProgram(programHash hashing.HashValue) (vmtype string, binary []byte, err error)
	GetContractRecord(contractHname isc.Hname) (ret *root.ContractRecord)
	Processors() *processors.Cache
}

type GasContext interface {
	GasBurnEnabled() bool
	GasBurnEnable(enable bool)
	GasBurn(burnCode gas.BurnCode, par ...uint64)
	GasEstimateMode() bool
}

type WaspCallContext interface {
	WaspContext
	GasContext
	isc.LogInterface
	Timestamp() time.Time
	CurrentContractAccountID() isc.AgentID
	Caller() isc.AgentID
	GetNativeTokens(agentID isc.AgentID) iotago.NativeTokenSum
	GetBaseTokensBalance(agentID isc.AgentID) iotago.BaseToken
	GetNativeTokenBalance(agentID isc.AgentID, nativeTokenID iotago.NativeTokenID) *big.Int
	Call(msg isc.Message, allowance *isc.Assets) dict.Dict
	ChainID() isc.ChainID
	ChainAccountID() (iotago.AccountID, bool)
	ChainOwnerID() isc.AgentID
	ChainInfo() *isc.ChainInfo
	CurrentContractHname() isc.Hname
	Params() *isc.Params
	ContractStateReaderWithGasBurn() kv.KVStoreReader
	SchemaVersion() isc.SchemaVersion
	GasBurned() gas.GasUnits
	GasBudgetLeft() gas.GasUnits
	GetAccountNFTs(agentID isc.AgentID) []iotago.NFTID
	GetNFTData(nftID iotago.NFTID) *isc.NFT
	L1API() iotago.API
	TokenInfo() *api.InfoResBaseToken
}
