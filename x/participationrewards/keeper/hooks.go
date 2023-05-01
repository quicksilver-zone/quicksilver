package keeper

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"

	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k *Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) error {
	return nil
}

func (k *Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) error {
	if epochIdentifier == epochstypes.EpochIdentifierEpoch {
		k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeConnection), func(index int64, data types.ProtocolData) (stop bool) {
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
				"epochblock",
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
			if errors.Is(err, types.ErrNothingToAllocate) {
				k.Logger(ctx).Info(err.Error())
			} else {
				k.Logger(ctx).Error(err.Error())
			}
		}

		k.Logger(ctx).Info("Triggering submodule hooks")
		for _, sub := range k.prSubmodules {
			sub.Hooks(ctx, k)
		}

		tvs, err := k.calcTokenValues(ctx)
		if err != nil {
			k.Logger(ctx).Error("unable to calculate token values", "error", err.Error())
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

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper.
type Hooks struct {
	k *Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

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
