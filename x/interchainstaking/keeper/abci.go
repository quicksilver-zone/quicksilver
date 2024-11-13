package keeper

import (
	"bytes"
	"fmt"
	"math"
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

const blockInterval = 30

type zoneItrFn func(index int64, zone *types.Zone) (stop bool)

// BeginBlocker of interchainstaking module.
func (k *Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if ctx.BlockHeight()%blockInterval == 0 {
		if err := k.GCCompletedRedelegations(ctx); err != nil {
			k.Logger(ctx).Error("error in GCCompletedRedelegations", "error", err)
		}

		//k.HandleMaturedUnbondings(ctx)
	}
	k.IterateZones(ctx, func(index int64, zone *types.Zone) (stop bool) {
		if ctx.BlockHeight()%30 == 0 {
			// for the tasks below, we cannot panic in begin blocker; as this will crash the chain.
			// and as failing here is not terminal panicking is not necessary, but we should log
			// as an error. we don't return on failure here as we still want to attempt the unrelated
			// tasks below.
			// commenting this out until we can revisit. in its current state it causes more issues than it fixes.

			if err := k.EnsureWithdrawalAddresses(ctx, zone); err != nil {
				k.Logger(ctx).Error("error in EnsureWithdrawalAddresses", "error", err.Error())
			}
			if err := k.HandleMaturedWithdrawals(ctx, zone); err != nil {
				k.Logger(ctx).Error("error in HandleMaturedWithdrawals", "error", err.Error())
			}
			if err := k.GCCompletedUnbondings(ctx, zone); err != nil {
				k.Logger(ctx).Error("error in GCCompletedUnbondings", "error", err.Error())
			}

			addressBytes, err := addressutils.AccAddressFromBech32(zone.DelegationAddress.Address, zone.AccountPrefix)
			if err != nil {
				k.Logger(ctx).Error("cannot decode bech32 delegation addr", "error", err.Error())
			}
			zone.DelegationAddress.IncrementBalanceWaitgroup()
			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				types.BankStoreKey,
				append(banktypes.CreateAccountBalancesPrefix(addressBytes), []byte(zone.BaseDenom)...),
				sdk.NewInt(-1),
				types.ModuleName,
				"accountbalance",
				0,
			)
		}

		connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, zone.ConnectionId)
		if !found {
			return false
		}

		consState, found := k.IBCKeeper.ClientKeeper.GetLatestClientConsensusState(ctx, connection.GetClientID())
		if !found {
			return false
		}

		tmConsState, ok := consState.(*tmtypes.ConsensusState)
		if !ok {
			return false
		}

		changedValSet := len(zone.IbcNextValidatorsHash) == 0 || !bytes.Equal(zone.IbcNextValidatorsHash, tmConsState.NextValidatorsHash.Bytes())
		if !changedValSet {
			return false
		}

		k.Logger(ctx).Info("IBC ValSet has changed; requerying valset")
		// trigger valset update.
		param := k.GetParam(ctx, types.KeyValidatorSetInterval)
		if param > math.MaxInt64 {
			k.Logger(ctx).Error("parameter value exceeds int64 range", "param", param)
			panic(fmt.Errorf("parameter value exceeds int64 range: %d", param))
		}
		period := int64(param)
		query := stakingtypes.QueryValidatorsRequest{}
		err := k.EmitValSetQuery(ctx, zone.ConnectionId, zone.ChainId, query, sdkmath.NewInt(period))
		if err != nil {
			k.Logger(ctx).Error("unable to trigger valset update query", "error", err.Error())
			// failing to emit the valset update is not terminal but constitutes
			// an error, as if this starts happening frequent it is something
			// we should investigate.
		}

		zone.IbcNextValidatorsHash = tmConsState.NextValidatorsHash.Bytes()
		k.SetZone(ctx, zone)
		return false
	})
}
