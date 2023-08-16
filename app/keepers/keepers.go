package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
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
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	packetforward "github.com/strangelove-ventures/packet-forward-middleware/v7/router"
	packetforwardkeeper "github.com/strangelove-ventures/packet-forward-middleware/v7/router/keeper"
	packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v7/router/types"

	"github.com/ingenuity-build/quicksilver/utils"

	appconfig "github.com/ingenuity-build/quicksilver/cmd/config"
	"github.com/ingenuity-build/quicksilver/wasmbinding"
	airdropkeeper "github.com/ingenuity-build/quicksilver/x/airdrop/keeper"
	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"
	claimsmanagerkeeper "github.com/ingenuity-build/quicksilver/x/claimsmanager/keeper"
	claimsmanagertypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	epochskeeper "github.com/ingenuity-build/quicksilver/x/epochs/keeper"
	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	interchainquerykeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	interchainquerytypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking"
	interchainstakingkeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	interchainstakingtypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	mintkeeper "github.com/ingenuity-build/quicksilver/x/mint/keeper"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards"
	participationrewardskeeper "github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	tokenfactorykeeper "github.com/ingenuity-build/quicksilver/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

type AppKeepers struct {
	// make scoped keepers public for test purposes
	ScopedIBCKeeper                      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper                 capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper            capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper                  capabilitykeeper.ScopedKeeper
	scopedWasmKeeper                     capabilitykeeper.ScopedKeeper // TODO: we can use this for testing

	// "Normal" keepers
	// 		SDK
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.BaseKeeper
	DistrKeeper           distrkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	GovKeeper             govkeeper.Keeper
	WasmKeeper            wasm.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// 		Quicksilver keepers
	EpochsKeeper               epochskeeper.Keeper
	MintKeeper                 mintkeeper.Keeper
	ClaimsManagerKeeper        claimsmanagerkeeper.Keeper
	InterchainstakingKeeper    *interchainstakingkeeper.Keeper
	InterchainQueryKeeper      interchainquerykeeper.Keeper
	ParticipationRewardsKeeper *participationrewardskeeper.Keeper
	AirdropKeeper              *airdropkeeper.Keeper
	TokenFactoryKeeper         tokenfactorykeeper.Keeper

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
	ICSModule           interchainstaking.IBCModule

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
	wasmConfig wasm.Config,
	wasmEnabledProposals []wasm.ProposalType,
	wasmOpts []wasm.Option,
	logger log.Logger,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()
	/*
		configure state listening capabilities using AppOptions
		we are doing nothing with the returned streamingServices and waitGroup in this case
	*/
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, appKeepers.keys); err != nil {
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
		wasmConfig,
		wasmEnabledProposals,
		wasmOpts,
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
	_ servertypes.AppOptions,
	wasmDir string,
	wasmConfig wasm.Config,
	wasmEnabledProposals []wasm.ProposalType,
	wasmOpts []wasm.Option,
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
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, appKeepers.keys[consensusparamtypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String())
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAControllerKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedWasmKeeper := appKeepers.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
	appKeepers.CapabilityKeeper.Seal()

	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights, appKeepers.keys[upgradetypes.StoreKey], appCodec, homePath, bApp, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// use custom account for contracts
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		appconfig.Bech32PrefixAccAddr,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		blockedAddresses,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		cdc,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
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
		appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
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
		&scopedICAControllerKeeper,
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
		&scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
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
		appKeepers.BankKeeper,
		appKeepers.ICAControllerKeeper,
		appKeepers.InterchainQueryKeeper,
		appKeepers.IBCKeeper,
		appKeepers.TransferKeeper,
		appKeepers.ClaimsManagerKeeper,
		appKeepers.GetSubspace(interchainstakingtypes.ModuleName),
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// interchainstakingModule := interchainstaking.NewAppModule(appCodec, app.InterchainstakingKeeper)
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
		appKeepers.BankKeeper.WithMintCoinsRestriction(tokenfactorytypes.NewTokenFactoryDenomMintCoinsRestriction()),
		appKeepers.DistrKeeper,
	)

	// Quicksilver Keepers
	appKeepers.EpochsKeeper = epochskeeper.NewKeeper(appCodec, appKeepers.keys[epochstypes.StoreKey])
	appKeepers.ParticipationRewardsKeeper.SetEpochsKeeper(appKeepers.EpochsKeeper)
	appKeepers.InterchainstakingKeeper.SetEpochsKeeper(&appKeepers.EpochsKeeper)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate,osmosis"
	wasmOpts = append(wasmbinding.RegisterCustomPlugins(&appKeepers.BankKeeper, &appKeepers.TokenFactoryKeeper), wasmOpts...)
	wasmOpts = append(wasmbinding.RegisterStargateQueries(*bApp.GRPCQueryRouter(), appCodec), wasmOpts...)
	appKeepers.WasmKeeper = wasm.NewKeeper(
		appCodec,
		appKeepers.keys[wasm.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.PacketForwardKeeper, // ISC4 Wrapper: fee IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	var icaControllerStack porttypes.IBCModule
	icaControllerStack = interchainstaking.NewIBCModule(appKeepers.InterchainstakingKeeper)
	appKeepers.ICSModule = icaControllerStack.(interchainstaking.IBCModule)
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, appKeepers.ICAControllerKeeper)

	icaHostIBCModule := icahost.NewIBCModule(appKeepers.ICAHostKeeper)

	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		appKeepers.PacketForwardKeeper,
		0,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasm.ModuleName, wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.PacketForwardKeeper)).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(interchainstakingtypes.ModuleName+icacontrollertypes.SubModuleName, icaControllerStack)
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
	// The gov proposal types can be individually enabled
	if len(wasmEnabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(appKeepers.WasmKeeper, wasmEnabledProposals))
	}

	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(interchainstakingtypes.RouterKey, interchainstaking.NewProposalHandler(appKeepers.InterchainstakingKeeper)).
		AddRoute(participationrewardstypes.RouterKey, participationrewards.NewProposalHandler(appKeepers.ParticipationRewardsKeeper))
	// add custom proposal routes here.

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

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
	appKeepers.ScopedICAHostKeeper = scopedICAHostKeeper
	appKeepers.scopedWasmKeeper = scopedWasmKeeper
}

// initParamsKeeper init params keeper and its subspaces.
func (appKeepers *AppKeepers) initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	// SDK subspaces
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	// ibc subspaces
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
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
	paramsKeeper.Subspace(wasm.ModuleName)
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
	appKeepers.InterchainstakingKeeper.SetHooks(
		interchainstakingtypes.NewMultiIcsHooks(
			appKeepers.ParticipationRewardsKeeper.Hooks(),
		),
	)
}
