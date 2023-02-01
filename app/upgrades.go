package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const (
	ProductionChainID = "quicksilver-2"
	InnuendoChainID   = "innuendo-4"
	DevnetChainID     = "quicktest-1"

	v010300UpgradeName = "v1.3.0"
	v010400UpgradeName = "v1.4.0"
)

func setUpgradeHandlers(app *Quicksilver) {
	app.UpgradeKeeper.SetUpgradeHandler(v010300UpgradeName, noOpHandler(app))
	app.UpgradeKeeper.SetUpgradeHandler(v010400UpgradeName, v010400UpgradeHandler(app))

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
		app.InterchainstakingKeeper.V1_2_IterateZones(ctx, func(index int64, legacyZone types.V1_2_Zone) (stop bool) {
			newZone := types.Zone{
				ConnectionId:                 legacyZone.ConnectionId,
				ChainId:                      legacyZone.ChainId,
				DepositAddress:               legacyZone.DepositAddress,
				WithdrawalAddress:            legacyZone.WithdrawalAddress,
				PerformanceAddress:           legacyZone.PerformanceAddress,
				DelegationAddress:            legacyZone.DelegationAddress,
				AccountPrefix:                legacyZone.AccountPrefix,
				LocalDenom:                   legacyZone.LocalDenom,
				BaseDenom:                    legacyZone.BaseDenom,
				RedemptionRate:               legacyZone.RedemptionRate,
				LastRedemptionRate:           legacyZone.LastRedemptionRate,
				Validators:                   legacyZone.Validators,
				AggregateIntent:              legacyZone.AggregateIntent,
				ReturnToSender:               false,
				LiquidityModule:              legacyZone.LiquidityModule,
				WithdrawalWaitgroup:          legacyZone.WithdrawalWaitgroup,
				IbcNextValidatorsHash:        legacyZone.IbcNextValidatorsHash,
				ValidatorSelectionAllocation: legacyZone.ValidatorSelectionAllocation,
				HoldingsAllocation:           legacyZone.HoldingsAllocation,
				Tvl:                          legacyZone.Tvl,
				UnbondingPeriod:              legacyZone.UnbondingPeriod,
				UnbondingEnabled:             false,
				Decimals:                     6,
				DepositsEnabled:              true,
			}
			app.InterchainstakingKeeper.SetZone(ctx, &newZone)
			return false
		})

		// upgrade receipts
		time := ctx.BlockTime()
		for _, r := range app.InterchainstakingKeeper.AllReceipts(ctx) {
			r.FirstSeen = &time
			r.Completed = &time
			app.InterchainstakingKeeper.SetReceipt(ctx, r)
		}
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}
