package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// BeginBlocker of interchainstaking module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	_ = k.Logger(ctx)
	// every N blocks, emit QueryAccountBalances event.
	k.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
		ctx.Logger().Info("ZoneInfo: %s; DepositAccount: %s", zoneInfo.Identifier, zoneInfo.DepositAddress)
		return false
	})
}
