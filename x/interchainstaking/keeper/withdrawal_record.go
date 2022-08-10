package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	// setting WithdrawStatusTokenize as 0 causes the value to be omitted when (un)marshalling :/
	WithdrawStatusTokenize int32 = iota + 1
	WithdrawStatusUnbond   int32 = iota + 1
	WithdrawStatusSend     int32 = iota + 1
)

func (k Keeper) AddWithdrawalRecord(ctx sdk.Context, zone *types.Zone, delegator string, validator string, recipient string, amount sdk.Coin, burnAmount sdk.Coin, hash string, completionTime time.Time) {
	record := &types.WithdrawalRecord{ChainId: zone.ChainId, Delegator: delegator, Validator: validator, Recipient: recipient, Amount: amount, Status: WithdrawStatusTokenize, BurnAmount: burnAmount, Txhash: hash, CompletionTime: completionTime}
	k.SetWithdrawalRecord(ctx, record)
}

func GetWithdrawalKey(chainID string, delegator string, txhash string) []byte {
	return append(types.KeyPrefixWithdrawalRecord, append([]byte(chainID), append([]byte(delegator), []byte(txhash)...)...)...)
}

///----------------------------------------------------------------

// GetWithdrawalRecord returns withdrawal record info by zone and delegator
func (k Keeper) GetWithdrawalRecord(ctx sdk.Context, zone *types.Zone, txhash string, delegator string, validator string) (types.WithdrawalRecord, bool) {
	record := types.WithdrawalRecord{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetWithdrawalKey(zone.ChainId, delegator, txhash))
	k.Logger(ctx).Error("Fetch key: ", "key", GetWithdrawalKey(zone.ChainId, delegator, txhash), "val", []byte(validator))
	bz := store.Get([]byte(validator))
	if bz == nil {
		return record, false
	}
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetWithdrawalRecord store the withdrawal record
func (k Keeper) SetWithdrawalRecord(ctx sdk.Context, record *types.WithdrawalRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetWithdrawalKey(record.ChainId, record.Delegator, record.Txhash))
	k.Logger(ctx).Error("Store key: ", "key", GetWithdrawalKey(record.ChainId, record.Delegator, record.Txhash), "val", []byte(record.Validator))
	bz := k.cdc.MustMarshal(record)
	store.Set([]byte(record.Validator), bz)
}

// DeleteWithdrawalRecord deletes withdrawal record
func (k Keeper) DeleteWithdrawalRecord(ctx sdk.Context, zone *types.Zone, txhash string, delegator string, validator string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetWithdrawalKey(zone.ChainId, delegator, txhash))
	k.Logger(ctx).Error("Delete key: ", "key", GetWithdrawalKey(zone.ChainId, delegator, txhash), "val", []byte(validator))
	store.Delete([]byte(validator))
}

// IteratePrefixedWithdrawalRecords iterate through all records with given prefix
func (k Keeper) IteratePrefixedWithdrawalRecords(ctx sdk.Context, prefixBytes []byte, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWithdrawalRecord)

	iterator := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		record := types.WithdrawalRecord{}
		k.cdc.MustUnmarshal(iterator.Value(), &record)

		stop := fn(i, record)

		if stop {
			break
		}
		i++
	}
}

// IterateWithdrawalRecords iterate through all records
func (k Keeper) IterateWithdrawalRecords(ctx sdk.Context, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	k.IteratePrefixedWithdrawalRecords(ctx, nil, fn)
}

// IterateZoneWithdrawalRecords iterate through records for a given zone
func (k Keeper) IterateZoneWithdrawalRecords(ctx sdk.Context, zone *types.Zone, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	k.IteratePrefixedWithdrawalRecords(ctx, []byte(zone.ChainId), fn)
}

// IterateZoneDelegatorWithdrawalRecords iterate through records for a given zone / delegator tuple
func (k Keeper) IterateZoneDelegatorWithdrawalRecords(ctx sdk.Context, zone *types.Zone, delegator string, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	k.IteratePrefixedWithdrawalRecords(ctx, []byte(zone.ChainId+delegator), fn)
}

// IterateWithdrawalRecords iterate through records for a given zone / delegator / hash treble
func (k Keeper) IterateZoneDelegatorHashWithdrawalRecords(ctx sdk.Context, zone *types.Zone, txhash string, delegator string, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	k.IteratePrefixedWithdrawalRecords(ctx, GetWithdrawalKey(zone.ChainId, delegator, txhash), fn)
}

// AllZoneDelegatorHashWithdrawalRecords returns every record in the store for the specified zone / delegator / hash treble
func (k Keeper) AllZoneDelegatorHashWithdrawalRecords(ctx sdk.Context, zone *types.Zone, txhash string, delegator string) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateZoneDelegatorHashWithdrawalRecords(ctx, zone, txhash, delegator, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// AllZoneDelegatorWithdrawalRecords returns every record in the store for the specified zone / delegator tuple
func (k Keeper) AllZoneDelegatorWithdrawalRecords(ctx sdk.Context, zone *types.Zone, delegator string) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateZoneDelegatorWithdrawalRecords(ctx, zone, delegator, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// AllZoneWithdrawalRecords returns every record in the store for the specified zone
func (k Keeper) AllZoneWithdrawalRecords(ctx sdk.Context, zone *types.Zone) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateZoneWithdrawalRecords(ctx, zone, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// AllWithdrawalRecords returns every record in the store for the specified zone
func (k Keeper) AllWithdrawalRecords(ctx sdk.Context) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateWithdrawalRecords(ctx, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}
