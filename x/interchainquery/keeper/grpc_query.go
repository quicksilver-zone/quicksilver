package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

var _ types.QuerySrvrServer = Keeper{}

// Queries returns information about registered zones.
func (k Keeper) Queries(c context.Context, req *types.QueryRequestsRequest) (*types.QueryRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var queries []types.Query
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)

	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(_, value []byte, accumulate bool) (bool, error) {
		var query types.Query
		if err := k.cdc.Unmarshal(value, &query); err != nil {
			return false, err
		}

		if query.ChainId == req.ChainId && (query.LastEmission.IsNil() || query.LastEmission.IsZero() || query.LastEmission.GTE(query.LastHeight)) {
			queries = append(queries, query)
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRequestsResponse{
		Queries:    queries,
		Pagination: pageRes,
	}, nil
}
