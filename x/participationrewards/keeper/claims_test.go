package keeper_test

import (
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) TestKeeper_NewClaim() {
	testAddress := utils.GenerateAccAddressForTest().String()
	type args struct {
		address    string
		chainID    string
		module     types.ClaimType
		srcChainID string
		amount     uint64
	}
	tests := []struct {
		name string
		args args
		want types.Claim
	}{
		{
			"blank",
			args{},
			types.Claim{},
		},
		{
			"valid",
			args{
				testAddress,
				suite.chainB.ChainID,
				types.ClaimTypeLiquidToken,
				"",
				5000000,
			},
			types.Claim{
				UserAddress:   testAddress,
				ChainId:       suite.chainB.ChainID,
				Module:        types.ClaimTypeLiquidToken,
				SourceChainId: "",
				Amount:        5000000,
			},
		},
	}

	k := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := k.NewClaim(suite.chainA.GetContext(), tt.args.address, tt.args.chainID, tt.args.module, tt.args.srcChainID, tt.args.amount)
			suite.Require().Equal(tt.want, got)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ClaimStore() {
	testAddress := utils.GenerateAccAddressForTest().String()
	setClaims := []types.Claim{
		// test user claim on chainB (using osmosis pool)
		{
			UserAddress:   testAddress,
			ChainId:       suite.chainB.ChainID,
			Module:        types.ClaimTypeOsmosisPool,
			SourceChainId: "osmosis-1",
			Amount:        5000000,
		},
		// test user claim on chainB (liquid)
		{
			UserAddress:   testAddress,
			ChainId:       suite.chainB.ChainID,
			Module:        types.ClaimTypeLiquidToken,
			SourceChainId: "",
			Amount:        5000000,
		},
		// random user claim on chainB (using osmosis pool)
		{
			UserAddress:   utils.GenerateAccAddressForTest().String(),
			ChainId:       suite.chainB.ChainID,
			Module:        types.ClaimTypeOsmosisPool,
			SourceChainId: "osmosis-1",
			Amount:        15000000,
		},
		// test user claim on "cosmoshub-4" (liquid)
		{
			UserAddress:   testAddress,
			ChainId:       "cosmoshub-4",
			Module:        types.ClaimTypeLiquidToken,
			SourceChainId: "",
			Amount:        10000000,
		},
		// random user claim on "cosmoshub-4" (liquid)
		{
			UserAddress:   utils.GenerateAccAddressForTest().String(),
			ChainId:       "cosmoshub-4",
			Module:        types.ClaimTypeLiquidToken,
			SourceChainId: "",
			Amount:        15000000,
		},
	}

	k := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper

	// no claim set
	var getClaim types.Claim
	var found bool

	getClaim, found = k.GetClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().False(found)

	// set claim
	k.SetClaim(suite.chainA.GetContext(), &setClaims[0])

	getClaim, found = k.GetClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().True(found)
	suite.Require().Equal(setClaims[0], getClaim)

	// delete claim
	k.DeleteClaim(suite.chainA.GetContext(), &getClaim)
	getClaim, found = k.GetClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().False(found)

	// iterators
	var claims []*types.Claim

	k.SetClaim(suite.chainA.GetContext(), &setClaims[0])
	k.SetClaim(suite.chainA.GetContext(), &setClaims[1])
	k.SetClaim(suite.chainA.GetContext(), &setClaims[2])
	k.SetClaim(suite.chainA.GetContext(), &setClaims[3])
	k.SetClaim(suite.chainA.GetContext(), &setClaims[4])

	claims = k.AllClaims(suite.chainA.GetContext())
	suite.Require().Equal(5, len(claims))

	claims = k.AllZoneClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(3, len(claims))

	claims = k.AllZoneClaims(suite.chainA.GetContext(), "cosmoshub-4")
	suite.Require().Equal(2, len(claims))

	claims = k.AllZoneUserClaims(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress)
	suite.Require().Equal(2, len(claims))

	claims = k.AllZoneUserClaims(suite.chainA.GetContext(), "cosmoshub-4", testAddress)
	suite.Require().Equal(1, len(claims))

	// archive (last epoch)
	k.ArchiveAndGarbageCollectClaims(suite.chainA.GetContext(), suite.chainB.ChainID)

	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().True(found)
	suite.Require().Equal(setClaims[0], getClaim)

	// "cosmoshub-4 was not archived so this should not be found"
	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	suite.Require().False(found)

	// set archive claim
	k.SetLastEpochClaim(suite.chainA.GetContext(), &setClaims[3])

	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	suite.Require().True(found)
	suite.Require().Equal(setClaims[3], getClaim)

	// delete archive claim
	k.DeleteLastEpochClaim(suite.chainA.GetContext(), &getClaim)
	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	suite.Require().False(found)

	// iterators
	claims = k.AllZoneClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(0, len(claims))

	claims = k.AllZoneLastEpochClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(3, len(claims))

	claims = k.AllZoneLastEpochUserClaims(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress)
	suite.Require().Equal(2, len(claims))

	// clear
	k.ClearClaims(suite.chainA.GetContext(), "cosmoshub-4")
	claims = k.AllZoneClaims(suite.chainA.GetContext(), "cosmoshub-4")
	suite.Require().Equal(0, len(claims))

	k.ArchiveAndGarbageCollectClaims(suite.chainA.GetContext(), suite.chainB.ChainID)

	claims = k.AllZoneClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(0, len(claims))

	claims = k.AllZoneLastEpochClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(0, len(claims))
}
