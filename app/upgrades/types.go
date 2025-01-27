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
	TestChainID       = "testchain1"

	// testnet upgrades
	V010405rc6UpgradeName   = "v1.4.5-rc6"
	V010405rc7UpgradeName   = "v1.4.5-rc7"
	V010407rc0UpgradeName   = "v1.4.7-rc0"
	V010407rc1UpgradeName   = "v1.4.7-rc1"
	V010407rc2UpgradeName   = "v1.4.7-rc2"
	V010500rc0UpgradeName   = "v1.5.0-rc0"
	V010500rc1UpgradeName   = "v1.5.0-rc1"
	V010503rc0UpgradeName   = "v1.5.3-rc0"
	V010600beta0UpgradeName = "v1.6.0-beta0"
	V010600beta1UpgradeName = "v1.6.0-beta1"
	V010600rc0UpgradeName   = "v1.6.0-rc0"
	V010600rc1UpgradeName   = "v1.6.0-rc1"
	V010601rc0UpgradeName   = "v1.6.1-rc0"
	V010601rc2UpgradeName   = "v1.6.1-rc2"
	V010601rc3UpgradeName   = "v1.6.1-rc3"
	V010601rc4UpgradeName   = "v1.6.1-rc4"

	// mainnet upgrades
	V010217UpgradeName = "v1.2.17"
	V010405UpgradeName = "v1.4.5"
	V010406UpgradeName = "v1.4.6"
	V010407UpgradeName = "v1.4.7"
	V010500UpgradeName = "v1.5.0"
	V010501UpgradeName = "v1.5.1"
	V010503UpgradeName = "v1.5.3"
	V010504UpgradeName = "v1.5.4"
	V010505UpgradeName = "v1.5.5"
	V010601UpgradeName = "v1.6.1"
	V010602UpgradeName = "v1.6.2"
	V010603UpgradeName = "v1.6.3"
	V010604UpgradeName = "v1.6.4"

	V010700UpgradeName = "v1.7.0"
	V010702UpgradeName = "v1.7.2"
	V010704UpgradeName = "v1.7.4"
	V010705UpgradeName = "v1.7.5"
	V010706UpgradeName = "v1.7.6"
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
