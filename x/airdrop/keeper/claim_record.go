package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

// GetClaimRecord returns the ClaimRecord of the given address for the given zone.
func (k Keeper) GetClaimRecord(ctx sdk.Context, chainID string, address string) (types.ClaimRecord, error) {
	cr := types.ClaimRecord{}

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return cr, err
	}

	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetKeyClaimRecord(chainID, addr))
	if len(b) == 0 {
		return cr, types.ErrClaimRecordNotFound
	}

	k.cdc.MustUnmarshal(b, &cr)
	return cr, nil
}

// SetClaimRecord creates/updates the given airdrop ClaimRecord.
func (k Keeper) SetClaimRecord(ctx sdk.Context, cr types.ClaimRecord) error {
	addr, err := sdk.AccAddressFromBech32(cr.Address)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&cr)
	store.Set(types.GetKeyClaimRecord(cr.ChainId, addr), b)

	return nil
}

// DeleteClaimRecord deletes the airdrop ClaimRecord of the given zone and address.
func (k Keeper) DeleteClaimRecord(ctx sdk.Context, chainID string, address string) error {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetKeyClaimRecord(chainID, addr))

	return nil
}

// IterateClaimRecords iterate through zone airdrop ClaimRecords.
func (k Keeper) IterateClaimRecords(ctx sdk.Context, chainID string, fn func(index int64, cr types.ClaimRecord) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixClaimRecord(chainID))
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		cr := types.ClaimRecord{}
		k.cdc.MustUnmarshal(iterator.Value(), &cr)

		stop := fn(i, cr)

		if stop {
			break
		}
		i++
	}
}

// AllClaimRecords returns all the claim records.
func (k Keeper) AllClaimRecords(ctx sdk.Context) []*types.ClaimRecord {
	crs := []*types.ClaimRecord{}

	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixClaimRecord)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		cr := types.ClaimRecord{}
		k.cdc.MustUnmarshal(iterator.Value(), &cr)

		crs = append(crs, &cr)
	}

	return crs
}

// AllZoneClaimRecords returns all the claim records of the given zone.
func (k Keeper) AllZoneClaimRecords(ctx sdk.Context, chainID string) []*types.ClaimRecord {
	crs := []*types.ClaimRecord{}
	k.IterateClaimRecords(ctx, chainID, func(_ int64, cr types.ClaimRecord) (stop bool) {
		crs = append(crs, &cr)
		return false
	})
	return crs
}

// ClearClaimRecords deletes all the claim records of the given zone.
func (k Keeper) ClearClaimRecords(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetPrefixClaimRecord(chainID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

// GetClaimableAmountForAction returns the amount claimable for the given
// action, by the given address, against the given zone.
func (k Keeper) GetClaimableAmountForAction(ctx sdk.Context, chainID string, address string, action types.Action) (uint64, error) {
	cr, err := k.GetClaimRecord(ctx, chainID, address)
	if err != nil {
		return 0, err
	}
	if cr.Address == "" {
		return 0, nil
	}

	// action completed, nothing to claim
	if _, exists := cr.ActionCompleted[int32(action)]; exists {
		return 0, nil
	}

	// get zone airdrop details
	zd, ok := k.GetZoneDrop(ctx, chainID)
	if !ok {
		return 0, types.ErrZoneDropNotFound
	}

	// airdrop has not started
	if ctx.BlockTime().Before(zd.StartTime) {
		return 0, nil
	}

	// airdrop has already ended
	if ctx.BlockTime().After(zd.StartTime.Add(zd.Duration).Add(zd.Decay)) {
		return 0, nil
	}

	// calculate action allocation:
	//   - zone drop action weight * claim record max allocation
	amount := zd.Actions[int32(action)].MulInt64(int64(cr.MaxAllocation)).TruncateInt64()

	// airdrop has not yet started to decay
	if ctx.BlockTime().Before(zd.StartTime.Add(zd.Duration)) {
		return uint64(amount), nil
	}

	// airdrop has started to decay, calculate claimable portion
	elapsedDecayTime := ctx.BlockTime().Sub(zd.StartTime.Add(zd.Duration))
	decayPercent := sdk.NewDec(elapsedDecayTime.Nanoseconds()).QuoInt64(zd.Decay.Nanoseconds())
	claimablePercent := sdk.OneDec().Sub(decayPercent)
	amount = claimablePercent.MulInt64(amount).TruncateInt64()

	return uint64(amount), nil
}

// GetClaimableAmountForUser returns the amount claimable for the given user
// against the given zone.
func (k Keeper) GetClaimableAmountForUser(ctx sdk.Context, chainID string, address string) (uint64, error) {
	cr, err := k.GetClaimRecord(ctx, chainID, address)
	if err != nil {
		return 0, err
	}
	if cr.Address == "" {
		return 0, nil
	}

	// get zone airdrop details
	zd, ok := k.GetZoneDrop(ctx, chainID)
	if !ok {
		return 0, types.ErrZoneDropNotFound
	}

	total := uint64(0)
	// we will only need the index as we will be calling GetClaimableAmountForAction
	for action := range zd.Actions {
		claimableForAction, err := k.GetClaimableAmountForAction(ctx, chainID, address, types.Action(action))
		if err != nil {
			return 0, err
		}
		total += claimableForAction
	}

	return total, nil
}

// Claim executes an airdrop claim for the given address on the given action
// against the given zone (chainID). It returns the claim amount or an error
// on failure.
//
// TODO: we also want to verify that the action was executed on the remote
// chain before we execute the claim...
func (k Keeper) Claim(ctx sdk.Context, chainID string, action types.Action, address string) (uint64, error) {
	// get zone airdrop details
	zd, ok := k.GetZoneDrop(ctx, chainID)
	if !ok {
		return 0, types.ErrZoneDropNotFound
	}

	// zone airdrop not active
	if !k.IsActiveZoneDrop(ctx, zd) {
		return 0, nil
	}

	// obtain claimable amount
	claimAmount, err := k.GetClaimableAmountForAction(ctx, chainID, address, action)
	if err != nil {
		return 0, err
	}

	// send coins to address
	zoneDropAccount := types.ModuleName + "." + chainID

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return 0, err
	}

	coins := sdk.NewCoins(
		sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(claimAmount)),
	)

	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zoneDropAccount, addr, coins); err != nil {
		return 0, err
	}

	// update claim record
	cr, err := k.GetClaimRecord(ctx, chainID, address)
	if err != nil {
		return 0, nil
	}

	cr.ActionCompleted[int32(action)] = &types.CompletedAction{
		CompleteTime: ctx.BlockTime(),
		ClaimAmount:  claimAmount,
	}

	if err = k.SetClaimRecord(ctx, cr); err != nil {
		return 0, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, address),
			sdk.NewAttribute("zone", chainID),
			sdk.NewAttribute(sdk.AttributeKeyAction, action.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
		),
	})

	return claimAmount, nil
}
