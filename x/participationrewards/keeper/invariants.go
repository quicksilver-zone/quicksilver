package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

const (
	paramsInvariantName = "params-total-distribution-proportions"
)

// RegisterInvariants registers all airdrop invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, paramsInvariantName, ParamsInvariant(k))
}

// AllInvariants runs all invariants of the module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg, broke := ParamsInvariant(k)(ctx)
		return msg, broke
	}
}

func ParamsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		params := k.GetParams(ctx)
		if !params.DistributionProportions.Total().Equal(sdk.OneDec()) {
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
