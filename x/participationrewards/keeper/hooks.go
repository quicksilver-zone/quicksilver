package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"

	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (*Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) error {
	return nil
}

func (k *Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) error {
	if epochIdentifier != epochstypes.EpochIdentifierEpoch {
		return nil
	}

	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeConnection), func(index int64, _ []byte, data types.ProtocolData) (stop bool) {
		blockQuery := tmservice.GetLatestBlockRequest{}
		bz := k.cdc.MustMarshal(&blockQuery)

		iConnectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, data.Data)
		if err != nil {
			k.Logger(ctx).Error("Error unmarshalling protocol data")
			return false
		}
		connectionData, _ := iConnectionData.(*types.ConnectionProtocolData)
		if connectionData.ChainID == ctx.ChainID() {
			return false
		}

		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"cosmos.base.tendermint.v1beta1.Service/GetLatestBlock",
			bz,
			sdk.NewInt(-1),
			types.ModuleName,
			SetEpochBlockCallbackID,
			0,
		)

		k.EventManagerKeeper.AddEvent(ctx, types.ModuleName, connectionData.ChainID, "get_epoch_height", "", emtypes.EventTypeICQGetLatestBlock, emtypes.EventStatusActive, nil, nil)
		return false
	})

	k.icsKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		// ensure we archive claims before we return!
		k.ClaimsManagerKeeper.ArchiveAndGarbageCollectClaims(ctx, zone.ChainId)
		// send validator performance query
		k.QueryValidatorDelegationPerformance(ctx, zone)
		return false
	})

	k.Logger(ctx).Info("setting self connection data...")
	err := k.UpdateSelfConnectionData(ctx)
	if err != nil {
		panic(err)
	}

	k.Logger(ctx).Info("allocate participation rewards...")

	// determine allocations splits the balance of the module between holding/usage and validatorSelection rewards.
	err = k.DetermineAllocations(
		ctx,
		k.GetModuleBalance(ctx),
		k.GetParams(ctx).DistributionProportions,
	)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	conditionGetEpochHeight, err := emtypes.NewConditionAll(ctx, emtypes.NewFieldValues(emtypes.NewFieldValue(emtypes.FieldIdentifier, "get_epoch_height", emtypes.FIELD_OPERATOR_EQUAL, true)), false)
	if err != nil {
		panic(err)
	}

	conditionValidatorPerformance, err := emtypes.NewConditionAll(ctx, emtypes.NewFieldValues(emtypes.NewFieldValue(emtypes.FieldIdentifier, "validator_performance", emtypes.FIELD_OPERATOR_EQUAL, true)), false)
	if err != nil {
		panic(err)
	}

	conditionSubmodulePre, err := emtypes.NewConditionAnd(ctx, conditionGetEpochHeight, conditionValidatorPerformance)
	if err != nil {
		panic(err)
	}

	// add event to ensure submodule hooks are called when the validator_performance and get_epoch_height calls have returned.
	k.EventManagerKeeper.AddEvent(ctx, types.ModuleName, "", "submodules", Submodules, emtypes.EventTypeSubmodules, emtypes.EventStatusPending, conditionSubmodulePre, nil)

	conditionSubmoduleComplete, err := emtypes.NewConditionAll(ctx, emtypes.NewFieldValues(emtypes.NewFieldValue(emtypes.FieldIdentifier, "submodule", emtypes.FIELD_OPERATOR_BEGINSWITH, true)), false)
	if err != nil {
		panic(err)
	}

	conditionCalcTokensPre, err := emtypes.NewConditionAnd(ctx, conditionSubmodulePre, conditionSubmoduleComplete)
	if err != nil {
		panic(err)
	}
	// add calc_tokens event to be triggered on satisfaction of all submodule*, validator_performance, and get_epoch_height calls events.
	k.EventManagerKeeper.AddEvent(ctx, types.ModuleName, "", "calc_tokens", CalculateValues, emtypes.EventTypeCalculateTvls, emtypes.EventStatusPending, conditionCalcTokensPre, nil)

	conditionCalcTokensComplete, err := emtypes.NewConditionAll(ctx, emtypes.NewFieldValues(emtypes.NewFieldValue(emtypes.FieldIdentifier, "calc_tokens", emtypes.FIELD_OPERATOR_EQUAL, true)), false)
	if err != nil {
		panic(err)
	}
	conditionDistributeRewardsPre, err := emtypes.NewConditionAnd(ctx, conditionCalcTokensPre, conditionCalcTokensComplete)
	if err != nil {
		panic(err)
	}

	// add distribute_rewards event to trigger on completion of get_epoch_height, validator_performance, submodule* and calc_token events.
	k.EventManagerKeeper.AddEvent(ctx, types.ModuleName, "", "distribute_rewards", DistributeRewards, emtypes.EventTypeDistributeRewards, emtypes.EventStatusPending, conditionDistributeRewardsPre, nil)

	return nil
}

func (k *Keeper) AfterZoneCreated(ctx sdk.Context, zone *icstypes.Zone) error {
	connectionPd := types.ConnectionProtocolData{
		ConnectionID:    zone.ConnectionId,
		ChainID:         zone.ChainId,
		LastEpoch:       0,
		Prefix:          zone.AccountPrefix,
		TransferChannel: zone.TransferChannel,
	}

	if err := connectionPd.ValidateBasic(); err != nil {
		return err
	}

	connectionPdBytes, err := json.Marshal(connectionPd)
	if err != nil {
		return err
	}

	k.SetProtocolData(ctx, connectionPd.GenerateKey(), &types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
		Data: connectionPdBytes,
	})

	localDenom := types.LiquidAllowedDenomProtocolData{
		ChainID:               ctx.ChainID(),
		RegisteredZoneChainID: zone.ChainId,
		IbcDenom:              zone.LocalDenom,
		QAssetDenom:           zone.LocalDenom,
	}

	if err := localDenom.ValidateBasic(); err != nil {
		return err
	}

	localDenomBytes, err := json.Marshal(localDenom)
	if err != nil {
		return err
	}

	k.SetProtocolData(ctx, localDenom.GenerateKey(), &types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)],
		Data: localDenomBytes,
	})

	return nil
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper.
type Hooks struct {
	k *Keeper
}

var (
	_ epochstypes.EpochHooks = Hooks{}
	_ icstypes.IcsHooks      = Hooks{}
)

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks.

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterZoneCreated(ctx sdk.Context, zone *icstypes.Zone) error {
	return h.k.AfterZoneCreated(ctx, zone)
}
