package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	claimsmanagertypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const (
	v001000UpgradeName = "v0.10.0"
	v001001UpgradeName = "v0.10.1"
	v001002UpgradeName = "v0.10.2"
	v001003UpgradeName = "v0.10.3"
)

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v001000UpgradeName, getv001000Upgrade(app))
	app.UpgradeKeeper.SetUpgradeHandler(v001001UpgradeName, getv001001Upgrade(app))
	app.UpgradeKeeper.SetUpgradeHandler(v001002UpgradeName, getv001002Upgrade(app))
	app.UpgradeKeeper.SetUpgradeHandler(v001003UpgradeName, getv001003Upgrade(app))

	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name {
	case v001000UpgradeName:

		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{claimsmanagertypes.ModuleName},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func getv001000Upgrade(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.0; no state transitions to apply.")
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func getv001001Upgrade(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.1; no state transitions to apply.")
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func getv001002Upgrade(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		switch ctx.ChainID() {
		case "quicktest-1": // devnet
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.2; reinstating queries for quickosmo-1.")

			// deposit
			zone, _ := app.InterchainstakingKeeper.GetZone(ctx, "quickosmo-1")
			balanceQuery := banktypes.QueryAllBalancesRequest{Address: zone.DepositAddress.Address}
			bz, _ := app.InterchainstakingKeeper.GetCodec().Marshal(&balanceQuery)

			app.InterchainstakingKeeper.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				"cosmos.bank.v1beta1.Query/AllBalances",
				bz,
				sdk.NewInt(int64(app.InterchainstakingKeeper.GetParam(ctx, types.KeyDepositInterval))),
				types.ModuleName,
				"allbalances",
				0,
			)

			// // quickgaia-1
			// zone, _ := app.InterchainstakingKeeper.GetZone(ctx, "quickgaia-1")
			// for _, addr := range []string{"cosmosvaloper1759teakrsvnx7rnur8ezc4qaq8669nhtgukm0x", "cosmosvaloper1jtjjyxtqk0fj85ud9cxk368gr8cjdsftvdt5jl", "cosmosvaloper1q86m0zq0p52h4puw5pg5xgc3c5e2mq52y6mth0"} {
			// 	app.InterchainstakingKeeper.SetPerformanceDelegation(ctx, &zone, icstypes.NewDelegation(zone.PerformanceAddress.Address, addr, sdk.NewInt64Coin(zone.BaseDenom, 10000)))
			// }
			// // quickosmo-1
			// zone, _ = app.InterchainstakingKeeper.GetZone(ctx, "quickosmo-1")
			// for _, addr := range []string{"osmovaloper10843vrvy6jh3t4mxt2fnvkwm7mwewhkdwcqmuj"} {
			// 	app.InterchainstakingKeeper.SetPerformanceDelegation(ctx, &zone, icstypes.NewDelegation(zone.PerformanceAddress.Address, addr, sdk.NewInt64Coin(zone.BaseDenom, 10000)))
			// }
			// // quickstar-1
			// zone, _ = app.InterchainstakingKeeper.GetZone(ctx, "quickstar-1")
			// for _, addr := range []string{"starsvaloper1d54z0yptatca3a05pqyvdv5jzpu8p5fmhkcw69", "starsvaloper13x85ct9jkmygqhyaf930c7x2564vqzq37kksc7", "starsvaloper1j0cjvx9u7kgs3gkqra8gmf7y396mv433d3zcut"} {
			// 	app.InterchainstakingKeeper.SetPerformanceDelegation(ctx, &zone, icstypes.NewDelegation(zone.PerformanceAddress.Address, addr, sdk.NewInt64Coin(zone.BaseDenom, 10000)))
			// }
			app.UpgradeKeeper.Logger(ctx).Info("state transitions complete.")

		case "innuendo-3":
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.2; removing osmo-test-4 zone.")
			app.InterchainstakingKeeper.DeleteZone(ctx, "osmo-test-4")
		default:
			// also do nothing
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.2; nothing to do.")
		}

		app.InterchainQueryKeeper.Logger(ctx).Info("removing legacy perfbalance queries.")

		for _, query := range app.InterchainQueryKeeper.AllQueries(ctx) {
			if query.CallbackId == "perfbalance" && query.Period.Equal(sdk.NewInt(-1)) {
				app.InterchainQueryKeeper.DeleteQuery(ctx, query.Id)
				app.InterchainQueryKeeper.Logger(ctx).Info("removed query", "id", query.Id, "chain", query.ChainId)
			}
		}

		app.InterchainQueryKeeper.Logger(ctx).Info("emitting v2 periodic perfbalance queries.")

		for _, zone := range app.InterchainstakingKeeper.AllZones(ctx) {
			zone := zone
			if err := app.InterchainstakingKeeper.EmitPerformanceBalanceQuery(ctx, &zone); err != nil {
				panic(err)
			}
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func getv001003Upgrade(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		switch ctx.ChainID() {
		case "innuendo-3":
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.3; removing defunct zones.")
			app.InterchainstakingKeeper.DeleteZone(ctx, "bitcanna-dev-5")
			app.InterchainstakingKeeper.DeleteZone(ctx, "fauxgaia-1")
			app.InterchainstakingKeeper.DeleteZone(ctx, "uni-5")
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.3; removing queries for defunct zones.")
			for _, query := range app.InterchainQueryKeeper.AllQueries(ctx) {
				if query.ChainId == "bitcanna-dev-5" || query.ChainId == "fauxgaia-1" || query.ChainId == "uni-5" {
					app.InterchainQueryKeeper.DeleteQuery(ctx, query.Id)
					app.InterchainQueryKeeper.Logger(ctx).Info("removed query", "id", query.Id, "chain", query.ChainId)
				}
			}

		default:
			// also do nothing
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.3; nothing to do.")
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
