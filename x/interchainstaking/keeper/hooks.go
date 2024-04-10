package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (*Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) error {
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
	if epochIdentifier != epochstypes.EpochIdentifierEpoch {
		k.IterateZones(ctx, func(index int64, zone *types.Zone) (stop bool) {
			vals := k.GetValidatorAddresses(ctx, zone.ChainId)
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
			return false
		})
		return nil
	}
	k.Logger(ctx).Info("handling epoch end", "epoch_identifier", epochIdentifier, "epoch_number", epochNumber)

	epochInfo := k.EpochsKeeper.GetEpochInfo(ctx, epochIdentifier)
	k.IterateZones(ctx, func(index int64, zone *types.Zone) (stop bool) {
		k.IterateZoneRedelegationRecords(ctx, zone.ChainId, func(index int64, key []byte, record types.RedelegationRecord) (stop bool) {
			unbondingPeriod := time.Duration(zone.UnbondingPeriod / 1_000_000_000)
			redelegationDuration := time.Duration(epochInfo.CurrentEpoch-record.EpochNumber) * epochInfo.Duration

			if redelegationDuration >= unbondingPeriod {
				k.DeleteRedelegationRecord(ctx, record.ChainId, record.Source, record.Destination, record.EpochNumber)
			}

			return false
		})

		if err := k.HandleMaturedUnbondings(ctx, zone); err != nil {
			k.Logger(ctx).Error("error in HandleMaturedUnbondings", "error", err.Error())
		}

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

		k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusUnbond, func(idx int64, record types.WithdrawalRecord) bool {
			if (record.Status == types.WithdrawStatusUnbond) && !record.Acknowledged && record.EpochNumber < epochNumber {
				record.Requeued = true
				k.UpdateWithdrawalRecordStatus(ctx, &record, types.WithdrawStatusQueued)
			}
			return false
		})

		if zone.GetWithdrawalWaitgroup() > 0 {
			zone.SetWithdrawalWaitgroup(k.Logger(ctx), 0, "epoch waitgroup was unexpected > 0")
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

		delegationParams := DelegatorDelegationsParams{
			ChainID:      zone.ChainId,
			ConnectionID: zone.ConnectionId,
			Request:      bz,
		}

		delegationPayload, err := json.Marshal(&delegationParams)
		if err != nil {
			panic(err)
		}
		delegationCondition1, err := emtypes.NewConditionAll(ctx,
			emtypes.NewFieldValues(
				emtypes.FieldEqual(emtypes.FieldChainID, zone.ChainId),
				emtypes.FieldEqual(emtypes.FieldModule, types.ModuleName),
				emtypes.FieldEqual(emtypes.FieldEventType, fmt.Sprintf("%d", emtypes.EventTypeICAUnbond)),
				emtypes.FieldBegins(emtypes.FieldIdentifier, types.EpochWithdrawalMemo(epochNumber)),
			),
			true,
		)
		if err != nil {
			panic(err)
		}
		delegationCondition2, err := emtypes.NewConditionAll(ctx,
			emtypes.NewFieldValues(
				emtypes.FieldEqual(emtypes.FieldChainID, zone.ChainId),
				emtypes.FieldEqual(emtypes.FieldModule, types.ModuleName),
				emtypes.FieldEqual(emtypes.FieldEventType, fmt.Sprintf("%d", emtypes.EventTypeICADelegate)),
				emtypes.FieldBegins(emtypes.FieldIdentifier, fmt.Sprintf("batch/%d", epochNumber)),
			),
			true,
		)
		if err != nil {
			panic(err)
		}

		delegationCondition, err := emtypes.NewConditionAnd(
			ctx,
			delegationCondition1,
			delegationCondition2,
		)
		if err != nil {
			panic(err)
		}

		k.EventManagerKeeper.AddEvent(
			ctx,
			types.ModuleName,
			zone.ChainId,
			"query_delegations_epoch",
			ICQEmitDelegatorDelegations,
			emtypes.EventTypeICQQueryDelegations,
			emtypes.EventStatusPending,
			delegationCondition,
			delegationPayload,
		)

		balancesQuery := banktypes.QueryAllBalancesRequest{Address: zone.DelegationAddress.Address}
		bz = k.cdc.MustMarshal(&balancesQuery)
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			"cosmos.bank.v1beta1.Query/AllBalances",
			bz,
			sdk.NewInt(-1),
			types.ModuleName,
			"delegationaccountbalances",
			0,
		)

		k.EventManagerKeeper.AddEvent(
			ctx,
			types.ModuleName,
			zone.ChainId,
			"query_delegationaccountbalances_epoch",
			"",
			emtypes.EventTypeICQAccountBalances,
			emtypes.EventStatusActive,
			nil,
			nil,
		)

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

		k.EventManagerKeeper.AddEvent(
			ctx,
			types.ModuleName,
			zone.ChainId,
			"query_rewards_epoch",
			"",
			emtypes.EventTypeICQQueryRewards,
			emtypes.EventStatusActive,
			nil,
			nil,
		)

		rrCondition1, err := emtypes.NewConditionAll(ctx,
			emtypes.NewFieldValues(
				emtypes.FieldEqual(emtypes.FieldChainID, zone.ChainId),
				emtypes.FieldEqual(emtypes.FieldModule, types.ModuleName),
				emtypes.FieldEqual(emtypes.FieldEventType, fmt.Sprintf("%d", emtypes.EventTypeICAWithdrawRewards)),
			),
			true,
		)
		if err != nil {
			panic(err)
		}
		rrCondition2, err := emtypes.NewConditionAll(ctx,
			emtypes.NewFieldValues(
				emtypes.FieldEqual(emtypes.FieldChainID, zone.ChainId),
				emtypes.FieldEqual(emtypes.FieldModule, types.ModuleName),
				emtypes.FieldEqual(emtypes.FieldEventType, fmt.Sprintf("%d", emtypes.EventTypeICQQueryDelegations)),
			),
			true,
		)
		if err != nil {
			panic(err)
		}
		rrCondition, err := emtypes.NewConditionAnd(ctx, rrCondition1, rrCondition2)
		if err != nil {
			panic(err)
		}
		rrConditionAnd, err := emtypes.NewConditionAnd(ctx, rrCondition, delegationCondition)
		if err != nil {
			panic(err)
		}

		k.EventManagerKeeper.AddEvent(
			ctx,
			types.ModuleName,
			zone.ChainId,
			"trigger_rr",
			TriggerCalculateRedemptionRate,
			emtypes.EventTypeICQQueryRewards,
			emtypes.EventStatusPending,
			rrConditionAnd,
			[]byte(zone.ChainId),
		)

		return false
	})

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
