package claimsmanager

import (
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/keeper"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
