package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/types"
)

var _ types.QueryServer = &Keeper{}

// Params returns params of the participationrewards module.
func (k *Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k *Keeper) ProtocolData(c context.Context, q *types.QueryProtocolDataRequest) (*types.QueryProtocolDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	out := []json.RawMessage{}

	pdType, exists := types.ProtocolDataType_value[q.Type]
	if !exists {
		return nil, types.ErrUnknownProtocolDataType
	}

	prefix := append(types.GetPrefixProtocolDataKey(types.ProtocolDataType(pdType)), []byte(q.Key)...)
	k.IteratePrefixedProtocolDatas(ctx, prefix, func(index int64, _ []byte, data types.ProtocolData) (stop bool) {
		out = append(out, data.Data)
		return false
	})

	return &types.QueryProtocolDataResponse{Data: out}, nil
}
