package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"

	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k *Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) error {
	return nil
}

func (k *Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) error {
	if epochIdentifier == epochstypes.EpochIdentifierEpoch {
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

			k.icsKeeper.ICQKeeper.MakeRequest(
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
		for _, sub := range k.prSubmodules {
			sub.Hooks(ctx, k)
		}

		tvs, err := k.CalcTokenValues(ctx)
		if err != nil {
			k.Logger(ctx).Error("unable to calculate token values", "error", err.Error())
			return nil
		}

		if allocation == nil {
			// if allocation is unset, then return early to avoid panic
			k.Logger(ctx).Error("nil allocation", "error", err.Error())
			return nil
		}

		if err := k.AllocateZoneRewards(ctx, tvs, *allocation); err != nil {
			k.Logger(ctx).Error(err.Error())
			return err
		}

		if !allocation.Lockup.IsZero() {
			// at genesis lockup will be disabled, and enabled when ICS is used.
			if err := k.AllocateLockupRewards(ctx, allocation.Lockup); err != nil {
				k.Logger(ctx).Error(err.Error())
				return err
			}
		}
	}
	return nil
}

func (k *Keeper) AfterZoneCreated(ctx sdk.Context, connectionID, chainID, accountPrefix string) error {
	connectionPd := types.ConnectionProtocolData{
		ConnectionID: connectionID,
		ChainID:      chainID,
		LastEpoch:    0,
		Prefix:       accountPrefix,
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

func (h Hooks) AfterZoneCreated(ctx sdk.Context, connectionID, chainID, accountPrefix string) error {
	return h.k.AfterZoneCreated(ctx, connectionID, chainID, accountPrefix)
}
