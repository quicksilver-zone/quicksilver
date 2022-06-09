package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var _ types.QueryServer = Keeper{}

// RegisteredZoneInfos returns information about registered zones.
func (k Keeper) RegisteredZoneInfos(c context.Context, req *types.QueryRegisteredZonesInfoRequest) (*types.QueryRegisteredZonesInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var zones []types.RegisteredZone
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var zone types.RegisteredZone
		if err := k.cdc.Unmarshal(value, &zone); err != nil {
			return err
		}
		zones = append(zones, zone)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRegisteredZonesInfoResponse{
		Zones:      zones,
		Pagination: pageRes,
	}, nil
}

// DepositAccount returns the deposit account address for the given zone.
func (k Keeper) DepositAccount(c context.Context, req *types.QueryDepositAccountForChainRequest) (*types.QueryDepositAccountForChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetRegisteredZoneInfo(ctx, req.GetChainId())
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
func (k Keeper) DelegatorIntent(c context.Context, req *types.QueryDelegatorIntentRequest) (*types.QueryDelegatorIntentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetRegisteredZoneInfo(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	intent, found := k.GetIntent(ctx, zone, req.DelegatorAddress, false)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no delegation intent specified for %s", req.GetChainId()))
	}

	return &types.QueryDelegatorIntentResponse{Intent: &intent}, nil
}

func (k Keeper) Delegations(c context.Context, req *types.QueryDelegationsRequest) (*types.QueryDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetRegisteredZoneInfo(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	delegations := k.GetAllDelegations(ctx, &zone)

	return &types.QueryDelegationsResponse{Delegations: delegations}, nil
}

func (k Keeper) DelegatorDelegations(c context.Context, req *types.QueryDelegatorDelegationsRequest) (*types.QueryDelegatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetRegisteredZoneInfo(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	_, addr, _ := bech32.DecodeAndConvert(req.DelegatorAddress)
	delegations := k.GetDelegatorDelegations(ctx, &zone, addr)

	return &types.QueryDelegatorDelegationsResponse{Delegations: delegations}, nil
}

func (k Keeper) ValidatorDelegations(c context.Context, req *types.QueryValidatorDelegationsRequest) (*types.QueryValidatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetRegisteredZoneInfo(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	_, valAddr, _ := bech32.DecodeAndConvert(req.ValidatorAddress)
	delegations := k.GetValidatorDelegations(ctx, &zone, valAddr)

	return &types.QueryValidatorDelegationsResponse{Delegations: delegations}, nil
}

func (k Keeper) DelegationPlans(c context.Context, req *types.QueryDelegationPlansRequest) (*types.QueryDelegationPlansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetRegisteredZoneInfo(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	delegationplans := k.GetAllDelegationPlans(ctx, &zone)

	return &types.QueryDelegationPlansResponse{Delegations: delegationplans}, nil
}
