package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/x/epochs/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.Quicksilver
	ctx         sdk.Context
	queryClient types.QueryClient
}

// Test helpers.

func (suite *KeeperTestSuite) DoSetupTest(t *testing.T) {
	t.Helper()

	checkTx := false
	suite.app = app.Setup(t, checkTx)

	suite.ctx = suite.app.BaseApp.NewContext(false)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, &suite.app.EpochsKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
