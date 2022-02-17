package keeper

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// BeginBlocker of interchainstaking module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	// every N blocks, emit QueryAccountBalances event.
	if ctx.BlockHeight()%5 == 0 {
		k.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
			ctx.Logger().Info("ZoneInfo: %s; DepositAccount: %s", zoneInfo.Identifier, zoneInfo.DepositAddress)
			balance_data, err := k.ICQKeeper.GetDatapoint(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.bank.v1beta1.Query/AllBalances", map[string]string{"address": zoneInfo.DepositAddress})
			if err != nil {
				ctx.Logger().Error("Unable to query balance for deposit account", zoneInfo.DepositAddress)
			}
			balance := sdk.Coins{}
			err = json.Unmarshal(balance_data.Value, &balance)
			if err != nil {
				ctx.Logger().Error("Unable to unmarshal balance for deposit account", zoneInfo.DepositAddress)
			}

			ctx.Logger().Info("Balance of deposit account", zoneInfo.DepositAddress, balance)

			return false
		})
	}
}
