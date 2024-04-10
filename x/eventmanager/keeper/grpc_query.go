package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

var _ types.QueryServer = Keeper{}

// Events returns information about registered zones.
func (k Keeper) Events(c context.Context, req *types.QueryEventsRequest) (*types.QueryEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var events []types.Event
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)

	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(_, value []byte, accumulate bool) (bool, error) {
		var event types.Event
		if err := k.cdc.Unmarshal(value, &event); err != nil {
			return false, err
		}

		if event.ChainId == req.ChainId {
			events = append(events, event)
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryEventsResponse{
		Events:     events,
		Pagination: pageRes,
	}, nil
}
