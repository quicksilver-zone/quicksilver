package keeper

import (
	"context"
	"fmt"
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

	if !found {
		k.Logger(ctx).Debug("query not found", "QueryID", msg.QueryId)

		return &types.MsgSubmitQueryResponseResponse{}, nil
	}

	// check if query was previously processed
	// - indicated by query.LastHeight matching current Block Height;
	if q.LastHeight.Int64() == ctx.BlockHeader().Height {
		k.Logger(ctx).Debug("ignoring duplicate query", "id", q.Id, "type", q.QueryType)
		// technically this is an error, but will cause the entire tx to fail
		// if we have one 'bad' message, so we can just no-op here.
		return &types.MsgSubmitQueryResponseResponse{}, nil
	}

	pathParts := strings.Split(q.QueryType, "/")
	if pathParts[len(pathParts)-1] == "key" {
		if err := utils.ValidateProofOps(ctx, k.IBCKeeper, q.ConnectionId, q.ChainId, msg.Height, pathParts[1], q.Request, msg.Result, msg.ProofOps); err != nil {
			k.Logger(ctx).Error("failed to validate proofops", "id", q.Id, "type", q.QueryType)
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
	callbackExecuted := false

	for _, key := range keys {
		module := k.callbacks[key]
		if module.Has(q.CallbackId) {
			err := module.Call(ctx, q.CallbackId, msg.Result, q)
			callbackExecuted = true
			if err != nil {
				// not edge case: proceed with regular error handling!
				if err != types.ErrSucceededNoDelete {
					k.Logger(ctx).Error("error in callback", "error", err, "msg", msg.QueryId, "result", msg.Result, "type", q.QueryType, "params", q.Request)
					return nil, err
				}
				// edge case: the callback has resent the same query (re-query)!
				// action:    set noDelete to true and continue (short circuit error handling)!
				noDelete = true
			}
			// we have executed a callback; only a single callback is expected per request, so break here.
			break
		}
	}

	if !callbackExecuted && q.CallbackId != "" {
		k.Logger(ctx).Error("callback expected but not found", "callbackId", q.CallbackId, "msg", msg.QueryId, "type", q.QueryType)
		return nil, fmt.Errorf("expected callback %s, but did not find it", q.CallbackId)
	}

	if q.Ttl > 0 {
		// don't store if ttl is 0
		if err := k.SetDatapointForID(ctx, msg.QueryId, msg.Result, sdk.NewInt(msg.Height)); err != nil {
			k.Logger(ctx).Error("failed to set datapoint", "id", q.Id, "type", q.QueryType)
			return nil, err
		}
	}

	// check for and delete non-repeating queries, update any other
	// - Period.IsNegative() indicates a single query;
	// - noDelete indicates a response that triggered a re-query;
	if q.Period.IsNegative() {
		if !noDelete {
			k.DeleteQuery(ctx, msg.QueryId)
		}
	} else {
		q.LastHeight = sdk.NewInt(ctx.BlockHeight())
		k.SetQuery(ctx, q)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgSubmitQueryResponseResponse{}, nil
}
