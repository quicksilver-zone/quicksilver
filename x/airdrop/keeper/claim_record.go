package keeper

import (
	"fmt"

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
	if !action.InBounds() {
		return 0, fmt.Errorf("%w, got %d", types.ErrActionOutOfBounds, action)
	}

	cr, err := k.GetClaimRecord(ctx, chainID, address)
	if err != nil {
		return 0, err
	}

	// action already completed, nothing to claim
	if _, exists := cr.ActionsCompleted[int32(action)]; exists {
		return 0, fmt.Errorf("%w: %s", types.ErrActionCompleted, types.Action_name[int32(action)])
	}

	// get zone airdrop details
	zd, ok := k.GetZoneDrop(ctx, cr.ChainId)
	if !ok {
		return 0, types.ErrZoneDropNotFound
	}

	// zone drop is expired, nothing to claim
	if k.IsExpiredZoneDrop(ctx, zd) {
		return 0, types.ErrZoneDropExpired
	}

	// calculate action allocation:
	//   - zone drop action weight * claim record max allocation
	// note: use int32(action)-1 as protobuf3 spec valid enum start at 1
	amount := zd.Actions[int32(action)-1].MulInt64(int64(cr.MaxAllocation)).TruncateInt64()

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

	// get zone airdrop details
	zd, ok := k.GetZoneDrop(ctx, cr.ChainId)
	if !ok {
		return 0, types.ErrZoneDropNotFound
	}

	total := uint64(0)
	// we will only need the index as we will be calling GetClaimableAmountForAction
	for i := range zd.Actions {
		// protobuf3 spec: valid enum start at 1
		action := i + 1

		claimableForAction, err := k.GetClaimableAmountForAction(ctx, cr.ChainId, cr.Address, types.Action(action))
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
func (k Keeper) Claim(
	ctx sdk.Context,
	chainID string,
	action types.Action,
	address string,
	proofs []*types.Proof,
) (uint64, error) {
	// check action in bounds
	if !action.InBounds() {
		return 0, fmt.Errorf("%w, got %d", types.ErrActionOutOfBounds, action)
	}

	// get zone airdrop details
	zd, ok := k.GetZoneDrop(ctx, chainID)
	if !ok {
		return 0, types.ErrZoneDropNotFound
	}

	// zone airdrop not active
	if !k.IsActiveZoneDrop(ctx, zd) {
		return 0, nil
	}

	// obtain claim record
	cr, err := k.GetClaimRecord(ctx, chainID, address)
	if err != nil {
		return 0, nil
	}

	return k.HandleClaim(ctx, cr, action, proofs)
}
