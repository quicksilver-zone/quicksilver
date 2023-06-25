package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var _ types.QueryServer = &Keeper{}

// Zones returns information about registered zones.
func (k *Keeper) Zones(c context.Context, req *types.QueryZonesRequest) (*types.QueryZonesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var zones []types.Zone
	var stats []*types.Statistics
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var zone types.Zone
		if err := k.cdc.Unmarshal(value, &zone); err != nil {
			return err
		}
		zones = append(zones, zone)
		zoneStats, err := k.CollectStatsForZone(ctx, &zone)
		if err != nil {
			return err
		}
		stats = append(stats, zoneStats)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryZonesResponse{
		Zones:      zones,
		Stats:      stats,
		Pagination: pageRes,
	}, nil
}

// Zone returns information about registered zones.
func (k *Keeper) Zone(c context.Context, req *types.QueryZoneRequest) (*types.QueryZoneResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetZone(ctx, req.ChainId)
	if !found {
		return nil, fmt.Errorf("no zone found for chain id %s", req.ChainId)
	}

	zoneStats, err := k.CollectStatsForZone(ctx, &zone)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryZoneResponse{
		Zone:  zone,
		Stats: zoneStats,
	}, nil
}

func (k Keeper) ZoneValidators(c context.Context, req *types.QueryZoneValidatorsRequest) (*types.QueryZoneValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var validators []types.Validator
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorsKey(req.ChainId))

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var validator types.Validator
		if err := k.cdc.Unmarshal(value, &validator); err != nil {
			return err
		}

		if req.Status == "" || req.Status == validator.Status {
			validators = append(validators, validator)
		}
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryZoneValidatorsResponse{Validators: validators, Pagination: pageRes}, nil
}

// DepositAccount returns the deposit account address for the given zone.
func (k *Keeper) DepositAccount(c context.Context, req *types.QueryDepositAccountForChainRequest) (*types.QueryDepositAccountForChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetZone(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	if zone.DepositAddress == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no deposit address registered yet for %s", req.GetChainId()))
	}

	return &types.QueryDepositAccountForChainResponse{
		DepositAccountAddress: zone.DepositAddress.Address,
	}, nil
}

// DelegatorIntent returns information about the delegation intent of the caller for the given zone.
func (k *Keeper) DelegatorIntent(c context.Context, req *types.QueryDelegatorIntentRequest) (*types.QueryDelegatorIntentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetZone(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	// we can ignore bool (found) as it always returns true
	// - see comment in GetDelegatorIntent
	intent, _ := k.GetDelegatorIntent(ctx, &zone, req.DelegatorAddress, false)

	return &types.QueryDelegatorIntentResponse{Intent: &intent}, nil
}

// DelegatorIntents returns information about the delegation intent of the given delegator for all zones.
func (k *Keeper) DelegatorIntents(c context.Context, req *types.QueryDelegatorIntentsRequest) (*types.QueryDelegatorIntentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	intents := []*types.DelegatorIntentsResponse{}

	k.IterateZones(ctx, func(_ int64, zone *types.Zone) bool {
		intent, _ := k.GetDelegatorIntent(ctx, zone, req.DelegatorAddress, false)
		if intent.Intents != nil {
			intents = append(intents, &types.DelegatorIntentsResponse{ChainId: zone.ChainId, Intent: &intent})
		}
		return false
	})

	return &types.QueryDelegatorIntentsResponse{Intents: intents}, nil
}

func (k *Keeper) Delegations(c context.Context, req *types.QueryDelegationsRequest) (*types.QueryDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetZone(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	delegations := make([]types.Delegation, 0)
	var sum int64

	k.IterateAllDelegations(ctx, &zone, func(delegation types.Delegation) (stop bool) {
		delegations = append(delegations, delegation)
		sum += delegation.Amount.Amount.Int64()
		return false
	})

	return &types.QueryDelegationsResponse{Delegations: delegations, Tvl: sum}, nil
}

func (k *Keeper) Receipts(c context.Context, req *types.QueryReceiptsRequest) (*types.QueryReceiptsResponse, error) {
	// TODO: implement pagination
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetZone(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	receipts := make([]types.Receipt, 0)

	k.IterateZoneReceipts(ctx, &zone, func(_ int64, receipt types.Receipt) (stop bool) {
		receipts = append(receipts, receipt)
		return false
	})

	return &types.QueryReceiptsResponse{Receipts: receipts}, nil
}

func (k *Keeper) TxStatus(c context.Context, req *types.QueryTxStatusRequest) (*types.QueryTxStatusResponse, error) {
	// TODO: implement pagination
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.GetTxHash() == "" {
		return nil, status.Error(codes.InvalidArgument, "tx hash cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	txReceipt, found := k.GetReceipt(ctx, types.GetReceiptKey(req.GetChainId(), req.GetTxHash()))
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no receipt found matching %s", req.TxHash))
	}

	return &types.QueryTxStatusResponse{Receipt: &txReceipt}, nil
}

func (k *Keeper) ZoneWithdrawalRecords(c context.Context, req *types.QueryWithdrawalRecordsRequest) (*types.QueryWithdrawalRecordsResponse, error) {
	// TODO: implement pagination
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetZone(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	withdrawalrecords := make([]types.WithdrawalRecord, 0)
	k.IterateZoneWithdrawalRecords(ctx, zone.ChainId, func(index int64, record types.WithdrawalRecord) (stop bool) {
		if record.Delegator == req.DelegatorAddress {
			withdrawalrecords = append(withdrawalrecords, record)
		}
		return false
	})

	return &types.QueryWithdrawalRecordsResponse{Withdrawals: withdrawalrecords}, nil
}

func (k *Keeper) WithdrawalRecords(c context.Context, req *types.QueryWithdrawalRecordsRequest) (*types.QueryWithdrawalRecordsResponse, error) {
	// TODO: implement pagination
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	withdrawalrecords := k.AllZoneWithdrawalRecords(ctx, req.ChainId)

	return &types.QueryWithdrawalRecordsResponse{Withdrawals: withdrawalrecords}, nil
}

func (k *Keeper) UserWithdrawalRecords(c context.Context, req *types.QueryUserWithdrawalRecordsRequest) (*types.QueryWithdrawalRecordsResponse, error) {
	// TODO: implement pagination
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if _, err := addressutils.AddressFromBech32(req.UserAddress, ""); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	withdrawalrecords := k.AllUserWithdrawalRecords(ctx, req.UserAddress)

	return &types.QueryWithdrawalRecordsResponse{Withdrawals: withdrawalrecords}, nil
}

func (k *Keeper) UnbondingRecords(c context.Context, req *types.QueryUnbondingRecordsRequest) (*types.QueryUnbondingRecordsResponse, error) {
	// TODO: implement pagination
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	unbondings := k.AllZoneUnbondingRecords(ctx, req.ChainId)

	return &types.QueryUnbondingRecordsResponse{Unbondings: unbondings}, nil
}

func (k *Keeper) RedelegationRecords(c context.Context, req *types.QueryRedelegationRecordsRequest) (*types.QueryRedelegationRecordsResponse, error) {
	// TODO: implement pagination
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	redelegations := k.ZoneRedelegationRecords(ctx, req.ChainId)

	return &types.QueryRedelegationRecordsResponse{Redelegations: redelegations}, nil
}

func (k *Keeper) MappedAccounts(c context.Context, req *types.QueryMappedAccountsRequest) (*types.QueryMappedAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	remoteAddressMap := make(map[string][]byte)
	addrBytes, err := addressutils.AccAddressFromBech32(req.Address, sdk.GetConfig().GetBech32AccountAddrPrefix())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Address")
	}

	k.IterateUserMappedAccounts(ctx, addrBytes, func(index int64, chainID string, remoteAddressBytes []byte) (stop bool) {
		remoteAddressMap[chainID] = remoteAddressBytes
		return false
	})

	return &types.QueryMappedAccountsResponse{RemoteAddressMap: remoteAddressMap}, nil
}
