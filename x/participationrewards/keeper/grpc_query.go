package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

var _ types.QueryServer = &Keeper{}

// Params returns params of the participationrewards module.
func (k *Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k *Keeper) ProtocolData(c context.Context, q *types.QueryProtocolDataRequest) (*types.QueryProtocolDataResponse, error) {
	return nil, nil
}
