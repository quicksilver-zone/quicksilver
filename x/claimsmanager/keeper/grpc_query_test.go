package keeper_test

import (
	"context"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

func (suite *KeeperTestSuite) TestKeeper_Queries() {
	k := suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper

	// now that we have a kepper set the chainID of chainB
	testClaims[0].ChainId = suite.chainB.ChainID
	testClaims[1].ChainId = suite.chainB.ChainID
	testClaims[2].ChainId = suite.chainB.ChainID
	testClaims[3].ChainId = suite.chainB.ChainID

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryClaimsRequest
		queryFn      func(context.Context, *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error)
		expectLength int
	}{
		{
			"Claims_chainB",
			func() {},
			&types.QueryClaimsRequest{
				ChainId: suite.chainB.ChainID,
			},
			k.Claims,
			4,
		},
		{
			"Claims_cosmoshub",
			func() {},
			&types.QueryClaimsRequest{
				ChainId: "cosmoshub-4",
			},
			k.Claims,
			2,
		},
		{
			"UserClaims_testAddress",
			func() {},
			&types.QueryClaimsRequest{
				Address: testAddress,
			},
			k.UserClaims,
			3,
		},
		{
			"LastEpochClaims_chainB",
			func() {
				k.ArchiveAndGarbageCollectClaims(suite.chainA.GetContext(), suite.chainB.ChainID)
			},
			&types.QueryClaimsRequest{
				ChainId: suite.chainB.ChainID,
			},
			k.LastEpochClaims,
			4,
		},
		{
			"LastEpochClaims_cosmoshub",
			func() {
			},
			&types.QueryClaimsRequest{
				ChainId: "cosmoshub-4",
			},
			k.LastEpochClaims,
			0, // none expected as this zone was not archived
		},
		{
			"UserLastEpochClaims_testAddress",
			func() {
			},
			&types.QueryClaimsRequest{
				Address: testAddress,
			},
			k.UserLastEpochClaims,
			2, // 2 expected from chainB, 1 ommited as it was not archived
		},
	}

	k.SetClaim(suite.chainA.GetContext(), &testClaims[0])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[1])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[2])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[3])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[4])
	k.SetClaim(suite.chainA.GetContext(), &testClaims[5])

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := tt.queryFn(suite.chainA.GetContext(), tt.req)
			suite.Require().NoError(err)
			suite.Require().NotNil(resp.Claims)
			suite.Require().Equal(tt.expectLength, len(resp.Claims))
		})
	}
}
