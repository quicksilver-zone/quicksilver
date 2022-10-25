package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

var _ types.QueryServer = Keeper{}

// Params returns params of the claimsmanager module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) Claims(c context.Context, q *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	out := []types.Claim{}
	k.Logger(ctx).Error("Claims query")
	k.IterateClaims(ctx, q.ChainId, func(_ int64, claim types.Claim) (stop bool) {
		k.Logger(ctx).Error("Claim", claim)

		out = append(out, claim)
		return false
	})
	k.Logger(ctx).Error("Romeo done.")

	return &types.QueryClaimsResponse{Claims: out}, nil
}

func (k Keeper) LastEpochClaims(c context.Context, q *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	out := []types.Claim{}

	k.IterateLastEpochClaims(ctx, q.ChainId, func(_ int64, claim types.Claim) (stop bool) {
		out = append(out, claim)
		return false
	})

	return &types.QueryClaimsResponse{Claims: out}, nil
}

func (k Keeper) UserClaims(c context.Context, q *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	out := []types.Claim{}

	k.IterateUserClaims(ctx, q.ChainId, q.Address, func(_ int64, claim types.Claim) (stop bool) {
		out = append(out, claim)
		return false
	})

	return &types.QueryClaimsResponse{Claims: out}, nil
}

func (k Keeper) UserLastEpochClaims(c context.Context, q *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	out := []types.Claim{}

	k.IterateLastEpochUserClaims(ctx, q.ChainId, q.Address, func(_ int64, claim types.Claim) (stop bool) {
		out = append(out, claim)
		return false
	})

	return &types.QueryClaimsResponse{Claims: out}, nil
}
