package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

const (
	RetryInterval = 25
)

// EndBlocker of interchainquery module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	_ = k.Logger(ctx)
	events := sdk.Events{}

	// emit events for periodic queries
	k.IteratePeriodicQueries(ctx, func(_ int64, queryInfo types.PeriodicQuery) (stop bool) {
		if queryInfo.LastHeight.Add(queryInfo.Period).Equal(sdk.NewInt(ctx.BlockHeight())) {
			k.Logger(ctx).Info("Interchainquery periodic event emitted", "id", queryInfo.Id)
			event := sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueQuery),
				sdk.NewAttribute(types.AttributeKeyQueryId, queryInfo.Id),
				sdk.NewAttribute(types.AttributeKeyChainId, queryInfo.ChainId),
				sdk.NewAttribute(types.AttributeKeyConnectionId, queryInfo.ConnectionId),
				sdk.NewAttribute(types.AttributeKeyType, queryInfo.QueryType),
			)

			for key, val := range queryInfo.GetQueryParameters() {
				event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyParams, fmt.Sprintf("%s:%s:%s", queryInfo.Id, key, val)))
			}

			events = append(events, event)
			queryInfo.LastHeight = sdk.NewInt(ctx.BlockHeight())
			k.SetPeriodicQuery(ctx, queryInfo)

		}
		return false
	})

	// emit events for single queries
	k.IterateSingleQueries(ctx, func(_ int64, queryInfo types.SingleQuery) (stop bool) {
		k.Logger(ctx).Info("Interchainquery single event found", "id", queryInfo.Id, "emitHeight", queryInfo.EmitHeight, "nextHeight", queryInfo.EmitHeight.Add(sdk.NewInt(RetryInterval)), "current", sdk.NewInt(ctx.BlockHeight()))
		if queryInfo.EmitHeight.Add(sdk.NewInt(RetryInterval)).LTE(sdk.NewInt(ctx.BlockHeight())) { // if it has been 25 blocks since we last emitted.
			k.Logger(ctx).Info("Interchainquery single event emitted", "id", queryInfo.Id)
			event := sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueQuery),
				sdk.NewAttribute(types.AttributeKeyQueryId, queryInfo.Id),
				sdk.NewAttribute(types.AttributeKeyChainId, queryInfo.ChainId),
				sdk.NewAttribute(types.AttributeKeyConnectionId, queryInfo.ConnectionId),
				sdk.NewAttribute(types.AttributeKeyType, queryInfo.QueryType),
			)

			for key, val := range queryInfo.GetQueryParameters() {
				event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyParams, fmt.Sprintf("%s:%s:%s", queryInfo.Id, key, val)))
			}

			events = append(events, event)
			queryInfo.EmitHeight = sdk.NewInt(ctx.BlockHeight())
			k.SetSingleQuery(ctx, queryInfo)
		}
		return false
	})

	if len(events) > 0 {
		ctx.EventManager().EmitEvents(events)
	}
	// garbage collection of DataPoints
}
