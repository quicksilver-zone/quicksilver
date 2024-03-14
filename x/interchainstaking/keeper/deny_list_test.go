package keeper_test

import (
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestStoreGetDeleteDenyList() {
	suite.Run("deny list - store / get / delete single", func() {
		suite.SetupTest()
		suite.setupTestZones()

		qApp := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := qApp.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)
		validator := addressutils.GenerateValAddressForTest()

		// Initially the deny list should be empty
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator.String())
		suite.False(found)
		// Add a validator to the deny list
		err := qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, types.Validator{
			ValoperAddress: validator.String(),
		})
		suite.NoError(err)

		// Ensure the deny list contains the validator

		denyList := qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 1)
		suite.Equal(validator.String(), denyList[0].ValoperAddress)

		val, found := qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator.String())
		suite.True(found)
		suite.Equal(validator.String(), val.ValoperAddress)

		// Remove the validator from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, val)
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator.String())
		suite.False(found)

		denyList = qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 0)
	})

	suite.Run("deny list - store / get / delete multiple", func() {
		suite.SetupTest()
		suite.setupTestZones()

		qApp := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := qApp.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)

		valAddr1 := addressutils.GenerateValAddressForTest()
		valAddr2 := addressutils.GenerateValAddressForTest()
		valAddr3 := addressutils.GenerateValAddressForTest()

		// Initially the deny list should be empty
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr1.String())
		suite.False(found)

		// Add three validators to the deny list
		err := qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, types.Validator{
			ValoperAddress: valAddr1.String(),
		})
		suite.NoError(err)
		err = qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, types.Validator{
			ValoperAddress: valAddr2.String(),
		})
		suite.NoError(err)
		err = qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, types.Validator{
			ValoperAddress: valAddr3.String(),
		})
		suite.NoError(err)

		// Ensure the deny list contains the three validators
		denyList := qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 3)

		val1, found := qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr1.String())
		suite.True(found)
		suite.Equal(valAddr1.String(), val1.ValoperAddress)

		val2, found := qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr2.String())
		suite.True(found)
		suite.Equal(valAddr2.String(), val2.ValoperAddress)

		val3, found := qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr3.String())
		suite.True(found)
		suite.Equal(valAddr3.String(), val3.ValoperAddress)

		// Remove the validator from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, val2)
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, val2.String())
		suite.False(found)

		// Ensure the deny list contains the two remaining validators
		denyList = qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 2)
	})
}
