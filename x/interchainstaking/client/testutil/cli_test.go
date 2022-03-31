package testutil

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/ingenuity-build/quicksilver/app"
	dbm "github.com/tendermint/tm-db"

	"github.com/stretchr/testify/suite"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	cfg.AppConstructor = NewQuicksilverConstructor()
	//types.RegisterInterfaces(cfg.InterfaceRegistry)
	suite.Run(t, NewIntegrationTestSuite(cfg))
}

func NewQuicksilverConstructor() network.AppConstructor {
	db := dbm.NewMemDB()
	return func(val network.Validator) servertypes.Application {
		return app.NewQuicksilver(
			val.Ctx.Logger,
			db, nil, true, map[int64]bool{},
			app.DefaultNodeHome,
			0,
			app.MakeEncodingConfig(),
			simapp.EmptyAppOptions{},
			baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			//baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}
