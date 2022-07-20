package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ingenuity-build/quicksilver/x/mint"
	mintkeeper "github.com/ingenuity-build/quicksilver/x/mint/keeper"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"

	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v4/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v4/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v4/modules/apps/29-fee/types"

	"github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v4/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/client"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibcmock "github.com/cosmos/ibc-go/v4/testing/mock"

	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v4/testing"

	"github.com/ingenuity-build/quicksilver/x/epochs"
	epochskeeper "github.com/ingenuity-build/quicksilver/x/epochs/keeper"
	epochstypes "github.com/ingenuity-build/quicksilver/x/epochs/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking"
	interchainstakingclient "github.com/ingenuity-build/quicksilver/x/interchainstaking/client"
	interchainstakingkeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	interchainstakingtypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	"github.com/ingenuity-build/quicksilver/x/interchainquery"
	interchainquerykeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	interchainquerytypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards"
	participationrewardskeeper "github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".quicksilverd")
}

const (
	// Name defines the application binary name
	Name        = "quicksilverd"
	upgradeName = "underpressure"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		mint.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.LegacyProposalHandler, upgradeclient.LegacyCancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler, ibcclientclient.UpgradeProposalHandler, interchainstakingclient.RegisterProposalHandler,
				interchainstakingclient.UpdateProposalHandler,
			},
		),

		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		ica.AppModuleBasic{},
		vesting.AppModuleBasic{},
		epochs.AppModuleBasic{},
		interchainstaking.AppModuleBasic{},
		interchainquery.AppModuleBasic{},
		participationrewards.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName: nil,
		distrtypes.ModuleName:      nil,
		// TODO: Remove Burner from mint - for dev/test only;
		minttypes.ModuleName:              {authtypes.Minter, authtypes.Burner},
		stakingtypes.BondedPoolName:       {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:    {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:               {authtypes.Burner},
		ibctransfertypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:               nil,
		interchainstakingtypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		interchainquerytypes.ModuleName:   nil,
		ibcfeetypes.ModuleName:            nil,
		// TODO: Remove Burner from participationrewards - for dev/test only;
		participationrewardstypes.ModuleName: {authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName:             true,
		interchainstakingtypes.ModuleName: true,
	}
)

var (
	_ servertypes.Application = (*Quicksilver)(nil)
	_ simapp.App              = (*Quicksilver)(nil)
	_ ibctesting.TestingApp   = (*Quicksilver)(nil)
)

// Quicksilver implements an extended ABCI application.
type Quicksilver struct {
	*baseapp.BaseApp

	// encoding
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper       authkeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	CapabilityKeeper    *capabilitykeeper.Keeper
	StakingKeeper       stakingkeeper.Keeper
	SlashingKeeper      slashingkeeper.Keeper
	DistrKeeper         distrkeeper.Keeper
	MintKeeper          mintkeeper.Keeper
	GovKeeper           govkeeper.Keeper
	CrisisKeeper        crisiskeeper.Keeper
	UpgradeKeeper       upgradekeeper.Keeper
	ParamsKeeper        paramskeeper.Keeper
	FeeGrantKeeper      feegrantkeeper.Keeper
	AuthzKeeper         authzkeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper        ibcfeekeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	EvidenceKeeper      evidencekeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper

	// Quicksilver keepers
	EpochsKeeper               epochskeeper.Keeper
	InterchainstakingKeeper    interchainstakingkeeper.Keeper
	InterchainQueryKeeper      interchainquerykeeper.Keeper
	ParticipationRewardsKeeper participationrewardskeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedIBCFeeKeeper        capabilitykeeper.ScopedKeeper
	ScopedFeeMockKeeper       capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedIBCMockKeeper       capabilitykeeper.ScopedKeeper
	ScopedICAMockKeeper       capabilitykeeper.ScopedKeeper

	// make IBC modules public for test purposes
	// these modules are never directly routed to by the IBC Router
	ICAAuthModule ibcmock.IBCModule
	FeeMockModule ibcmock.IBCModule

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// the configurator
	configurator module.Configurator

	// tps counter
	tpsCounter *tpsCounter
}

// NewQuicksilver returns a reference to a new initialized Quicksilver application.
func NewQuicksilver(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Quicksilver {
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	// NOTE we use custom transaction decoder that supports the sdk.Tx interface instead of sdk.StdTx
	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		// SDK keys
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distrtypes.StoreKey, minttypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, capabilitytypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey,
		// ibc keys
		ibchost.StoreKey, ibctransfertypes.StoreKey,
		icacontrollertypes.StoreKey,
		icahosttypes.StoreKey,
		// quicksilver keys
		epochstypes.StoreKey,
		interchainstakingtypes.StoreKey,
		interchainquerytypes.StoreKey,
		participationrewardstypes.StoreKey,
	)

	// Add the transient store key
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &Quicksilver{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// init params keeper and subspaces
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedInterchainStakingKeeper := app.CapabilityKeeper.ScopeToModule(interchainstakingtypes.ModuleName)
	scopedIBCMockKeeper := app.CapabilityKeeper.ScopeToModule(ibcmock.ModuleName)
	scopedICAMockKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	// use custom account for contracts
	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms, sdk.Bech32MainPrefix,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName,
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		app.DistrKeeper, app.EpochsKeeper, authtypes.FeeCollectorName,
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)

	// Mock Module Stack

	// Mock Module setup for testing IBC and also acts as the interchain accounts authentication module
	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// not replicate if you do not need to test core IBC or light clients.
	mockModule := ibcmock.NewAppModule(&app.IBCKeeper.PortKeeper)

	// The mock module is used for testing IBC
	mockIBCModule := ibcmock.NewIBCModule(&mockModule, ibcmock.NewMockIBCApp(ibcmock.ModuleName, scopedIBCMockKeeper))

	// Create Transfer Stack
	// SendPacket, since it is originating from the application to core IBC:
	// transferKeeper.SendPacket -> fee.SendPacket -> channel.SendPacket

	// RecvPacket, message that originates from core IBC and goes down to app, the flow is the other way
	// channel.RecvPacket -> fee.OnRecvPacket -> transfer.OnRecvPacket

	// transfer stack contains (from top to bottom):
	// - IBC Fee Middleware
	// - Transfer

	// create IBC module from bottom to top of stack
	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, app.IBCFeeKeeper)

	// Add transfer stack to IBC Router

	// Create Interchain Accounts Stack
	// SendPacket, since it is originating from the application to core IBC:
	// icaAuthModuleKeeper.SendTx -> icaController.SendPacket -> fee.SendPacket -> channel.SendPacket

	// initialize ICA module with mock module as the authentication module on the controller side
	var icaControllerStack porttypes.IBCModule
	icaControllerStack = ibcmock.NewIBCModule(&mockModule, ibcmock.NewMockIBCApp("", scopedICAMockKeeper))
	app.ICAAuthModule = icaControllerStack.(ibcmock.IBCModule)
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, app.ICAControllerKeeper)
	icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerStack, app.IBCFeeKeeper)

	// RecvPacket, message that originates from core IBC and goes down to app, the flow is:
	// channel.RecvPacket -> fee.OnRecvPacket -> icaHost.OnRecvPacket

	var icaHostStack porttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(app.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, app.IBCFeeKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostStack).
		AddRoute(interchainstakingtypes.ModuleName, icaControllerStack).
		AddRoute(ibcmock.ModuleName, mockIBCModule).
		AddRoute(ibcmock.ModuleName+icacontrollertypes.SubModuleName, icaControllerStack) // ica with mock auth module stack route to ica (top level of middleware stack)
	app.IBCKeeper.SetRouter(ibcRouter)

	// Add host, controller & ica auth modules to IBC router
	// the ICA Controller middleware needs to be explicitly added to the IBC Router because the
	// ICA controller module owns the port capability for ICA. The ICA authentication module
	// owns the channel capability.

	// ICA Controller keeper
	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec, keys[icacontrollertypes.StoreKey], app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper, app.MsgServiceRouter(),
	)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey], app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, scopedICAHostKeeper, app.MsgServiceRouter(),
	)
	icaModule := ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper)

	app.InterchainQueryKeeper = interchainquerykeeper.NewKeeper(appCodec, keys[interchainquerytypes.StoreKey], app.IBCKeeper)
	interchainQueryModule := interchainquery.NewAppModule(appCodec, app.InterchainQueryKeeper)

	app.InterchainstakingKeeper = interchainstakingkeeper.NewKeeper(
		appCodec,
		keys[interchainstakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.ICAControllerKeeper,
		scopedInterchainStakingKeeper,
		app.InterchainQueryKeeper,
		*app.IBCKeeper,
		app.GetSubspace(interchainstakingtypes.ModuleName),
	)

	interchainstakingModule := interchainstaking.NewAppModule(appCodec, app.InterchainstakingKeeper)

	// interchainstakingIBCModule := interchainstaking.NewIBCModule(app.InterchainstakingKeeper)

	app.InterchainQueryKeeper.SetCallbackHandler(interchainstakingtypes.ModuleName, app.InterchainstakingKeeper.CallbackHandler())

	app.ParticipationRewardsKeeper = participationrewardskeeper.NewKeeper(
		appCodec,
		keys[participationrewardstypes.StoreKey],
		app.GetSubspace(participationrewardstypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.InterchainQueryKeeper,
		app.InterchainstakingKeeper,
		authtypes.FeeCollectorName,
	)
	participationrewardsModule := participationrewards.NewAppModule(appCodec, app.ParticipationRewardsKeeper)

	app.InterchainQueryKeeper.SetCallbackHandler(participationrewardstypes.ModuleName, app.ParticipationRewardsKeeper.CallbackHandler())

	// Quicksilver Keepers
	epochsKeeper := epochskeeper.NewKeeper(appCodec, keys[epochstypes.StoreKey])
	app.EpochsKeeper = *epochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(
			app.MintKeeper.Hooks(),
			app.InterchainstakingKeeper.Hooks(),
			app.ParticipationRewardsKeeper.Hooks(),
		),
	)

	// icaControllerIBCModule := icacontroller.NewIBCMiddleware(interchainstakingIBCModule, app.ICAControllerKeeper)
	// icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(interchainstakingtypes.RouterKey, interchainstaking.NewProposalHandler(app.InterchainstakingKeeper))
	// add custom proposal routes here.

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://github.com/cosmos/cosmos-sdk/blob/release/v0.46.x/x/gov/spec/01_concepts.md#proposal-messages

	govConfig := govtypes.DefaultConfig()
	/*
		Example of setting gov params:
		govConfig.MaxMetadataLen = 10000
	*/
	govKeeper := govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter, app.MsgServiceRouter(), govConfig,
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		// SDK app modules
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, app.BankKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		// ibc modules
		ibc.NewAppModule(app.IBCKeeper),
		// transferModule,
		icaModule,
		// Quicksilver app modules
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		interchainstakingModule,
		interchainQueryModule,
		participationrewardsModule,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: upgrade module must go first to handle software upgrades.
	// NOTE: staking module is required if HistoricalEntries param > 0.
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		// Note: epochs' begin should be "real" start of epochs, we keep epochs beginblock at the beginning
		epochstypes.ModuleName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,
		interchainstakingtypes.ModuleName,
		interchainquerytypes.ModuleName, // check ordering here.
		// no-op modules
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		participationrewardstypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
	)

	// NOTE: fee market module must go last in order to retrieve the block gas used.
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		// Note: epochs' endblock should be "real" end of epochs, we keep epochs endblock at the end
		interchainquerytypes.ModuleName,
		epochstypes.ModuleName,
		// no-op modules
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		interchainstakingtypes.ModuleName,
		participationrewardstypes.ModuleName,
		// currently no-op.
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		// SDK modules
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		icatypes.ModuleName,
		// Quicksilver modules
		epochstypes.ModuleName,
		interchainstakingtypes.ModuleName,
		interchainquerytypes.ModuleName,
		participationrewardstypes.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	// handle upgrades here/
	// app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	// 	return app.mm.RunMigrations(ctx, app.configurator, vm)
	// })
	app.mm.RegisterServices(app.configurator)

	// // add test gRPC service for testing gRPC queries in isolation
	// // testdata.RegisterTestServiceServer(app.GRPCQueryRouter(), testdata.TestServiceImpl{})

	// // create the simulation manager and define the order of the modules for deterministic simulations

	// // NOTE: this is not required apps that don't use the simulator for fuzz testing
	// // transactions
	// app.sm = module.NewSimulationManager(
	// 	auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
	// 	bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
	// 	capability.NewAppModule(appCodec, *app.CapabilityKeeper),
	// 	gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
	// 	mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
	// 	staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
	// 	distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
	// 	slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
	// 	params.NewAppModule(app.ParamsKeeper),
	// 	evidence.NewAppModule(app.EvidenceKeeper),
	// 	feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
	// 	authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
	// 	// ibc modules
	// 	ibc.NewAppModule(app.IBCKeeper),
	// 	transferModule,
	// 	//icaModule,
	// 	// Quicksilver app modules
	// 	epochs.NewAppModule(appCodec, app.EpochsKeeper),
	// 	interchainstakingModule,
	// )

	// app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	options := HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			FeegrantKeeper:  app.FeeGrantKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
		IBCKeeper: app.IBCKeeper,
	}

	app.SetAnteHandler(NewAnteHandler(options))
	app.SetEndBlocker(app.EndBlocker)

	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("nothing to see here o.O")

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedICAControllerKeeper = scopedICAControllerKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper

	// Finally start the tpsCounter.
	app.tpsCounter = newTPSCounter(logger)
	go func() {
		// Unfortunately golangci-lint is so pedantic
		// so we have to ignore this error explicitly.
		_ = app.tpsCounter.start(context.Background())
	}()

	return app
}

// Name returns the name of the App
func (app *Quicksilver) Name() string { return app.BaseApp.Name() }

// BeginBlocker updates every begin block
func (app *Quicksilver) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker updates every end block
func (app *Quicksilver) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// We are intentionally decomposing the DeliverTx method so as to calculate the transactions per second.
func (app *Quicksilver) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	defer func() {
		// TODO: Record the count along with the code and or reason so as to display
		// in the transactions per second live dashboards.
		if res.IsErr() {
			app.tpsCounter.incrementFailure()
		} else {
			app.tpsCounter.incrementSuccess()
		}
	}()

	return app.BaseApp.DeliverTx(req)
}

// InitChainer updates at chain initialization
func (app *Quicksilver) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads state at a particular height
func (app *Quicksilver) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *Quicksilver) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *Quicksilver) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns Quicksilver's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Quicksilver) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Quicksilver's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Quicksilver) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Quicksilver's InterfaceRegistry
func (app *Quicksilver) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Quicksilver) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Quicksilver) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *Quicksilver) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Quicksilver) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *Quicksilver) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *Quicksilver) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

func (app *Quicksilver) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *Quicksilver) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// IBC Go TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *Quicksilver) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
func (app *Quicksilver) GetStakingKeeper() stakingkeeper.Keeper {
	return app.StakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *Quicksilver) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *Quicksilver) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *Quicksilver) GetTxConfig() client.TxConfig {
	cfg := MakeEncodingConfig()
	return cfg.TxConfig
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(
	appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// SDK subspaces
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	// ibc subspaces
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	// quicksilver subspaces
	paramsKeeper.Subspace(interchainstakingtypes.ModuleName)
	paramsKeeper.Subspace(interchainquerytypes.ModuleName)
	paramsKeeper.Subspace(participationrewardstypes.ModuleName)
	return paramsKeeper
}
