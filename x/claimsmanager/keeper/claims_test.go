package keeper_test

import (
	"cosmossdk.io/math"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
)

var testClaims = []types.Claim{
	// test user claim on chainB (using osmosis pool)
	{
		UserAddress: testAddress,
		// ChainID:       suite.chainB.ChainID,
		Module:        types.ClaimTypeOsmosisPool,
		SourceChainId: "osmosis-1",
		Amount:        math.NewInt(5000000),
	},
	// test user claim on chainB (liquid)
	{
		UserAddress: testAddress,
		// ChainID:       suite.chainB.ChainID,
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "",
		Amount:        math.NewInt(5000000),
	},
	// random user claim on chainB (using osmosis pool)
	{
		UserAddress: addressutils.GenerateAccAddressForTest().String(),
		// ChainID:       suite.chainB.ChainID,
		Module:        types.ClaimTypeOsmosisPool,
		SourceChainId: "osmosis-1",
		Amount:        math.NewInt(15000000),
	},
	// zero value claim
	{
		UserAddress: "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
		// ChainID:       suite.chainB.ChainID,
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "osmosis-1",
		Amount:        math.ZeroInt(),
	},
	// test user claim on "cosmoshub-4" (liquid)
	{
		UserAddress:   testAddress,
		ChainId:       "cosmoshub-4",
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "",
		Amount:        math.NewInt(10000000),
	},
	// random user claim on "cosmoshub-4" (liquid)
	{
		UserAddress:   addressutils.GenerateAccAddressForTest().String(),
		ChainId:       "cosmoshub-4",
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "",
		Amount:        math.NewInt(15000000),
	},
}

func (suite *KeeperTestSuite) TestKeeper_NewClaim() {
	type args struct {
		address    string
		chainID    string
		module     types.ClaimType
		srcChainID string
		amount     math.Int
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
				math.NewInt(5000000),
			},
			types.Claim{
				UserAddress:   testAddress,
				ChainId:       suite.chainB.ChainID,
				Module:        types.ClaimTypeLiquidToken,
				SourceChainId: "",
				Amount:        math.NewInt(5000000),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := types.NewClaim(tt.args.address, tt.args.chainID, tt.args.module, tt.args.srcChainID, tt.args.amount)
			suite.Require().Equal(tt.want, got)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ClaimStore() {
	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper

	testClaims[0].ChainId = suite.chainB.ChainID
	testClaims[1].ChainId = suite.chainB.ChainID
	testClaims[2].ChainId = suite.chainB.ChainID
	testClaims[3].ChainId = suite.chainB.ChainID

	// no claim set
	var getClaim types.Claim
	_, found := k.GetClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().False(found)

	// set claim
	k.SetClaim(suite.chainA.GetContext(), &testClaims[0])

	getClaim, found = k.GetClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().True(found)
	suite.Require().Equal(testClaims[0], getClaim)

	// delete claim
	k.DeleteClaim(suite.chainA.GetContext(), &getClaim)
	getClaim, found = k.GetClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().False(found)

	// iterators
	var claims []*types.Claim

	k.SetClaim(suite.chainA.GetContext(), &testClaims[0])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[1])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[2])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[3])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[4])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[5])

	claims = k.AllClaims(suite.chainA.GetContext())
	suite.Require().Equal(6, len(claims))

	claims = k.AllZoneClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(4, len(claims))

	claims = k.AllZoneClaims(suite.chainA.GetContext(), "cosmoshub-4")
	suite.Require().Equal(2, len(claims))

	// archive (last epoch)
	k.ArchiveAndGarbageCollectClaims(suite.chainA.GetContext(), suite.chainB.ChainID)

	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	suite.Require().True(found)
	suite.Require().Equal(testClaims[0], getClaim)

	// "cosmoshub-4 was not archived so this should not be found"
	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	suite.Require().False(found)

	// set archive claim
	k.SetLastEpochClaim(suite.chainA.GetContext(), &testClaims[4])

	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	suite.Require().True(found)
	suite.Require().Equal(testClaims[4], getClaim)

	// delete archive claim
	k.DeleteLastEpochClaim(suite.chainA.GetContext(), &getClaim)
	getClaim, found = k.GetLastEpochClaim(suite.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	suite.Require().False(found)

	// iterators
	// we expect none as claims have been archived
	claims = k.AllZoneClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(0, len(claims))

	// we expect the archived claims for chainB
	claims = k.AllZoneLastEpochClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(4, len(claims))

	// clear
	k.ClearClaims(suite.chainA.GetContext(), "cosmoshub-4")
	// we expect none as claims have been cleared
	claims = k.AllZoneClaims(suite.chainA.GetContext(), "cosmoshub-4")
	suite.Require().Equal(0, len(claims))

	// we archive current claims (none) to ensure the last epoch claims are correctly set
	k.ArchiveAndGarbageCollectClaims(suite.chainA.GetContext(), suite.chainB.ChainID)

	// we expect none as claims have been archived
	claims = k.AllZoneClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(0, len(claims))

	// we expect none as no current claims existed when we archived
	claims = k.AllZoneLastEpochClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().Equal(0, len(claims))
}

func (suite *KeeperTestSuite) TestIterateUserClaims() {
	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper

	testClaims[0].ChainId = suite.chainB.ChainID
	testClaims[1].ChainId = suite.chainB.ChainID
	testClaims[2].ChainId = suite.chainB.ChainID
	testClaims[3].ChainId = suite.chainB.ChainID

	k.SetClaim(suite.chainA.GetContext(), &testClaims[0])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[1])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[2])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[3])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[4])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[5])

	type args struct {
		chainID string
		address string
		fn      func(index int64, data types.Claim) (stop bool)
	}

	tests := []struct {
		name        string
		args        args
		expectedLen int
	}{
		{
			"blank",
			args{},
			0,
		},
		{
			"valid",
			args{
				chainID: suite.chainB.ChainID,
				address: testAddress,
			},
			2,
		},
		{
			"bad_chain_id",
			args{
				chainID: "badchain",
				address: testAddress,
			},
			0,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			var output []types.Claim
			fn := func(idx int64, data types.Claim) (stop bool) {
				output = append(output, data)
				return false
			}

			tt.args.fn = fn

			k.IterateUserClaims(suite.chainA.GetContext(), tt.args.chainID, tt.args.address, tt.args.fn)
			suite.Require().Equal(tt.expectedLen, len(output))
			output = nil
		})
	}
}

func (suite *KeeperTestSuite) TestAllZoneUserClaims() {
	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper

	testClaims[0].ChainId = suite.chainB.ChainID
	testClaims[1].ChainId = suite.chainB.ChainID

	k.SetClaim(suite.chainA.GetContext(), &testClaims[0])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[1])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[2])

	allClaims := k.AllZoneUserClaims(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress)
	suite.Require().Equal(2, len(allClaims))
}

// func (suite *KeeperTestSuite) TestKeeper_IterateLastEpochUserClaims() {
// 	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper

// 	setClaims[0].ChainID = suite.chainB.ChainID
// 	setClaims[1].ChainID = suite.chainB.ChainID
// 	setClaims[2].ChainID = suite.chainB.ChainID

// 	k.SetLastEpochClaim(suite.chainA.GetContext(), &setClaims[0])
// 	k.SetLastEpochClaim(suite.chainA.GetContext(), &setClaims[1])
// 	k.SetLastEpochClaim(suite.chainA.GetContext(), &setClaims[2])
// 	k.SetLastEpochClaim(suite.chainA.GetContext(), &setClaims[3])
// 	k.SetLastEpochClaim(suite.chainA.GetContext(), &setClaims[4])

// 	type args struct {
// 		chainID string
// 		address string
// 		fn      func(index int64, data types.Claim) (stop bool)
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			"blank",
// 			args{},
// 		},
// 		{
// 			"bad_chain_id",
// 			args{
// 				chainID: "badchain",
// 				address: testAddress,
// 				fn: func(idx int64, data types.Claim) (stop bool) {
// 					fmt.Printf("iterator [%d]: %v\n", idx, data)
// 					return false
// 				},
// 			},
// 		},
// 		{
// 			"suite.chainB.ChainID",
// 			args{
// 				chainID: suite.chainB.ChainID,
// 				address: testAddress,
// 				fn: func(idx int64, data types.Claim) (stop bool) {
// 					fmt.Printf("iterator [%d]: %v\n", idx, data)
// 					return false
// 				},
// 			},
// 		},
// 		{
// 			"chainId_cosmoshub-4",
// 			args{
// 				chainID: "cosmoshub-4",
// 				address: testAddress,
// 				fn: func(idx int64, data types.Claim) (stop bool) {
// 					fmt.Printf("iterator [%d]: %v\n", idx, data)
// 					return false
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		suite.Run(tt.name, func() {
// 			k.IterateLastEpochUserClaims(suite.chainA.GetContext(), tt.args.chainID, tt.args.address, tt.args.fn)
// 		})
// 	}
// }
