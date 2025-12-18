package upgrades

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

// upgrade name consts: vMMmmppUpgradeName (M=Major, m=minor, p=patch).
const (
	ProductionChainID = "quicksilver-2"
	RhyeChainID       = "rhye-2"
	DevnetChainID     = "magic-2"
	TestChainID       = "testchain1-1"

	// mainnet
	V010700UpgradeName = "v1.7.0"
	V010702UpgradeName = "v1.7.2"
	V010704UpgradeName = "v1.7.4"
	V010705UpgradeName = "v1.7.5"
	V010706UpgradeName = "v1.7.6"
	V010707UpgradeName = "v1.7.7"

	V010800r1UpgradeName = "v1.8.0-rc.1"
	V010800UpgradeName   = "v1.8.0"
	V010801UpgradeName   = "v1.8.1"

	V010900UpgradeName = "v1.9.0"

	V0101000rc0UpgradeName = "v1.10.0-rc.0"
	V0101000UpgradeName    = "v1.10.0"
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

// nolint:all //function useful for writing network specific upgrade handlers
func isTest(ctx sdk.Context) bool {
	return ctx.ChainID() == TestChainID
}

// nolint:all //function useful for writing network specific upgrade handlers
func isDevnet(ctx sdk.Context) bool {
	return ctx.ChainID() == DevnetChainID
}

// nolint:all //function useful for writing network specific upgrade handlers
func isTestnet(ctx sdk.Context) bool {
	return ctx.ChainID() == RhyeChainID
}

// nolint:all //function useful for writing network specific upgrade handlers
func isMainnet(ctx sdk.Context) bool {
	return ctx.ChainID() == ProductionChainID
}
