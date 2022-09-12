package keeper

import (
	"encoding/binary"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	// setting WithdrawStatusTokenize as 0 causes the value to be omitted when (un)marshalling :/
	WithdrawStatusTokenize int32 = iota + 1
	WithdrawStatusQueued   int32 = iota + 1
	WithdrawStatusUnbond   int32 = iota + 1
	WithdrawStatusSend     int32 = iota + 1
)

func deprotoizeIntMap(m map[string]sdk.Int) map[string]int64 {
	n := make(map[string]int64, 0)
	for _, j := range utils.Keys(m) {
		n[j] = m[j].Int64()
	}
	return n
}

func (k Keeper) AddWithdrawalRecord(ctx sdk.Context, zone types.Zone, delegator string, distribution map[string]sdk.Int, recipient string, amount sdk.Coins, burnAmount sdk.Coin, hash string, status int32, completionTime time.Time) {
	record := &types.WithdrawalRecord{ChainId: zone.ChainId, Delegator: delegator, Distribution: deprotoizeIntMap(distribution), Recipient: recipient, Amount: amount, Status: status, BurnAmount: burnAmount, Txhash: hash, CompletionTime: completionTime}
	k.SetWithdrawalRecord(ctx, record)
}

func GetWithdrawalKey(chainID string, status int32) []byte {
	statusBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(statusBytes, uint64(status))
	return append(types.KeyPrefixWithdrawalRecord, append([]byte(chainID), statusBytes...)...)
}

///----------------------------------------------------------------

// GetWithdrawalRecord returns withdrawal record info by zone and delegator
func (k Keeper) GetWithdrawalRecord(ctx sdk.Context, zone *types.Zone, txhash string, status int32) (types.WithdrawalRecord, bool) {
	record := types.WithdrawalRecord{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetWithdrawalKey(zone.ChainId, status))
	bz := store.Get([]byte(txhash))
	if bz == nil {
		return record, false
	}
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetWithdrawalRecord store the withdrawal record
func (k Keeper) SetWithdrawalRecord(ctx sdk.Context, record *types.WithdrawalRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetWithdrawalKey(record.ChainId, record.Status))
	bz := k.cdc.MustMarshal(record)
	store.Set([]byte(record.Txhash), bz)
}

// DeleteWithdrawalRecord deletes withdrawal record
func (k Keeper) DeleteWithdrawalRecord(ctx sdk.Context, zone *types.Zone, txhash string, status int32) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetWithdrawalKey(zone.ChainId, status))
	store.Delete([]byte(txhash))
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
func (k Keeper) IterateZoneStatusWithdrawalRecords(ctx sdk.Context, zone *types.Zone, status int32, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	statusBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(statusBytes, uint64(status))
	key := append([]byte(zone.ChainId), statusBytes...)
	k.IteratePrefixedWithdrawalRecords(ctx, key, fn)
}

// AllZoneDelegatorWithdrawalRecords returns every record in the store for the specified zone / delegator tuple
func (k Keeper) AllZoneStatusWithdrawalRecords(ctx sdk.Context, zone *types.Zone, status int32) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateZoneStatusWithdrawalRecords(ctx, zone, status, func(_ int64, record types.WithdrawalRecord) (stop bool) {
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
