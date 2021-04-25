package core

import (
	"bytes"
	"fmt"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"sort"
)

func init() {
	//printWellKnownHnames()
}

// for debugging
func PrintWellKnownHnames() {
	fmt.Printf("--------------- well known hnames ------------------\n")
	hashes := make([]hashing.HashValue, 0)
	for _, rec := range AllCoreContracts {
		hashes = append(hashes, rec.ProgramHash)
	}
	sort.Slice(hashes, func(i, j int) bool {
		return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
	})
	for _, h := range hashes {
		rec := AllCoreContracts[h]
		fmt.Printf("    %10d, %10s: '%s'\n", rec.Hname(), rec.Hname(), rec.Name)
	}
	fmt.Printf("    %10d, %10s: '%s'\n", coretypes.EntryPointInit, coretypes.EntryPointInit, coretypes.FuncInit)
	fmt.Printf("    %10d, %10s: '%s'\n", coretypes.Hn("test_sandbox"), coretypes.Hn("test_sandbox"), "test_sandbox")
	fmt.Printf("--------------- well known hnames ------------------\n")
}
