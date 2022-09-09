package airdrop

import (
	"fmt"

	"github.com/ingenuity-build/quicksilver/x/airdrop/keeper"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the airdrop module's state from a provided genesis
// state.
// This function will panic on failure.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	sum := uint64(0)
	zsum := make(map[string]uint64)
	for _, cr := range genState.ClaimRecords {
		zsum[cr.ChainId] += cr.MaxAllocation
		sum += cr.MaxAllocation

		if err := k.SetClaimRecord(ctx, *cr); err != nil {
			panic(err)
		}
	}

	moduleBalance := k.GetModuleAccountBalance(ctx)
	if sum > moduleBalance.Amount.Uint64() {
		panic(fmt.Sprintf("insufficient airdrop module account balance for airdrop module account %s, expected %d", k.GetModuleAccountAddress(ctx), sum))
	}

	for _, zd := range genState.ZoneDrops {
		zs, ok := zsum[zd.ChainId]
		if !ok {
			panic("zone sum not found")
		}

		if zs != zd.Allocation {
			panic(fmt.Sprintf("zone sum does not match zone allocation; got %d, allocated %d", zs, zd.Allocation))
		}

		zonedropAddress := k.GetZoneDropAccountAddress(zd.ChainId)

		err := k.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			zonedropAddress,
			sdk.NewCoins(
				sdk.NewCoin(
					k.BondDenom(ctx),
					sdk.NewIntFromUint64(zd.Allocation),
				),
			),
		)
		if err != nil {
			panic(err)
		}

		k.SetZoneDrop(ctx, *zd)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	zoneDrops := k.AllZoneDrops(ctx)
	claimRecords := k.AllClaimRecords(ctx)

	return types.NewGenesisState(params, zoneDrops, claimRecords)
}
