package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
)

const (
	paramsInvariantName = "params-total-distribution-proportions"
)

// RegisterInvariants registers all participtationrewards invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k *Keeper) {
	ir.RegisterRoute(types.ModuleName, paramsInvariantName, ParamsInvariant(k))
}

// AllInvariants runs all invariants of the module.
func AllInvariants(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg, broke := ParamsInvariant(k)(ctx)
		return msg, broke
	}
}

func ParamsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		params := k.GetParams(ctx)
		if !params.DistributionProportions.Total().Equal(sdkmath.LegacyOneDec()) {
			return sdk.FormatInvariant(
				types.ModuleName,
				paramsInvariantName,
				"\tdistribution total proportions are not equal to 1.0",
			), true
		}
		return sdk.FormatInvariant(
			types.ModuleName,
			paramsInvariantName,
			"\tdistribution total proportions are equal to 1.0",
		), false
	}
}
