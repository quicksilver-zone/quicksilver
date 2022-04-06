package keeper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type zoneItrFn func(index int64, zoneInfo types.RegisteredZone) (stop bool)

// BeginBlocker of interchainstaking module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if ctx.BlockHeight()%types.ValidatorSetInterval == 0 {
		k.IterateRegisteredZones(ctx, k.validatorSetInterval(ctx))
	}

	// every N blocks, emit QueryAccountBalances event.
	if ctx.BlockHeight()%types.DepositInterval == 0 {
		k.IterateRegisteredZones(ctx, k.depositInterval(ctx))
	}

	if ctx.BlockHeight()%types.DelegateInterval == 0 {
		k.IterateRegisteredZones(ctx, k.delegateInterval(ctx))
	}

	if ctx.BlockHeight()%types.DelegateDelegationsInterval == 0 {
		k.IterateRegisteredZones(ctx, k.delegateDelegationsInterval(ctx))
	}
}
