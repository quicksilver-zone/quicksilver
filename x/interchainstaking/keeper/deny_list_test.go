package keeper_test

func (suite *KeeperTestSuite) TestStoreGetDeleteDenyList() {
	suite.Run("deny list - store / get / delete single", func() {
		suite.SetupTest()
		suite.setupTestZones()

		qApp := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := qApp.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)
		vals := qApp.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
		suite.Len(vals, 4)
		val := vals[0]
		validator := vals[0].ValoperAddress
		// Initially the deny list should be empty
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator)
		suite.False(found)
		// Add a validator to the deny list
		err := qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, val)
		suite.NoError(err)

		// Ensure the deny list contains the validator

		denyList := qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 1)
		suite.Equal(validator, denyList[0])

		val, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator)
		suite.True(found)
		suite.Equal(validator, val.ValoperAddress)

		// Remove the validator from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, val)
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator)
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
		vals := qApp.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
		suite.Len(vals, 4)
		val1, val2, val3 := vals[0], vals[1], vals[2]
		valAddr1, valAddr2, valAddr3 := val1.ValoperAddress, val2.ValoperAddress, val3.ValoperAddress

		// Initially the deny list should be empty
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr1)
		suite.False(found)

		// Add three validators to the deny list
		err := qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, val1)
		suite.NoError(err)
		err = qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, val2)
		suite.NoError(err)
		err = qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, val3)
		suite.NoError(err)

		denyList := qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 3)

		denyVal1, found := qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr1)
		suite.True(found)
		suite.Equal(valAddr1, denyVal1.ValoperAddress)

		denyVal2, found := qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr2)
		suite.True(found)
		suite.Equal(valAddr2, denyVal2.ValoperAddress)

		denyVal3, found := qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr3)
		suite.True(found)
		suite.Equal(valAddr3, denyVal3.ValoperAddress)

		// Remove the validator from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, val2)
		_, found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr2)
		suite.False(found)

		// Ensure the deny list contains the two remaining validators
		denyList = qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 2)
		suite.NotContains(denyList, valAddr2)
		suite.Contains(denyList, valAddr1)
		suite.Contains(denyList, valAddr3)

		// Remove the remaining two validators from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, val1)
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, val3)

		// Ensure the deny list is empty
		denyList = qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 0)
	})
}
