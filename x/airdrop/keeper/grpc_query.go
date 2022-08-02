package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

var _ types.QueryServer = Keeper{}

// Params returns params of the airdrop module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// ZoneDrop returns the details of the specified zone airdrop.
func (k Keeper) ZoneDrop(c context.Context, req *types.QueryZoneDropRequest) (*types.QueryZoneDropResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	zd, ok := k.GetZoneDrop(ctx, req.ChainId)
	if !ok {
		return nil, types.ErrZoneDropNotFound
	}

	return &types.QueryZoneDropResponse{ZoneDrop: zd}, nil
}

// AccountBalance returns the airdrop module account balance of the specified zone.
func (k Keeper) AccountBalance(c context.Context, req *types.QueryAccountBalanceRequest) (*types.QueryAccountBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	ab := k.GetZoneDropAccountBalance(ctx, req.ChainId)

	return &types.QueryAccountBalanceResponse{
		AccountBalance: &ab,
	}, nil
}

// ZoneDrops returns all zone airdrops of the specified status.
func (k Keeper) ZoneDrops(c context.Context, req *types.QueryZoneDropsRequest) (*types.QueryZoneDropsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var zds []types.ZoneDrop
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZoneDrop)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var zd types.ZoneDrop
		if err := k.cdc.Unmarshal(value, &zd); err != nil {
			return err
		}

		switch req.Status {
		case types.StatusActive:
			if k.IsActiveZoneDrop(ctx, zd) {
				zds = append(zds, zd)
			}
		case types.StatusFuture:
			if k.IsFutureZoneDrop(ctx, zd) {
				zds = append(zds, zd)
			}
		case types.StatusExpired:
			if k.IsExpiredZoneDrop(ctx, zd) {
				zds = append(zds, zd)
			}
		default:
			// unknown status no-op
		}

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryZoneDropsResponse{
		ZoneDrops:  zds,
		Pagination: pageRes,
	}, nil
}

// ClaimRecord returns the claim record that corresponds to the given zone and address.
func (k Keeper) ClaimRecord(c context.Context, req *types.QueryClaimRecordRequest) (*types.QueryClaimRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	cr, err := k.GetClaimRecord(ctx, req.ChainId, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.QueryClaimRecordResponse{ClaimRecord: &cr}, nil
}

// ClaimRecords returns all the claim records of the given zone.
func (k Keeper) ClaimRecords(c context.Context, req *types.QueryClaimRecordsRequest) (*types.QueryClaimRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var crs []types.ClaimRecord
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetPrefixClaimRecord(req.ChainId))

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var cr types.ClaimRecord
		if err := k.cdc.Unmarshal(value, &cr); err != nil {
			return err
		}
		crs = append(crs, cr)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryClaimRecordsResponse{
		ClaimRecords: crs,
		Pagination:   pageRes,
	}, nil
}
