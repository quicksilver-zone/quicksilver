package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	supplytypes "github.com/quicksilver-zone/quicksilver/x/supply/types"
)

func (k Keeper) AllocateHoldingsRewards(ctx sdk.Context) error {
	// obtain and iterate all claim records for each zone
	k.icsKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		k.Logger(ctx).Info("zones", "zone", zone.ChainId)
		userAllocations, remaining := k.CalcUserHoldingsAllocations(ctx, zone)

		if err := k.DistributeToUsersFromModule(ctx, userAllocations); err != nil {
			k.Logger(ctx).Error("failed to distribute to users", "ua", userAllocations, "err", err)
			return false
		}

		if remaining.IsPositive() {
			k.Logger(ctx).Error("remaining amount to return to incentives pool", "remainder", remaining, "pool balance", k.GetModuleBalance(ctx))
			// send unclaimed remainder to incentives pool
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, supplytypes.AirdropAccount, sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), remaining))); err != nil {
				k.Logger(ctx).Error("failed to send remaining amount to return to incentives pool", "remainder", remaining, "pool balance", k.GetModuleBalance(ctx), "err", err)
				return false
			}
		}

		return false
	})

	return nil
}

// CalcUserHoldingsAllocations calculates allocations per user for a given zone, based upon claims submitted and zone.
func (k Keeper) CalcUserHoldingsAllocations(ctx sdk.Context, zone *icstypes.Zone) ([]types.UserAllocation, math.Int) {
	k.Logger(ctx).Info("CalcUserHoldingsAllocations", "zone", zone.ChainId, "allocations", zone.HoldingsAllocation)

	userAllocations := make([]types.UserAllocation, 0)

	supply := k.bankKeeper.GetSupply(ctx, zone.LocalDenom)

	if zone.HoldingsAllocation == 0 || !supply.Amount.IsPositive() {
		k.Logger(ctx).Info("holdings allocation is zero, nothing to allocate")
		return userAllocations, math.NewIntFromUint64(zone.HoldingsAllocation)
	}

	// calculate user totals and zone total (held assets)
	zoneAmount := math.ZeroInt()
	userAmountsMap := make(map[string]math.Int)

	k.ClaimsManagerKeeper.IterateLastEpochClaims(ctx, zone.ChainId, func(_ int64, claim cmtypes.Claim) (stop bool) {
		k.Logger(ctx).Info(
			"claim",
			"type", cmtypes.ClaimType_name[int32(claim.Module)],
			"user", claim.UserAddress,
			"zone", claim.ChainId,
			"amount", claim.Amount,
		)

		if _, exists := userAmountsMap[claim.UserAddress]; !exists {
			userAmountsMap[claim.UserAddress] = math.ZeroInt()
		}

		userAmountsMap[claim.UserAddress] = userAmountsMap[claim.UserAddress].Add(claim.Amount)

		// total zone assets held remotely
		zoneAmount = zoneAmount.Add(claim.Amount)

		return false
	})

	if !zoneAmount.IsPositive() {
		k.Logger(ctx).Info("zero claims for zone", "zone", zone.ChainId)
		return userAllocations, math.NewIntFromUint64(zone.HoldingsAllocation)
	}

	zoneAllocation := math.NewIntFromUint64(zone.HoldingsAllocation)
	tokensPerAsset := sdk.NewDecFromInt(zoneAllocation).Quo(sdk.NewDecFromInt(supply.Amount))

	k.Logger(ctx).Info("tokens per asset", "zone", zone.ChainId, "tpa", tokensPerAsset)

	for _, address := range utils.Keys(userAmountsMap) {
		amount := userAmountsMap[address]
		userAllocation := sdk.NewDecFromInt(amount).Mul(tokensPerAsset).TruncateInt()
		allocation := types.UserAllocation{
			Address: address,
			Amount:  sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), userAllocation),
		}
		userAllocations = append(userAllocations, allocation)
		zoneAllocation = zoneAllocation.Sub(userAllocation)
		if zoneAllocation.LT(sdk.ZeroInt()) {
			panic("user allocation overflow")
		}

	}

	return userAllocations, zoneAllocation
}
