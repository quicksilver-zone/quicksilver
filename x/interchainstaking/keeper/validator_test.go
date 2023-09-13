package keeper_test

import (
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *KeeperTestSuite) TestStoreGetDeleteValidator() {
	suite.Run("validator - store / get / delete", func() {
		suite.SetupTest()
		suite.setupTestZones()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)

		validator := addressutils.GenerateValAddressForTest()

		valAddrBytes, err := addressutils.ValAddressFromBech32(validator.String(), zone.GetValoperPrefix())
		suite.NoError(err)
		_, found = app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		suite.False(found)

		count := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))

		newValidator := types.Validator{
			ValoperAddress:  validator.String(),
			CommissionRate:  sdk.NewDec(5.0),
			DelegatorShares: sdk.NewDec(1000.0),
			VotingPower:     sdk.NewInt(1000),
			Status:          stakingtypes.BondStatusBonded,
			Score:           sdk.NewDec(0),
		}
		app.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, newValidator)

		count2 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))

		suite.Equal(count+1, count2)

		fetchedValidator, found := app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		suite.True(found)
		suite.Equal(newValidator, fetchedValidator)

		app.InterchainstakingKeeper.DeleteValidator(ctx, zone.ChainId, valAddrBytes)

		count3 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))
		suite.Equal(count, count3)
	})
}
