package keepers

import (
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v6/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v6/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v6/packetforward/types"
	tmos "github.com/tendermint/tendermint/libs/os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	"github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v6/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"

	appconfig "github.com/quicksilver-zone/quicksilver/cmd/config"
	"github.com/quicksilver-zone/quicksilver/utils"
	airdropkeeper "github.com/quicksilver-zone/quicksilver/x/airdrop/keeper"
	airdroptypes "github.com/quicksilver-zone/quicksilver/x/airdrop/types"
	claimsmanagerkeeper "github.com/quicksilver-zone/quicksilver/x/claimsmanager/keeper"
	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	epochskeeper "github.com/quicksilver-zone/quicksilver/x/epochs/keeper"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	interchainquerykeeper "github.com/quicksilver-zone/quicksilver/x/interchainquery/keeper"
	interchainquerytypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking"
	interchainstakingkeeper "github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
	interchainstakingtypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	mintkeeper "github.com/quicksilver-zone/quicksilver/x/mint/keeper"
	minttypes "github.com/quicksilver-zone/quicksilver/x/mint/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards"
	participationrewardskeeper "github.com/quicksilver-zone/quicksilver/x/participationrewards/keeper"
	participationrewardstypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	supplykeeper "github.com/quicksilver-zone/quicksilver/x/supply/keeper"
)

type AppKeepers struct {
	// make scoped keepers public for test purposes
	ScopedIBCKeeper                      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper                 capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper            capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper                  capabilitykeeper.ScopedKeeper
	ScopedInterchainStakingAccountKeeper capabilitykeeper.ScopedKeeper

	// "Normal" keepers
	// 		SDK
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.BaseKeeper
	DistrKeeper      distrkeeper.Keeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	GovKeeper        govkeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper

	// 		Quicksilver keepers
	EpochsKeeper               epochskeeper.Keeper
	MintKeeper                 mintkeeper.Keeper
	ClaimsManagerKeeper        claimsmanagerkeeper.Keeper
	InterchainstakingKeeper    *interchainstakingkeeper.Keeper
	InterchainQueryKeeper      interchainquerykeeper.Keeper
	ParticipationRewardsKeeper *participationrewardskeeper.Keeper
	AirdropKeeper              *airdropkeeper.Keeper
	SupplyKeeper               supplykeeper.Keeper

	// 		IBC keepers
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
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
	invCheckPeriod uint,
	appOpts servertypes.AppOptions,
	supplyEndpointEnabled bool,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	/*
		configure state listening capabilities using AppOptions
		we are doing nothing with the returned streamingServices and waitGroup in this case
	*/
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, appKeepers.keys); err != nil {
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
		invCheckPeriod,
		appOpts,
		supplyEndpointEnabled,
	)

	appKeepers.SetupHooks()

	return appKeepers
}

// InitKeepers initializes all keepers.
func (appKeepers *AppKeepers) InitKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	maccPerms map[string][]string,
	blockedAddresses map[string]bool,
	skipUpgradeHeights map[int64]bool,
	mock bool,
	homePath string,
	invCheckPeriod uint,
	_ servertypes.AppOptions,
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

	appKeepers.ParamsKeeper = appKeepers.initParamsKeeper(appCodec, cdc, appKeepers.keys[paramstypes.StoreKey], appKeepers.tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAControllerKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedInterchainStakingKeeper := appKeepers.CapabilityKeeper.ScopeToModule(interchainstakingtypes.ModuleName)

	appKeepers.CapabilityKeeper.Seal()

	// TODO: Make a SetInvCheckPeriod fn on CrisisKeeper.
	// IMO, its bad design atm that it requires this in state machine initialization
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName), invCheckPeriod, appKeepers.BankKeeper, authtypes.FeeCollectorName,
	)

	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights, appKeepers.keys[upgradetypes.StoreKey], appCodec, homePath, bApp, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// use custom account for contracts
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		appKeepers.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
		appconfig.Bech32PrefixAccAddr,
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(banktypes.ModuleName),
		blockedAddresses,
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.GetSubspace(stakingtypes.ModuleName),
	)
	appKeepers.StakingKeeper = &stakingKeeper

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.GetSubspace(distrtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
	)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.GetSubspace(minttypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistrKeeper,
		&appKeepers.EpochsKeeper,
		authtypes.FeeCollectorName,
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.GetSubspace(slashingtypes.ModuleName),
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec, bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		stakingKeeper,
		appKeepers.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// RouterKeeper must be created before TransferKeeper
	appKeepers.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[packetforwardtypes.StoreKey],
		appKeepers.GetSubspace(packetforwardtypes.ModuleName),
		appKeepers.TransferKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.DistrKeeper,
		appKeepers.BankKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
	)

	// Create Transfer Keepers
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.PacketForwardKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		scopedTransferKeeper,
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
	appKeepers.PacketForwardModule = packetforward.NewAppModule(appKeepers.PacketForwardKeeper)

	// ICA Keepers
	appKeepers.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icacontrollertypes.StoreKey],
		appKeepers.GetSubspace(icacontrollertypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper,
		bApp.MsgServiceRouter(),
	)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)
	appKeepers.ICAModule = ica.NewAppModule(&appKeepers.ICAControllerKeeper, &appKeepers.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(appKeepers.ICAHostKeeper)

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
		appKeepers.AuthzKeeper,
		appKeepers.BankKeeper,
		appKeepers.ICAControllerKeeper,
		&scopedInterchainStakingKeeper,
		appKeepers.InterchainQueryKeeper,
		appKeepers.IBCKeeper,
		appKeepers.TransferKeeper,
		appKeepers.ClaimsManagerKeeper,
		appKeepers.GetSubspace(interchainstakingtypes.ModuleName),
	)

	// interchainstakingModule := interchainstaking.NewAppModule(appCodec, app.InterchainstakingKeeper)

	interchainstakingIBCModule := interchainstaking.NewIBCModule(appKeepers.InterchainstakingKeeper)

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

	// Quicksilver Keepers
	appKeepers.EpochsKeeper = epochskeeper.NewKeeper(appCodec, appKeepers.keys[epochstypes.StoreKey])
	appKeepers.ParticipationRewardsKeeper.SetEpochsKeeper(appKeepers.EpochsKeeper)
	appKeepers.InterchainstakingKeeper.SetEpochsKeeper(appKeepers.EpochsKeeper)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks

	icaControllerIBCModule := icacontroller.NewIBCMiddleware(interchainstakingIBCModule, appKeepers.ICAControllerKeeper)

	var ibcStack porttypes.IBCModule
	ibcStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	ibcStack = packetforward.NewIBCMiddleware(
		ibcStack,
		appKeepers.PacketForwardKeeper,
		0,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)
	ibcStack = interchainstaking.NewTransferMiddleware(ibcStack, appKeepers.InterchainstakingKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.
		AddRoute(ibctransfertypes.ModuleName, ibcStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(interchainstakingtypes.ModuleName, icaControllerIBCModule)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	appKeepers.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)

	govConfig := govtypes.DefaultConfig()

	// register the proposal types
	govRouter := govv1beta1.NewRouter()

	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(appKeepers.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(interchainstakingtypes.RouterKey, interchainstaking.NewProposalHandler(appKeepers.InterchainstakingKeeper)).
		AddRoute(participationrewardstypes.RouterKey, participationrewards.NewProposalHandler(appKeepers.ParticipationRewardsKeeper))
	// add custom proposal routes here.

	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.GetSubspace(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		govRouter,
		bApp.MsgServiceRouter(),
		govConfig,
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
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName).WithKeyTable(packetforwardtypes.ParamKeyTable())
	// quicksilver subspaces
	paramsKeeper.Subspace(claimsmanagertypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(interchainstakingtypes.ModuleName)
	paramsKeeper.Subspace(interchainquerytypes.ModuleName)
	paramsKeeper.Subspace(participationrewardstypes.ModuleName)
	paramsKeeper.Subspace(airdroptypes.ModuleName)

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
