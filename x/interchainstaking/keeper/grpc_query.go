package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var _ types.QueryServer = &Keeper{}

// Zones returns information about registered zones.
func (k *Keeper) Zones(c context.Context, req *types.QueryZonesInfoRequest) (*types.QueryZonesInfoResponse, error) {
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
		stats = append(stats, k.CollectStatsForZone(ctx, &zone))
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryZonesInfoResponse{
		Zones:      zones,
		Stats:      stats,
		Pagination: pageRes,
	}, nil
}

// Zone returns information about registered zones.
func (k *Keeper) Zone(c context.Context, req *types.QueryZoneInfoRequest) (*types.QueryZoneInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetZone(ctx, req.ChainId)
	if !found {
		return nil, fmt.Errorf("no zone found for chain id %s", req.ChainId)
	}

	return &types.QueryZoneInfoResponse{
		Zone:  zone,
		Stats: k.CollectStatsForZone(ctx, &zone),
	}, nil
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
	// - see comment in GetIntent
	intent, _ := k.GetIntent(ctx, &zone, req.DelegatorAddress, false)

	return &types.QueryDelegatorIntentResponse{Intent: &intent}, nil
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
