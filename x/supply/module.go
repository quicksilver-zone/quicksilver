package supply

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/quicksilver-zone/quicksilver/x/supply/keeper"
	"github.com/quicksilver-zone/quicksilver/x/supply/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the capability module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the participationrewards module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers a legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
}

// RegisterInterfaces registers the module's interface types
func (AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
}

// DefaultGenesis returns the participationrewards module's default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the participationrewards module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

// GetTxCmd returns the participationrewards module's root tx command.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the participationrewards module's root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the participationrewards module.
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule return a new AppModule
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		keeper: keeper,
	}
}

// RegisterInvariants registers the participationrewards module invariants.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// RegisterServices registers a gRPC query service to respond to the
// module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(am.keeper))
}

// InitGenesis performs the participationrewards module's genesis
// initialization It returns no validator updates.
func (AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the participationrewards module's exported genesis state as raw JSON bytes.
func (AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return nil
}

// BeginBlock executes all ABCI BeginBlock logic respective to the participationrewards module.
func (AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock executes all ABCI EndBlock logic respective to the participationrewards module. It
// returns no validator updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }
