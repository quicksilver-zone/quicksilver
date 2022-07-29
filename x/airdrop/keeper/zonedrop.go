package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

// CreateZoneDropAccount creates a zone specific module account.
func (k Keeper) CreateZoneDropAccount(ctx sdk.Context, chainId string) {
	name := types.ModuleName + "." + chainId
	moduleAcc := authtypes.NewEmptyModuleAccount(name, "")
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
}

// GetZoneDropAccountAddress returns the zone airdrop account address.
func (k Keeper) GetZoneDropAccountAddress(ctx sdk.Context, chainId string) sdk.AccAddress {
	name := types.ModuleName + "." + chainId
	return k.accountKeeper.GetModuleAddress(name)
}

// GetZoneDropAccountBalance gets the zone airdrop account coin balance.
func (k Keeper) GetZoneDropAccountBalance(ctx sdk.Context, chainId string) sdk.Coin {
	zonedropAccAddr := k.GetZoneDropAccountAddress(ctx, chainId)
	return k.BankKeeper.GetBalance(ctx, zonedropAccAddr, k.StakingKeeper.BondDenom(ctx))
}

// GetZoneDrop returns airdrop details for the zone identified by chainId.
func (k Keeper) GetZoneDrop(ctx sdk.Context, chainId string) (types.ZoneDrop, bool) {
	zd := types.ZoneDrop{}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetKeyZoneDrop(chainId))
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

// DeleteZoneDrop deletes the airdrop of the zone identified by chainId.
func (k Keeper) DeleteZoneDrop(ctx sdk.Context, chainId string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetKeyZoneDrop(chainId))
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
	if bt.After(zd.StartTime.Add(zd.Duration).Add(zd.Decay)) {
		return false
	}

	return true
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
	if bt.After(zd.StartTime) {
		return false
	}

	return true
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
	if bt.Before(zd.StartTime.Add(zd.Duration).Add(zd.Decay)) {
		return false
	}

	return true
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
func (k Keeper) EndZoneDrop(ctx sdk.Context, chainId string) error {
	if err := k.returnUnclaimedZoneDropTokens(ctx, chainId); err != nil {
		return err
	}
	k.ClearClaimRecords(ctx, chainId)

	zd, ok := k.GetZoneDrop(ctx, chainId)
	if !ok {
		return types.ErrZoneDropNotFound
	}

	zd.IsConcluded = true
	k.SetZoneDrop(ctx, zd)

	return nil
}

// returnUnclaimedZoneDropTokens returns all unclaimed zone airdrop tokens to
// the airdrop module account.
func (k Keeper) returnUnclaimedZoneDropTokens(ctx sdk.Context, chainId string) error {
	zonedropAccountAddress := k.GetZoneDropAccountAddress(ctx, chainId)
	zonedropAccountBalance := k.GetZoneDropAccountBalance(ctx, chainId)
	airdropAccountAddress := k.GetModuleAccountAddress(ctx)
	return k.BankKeeper.SendCoinsFromModuleToModule(
		ctx,
		zonedropAccountAddress.String(),
		airdropAccountAddress.String(),
		sdk.NewCoins(zonedropAccountBalance),
	)
}
