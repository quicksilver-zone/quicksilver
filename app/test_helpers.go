package app

// DONTCOVER

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	ibctesting "github.com/cosmos/ibc-go/v7/testing"
)

// EmptyAppOptions is a stub implementing AppOptions.
type EmptyAppOptions struct{}

// Get implements AppOptions.
func (EmptyAppOptions) Get(_ string) interface{} {
	return nil
}

// Setup initializes a new Quicksilver. A Nop logger is set in Quicksilver.
func Setup(t *testing.T, isCheckTx bool) *Quicksilver {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := mock.NewPV()
	senderPubkey := senderPrivKey.PrivKey.PubKey()

	acc := authtypes.NewBaseAccount(senderPubkey.Address().Bytes(), senderPubkey, 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	db := dbm.NewMemDB()
	app := NewQuicksilver(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		MakeEncodingConfig(),
		EmptyAppOptions{},
		false,
		false,
		"",
		baseapp.SetChainID("quicksilver-1"),
	)

	genesisState := NewDefaultGenesisState()
	genesisState = GenesisStateWithValSet(t, app, genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)

	if !isCheckTx {
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				ChainId:    "quicksilver-1",
				Validators: []abci.ValidatorUpdate{},
				// ConsensusParams: DefaultConsensusParams,
				AppStateBytes: stateBytes,
			},
		)
	}

	return app
}

func GetAppWithContext(t *testing.T, init bool) (*Quicksilver, sdk.Context) {
	t.Helper()

	app := Setup(t, !init)
	ctx := app.NewContext(false, tmproto.Header{Height: 1, ChainID: "mercury-1", Time: time.Now().UTC()})
	return app, ctx
}

// SetupTestingApp initializes the IBC-go testing application.
func SetupTestingApp() (testApp ibctesting.TestingApp, genesisState map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	app := NewQuicksilver(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		5,
		MakeEncodingConfig(),
		EmptyAppOptions{},
		true, // set mock state to true
		false,
		"",
	)
	return app, NewDefaultGenesisState()
}

// GenesisStateWithValSet creates a quicksilver genesis state with the given validator set.
func GenesisStateWithValSet(
	t *testing.T,
	app *Quicksilver, genesisState GenesisState,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) GenesisState {
	t.Helper()
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			Jailed:          false,
			Status:          stakingtypes.Bonded,
			Tokens:          bondAmt,
			DelegatorShares: sdk.OneDec(),
			Description:     stakingtypes.Description{},
			UnbondingHeight: int64(0),
			UnbondingTime:   time.Unix(0, 0).UTC(),
			Commission:      stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	return genesisState
}
