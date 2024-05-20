package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"

	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
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
		return false
	})

	k.Logger(ctx).Info("setting self connection data...")
	err := k.UpdateSelfConnectionData(ctx)
	if err != nil {
		panic(err)
	}

	k.Logger(ctx).Info("distribute participation rewards...")

	allocation, err := types.GetRewardsAllocations(
		k.GetModuleBalance(ctx),
		k.GetParams(ctx).DistributionProportions,
	)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	k.Logger(ctx).Info("Triggering submodule hooks")
	for _, sub := range k.PrSubmodules {
		sub.Hooks(ctx, k)
	}

	// ensure we archive claims before we return!
	k.icsKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		k.ClaimsManagerKeeper.ArchiveAndGarbageCollectClaims(ctx, zone.ChainId)
		return false
	})

	tvs, err := k.CalcTokenValues(ctx)
	if err != nil {
		k.Logger(ctx).Error("unable to calculate token values", "error", err.Error())
		return nil
	}

	if allocation == nil {
		// if allocation is unset, then return early to avoid panic
		k.Logger(ctx).Error("nil allocation")
		return nil
	}

	if err := k.AllocateZoneRewards(ctx, tvs, *allocation); err != nil {
		k.Logger(ctx).Error("unable to allocate: tvl is zero", "error", err.Error())
		return nil
	}

	// TODO: remove 'lockup' allocation logic.
	// if !allocation.Lockup.IsZero() {
	// 	// at genesis lockup will be disabled, and enabled when ICS is used.
	// 	if err := k.AllocateLockupRewards(ctx, allocation.Lockup); err != nil {
	// 		k.Logger(ctx).Error(err.Error())
	// 		return err
	// 	}
	// }
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

	// channel for the host chain
	channel, found := k.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", zone.TransferChannel)
	if !found {
		return fmt.Errorf("channel not found: %s", zone.TransferChannel)
	}

	// Create LiquidAllowedDenomProtocolData for the host zone
	hostZoneDenom := types.LiquidAllowedDenomProtocolData{
		ChainID:               zone.ChainId,
		RegisteredZoneChainID: zone.ChainId,
		IbcDenom:              utils.DeriveIbcDenom(transfertypes.PortID, channel.Counterparty.ChannelId, transfertypes.PortID, zone.TransferChannel, zone.LocalDenom),
		QAssetDenom:           zone.LocalDenom,
	}
	if err := hostZoneDenom.ValidateBasic(); err != nil {
		return err
	}
	hostZoneDenomBytes, err := json.Marshal(hostZoneDenom)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, hostZoneDenom.GenerateKey(), &types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)],
		Data: hostZoneDenomBytes,
	})

	// Fetch OsmosisParamsProtocolData and create LiquidAllowedDenomProtocolData for Osmosis
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeOsmosisParams), func(index int64, _ []byte, data types.ProtocolData) (stop bool) {
		osmosisParams, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisParams, data.Data)
		if err != nil {
			return false
		}
		osmosisParamsData, ok := osmosisParams.(*types.OsmosisParamsProtocolData)
		if !ok {
			return false
		}
		osmosisConnectionData, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, osmosisParamsData.ChainID)
		if !found {
			return false
		}
		connectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, osmosisConnectionData.Data)
		if err != nil {
			return false
		}
		osmosisChannel := connectionData.(*types.ConnectionProtocolData).TransferChannel

		// channel for the osmosis chain
		channel, found := k.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", osmosisChannel)
		if !found {
			return false
		}
		osmosisDenom := types.LiquidAllowedDenomProtocolData{
			ChainID:               osmosisParamsData.ChainID,
			RegisteredZoneChainID: zone.ChainId,
			IbcDenom:              utils.DeriveIbcDenom(transfertypes.PortID, channel.Counterparty.ChannelId, transfertypes.PortID, osmosisChannel, zone.LocalDenom),
			QAssetDenom:           zone.LocalDenom,
		}
		if err := osmosisDenom.ValidateBasic(); err != nil {
			return false
		}
		osmosisDenomBytes, err := json.Marshal(osmosisDenom)
		if err != nil {
			return false
		}
		k.SetProtocolData(ctx, osmosisDenom.GenerateKey(), &types.ProtocolData{
			Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)],
			Data: osmosisDenomBytes,
		})
		return true
	})

	// Fetch UmeeParamsProtocolData and create LiquidAllowedDenomProtocolData for Umee
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeParams), func(index int64, _ []byte, data types.ProtocolData) (stop bool) {
		umeeParams, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeParams, data.Data)
		if err != nil {
			return false
		}
		umeeParamsData, ok := umeeParams.(*types.UmeeParamsProtocolData)
		if !ok {
			return false
		}
		umeeConnectionData, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, umeeParamsData.ChainID)
		if !found {
			return false
		}
		connectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, umeeConnectionData.Data)
		if err != nil {
			return false
		}

		umeeChannel := connectionData.(*types.ConnectionProtocolData).TransferChannel
		// channel for the umee chain
		channel, found := k.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", umeeChannel)
		if !found {
			return false
		}
		umeeDenom := types.LiquidAllowedDenomProtocolData{
			ChainID:               umeeParamsData.ChainID,
			RegisteredZoneChainID: zone.ChainId,
			IbcDenom:              utils.DeriveIbcDenom(transfertypes.PortID, channel.Counterparty.ChannelId, transfertypes.PortID, umeeChannel, zone.LocalDenom),
			QAssetDenom:           zone.LocalDenom,
		}
		if err := umeeDenom.ValidateBasic(); err != nil {
			return false
		}
		umeeDenomBytes, err := json.Marshal(umeeDenom)
		if err != nil {
			return false
		}
		k.SetProtocolData(ctx, umeeDenom.GenerateKey(), &types.ProtocolData{
			Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)],
			Data: umeeDenomBytes,
		})
		return true
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
