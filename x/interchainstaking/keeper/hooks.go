package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	// every epoch
	if epochIdentifier == "epoch" {
		k.Logger(ctx).Info("handling epoch end")

		k.IterateZones(ctx, func(index int64, zoneInfo types.Zone) (stop bool) {
			k.Logger(ctx).Info("taking a snapshot of intents")
			err := k.AggregateIntents(ctx, &zoneInfo)
			if err != nil {
				// we can and need not panic here; logging the error is sufficient.
				// an error here is not expected, but also not terminal.
				// we don't return on failure here as we still want to attempt
				// the unrelated tasks below.
				k.Logger(ctx).Error("encountered a problem aggregating intents; leaving aggregated intents unchanged since last epoch", "error", err)
			}

			if zoneInfo.DelegationAddress == nil {
				// we have reached the end of the epoch and the delegation address is nil.
				// This shouldn't happen in normal operation, but can if the zone was registered right on the epoch boundary.
				return false
			}

			if zoneInfo.WithdrawalWaitgroup > 0 {
				k.Logger(ctx).Error("epoch waitgroup was unexpected > 0; this means we did not process the previous epoch!")
				zoneInfo.WithdrawalWaitgroup = 0
			}

			// OnChanOpenAck calls SetWithdrawalAddress (see ibc_module.go)
			k.Logger(ctx).Info("Withdrawing rewards")

			delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zoneInfo.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zoneInfo.Validators))}}
			bz := k.cdc.MustMarshal(&delegationQuery)

			k.ICQKeeper.MakeRequest(
				ctx,
				zoneInfo.ConnectionId,
				zoneInfo.ChainId,
				"cosmos.staking.v1beta1.Query/DelegatorDelegations",
				bz,
				sdk.NewInt(-1),
				types.ModuleName,
				"delegations",
				0,
			)

			rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: zoneInfo.DelegationAddress.Address}
			bz = k.cdc.MustMarshal(&rewardsQuery)

			k.ICQKeeper.MakeRequest(
				ctx,
				zoneInfo.ConnectionId,
				zoneInfo.ChainId,
				"cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
				bz,
				sdk.NewInt(-1),
				types.ModuleName,
				"rewards",
				0,
			)

			// increment the WithdrawalWaitgroup
			// this allows us to track the response for every protocol delegator
			// WithdrawalWaitgroup is decremented in RewardsCallback
			zoneInfo.WithdrawalWaitgroup++
			k.Logger(ctx).Info("Incrementing waitgroup for delegation", "value", zoneInfo.WithdrawalWaitgroup)
			k.SetZone(ctx, &zoneInfo)

			return false
		})
	}
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
