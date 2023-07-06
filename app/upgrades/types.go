package upgrades

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/ingenuity-build/quicksilver/app/keepers"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch).
const (
	ProductionChainID     = "quicksilver-2"
	RhyeChainID           = "rhye-1"
	DevnetChainID         = "magic-1"
	TestChainID           = "testchain1"
	OsmosisTestnetChainID = "osmo-test-5"

	V010402rc1UpgradeName   = "v1.4.2-rc1"
	V010402rc2UpgradeName   = "v1.4.2-rc2"
	V010402rc3UpgradeName   = "v1.4.2-rc3"
	V010402rc4UpgradeName   = "v1.4.2-rc4"
	V010402rc5UpgradeName   = "v1.4.2-rc5"
	V010402rc6UpgradeName   = "v1.4.2-rc6"
	V010402rc7UpgradeName   = "v1.4.2-rc7"
	V010403rc0UpgradeName   = "v1.4.3-rc0"
	V010404beta0UpgradeName = "v1.4.4-beta.0"
	V010404beta1UpgradeName = "v1.4.4-beta.1"
	V010404beta5UpgradeName = "v1.4.4-beta.5"
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

//nolint:all //function useful for writing network specific upgrade handlers
func isTest(ctx sdk.Context) bool {
	return ctx.ChainID() == TestChainID
}

//nolint:all //function useful for writing network specific upgrade handlers
func isDevnet(ctx sdk.Context) bool {
	return ctx.ChainID() == DevnetChainID
}

//nolint:all //function useful for writing network specific upgrade handlers
func isTestnet(ctx sdk.Context) bool {
	return ctx.ChainID() == RhyeChainID
}

//nolint:all //function useful for writing network specific upgrade handlers
func isMainnet(ctx sdk.Context) bool {
	return ctx.ChainID() == ProductionChainID
}
