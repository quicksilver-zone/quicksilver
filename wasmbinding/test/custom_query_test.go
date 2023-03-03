package wasmbinding

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/wasmbinding/bindings"
)

// we must pay this many uosmo for every pool we create
var poolFee int64 = 1000000000

var defaultFunds = sdk.NewCoins(
	sdk.NewInt64Coin("qck", 333000000),
	sdk.NewInt64Coin("umai", 555000000+2*poolFee),
	sdk.NewInt64Coin("uck", 999000000),
)

func SetupCustomApp(t *testing.T, addr sdk.AccAddress) (*app.Quicksilver, sdk.Context) {
	quicksilverApp, ctx := CreateTestInput(t)
	wasmKeeper := quicksilverApp.WasmKeeper

	storeReflectCode(t, ctx, quicksilverApp, addr)

	cInfo := wasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)

	return quicksilverApp, ctx
}

func TestQueryFullDenom(t *testing.T) {
	actor := RandomAccountAddress()
	quicksilverApp, ctx := SetupCustomApp(t, actor)

	reflect := instantiateReflectContract(t, ctx, quicksilverApp, actor)
	require.NotEmpty(t, reflect)

	// query full denom
	query := bindings.QuickSilverQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "ustart",
		},
	}
	resp := bindings.FullDenomResponse{}
	queryCustom(t, ctx, quicksilverApp, reflect, query, &resp)

	expected := fmt.Sprintf("factory/%s/ustart", reflect.String())
	require.EqualValues(t, expected, resp.Denom)
}

type ReflectQuery struct {
	Chain *ChainRequest `json:"chain,omitempty"`
}

type ChainRequest struct {
	Request wasmvmtypes.QueryRequest `json:"request"`
}

type ChainResponse struct {
	Data []byte `json:"data"`
}

func queryCustom(t *testing.T, ctx sdk.Context, quicksilver *app.Quicksilver, contract sdk.AccAddress, request bindings.QuickSilverQuery, response interface{}) {
	msgBz, err := json.Marshal(request)
	require.NoError(t, err)

	query := ReflectQuery{
		Chain: &ChainRequest{
			Request: wasmvmtypes.QueryRequest{Custom: msgBz},
		},
	}
	queryBz, err := json.Marshal(query)
	require.NoError(t, err)

	resBz, err := quicksilver.WasmKeeper.QuerySmart(ctx, contract, queryBz)
	require.NoError(t, err)
	var resp ChainResponse
	err = json.Unmarshal(resBz, &resp)
	require.NoError(t, err)
	err = json.Unmarshal(resp.Data, response)
	require.NoError(t, err)
}

func storeReflectCode(t *testing.T, ctx sdk.Context, quicksilverApp *app.Quicksilver, addr sdk.AccAddress) {
	govKeeper := quicksilverApp.GovKeeper
	wasmCode, err := os.ReadFile("../testdata/osmo_reflect.wasm")
	govAddress := govKeeper.GetGovernanceAccount(ctx).GetAddress().String()

	require.NoError(t, err)

	src := wasmtypes.StoreCodeProposalFixture(func(p *wasmtypes.StoreCodeProposal) {
		p.RunAs = addr.String()
		p.WASMByteCode = wasmCode
	})

	msgContent, err := govv1.NewLegacyContent(src, govAddress)
	require.NoError(t, err)

	// when stored
	_, err = govKeeper.SubmitProposal(ctx, []sdk.Msg{msgContent}, "testing123")
	require.NoError(t, err)

	// and proposal execute
	em := sdk.NewEventManager()
	handler := govKeeper.LegacyRouter().GetRoute(src.ProposalRoute())
	err = handler(ctx.WithEventManager(em), src)
	require.NoError(t, err)
}

func instantiateReflectContract(t *testing.T, ctx sdk.Context, quicksilverApp *app.Quicksilver, funder sdk.AccAddress) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(quicksilverApp.WasmKeeper)
	codeID := uint64(1)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	require.NoError(t, err)

	return addr
}

func fundAccount(t *testing.T, ctx sdk.Context, quicksilver *app.Quicksilver, addr sdk.AccAddress, coins sdk.Coins) {
	err := FundAccount(
		quicksilver.BankKeeper,
		ctx,
		addr,
		coins,
	)
	require.NoError(t, err)
}

func FundAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}
