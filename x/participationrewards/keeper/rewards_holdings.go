package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k Keeper) AllocateHoldingsRewards(ctx sdk.Context) error {
	k.Logger(ctx).Info("allocateHoldingsRewards")

	// obtain and iterate all claim records for each zone
	for i, zone := range k.icsKeeper.AllZones(ctx) {
		k.Logger(ctx).Info("zones", "i", i, "zone", zone.ChainId)
		userAllocations := k.CalcUserHoldingsAllocations(ctx, zone)

		if err := k.distributeToUsers(ctx, userAllocations); err != nil {
			// we might want to do a soft fail here so that all zones are not affected...
			return err
		}

		k.icsKeeper.ClaimsManagerKeeper.ArchiveAndGarbageCollectClaims(ctx, zone.ChainId)
	}

	return nil
}

func (k Keeper) CalcUserHoldingsAllocations(ctx sdk.Context, zone icstypes.Zone) []userAllocation {
	k.Logger(ctx).Info("calcUserHoldingsAllocations", "zone", zone.ChainId, "allocations", zone.HoldingsAllocation)

	userAllocations := make([]userAllocation, 0)

	if zone.HoldingsAllocation == 0 {
		k.Logger(ctx).Info("holdings allocation is zero, nothing to allocate")
		return userAllocations
	}

	// calculate user totals and zone total (held assets)
	zoneAmount := math.ZeroInt()
	userAmountsMap := make(map[string]math.Int)

	k.icsKeeper.ClaimsManagerKeeper.IterateClaims(ctx, zone.ChainId, func(_ int64, claim cmtypes.Claim) (stop bool) {
		// we can suppress the error here as the address is from claim
		// state that is verified.
		// userAccount, _ := sdk.AccAddressFromBech32(claim.UserAddress)
		// calculate user held amount
		// total = local + remote
		// local amount here uses the current epoch balance which is not aligned
		// with claims that are against the previous epoch
		// local := k.bankKeeper.GetBalance(ctx, userAccount, zone.LocalDenom).Amount
		// remote := math.NewIntFromUint64(claim.Amount)
		// total := local.Add(remote)

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
		return userAllocations
	}

	// calculate user held proportions and apply limit
	limit := sdk.MustNewDecFromStr("0.02")
	adjustedZoneAmount := math.ZeroInt()
	for _, address := range utils.Keys(userAmountsMap) {
		userAmount := userAmountsMap[address]
		userPortion := sdk.NewDecFromInt(userAmount).Quo(sdk.NewDecFromInt(zoneAmount))
		// check for and apply limit
		if userPortion.GT(limit) {
			userAmount = sdk.NewDecFromInt(zoneAmount).Mul(limit).TruncateInt()
			userAmountsMap[address] = userAmount
		}
		adjustedZoneAmount = adjustedZoneAmount.Add(userAmount)
	}
	k.Logger(ctx).Info("rewards limit adjustment", "zoneAmount", zoneAmount, "adjustedZoneAmount", adjustedZoneAmount)

	allocation := sdk.NewDecFromInt(math.NewIntFromUint64(zone.HoldingsAllocation))
	tokensPerAsset := allocation.Quo(sdk.NewDecFromInt(adjustedZoneAmount))
	k.Logger(ctx).Info("tokens per asset", "zone", zone.ChainId, "tpa", tokensPerAsset)

	for _, address := range utils.Keys(userAmountsMap) {
		amount := userAmountsMap[address]
		allocation := userAllocation{
			Address: address,
			Amount:  sdk.NewDecFromInt(amount).Mul(tokensPerAsset).TruncateInt(),
		}
		userAllocations = append(userAllocations, allocation)
	}

	return userAllocations
}
