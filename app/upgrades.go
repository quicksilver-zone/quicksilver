package app

import (
	"errors"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
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
		app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone types.Zone) (stop bool) {
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
		if ctx.ChainID() == "innuendo-5" {

			// clear uni-5 unbondings
			app.InterchainstakingKeeper.IteratePrefixedUnbondingRecords(ctx, []byte("uni-5"), func(_ int64, record types.UnbondingRecord) (stop bool) {
				app.InterchainstakingKeeper.DeleteUnbondingRecord(ctx, record.ChainId, record.Validator, record.EpochNumber)
				return false
			})

			// clear uni-5 redelegations
			app.InterchainstakingKeeper.IteratePrefixedRedelegationRecords(ctx, []byte("uni-5"), func(_ int64, _ []byte, record types.RedelegationRecord) (stop bool) {
				app.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, record.ChainId, record.Source, record.Destination, record.EpochNumber)
				return false
			})

			// remove uni-5 zone and related records
			app.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone types.Zone) (stop bool) {
				if zone.ChainId == "uni-5" {
					// remove uni-5 delegation records
					app.InterchainstakingKeeper.IterateAllDelegations(ctx, &zone, func(delegation types.Delegation) (stop bool) {
						err := app.InterchainstakingKeeper.RemoveDelegation(ctx, &zone, delegation)
						if err != nil {
							panic(err)
						}
						return false
					})

					// remove uni-5 performance delegation records
					app.InterchainstakingKeeper.IterateAllPerformanceDelegations(ctx, &zone, func(delegation types.Delegation) (stop bool) {
						err := app.InterchainstakingKeeper.RemoveDelegation(ctx, &zone, delegation)
						if err != nil {
							panic(err)
						}
						return false
					})
					// remove uni-5 receipts
					app.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
						app.InterchainstakingKeeper.DeleteReceipt(ctx, icskeeper.GetReceiptKey(receiptInfo.ChainId, receiptInfo.Txhash))
						return false
					})

					// remove zone withdrawal records
					app.InterchainstakingKeeper.IterateZoneWithdrawalRecords(ctx, zone.ChainId, func(index int64, record types.WithdrawalRecord) (stop bool) {
						app.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, zone.ChainId, record.Txhash, record.Status)
						return false
					})

					app.InterchainstakingKeeper.DeleteZone(ctx, zone.ChainId)

				}
				return false
			})

			// remove uni-5 queries in state
			app.InterchainQueryKeeper.IterateQueries(ctx, func(_ int64, queryInfo icqtypes.Query) (stop bool) {
				if queryInfo.ChainId == "uni-5" {
					app.InterchainQueryKeeper.DeleteQuery(ctx, queryInfo.Id)
				}
				return false
			})

			// burn uqjunox
			addr1, err := AccAddressFromBech32("quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq", "quick")
			if err != nil {
				return nil, err
			}
			addr2, err := AccAddressFromBech32("quick16x03wcp37kx5e8ehckjxvwcgk9j0cqnhcccnty", "quick")
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

			err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, types.EscrowModuleAccount, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(400000))))
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

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address, prefix string) (addr sdk.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.AccAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
