package app_test

import (
	"encoding/json"
	"os"
	"testing"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cometlog "github.com/cometbft/cometbft/libs/log"
	"github.com/quicksilver-zone/quicksilver/v7/app"
	"github.com/rs/zerolog"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"

	cosmossecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestQuicksilverExport(t *testing.T) {
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	senderPubKey := cosmossecp256k1.PubKey{
		Key: senderPrivKey.PubKey().Bytes(),
	}

	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), &senderPubKey, 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000000000000))),
	}
	db := dbm.NewMemDB()
	quicksilver := app.NewQuicksilver(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		app.MakeEncodingConfig(),
		app.EmptyAppOptions{},
		false,
		false,
		app.GetWasmOpts(app.EmptyAppOptions{}),
	)

	genesisState := app.NewDefaultGenesisState()
	genesisState = app.GenesisStateWithValSet(t, quicksilver, genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)

	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	quicksilver.InitChain(
		&abci.RequestInitChain{
			ChainId:       "quicksilver-1",
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	quicksilver.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := app.NewQuicksilver(log.NewCustomLogger(cometlog.NewSyncWriter(os.Stdout).(zerolog.Logger)),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		app.MakeEncodingConfig(),
		app.EmptyAppOptions{},
		false,
		false,
		app.GetWasmOpts(app.EmptyAppOptions{}),
	)
	_, err = app2.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
