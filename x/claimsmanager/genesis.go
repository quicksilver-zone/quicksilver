package claimsmanager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/keeper"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
)

// InitGenesis initializes the claimsmanager module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, claim := range genState.Claims {
		k.SetClaim(ctx, claim)
	}
}

// ExportGenesis returns the claimsmanager module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Claims: k.AllClaims(ctx),
	}
}
