package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) AllocateHoldingsRewards(ctx sdk.Context) error {
	// obtain and iterate all claim records for each zone
	k.icsKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		k.Logger(ctx).Info("zones", "zone", zone.ChainId)

		customeZone := cmtypes.CustomeZone{
			ChainId:            zone.ChainId,
			HoldingsAllocation: zone.HoldingsAllocation,
			LocalDenom:         zone.LocalDenom,
		}
		userAllocations, remaining := k.icsKeeper.ClaimsManagerKeeper.CalcUserHoldingsAllocations(ctx, customeZone)

		if err := k.DistributeToUsers(ctx, userAllocations); err != nil {
			k.Logger(ctx).Error("failed to distribute to users", "ua", userAllocations, "err", err)
			// we might want to do a soft fail here so that all zones are not affected...
			return false
		}

		if remaining.IsPositive() {
			k.Logger(ctx).Error("remaining amount to return to incentives pool", "remainder", remaining, "pool balance", k.GetModuleBalance(ctx))
			// send unclaimed remainder to incentives pool
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, airdroptypes.ModuleName, sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), remaining))); err != nil {
				k.Logger(ctx).Error("failed to send remaining amount to return to incentives pool", "remainder", remaining, "pool balance", k.GetModuleBalance(ctx), "err", err)
				return false
			}
		}

		k.icsKeeper.ClaimsManagerKeeper.ArchiveAndGarbageCollectClaims(ctx, zone.ChainId)
		return false
	})

	return nil
}
