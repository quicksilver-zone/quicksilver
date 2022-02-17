package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: &keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SubmitQueryResponse(goCtx context.Context, msg *types.MsgSubmitQueryResponse) (*types.MsgSubmitQueryResponseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	q, found := k.GetPeriodicQuery(ctx, msg.QueryId)
	if found {
		q.LastHeight = sdk.NewInt(ctx.BlockHeight())
		k.SetPeriodicQuery(ctx, q)
		k.SetDatapointForId(ctx, msg.QueryId, msg.Result, sdk.NewInt(msg.Height))

	} else {
		_, found2 := k.GetSingleQuery(ctx, msg.QueryId)
		if found2 {
			k.DeleteSingleQuery(ctx, msg.QueryId)
			k.SetDatapointForId(ctx, msg.QueryId, msg.Result, sdk.NewInt(msg.Height))
		} else {
			return nil, fmt.Errorf("query object no longer exists; likely deleted since query was requested")
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgSubmitQueryResponseResponse{}, nil
}
