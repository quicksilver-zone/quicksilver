package keeper

import (
	"encoding/hex"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
)

const DefaultTTL = 1000

// EndBlocker of interchainquery module.
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	_ = k.Logger(ctx)
	events := sdk.Events{}
	// emit events for periodic queries
	k.IterateQueries(ctx, func(_ int64, queryInfo types.Query) (stop bool) {
		// if !periodic, and emission was too long ago, delete
		if queryInfo.Period.IsNegative() && // not periodic
			!queryInfo.LastEmission.IsNil() && // has been emitted
			!queryInfo.LastEmission.IsZero() && // has been emitted
			queryInfo.LastEmission.LT(math.NewInt(ctx.BlockHeight()-DefaultTTL)) { // TODO: add per query TTL
			k.DeleteQuery(ctx, queryInfo.Id)
			return false
		}

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
}
