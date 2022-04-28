package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetRegsiteredZoneInfo returns zone info by chain_id
func (k Keeper) GetRegisteredZoneInfo(ctx sdk.Context, chain_id string) (types.RegisteredZone, bool) {
	zone := types.RegisteredZone{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	bz := store.Get([]byte(chain_id))
	if len(bz) == 0 {
		return zone, false
	}

	k.cdc.MustUnmarshal(bz, &zone)
	return zone, true
}

// SetRegisteredZone set zone info
func (k Keeper) SetRegisteredZone(ctx sdk.Context, zone types.RegisteredZone) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	bz := k.cdc.MustMarshal(&zone)
	store.Set([]byte(zone.ChainId), bz)
}

// DeleteRegisteredZone delete zone info
func (k Keeper) DeleteRegisteredZone(ctx sdk.Context, chain_id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	ctx.Logger().Error(fmt.Sprintf("Removing chain: %s", chain_id))
	store.Delete([]byte(chain_id))
}

// IterateRegisteredZones iterate through zones
func (k Keeper) IterateRegisteredZones(ctx sdk.Context, fn func(index int64, zoneInfo types.RegisteredZone) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		zone := types.RegisteredZone{}
		k.cdc.MustUnmarshal(iterator.Value(), &zone)

		stop := fn(i, zone)

		if stop {
			break
		}
		i++
	}
}

// AllRegisteredZonesInfos returns every zoneInfo in the store
func (k Keeper) AllRegisteredZones(ctx sdk.Context) []types.RegisteredZone {
	zones := []types.RegisteredZone{}
	k.IterateRegisteredZones(ctx, func(_ int64, zoneInfo types.RegisteredZone) (stop bool) {
		zones = append(zones, zoneInfo)
		return false
	})
	return zones
}

func (k Keeper) GetZoneForDelegateAccount(ctx sdk.Context, address string) *types.RegisteredZone {
	var zone *types.RegisteredZone
	k.IterateRegisteredZones(ctx, func(_ int64, zoneInfo types.RegisteredZone) (stop bool) {
		for _, ica := range zoneInfo.DelegationAddresses {
			if ica.Address == address {
				zone = &zoneInfo
				return true
			}
		}
		return false
	})
	return zone
}

func (k Keeper) GetICAForDelegateAccount(ctx sdk.Context, address string) *types.ICAAccount {
	var ica *types.ICAAccount
	k.IterateRegisteredZones(ctx, func(_ int64, zoneInfo types.RegisteredZone) (stop bool) {
		for _, delegateAccount := range zoneInfo.DelegationAddresses {
			if delegateAccount.Address == address {
				ica = delegateAccount
				return true
			}
		}
		return false
	})
	return ica
}
func (k Keeper) DetermineValidatorsForDelegation(ctx sdk.Context, zone types.RegisteredZone, amount sdk.Coin) (map[string]sdk.Coin, error) {
	out := make(map[string]sdk.Coin)

	coinAmount := amount.Amount
	aggregateIntents := zone.GetAggregateIntent()

	for valoper, intent := range aggregateIntents {
		if !coinAmount.IsZero() {
			// while there is some balance left to distribute
			// calculate the int value of weight * amount to distribute.
			thisAmount := intent.Weight.MulInt(amount.Amount).TruncateInt()
			// set distrubtion amount
			out[valoper] = sdk.Coin{Denom: amount.Denom, Amount: thisAmount}
			// reduce outstanding pool
			coinAmount = coinAmount.Sub(thisAmount)
		}
	}
	for valoper := range aggregateIntents {
		// handle leftover amount in pool (add blindly to first validator)
		out[valoper] = out[valoper].AddAmount(coinAmount)
		break
	}

	k.Logger(ctx).Info("Validator weightings without diffs", "weights", out)

	// calculate diff between current state and intended state.
	diffs := zone.DetermineStateIntentDiff(aggregateIntents)

	// apply diff to distrubtion of delegation.
	out, remaining := zone.ApplyDiffsToDistribution(out, diffs)
	if !remaining.IsZero() {
		for valoper, intent := range aggregateIntents {
			thisAmount := intent.Weight.MulInt(remaining).TruncateInt()
			thisOutAmount, ok := out[valoper]
			if !ok {
				thisOutAmount = sdk.NewCoin(amount.Denom, sdk.ZeroInt())
			}

			out[valoper] = thisOutAmount.AddAmount(thisAmount)
			remaining = remaining.Sub(thisAmount)
		}
		for valoper := range aggregateIntents {
			// handle leftover amount.
			out[valoper] = out[valoper].AddAmount(remaining)
			break
		}
	}

	k.Logger(ctx).Info("Determined validators from aggregated intents +/- rebalance diffs", "amount", amount.Amount, "out", out)
	return out, nil
}

func (k Keeper) SetAccountBalance(ctx sdk.Context, zone types.RegisteredZone, address string, queryResult []byte) error {
	queryRes := banktypes.QueryAllBalancesResponse{}
	err := k.cdc.UnmarshalJSON(queryResult, &queryRes)
	if err != nil {
		k.Logger(ctx).Error("Unable to unmarshal validators info for zone", "zone", zone.ChainId, "err", err)
		return err
	}

	switch address {
	case zone.DepositAddress.Address:
		zone.DepositAddress.Balance = queryRes.Balances
	case zone.FeeAddress.Address:
		zone.FeeAddress.Balance = queryRes.Balances
	case zone.WithdrawalAddress.Address:
		zone.WithdrawalAddress.Balance = queryRes.Balances
	default:
		icaAccount, err := zone.GetDelegationAccountByAddress(address)
		if err != nil {
			return err
		}
		icaAccount.Balance = queryRes.Balances
	}
	k.SetRegisteredZone(ctx, zone)
	return nil
}
