package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {
	// delete expired events, and allow triggers to be executed.
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	k.IteratePrefixedEvents(ctx, nil, func(_ int64, eventInfo types.Event) (stop bool) {
		if eventInfo.ExpiryTime.Before(ctx.BlockTime()) {
			k.Logger(ctx).Info("deleting expired event", "module", eventInfo.Module, "id", eventInfo.Identifier, "chain", eventInfo.ChainId)
			k.DeleteEvent(ctx, eventInfo.Module, eventInfo.ChainId, eventInfo.Identifier)
		}
		return false
	})

	k.TriggerAll(ctx)
}
