package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/mint/keeper"
	"github.com/ingenuity-build/quicksilver/x/mint/types"
)

// InitGenesis new mint genesis.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak types.AccountKeeper, data *types.GenesisState) {
	data.Minter.EpochProvisions = data.Params.GenesisEpochProvisions
	k.SetMinter(ctx, data.Minter)
	k.SetParams(ctx, data.Params)

	if !ak.HasAccount(ctx, ak.GetModuleAddress(types.ModuleName)) {
		ak.GetModuleAccount(ctx, types.ModuleName)
	}

	k.SetLastReductionEpochNum(ctx, data.ReductionStartedEpoch)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)
	lastReductionEpoch := k.GetLastReductionEpochNum(ctx)
	return types.NewGenesis(minter, params, lastReductionEpoch)
}
