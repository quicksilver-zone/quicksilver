package app

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const (
	ProductionChainID = "quicksilver-2"
	InnuendoChainID   = "innuendo-5"
	DevnetChainID     = "quicktest-1"
	TestChainID       = "testchain1"

	v010300UpgradeName    = "v1.3.0"
	v010400UpgradeName    = "v1.4.0"
	v010400rc5UpgradeName = "v1.4.0-rc5"
)

func isTest(ctx sdk.Context) bool {
	return ctx.ChainID() == TestChainID
}

func isDevnet(ctx sdk.Context) bool {
	return ctx.ChainID() == DevnetChainID
}

func isTestnet(ctx sdk.Context) bool {
	return ctx.ChainID() == InnuendoChainID
}

//nolint:all //function useful for writing network specific upgrade handlers
func isMainnet(ctx sdk.Context) bool {
	return ctx.ChainID() == ProductionChainID
}

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v010300UpgradeName, noOpHandler(app))
	app.UpgradeKeeper.SetUpgradeHandler(v010400UpgradeName, v010400UpgradeHandler(app))
	app.UpgradeKeeper.SetUpgradeHandler(v010400rc5UpgradeName, v010400rc5UpgradeHandler(app))

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

	switch upgradeInfo.Name { //nolint:gocritic
	// case v001000UpgradeName:

	// 	storeUpgrades = &storetypes.StoreUpgrades{
	// 		Added: []string{claimsmanagertypes.ModuleName},
	// 	}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func noOpHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func v010400UpgradeHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// upgrade zones
		app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone icstypes.Zone) (stop bool) {
			zone.DepositsEnabled = true
			zone.ReturnToSender = false
			zone.UnbondingEnabled = false
			zone.Decimals = 6
			app.InterchainstakingKeeper.SetZone(ctx, &zone)
			return false
		})

		// upgrade receipts
		time := ctx.BlockTime()
		for _, r := range app.InterchainstakingKeeper.AllReceipts(ctx) {
			r.FirstSeen = &time
			r.Completed = &time
			app.InterchainstakingKeeper.SetReceipt(ctx, r)
		}
		if isTestnet(ctx) || isTest(ctx) {

			app.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "uni-5")

			// burn uqjunox
			addr1, err := utils.AccAddressFromBech32("quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq", "quick")
			if err != nil {
				return nil, err
			}
			addr2, err := utils.AccAddressFromBech32("quick16x03wcp37kx5e8ehckjxvwcgk9j0cqnhcccnty", "quick")
			if err != nil {
				return nil, err
			}

			err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(1600000))))
			if err != nil {
				return nil, err
			}

			err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr2, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(200000000))))
			if err != nil {
				return nil, err
			}

			err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.EscrowModuleAccount, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(400000))))
			if err != nil {
				return nil, err
			}

			err = app.BankKeeper.BurnCoins(ctx, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(202000000))))
			if err != nil {
				return nil, err
			}
		}
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

func v010400rc5UpgradeHandler(app *Quicksilver) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) {
			app.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "regen-redwood-1")
			// re-register regen-redwood-1 with new connection
			regenProp := icstypes.NewRegisterZoneProposal("register regen-redwood-1 zone",
				"register regen-redwood-1  (regen-testnet) zone with multisend and lsm disabled",
				"connection-8",
				"uregen",
				"uqregen",
				"regen",
				false,
				true,
				true,
				false,
				6)
			err := icskeeper.HandleRegisterZoneProposal(ctx, app.InterchainstakingKeeper, regenProp)
			if err != nil {
				return nil, err
			}
		}

		// remove expired failed redelegation records
		app.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(_ int64, key []byte, record icstypes.RedelegationRecord) (stop bool) {
			if record.CompletionTime.Equal(time.Time{}) {
				app.InterchainstakingKeeper.DeleteRedelegationRecordByKey(ctx, key)
			}
			return false
		})

		// remove and refund failed unbondings
		app.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(index int64, record icstypes.WithdrawalRecord) (stop bool) {
			if record.Status == icskeeper.WithdrawStatusUnbond && record.CompletionTime.Equal(time.Time{}) {
				delegatorAcc, err := utils.AccAddressFromBech32(record.Delegator, "quick")
				if err != nil {
					panic(err)
				}
				if err = app.InterchainstakingKeeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.EscrowModuleAccount, delegatorAcc, sdk.NewCoins(record.BurnAmount)); err != nil {
					panic(err)
				}
				app.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
			}
			return false
		})

		if isTestnet(ctx) || isDevnet(ctx) {
			app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zoneInfo icstypes.Zone) (stop bool) {
				app.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, zoneInfo)
				return false
			})
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
