package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type zoneItrFn func(index int64, zoneInfo types.RegisteredZone) (stop bool)

// BeginBlocker of interchainstaking module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if ctx.BlockHeight()%int64(k.GetParam(ctx, types.KeyDepositInterval)) == 0 {
		k.IterateRegisteredZones(ctx, k.depositInterval(ctx))
	}

	if ctx.BlockHeight()%int64(k.GetParam(ctx, types.KeyDelegationsInterval)) == 0 {
		k.IterateRegisteredZones(ctx, k.delegateDelegationsInterval(ctx))
	}
}
