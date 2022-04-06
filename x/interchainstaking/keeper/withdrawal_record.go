package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	// setting WITHDRAW_STATUS_TOKENIZE as 0 causes the value to be omitted when (un)marshalling :/
	WITHDRAW_STATUS_TOKENIZE int32 = iota + 1
	WITHDRAW_STATUS_SEND     int32 = iota + 1
)

func (k Keeper) AddWithdrawalRecord(ctx sdk.Context, delegator string, validator string, recipient string, amount sdk.Coin) {
	record := &types.WithdrawalRecord{Delegator: delegator, Validator: validator, Recipient: recipient, Amount: amount, Status: WITHDRAW_STATUS_TOKENIZE}
	k.SetWithdrawalRecord(ctx, record)
}

///----------------------------------------------------------------

// GetWithdrawalRecord returns withdrawal record info by zone and delegator
func (k Keeper) GetWithdrawalRecord(ctx sdk.Context, delegator string, validator string, recipient string) (types.WithdrawalRecord, bool) {
	record := types.WithdrawalRecord{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixWithdrawalRecord, []byte(delegator)...))
	bz := store.Get([]byte(validator + recipient))
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetWithdrawalRecord store the withdrawal record
func (k Keeper) SetWithdrawalRecord(ctx sdk.Context, record *types.WithdrawalRecord) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixWithdrawalRecord, []byte(record.Delegator)...))
	bz := k.cdc.MustMarshal(record)
	store.Set([]byte(record.Validator+record.Recipient), bz)
}

// DeleteWithdrawalRecord deletes withdrawal record
func (k Keeper) DeleteWithdrawalRecord(ctx sdk.Context, delegator string, validator string, recipient string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixWithdrawalRecord, []byte(delegator)...))
	store.Delete([]byte(validator + recipient))
}

// IterateWithdrawalRecords iterate through records for a given zone
func (k Keeper) IterateWithdrawalRecords(ctx sdk.Context, delegator string, fn func(index int64, record types.WithdrawalRecord) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixWithdrawalRecord, []byte(delegator)...))

	iterator := sdk.KVStorePrefixIterator(store, nil)
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

// AllWithdrawalRecords returns every record in the store for the specified zone
func (k Keeper) AllWithdrawalRecords(ctx sdk.Context, delegator string) []types.WithdrawalRecord {
	records := []types.WithdrawalRecord{}
	k.IterateWithdrawalRecords(ctx, delegator, func(_ int64, record types.WithdrawalRecord) (stop bool) {
		records = append(records, record)
		return false
	})
	return records
}
