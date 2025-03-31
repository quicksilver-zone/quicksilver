package app

import (
	"fmt"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	crsistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	"github.com/quicksilver-zone/quicksilver/app/upgrades"
)

const (
	wasmModuleName    = "wasm"
	tfModuleName      = "tokenfactory"
	airdropModuleName = "airdrop"
)

func (app *Quicksilver) setUpgradeHandlers() {
	for _, upgrade := range upgrades.Upgrades() {
		if upgrade.UpgradeName == upgrades.V010800UpgradeName {
			_, err := ibctmmigrations.PruneExpiredConsensusStates(app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()}), app.appCodec, app.IBCKeeper.ClientKeeper)
			if err != nil {
				panic(fmt.Errorf("failed to prune expired consensus states: %w", err))
			}
		}

		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.mm,
				app.configurator,
				&app.AppKeepers,
			),
		)
	}
}

func (app *Quicksilver) setUpgradeStoreLoaders() {
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

	switch upgradeInfo.Name { // nolint:gocritic

	case upgrades.V010706UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Deleted: []string{airdropModuleName},
		}
	case upgrades.V010800UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{crsistypes.ModuleName},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
