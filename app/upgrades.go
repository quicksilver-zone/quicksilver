package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	claimsmanagertypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const (
	v001000UpgradeName = "v0.10.0"
	v001001UpgradeName = "v0.10.1"
	v001002UpgradeName = "v0.10.2"
	v001003UpgradeName = "v0.10.3"
	v001004UpgradeName = "v0.10.4"

	InnuendoChainID = "innuendo-3"
)

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v001000UpgradeName, getv001000Upgrade(app))
	app.UpgradeKeeper.SetUpgradeHandler(v001001UpgradeName, getv001001Upgrade(app))
	app.UpgradeKeeper.SetUpgradeHandler(v001002UpgradeName, getv001002Upgrade(app))
	app.UpgradeKeeper.SetUpgradeHandler(v001003UpgradeName, getv001003Upgrade(app))
	app.UpgradeKeeper.SetUpgradeHandler(v001004UpgradeName, getv001004Upgrade(app))

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
				sdk.NewInt(int64(app.InterchainstakingKeeper.GetParam(ctx, icstypes.KeyDepositInterval))),
				icstypes.ModuleName,
				"allbalances",
				0,
			)

			app.UpgradeKeeper.Logger(ctx).Info("state transitions complete.")

		case InnuendoChainID:
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
		case InnuendoChainID:
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

func getv001004Upgrade(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		switch ctx.ChainID() {
		case InnuendoChainID:
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.4; removing withdrawal records for previously removed zones.")
			app.InterchainstakingKeeper.IteratePrefixedWithdrawalRecords(ctx, []byte("fauxgaia-1"), func(_ int64, record icstypes.WithdrawalRecord) bool {
				app.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, "fauxgaia-1", record.Txhash, record.Status)
				return false
			})

			app.InterchainstakingKeeper.IteratePrefixedWithdrawalRecords(ctx, []byte("bitcanna-dev-5"), func(_ int64, record icstypes.WithdrawalRecord) bool {
				app.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, "bitcanna-dev-5", record.Txhash, record.Status)
				return false
			})

			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.4; removing unbonding records for previously removed zones.")
			app.InterchainstakingKeeper.IteratePrefixedUnbondingRecords(ctx, []byte("fauxgaia-1"), func(_ int64, record icstypes.UnbondingRecord) bool {
				app.InterchainstakingKeeper.DeleteUnbondingRecord(ctx, "fauxgaia-1", record.Validator, record.EpochNumber)
				return false
			})

			app.InterchainstakingKeeper.IteratePrefixedUnbondingRecords(ctx, []byte("bitcanna-dev-5"), func(_ int64, record icstypes.UnbondingRecord) bool {
				app.InterchainstakingKeeper.DeleteUnbondingRecord(ctx, "bitcanna-dev-5", record.Validator, record.EpochNumber)
				return false
			})

			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.4; removing delegation records for previously removed zones.")
			fgZone, _ := app.InterchainstakingKeeper.GetZone(ctx, "fauxgaia-1")
			app.InterchainstakingKeeper.IterateAllDelegations(ctx, &fgZone, func(record icstypes.Delegation) (stop bool) {
				if err := app.InterchainstakingKeeper.RemoveDelegation(ctx, &fgZone, record); err != nil {
					panic(err)
				}
				return false
			})

			bcZone, _ := app.InterchainstakingKeeper.GetZone(ctx, "bitcanna-dev-5")
			app.InterchainstakingKeeper.IterateAllDelegations(ctx, &bcZone, func(record icstypes.Delegation) (stop bool) {
				if err := app.InterchainstakingKeeper.RemoveDelegation(ctx, &bcZone, record); err != nil {
					panic(err)
				}
				return false
			})

			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.4; tidy up withdrawal records pertaining to withdrawal for jailed validators bug.")
			app.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(_ int64, record icstypes.WithdrawalRecord) bool {
				if record.Status == 3 && record.CompletionTime.String() == "0001-01-01T00:00:00Z" {
					app.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
					// unbonding never happened here. credit burn_amount back to delegator.
					if err := app.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, sdk.MustAccAddressFromBech32(record.Delegator), sdk.NewCoins(record.BurnAmount)); err != nil {
						panic(err)
					}
				}

				if record.Status == 4 && record.CompletionTime.Before(ctx.BlockTime()) {
					app.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
					// unbonding completed, burn qAtoms to restore balance.
					if err := app.BankKeeper.BurnCoins(ctx, icstypes.ModuleName, sdk.NewCoins(record.BurnAmount)); err != nil {
						panic(err)
					}
				}
				return false
			})

			app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zoneInfo icstypes.Zone) (stop bool) {
				app.UpgradeKeeper.Logger(ctx).Info("re-asserting redemption rate after upgrade.")
				app.InterchainstakingKeeper.UpdateRedemptionRateNoBounds(ctx, zoneInfo)
				return false
			})

		default:
			// also do nothing
			app.UpgradeKeeper.Logger(ctx).Info("upgrade to v0.10.4; nothing to do.")
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
