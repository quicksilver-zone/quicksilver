package keeper_test

import (
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

var testClaims = []types.Claim{
	// test user claim on chainB (using osmosis pool)
	{
		UserAddress: testAddress,
		// ChainId:       suite.chainB.ChainID,
		Module:        types.ClaimTypeOsmosisPool,
		SourceChainId: "osmosis-1",
		Amount:        5000000,
	},
	// test user claim on chainB (liquid)
	{
		UserAddress: testAddress,
		// ChainId:       suite.chainB.ChainID,
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "",
		Amount:        5000000,
	},
	// random user claim on chainB (using osmosis pool)
	{
		UserAddress: utils.GenerateAccAddressForTest().String(),
		// ChainId:       suite.chainB.ChainID,
		Module:        types.ClaimTypeOsmosisPool,
		SourceChainId: "osmosis-1",
		Amount:        15000000,
	},
	// zero value claim
	{
		UserAddress: "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
		// ChainId:       suite.chainB.ChainID,
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "osmosis-1",
		Amount:        0,
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

func (suite *KeeperTestSuite) TestKeeper_NewClaim() {
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

	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := k.NewClaim(suite.chainA.GetContext(), tt.args.address, tt.args.chainID, tt.args.module, tt.args.srcChainID, tt.args.amount)
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
	var found bool

	getClaim, found = k.GetClaim(suite.chainA.GetContext(), suite.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
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

// func (suite *KeeperTestSuite) TestKeeper_IterateLastEpochUserClaims() {
// 	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper

// 	setClaims[0].ChainId = suite.chainB.ChainID
// 	setClaims[1].ChainId = suite.chainB.ChainID
// 	setClaims[2].ChainId = suite.chainB.ChainID

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
