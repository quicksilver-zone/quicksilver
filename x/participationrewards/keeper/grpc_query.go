package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

var _ types.QueryServer = Keeper{}

// Params returns params of the participationrewards module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) ProtocolData(c context.Context, q *types.QueryProtocolDataRequest) (*types.QueryProtocolDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	out := []json.RawMessage{}
	k.IterateProtocolDatas(ctx, q.Protocol, func(index int64, data types.ProtocolData) (stop bool) {
		out = append(out, data.Data)
		return false
	})

	return &types.QueryProtocolDataResponse{Data: out}, nil
}
