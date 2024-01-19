package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/std"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/x/tx/signing"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"

	"cosmossdk.io/log"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	tmos "github.com/cometbft/cometbft/libs/os"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	"github.com/spf13/cast"

	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctestingtypes "github.com/cosmos/ibc-go/v8/testing/types"

	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/quicksilver-zone/quicksilver/v7/app/keepers"
	"github.com/quicksilver-zone/quicksilver/v7/docs"
	airdroptypes "github.com/quicksilver-zone/quicksilver/v7/x/airdrop/types"
	interchainstakingtypes "github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
)

func Init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultNodeHome = filepath.Join(userHomeDir, ".quicksilverd")
}

const (
	// Name defines the application binary name.
	Name         = "quicksilverd"
	Bech32Prefix = "quicksilver"
)

// These constants are derived from the above variables.
// These are the ones we will want to use in the code, based on
// any overrides above

var (
	// DefaultNodeHome default home directories for the application daemon.
	DefaultNodeHome string

	// module accounts that are allowed to receive tokens.
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName:             true,
		interchainstakingtypes.ModuleName: true,
		airdroptypes.ModuleName:           true,
	}
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

var (
	_ runtime.AppI            = (*Quicksilver)(nil)
	_ servertypes.Application = (*Quicksilver)(nil)
)

// Quicksilver implements an extended ABCI application.
type Quicksilver struct {
	*baseapp.BaseApp
	keepers.AppKeepers
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// the module manager
	mm                 *module.Manager
	BasicModuleManager module.BasicManager

	// simulation manager
	sm *module.SimulationManager

	// the configurator
	configurator module.Configurator

	tpsCounter *tpsCounter
	once       sync.Once
}

// NewQuicksilver returns a reference to a new initialized Quicksilver application.
func NewQuicksilver(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	appOpts servertypes.AppOptions,
	mock bool,
	enableSupplyEndpoint bool,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Quicksilver {
	interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	appCodec := codec.NewProtoCodec(interfaceRegistry)
	cdc := codec.NewLegacyAmino()
	txConfig := authtx.NewTxConfig(appCodec, authtx.DefaultSignModes)

	std.RegisterLegacyAminoCodec(cdc)
	std.RegisterInterfaces(interfaceRegistry)
	// NOTE we use custom transaction decoder that supports the sdk.Tx interface instead of sdk.StdTx
	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		txConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	app := &Quicksilver{
		BaseApp:           bApp,
		legacyAmino:       cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		txConfig:          txConfig,
	}

	wasmDir := filepath.Join(homePath, "data")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	app.AppKeepers = keepers.NewAppKeepers(
		appCodec,
		bApp,
		cdc,
		maccPerms,
		app.BlockedAddrs(),
		skipUpgradeHeights,
		mock,
		homePath,
		appOpts,
		wasmDir,
		enableSupplyEndpoint,
	)

	// ****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))
	app.mm = module.NewManager(
		append(
			appModules(app, appCodec, skipGenesisInvariants),
			genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app, txConfig),
		)...)

	// BasicModuleManager defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	// By default it is composed of all the module from the module manager.
	// Additionally, app module basics can be overwritten by passing them as argument.
	app.BasicModuleManager = module.NewBasicManagerFromManager(
		app.mm,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
		})
	app.BasicModuleManager.RegisterLegacyAminoCodec(cdc)
	app.BasicModuleManager.RegisterInterfaces(interfaceRegistry)

	// NOTE: upgrade module is required to be prioritized
	app.mm.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)
	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.mm.SetOrderEndBlockers(orderEndBlockers()...)
	app.mm.SetOrderInitGenesis(orderInitBlockers()...)
	app.mm.SetOrderExportGenesis(orderInitBlockers()...)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	// TODO: in this commit they just removed the RegisterRoutes method https://github.com/cosmos/cosmos-sdk/commit/3a097012b59413641ac92f18f226c5d6b674ae42

	// app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	// // add test gRPC service for testing gRPC queries in isolation
	// // testdata.RegisterTestServiceServer(app.GRPCQueryRouter(), testdata.TestServiceImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(simulationModules(app, appCodec)...)
	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	// initialize BaseApp
	options := HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			FeegrantKeeper:  app.FeeGrantKeeper,
			SignModeHandler: txConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
		WasmConfig:        wasmConfig,
		TxCounterStoreKey: runtime.NewKVStoreService(app.AppKeepers.GetKey(wasm.StoreKey)),
		IBCKeeper:         app.IBCKeeper,
	}

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(NewAnteHandler(options))
	app.SetEndBlocker(app.EndBlocker)

	// handle upgrades here
	app.setUpgradeHandlers()
	app.setUpgradeStoreLoaders()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	// Finally start the tpsCounter.
	app.tpsCounter = newTPSCounter(logger)
	go func() {
		// Unfortunately golangci-lint is so pedantic
		// so we have to ignore this error explicitly.
		_ = app.tpsCounter.start(context.Background()) // nolint:errcheck
	}()

	return app
}

// Name returns the name of the App.
func (app *Quicksilver) Name() string { return app.BaseApp.Name() }

// BeginBlocker updates every begin block.
func (app *Quicksilver) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	if ctx.ChainID() == "quicksilver-2" && ctx.BlockHeight() == 235001 {
		zone, found := app.InterchainstakingKeeper.GetZone(ctx, "stargaze-1")
		if !found {
			panic("ERROR: unable to find expected stargaze-1 zone")
		}
		app.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, &zone)
	}

	return app.mm.BeginBlock(ctx)
}

// EndBlocker updates every end block.
func (app *Quicksilver) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// // TODO: Figure out how to reimplement this, new version of cosmos-sdk doesn't have DeliverTx exposed
// ref:https://github.com/cosmos/cosmos-sdk/blob/release/v0.50.x/UPGRADING.md#baseapp

// DeliverTx calls BaseApp.DeliverTx and calculates transactions per second.
// func (app *Quicksilver) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
// 	defer func() {
// 		// TODO: Record the count along with the code and or reason so as to display
// 		// in the transactions per second live dashboards.
// 		if res.IsErr() {
// 			app.tpsCounter.incrementFailure()
// 		} else {
// 			app.tpsCounter.incrementSuccess()
// 		}
// 	}()

// 	return app.BaseApp.DeliverTx(req)
// }

// InitChainer application update at chain initialization
func (app *Quicksilver) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	if err != nil {
		panic(err)
	}
	response, err := app.mm.InitGenesis(ctx, app.appCodec, genesisState)
	return response, err
}

// LoadHeight loads state at a particular height.
func (app *Quicksilver) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (*Quicksilver) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (*Quicksilver) BlockedAddrs() map[string]bool {
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
	return app.legacyAmino
}

// AppCodec returns Quicksilver's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Quicksilver) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Quicksilver's InterfaceRegistry.
func (app *Quicksilver) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Quicksilver) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface.
func (app *Quicksilver) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (*Quicksilver) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API
	if apiConfig.Swagger {
		apiSvr.Router.Handle("/swagger.yml", http.FileServer(http.FS(docs.Swagger)))
		apiSvr.Router.HandleFunc("/", docs.Handler(Name, "/swagger.yml"))
	}
}

func (app *Quicksilver) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *Quicksilver) RegisterTendermintService(clientCtx client.Context) {
	cmtApp := server.NewCometABCIWrapper(app)

	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		cmtApp.Query,
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
	return app.txConfig
}

// GetMaccPerms returns a copy of the module account permissions.
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

func (app *Quicksilver) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// AutoCliOpts returns the autocli options for the app.
func (app *Quicksilver) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(app.mm.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

func (app *Quicksilver) FinalizeBlock(req *abci.RequestFinalizeBlock) (*abci.ResponseFinalizeBlock, error) {
	// when skipping sdk 47 for sdk 50, the upgrade handler is called too late in BaseApp
	// this is a hack to ensure that the migration is executed when needed and not panics
	app.once.Do(func() {
		ctx := app.NewUncachedContext(false, tmproto.Header{})
		if _, err := app.ConsensusParamsKeeper.Params(ctx, &consensusparamtypes.QueryParamsRequest{}); err != nil {
			// prevents panic: consensus key is nil: collections: not found: key 'no_key' of type github.com/cosmos/gogoproto/tendermint.types.ConsensusParams
			// sdk 47:
			// Migrate Tendermint consensus parameters from x/params module to a dedicated x/consensus module.
			// see https://github.com/cosmos/cosmos-sdk/blob/v0.47.0/simapp/upgrades.go#L66
			baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace)
			err := baseapp.MigrateParams(sdk.UnwrapSDKContext(ctx), baseAppLegacySS, app.ConsensusParamsKeeper.ParamsStore)
			if err != nil {
				panic(err)
			}
		}
	})

	return app.BaseApp.FinalizeBlock(req)
}
