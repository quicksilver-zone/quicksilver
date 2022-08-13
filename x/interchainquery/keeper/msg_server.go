package keeper

import (
	"context"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/utils"
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
	q, found := k.GetQuery(ctx, msg.QueryId)
	// if found && q.LastHeight.Int64() != ctx.BlockHeader().Height {
	if found {
		pathParts := strings.Split(q.QueryType, "/")
		if pathParts[len(pathParts)-1] == "key" {
			if err := utils.ValidateProofOps(ctx, k.IBCKeeper, q.ConnectionId, q.ChainId, msg.Height, pathParts[1], q.Request, msg.Result, msg.ProofOps); err != nil {
				return nil, err
			}
		}

		noDelete := false
		// execute registered callbacks.

		keys := []string{}
		for k := range k.callbacks {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, key := range keys {
			module := k.callbacks[key]
			if module.Has(q.CallbackId) {
				err := module.Call(ctx, q.CallbackId, msg.Result, q)
				if err != nil {
					// handle edge case; callback has resent the same query!
					// set noDelete to true and short circuit error handling!
					if err == types.ErrSucceededNoDelete {
						noDelete = true
					} else {
						k.Logger(ctx).Error("error in callback", "error", err, "msg", msg.QueryId, "result", msg.Result, "type", q.QueryType, "params", q.Request)
						return nil, err
					}
				}
			}
		}

		if q.Ttl > 0 {
			// don't store if ttl is 0
			if err := k.SetDatapointForID(ctx, msg.QueryId, msg.Result, sdk.NewInt(msg.Height)); err != nil {
				return nil, err
			}
		}

		if q.Period.IsNegative() {
			if !noDelete {
				k.DeleteQuery(ctx, msg.QueryId)
			}
		} else {
			q.LastHeight = sdk.NewInt(ctx.BlockHeight())
			k.SetQuery(ctx, q)
		}

	} else {
		k.Logger(ctx).Info("Ignoring duplicate query")
		return &types.MsgSubmitQueryResponseResponse{}, nil // technically this is an error, but will cause the entire tx to fail if we have one 'bad' message, so we can just no-op here.
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgSubmitQueryResponseResponse{}, nil
}
