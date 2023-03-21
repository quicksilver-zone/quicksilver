package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k Keeper) AllocateHoldingsRewards(ctx sdk.Context) error {
	k.Logger(ctx).Info("allocateHoldingsRewards")

	// obtain and iterate all claim records for each zone
	k.icsKeeper.IterateZones(ctx, func(_ int64, zone *icstypes.Zone) (stop bool) {
		k.Logger(ctx).Info("zones", "zone", zone.ChainId)
		userAllocations, remaining := k.CalcUserHoldingsAllocations(ctx, zone)

		if err := k.distributeToUsers(ctx, userAllocations); err != nil {
			k.Logger(ctx).Error("failed to distribute to users", "ua", userAllocations, "err", err)
			// we might want to do a soft fail here so that all zones are not affected...
			return false
		}

		if remaining.IsPositive() {
			k.Logger(ctx).Error("remaining amount to return to incentives pool", "remainder", remaining, "pool balance", k.GetModuleBalance(ctx))
			// send unclaimed remainder to incentives pool
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, airdroptypes.ModuleName, sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), remaining))); err != nil {
				k.Logger(ctx).Error("failed to send remaining amount to return to incentives pool", "remainder", remaining, "pool balance", k.GetModuleBalance(ctx), "err", err)
				return false
			}
		}

		k.icsKeeper.ClaimsManagerKeeper.ArchiveAndGarbageCollectClaims(ctx, zone.ChainId)
		return false
	})

	return nil
}

// calculate allocations per user for a given zone, based upon claims submitted and zone
func (k Keeper) CalcUserHoldingsAllocations(ctx sdk.Context, zone *icstypes.Zone) ([]UserAllocation, math.Int) {
	k.Logger(ctx).Info("calcUserHoldingsAllocations", "zone", zone.ChainId, "allocations", zone.HoldingsAllocation)

	userAllocations := make([]UserAllocation, 0)
	supply := k.bankKeeper.GetSupply(ctx, zone.LocalDenom)

	if zone.HoldingsAllocation == 0 || supply.Amount.IsZero() {
		k.Logger(ctx).Info("holdings allocation is zero, nothing to allocate")
		return userAllocations, math.NewIntFromUint64(zone.HoldingsAllocation)
	}

	// calculate user totals and zone total (held assets)
	zoneAmount := math.ZeroInt()
	userAmountsMap := make(map[string]math.Int)

	k.icsKeeper.ClaimsManagerKeeper.IterateClaims(ctx, zone.ChainId, func(_ int64, claim cmtypes.Claim) (stop bool) {
		amount := math.NewIntFromUint64(claim.Amount)
		k.Logger(ctx).Info(
			"claim",
			"type", cmtypes.ClaimType_name[int32(claim.Module)],
			"user", claim.UserAddress,
			"zone", claim.ChainId,
			"amount", amount,
		)

		if _, exists := userAmountsMap[claim.UserAddress]; !exists {
			userAmountsMap[claim.UserAddress] = math.ZeroInt()
		}

		userAmountsMap[claim.UserAddress] = userAmountsMap[claim.UserAddress].Add(amount)

		// total zone assets held remotely
		zoneAmount = zoneAmount.Add(amount)

		return false
	})

	if zoneAmount.IsZero() {
		k.Logger(ctx).Info("zero claims for zone", "zone", zone.ChainId)
		return userAllocations, math.NewIntFromUint64(zone.HoldingsAllocation)
	}

	zoneAllocation := math.NewIntFromUint64(zone.HoldingsAllocation)
	tokensPerAsset := sdk.NewDecFromInt(zoneAllocation).Quo(sdk.NewDecFromInt(supply.Amount))

	k.Logger(ctx).Info("tokens per asset", "zone", zone.ChainId, "tpa", tokensPerAsset)

	for _, address := range utils.Keys(userAmountsMap) {
		amount := userAmountsMap[address]
		userAllocation := sdk.NewDecFromInt(amount).Mul(tokensPerAsset).TruncateInt()
		allocation := UserAllocation{
			Address: address,
			Amount:  userAllocation,
		}
		userAllocations = append(userAllocations, allocation)
		zoneAllocation = zoneAllocation.Sub(userAllocation)
		if zoneAllocation.LT(sdk.ZeroInt()) {
			panic("user allocation overflow")
		}
	}

	return userAllocations, zoneAllocation
}
