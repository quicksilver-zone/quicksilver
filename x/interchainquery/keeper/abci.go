package keeper

import (
	"encoding/hex"
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
	k.IterateQueries(ctx, func(_ int64, queryInfo types.Query) (stop bool) {
		// if queryInfo.QueryType == "ibc.ClientUpdate" && queryInfo.LastEmission.AddRaw(1000).LTE(sdk.NewInt(ctx.BlockHeight())) {
		// 	k.DeleteQuery(ctx, queryInfo.Id)
		// 	k.Logger(ctx).Error("Deleting stale query")
		// 	return false
		// }
		if queryInfo.LastEmission.IsNil() || queryInfo.LastEmission.IsZero() || queryInfo.LastEmission.Add(queryInfo.Period).Equal(sdk.NewInt(ctx.BlockHeight())) {
			k.Logger(ctx).Debug("Interchainquery event emitted", "id", queryInfo.Id)
			event := sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueQuery),
				sdk.NewAttribute(types.AttributeKeyQueryID, queryInfo.Id),
				sdk.NewAttribute(types.AttributeKeyChainID, queryInfo.ChainId),
				sdk.NewAttribute(types.AttributeKeyConnectionID, queryInfo.ConnectionId),
				sdk.NewAttribute(types.AttributeKeyType, queryInfo.QueryType),
				// TODO: add height to request type
				sdk.NewAttribute(types.AttributeKeyHeight, "0"),
				sdk.NewAttribute(types.AttributeKeyRequest, hex.EncodeToString(queryInfo.Request)),
			)

			events = append(events, event)
			queryInfo.LastEmission = sdk.NewInt(ctx.BlockHeight())
			k.SetQuery(ctx, queryInfo)

		}
		return false
	})

	if len(events) > 0 {
		ctx.EventManager().EmitEvents(events)
	}

	k.IterateDatapoints(ctx, func(_ int64, dp types.DataPoint) bool {
		q, found := k.GetQuery(ctx, dp.Id)
		if !found {
			// query was removed; delete datapoint
			k.DeleteDatapoint(ctx, dp.Id)
		} else if dp.LocalHeight.Int64()+int64(q.Ttl) > ctx.BlockHeader().Height {
			// gc old data
			k.DeleteDatapoint(ctx, dp.Id)
		}

		return false
	})
}
