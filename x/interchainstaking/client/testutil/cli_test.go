package testutil

import (
	"testing"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	"github.com/stretchr/testify/suite"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	cfg.AppConstructor = NewQuicksilverConstructor()
	types.RegisterInterfaces(cfg.InterfaceRegistry)
	suite.Run(t, NewIntegrationTestSuite(cfg))
}

func NewQuicksilverConstructor() network.AppConstructor {
	//db := dbm.NewMemDB()
	return func(val network.Validator) servertypes.Application {
		return app.Setup(false)
		/*return app.NewQuicksilver(
			val.Ctx.Logger,
			db, nil, true, map[int64]bool{},
			app.DefaultNodeHome,
			0,
			app.MakeEncodingConfig(),
			simapp.EmptyAppOptions{},
			baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)*/
	}
}
