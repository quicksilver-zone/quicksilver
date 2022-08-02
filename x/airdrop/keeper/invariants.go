package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

const (
	claimRecordInvariantName = "claim-record-max-allocation"
)

// RegisterInvariants registers all airdrop invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, claimRecordInvariantName, ClaimRecordInvariant(k))
}

// AllInvariants runs all invariants of the module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg, broke := ClaimRecordInvariant(k)(ctx)
		return msg, broke
	}
}

func ClaimRecordInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		crs := k.AllClaimRecords(ctx)
		for _, cr := range crs {
			sum := uint64(0)
			for _, ca := range cr.ActionCompleted {
				sum += ca.ClaimAmount
			}
			if cr.MaxAllocation < sum {
				return sdk.FormatInvariant(
					types.ModuleName,
					claimRecordInvariantName,
					"\tclaim record completed actions exceed max allocation",
				), true
			}
		}
		return sdk.FormatInvariant(
			types.ModuleName,
			claimRecordInvariantName,
			"\tall claim records completed actions are less than or equal to max allocations",
		), false
	}
}
