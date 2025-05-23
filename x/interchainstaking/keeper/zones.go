package keeper

import (
	"errors"
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// GetZone returns zone info by chainID.
func (k *Keeper) GetZone(ctx sdk.Context, chainID string) (types.Zone, bool) {
	zone := types.Zone{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	bz := store.Get([]byte(chainID))
	if len(bz) == 0 {
		return zone, false
	}

	k.cdc.MustUnmarshal(bz, &zone)
	return zone, true
}

// SetZone set zone info.
func (k *Keeper) SetZone(ctx sdk.Context, zone *types.Zone) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	bz := k.cdc.MustMarshal(zone)
	store.Set([]byte(zone.ChainId), bz)
}

// DeleteZone delete zone info.
func (k *Keeper) DeleteZone(ctx sdk.Context, chainID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)
	store.Delete([]byte(chainID))
}

// IterateZones iterate through zones.
func (k *Keeper) IterateZones(ctx sdk.Context, fn func(index int64, zone *types.Zone) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		zone := types.Zone{}
		k.cdc.MustUnmarshal(iterator.Value(), &zone)

		stop := fn(i, &zone)

		if stop {
			break
		}
		i++
	}
}

// GetAddressZoneMapping returns zone <-> address mapping.
func (k *Keeper) GetAddressZoneMapping(ctx sdk.Context, address string) (string, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAddressZoneMapping)
	bz := store.Get([]byte(address))
	if len(bz) == 0 {
		return "", false
	}
	return string(bz), true
}

// SetAddressZoneMapping set zone <-> address mapping.
func (k *Keeper) SetAddressZoneMapping(ctx sdk.Context, address, chainID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAddressZoneMapping)
	store.Set([]byte(address), []byte(chainID))
}

// DeleteAddressZoneMapping delete zone info.
func (k *Keeper) DeleteAddressZoneMapping(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAddressZoneMapping)
	store.Delete([]byte(address))
}

func (k *Keeper) GetDelegatedAmount(ctx sdk.Context, zone *types.Zone) sdk.Coin {
	out := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	k.IterateAllDelegations(ctx, zone.ChainId, func(delegation types.Delegation) (stop bool) {
		out = out.Add(delegation.Amount)
		return false
	})
	return out
}

// GetUnbondingTokensAndCount returns the total amount of unbonding tokens and the count of unbonding for a given zone.
func (k *Keeper) GetUnbondingTokens(ctx sdk.Context, zone *types.Zone) sdk.Coin {
	out := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	k.IterateUnbondingRecords(ctx, func(index int64, wr types.UnbondingRecord) (stop bool) {
		if wr.ChainId != zone.ChainId {
			return false
		}
		amount := wr.Amount
		if !amount.IsNegative() {
			out = out.Add(amount)
		}
		return false
	})
	return out
}

// GetWithdrawingTokensAndCount return the total amount of unbonding tokens and the count of unbonding for a given zone.
func (k *Keeper) GetWithdrawnTokensAndCount(ctx sdk.Context, zone *types.Zone) (sdk.Coin, uint32) {
	out := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	var count uint32
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusUnbond, func(index int64, wr types.WithdrawalRecord) (stop bool) {
		amount := wr.Amount[0]
		if !amount.IsNegative() {
			out = out.Add(amount)
		}
		count++
		return false
	})
	return out, count
}

func (k *Keeper) GetQueuedTokensAndCount(ctx sdk.Context, zone *types.Zone) (sdk.Coin, uint32) {
	out := sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt())
	var count uint32
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusQueued, func(index int64, wr types.WithdrawalRecord) (stop bool) {
		if !wr.BurnAmount.IsNegative() {
			out = out.Add(wr.BurnAmount)
		}
		count++
		return false
	})
	return out, count
}

func (k *Keeper) GetUnbondRecordCount(ctx sdk.Context, zone *types.Zone) uint32 {
	var count uint32
	k.IteratePrefixedUnbondingRecords(ctx, []byte(zone.ChainId), func(_ int64, record types.UnbondingRecord) (stop bool) {
		count++
		return false
	})
	return count
}

// AllZones returns every Zone in the store.
func (k *Keeper) AllZones(ctx sdk.Context) []types.Zone {
	var zones []types.Zone
	k.IterateZones(ctx, func(_ int64, zone *types.Zone) (stop bool) {
		zones = append(zones, *zone)
		return false
	})
	return zones
}

// GetZoneFromConnectionID determines the zone from the connection ID
func (k *Keeper) GetZoneFromConnectionID(ctx sdk.Context, connectionID string) (*types.Zone, error) {
	chainID, err := k.GetChainID(ctx, connectionID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch zone from connection id: %w", err)
	}
	zone, found := k.GetZone(ctx, chainID)
	if !found {
		err := fmt.Errorf("unable to fetch zone from connection id: not found for chainID %s", chainID)
		k.Logger(ctx).Error(err.Error())
		return nil, err
	}
	return &zone, nil
}

func (k *Keeper) GetZoneForAccount(ctx sdk.Context, address string) (*types.Zone, bool) {
	chainID, found := k.GetAddressZoneMapping(ctx, address)
	if !found {
		return nil, false
	}

	zone, found := k.GetZone(ctx, chainID)
	return &zone, found
}

// GetZoneForDelegateAccount determines the zone for a given address.
func (k *Keeper) GetZoneForDelegateAccount(ctx sdk.Context, address string) (*types.Zone, bool) {
	zone, found := k.GetZoneForAccount(ctx, address)
	if !found {
		return nil, false // address not found
	}
	if zone.DelegationAddress != nil && address == zone.DelegationAddress.Address {
		return zone, true // address found and is delegate Account
	}
	return nil, false // address found, but not delegate account
}

func (k *Keeper) GetZoneForPerformanceAccount(ctx sdk.Context, address string) (*types.Zone, bool) {
	zone, found := k.GetZoneForAccount(ctx, address)
	if !found {
		return nil, false // address not found
	}
	if zone.PerformanceAddress != nil && address == zone.PerformanceAddress.Address {
		return zone, true // address found and is performance Account
	}
	return nil, false // address found, but not performance account
}

func (k *Keeper) GetZoneForDepositAccount(ctx sdk.Context, address string) (*types.Zone, bool) {
	zone, found := k.GetZoneForAccount(ctx, address)
	if !found {
		return nil, false // address not found
	}
	if zone.DepositAddress != nil && address == zone.DepositAddress.Address {
		return zone, true // address found and is deposit Account
	}
	return nil, false // address found, but not deposit account
}

func (k *Keeper) GetZoneForWithdrawalAccount(ctx sdk.Context, address string) (*types.Zone, bool) {
	zone, found := k.GetZoneForAccount(ctx, address)
	if !found {
		return nil, false // address not found
	}
	if zone.WithdrawalAddress != nil && address == zone.WithdrawalAddress.Address {
		return zone, true // address found and is withdrawal Account
	}
	return nil, false // address found, but not withdrawal account
}

func (k *Keeper) GetICAAccountForAddress(ctx sdk.Context, address string) (*types.ICAAccount, *types.Zone, error) {
	zone, found := k.GetZoneForAccount(ctx, address)
	if !found {
		return nil, nil, errors.New("address not found") // address not found
	}

	switch {
	case zone.DepositAddress != nil && address == zone.DepositAddress.Address:
		return zone.DepositAddress, zone, nil
	case zone.WithdrawalAddress != nil && address == zone.WithdrawalAddress.Address:
		return zone.WithdrawalAddress, zone, nil
	case zone.DelegationAddress != nil && address == zone.DelegationAddress.Address:
		return zone.DelegationAddress, zone, nil
	case zone.PerformanceAddress != nil && address == zone.PerformanceAddress.Address:
		return zone.PerformanceAddress, zone, nil
	default:
		return nil, zone, errors.New("unexpected account type")
	}
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

	if zone.PerformanceAddress == nil {
		k.Logger(ctx).Info("Performance address not set")
		return nil
	}

	withdrawalAddress := zone.WithdrawalAddress.Address

	if zone.DepositAddress.WithdrawalAddress != withdrawalAddress {
		msg := distrtypes.MsgSetWithdrawAddress{DelegatorAddress: zone.DepositAddress.Address, WithdrawAddress: withdrawalAddress}
		err := k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.DepositAddress, "", zone.MessagesPerTx)
		if err != nil {
			return err
		}
	}

	if zone.DelegationAddress.WithdrawalAddress != withdrawalAddress {
		msg := distrtypes.MsgSetWithdrawAddress{DelegatorAddress: zone.DelegationAddress.Address, WithdrawAddress: withdrawalAddress}
		err := k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.DelegationAddress, "", zone.MessagesPerTx)
		if err != nil {
			return err
		}
	}

	// set withdrawal address for performance address, if it exists
	if zone.PerformanceAddress != nil && zone.PerformanceAddress.WithdrawalAddress != withdrawalAddress {
		msg := distrtypes.MsgSetWithdrawAddress{DelegatorAddress: zone.PerformanceAddress.Address, WithdrawAddress: withdrawalAddress}
		err := k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.PerformanceAddress, "", zone.MessagesPerTx)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetAccountBalanceForDenom sets the balance on an account for a given denomination.
func (k *Keeper) SetAccountBalanceForDenom(ctx sdk.Context, zone *types.Zone, address string, coin sdk.Coin) error {
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
		existing := zone.PerformanceAddress.Balance.AmountOf(coin.Denom)
		err = zone.PerformanceAddress.SetBalance(zone.PerformanceAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))...).Add(coin)) // reset this denom
		if err != nil {
			return err
		}
		err = zone.PerformanceAddress.DecrementBalanceWaitgroup()
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("Matched performance address", "address", address, "wg", zone.PerformanceAddress.BalanceWaitgroup, "balance", zone.PerformanceAddress.Balance)
	case zone.DelegationAddress != nil && address == zone.DelegationAddress.Address:
		existing := zone.DelegationAddress.Balance.AmountOf(coin.Denom)
		err = zone.DelegationAddress.SetBalance(zone.DelegationAddress.Balance.Sub(sdk.NewCoins(sdk.NewCoin(coin.Denom, existing))...).Add(coin)) // reset this denom
		if err != nil {
			return err
		}
		err = zone.DelegationAddress.DecrementBalanceWaitgroup()
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("Matched delegation address", "address", address, "wg", zone.DelegationAddress.BalanceWaitgroup, "balance", zone.DelegationAddress.Balance)
	default:
		panic("unexpected")
	}
	k.SetZone(ctx, zone)
	return nil
}

// SetAccountBalance triggers provable KV queries to prove an AllBalances query.
func (k *Keeper) SetAccountBalance(ctx sdk.Context, zone types.Zone, address string, queryResult []byte) error {
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
				append(data, coin.Denom...),
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
			append(data, coin.Denom...),
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

func (k *Keeper) UpdatePerformanceDelegations(ctx sdk.Context, zone types.Zone) error {
	k.Logger(ctx).Info("Initialize performance delegations")

	delegations := k.GetAllPerformanceDelegations(ctx, zone.ChainId)
	validatorsToDelegate := []string{}
OUTER:
	for _, v := range k.GetActiveValidators(ctx, zone.ChainId) {
		for _, d := range delegations {
			if d.ValidatorAddress == v.ValoperAddress {
				continue OUTER
			}
		}
		validatorsToDelegate = append(validatorsToDelegate, v.ValoperAddress)
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

	msgs := make([]sdk.Msg, len(validatorsToDelegate))
	for i, val := range validatorsToDelegate {
		k.Logger(ctx).Info(
			"performance delegation",
			"zone", zone.ChainId,
			"validator", val,
			"amount", amount,
		)
		msgs[i] = &stakingtypes.MsgDelegate{
			DelegatorAddress: zone.PerformanceAddress.GetAddress(),
			ValidatorAddress: val,
			Amount:           amount,
		}
	}

	if len(msgs) > 0 {
		return k.SubmitTx(ctx, msgs, zone.PerformanceAddress, "", zone.MessagesPerTx)
	}
	return nil
}

func (k *Keeper) CollectStatsForZone(ctx sdk.Context, zone *types.Zone) (*types.Statistics, error) {
	out := &types.Statistics{}
	out.ChainId = zone.ChainId
	out.Deposited = sdkmath.ZeroInt()
	out.Delegated = k.GetDelegatedAmount(ctx, zone).Amount
	userMap := map[string]bool{}
	k.IterateZoneReceipts(ctx, zone.ChainId, func(_ int64, receipt types.Receipt) bool {
		for _, coin := range receipt.Amount {
			out.Deposited = out.Deposited.Add(coin.Amount)
			if _, found := userMap[receipt.Sender]; !found {
				userMap[receipt.Sender] = true
				out.Depositors++
			}
			out.Deposits++
		}
		return false
	})
	out.Supply = k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount
	distance, err := k.DistanceToTarget(ctx, zone)
	if err != nil {
		return nil, err
	}
	out.DistanceToTarget = fmt.Sprintf("%f", distance)

	// Unbonding info
	out.UnbondingAmount, out.UnbondingCount = k.GetWithdrawnTokensAndCount(ctx, zone)
	out.QueuedAmount, out.QueuedCount = k.GetQueuedTokensAndCount(ctx, zone)
	out.UnbondRecordCount = k.GetUnbondRecordCount(ctx, zone)
	return out, nil
}

func (k *Keeper) RemoveZoneAndAssociatedRecords(ctx sdk.Context, chainID string) {
	// remove zone and related records
	zone, ok := k.GetZone(ctx, chainID)
	if !ok {
		panic("cannot find zone for deletion")
	}

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

	// remove delegation records
	k.IterateAllDelegations(ctx, chainID, func(delegation types.Delegation) (stop bool) {
		err := k.RemoveDelegation(ctx, chainID, delegation)
		if err != nil {
			panic(err)
		}
		return false
	})

	// remove performance delegation records
	k.IterateAllPerformanceDelegations(ctx, chainID, func(delegation types.Delegation) (stop bool) {
		err := k.RemovePerformanceDelegation(ctx, chainID, delegation)
		if err != nil {
			panic(err)
		}
		return false
	})

	// remove receipts
	k.IterateZoneReceipts(ctx, chainID, func(index int64, receiptInfo types.Receipt) (stop bool) {
		k.DeleteReceipt(ctx, chainID, receiptInfo.Txhash)
		return false
	})

	// remove withdrawal records
	k.IterateZoneWithdrawalRecords(ctx, chainID, func(index int64, record types.WithdrawalRecord) (stop bool) {
		k.DeleteWithdrawalRecord(ctx, chainID, record.Txhash, record.Status)
		return false
	})

	// remove validators
	k.IterateValidators(ctx, chainID, func(index int64, validator types.Validator) (stop bool) {
		valAddr, err := validator.GetAddressBytes()
		if err != nil {
			panic(err)
		}
		k.DeleteValidator(ctx, chainID, valAddr)
		return false
	})

	k.IteratePortConnections(ctx, func(pc types.PortConnectionTuple) (stop bool) {
		if pc.ConnectionId == zone.ConnectionId {
			k.DeleteConnectionForPort(ctx, pc.PortId)
		}
		return false
	})

	k.DeleteDenomZoneMapping(ctx, zone.LocalDenom)

	k.DeleteZone(ctx, zone.ChainId)

	// remove queries in state
	k.ICQKeeper.IterateQueries(ctx, func(_ int64, queryInfo icqtypes.Query) (stop bool) {
		if queryInfo.ChainId == chainID {
			k.ICQKeeper.DeleteQuery(ctx, queryInfo.Id)
		}
		return false
	})

	// remove claims
	k.ClaimsManagerKeeper.IterateClaims(ctx, chainID, func(index int64, data claimsmanagertypes.Claim) (stop bool) {
		k.ClaimsManagerKeeper.DeleteClaim(ctx, &data)
		return false
	})

	k.ClaimsManagerKeeper.IterateLastEpochClaims(ctx, chainID, func(index int64, data claimsmanagertypes.Claim) (stop bool) {
		k.ClaimsManagerKeeper.DeleteClaim(ctx, &data)
		return false
	})
}

func (k *Keeper) CurrentDelegationsAsIntent(ctx sdk.Context, zone *types.Zone) types.ValidatorIntents {
	currentDelegations := k.GetAllDelegations(ctx, zone.ChainId)
	intents := make(types.ValidatorIntents, 0)
	for _, d := range currentDelegations {
		intents = append(intents, &types.ValidatorIntent{ValoperAddress: d.ValidatorAddress, Weight: sdk.NewDecFromInt(d.Amount.Amount)})
	}

	return intents.Normalize()
}

func (k *Keeper) DistanceToTarget(ctx sdk.Context, zone *types.Zone) (float64, error) {
	current := k.CurrentDelegationsAsIntent(ctx, zone)
	target, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return 0, err
	}
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
	return math.Sqrt(psqrtf), nil
}

// DefaultAggregateIntents determines the default aggregate intent (for epoch 0).
func (k *Keeper) DefaultAggregateIntents(ctx sdk.Context, chainID string) types.ValidatorIntents {
	out := make(types.ValidatorIntents, 0)
	k.IterateValidators(ctx, chainID, func(index int64, validator types.Validator) (stop bool) {
		if validator.CommissionRate.LTE(sdk.NewDecWithPrec(5, 1)) { // 50%; make this a param.
			if !validator.Jailed && !validator.Tombstoned && validator.Status == stakingtypes.BondStatusBonded {
				out = append(out, &types.ValidatorIntent{ValoperAddress: validator.GetValoperAddress(), Weight: sdk.OneDec()})
			}
		}
		return false
	})

	valCount := sdk.NewInt(int64(len(out)))

	// normalise the array (divide everything by length of intent list)
	for idx, intent := range out.Sort() {
		out[idx].Weight = intent.Weight.Quo(sdk.NewDecFromInt(valCount))
	}

	return out
}
