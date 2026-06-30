package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/epochs/types"
)

// BeginBlocker of epochs module.
func (k *Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
}

func startInitialEpoch(epochInfo types.EpochInfo) types.EpochInfo {
	epochInfo.EpochCountingStarted = true
	epochInfo.CurrentEpoch = 1
	epochInfo.CurrentEpochStartTime = epochInfo.StartTime
	return epochInfo
}

func endEpoch(epochInfo types.EpochInfo) types.EpochInfo {
	epochInfo.CurrentEpoch++
	epochInfo.CurrentEpochStartTime = epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
	return epochInfo
}
