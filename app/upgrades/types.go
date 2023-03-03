package upgrades

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/ingenuity-build/quicksilver/app/keepers"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch)
const (
	ProductionChainID = "quicksilver-2"
	InnuendoChainID   = "innuendo-5"
	DevnetChainID     = "quicktest-1"
	TestChainID       = "testchain1"

	V010300UpgradeName    = "v1.3.0"
	V010400UpgradeName    = "v1.4.0"
	V010400rc6UpgradeName = "v1.4.0-rc6"
	V010400rc7UpgradeName = "v1.4.0-rc7"
	V010400rc8UpgradeName = "v1.4.0-rc8"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string

	// CreateUpgradeHandler defines the function that creates an upgrade handler
	CreateUpgradeHandler func(*module.Manager, module.Configurator, *keepers.AppKeepers) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades storetypes.StoreUpgrades
}

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
