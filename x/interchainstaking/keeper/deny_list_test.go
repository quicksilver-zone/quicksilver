package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

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
		validator := sdk.ValAddress(val.ValoperAddress)
		// Initially the deny list should be empty
		found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator)
		suite.False(found)
		// Add a validator to the deny list
		err := qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, validator)
		suite.NoError(err)

		// Ensure the deny list contains the validator

		denyList := qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 1)
		suite.Equal(addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), validator), denyList[0])

		found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator)
		suite.True(found)

		// Remove the validator from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, validator)
		found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, validator)
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
		valAddr := make([]sdk.ValAddress, 4)
		for i, v := range vals {
			valAddr[i] = sdk.ValAddress(v.ValoperAddress)
		}
		// Initially the deny list should be empty
		found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr[0])
		suite.False(found)

		// Add three validators to the deny list
		err := qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, valAddr[0])
		suite.NoError(err)
		err = qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, valAddr[1])
		suite.NoError(err)
		err = qApp.InterchainstakingKeeper.SetZoneValidatorToDenyList(ctx, zone.ChainId, valAddr[2])
		suite.NoError(err)

		denyList := qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 3)

		// Ensure the deny list contains the three validators
		suite.Contains(denyList, addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), valAddr[0]))
		suite.Contains(denyList, addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), valAddr[1]))
		suite.Contains(denyList, addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), valAddr[2]))

		// Remove the validator from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, valAddr[1])
		found = qApp.InterchainstakingKeeper.GetDeniedValidatorInDenyList(ctx, zone.ChainId, valAddr[1])
		suite.False(found)

		// Ensure the deny list contains the two remaining validators
		denyList = qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 2)
		suite.NotContains(denyList, addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), valAddr[1]))
		suite.Contains(denyList, addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), valAddr[0]))
		suite.Contains(denyList, addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), valAddr[2]))

		// Remove the remaining two validators from the deny list
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, valAddr[0])
		qApp.InterchainstakingKeeper.RemoveValidatorFromDenyList(ctx, zone.ChainId, valAddr[2])

		// Ensure the deny list is empty
		denyList = qApp.InterchainstakingKeeper.GetZoneValidatorDenyList(ctx, zone.ChainId)
		suite.Len(denyList, 0)
	})
}
