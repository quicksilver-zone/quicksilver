package keeper

import (
	"time"

	types "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


// BeginBlocker of interchainstaking module.
func (k *Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
}
