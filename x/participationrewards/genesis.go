package participationrewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/keeper"
	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, kpd := range genState.ProtocolData {
		k.SetProtocolData(ctx, []byte(kpd.Key), kpd.ProtocolData)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:       k.GetParams(ctx),
		ProtocolData: k.AllKeyedProtocolDatas(ctx),
	}
}
