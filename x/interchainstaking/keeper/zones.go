package keeper

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetZone returns zone info by chainID
func (k Keeper) GetZone(ctx sdk.Context, chainID string) (types.Zone, bool) {
	zone := types.Zone{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	bz := store.Get([]byte(chainID))
	if len(bz) == 0 {
		return zone, false
	}

	k.cdc.MustUnmarshal(bz, &zone)
	return zone, true
}

// SetZone set zone info
func (k Keeper) SetZone(ctx sdk.Context, zone *types.Zone) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	bz := k.cdc.MustMarshal(zone)
	store.Set([]byte(zone.ChainId), bz)
}

// DeleteZone delete zone info
func (k Keeper) DeleteZone(ctx sdk.Context, chainID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	store.Delete([]byte(chainID))
}

// IterateZones iterate through zones
func (k Keeper) IterateZones(ctx sdk.Context, fn func(index int64, zoneInfo types.Zone) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		zone := types.Zone{}
		k.cdc.MustUnmarshal(iterator.Value(), &zone)

		stop := fn(i, zone)

		if stop {
			break
		}
		i++
	}
}

func (k Keeper) GetDelegatedAmount(ctx sdk.Context, zone *types.Zone) sdk.Coin {
	out := sdk.Coin{}
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) (stop bool) {
		out = out.Add(delegation.Amount)
		return false
	})
	return out
}

// AllZonesInfos returns every zoneInfo in the store
func (k Keeper) AllZones(ctx sdk.Context) []types.Zone {
	zones := []types.Zone{}
	k.IterateZones(ctx, func(_ int64, zoneInfo types.Zone) (stop bool) {
		zones = append(zones, zoneInfo)
		return false
	})
	return zones
}

// GetZoneFromContext determines the zone from the current context
func (k Keeper) GetZoneFromContext(ctx sdk.Context) (*types.Zone, error) {
	chainID, err := k.GetChainIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch zone from context: %w", err)
	}
	zone, found := k.GetZone(ctx, chainID)
	if !found {
		err := fmt.Errorf("unable to fetch zone from context: not found for chainID %s", chainID)
		k.Logger(ctx).Error(err.Error())
		return nil, err
	}
	return &zone, nil
}

// GetZoneForDelegateAccount determines the zone for a given address.
func (k Keeper) GetZoneForDelegateAccount(ctx sdk.Context, address string) *types.Zone {
	var zone *types.Zone
	k.IterateZones(ctx, func(_ int64, zoneInfo types.Zone) (stop bool) {
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

func (k Keeper) GetZoneForPerformanceAccount(ctx sdk.Context, address string) *types.Zone {
	var zone *types.Zone
	k.IterateZones(ctx, func(_ int64, zoneInfo types.Zone) (stop bool) {
		if zoneInfo.PerformanceAddress.Address == address {
			zone = &zoneInfo
			return true
		}
		return false
	})
	return zone
}

// GetZoneForDelegateAccount determines the zone, and returns the ICAAccount for a given address.
func (k Keeper) GetICAForDelegateAccount(ctx sdk.Context, address string) (*types.Zone, *types.ICAAccount) {
	var ica *types.ICAAccount
	var zone *types.Zone
	k.IterateZones(ctx, func(_ int64, zoneInfo types.Zone) (stop bool) {
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

func (k *Keeper) EnsureWithdrawalAddresses(ctx sdk.Context, zone *types.Zone) error {
	if zone.WithdrawalAddress == nil {
		k.Logger(ctx).Info("Withdrawal address not set")
		return nil
	}
	withdrawalAddress := zone.WithdrawalAddress.Address

	for _, da := range zone.GetDelegationAccounts() {
		if da.WithdrawalAddress != zone.WithdrawalAddress.Address {
			msg := distrTypes.MsgSetWithdrawAddress{DelegatorAddress: da.Address, WithdrawAddress: withdrawalAddress}
			err := k.SubmitTx(ctx, []sdk.Msg{&msg}, da, "")
			if err != nil {
				return err
			}
		}
	}

	// set withdrawal address for performance address, if it exists
	if zone.PerformanceAddress != nil && zone.PerformanceAddress.WithdrawalAddress != withdrawalAddress {
		msg := distrTypes.MsgSetWithdrawAddress{DelegatorAddress: zone.PerformanceAddress.Address, WithdrawAddress: withdrawalAddress}
		err := k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.PerformanceAddress, "")
		if err != nil {
			return err
		}
	}
	return nil
}

// SetAccountBalanceForDenom sets the balance on an account for a given denominination.
func SetAccountBalanceForDenom(k Keeper, ctx sdk.Context, zone types.Zone, address string, coin sdk.Coin) error {
	// ? is this switch statement still required ?
	// prior to callback we had no way to distinguish the originator
	// with the query type in setAccountCb this is probably superfluous...
	switch {
	case zone.DepositAddress != nil && address == zone.DepositAddress.Address:
		existing := zone.DepositAddress.Balance.AmountOf(coin.Denom)
		zone.DepositAddress.SetBalance(zone.DepositAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))).Add(coin)) // reset this denom
		zone.DepositAddress.DecrementBalanceWaitgroup()
		k.Logger(ctx).Info("Matched deposit address", "address", address, "wg", zone.DepositAddress.BalanceWaitgroup, "balance", zone.DepositAddress.Balance)
		if zone.DepositAddress.BalanceWaitgroup == 0 {
			k.depositInterval(ctx)(0, zone)
		}
	case zone.WithdrawalAddress != nil && address == zone.WithdrawalAddress.Address:
		existing := zone.WithdrawalAddress.Balance.AmountOf(coin.Denom)
		zone.WithdrawalAddress.SetBalance(zone.WithdrawalAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))).Add(coin)) // reset this denom
		zone.WithdrawalAddress.DecrementBalanceWaitgroup()
		k.Logger(ctx).Info("Matched withdrawal address", "address", address, "wg", zone.WithdrawalAddress.BalanceWaitgroup, "balance", zone.WithdrawalAddress.Balance)
	case zone.PerformanceAddress != nil && address == zone.PerformanceAddress.Address:
		k.Logger(ctx).Info("Matched performance address")
	default:
		panic("do i get called?")
		icaAccount, err := zone.GetDelegationAccountByAddress(address)
		if err != nil {
			return err
		}
		existing := icaAccount.Balance.AmountOf(coin.Denom)
		k.Logger(ctx).Info("Matched delegate address", "address", address, "wg", icaAccount.BalanceWaitgroup, "balance", icaAccount.Balance)

		icaAccount.Balance = icaAccount.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))) // zero this denom

		// TODO: figure out how this impacts delegations in progress / race conditions (in most cases, the duplicate delegation will just fail)
		if !icaAccount.Balance.Empty() {
			claims := k.AllZoneDelegatorWithdrawalRecords(ctx, &zone, icaAccount.Address)
			if len(claims) > 0 {
				// should we reconcile here?
				k.Logger(ctx).Info("Outstanding Withdrawal Claims", "count", len(claims))
				for _, claim := range claims {
					if claim.Status == WithdrawStatusTokenize || claim.Status == WithdrawStatusUnbond {
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
				valPlan, err := types.DelegationPlanFromGlobalIntent(k.GetDelegatedAmount(ctx, &zone), k.GetDelegationBinsMap(ctx, &zone), coin, zone.GetAggregateIntentOrDefault())
				if err != nil {
					return err
				}
				err = k.Delegate(ctx, zone, icaAccount, valPlan)
				if err != nil {
					return err
				}
			}
		}

		icaAccount.BalanceWaitgroup--

	}
	k.SetZone(ctx, &zone)
	return nil
}

// SetAccountBalance triggers provable KV queries to prove an AllBalances query.
func (k Keeper) SetAccountBalance(ctx sdk.Context, zone types.Zone, address string, queryResult []byte) error {
	queryRes := banktypes.QueryAllBalancesResponse{}
	err := k.cdc.Unmarshal(queryResult, &queryRes)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal balance", "zone", zone.ChainId, "err", err)
		return err
	}
	_, addr, _ := bech32.DecodeAndConvert(address)
	data := banktypes.CreateAccountBalancesPrefix(addr)

	var icaAccount *types.ICAAccount

	switch {
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
			k.Logger(ctx).Info("Querying for value", "key", key, "denom", coin.Denom) // debug?
			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				key,
				append(data, []byte(coin.Denom)...),
				sdk.NewInt(-1),
				types.ModuleName,
				"accountbalance",
				0,
			)
			icaAccount.BalanceWaitgroup++

		}
	}

	for _, coin := range queryRes.Balances {
		key := "store/bank/key"
		k.Logger(ctx).Info("Querying for value", "key", key, "denom", coin.Denom) // debug?
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			key,
			append(data, []byte(coin.Denom)...),
			sdk.NewInt(-1),
			types.ModuleName,
			"accountbalance",
			0,
		)
		icaAccount.BalanceWaitgroup++
	}

	k.SetZone(ctx, &zone)
	return nil
}

type (
	RedemptionTarget  types.DelegationPlan
	RedemptionTargets []RedemptionTarget
)

func (r RedemptionTargets) Sorted() RedemptionTargets {
	sort.SliceStable(r, func(i, j int) bool {
		return fmt.Sprintf("%s%s", r[i].DelegatorAddress, r[i].ValidatorAddress) < fmt.Sprintf("%s%s", r[j].DelegatorAddress, r[j].ValidatorAddress)
	})
	return r
}

func (r RedemptionTargets) Get(delAddr string, valAddr string) *RedemptionTarget {
	for _, rt := range r.Sorted() {
		if rt.DelegatorAddress == delAddr && rt.ValidatorAddress == valAddr {
			return &rt
		}
	}
	return nil
}

func (r RedemptionTargets) Add(delAddr string, valAddr string, amount sdk.Coins) RedemptionTargets {
	for _, rt := range r.Sorted() {
		if rt.DelegatorAddress == delAddr && rt.ValidatorAddress == valAddr {
			rt.Value = rt.Value.Add(amount...)
			return r
		}
	}
	return append(r, RedemptionTarget{ValidatorAddress: valAddr, DelegatorAddress: delAddr, Value: amount})
}

func ApplyDeltasToIntent(requests types.Allocations, deltas types.Diffs, currentState types.Allocations) types.Allocations {
OUT:
	for fromIdx := 0; fromIdx < len(deltas) && deltas[fromIdx].Amount.LT(sdk.ZeroInt()); {
		for idx := len(deltas) - 1; idx > fromIdx; idx-- {
			if idx == fromIdx {
				continue
			}
			if intent := requests.Get(deltas[idx].Valoper); intent != nil {
				var remainder sdk.Coins
				toSub := deltas[fromIdx].Amount.Abs()
				requests, remainder = requests.Sub(sdk.Coins{sdk.NewCoin(types.GenericToken, toSub)}, intent.Address)
				requests = requests.Allocate(deltas[fromIdx].Valoper, sdk.Coins{sdk.NewCoin(types.GenericToken, toSub)}.Sub(remainder))
				deltas[fromIdx].Amount = remainder.AmountOf(types.GenericToken).Neg()
				if deltas[fromIdx].Amount.Equal(sdk.ZeroInt()) {
					fromIdx++
					continue OUT
				}
			}
		}
		if !deltas[fromIdx].Amount.Equal(sdk.ZeroInt()) {
			break
		}

	}

	return SatisfyRequestsForBins(requests, currentState, deltas).Sorted()
}

func SatisfyRequestsForBins(requests types.Allocations, bins types.Allocations, deltas types.Diffs) types.Allocations {
	for dIdx, delta := range deltas {
		maxWithdrawableForDenom := bins.SumForDenom(delta.Valoper)
		r := requests.Get(delta.Valoper)
		if r != nil {
			if r.Amount.AmountOf(types.GenericToken).GT(maxWithdrawableForDenom) {
				toReallocate := sdk.Coins{sdk.Coin{Denom: types.GenericToken, Amount: r.Amount.AmountOf(types.GenericToken).Sub(maxWithdrawableForDenom)}}
				requests, _ = requests.Sub(toReallocate, r.Address)
				if dIdx >= len(deltas)-1 {
					requests = requests.Allocate(deltas[0].Valoper, toReallocate)
					return SatisfyRequestsForBins(requests, bins, deltas)
				}
				requests = requests.Allocate(deltas[dIdx+1].Valoper, toReallocate)
			}
		}
	}
	return requests
}

func (k *Keeper) GetRedemptionTargets(ctx sdk.Context, zone types.Zone, requests types.Allocations) RedemptionTargets {
	out := RedemptionTargets{}

	bins := k.GetDelegationBinsMap(ctx, &zone)

	deltas := types.DetermineIntentDelta(bins, k.GetDelegatedAmount(ctx, &zone).Amount, zone.GetAggregateIntentOrDefault())

	requests = ApplyDeltasToIntent(requests, deltas, bins)

	for _, allocation := range requests.Sorted() {

		valoper := allocation.Address
		remainingTokens := allocation.Amount.AmountOf(types.GenericToken)

		_, valAddr, _ := bech32.DecodeAndConvert(valoper)

		delegations := k.GetValidatorDelegations(ctx, &zone, valAddr)
		sort.SliceStable(delegations, func(i, j int) bool {
			return delegations[i].Amount.Amount.LT(delegations[j].Amount.Amount)
		})

		for _, delegation := range delegations {
			if delegation.Amount.Amount.GTE(remainingTokens) {
				out = out.Add(delegation.DelegationAddress, delegation.ValidatorAddress, sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, remainingTokens)))
				remainingTokens = sdk.ZeroInt()
				break
			} else {
				val := delegation.Amount.Amount
				remainingTokens = remainingTokens.Sub(val)
				out = out.Add(delegation.DelegationAddress, delegation.ValidatorAddress, sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, val)))
			}
		}

		if remainingTokens.GT(sdk.ZeroInt()) {
			panic("redemption with remaining amount:" + remainingTokens.String())
		}

	}

	return out
}

func (k Keeper) InitPerformanceDelegations(ctx sdk.Context, zone types.Zone, response []byte) error {
	k.Logger(ctx).Info("Initialize performance delegations")

	resp := banktypes.QueryAllBalancesResponse{}
	err := k.cdc.Unmarshal(response, &resp)
	if err != nil {
		return err
	}
	k.Logger(ctx).Info("Performance Balance", "Account", zone.PerformanceAddress, "Balances", resp.Balances)

	if resp.Balances.IsZero() {
		// if zero balance, retrigger the query.
		if err := k.EmitPerformanceBalanceQuery(ctx, &zone); err != nil {
			return err
		}
		k.Logger(ctx).Info("performance account has a zero balance; requerying")
		return icqtypes.ErrSucceededNoDelete
	}

	amount := sdk.NewCoin(zone.BaseDenom, sdk.NewInt(10000))
	minBalance := sdk.NewInt(int64(len(zone.Validators)) * amount.Amount.Int64())
	balance := resp.Balances.AmountOfNoDenomValidation(zone.BaseDenom)
	if balance.LT(minBalance) {
		return fmt.Errorf(
			"performance account has an insufficient balance, got %v, expected at least %v",
			balance,
			minBalance,
		)
	}

	// send delegations to validators
	k.Logger(ctx).Info("send performance delegations", "zone", zone.ChainId)
	var msgs []sdk.Msg
	for _, val := range zone.Validators {
		k.Logger(ctx).Info(
			"performance delegation",
			"zone", zone.ChainId,
			"validator", val.ValoperAddress,
			"amount", amount,
		)
		msgs = append(msgs, &stakingtypes.MsgDelegate{
			DelegatorAddress: zone.PerformanceAddress.GetAddress(),
			ValidatorAddress: val.GetValoperAddress(),
			Amount:           amount,
		})
	}

	return k.SubmitTx(ctx, msgs, zone.PerformanceAddress, "")
}
