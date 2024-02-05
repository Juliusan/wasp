package isc

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/minio/blake2b-simd"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/util/rwutil"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

type offLedgerSignature struct {
	publicKey *cryptolib.PublicKey
	signature []byte
}

type OffLedgerRequestData struct {
	allowance *Assets
	chainID   ChainID
	msg       Message
	gasBudget gas.GasUnits
	nonce     uint64
	signature offLedgerSignature
}

var (
	_ Request                  = new(OffLedgerRequestData)
	_ OffLedgerRequest         = new(OffLedgerRequestData)
	_ UnsignedOffLedgerRequest = new(OffLedgerRequestData)
	_ Calldata                 = new(OffLedgerRequestData)
	_ Features                 = new(OffLedgerRequestData)
)

func NewOffLedgerRequest(
	chainID ChainID,
	msg Message,
	nonce uint64,
	gasBudget gas.GasUnits,
) UnsignedOffLedgerRequest {
	return &OffLedgerRequestData{
		chainID:   chainID,
		msg:       msg,
		nonce:     nonce,
		allowance: NewEmptyAssets(),
		gasBudget: gasBudget,
	}
}

func (req *OffLedgerRequestData) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	req.readEssence(rr)
	req.signature.publicKey = cryptolib.NewEmptyPublicKey()
	rr.Read(req.signature.publicKey)
	req.signature.signature = rr.ReadBytes()
	return rr.Err
}

func (req *OffLedgerRequestData) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	req.writeEssence(ww)
	ww.Write(req.signature.publicKey)
	ww.WriteBytes(req.signature.signature)
	return ww.Err
}

func (req *OffLedgerRequestData) readEssence(rr *rwutil.Reader) {
	rr.ReadKindAndVerify(rwutil.Kind(requestKindOffLedgerISC))
	rr.Read(&req.chainID)
	rr.Read(&req.msg.Target.Contract)
	rr.Read(&req.msg.Target.EntryPoint)
	req.msg.Params = dict.New()
	rr.Read(&req.msg.Params)
	req.nonce = rr.ReadAmount64()
	req.gasBudget = gas.GasUnits(rr.ReadGas64())
	req.allowance = NewEmptyAssets()
	rr.Read(req.allowance)
}

func (req *OffLedgerRequestData) writeEssence(ww *rwutil.Writer) {
	ww.WriteKind(rwutil.Kind(requestKindOffLedgerISC))
	ww.Write(&req.chainID)
	ww.Write(&req.msg.Target.Contract)
	ww.Write(&req.msg.Target.EntryPoint)
	ww.Write(&req.msg.Params)
	ww.WriteAmount64(req.nonce)
	ww.WriteGas64(uint64(req.gasBudget))
	ww.Write(req.allowance)
}

// Allowance from the sender's anchor to the target smart contract. Nil mean no Allowance
func (req *OffLedgerRequestData) Allowance() *Assets {
	return req.allowance
}

// Assets is attached assets to the UTXO. Nil for off-ledger
func (req *OffLedgerRequestData) Assets() *Assets {
	return NewEmptyAssets()
}

func (req *OffLedgerRequestData) Bytes() []byte {
	return rwutil.WriteToBytes(req)
}

func (req *OffLedgerRequestData) Message() Message {
	return req.msg
}

func (req *OffLedgerRequestData) ChainID() ChainID {
	return req.chainID
}

func (req *OffLedgerRequestData) EssenceBytes() []byte {
	ww := rwutil.NewBytesWriter()
	req.writeEssence(ww)
	return ww.Bytes()
}

func (req *OffLedgerRequestData) messageToSign() []byte {
	ret := blake2b.Sum256(req.EssenceBytes())
	return ret[:]
}

func (req *OffLedgerRequestData) Expiry() (iotago.SlotIndex, iotago.Address, bool) {
	return 0, nil, false
}

func (req *OffLedgerRequestData) GasBudget() (gasBudget uint64, isEVM bool) {
	return uint64(req.gasBudget), false
}

// ID returns request id for this request
// index part of request id is always 0 for off ledger requests
// note that request needs to have been signed before this value is
// considered valid
func (req *OffLedgerRequestData) ID() (requestID RequestID) {
	return NewRequestID(iotago.TransactionIDRepresentingData(0, req.Bytes()), 0)
}

func (req *OffLedgerRequestData) IsOffLedger() bool {
	return true
}

func (req *OffLedgerRequestData) NFT() *NFT {
	return nil
}

// Nonce incremental nonce used for replay protection
func (req *OffLedgerRequestData) Nonce() uint64 {
	return req.nonce
}

func (req *OffLedgerRequestData) ReturnAmount() (iotago.BaseToken, bool) {
	return 0, false
}

func (req *OffLedgerRequestData) SenderAccount() AgentID {
	return NewAgentID(req.signature.publicKey.AsEd25519Address())
}

// Sign signs the essence
func (req *OffLedgerRequestData) Sign(key *cryptolib.KeyPair) OffLedgerRequest {
	req.signature = offLedgerSignature{
		publicKey: key.GetPublicKey(),
		signature: key.GetPrivateKey().Sign(req.messageToSign()),
	}
	return req
}

func (req *OffLedgerRequestData) String(bech32HRP iotago.NetworkPrefix) string {
	return fmt.Sprintf("offLedgerRequestData::{ ID: %s, sender: %s, target: %s, entrypoint: %s, Params: %s, nonce: %d }",
		req.ID().String(),
		req.SenderAccount().Bech32(bech32HRP),
		req.msg.Target.Contract.String(),
		req.msg.Target.EntryPoint.String(),
		req.msg.Params.String(),
		req.nonce,
	)
}

func (req *OffLedgerRequestData) TargetAddress() iotago.Address {
	return req.chainID.AsAddress()
}

func (req *OffLedgerRequestData) TimeLock() (iotago.SlotIndex, bool) {
	return 0, false
}

func (req *OffLedgerRequestData) Timestamp() time.Time {
	// no request TX, return zero time
	return time.Time{}
}

// VerifySignature verifies essence signature
func (req *OffLedgerRequestData) VerifySignature() error {
	if !req.signature.publicKey.Verify(req.messageToSign(), req.signature.signature) {
		return errors.New("invalid signature")
	}
	return nil
}

func (req *OffLedgerRequestData) WithAllowance(allowance *Assets) UnsignedOffLedgerRequest {
	req.allowance = allowance.Clone()
	return req
}

func (req *OffLedgerRequestData) WithGasBudget(gasBudget gas.GasUnits) UnsignedOffLedgerRequest {
	req.gasBudget = gasBudget
	return req
}

func (req *OffLedgerRequestData) WithNonce(nonce uint64) UnsignedOffLedgerRequest {
	req.nonce = nonce
	return req
}

// WithSender can be used to estimate gas, without a signature
func (req *OffLedgerRequestData) WithSender(sender *cryptolib.PublicKey) UnsignedOffLedgerRequest {
	req.signature = offLedgerSignature{
		publicKey: sender,
		signature: []byte{},
	}
	return req
}

func (*OffLedgerRequestData) EVMTransaction() *types.Transaction {
	return nil
}
