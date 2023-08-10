package keeper_test

import (
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

var testClaims = []types.Claim{
	// test user claim on chainB (using osmosis pool)
	{
		UserAddress: testAddress,
		// ChainID:       suite.chainB.ChainID,
		Module:        types.ClaimTypeOsmosisPool,
		SourceChainId: "osmosis-1",
		Amount:        5000000,
	},
	// test user claim on chainB (liquid)
	{
		UserAddress: testAddress,
		// ChainID:       suite.chainB.ChainID,
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "",
		Amount:        5000000,
	},
	// random user claim on chainB (using osmosis pool)
	{
		UserAddress: addressutils.GenerateAccAddressForTest().String(),
		// ChainID:       suite.chainB.ChainID,
		Module:        types.ClaimTypeOsmosisPool,
		SourceChainId: "osmosis-1",
		Amount:        15000000,
	},
	// zero value claim
	{
		UserAddress: "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
		// ChainID:       suite.chainB.ChainID,
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
		UserAddress:   addressutils.GenerateAccAddressForTest().String(),
		ChainId:       "cosmoshub-4",
		Module:        types.ClaimTypeLiquidToken,
		SourceChainId: "",
		Amount:        15000000,
	},
}

func (s *KeeperTestSuite) TestKeeper_NewClaim() {
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
				s.chainB.ChainID,
				types.ClaimTypeLiquidToken,
				"",
				5000000,
			},
			types.Claim{
				UserAddress:   testAddress,
				ChainId:       s.chainB.ChainID,
				Module:        types.ClaimTypeLiquidToken,
				SourceChainId: "",
				Amount:        5000000,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := types.NewClaim(tt.args.address, tt.args.chainID, tt.args.module, tt.args.srcChainID, tt.args.amount)
			s.Require().Equal(tt.want, got)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_ClaimStore() {
	k := s.GetQuicksilverApp(s.chainA).ClaimsManagerKeeper
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper

	testClaims[0].ChainId = s.chainB.ChainID
	testClaims[1].ChainId = s.chainB.ChainID
	testClaims[2].ChainId = s.chainB.ChainID
	testClaims[3].ChainId = s.chainB.ChainID

	// no claim set
	var getClaim types.Claim
	_, found := k.GetClaim(s.chainA.GetContext(), s.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	s.Require().False(found)

	// set claim
	k.SetClaim(s.chainA.GetContext(), &testClaims[0])

	getClaim, found = k.GetClaim(s.chainA.GetContext(), s.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	s.Require().True(found)
	s.Require().Equal(testClaims[0], getClaim)

	// delete claim
	k.DeleteClaim(s.chainA.GetContext(), &getClaim)
	getClaim, found = k.GetClaim(s.chainA.GetContext(), s.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	s.Require().False(found)

	// iterators
	var claims []*types.Claim

	k.SetClaim(s.chainA.GetContext(), &testClaims[0])
	k.SetClaim(s.chainA.GetContext(), &testClaims[1])
	k.SetClaim(s.chainA.GetContext(), &testClaims[2])
	k.SetClaim(s.chainA.GetContext(), &testClaims[3])
	k.SetClaim(s.chainA.GetContext(), &testClaims[4])
	k.SetClaim(s.chainA.GetContext(), &testClaims[5])

	claims = k.AllClaims(s.chainA.GetContext())
	s.Require().Equal(6, len(claims))

	claims = k.AllZoneClaims(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().Equal(4, len(claims))

	claims = k.AllZoneClaims(s.chainA.GetContext(), "cosmoshub-4")
	s.Require().Equal(2, len(claims))

	zone, found := icsKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().True(found)

	// archive (last epoch)
	k.ArchiveAndGarbageCollectClaims(s.chainA.GetContext(), &zone)

	getClaim, found = k.GetLastEpochClaim(s.chainA.GetContext(), s.chainB.ChainID, testAddress, types.ClaimTypeOsmosisPool, "osmosis-1")
	s.Require().True(found)
	s.Require().Equal(testClaims[0], getClaim)

	// "cosmoshub-4 was not archived so this should not be found"
	getClaim, found = k.GetLastEpochClaim(s.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	s.Require().False(found)

	// set archive claim
	k.SetLastEpochClaim(s.chainA.GetContext(), &testClaims[4])

	getClaim, found = k.GetLastEpochClaim(s.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	s.Require().True(found)
	s.Require().Equal(testClaims[4], getClaim)

	// delete archive claim
	k.DeleteLastEpochClaim(s.chainA.GetContext(), &getClaim)
	getClaim, found = k.GetLastEpochClaim(s.chainA.GetContext(), "cosmoshub-4", testAddress, types.ClaimTypeLiquidToken, "")
	s.Require().False(found)

	// iterators
	// we expect none as claims have been archived
	claims = k.AllZoneClaims(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().Equal(0, len(claims))

	// we expect the archived claims for chainB
	claims = k.AllZoneLastEpochClaims(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().Equal(4, len(claims))

	// clear
	k.ClearClaims(s.chainA.GetContext(), "cosmoshub-4")
	// we expect none as claims have been cleared
	claims = k.AllZoneClaims(s.chainA.GetContext(), "cosmoshub-4")
	s.Require().Equal(0, len(claims))

	zone, found = icsKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().True(found)

	// we archive current claims (none) to ensure the last epoch claims are correctly set
	k.ArchiveAndGarbageCollectClaims(s.chainA.GetContext(), &zone)

	// we expect none as claims have been archived
	claims = k.AllZoneClaims(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().Equal(0, len(claims))

	// we expect none as no current claims existed when we archived
	claims = k.AllZoneLastEpochClaims(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().Equal(0, len(claims))
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
