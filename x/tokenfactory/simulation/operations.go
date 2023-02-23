package simulation

import (
	"github.com/ingenuity-build/quicksilver/app"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdksimtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/ingenuity-build/quicksilver/simulation/simtypes"
	"github.com/ingenuity-build/quicksilver/x/tokenfactory/keeper"
	"github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

const (
	OpWeightMsgCreateDenom               = "op_weight_msg_create_denom"
	OpWeightMsgMint                      = "op_weight_msg_mint"
	OpWeightMsgBurn                      = "op_weight_msg_burn"
	OpWeightMsgChangeAdmin               = "op_weight_msg_change_admin"
	OpWeightMsgSetDenomMetadata          = "op_weight_msg_set_denom_metadata"
	DefaultWeightMsgCreateDenom      int = 100
	DefaultWeightMsgMint             int = 100
	DefaultWeightMsgBurn             int = 100
	DefaultWeightMsgChangeAdmin      int = 100
	DefaultWeightMsgSetDenomMetadata int = 100
)

var (
	TypeMsgCreateDenom      = sdk.MsgTypeURL(&types.MsgCreateDenom{})
	TypeMsgMint             = sdk.MsgTypeURL(&types.MsgMint{})
	TypeMsgBurn             = sdk.MsgTypeURL(&types.MsgBurn{})
	TypeMsgChangeAdmin      = sdk.MsgTypeURL(&types.MsgChangeAdmin{})
	TypeMsgSetDenomMetadata = sdk.MsgTypeURL(&types.MsgSetDenomMetadata{})
)

func WeightedOperations(
	registry codectypes.InterfaceRegistry,
	appParams sdksimtypes.AppParams,
	cdc codec.JSONCodec,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreateDenom      int
		weightMsgMint             int
		weightMsgBurn             int
		weightMsgChangeAdmin      int
		weightMsgSetDenomMetadata int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateDenom, &weightMsgCreateDenom, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDenom = DefaultWeightMsgCreateDenom
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgMint, &weightMsgMint, nil,
		func(_ *rand.Rand) {
			weightMsgMint = DefaultWeightMsgMint
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgBurn, &weightMsgBurn, nil,
		func(_ *rand.Rand) {
			weightMsgBurn = DefaultWeightMsgBurn
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgChangeAdmin, &weightMsgChangeAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgChangeAdmin = DefaultWeightMsgChangeAdmin
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgSetDenomMetadata, &weightMsgSetDenomMetadata, nil,
		func(_ *rand.Rand) {
			weightMsgSetDenomMetadata = DefaultWeightMsgSetDenomMetadata
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateDenom,
			SimulateMsgCreateDenom(codec.NewProtoCodec(registry), ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgMint,
			SimulateMsgMint(codec.NewProtoCodec(registry), ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBurn,
			SimulateMsgBurn(codec.NewProtoCodec(registry), ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgChangeAdmin,
			SimulateMsgChangeAdmin(codec.NewProtoCodec(registry), ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSetDenomMetadata,
			SimulateMsgSetDenomMetadata(codec.NewProtoCodec(registry), ak, bk, k),
		),
	}
}

// SimulateMsgCreateDenom generates a MsgCreateDenom with random values.
func SimulateMsgCreateDenom(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, bApp *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		minCoins := k.GetParams(ctx).DenomCreationFee
		acc, err := simtypes.RandomSimAccountWithMinCoins(ctx, r, accs, minCoins, bk)
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgCreateDenom, "no account with balance found"), nil, nil

		}

		msg := &types.MsgCreateDenom{
			Sender:   acc.Address.String(),
			Subdenom: simtypes.RandStringOfLength(r, types.MaxSubdenomLength),
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             bApp,
			TxGen:           app.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgCreateDenom,
			CoinsSpentInMsg: minCoins,
			Context:         sdk.Context{},
			SimAccount:      sdksimtypes.Account{},
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgMint generates a MsgMint with random values.
func SimulateMsgMint(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgMint, "TODO"), nil, nil

	}
}

// SimulateMsgBurn generates a MsgBurn with random values.
func SimulateMsgBurn(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgBurn, "TODO"), nil, nil
	}
}

// SimulateMsgChangeAdmin generates a MsgChangeAdmin with random values.
func SimulateMsgChangeAdmin(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "TODO"), nil, nil
	}
}

// SimulateMsgSetDenomMetadata generates a MsgSetDenomMetadata with random values.
func SimulateMsgSetDenomMetadata(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgSetDenomMetadata, "TODO"), nil, nil
	}
}
