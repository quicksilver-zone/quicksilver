package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
)

// GetRedelegationRecord returns Redelegation record info by zone, validator and epoch.
func (k *Keeper) GetRedelegationRecord(ctx sdk.Context, chainID, source, destination string, epochNumber int64) (types.RedelegationRecord, bool) {
	record := types.RedelegationRecord{}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := store.Get(types.GetRedelegationKey(chainID, source, destination, epochNumber))
	if bz == nil {
		return record, false
	}
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetRedelegationRecord store the Redelegation record.
func (k *Keeper) SetRedelegationRecord(ctx sdk.Context, record types.RedelegationRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := k.cdc.MustMarshal(&record)
	store.Set(types.GetRedelegationKey(record.ChainId, record.Source, record.Destination, record.EpochNumber), bz)
}

// DeleteRedelegationRecord deletes Redelegation record.
func (k *Keeper) DeleteRedelegationRecord(ctx sdk.Context, chainID, source, destination string, epochNumber int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	store.Delete(types.GetRedelegationKey(chainID, source, destination, epochNumber))
}

// DeleteRedelegationRecordByKey deletes Redelegation record.
func (k *Keeper) DeleteRedelegationRecordByKey(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	store.Delete(key)
}

// IteratePrefixedRedelegationRecords iterate through all records with given prefix.
func (k *Keeper) IteratePrefixedRedelegationRecords(ctx sdk.Context, prefixBytes []byte, fn func(index int64, key []byte, record types.RedelegationRecord) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRedelegationRecord)

	iterator := storetypes.KVStorePrefixIterator(store, prefixBytes)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		record := types.RedelegationRecord{}
		k.cdc.MustUnmarshal(iterator.Value(), &record)

		stop := fn(i, iterator.Key(), record)

		if stop {
			break
		}
		i++
	}
}

// IterateRedelegationRecords iterate through all records.
func (k *Keeper) IterateRedelegationRecords(ctx sdk.Context, fn func(index int64, key []byte, record types.RedelegationRecord) (stop bool)) {
	k.IteratePrefixedRedelegationRecords(ctx, nil, fn)
}

// AllRedelegationRecords returns every record in the store for the specified zone.
func (k *Keeper) AllRedelegationRecords(ctx sdk.Context) []types.RedelegationRecord {
	records := []types.RedelegationRecord{}
	k.IterateRedelegationRecords(ctx, func(_ int64, _ []byte, record types.RedelegationRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// ZoneRedelegationRecords returns every record in the store for the specified zone.
func (k *Keeper) ZoneRedelegationRecords(ctx sdk.Context, chainID string) []types.RedelegationRecord {
	records := []types.RedelegationRecord{}
	k.IteratePrefixedRedelegationRecords(ctx, []byte(chainID), func(_ int64, _ []byte, record types.RedelegationRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// IterateZoneRedelegationRecords iterates through every record in the store for the specified zone.
func (k *Keeper) IterateZoneRedelegationRecords(ctx sdk.Context, zoneID string, fn func(index int64, key []byte, record types.RedelegationRecord) (stop bool)) {
	k.IteratePrefixedRedelegationRecords(ctx, []byte(zoneID), fn)
}
