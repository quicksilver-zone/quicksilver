package keeper

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type zoneItrFn func(index int64, zoneInfo types.Zone) (stop bool)

// BeginBlocker of interchainstaking module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	k.IterateZones(ctx, func(index int64, zone types.Zone) (stop bool) {
		if ctx.BlockHeight()%10 == 0 {
			if err := k.EnsureWithdrawalAddresses(ctx, &zone); err != nil {
				k.Logger(ctx).Error(err.Error())
				panic(err)
			}
		}
		connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, zone.ConnectionId)
		if found {
			consState, found := k.IBCKeeper.ClientKeeper.GetLatestClientConsensusState(ctx, connection.GetClientID())
			if found {
				tmConsState, ok := consState.(*tmtypes.ConsensusState)
				if ok {
					if len(zone.IbcNextValidatorsHash) == 0 || !bytes.Equal(zone.IbcNextValidatorsHash, tmConsState.NextValidatorsHash.Bytes()) {
						k.Logger(ctx).Info("IBC ValSet has changed; requerying valset")
						// trigger valset update.
						err := k.EmitValsetRequery(ctx, zone.ConnectionId, zone.ChainId)
						if err != nil {
							k.Logger(ctx).Error("unable to trigger valset update query")
						}
						zone.IbcNextValidatorsHash = tmConsState.NextValidatorsHash.Bytes()
						k.SetZone(ctx, &zone)
					}
				}
			}
		}
		return false
	})
}
