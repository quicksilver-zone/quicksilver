package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	sdksimtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/ingenuity-build/quicksilver/osmosis-types/osmoutils"
	"github.com/ingenuity-build/quicksilver/test/simulation/simtypes"
	"github.com/ingenuity-build/quicksilver/x/tokenfactory/keeper"
	"github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

const (
	OpWeightMsgCreateDenom               = "op_weight_msg_create_denom"       //nolint:gosec // not credentials
	OpWeightMsgMint                      = "op_weight_msg_mint"               //nolint:gosec // not credentials
	OpWeightMsgBurn                      = "op_weight_msg_burn"               //nolint:gosec // not credentials
	OpWeightMsgChangeAdmin               = "op_weight_msg_change_admin"       //nolint:gosec // not credentials
	OpWeightMsgSetDenomMetadata          = "op_weight_msg_set_denom_metadata" //nolint:gosec // not credentials
	DefaultWeightMsgCreateDenom      int = 50
	DefaultWeightMsgMint             int = 10
	DefaultWeightMsgBurn             int = 10
	DefaultWeightMsgChangeAdmin      int = 5
	DefaultWeightMsgSetDenomMetadata int = 5
)

var (
	TypeMsgCreateDenom      = sdk.MsgTypeURL(&types.MsgCreateDenom{})
	TypeMsgMint             = sdk.MsgTypeURL(&types.MsgMint{})
	TypeMsgBurn             = sdk.MsgTypeURL(&types.MsgBurn{})
	TypeMsgChangeAdmin      = sdk.MsgTypeURL(&types.MsgChangeAdmin{})
	TypeMsgSetDenomMetadata = sdk.MsgTypeURL(&types.MsgSetDenomMetadata{})
)

func WeightedOperations(
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
			SimulateMsgCreateDenom(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgMint,
			SimulateMsgMint(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBurn,
			SimulateMsgBurn(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgChangeAdmin,
			SimulateMsgChangeAdmin(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSetDenomMetadata,
			SimulateMsgSetDenomMetadata(ak, bk, k),
		),
	}
}

// SimulateMsgCreateDenom generates a MsgCreateDenom with random values.
func SimulateMsgCreateDenom(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
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
			TxGen:           moduletestutil.MakeTestTxConfig(),
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgCreateDenom,
			CoinsSpentInMsg: minCoins,
			Context:         ctx,
			SimAccount:      acc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgMint generates a MsgMint with random values.
func SimulateMsgMint(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, bApp *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		acc, senderExists := simtypes.RandomSimAccountWithConstraint(r, accountCreatedTokenFactoryDenom(k, ctx), accs)
		if !senderExists {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgMint, "no account with tokenfactory denom found"), nil, nil
		}

		denom, addr, err := getTokenFactoryDenomAndItsAdmin(k, ctx, r, acc)
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgMint, "error finding denom and admin"), nil, err
		}
		if addr == nil || addr.String() != acc.Address.String() {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgMint, "account is not admin"), nil, nil
		}

		mintAmount, err := simtypes.RandPositiveInt(r, sdk.NewIntFromUint64(1000_000000))
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgMint, "error creating random sdkmath.Int"), nil, err
		}

		msg := &types.MsgMint{
			Sender: addr.String(),
			Amount: sdk.NewCoin(denom, mintAmount),
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             bApp,
			TxGen:           moduletestutil.MakeTestTxConfig(),
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgMint,
			CoinsSpentInMsg: nil,
			Context:         ctx,
			SimAccount:      acc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgBurn generates a MsgBurn with random values.
func SimulateMsgBurn(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, bApp *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		acc, senderExists := simtypes.RandomSimAccountWithConstraint(r, accountCreatedTokenFactoryDenom(k, ctx), accs)
		if !senderExists {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgBurn, "no account with tokenfactory denom found"), nil, nil
		}

		denom, addr, err := getTokenFactoryDenomAndItsAdmin(k, ctx, r, acc)
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgBurn, "error finding denom and admin"), nil, err
		}
		if addr == nil || addr.String() != acc.Address.String() {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgBurn, "account is not admin"), nil, nil
		}

		spendable := bk.SpendableCoins(ctx, addr)
		denomBal := spendable.AmountOf(denom)

		if denomBal.IsZero() {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgBurn, "addr does not have enough balance to burn"), nil, nil
		}

		burnAmount, err := simtypes.RandPositiveInt(r, denomBal)
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgMint, "error creating random sdkmath.Int"), nil, err
		}

		burnCoin := sdk.NewCoin(denom, burnAmount)
		msg := &types.MsgBurn{
			Sender: addr.String(),
			Amount: burnCoin,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             bApp,
			TxGen:           moduletestutil.MakeTestTxConfig(),
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgBurn,
			CoinsSpentInMsg: sdk.NewCoins(burnCoin),
			Context:         ctx,
			SimAccount:      acc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgChangeAdmin generates a MsgChangeAdmin with random values.
func SimulateMsgChangeAdmin(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, bApp *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		acc, senderExists := simtypes.RandomSimAccountWithConstraint(r, accountCreatedTokenFactoryDenom(k, ctx), accs)
		if !senderExists {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "no account with tokenfactory denom found"), nil, nil
		}

		denom, addr, err := getTokenFactoryDenomAndItsAdmin(k, ctx, r, acc)
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "error finding denom and admin"), nil, err
		}
		if addr == nil || addr.String() != acc.Address.String() {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "account is not admin"), nil, nil
		}

		newAdmin := simtypes.RandomSimAccount(r, accs)
		if newAdmin.Address.String() == addr.String() {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "denom has no admin"), nil, nil
		}

		msg := &types.MsgChangeAdmin{
			Sender:   addr.String(),
			Denom:    denom,
			NewAdmin: newAdmin.Address.String(),
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             bApp,
			TxGen:           moduletestutil.MakeTestTxConfig(),
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgChangeAdmin,
			CoinsSpentInMsg: nil,
			Context:         ctx,
			SimAccount:      acc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgSetDenomMetadata generates a MsgSetDenomMetadata with random values.
func SimulateMsgSetDenomMetadata(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) sdksimtypes.Operation {
	return func(
		r *rand.Rand, bApp *baseapp.BaseApp, ctx sdk.Context, accs []sdksimtypes.Account, chainID string,
	) (sdksimtypes.OperationMsg, []sdksimtypes.FutureOperation, error) {
		acc, senderExists := simtypes.RandomSimAccountWithConstraint(r, accountCreatedTokenFactoryDenom(k, ctx), accs)
		if !senderExists {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "no account with tokenfactory denom found"), nil, nil
		}

		denom, addr, err := getTokenFactoryDenomAndItsAdmin(k, ctx, r, acc)
		if err != nil {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "error finding denom and admin"), nil, err
		}
		if addr == nil || addr.String() != acc.Address.String() {
			return sdksimtypes.NoOpMsg(types.ModuleName, TypeMsgChangeAdmin, "account is not admin"), nil, nil
		}

		msg := &types.MsgSetDenomMetadata{
			Sender: addr.String(),
			Metadata: banktypes.Metadata{
				Description: simtypes.RandStringOfLength(r, 10),
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    denom,
						Exponent: 0,
						Aliases: []string{
							simtypes.RandStringOfLength(r, 4),
							simtypes.RandStringOfLength(r, 4),
						},
					},
				},
				Base:    denom,
				Display: denom,
				Name:    simtypes.RandStringOfLength(r, 10),
				Symbol:  simtypes.RandStringOfLength(r, 4),
				URI:     "",
				URIHash: "",
			},
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             bApp,
			TxGen:           moduletestutil.MakeTestTxConfig(),
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgSetDenomMetadata,
			CoinsSpentInMsg: nil,
			Context:         ctx,
			SimAccount:      acc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func accountCreatedTokenFactoryDenom(k keeper.Keeper, ctx sdk.Context) simtypes.SimAccountConstraint {
	return func(acc sdksimtypes.Account) bool {
		store := k.GetCreatorPrefixStore(ctx, acc.Address.String())
		iterator := store.Iterator(nil, nil)
		defer iterator.Close()
		return iterator.Valid()
	}
}

func getTokenFactoryDenomAndItsAdmin(k keeper.Keeper, ctx sdk.Context, r *rand.Rand, acc sdksimtypes.Account) (string, sdk.AccAddress, error) {
	store := k.GetCreatorPrefixStore(ctx, acc.Address.String())
	denoms := osmoutils.GatherAllKeysFromStore(store)
	denom := simtypes.RandSelect(r, denoms...)

	authData, err := k.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return "", nil, err
	}
	admin := authData.Admin
	addr, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return "", nil, err
	}
	return denom, addr, nil
}
