package keeper

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
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
	out := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) (stop bool) {
		out = out.Add(delegation.Amount)
		return false
	})
	return out
}

func (k Keeper) GetUnbondingAmount(ctx sdk.Context, zone *types.Zone) sdk.Coin {
	out := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, WithdrawStatusUnbond, func(index int64, wr types.WithdrawalRecord) (stop bool) {
		out = out.Add(wr.Amount[0])
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
		if zoneInfo.DelegationAddress != nil && zoneInfo.DelegationAddress.Address == address {
			zone = &zoneInfo
			return true
		}
		return false
	})
	return zone
}

func (k Keeper) GetZoneForPerformanceAccount(ctx sdk.Context, address string) *types.Zone {
	var zone *types.Zone
	k.IterateZones(ctx, func(_ int64, zoneInfo types.Zone) (stop bool) {
		if zoneInfo.PerformanceAddress != nil && zoneInfo.PerformanceAddress.Address == address {
			zone = &zoneInfo
			return true
		}
		return false
	})
	return zone
}

func (k Keeper) GetZoneForDepositAccount(ctx sdk.Context, address string) *types.Zone {
	var zone *types.Zone
	k.IterateZones(ctx, func(_ int64, zoneInfo types.Zone) (stop bool) {
		if zoneInfo.DepositAddress != nil && zoneInfo.DepositAddress.Address == address {
			zone = &zoneInfo
			return true
		}
		return false
	})
	return zone
}

func (k Keeper) EnsureICAsActive(ctx sdk.Context, zone *types.Zone) error {
	k.Logger(ctx).Info("Ensuring ICAs for zone", "zone", zone.ChainId)
	if err := k.EnsureICAActive(ctx, zone, zone.DepositAddress); err != nil {
		return err
	}
	if err := k.EnsureICAActive(ctx, zone, zone.DelegationAddress); err != nil {
		return err
	}
	if err := k.EnsureICAActive(ctx, zone, zone.PerformanceAddress); err != nil {
		return err
	}
	if err := k.EnsureICAActive(ctx, zone, zone.WithdrawalAddress); err != nil {
		return err
	}
	return nil
}

func (k Keeper) EnsureICAActive(ctx sdk.Context, zone *types.Zone, account *types.ICAAccount) error {
	if account == nil {
		k.Logger(ctx).Info("Account does not exist")
		// address has not been set yet. nothing to check.
		return nil
	}

	if _, found := k.ICAControllerKeeper.GetOpenActiveChannel(ctx, zone.ConnectionId, account.GetPortName()); found {
		k.Logger(ctx).Info("Account is active", "account", account.Address)
		// channel is active. all is well :)
		return nil
	}

	// channel is not active; attempt reopen.
	k.Logger(ctx).Error("channel is inactive. attempting to reopen.", "connection", zone.ConnectionId, "port", account.GetPortName())
	return k.ICAControllerKeeper.RegisterInterchainAccount(ctx, zone.ConnectionId, strings.TrimPrefix(account.GetPortName(), icatypes.PortPrefix), "")
}

func (k *Keeper) EnsureWithdrawalAddresses(ctx sdk.Context, zone *types.Zone) error {
	if zone.WithdrawalAddress == nil {
		k.Logger(ctx).Info("Withdrawal address not set")
		return nil
	}

	if zone.DelegationAddress == nil {
		k.Logger(ctx).Info("Delegation address not set")
		return nil
	}

	if zone.DepositAddress == nil {
		k.Logger(ctx).Info("Deposit address not set")
		return nil
	}

	withdrawalAddress := zone.WithdrawalAddress.Address

	if zone.DepositAddress.WithdrawalAddress != withdrawalAddress {
		msg := distrTypes.MsgSetWithdrawAddress{DelegatorAddress: zone.DepositAddress.Address, WithdrawAddress: withdrawalAddress}
		err := k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.DepositAddress, "")
		if err != nil {
			return err
		}
	}

	if zone.DelegationAddress.WithdrawalAddress != withdrawalAddress {
		msg := distrTypes.MsgSetWithdrawAddress{DelegatorAddress: zone.DelegationAddress.Address, WithdrawAddress: withdrawalAddress}
		err := k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.DelegationAddress, "")
		if err != nil {
			return err
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
	var err error
	switch {
	case zone.DepositAddress != nil && address == zone.DepositAddress.Address:
		existing := zone.DepositAddress.Balance.AmountOf(coin.Denom)
		err = zone.DepositAddress.SetBalance(zone.DepositAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))...).Add(coin)) // reset this denom
		if err != nil {
			return err
		}
		err = zone.DepositAddress.DecrementBalanceWaitgroup()
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("Matched deposit address", "address", address, "wg", zone.DepositAddress.BalanceWaitgroup, "balance", zone.DepositAddress.Balance)
		if zone.DepositAddress.BalanceWaitgroup == 0 {
			k.depositInterval(ctx)(0, zone)
		}
	case zone.WithdrawalAddress != nil && address == zone.WithdrawalAddress.Address:
		existing := zone.WithdrawalAddress.Balance.AmountOf(coin.Denom)
		err = zone.WithdrawalAddress.SetBalance(zone.WithdrawalAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))...).Add(coin)) // reset this denom
		if err != nil {
			return err
		}
		err = zone.WithdrawalAddress.DecrementBalanceWaitgroup()
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("Matched withdrawal address", "address", address, "wg", zone.WithdrawalAddress.BalanceWaitgroup, "balance", zone.WithdrawalAddress.Balance)
	case zone.PerformanceAddress != nil && address == zone.PerformanceAddress.Address:
		k.Logger(ctx).Info("Matched performance address")
	default:
		panic("unexpected")
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
	_, addr, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return err
	}
	data := banktypes.CreateAccountBalancesPrefix(addr)

	var icaAccount *types.ICAAccount

	switch {
	case zone.DepositAddress != nil && address == zone.DepositAddress.Address:
		icaAccount = zone.DepositAddress
	case zone.WithdrawalAddress != nil && address == zone.WithdrawalAddress.Address:
		icaAccount = zone.WithdrawalAddress
	case zone.DelegationAddress != nil && address == zone.DelegationAddress.Address:
		icaAccount = zone.DelegationAddress
	case zone.PerformanceAddress != nil && address == zone.PerformanceAddress.Address:
		icaAccount = zone.PerformanceAddress
	default:
		return errors.New("unexpected address")
	}

	if icaAccount == nil {
		return fmt.Errorf("unable to determine account for address %s", address)
	}

	for _, coin := range icaAccount.Balance {
		if queryRes.Balances.AmountOf(coin.Denom).Equal(sdk.ZeroInt()) {
			// coin we used to have is now zero - also validate this.
			k.Logger(ctx).Info("Querying for value", "key", types.BankStoreKey, "denom", coin.Denom) // debug?
			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				types.BankStoreKey,
				append(data, []byte(coin.Denom)...),
				sdk.NewInt(-1),
				types.ModuleName,
				"accountbalance",
				0,
			)
			icaAccount.IncrementBalanceWaitgroup()

		}
	}

	for _, coin := range queryRes.Balances {
		k.Logger(ctx).Info("Querying for value", "key", types.BankStoreKey, "denom", coin.Denom) // debug?
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			types.BankStoreKey,
			append(data, []byte(coin.Denom)...),
			sdk.NewInt(-1),
			types.ModuleName,
			"accountbalance",
			0,
		)
		icaAccount.IncrementBalanceWaitgroup()
	}

	k.SetZone(ctx, &zone)
	return nil
}

func (k Keeper) UpdatePerformanceDelegations(ctx sdk.Context, zone types.Zone, response []byte) error {
	k.Logger(ctx).Info("Initialize performance delegations")

	delegations := k.GetAllPerformanceDelegations(ctx, &zone)
	validatorsToDelegate := []string{}
OUTER:
	for _, v := range zone.GetBondedValidatorAddressesAsSlice() {
		for _, d := range delegations {
			if d.ValidatorAddress == v {
				continue OUTER
			}
		}
		validatorsToDelegate = append(validatorsToDelegate, v)
	}

	amount := sdk.NewCoin(zone.BaseDenom, sdk.NewInt(10000))
	minBalance := sdk.NewInt(int64(len(validatorsToDelegate)) * amount.Amount.Int64())
	balance := zone.PerformanceAddress.Balance.AmountOfNoDenomValidation(zone.BaseDenom)
	if balance.LT(minBalance) {
		k.Logger(ctx).Error(
			fmt.Sprintf(
				"performance account has an insufficient balance, got %v, expected at least %v",
				balance,
				minBalance,
			),
		)
		return nil // don't error here, as we don't want the underlying tx to fail.
	}

	// send delegations to validators
	k.Logger(ctx).Info("send performance delegations", "zone", zone.ChainId)
	var msgs []sdk.Msg
	for _, val := range validatorsToDelegate {
		k.Logger(ctx).Info(
			"performance delegation",
			"zone", zone.ChainId,
			"validator", val,
			"amount", amount,
		)
		msgs = append(msgs, &stakingtypes.MsgDelegate{
			DelegatorAddress: zone.PerformanceAddress.GetAddress(),
			ValidatorAddress: val,
			Amount:           amount,
		})
	}

	if len(msgs) > 0 {
		return k.SubmitTx(ctx, msgs, zone.PerformanceAddress, "")
	}
	return nil
}

func (k *Keeper) CollectStatsForZone(ctx sdk.Context, zone *types.Zone) *types.Statistics {
	out := &types.Statistics{}
	out.ChainId = zone.ChainId
	out.Delegated = k.GetDelegatedAmount(ctx, zone).Amount.Int64()
	userMap := map[string]bool{}
	k.IterateZoneReceipts(ctx, zone, func(_ int64, receipt types.Receipt) bool {
		for _, coin := range receipt.Amount {
			out.Deposited += coin.Amount.Int64()
			if _, found := userMap[receipt.Sender]; !found {
				userMap[receipt.Sender] = true
				out.Depositors++
			}
			out.Deposits++
		}
		return false
	})
	out.Supply = k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.Int64()
	out.DistanceToTarget = fmt.Sprintf("%f", k.DistanceToTarget(ctx, zone))
	return out
}

func (k *Keeper) RemoveZoneAndAssociatedRecords(ctx sdk.Context, chainID string) {
	// clear unbondings
	k.IteratePrefixedUnbondingRecords(ctx, []byte(chainID), func(_ int64, record types.UnbondingRecord) (stop bool) {
		k.DeleteUnbondingRecord(ctx, record.ChainId, record.Validator, record.EpochNumber)
		return false
	})

	// clear redelegations
	k.IteratePrefixedRedelegationRecords(ctx, []byte(chainID), func(_ int64, _ []byte, record types.RedelegationRecord) (stop bool) {
		k.DeleteRedelegationRecord(ctx, record.ChainId, record.Source, record.Destination, record.EpochNumber)
		return false
	})

	// remove zone and related records
	k.IterateZones(ctx, func(index int64, zone types.Zone) (stop bool) {
		if zone.ChainId == chainID {
			// remove uni-5 delegation records
			k.IterateAllDelegations(ctx, &zone, func(delegation types.Delegation) (stop bool) {
				err := k.RemoveDelegation(ctx, &zone, delegation)
				if err != nil {
					panic(err)
				}
				return false
			})

			// remove performance delegation records
			k.IterateAllPerformanceDelegations(ctx, &zone, func(delegation types.Delegation) (stop bool) {
				err := k.RemoveDelegation(ctx, &zone, delegation)
				if err != nil {
					panic(err)
				}
				return false
			})
			// remove receipts
			k.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
				k.DeleteReceipt(ctx, GetReceiptKey(receiptInfo.ChainId, receiptInfo.Txhash))
				return false
			})

			// remove withdrawal records
			k.IterateZoneWithdrawalRecords(ctx, zone.ChainId, func(index int64, record types.WithdrawalRecord) (stop bool) {
				k.DeleteWithdrawalRecord(ctx, zone.ChainId, record.Txhash, record.Status)
				return false
			})

			k.DeleteZone(ctx, zone.ChainId)

		}
		return false
	})

	// remove queries in state
	k.ICQKeeper.IterateQueries(ctx, func(_ int64, queryInfo icqtypes.Query) (stop bool) {
		if queryInfo.ChainId == chainID {
			k.ICQKeeper.DeleteQuery(ctx, queryInfo.Id)
		}
		return false
	})
}

func (k *Keeper) CurrentDelegationsAsIntent(ctx sdk.Context, zone *types.Zone) types.ValidatorIntents {
	currentDelegations := k.GetAllDelegations(ctx, zone)
	intents := make(types.ValidatorIntents, 0)
	for _, d := range currentDelegations {
		intents = append(intents, &types.ValidatorIntent{ValoperAddress: d.ValidatorAddress, Weight: sdk.NewDecFromInt(d.Amount.Amount)})
	}

	return intents.Normalize()
}

func (k *Keeper) DistanceToTarget(ctx sdk.Context, zone *types.Zone) float64 {
	current := k.CurrentDelegationsAsIntent(ctx, zone)
	target := zone.GetAggregateIntentOrDefault()
	preSqRt := sdk.ZeroDec()

	for _, valoper := range zone.Validators {
		c := current.MustGetForValoper(valoper.ValoperAddress)
		t := target.MustGetForValoper(valoper.ValoperAddress)
		v := c.Weight.SubMut(t.Weight)
		preSqRt = preSqRt.AddMut(v.Mul(v))
	}

	psqrtf, err := preSqRt.Float64()
	if err != nil {
		panic("this value should never be greater than 64-bit dec!")
	}
	return math.Sqrt(psqrtf)
}
