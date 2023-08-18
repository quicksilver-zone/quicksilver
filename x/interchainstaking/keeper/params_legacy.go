package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) GetParamLegacy(ctx sdk.Context, key []byte) uint64 {
	var out uint64
	k.paramStore.Get(ctx, key, &out)
	return out
}

func (k *Keeper) GetUnbondingEnabledLegacy(ctx sdk.Context) bool {
	var out bool
	k.paramStore.Get(ctx, types.KeyUnbondingEnabled, &out)
	return out
}

func (k *Keeper) GetCommissionRateLegacy(ctx sdk.Context) sdk.Dec {
	var out sdk.Dec
	k.paramStore.Get(ctx, types.KeyCommissionRate, &out)
	return out
}

// MigrateParams fetches params, adds ClaimsEnabled field and re-sets params.
func (k *Keeper) MigrateParams(ctx sdk.Context) {
	params := types.Params{}
	params.DepositInterval = k.GetParamLegacy(ctx, types.KeyDepositInterval)
	params.CommissionRate = k.GetCommissionRateLegacy(ctx)
	params.ValidatorsetInterval = k.GetParamLegacy(ctx, types.KeyValidatorSetInterval)
	params.UnbondingEnabled = false

	k.paramStore.SetParamSet(ctx, &params)
}
