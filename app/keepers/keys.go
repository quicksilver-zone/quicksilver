package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v7/router/types"

	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"
	claimsmanagertypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	interchainquerytypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	interchainstakingtypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

// TODO: We need to automate this, by bundling with a module struct...
func KVStoreKeys() []string {
	return []string{
		// SDK keys
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		capabilitytypes.StoreKey,
		feegrant.StoreKey,
		authzkeeper.StoreKey,
		consensusparamtypes.StoreKey,
		// ibc keys
		ibcexported.StoreKey,
		ibctransfertypes.StoreKey,
		icacontrollertypes.StoreKey,
		icahosttypes.StoreKey,
		packetforwardtypes.StoreKey,
		// quicksilver keys
		minttypes.StoreKey,
		claimsmanagertypes.StoreKey,
		epochstypes.StoreKey,
		interchainstakingtypes.StoreKey,
		interchainquerytypes.StoreKey,
		participationrewardstypes.StoreKey,
		airdroptypes.StoreKey,
		wasm.StoreKey,
		tokenfactorytypes.StoreKey,
	}
}

// GenerateKeys generates new keys (KV Store, Transient store, and memory store).
func (appKeepers *AppKeepers) GenerateKeys() {
	// Define what keys will be used in the cosmos-sdk key/value store.
	// Cosmos-SDK modules each have a "key" that allows the application to reference what they've stored on the chain.
	appKeepers.keys = sdk.NewKVStoreKeys(KVStoreKeys()...)

	// Define transient store keys
	appKeepers.tkeys = sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	// MemKeys are for information that is stored only in RAM.
	appKeepers.memKeys = sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
}

// GetSubspace gets existing substore from keeper.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// GetKVStoreKey gets KV Store keys.
func (appKeepers *AppKeepers) GetKVStoreKey() map[string]*storetypes.KVStoreKey {
	return appKeepers.keys
}

// GetTransientStoreKey gets Transient Store keys.
func (appKeepers *AppKeepers) GetTransientStoreKey() map[string]*storetypes.TransientStoreKey {
	return appKeepers.tkeys
}

// GetMemoryStoreKey get memory Store keys.
func (appKeepers *AppKeepers) GetMemoryStoreKey() map[string]*storetypes.MemoryStoreKey {
	return appKeepers.memKeys
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetKey(storeKey string) *storetypes.KVStoreKey {
	return appKeepers.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return appKeepers.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (appKeepers *AppKeepers) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return appKeepers.memKeys[storeKey]
}
