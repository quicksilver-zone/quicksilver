package ibc_hooks

import (
	"encoding/json"

	"math/rand"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ingenuity-build/quicksilver/x/ibc-hooks/simulation"
	"github.com/spf13/cobra"

	"github.com/ingenuity-build/quicksilver/x/ibc-hooks/client/cli"
	"github.com/ingenuity-build/quicksilver/x/ibc-hooks/types"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/ingenuity-build/quicksilver/third-party-chains/osmosis-types/osmoutils"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the ibc-hooks module.
type AppModuleBasic struct {
	cdc codec.Codec
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the ibc-hooks module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the ibc-hooks module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers the module's interface types.
func (b AppModuleBasic) RegisterInterfaces(_ cdctypes.InterfaceRegistry) {}

// DefaultGenesis returns default genesis state as raw bytes for the
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	emptyString := "{}"
	return []byte(emptyString)
}

// ValidateGenesis performs genesis state validation for the ibc-hooks module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the ibc-hooks module.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the ibc-hooks module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

// GetTxCmd returns no root tx command for the ibc-hooks module.
func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

// GetQueryCmd returns the root query command for the ibc-hooks module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// ___________________________________________________________________________

// AppModule implements an application module for the ibc-hooks module.
type AppModule struct {
	AppModuleBasic

	authKeeper osmoutils.AccountKeeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(cdc codec.Codec, ak osmoutils.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		authKeeper:     ak,
	}
}

// NewAppModuleBasic return a new AppModuleBasic.
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the ibc-hooks module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the ibc-hooks module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// RegisterServices registers a gRPC query service to respond to the
// module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
}

// InitGenesis performs genesis initialization for the ibc-hooks module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return json.RawMessage([]byte("{}"))
}

// BeginBlock returns the begin blocker for the ibc-hooks module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the ibc-hooks module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the pool-incentives module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RandomizedParams creates randomized mint param changes for the simulator.
func (AppModule) RandomizedParams(_ *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{}
}

// RegisterStoreDecoder registers a decoder for supply module's types.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil // TODO add
}
