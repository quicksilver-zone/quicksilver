package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) Test_msgServer_SubmitClaim() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	msg := types.MsgSubmitClaim{}
	tests := []struct {
		name     string
		malleate func()
		want     *types.MsgSubmitClaimResponse
		wantErr  bool
	}{
		{
			"blank",
			func() {},
			nil,
			true,
		},
		{
			"invalid_height",
			func() {
				msg = types.MsgSubmitClaim{
					UserAddress: utils.GenerateAccAddressForTest().String(),
					Zone:        suite.chainB.ChainID,
					ClaimType:   cmtypes.ClaimTypeOsmosisPool,
					Proofs: []*cmtypes.Proof{
						{
							Key:       []byte{1, 2, 3, 4, 5},
							Data:      []byte{0, 0, 1, 1, 2, 3, 4, 5},
							ProofOps:  &crypto.ProofOps{},
							Height:    123,
							ProofType: "lockup",
						},
					},
				}
			},
			nil,
			true,
		},
		{
			"test",
			func() {
				msg = types.MsgSubmitClaim{
					UserAddress: utils.GenerateAccAddressForTest().String(),
					Zone:        suite.chainB.ChainID,
					ClaimType:   cmtypes.ClaimTypeOsmosisPool,
					Proofs: []*cmtypes.Proof{
						{
							Key:       []byte{1, 2, 3, 4, 5},
							Data:      []byte{0, 0, 1, 1, 2, 3, 4, 5},
							ProofOps:  &crypto.ProofOps{},
							Height:    0,
							ProofType: "lockup",
						},
					},
				}
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			tt.malleate()

			k := keeper.NewMsgServerImpl(appA.ParticipationRewardsKeeper)
			resp, err := k.SubmitClaim(sdk.WrapSDKContext(suite.chainA.GetContext()), &msg)
			if tt.wantErr {
				suite.Require().Error(err)
				suite.Require().Nil(resp)
				suite.T().Logf("Error: %v", err)
				return
			}

			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.want, resp)
		})
	}
}
