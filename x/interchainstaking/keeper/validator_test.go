package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestStoreGetDeleteValidator() {
	suite.Run("validator - store / get / delete", func() {
		suite.SetupTest()
		suite.setupTestZones()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.Require().True(found)

		validator := addressutils.GenerateValAddressForTest()

		valAddrBytes, err := addressutils.ValAddressFromBech32(validator.String(), zone.GetValoperPrefix())
		suite.Require().NoError(err)
		_, found = app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainID(), valAddrBytes)
		suite.Require().False(found)

		count := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainID()))

		newValidator := types.Validator{
			ValoperAddress:  validator.String(),
			CommissionRate:  sdk.NewDec(5.0),
			DelegatorShares: sdk.NewDec(1000.0),
			VotingPower:     sdk.NewInt(1000),
			Status:          stakingtypes.BondStatusBonded,
			Score:           sdk.NewDec(0),
		}
		app.InterchainstakingKeeper.SetValidator(ctx, zone.ChainID(), newValidator)

		count2 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainID()))

		suite.Require().Equal(count+1, count2)

		fetchedValidator, found := app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainID(), valAddrBytes)
		suite.Require().True(found)
		suite.Require().Equal(newValidator, fetchedValidator)

		app.InterchainstakingKeeper.DeleteValidator(ctx, zone.ChainID(), valAddrBytes)

		count3 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainID()))
		suite.Require().Equal(count, count3)
	})
}
