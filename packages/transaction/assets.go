package transaction

import (
	"fmt"
	"math/big"

	"github.com/samber/lo"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/vm"
	"github.com/iotaledger/wasp/packages/isc"
)

type AssetsWithMana struct {
	isc.Assets
	Mana iotago.Mana
}

func NewAssetsWithMana(assets *isc.Assets, mana iotago.Mana) *AssetsWithMana {
	return &AssetsWithMana{Assets: *assets, Mana: mana}
}

func NewEmptyAssetsWithMana() *AssetsWithMana {
	return NewAssetsWithMana(isc.NewEmptyAssets(), 0)
}

func (a *AssetsWithMana) String() string {
	ret := a.Assets.String()
	if a.Mana > 0 {
		ret += fmt.Sprintf("\n Mana: %d", a.Mana)
	}
	return ret
}

func (a *AssetsWithMana) Geq(b *AssetsWithMana) bool {
	if !a.Assets.Geq(&b.Assets) {
		return false
	}
	return a.Mana > b.Mana
}

func (a *AssetsWithMana) Equals(b *AssetsWithMana) bool {
	return a.Assets.Equals(&b.Assets) && a.Mana == b.Mana
}

func (a *AssetsWithMana) Add(b *AssetsWithMana) {
	a.Assets.Add(&b.Assets)
	a.Mana += b.Mana
}

func (a *AssetsWithMana) Clone() *AssetsWithMana {
	return &AssetsWithMana{
		Assets: *a.Assets.Clone(),
		Mana:   a.Mana,
	}
}

func MustSingleNativeToken(a *isc.FungibleTokens) (iotago.NativeTokenID, *big.Int, bool) {
	if len(a.NativeTokens) > 1 {
		panic("expected at most 1 native token")
	}
	for id, n := range a.NativeTokens {
		return id, n, true
	}
	return iotago.NativeTokenID{}, nil, false
}

func AssetsAndAvailableManaFromOutput(
	oID iotago.OutputID,
	o iotago.Output,
	slotIndex iotago.SlotIndex,
	l1 iotago.APIProvider,
) (*AssetsWithMana, error) {
	l1API := l1.APIForSlot(slotIndex)
	assets := isc.AssetsFromOutput(o, oID)
	mana, err := vm.TotalManaIn(
		l1API.ManaDecayProvider(),
		l1API.StorageScoreStructure(),
		slotIndex,
		vm.InputSet{oID: o},
		vm.RewardsInputSet{},
	)
	if err != nil {
		return nil, err
	}
	return NewAssetsWithMana(assets, mana), nil
}

func AdjustToMinimumStorageDeposit[T iotago.Output](out T, l1API iotago.API) T {
	storageDeposit := lo.Must(l1API.StorageScoreStructure().MinDeposit(out))
	if out.BaseTokenAmount() >= storageDeposit {
		return out
	}
	switch out := iotago.Output(out).(type) {
	case *iotago.AnchorOutput:
		out.Amount = storageDeposit
	case *iotago.BasicOutput:
		out.Amount = storageDeposit
	case *iotago.FoundryOutput:
		out.Amount = storageDeposit
	case *iotago.NFTOutput:
		out.Amount = storageDeposit
	default:
		panic(fmt.Sprintf("no handler for output type %T", out))
	}
	return out
}