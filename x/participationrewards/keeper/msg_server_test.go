package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	osmolockup "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
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
					SrcZone:     suite.chainB.ChainID,
					ClaimType:   cmtypes.ClaimTypeOsmosisPool,
					Proofs: []*cmtypes.Proof{
						{
							Key:       []byte{1},
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
			"invalid_osmosis_user",
			func() {
				userAddress := utils.GenerateAccAddressForTest()
				osmoAddress := utils.GenerateAccAddressForTestWithPrefix("osmo")
				lockedResp := osmolockup.LockedResponse{
					Lock: &osmolockup.PeriodLock{
						ID:       1,
						Owner:    osmoAddress,
						Duration: time.Hour * 72,
						Coins: sdk.NewCoins(
							sdk.NewCoin(
								"gamm/1",
								math.NewInt(10000000),
							),
						),
					},
				}
				bz, err := lockedResp.Marshal()
				suite.Require().NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: userAddress.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "osmosis-1",
					ClaimType:   cmtypes.ClaimTypeOsmosisPool,
					Proofs: []*cmtypes.Proof{
						{
							Key:       []byte{2, 1},
							Data:      bz,
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
		{
			"invalid_osmosis_pool",
			func() {
				userAddress := utils.GenerateAccAddressForTest()
				osmoAddress := utils.ConvertAccAddressForTestUsingPrefix(userAddress, "osmo")
				lockedResp := osmolockup.LockedResponse{
					Lock: &osmolockup.PeriodLock{
						ID:       1,
						Owner:    osmoAddress,
						Duration: time.Hour * 72,
						Coins: sdk.NewCoins(
							sdk.NewCoin(
								"gamm/2",
								math.NewInt(10000000),
							),
						),
					},
				}
				bz, err := lockedResp.Marshal()
				suite.Require().NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: userAddress.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "osmosis-1",
					ClaimType:   cmtypes.ClaimTypeOsmosisPool,
					Proofs: []*cmtypes.Proof{
						{
							Key:       []byte{2, 1},
							Data:      bz,
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
		{
			"valid_osmosis",
			func() {
				userAddress := utils.GenerateAccAddressForTest()
				osmoAddress := utils.ConvertAccAddressForTestUsingPrefix(userAddress, "osmo")
				lockedResp := osmolockup.LockedResponse{
					Lock: &osmolockup.PeriodLock{
						ID:       1,
						Owner:    osmoAddress,
						Duration: time.Hour * 72,
						Coins: sdk.NewCoins(
							sdk.NewCoin(
								"gamm/1",
								math.NewInt(10000000),
							),
						),
					},
				}
				bz, err := lockedResp.Marshal()
				suite.Require().NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: userAddress.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "osmosis-1",
					ClaimType:   cmtypes.ClaimTypeOsmosisPool,
					Proofs: []*cmtypes.Proof{
						{
							Key:       []byte{2, 1},
							Data:      bz,
							ProofOps:  &crypto.ProofOps{},
							Height:    0,
							ProofType: "lockup",
						},
					},
				}
			},
			&types.MsgSubmitClaimResponse{},
			false,
		},
		// {
		// 	"valid_liquid",
		// 	func() {
		// 		address := utils.GenerateAccAddressForTest()
		// 		key := append(address, []byte("uqatom")...)
		// 		msg = types.MsgSubmitClaim{
		// 			UserAddress: address.String(),
		// 			Zone:        suite.chainB.ChainID,
		// 			SrcZone:     "osmosis-1",
		// 			ClaimType:   cmtypes.ClaimTypeLiquidToken,
		// 			Proofs: []*cmtypes.Proof{
		// 				{
		// 					Key:       key,
		// 					Data:      []byte{0, 0, 1, 1, 2, 3, 4, 5},
		// 					ProofOps:  &crypto.ProofOps{},
		// 					Height:    0,
		// 					ProofType: "lockup",
		// 				},
		// 			},
		// 		}
		// 	},
		// 	&types.MsgSubmitClaimResponse{},
		// 	false,
		// },
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
