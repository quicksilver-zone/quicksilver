package keeper

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
