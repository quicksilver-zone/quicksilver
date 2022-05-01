package keeper

import (
	"fmt"
	"strings"

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

func (k Keeper) GetICAForDelegateAccount(ctx sdk.Context, address string) (*types.RegisteredZone, *types.ICAAccount) {
	var ica *types.ICAAccount
	var zone *types.RegisteredZone
	k.IterateRegisteredZones(ctx, func(_ int64, zoneInfo types.RegisteredZone) (stop bool) {
		for _, delegateAccount := range zoneInfo.DelegationAddresses {
			if delegateAccount.Address == address {
				ica = delegateAccount
				zone = &zoneInfo
				return true
			}
		}
		return false
	})
	return zone, ica
}
func (k Keeper) DetermineValidatorsForDelegation(ctx sdk.Context, zone types.RegisteredZone, amount sdk.Coin) (map[string]sdk.Coin, error) {
	out := make(map[string]sdk.Coin)

	coinAmount := amount.Amount
	aggregateIntents := zone.GetAggregateIntent()

	if len(aggregateIntents) == 0 {
		aggregateIntents = defaultAggregateIntents(ctx, zone)
	}

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

func defaultAggregateIntents(ctx sdk.Context, zone types.RegisteredZone) map[string]*types.ValidatorIntent {
	out := make(map[string]*types.ValidatorIntent)
	for _, val := range zone.GetValidators() {
		if val.CommissionRate.LTE(sdk.NewDecWithPrec(5, 1)) { // 50%; make this a param.
			out[val.GetValoperAddress()] = &types.ValidatorIntent{ValoperAddress: val.GetValoperAddress(), Weight: sdk.OneDec()}
		}
	}

	valCount := sdk.NewInt(int64(len(out)))

	// normalise the array (divide everything by length of intent list)
	for key, val := range out {
		val.Weight = val.Weight.Quo(sdk.NewDecFromInt(valCount))
		out[key] = val
	}

	return out
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
		// TODO: figure out how this impacts delegations in progress / race conditions (in most cases, the duplicate delegation will just fail)
		if !queryRes.Balances.Empty() {
			icaAccount.Balance = queryRes.Balances
			claims := k.AllWithdrawalRecords(ctx, icaAccount.Address)
			if len(claims) > 0 {
				// should we reconcile here?
				k.Logger(ctx).Info("Outstanding Withdrawal Claims", "count", len(claims))
				for _, claim := range claims {
					if claim.Status == WITHDRAW_STATUS_TOKENIZE {
						// if the claim has tokenize status AND then remove any coins in the balance that match that validator.
						// so we don't try to re-delegate any recently redeemed tokens that haven't been sent yet.
						for _, coin := range queryRes.Balances {
							if strings.HasPrefix(coin.Denom, claim.Validator) {
								k.Logger(ctx).Info("Ignoring denom this iteration", "denom", coin.GetDenom())
								queryRes.Balances = queryRes.Balances.Sub(sdk.NewCoins(coin))
							}
						}
					}
				}
			}
			if !queryRes.Balances.Empty() && !queryRes.Balances.IsZero() {
				k.Logger(ctx).Info("Delegate account balance is non-zero; delegating!", "to_delegate", queryRes.Balances)
				err := k.Delegate(ctx, zone, icaAccount)
				if err != nil {
					return err
				}
			}
		}
	}
	k.SetRegisteredZone(ctx, zone)
	return nil
}
