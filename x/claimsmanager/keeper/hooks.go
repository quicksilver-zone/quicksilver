package keeper

import (
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	if epochIdentifier == epochstypes.EpochIdentifierEpoch && epochNumber > 1 {
		if err := k.StoreSelfConsensusState(ctx, "epoch"); err != nil {
			k.Logger(ctx).Error("unable to store consensus state", "error", err)
			return err
		}
	}
	return nil
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) error {

	// add the logic to update the heights of all the chains here in the claimable event
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

	return nil
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper.
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks.
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
