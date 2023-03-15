package keeper

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type zoneItrFn func(index int64, zoneInfo types.Zone) (stop bool)

// BeginBlocker of interchainstaking module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	// post upgrade-v1.2.5 processing
	if ctx.BlockHeight() == 14540740 {
		zone, found := k.GetZone(ctx, "regen-1")
		if found {
			k.IterateReceipts(ctx, func(_ int64, receiptInfo types.Receipt) (stop bool) {
				if receiptInfo.ChainId == "regen-1" && receiptInfo.Completed == nil {
					sendMsg := banktypes.MsgSend{
						FromAddress: "",
						ToAddress:   "",
						Amount:      receiptInfo.Amount,
					}
					err := k.handleSendToDelegate(ctx, &zone, &sendMsg, receiptInfo.Txhash)
					if err != nil {
						k.Logger(ctx).Error("error in processing Pending delegations for regen-1 ", "error", err)
					}

				}
				return false
			})
		}

	}
	if ctx.BlockHeight()%30 == 0 {
		if err := k.GCCompletedRedelegations(ctx); err != nil {
			k.Logger(ctx).Error("error in GCCompletedRedelegations", "error", err)
		}
	}
	k.IterateZones(ctx, func(index int64, zone types.Zone) (stop bool) {
		if ctx.BlockHeight()%30 == 0 {
			// for the tasks below, we cannot panic in begin blocker; as this will crash the chain.
			// and as failing here is not terminal we panicking is not necessary, but we should log
			// as an error. we don't return on failure here as we still want to attempt the unrelated
			// tasks below.
			// commenting this out until we can revisit. in it's current state it causes more issues than it fixes.
			// if err := k.EnsureICAsActive(ctx, &zone); err != nil {
			// 	k.Logger(ctx).Error("error in EnsureICAsActive", "error", err)
			// }

			if err := k.EnsureWithdrawalAddresses(ctx, &zone); err != nil {
				k.Logger(ctx).Error("error in EnsureWithdrawalAddresses", "error", err)
			}
			if err := k.HandleMaturedUnbondings(ctx, &zone); err != nil {
				k.Logger(ctx).Error("error in HandleMaturedUnbondings", "error", err)
			}
			if err := k.GCCompletedUnbondings(ctx, &zone); err != nil {
				k.Logger(ctx).Error("error in GCCompletedUnbondings", "error", err)
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
							k.Logger(ctx).Error("unable to trigger valset update query", "error", err)
							// failing to emit the valset update is not terminal but constitutes
							// an error, as if this starts happening frequent it is something
							// we should investigate.
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
