package keepers

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	icacontroller "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardexported "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/exported"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/spf13/cast"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/quicksilver-zone/quicksilver/v7/utils"
	airdropkeeper "github.com/quicksilver-zone/quicksilver/v7/x/airdrop/keeper"
	airdroptypes "github.com/quicksilver-zone/quicksilver/v7/x/airdrop/types"
	claimsmanagerkeeper "github.com/quicksilver-zone/quicksilver/v7/x/claimsmanager/keeper"
	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/v7/x/claimsmanager/types"
	epochskeeper "github.com/quicksilver-zone/quicksilver/v7/x/epochs/keeper"
	epochstypes "github.com/quicksilver-zone/quicksilver/v7/x/epochs/types"
	interchainquerykeeper "github.com/quicksilver-zone/quicksilver/v7/x/interchainquery/keeper"
	interchainquerytypes "github.com/quicksilver-zone/quicksilver/v7/x/interchainquery/types"
	interchainstakingkeeper "github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/keeper"
	interchainstakingtypes "github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
	mintkeeper "github.com/quicksilver-zone/quicksilver/v7/x/mint/keeper"
	minttypes "github.com/quicksilver-zone/quicksilver/v7/x/mint/types"
	participationrewardskeeper "github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/keeper"
	participationrewardstypes "github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/types"
	supplykeeper "github.com/quicksilver-zone/quicksilver/v7/x/supply/keeper"
	tokenfactorykeeper "github.com/quicksilver-zone/quicksilver/v7/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/quicksilver-zone/quicksilver/v7/x/tokenfactory/types"
)

type AppKeepers struct {
	// make scoped keepers public for test purposes
	ScopedIBCKeeper                      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper                 capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper            capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper                  capabilitykeeper.ScopedKeeper
	ScopedInterchainStakingAccountKeeper capabilitykeeper.ScopedKeeper
	scopedWasmKeeper                     capabilitykeeper.ScopedKeeper // TODO: we can use this for testing
	ScopedIBCFeeKeeper                   capabilitykeeper.ScopedKeeper

	// "Normal" keepers
	// 		SDK
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.BaseKeeper
	DistrKeeper      distrkeeper.Keeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	GovKeeper        govkeeper.Keeper
	WasmKeeper       wasmkeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper

	// 		Quicksilver keepers
	EpochsKeeper               epochskeeper.Keeper
	MintKeeper                 mintkeeper.Keeper
	ClaimsManagerKeeper        claimsmanagerkeeper.Keeper
	InterchainstakingKeeper    *interchainstakingkeeper.Keeper
	InterchainQueryKeeper      interchainquerykeeper.Keeper
	ParticipationRewardsKeeper *participationrewardskeeper.Keeper
	AirdropKeeper              *airdropkeeper.Keeper
	TokenFactoryKeeper         tokenfactorykeeper.Keeper
	SupplyKeeper               supplykeeper.Keeper
	ConsensusParamsKeeper      consensusparamkeeper.Keeper

	// 		IBC keepers
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper        ibcfeekeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	PacketForwardKeeper *packetforwardkeeper.Keeper

	// Modules
	ICAModule           ica.AppModule
	TransferModule      transfer.AppModule
	PacketForwardModule packetforward.AppModule

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey
}

func NewAppKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	blockedAddresses map[string]bool,
	skipUpgradeHeights map[int64]bool,
	mock bool,
	homePath string,
	appOpts servertypes.AppOptions,
	wasmDir string,
	supplyEndpointEnabled bool,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	/*
		configure state listening capabilities using AppOptions
		we are doing nothing with the returned streamingServices and waitGroup in this case
	*/
	if err := bApp.RegisterStreamingServices(appOpts, appKeepers.keys); err != nil {
		tmos.Exit(err.Error())
	}

	appKeepers.InitKeepers(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		blockedAddresses,
		skipUpgradeHeights,
		mock,
		homePath,
		appOpts,
		wasmDir,
		supplyEndpointEnabled,
	)

	appKeepers.SetupHooks()

	return appKeepers
}

// InitKeepers initializes all keepers.
func (appKeepers *AppKeepers) InitKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	blockedAddresses map[string]bool,
	skipUpgradeHeights map[int64]bool,
	mock bool,
	homePath string,
	appOpts servertypes.AppOptions,
	wasmDir string,
	supplyEndpointEnabled bool,
) {
	// Add 'normal' keepers
	proofOpsFn := utils.ValidateProofOps
	if mock {
		proofOpsFn = utils.MockProofOps
	}

	selfProofOpsFn := utils.ValidateSelfProofOps
	if mock {
		selfProofOpsFn = utils.MockSelfProofOps
	}

	appKeepers.ParamsKeeper = appKeepers.initParamsKeeper(appCodec, legacyAmino, appKeepers.keys[paramstypes.StoreKey], appKeepers.tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec,
		runtime.NewKVStoreService(appKeepers.keys[consensusparamtypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(appKeepers.ConsensusParamsKeeper.ParamsStore)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAControllerKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedInterchainStakingKeeper := appKeepers.CapabilityKeeper.ScopeToModule(interchainstakingtypes.ModuleName)
	scopedWasmKeeper := appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	appKeepers.CapabilityKeeper.Seal()

	// TODO: Make a SetInvCheckPeriod fn on CrisisKeeper.
	// IMO, its bad design atm that it requires this in state machine initialization
	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(appCodec, runtime.NewKVStoreService(appKeepers.keys[crisistypes.StoreKey]), invCheckPeriod,
		appKeepers.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(), appKeepers.AccountKeeper.AddressCodec())
	// // get skipUpgradeHeights from the app options
	// skipUpgradeHeights := map[int64]bool{}
	// for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
	// 	skipUpgradeHeights[int64(h)] = true
	// }
	// homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	// set the governance module account as the authority for conducting upgrades
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(appKeepers.keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// use custom account for contracts
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[banktypes.StoreKey]),
		appKeepers.AccountKeeper,
		blockedAddresses,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		bApp.Logger(),
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[stakingtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(appCodec,
		runtime.NewKVStoreService(appKeepers.keys[distrtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	mintparamsSubspace, found := appKeepers.ParamsKeeper.GetSubspace(minttypes.ModuleName)
	if !found {
		panic("mint subspaces is not found")
	}
	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		mintparamsSubspace,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistrKeeper,
		&appKeepers.EpochsKeeper,
		authtypes.FeeCollectorName,
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(appKeepers.keys[slashingtypes.StoreKey]),
		appKeepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[feegrant.StoreKey]),
		appKeepers.AccountKeeper,
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[authzkeeper.StoreKey]),
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcexported.ModuleName],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		scopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// RouterKeeper must be created before TransferKeeper
	appKeepers.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[packetforwardtypes.StoreKey],
		appKeepers.TransferKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.DistrKeeper,
		appKeepers.BankKeeper,
		// TODO: https://github.com/osmosis-labs/osmosis/blob/891866b619754dddf13b871b394140bfa17c5025/app/keepers/keepers.go#L617-L618
		appKeepers.IBCKeeper.ChannelKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create Transfer Keepers
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.PacketForwardKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		scopedTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.SupplyKeeper = supplykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icacontrollertypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		utils.Keys[[]string](maccPerms),
		supplyEndpointEnabled,
	)
	appKeepers.PacketForwardKeeper.SetTransferKeeper(appKeepers.TransferKeeper)
	appKeepers.TransferModule = transfer.NewAppModule(appKeepers.TransferKeeper)

	// TODO: Not sure
	appKeepers.PacketForwardModule = packetforward.NewAppModule(appKeepers.PacketForwardKeeper, packetforwardexported.Subspace(appKeepers.GetSubspace(packetforwardtypes.ModuleName)))

	// ICA Keepers
	appKeepers.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icacontrollertypes.StoreKey],
		appKeepers.GetSubspace(icacontrollertypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper,
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.ICAModule = ica.NewAppModule(&appKeepers.ICAControllerKeeper, &appKeepers.ICAHostKeeper)

	appKeepers.ClaimsManagerKeeper = claimsmanagerkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[claimsmanagertypes.StoreKey],
		appKeepers.IBCKeeper,
	)

	// claimsmanagerModule := claimsmanager.NewAppModule(appCodec, appKeepers.ClaimsManagerKeeper)

	appKeepers.InterchainQueryKeeper = interchainquerykeeper.NewKeeper(appCodec, appKeepers.keys[interchainquerytypes.StoreKey], appKeepers.IBCKeeper)

	// interchainQueryModule := interchainquery.NewAppModule(appCodec, appKeepers.InterchainQueryKeeper)

	appKeepers.InterchainstakingKeeper = interchainstakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[interchainstakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.InterchainstakingKeeper.BankKeeper,
		appKeepers.ICAControllerKeeper,
		&scopedInterchainStakingKeeper,
		appKeepers.InterchainQueryKeeper,
		appKeepers.IBCKeeper,
		appKeepers.TransferKeeper,
		appKeepers.ClaimsManagerKeeper,
		appKeepers.GetSubspace(interchainstakingtypes.ModuleName),
	)

	// interchainstakingModule := interchainstaking.NewAppModule(appCodec, app.InterchainstakingKeeper)

	// interchainstakingIBCModule := interchainstaking.NewIBCModule(appKeepers.InterchainstakingKeeper)

	appKeepers.ParticipationRewardsKeeper = participationrewardskeeper.NewKeeper(
		appCodec,
		appKeepers.keys[participationrewardstypes.StoreKey],
		appKeepers.GetSubspace(participationrewardstypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.IBCKeeper,
		&appKeepers.InterchainQueryKeeper,
		appKeepers.InterchainstakingKeeper,
		appKeepers.ClaimsManagerKeeper,
		authtypes.FeeCollectorName,
		proofOpsFn,
		selfProofOpsFn,
	)

	if err := appKeepers.InterchainQueryKeeper.SetCallbackHandler(interchainstakingtypes.ModuleName, appKeepers.InterchainstakingKeeper.CallbackHandler()); err != nil {
		panic(err)
	}

	// participationrewardsModule := participationrewards.NewAppModule(appCodec, appKeepers.ParticipationRewardsKeeper)

	if err := appKeepers.InterchainQueryKeeper.SetCallbackHandler(participationrewardstypes.ModuleName, appKeepers.ParticipationRewardsKeeper.CallbackHandler()); err != nil {
		panic(err)
	}

	appKeepers.TokenFactoryKeeper = tokenfactorykeeper.NewKeeper(
		appKeepers.keys[tokenfactorytypes.StoreKey],
		appKeepers.GetSubspace(tokenfactorytypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper.WithMintCoinsRestriction(tokenfactorytypes.NewTokenFactoryDenomMintCoinsRestriction),
		appKeepers.DistrKeeper,
	)

	// Quicksilver Keepers
	appKeepers.EpochsKeeper = epochskeeper.NewKeeper(appCodec, appKeepers.keys[epochstypes.StoreKey])
	appKeepers.ParticipationRewardsKeeper.SetEpochsKeeper(appKeepers.EpochsKeeper)
	appKeepers.InterchainstakingKeeper.SetEpochsKeeper(appKeepers.EpochsKeeper)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	availableCapabilities := "iterator,staking,stargate,osmosis"
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}
	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[wasmtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		availableCapabilities,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		[]wasmkeeper.Option{}...,
	)

	// icaControllerIBCModule := icacontroller.NewIBCMiddleware(interchainstakingIBCModule, appKeepers.ICAControllerKeeper)
	// Create Transfer Stack
	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, appKeepers.IBCFeeKeeper)

	var icaControllerStack porttypes.IBCModule

	var noAuthzModule porttypes.IBCModule
	icaControllerStack = icacontroller.NewIBCMiddleware(noAuthzModule, appKeepers.ICAControllerKeeper)
	icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerStack, appKeepers.IBCFeeKeeper)

	var ibcStack porttypes.IBCModule
	ibcStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	ibcStack = packetforward.NewIBCMiddleware(
		ibcStack,
		appKeepers.PacketForwardKeeper,
		0,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)
	// RecvPacket, message that originates from core IBC and goes down to app, the flow is:
	// channel.RecvPacket -> fee.OnRecvPacket -> icaHost.OnRecvPacket
	var icaHostStack porttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(appKeepers.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, appKeepers.IBCFeeKeeper)

	// Create fee enabled wasm ibc Stack
	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, appKeepers.IBCFeeKeeper)

	// Create static IBC router, add app routes, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmtypes.ModuleName, wasmStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostStack)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	appKeepers.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[evidencetypes.StoreKey]),
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
		appKeepers.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper))
	govConfig := govtypes.DefaultConfig()
	/*
		Example of setting gov params:
		govConfig.MaxMetadataLen = 10000
	*/
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[govtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.DistrKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Set legacy router for backwards compatibility with gov v1beta1
	govKeeper.SetLegacyRouter(govRouter)

	// IBC Fee Module keeper
	appKeepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcfeetypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	appKeepers.AirdropKeeper = airdropkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[airdroptypes.StoreKey],
		appKeepers.GetSubspace(airdroptypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.GovKeeper,
		appKeepers.IBCKeeper,
		appKeepers.InterchainstakingKeeper,
		appKeepers.ParticipationRewardsKeeper,
		proofOpsFn,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	// airdropModule := airdrop.NewAppModule(appCodec, appKeepers.AirdropKeeper)

	appKeepers.ScopedIBCKeeper = scopedIBCKeeper
	appKeepers.ScopedTransferKeeper = scopedTransferKeeper
	appKeepers.ScopedICAControllerKeeper = scopedICAControllerKeeper
	appKeepers.ScopedInterchainStakingAccountKeeper = scopedInterchainStakingKeeper
	appKeepers.ScopedICAHostKeeper = scopedICAHostKeeper
	appKeepers.scopedWasmKeeper = scopedWasmKeeper
}

// initParamsKeeper init params keeper and its subspaces.
func (*AppKeepers) initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// SDK subspaces
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	// ibc subspaces
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(ibchost.SubModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName).WithKeyTable(packetforwardtypes.ParamKeyTable())
	// quicksilver subspaces
	paramsKeeper.Subspace(claimsmanagertypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(interchainstakingtypes.ModuleName)
	paramsKeeper.Subspace(interchainquerytypes.ModuleName)
	paramsKeeper.Subspace(participationrewardstypes.ModuleName)
	paramsKeeper.Subspace(airdroptypes.ModuleName)
	paramsKeeper.Subspace(tokenfactorytypes.ModuleName)
	// wasm subspace
	paramsKeeper.Subspace(wasmtypes.ModuleName)

	return paramsKeeper
}

// SetupHooks sets up hooks for modules.
func (appKeepers *AppKeepers) SetupHooks() {
	// For every module that has hooks set on it,
	// you must check InitNormalKeepers to ensure that its not passed by de-reference
	// e.g. *app.StakingKeeper doesn't appear

	// Recall that SetHooks is a mutative call.
	appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)

	appKeepers.EpochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(
			appKeepers.MintKeeper.Hooks(),
			appKeepers.ClaimsManagerKeeper.Hooks(),
			appKeepers.InterchainstakingKeeper.Hooks(),
			appKeepers.ParticipationRewardsKeeper.Hooks(),
		),
	)

	appKeepers.GovKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// insert governance hooks receivers here
		),
	)

	appKeepers.InterchainstakingKeeper.SetHooks(
		interchainstakingtypes.NewMultiIcsHooks(
			appKeepers.ParticipationRewardsKeeper.Hooks(),
		),
	)
}
