package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestGetUnlockedTokensForZoneAllHaveDelegation() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	vals := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	suite.True(found)

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[0].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[1].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[2].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[3].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))

	availPerVal, total, err := quicksilver.InterchainstakingKeeper.GetUnlockedTokensForZone(ctx, &zone)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(400), total)
	for _, x := range availPerVal {
		suite.Equal(sdkmath.NewInt(100), x)
	}
}

func (suite *KeeperTestSuite) TestGetUnlockedTokensForZoneNotAllHaveDelegation() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	vals := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	suite.True(found)

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[0].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[1].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[2].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))

	availPerVal, total, err := quicksilver.InterchainstakingKeeper.GetUnlockedTokensForZone(ctx, &zone)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(300), total)

	// ensure all vals exist in list, even if no delegation
	for _, v := range vals {
		_, found := availPerVal[v.ValoperAddress]
		suite.True(found)
	}
}

func (suite *KeeperTestSuite) TestGetUnlockedTokensForZoneWithRedelegation() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	vals := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	suite.True(found)

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[0].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[1].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[2].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, types.NewDelegation(zone.DelegationAddress.Address, vals[3].ValoperAddress, sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(100))))

	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, types.RedelegationRecord{ChainId: zone.ChainId, EpochNumber: 1, Source: vals[0].ValoperAddress, Destination: vals[2].ValoperAddress, Amount: 5})
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, types.RedelegationRecord{ChainId: zone.ChainId, EpochNumber: 2, Source: vals[0].ValoperAddress, Destination: vals[2].ValoperAddress, Amount: 5})
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, types.RedelegationRecord{ChainId: zone.ChainId, EpochNumber: 2, Source: vals[1].ValoperAddress, Destination: vals[2].ValoperAddress, Amount: 5})

	availPerVal, total, err := quicksilver.InterchainstakingKeeper.GetUnlockedTokensForZone(ctx, &zone)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(385), total)
	suite.Equal(sdkmath.NewInt(85), availPerVal[vals[2].ValoperAddress])
}
