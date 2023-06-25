package keeper

import (
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) GetNextWithdrawalRecordSequence(ctx sdk.Context) (sequence uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := store.Get(types.KeyPrefixRequeuedWithdrawalRecordSeq)
	if bz == nil {
		bz := make([]byte, 8)
		binary.BigEndian.PutUint64(bz, uint64(2))
		store.Set(types.KeyPrefixRequeuedWithdrawalRecordSeq, bz)
		return 1
	}
	sequence = binary.BigEndian.Uint64(bz)
	binary.BigEndian.PutUint64(bz, sequence+1)
	store.Set(types.KeyPrefixRequeuedWithdrawalRecordSeq, bz)
	return sequence
}

func (k *Keeper) AddWithdrawalRecord(ctx sdk.Context, chainID, delegator string, distribution []*types.Distribution, recipient string, amount sdk.Coins, burnAmount sdk.Coin, hash string, status int32, completionTime time.Time) {
	record := types.WithdrawalRecord{ChainId: chainID, Delegator: delegator, Distribution: distribution, Recipient: recipient, Amount: amount, Status: status, BurnAmount: burnAmount, Txhash: hash, CompletionTime: completionTime}
	k.Logger(ctx).Error("addWithdrawalRecord", "record", record)
	k.SetWithdrawalRecord(ctx, record)
}

///----------------------------------------------------------------

// GetWithdrawalRecord returns withdrawal record info by zone and delegator.
func (k *Keeper) GetWithdrawalRecord(ctx sdk.Context, chainID, txhash string, status int32) (types.WithdrawalRecord, bool) {
	record := types.WithdrawalRecord{}

	key, err := hex.DecodeString(txhash)
	if err != nil {
		return record, false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetWithdrawalKey(chainID, status))
	bz := store.Get(key)
	if bz == nil {
		return record, false
	}
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetWithdrawalRecord store the withdrawal record.
func (k *Keeper) SetWithdrawalRecord(ctx sdk.Context, record types.WithdrawalRecord) {
	key, err := hex.DecodeString(record.Txhash)
	if err != nil {
		panic(err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetWithdrawalKey(record.ChainId, record.Status))
	bz := k.cdc.MustMarshal(&record)
	store.Set(key, bz)
}

func (k *Keeper) UpdateWithdrawalRecordStatus(ctx sdk.Context, withdrawal *types.WithdrawalRecord, newStatus int32) {
	k.DeleteWithdrawalRecord(ctx, withdrawal.ChainId, withdrawal.Txhash, withdrawal.Status)
	withdrawal.Status = newStatus
	k.SetWithdrawalRecord(ctx, *withdrawal)
}

// DeleteWithdrawalRecord deletes withdrawal record.
func (k *Keeper) DeleteWithdrawalRecord(ctx sdk.Context, chainID, txhash string, status int32) {
	key, err := hex.DecodeString(txhash)
	if err != nil {
		panic(err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetWithdrawalKey(chainID, status))
	store.Delete(key)
}

// IteratePrefixedWithdrawalRecords iterate through all records with given prefix.
func (k *Keeper) IteratePrefixedWithdrawalRecords(ctx sdk.Context, prefixBytes []byte, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
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

// IterateWithdrawalRecords iterate through all records.
func (k *Keeper) IterateWithdrawalRecords(ctx sdk.Context, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	k.IteratePrefixedWithdrawalRecords(ctx, nil, fn)
}

// IterateZoneWithdrawalRecords iterate through records for a given zone.
func (k *Keeper) IterateZoneWithdrawalRecords(ctx sdk.Context, chainID string, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	k.IteratePrefixedWithdrawalRecords(ctx, []byte(chainID), fn)
}

// IterateZoneStatusWithdrawalRecords iterate through records for a given zone / delegator tuple.
func (k *Keeper) IterateZoneStatusWithdrawalRecords(ctx sdk.Context, chainID string, status int32, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	statusBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(statusBytes, uint64(status))
	key := append([]byte(chainID), statusBytes...)
	k.IteratePrefixedWithdrawalRecords(ctx, key, fn)
}

// AllWithdrawalRecords returns every record in the store for the specified zone.
func (k *Keeper) AllWithdrawalRecords(ctx sdk.Context) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateWithdrawalRecords(ctx, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// AllUserWithdrawalRecords returns every record in the store for the specified user.
func (k *Keeper) AllUserWithdrawalRecords(ctx sdk.Context, address string) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateWithdrawalRecords(ctx, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		if record.Delegator == address {
			records = append(records, record)
		}
		return false
	})
	return records
}

// AllZoneWithdrawalRecords returns every record in the store for the specified zone.
func (k *Keeper) AllZoneWithdrawalRecords(ctx sdk.Context, chainID string) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateZoneWithdrawalRecords(ctx, chainID, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// GetUnbondingRecord returns unbonding record info by zone, validator and epoch.
func (k *Keeper) GetUnbondingRecord(ctx sdk.Context, chainID, validator string, epochNumber int64) (types.UnbondingRecord, bool) {
	record := types.UnbondingRecord{}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := store.Get(types.GetUnbondingKey(chainID, validator, epochNumber))
	if bz == nil {
		return record, false
	}
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetUnbondingRecord store the unbonding record.
func (k *Keeper) SetUnbondingRecord(ctx sdk.Context, record types.UnbondingRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := k.cdc.MustMarshal(&record)
	store.Set(types.GetUnbondingKey(record.ChainId, record.Validator, record.EpochNumber), bz)
}

// DeleteUnbondingRecord deletes unbonding record.
func (k *Keeper) DeleteUnbondingRecord(ctx sdk.Context, chainID, validator string, epochNumber int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	store.Delete(types.GetUnbondingKey(chainID, validator, epochNumber))
}

// IteratePrefixedUnbondingRecords iterate through all records with given prefix.
func (k *Keeper) IteratePrefixedUnbondingRecords(ctx sdk.Context, prefixBytes []byte, fn func(index int64, record types.UnbondingRecord) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUnbondingRecord)

	iterator := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		record := types.UnbondingRecord{}
		k.cdc.MustUnmarshal(iterator.Value(), &record)

		stop := fn(i, record)

		if stop {
			break
		}
		i++
	}
}

// IterateUnbondingRecords iterate through all records.
func (k *Keeper) IterateUnbondingRecords(ctx sdk.Context, fn func(index int64, record types.UnbondingRecord) (stop bool)) {
	k.IteratePrefixedUnbondingRecords(ctx, nil, fn)
}

// AllUnbondingRecords returns every record in the store.
func (k *Keeper) AllUnbondingRecords(ctx sdk.Context) []types.UnbondingRecord {
	var records []types.UnbondingRecord
	k.IterateUnbondingRecords(ctx, func(_ int64, record types.UnbondingRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}

// AllZoneUnbondingRecords returns every record in the store for the specified zone.
func (k *Keeper) AllZoneUnbondingRecords(ctx sdk.Context, chainID string) []types.UnbondingRecord {
	var records []types.UnbondingRecord
	k.IteratePrefixedUnbondingRecords(ctx, []byte(chainID), func(_ int64, record types.UnbondingRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}
