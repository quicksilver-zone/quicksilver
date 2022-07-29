package airdrop

import (
	"github.com/ingenuity-build/quicksilver/x/airdrop/keeper"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the airdrop module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	sum := uint64(0)
	zsum := make(map[string]uint64)
	for _, cr := range genState.ClaimRecords {
		if _, ok := zsum[cr.ChainId]; !ok {
			zsum[cr.ChainId] = 0
		}
		zsum[cr.ChainId] += cr.MaxAllocation
		sum += cr.MaxAllocation

		if err := k.SetClaimRecord(ctx, *cr); err != nil {
			panic(err)
		}
	}

	moduleBalance := k.GetModuleAccountBalance(ctx)
	if sum > moduleBalance.Amount.Uint64() {
		panic("insufficient airdrop module account balance")
	}

	moduleAddress := k.GetModuleAccountAddress(ctx)

	for _, zd := range genState.ZoneDrops {
		zs, ok := zsum[zd.ChainId]
		if !ok {
			panic("zone sum not found")
		}

		if zs != zd.Allocation {
			panic("zone sum does not match zone allocation")
		}

		zonedropAddress := k.GetZoneDropAccountAddress(ctx, zd.ChainId)
		err := k.BankKeeper.SendCoinsFromModuleToModule(
			ctx,
			moduleAddress.String(),
			zonedropAddress.String(),
			sdk.NewCoins(
				sdk.NewCoin(
					k.StakingKeeper.BondDenom(ctx),
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
