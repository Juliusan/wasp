package wallet

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/log"
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
)

type WalletConfig struct {
	Seed []byte
}

type Wallet struct {
	KeyPair *cryptolib.KeyPair
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new wallet",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		seed := cryptolib.NewSeed()
		seedString := hexutil.Encode(seed[:])
		viper.Set("wallet.seed", seedString)
		log.Check(viper.WriteConfig())

		log.Printf("Initialized wallet seed in %s\n", config.ConfigPath)
		log.Printf("\nIMPORTANT: wasp-cli is alpha phase. The seed is currently being stored " +
			"in a plain text file which is NOT secure. Do not use this seed to store funds " +
			"in the mainnet!\n")
		log.Verbosef("\nSeed: %s\n", seedString)
	},
}

func Load() *Wallet {
	seedb58 := viper.GetString("wallet.seed")
	if seedb58 == "" {
		log.Fatalf("call `init` first")
	}
	seedBytes, err := hexutil.Decode(seedb58)
	log.Check(err)
	seed := cryptolib.NewSeedFromBytes(seedBytes)
	kp := cryptolib.NewKeyPairFromSeed(seed.SubSeed(uint64(addressIndex)))
	return &Wallet{KeyPair: kp}
}

var addressIndex int

func (w *Wallet) PrivateKey() *cryptolib.PrivateKey {
	return w.KeyPair.GetPrivateKey()
}

func (w *Wallet) Address() iotago.Address {
	return w.KeyPair.GetPublicKey().AsEd25519Address()
}
