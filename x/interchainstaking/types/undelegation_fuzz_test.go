package types_test

import (
	"encoding/json"
	"testing"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func FuzzDetermineAllocationsForUndelegation(f *testing.F) {
	if testing.Short() {
		f.Skip("Running in -short mode")
	}

	seeds := detRedemptionTests(vals)
	for _, seed := range seeds {
		blob, err := json.Marshal(seed)
		if err == nil {
			f.Add(blob)
		}
	}

	f.Fuzz(func(_ *testing.T, inputJSON []byte) {
		drt := new(detRedemptionTest)
		if err := json.Unmarshal(inputJSON, drt); err != nil {
			return
		}
		_, _ = types.DetermineAllocationsForUndelegation(
			drt.CurrentAllocations,
			map[string]bool{},
			sum(drt.CurrentAllocations),
			drt.TargetAllocations,
			drt.Unlocked,
			drt.Amount)
	})
}

func FuzzDetermineAllocationsForUndelegationPredef(f *testing.F) {
	if testing.Short() {
		f.Skip("Running in -short mode")
	}

	seeds := detRedemptionTests(vals)
	for _, seed := range seeds {
		blob, err := json.Marshal(seed)
		if err == nil {
			f.Add(blob)
		}
	}

	f.Fuzz(func(_ *testing.T, inputJSON []byte) {
		drt := new(detRedemptionTest)
		if err := json.Unmarshal(inputJSON, drt); err != nil {
			return
		}
		_, _ = types.DetermineAllocationsForUndelegationPredef(
			drt.CurrentAllocations,
			map[string]bool{},
			sum(drt.CurrentAllocations),
			drt.TargetAllocations,
			drt.Unlocked,
			drt.Amount)
	})
}
