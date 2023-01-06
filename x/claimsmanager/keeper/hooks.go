package keeper

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclitypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	if epochIdentifier == "epoch" && epochNumber > 1 {
		if strings.Contains(ctx.ChainID(), "-") {
			revisionNum, err := strconv.ParseUint(strings.Split(ctx.ChainID(), "-")[1], 10, 64)
			if err != nil {
				k.Logger(ctx).Error("Error getting revision number for client ")
			}

			height := ibcclitypes.Height{
				RevisionNumber: revisionNum,
				RevisionHeight: uint64(ctx.BlockHeight() - 1),
			}

			selfConsState, err := k.IBCKeeper.ClientKeeper.GetSelfConsensusState(ctx, height)
			if err != nil {
				k.Logger(ctx).Error("Error getting self consensus state of previous height")
			}

			state := selfConsState.(*ibctmtypes.ConsensusState)
			k.SetSelfConsensusState(ctx, *state)
		}
	}
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {}

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
