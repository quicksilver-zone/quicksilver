package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var _ types.QueryServer = Keeper{}

// RegisteredZoneInfos provide running epochInfos
func (k Keeper) RegisteredZoneInfos(c context.Context, req *types.QueryRegisteredZonesInfoRequest) (*types.QueryRegisteredZonesInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var epochs []types.RegisteredZone
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var epoch types.RegisteredZone
		if err := k.cdc.Unmarshal(value, &epoch); err != nil {
			return err
		}
		epochs = append(epochs, epoch)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRegisteredZonesInfoResponse{
		Zones:      epochs,
		Pagination: pageRes,
	}, nil
}
