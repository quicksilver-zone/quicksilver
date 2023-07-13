package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) error {
	return nil
}

// AfterEpochEnd is called after any registered epoch ends.
// calls:
//
//	k.AggregateDelegatorIntents
//	k.HandleQueuedUnbondings
//	k.Rebalance
//
// and re-queries icq for new zone info.
func (k *Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	// every epoch
	if epochIdentifier == epochstypes.EpochIdentifierEpoch {
		k.Logger(ctx).Info("handling epoch end", "epoch_identifier", epochIdentifier, "epoch_number", epochNumber)

		k.IterateZones(ctx, func(index int64, zone *types.Zone) (stop bool) {
			k.Logger(ctx).Info(
				"taking a snapshot of delegator intents",
				"epoch_identifier", epochIdentifier,
				"epoch_number", epochNumber,
			)
			err := k.AggregateDelegatorIntents(ctx, zone)
			if err != nil {
				// we can and need not panic here; logging the error is sufficient.
				// an error here is not expected, but also not terminal.
				// we don't return on failure here as we still want to attempt
				// the unrelated tasks below.
				k.Logger(ctx).Error(
					"encountered a problem aggregating intents; leaving aggregated intents unchanged since last epoch",
					"error", err.Error(),
					"chain_id", zone.ChainId,
					"epoch_identifier", epochIdentifier,
					"epoch_number", epochNumber,
				)
			}

			if zone.DelegationAddress == nil {
				// we have reached the end of the epoch and the delegation address is nil.
				// This shouldn't happen in normal operation, but can if the zone was registered right on the epoch boundary.
				return false
			}

			if err := k.HandleQueuedUnbondings(ctx, zone, epochNumber); err != nil {
				// we can and need not panic here; logging the error is sufficient.
				// an error here is not expected, but also not terminal.
				// we don't return on failure here as we still want to attempt
				// the unrelated tasks below.
				k.Logger(ctx).Error(
					"encountered a problem handling queued unbondings",
					"error", err.Error(),
					"chain_id", zone.ChainId,
					"epoch_identifier", epochIdentifier,
					"epoch_number", epochNumber,
				)
			}

			err = k.Rebalance(ctx, zone, epochNumber)
			if err != nil {
				// we can and need not panic here; logging the error is sufficient.
				// an error here is not expected, but also not terminal.
				// we don't return on failure here as we still want to attempt
				// the unrelated tasks below.
				k.Logger(ctx).Error(
					"encountered a problem rebalancing",
					"error", err.Error(),
					"chain_id", zone.ChainId,
					"epoch_identifier", epochIdentifier,
					"epoch_number", epochNumber,
				)
			}

			if zone.WithdrawalWaitgroup > 0 {
				k.Logger(ctx).Error(
					"epoch waitgroup was unexpected > 0; this means we did not process the previous epoch!",
					"chain_id", zone.ChainId,
					"epoch_identifier", epochIdentifier,
					"epoch_number", epochNumber,
				)
				zone.WithdrawalWaitgroup = 0
			}

			// OnChanOpenAck calls SetWithdrawalAddress (see ibc_module.go)
			k.Logger(ctx).Info(
				"withdrawing rewards",
				"chain_id", zone.ChainId,
				"epoch_identifier", epochIdentifier,
				"epoch_number", epochNumber,
			)

			vals := k.GetValidators(ctx, zone.ChainId)
			delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(vals))}}
			bz := k.cdc.MustMarshal(&delegationQuery)

			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				"cosmos.staking.v1beta1.Query/DelegatorDelegations",
				bz,
				sdk.NewInt(-1),
				types.ModuleName,
				"delegations",
				0,
			)

			addressBytes, err := addressutils.AccAddressFromBech32(zone.DelegationAddress.Address, zone.AccountPrefix)
			if err != nil {
				k.Logger(ctx).Error("cannot decode bech32 delegation addr")
				return false
			}
			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				types.BankStoreKey,
				append(banktypes.CreateAccountBalancesPrefix(addressBytes), []byte(zone.BaseDenom)...),
				sdk.NewInt(-1),
				types.ModuleName,
				"delegationaccountbalance",
				0,
			)
			// increment waitgroup; decremented in delegationaccountbalance callback
			zone.WithdrawalWaitgroup++

			rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: zone.DelegationAddress.Address}
			bz = k.cdc.MustMarshal(&rewardsQuery)

			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
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
			zone.WithdrawalWaitgroup++
			k.Logger(ctx).Info("Incrementing waitgroup for delegation",
				"value", zone.WithdrawalWaitgroup,
				"chain_id", zone.ChainId,
				"epoch_identifier", epochIdentifier,
				"epoch_number", epochNumber,
			)
			k.SetZone(ctx, zone)

			return false
		})
	}
	return nil
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for interchainstaking keeper.
type Hooks struct {
	k *Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
