package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
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

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/wasmbinding"

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

	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ingenuity-build/quicksilver/x/mint"
	mintkeeper "github.com/ingenuity-build/quicksilver/x/mint/keeper"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
	"github.com/ingenuity-build/quicksilver/x/tokenfactory"

	ica "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v5/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v5/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v5/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v5/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v5/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v5/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"

	ibctestingtypes "github.com/cosmos/ibc-go/v5/testing/types"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager"
	claimsmanagerkeeper "github.com/ingenuity-build/quicksilver/x/claimsmanager/keeper"
	claimsmanagertypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"

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
	participationrewardsclient "github.com/ingenuity-build/quicksilver/x/participationrewards/client"
	participationrewardskeeper "github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop"
	airdropkeeper "github.com/ingenuity-build/quicksilver/x/airdrop/keeper"
	airdroptypes "github.com/ingenuity-build/quicksilver/x/airdrop/types"

	tokenfactorykeeper "github.com/ingenuity-build/quicksilver/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"

	"github.com/CosmWasm/wasmd/x/wasm"
)

func Init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultNodeHome = filepath.Join(userHomeDir, ".quicksilverd")
}

const (
	// Name defines the application binary name
	Name = "quicksilverd"
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
				ibcclientclient.UpdateClientProposalHandler, ibcclientclient.UpgradeProposalHandler, interchainstakingclient.RegisterProposalHandler, interchainstakingclient.UpdateProposalHandler,
				participationrewardsclient.AddProtocolDataProposalHandler,
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
		claimsmanager.AppModuleBasic{},
		epochs.AppModuleBasic{},
		interchainstaking.AppModuleBasic{},
		interchainquery.AppModuleBasic{},
		participationrewards.AppModuleBasic{},
		airdrop.AppModuleBasic{},
		tokenfactory.AppModuleBasic{},
		wasm.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:                 nil,
		distrtypes.ModuleName:                      nil,
		minttypes.ModuleName:                       {authtypes.Minter},
		stakingtypes.BondedPoolName:                {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:             {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:                        {authtypes.Burner},
		ibctransfertypes.ModuleName:                {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:                        nil,
		claimsmanagertypes.ModuleName:              nil,
		interchainstakingtypes.ModuleName:          {authtypes.Minter},
		interchainstakingtypes.EscrowModuleAccount: {authtypes.Burner},
		interchainquerytypes.ModuleName:            nil,
		participationrewardstypes.ModuleName:       nil,
		airdroptypes.ModuleName:                    nil,
		wasm.ModuleName:                            {authtypes.Burner},
		tokenfactorytypes.ModuleName:               {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName:             true,
		interchainstakingtypes.ModuleName: true,
		airdroptypes.ModuleName:           true,
	}
)

var (
	_ servertypes.Application = (*Quicksilver)(nil)
	_ simapp.App              = (*Quicksilver)(nil)
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
	BankKeeper          bankkeeper.BaseKeeper
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
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	EvidenceKeeper      evidencekeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper

	// Quicksilver keepers
	EpochsKeeper               epochskeeper.Keeper
	ClaimsManagerKeeper        claimsmanagerkeeper.Keeper
	InterchainstakingKeeper    interchainstakingkeeper.Keeper
	InterchainQueryKeeper      interchainquerykeeper.Keeper
	ParticipationRewardsKeeper participationrewardskeeper.Keeper
	AirdropKeeper              airdropkeeper.Keeper
	TokenFactoryKeeper         *tokenfactorykeeper.Keeper
	// make scoped keepers public for test purposes
	ScopedIBCKeeper                      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper                 capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper            capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper                  capabilitykeeper.ScopedKeeper
	ScopedInterchainStakingAccountKeeper capabilitykeeper.ScopedKeeper

	WasmKeeper       wasm.Keeper
	scopedWasmKeeper capabilitykeeper.ScopedKeeper //nolint:unused //TODO: we can use this for testing

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// the configurator
	configurator module.Configurator

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
	enabledProposals []wasm.ProposalType,
	appOpts servertypes.AppOptions,
	wasmOpts []wasm.Option,
	mock bool,
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
		claimsmanagertypes.StoreKey,
		epochstypes.StoreKey,
		interchainstakingtypes.StoreKey,
		interchainquerytypes.StoreKey,
		participationrewardstypes.StoreKey,
		airdroptypes.StoreKey,
		wasm.StoreKey,
		tokenfactorytypes.StoreKey,
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
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	proofOpsFn := utils.ValidateProofOps
	if mock {
		proofOpsFn = utils.MockProofOps
	}

	// use custom account for contracts
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
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	// ICA Keepers
	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec, keys[icacontrollertypes.StoreKey], app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper, app.MsgServiceRouter(),
	)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey], app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, scopedICAHostKeeper, app.MsgServiceRouter(),
	)

	icaModule := ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper)

	app.ClaimsManagerKeeper = claimsmanagerkeeper.NewKeeper(
		appCodec,
		keys[claimsmanagertypes.StoreKey],
	)

	claimsmanagerModule := claimsmanager.NewAppModule(appCodec, app.ClaimsManagerKeeper)

	app.InterchainQueryKeeper = interchainquerykeeper.NewKeeper(appCodec, keys[interchainquerytypes.StoreKey], app.IBCKeeper)
	interchainQueryModule := interchainquery.NewAppModule(appCodec, app.InterchainQueryKeeper)

	app.InterchainstakingKeeper = interchainstakingkeeper.NewKeeper(
		appCodec,
		keys[interchainstakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.ICAControllerKeeper,
		&scopedInterchainStakingKeeper,
		app.InterchainQueryKeeper,
		*app.IBCKeeper,
		app.TransferKeeper,
		app.ClaimsManagerKeeper,
		app.GetSubspace(interchainstakingtypes.ModuleName),
	)
	interchainstakingModule := interchainstaking.NewAppModule(appCodec, app.InterchainstakingKeeper)

	interchainstakingIBCModule := interchainstaking.NewIBCModule(app.InterchainstakingKeeper)

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
		proofOpsFn,
	)
	if err := app.InterchainQueryKeeper.SetCallbackHandler(interchainstakingtypes.ModuleName, app.InterchainstakingKeeper.CallbackHandler()); err != nil {
		panic(err)
	}

	participationrewardsModule := participationrewards.NewAppModule(appCodec, app.ParticipationRewardsKeeper)

	if err := app.InterchainQueryKeeper.SetCallbackHandler(participationrewardstypes.ModuleName, app.ParticipationRewardsKeeper.CallbackHandler()); err != nil {
		panic(err)
	}

	tokenFactoryKeeper := tokenfactorykeeper.NewKeeper(
		keys[tokenfactorytypes.StoreKey],
		app.GetSubspace(tokenfactorytypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper.WithMintCoinsRestriction(tokenfactorytypes.NewTokenFactoryDenomMintCoinsRestriction()),
		app.DistrKeeper,
	)

	app.TokenFactoryKeeper = &tokenFactoryKeeper
	// Quicksilver Keepers
	epochsKeeper := epochskeeper.NewKeeper(appCodec, keys[epochstypes.StoreKey])
	app.EpochsKeeper = *epochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(
			app.MintKeeper.Hooks(),
			app.InterchainstakingKeeper.Hooks(),
			app.ParticipationRewardsKeeper.Hooks(),
		),
	)

	app.ParticipationRewardsKeeper.SetEpochsKeeper(app.EpochsKeeper)

	// this line is used by starport scaffolding # stargate/app/keeperDefinition
	wasmDir := filepath.Join(homePath, "data")

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate,osmosis"
	wasmOpts = append(wasmbinding.RegisterCustomPlugins(&app.BankKeeper, app.TokenFactoryKeeper), wasmOpts...)
	wasmOpts = append(wasmbinding.RegisterStargateQueries(*bApp.GRPCQueryRouter(), appCodec), wasmOpts...)

	app.WasmKeeper = wasm.NewKeeper(
		appCodec,
		keys[wasm.StoreKey],
		app.GetSubspace(wasm.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)

	icaControllerIBCModule := icacontroller.NewIBCMiddleware(interchainstakingIBCModule, app.ICAControllerKeeper)
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
		AddRoute(wasm.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper)).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(interchainstakingtypes.ModuleName, icaControllerIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	govConfig := govtypes.DefaultConfig()

	// register the proposal types
	govRouter := govv1beta1.NewRouter()

	// The gov proposal types can be individually enabled
	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, enabledProposals))
	}

	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(interchainstakingtypes.RouterKey, interchainstaking.NewProposalHandler(app.InterchainstakingKeeper)).
		AddRoute(participationrewardstypes.RouterKey, participationrewards.NewProposalHandler(app.ParticipationRewardsKeeper))
	// add custom proposal routes here.

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter, app.MsgServiceRouter(), govConfig,
	)

	app.AirdropKeeper = airdropkeeper.NewKeeper(
		appCodec,
		keys[airdroptypes.StoreKey],
		app.GetSubspace(airdroptypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.GovKeeper,
		app.InterchainstakingKeeper,
		app.InterchainQueryKeeper,
		app.ParticipationRewardsKeeper,
		proofOpsFn,
	)
	airdropModule := airdrop.NewAppModule(appCodec, app.AirdropKeeper)

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
		transferModule,
		icaModule,
		// Quicksilver app modules
		claimsmanagerModule,
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		interchainstakingModule,
		interchainQueryModule,
		participationrewardsModule,
		airdropModule,
		tokenfactory.NewAppModule(*app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
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
		claimsmanagertypes.ModuleName,
		participationrewardstypes.ModuleName,
		airdroptypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		tokenfactorytypes.ModuleName,
		wasm.ModuleName,
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
		claimsmanagertypes.ModuleName,
		interchainstakingtypes.ModuleName,
		participationrewardstypes.ModuleName,
		airdroptypes.ModuleName,
		tokenfactorytypes.ModuleName,
		wasm.ModuleName,
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
		claimsmanagertypes.ModuleName,
		interchainstakingtypes.ModuleName,
		interchainquerytypes.ModuleName,
		participationrewardstypes.ModuleName,
		airdroptypes.ModuleName,
		tokenfactorytypes.ModuleName,
		// wasmd
		wasm.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	// handle upgrades here
	setUpgradeHandlers(app)

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
		WasmConfig:        wasmConfig,
		TxCounterStoreKey: keys[wasm.StoreKey],
		IBCKeeper:         app.IBCKeeper,
	}

	app.SetAnteHandler(NewAnteHandler(options))
	app.SetEndBlocker(app.EndBlocker)

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
	if ctx.ChainID() == "quicksilver-2" && ctx.BlockHeight() == 235001 {
		zone, found := app.InterchainstakingKeeper.GetZone(ctx, "stargaze-1")
		if !found {
			panic("ERROR: unable to find expected stargaze-1 zone")
		}
		app.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, zone)
	}

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
func (app *Quicksilver) GetStakingKeeper() ibctestingtypes.StakingKeeper {
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
	paramsKeeper.Subspace(claimsmanagertypes.ModuleName)
	paramsKeeper.Subspace(interchainstakingtypes.ModuleName)
	paramsKeeper.Subspace(interchainquerytypes.ModuleName)
	paramsKeeper.Subspace(participationrewardstypes.ModuleName)
	paramsKeeper.Subspace(airdroptypes.ModuleName)
	paramsKeeper.Subspace(tokenfactorytypes.ModuleName)
	// wasm subspace
	paramsKeeper.Subspace(wasm.ModuleName)

	return paramsKeeper
}

func GetWasmOpts(appOpts servertypes.AppOptions) []wasm.Option {
	var wasmOpts []wasm.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	wasmOpts = append(wasmOpts, wasmkeeper.WithGasRegister(NewWasmGasRegister()))

	return wasmOpts
}

const (

	// PROPOSALS
	ProposalsEnabled = "true"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificProposals = ""
)

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}
