package keeper

import (
	"bytes"
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Claims(c context.Context, req *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var claims []types.Claim
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixClaim)

	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(_, value []byte, accumulate bool) (bool, error) {
		var claim types.Claim
		if err := k.cdc.Unmarshal(value, &claim); err != nil {
			return false, err
		}

		if claim.ChainId == req.ChainId {
			claims = append(claims, claim)
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryClaimsResponse{
		Claims:     claims,
		Pagination: pageRes,
	}, nil
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

	addrBytes := []byte(q.Address)

	k.IterateAllClaims(ctx, func(_ int64, key []byte, claim types.Claim) (stop bool) {
		// check for the presence of the addr bytes in the key.
		// first prefix byte is 0x00; so cater for that! Then + 1 to skip the separator.
		idx := bytes.Index(key[1:], []byte{0x00}) + 1 + 1
		if bytes.Equal(key[idx:idx+len(addrBytes)], addrBytes) {
			out = append(out, claim)
		}
		return false
	})
	return &types.QueryClaimsResponse{Claims: out}, nil
}

func (k Keeper) UserLastEpochClaims(c context.Context, q *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	out := []types.Claim{}

	addrBytes := []byte(q.Address)
	k.IterateAllLastEpochClaims(ctx, func(_ int64, key []byte, claim types.Claim) (stop bool) {
		// check for the presence of the addr bytes in the key.
		// First byte is 0x01 here, so no need to consider it; + 1 to skip the separator.
		idx := bytes.Index(key, []byte{0x00}) + 1
		if bytes.Equal(key[idx:idx+len(addrBytes)], addrBytes) {
			out = append(out, claim)
		}
		return false
	})

	return &types.QueryClaimsResponse{Claims: out}, nil
}
