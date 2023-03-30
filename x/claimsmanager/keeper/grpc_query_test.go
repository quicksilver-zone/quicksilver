package keeper_test

import (
	"context"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

func (s *KeeperTestSuite) TestKeeper_Queries() {
	k := s.GetQuicksilverApp(s.chainA).ClaimsManagerKeeper

	// now that we have a kepper set the chainID of chainB
	testClaims[0].ChainId = s.chainB.ChainID
	testClaims[1].ChainId = s.chainB.ChainID
	testClaims[2].ChainId = s.chainB.ChainID
	testClaims[3].ChainId = s.chainB.ChainID

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
				ChainId: s.chainB.ChainID,
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
				k.ArchiveAndGarbageCollectClaims(s.chainA.GetContext(), s.chainB.ChainID)
			},
			&types.QueryClaimsRequest{
				ChainId: s.chainB.ChainID,
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
			2, // 2 expected from chainB, 1 omitted as it was not archived
		},
	}

	k.SetClaim(s.chainA.GetContext(), &testClaims[0])
	k.SetClaim(s.chainA.GetContext(), &testClaims[1])
	k.SetClaim(s.chainA.GetContext(), &testClaims[2])
	k.SetClaim(s.chainA.GetContext(), &testClaims[3])
	k.SetClaim(s.chainA.GetContext(), &testClaims[4])
	k.SetClaim(s.chainA.GetContext(), &testClaims[5])

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := tt.queryFn(s.chainA.GetContext(), tt.req)
			s.Require().NoError(err)
			s.Require().NotNil(resp.Claims)
			s.Require().Equal(tt.expectLength, len(resp.Claims))
		})
	}
}
