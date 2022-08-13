package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

// CreateZoneDropAccount creates a zone specific module account.
func (k Keeper) CreateZoneDropAccount(ctx sdk.Context, chainID string) {
	name := types.ModuleName + "." + chainID
	moduleAcc := authtypes.NewEmptyModuleAccount(name, "")
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
}

// GetZoneDropAccountAddress returns the zone airdrop account address.
func (k Keeper) GetZoneDropAccountAddress(ctx sdk.Context, chainID string) sdk.AccAddress {
	name := types.ModuleName + "." + chainID
	return k.accountKeeper.GetModuleAddress(name)
}

// GetZoneDropAccountBalance gets the zone airdrop account coin balance.
func (k Keeper) GetZoneDropAccountBalance(ctx sdk.Context, chainID string) sdk.Coin {
	zonedropAccAddr := k.GetZoneDropAccountAddress(ctx, chainID)
	return k.bankKeeper.GetBalance(ctx, zonedropAccAddr, k.stakingKeeper.BondDenom(ctx))
}

// GetZoneDrop returns airdrop details for the zone identified by chainID.
func (k Keeper) GetZoneDrop(ctx sdk.Context, chainID string) (types.ZoneDrop, bool) {
	zd := types.ZoneDrop{}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetKeyZoneDrop(chainID))
	if len(b) == 0 {
		return zd, false
	}

	k.cdc.MustUnmarshal(b, &zd)
	return zd, true
}

// SetZoneDrop creates/updates the given zone airdrop (ZoneDrop).
func (k Keeper) SetZoneDrop(ctx sdk.Context, zd types.ZoneDrop) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&zd)
	store.Set(types.GetKeyZoneDrop(zd.ChainId), b)
}

// DeleteZoneDrop deletes the airdrop of the zone identified by chainID.
func (k Keeper) DeleteZoneDrop(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetKeyZoneDrop(chainID))
}

// IterateZoneDrops iterate through zone airdrops.
func (k Keeper) IterateZoneDrops(ctx sdk.Context, fn func(index int64, zoneInfo types.ZoneDrop) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixZoneDrop)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		zd := types.ZoneDrop{}
		k.cdc.MustUnmarshal(iterator.Value(), &zd)

		stop := fn(i, zd)

		if stop {
			break
		}
		i++
	}
}

// AllZoneDrops returns all zone airdrops (active, future, expired).
func (k Keeper) AllZoneDrops(ctx sdk.Context) []*types.ZoneDrop {
	zds := []*types.ZoneDrop{}
	k.IterateZoneDrops(ctx, func(_ int64, zd types.ZoneDrop) (stop bool) {
		zds = append(zds, &zd)
		return false
	})
	return zds
}

// AllActiveZoneDrops returns all active zone airdrops.
func (k Keeper) AllActiveZoneDrops(ctx sdk.Context) []types.ZoneDrop {
	zds := []types.ZoneDrop{}
	k.IterateZoneDrops(ctx, func(_ int64, zd types.ZoneDrop) (stop bool) {
		if k.IsActiveZoneDrop(ctx, zd) {
			zds = append(zds, zd)
		}
		return false
	})
	return zds
}

// IsActiveZoneDrop returns true if the zone airdrop is currently active.
func (k Keeper) IsActiveZoneDrop(ctx sdk.Context, zd types.ZoneDrop) bool {
	bt := ctx.BlockTime()

	// Zone airdrop has not yet started
	if bt.Before(zd.StartTime) {
		return false
	}

	// Zone airdrop has expired
	return !bt.After(zd.StartTime.Add(zd.Duration).Add(zd.Decay))
}

// AllFutureZoneDrops returns all future zone airdrops.
func (k Keeper) AllFutureZoneDrops(ctx sdk.Context) []types.ZoneDrop {
	zds := []types.ZoneDrop{}
	k.IterateZoneDrops(ctx, func(_ int64, zd types.ZoneDrop) (stop bool) {
		if k.IsFutureZoneDrop(ctx, zd) {
			zds = append(zds, zd)
		}
		return false
	})
	return zds
}

// IsFutureZoneDrop returns true if the zone airdrop is in the future.
func (k Keeper) IsFutureZoneDrop(ctx sdk.Context, zd types.ZoneDrop) bool {
	bt := ctx.BlockTime()

	// Zone airdrop has already started
	return !bt.After(zd.StartTime)
}

// AllExpiredZoneDrops returns all expired zone airdrops.
func (k Keeper) AllExpiredZoneDrops(ctx sdk.Context) []types.ZoneDrop {
	zds := []types.ZoneDrop{}
	k.IterateZoneDrops(ctx, func(_ int64, zd types.ZoneDrop) (stop bool) {
		if k.IsExpiredZoneDrop(ctx, zd) {
			zds = append(zds, zd)
		}
		return false
	})
	return zds
}

// IsExpiredZoneDrop returns true if the zone airdrop has already expired.
func (k Keeper) IsExpiredZoneDrop(ctx sdk.Context, zd types.ZoneDrop) bool {
	bt := ctx.BlockTime()

	// Zone airdrop has not yet expired
	return !bt.Before(zd.StartTime.Add(zd.Duration).Add(zd.Decay))
}

// UnconcludedAirdrops returns all expired zone airdrops that have not yet been
// concluded.
func (k Keeper) UnconcludedAirdrops(ctx sdk.Context) []types.ZoneDrop {
	zds := []types.ZoneDrop{}
	k.IterateZoneDrops(ctx, func(_ int64, zd types.ZoneDrop) (stop bool) {
		if k.IsExpiredZoneDrop(ctx, zd) {
			if !zd.IsConcluded {
				zds = append(zds, zd)
			}
		}
		return false
	})
	return zds
}

// EndZoneDrop concludes a zone airdrop. It deletes all ClaimRecords for the
// given zone.
func (k Keeper) EndZoneDrop(ctx sdk.Context, chainID string) error {
	if err := k.returnUnclaimedZoneDropTokens(ctx, chainID); err != nil {
		return err
	}
	k.ClearClaimRecords(ctx, chainID)

	zd, ok := k.GetZoneDrop(ctx, chainID)
	if !ok {
		return types.ErrZoneDropNotFound
	}

	zd.IsConcluded = true
	k.SetZoneDrop(ctx, zd)

	return nil
}

// returnUnclaimedZoneDropTokens returns all unclaimed zone airdrop tokens to
// the airdrop module account.
func (k Keeper) returnUnclaimedZoneDropTokens(ctx sdk.Context, chainID string) error {
	zonedropAccountAddress := k.GetZoneDropAccountAddress(ctx, chainID)
	zonedropAccountBalance := k.GetZoneDropAccountBalance(ctx, chainID)
	airdropAccountAddress := k.GetModuleAccountAddress(ctx)
	return k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		zonedropAccountAddress.String(),
		airdropAccountAddress.String(),
		sdk.NewCoins(zonedropAccountBalance),
	)
}
