package keeper

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
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

func (k Keeper) GetZoneFromContext(ctx sdk.Context) (*types.RegisteredZone, error) {
	chainId, err := k.GetChainIdFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch zone from context: %w", err)
	}
	zone, found := k.GetRegisteredZoneInfo(ctx, chainId)
	if !found {
		err := fmt.Errorf("unable to fetch zone from context: not found for chainId %s", chainId)
		k.Logger(ctx).Error(err.Error())
		return nil, err
	}
	return &zone, nil
}

func (k Keeper) GetZoneForDelegateAccount(ctx sdk.Context, address string) *types.RegisteredZone {
	var zone *types.RegisteredZone
	k.IterateRegisteredZones(ctx, func(_ int64, zoneInfo types.RegisteredZone) (stop bool) {
		for _, ica := range zoneInfo.GetDelegationAccounts() {
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
		for _, delegateAccount := range zoneInfo.GetDelegationAccounts() {
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
func (k Keeper) DetermineValidatorsForDelegation(ctx sdk.Context, zone types.RegisteredZone, amount sdk.Coin) ([]string, map[string]sdk.Coin, error) {
	out := make(map[string]sdk.Coin)

	coinAmount := amount.Amount
	aggregateIntents := zone.GetAggregateIntent()

	if len(aggregateIntents) == 0 {
		aggregateIntents = defaultAggregateIntents(ctx, zone)
	}

	keys := make([]string, 0)
	for valoper, intent := range aggregateIntents {
		keys = append(keys, valoper)
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

	sort.Strings(keys)
	v0 := keys[0]
	out[v0] = out[v0].AddAmount(coinAmount)

	k.Logger(ctx).Info("Validator weightings without diffs", "weights", out)

	// calculate diff between current state and intended state.
	//diffs := zone.DetermineStateIntentDiff(aggregateIntents)

	// apply diff to distrubtion of delegation.
	// out, remaining := zone.ApplyDiffsToDistribution(out, diffs)
	// if !remaining.IsZero() {
	// 	for _, valoper := range keys {
	// 		intent := aggregateIntents[valoper]
	// 		thisAmount := intent.Weight.MulInt(remaining).TruncateInt()
	// 		thisOutAmount, ok := out[valoper]
	// 		if !ok {
	// 			thisOutAmount = sdk.NewCoin(amount.Denom, sdk.ZeroInt())
	// 		}

	// 		out[valoper] = thisOutAmount.AddAmount(thisAmount)
	// 		remaining = remaining.Sub(thisAmount)
	// 	}

	// 	v0 := keys[0]
	// 	out[v0] = out[v0].AddAmount(remaining)
	// }

	//k.Logger(ctx).Info("Determined validators from aggregated intents +/- rebalance diffs", "amount", amount.Amount, "out", out)

	return keys, out, nil
}

func defaultAggregateIntents(ctx sdk.Context, zone types.RegisteredZone) map[string]*types.ValidatorIntent {
	out := make(map[string]*types.ValidatorIntent)
	for _, val := range zone.GetValidatorsSorted() {
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

var setAccountCb Callback = func(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	balancesStore := []byte(query.Request[1:])
	accAddr, err := banktypes.AddressFromBalancesStore(balancesStore)
	if err != nil {
		return err
	}

	coin := sdk.Coin{}
	err = k.cdc.Unmarshal(args, &coin)
	if err != nil {
		k.Logger(ctx).Error("Unable to unmarshal balance info for zone", "zone", zone.ChainId, "err", err)
		return err
	}

	if coin.IsNil() {
		denom := ""

		for i := 0; i < len(query.Request)-len(accAddr); i++ {
			if bytes.Equal(query.Request[i:i+len(accAddr)], accAddr) {
				denom = string(query.Request[i+len(accAddr):])
				k.Logger(ctx).Error("denom found", "denom", denom)
				break
			}

		}
		// if balance is nil, the response sent back is nil, so we don't receive the denom. Override that now.
		coin = sdk.NewCoin(denom, sdk.ZeroInt())
	}

	address, err := bech32.ConvertAndEncode(zone.AccountPrefix, accAddr)
	if err != nil {
		return err
	}

	return SetAccountBalanceForDenom(k, ctx, zone, address, coin)
}

func SetAccountBalanceForDenom(k Keeper, ctx sdk.Context, zone types.RegisteredZone, address string, coin sdk.Coin) error {

	switch true {
	case zone.DepositAddress != nil && address == zone.DepositAddress.Address:
		existing := zone.DepositAddress.Balance.AmountOf(coin.Denom)
		zone.DepositAddress.Balance = zone.DepositAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))).Add(coin) // reset this denom
		zone.DepositAddress.BalanceWaitgroup = zone.DepositAddress.BalanceWaitgroup - 1
		k.Logger(ctx).Info("Matched deposit address", "address", address, "wg", zone.DepositAddress.BalanceWaitgroup, "balance", zone.DepositAddress.Balance)
		if zone.DepositAddress.BalanceWaitgroup == 0 {
			k.depositInterval(ctx)(0, zone)
		}
	case zone.WithdrawalAddress != nil && address == zone.WithdrawalAddress.Address:
		existing := zone.WithdrawalAddress.Balance.AmountOf(coin.Denom)
		zone.WithdrawalAddress.Balance = zone.WithdrawalAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))).Add(coin) // reset this denom
		zone.WithdrawalAddress.BalanceWaitgroup = zone.WithdrawalAddress.BalanceWaitgroup - 1
		k.Logger(ctx).Info("Matched withdrawal address", "address", address, "wg", zone.WithdrawalAddress.BalanceWaitgroup, "balance", zone.WithdrawalAddress.Balance)
	default:
		icaAccount, err := zone.GetDelegationAccountByAddress(address)
		if err != nil {
			return err
		}
		existing := icaAccount.Balance.AmountOf(coin.Denom)
		k.Logger(ctx).Info("Matched delegate address", "address", address, "wg", icaAccount.BalanceWaitgroup, "balance", icaAccount.Balance)

		icaAccount.Balance = icaAccount.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))) // zero this denom

		// TODO: figure out how this impacts delegations in progress / race conditions (in most cases, the duplicate delegation will just fail)
		if !icaAccount.Balance.Empty() {
			claims := k.AllWithdrawalRecords(ctx, icaAccount.Address)
			if len(claims) > 0 {
				// should we reconcile here?
				k.Logger(ctx).Info("Outstanding Withdrawal Claims", "count", len(claims))
				for _, claim := range claims {
					if claim.Status == WITHDRAW_STATUS_TOKENIZE {
						// if the claim has tokenize status AND then remove any coins in the balance that match that validator.
						// so we don't try to re-delegate any recently redeemed tokens that haven't been sent yet.
						if strings.HasPrefix(coin.Denom, claim.Validator) {
							k.Logger(ctx).Info("Ignoring denom this iteration", "denom", coin.GetDenom())
							coin = coin.Sub(claim.Amount)
						}
					}
				}
			}
		}

		icaAccount.Balance = icaAccount.Balance.Add(coin)
		k.Logger(ctx).Info("Matched delegate address", "address", address, "wg", icaAccount.BalanceWaitgroup, "balance", icaAccount.Balance)

		if zone.WithdrawalAddress.BalanceWaitgroup == 0 {
			if !icaAccount.Balance.Empty() {
				k.Logger(ctx).Info("Delegate account balance is non-zero; delegating!", "to_delegate", icaAccount.Balance)
				err := k.Delegate(ctx, zone, icaAccount)
				if err != nil {
					return err
				}
			}
		}

		icaAccount.BalanceWaitgroup = icaAccount.BalanceWaitgroup - 1

	}
	k.SetRegisteredZone(ctx, zone)
	return nil
}

func (k Keeper) SetAccountBalance(ctx sdk.Context, zone types.RegisteredZone, address string, queryResult []byte) error {
	queryRes := banktypes.QueryAllBalancesResponse{}
	err := k.cdc.Unmarshal(queryResult, &queryRes)
	if err != nil {
		k.Logger(ctx).Error("Unable to unmarshal balance", "zone", zone.ChainId, "err", err)
		return err
	}
	_, addr, _ := bech32.DecodeAndConvert(address)
	data := banktypes.CreateAccountBalancesPrefix(addr)

	var icaAccount *types.ICAAccount

	switch true {
	case zone.DepositAddress != nil && address == zone.DepositAddress.Address:
		icaAccount = zone.DepositAddress
	case zone.WithdrawalAddress != nil && address == zone.WithdrawalAddress.Address:
		icaAccount = zone.WithdrawalAddress
	default:
		icaAccount, err = zone.GetDelegationAccountByAddress(address)
		if err != nil {
			return err
		}
	}

	if icaAccount == nil {
		return fmt.Errorf("unable to determine account for address %s", address)
	}

	for _, coin := range zone.DepositAddress.Balance {
		if queryRes.Balances.AmountOf(coin.Denom).Equal(sdk.ZeroInt()) {
			// coin we used to have is now zero - also validate this.
			key := "store/bank/key"
			k.Logger(ctx).Error("Querying for value", "key", key, "denom", coin.Denom)
			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				key,
				append(data, []byte(coin.Denom)...),
				sdk.NewInt(-1),
				types.ModuleName,
				setAccountCb,
			)
			icaAccount.BalanceWaitgroup += 1

		}

	}

	for _, coin := range queryRes.Balances {
		key := "store/bank/key"
		k.Logger(ctx).Error("Querying for value", "key", key, "denom", coin.Denom)
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			key,
			append(data, []byte(coin.Denom)...),
			sdk.NewInt(-1),
			types.ModuleName,
			setAccountCb,
		)
		icaAccount.BalanceWaitgroup += 1
	}

	k.SetRegisteredZone(ctx, zone)
	return nil
}
