package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// userAmount is an internal struct to track transient state for rewards
// distribution. It contains the user address and held amount.
type userAmount struct {
	Address string
	Amount  sdk.Int
}

func (k Keeper) allocateHoldingsRewards(ctx sdk.Context) error {
	k.Logger(ctx).Info("allocateHoldingsRewards")

	// obtain and iterate all claim records for each zone
	for i, zone := range k.icsKeeper.AllZones(ctx) {
		k.Logger(ctx).Info("zones", "i", i, "zone", zone.ChainId)
		userAllocations := k.calcUserHoldingsAllocations(ctx, zone)

		if err := k.distributeToUsers(ctx, userAllocations); err != nil {
			// we might want to do a soft fail here so that all zones are not affected...
			return err
		}

		k.ClearClaims(ctx, zone.ChainId)
	}

	return nil
}

func (k Keeper) calcUserHoldingsAllocations(ctx sdk.Context, zone icstypes.Zone) []userAllocation {
	k.Logger(ctx).Info("calcUserHoldingsAllocations", "zone", zone.ChainId, "allocations", zone.HoldingsAllocation)

	userAllocations := make([]userAllocation, 0)

	if zone.HoldingsAllocation.IsZero() {
		k.Logger(ctx).Info("holdings allocation is zero, nothing to allocate")
		return userAllocations
	}

	// get zone claims
	claims := k.AllZoneClaims(ctx, zone.ChainId)

	// calculate user totals and zone total (held assets)
	zoneAmount := sdk.ZeroInt()
	userAmounts := make([]userAmount, len(claims))
	for i, claim := range claims {
		// we can suppress the error here as the address is from claim
		// state that is verified.
		userAccount, _ := sdk.AccAddressFromBech32(claim.UserAddress)
		// calculate user held amount
		// total = local + remote
		local := k.bankKeeper.GetBalance(ctx, userAccount, zone.LocalDenom).Amount
		remote := sdk.NewIntFromUint64(claim.Amount)
		total := local.Add(remote)
		k.Logger(ctx).Info("user amount for zone", "user", claim.UserAddress, "zone", claim.ChainId, "held", total)
		userAmounts[i] = userAmount{
			Address: claim.UserAddress,
			Amount:  total,
		}

		zoneAmount = zoneAmount.Add(total)
	}

	if zoneAmount.IsZero() {
		k.Logger(ctx).Info("zero claims for zone", "zone", zone.ChainId)
		return userAllocations
	}

	// calculate user held proportions and apply limit
	limit := sdk.MustNewDecFromStr("0.02")
	adjustedZoneAmount := sdk.ZeroInt()
	for i, userAmount := range userAmounts {
		userPortion := userAmount.Amount.ToDec().Quo(zoneAmount.ToDec())
		// check for and apply limit
		if userPortion.GT(limit) {
			userAmount.Amount = zoneAmount.ToDec().Mul(limit).TruncateInt()
			userAmounts[i] = userAmount
		}
		adjustedZoneAmount = adjustedZoneAmount.Add(userAmount.Amount)
	}
	k.Logger(ctx).Info("rewards limit adjustment", "zoneAmount", zoneAmount, "adjustedZoneAmount", adjustedZoneAmount)

	tokensPerAsset := zone.HoldingsAllocation.AmountOfNoDenomValidation(k.stakingKeeper.BondDenom(ctx)).ToDec().Quo(adjustedZoneAmount.ToDec())
	k.Logger(ctx).Info("tokens per asset", "zone", zone.ChainId, "tpa", tokensPerAsset)

	for _, ua := range userAmounts {
		allocation := userAllocation{
			Address: ua.Address,
			Coins: sdk.NewCoins(
				sdk.NewCoin(
					k.stakingKeeper.BondDenom(ctx),
					ua.Amount.ToDec().Mul(tokensPerAsset).TruncateInt(),
				),
			),
		}
		userAllocations = append(userAllocations, allocation)
	}

	return userAllocations
}
