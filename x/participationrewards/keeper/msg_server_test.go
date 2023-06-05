package keeper_test

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/ingenuity-build/quicksilver/app"
	osmolockup "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) Test_msgServer_SubmitClaim() {
	// TODO: these tests ought to validate the error received.
	appA := suite.GetQuicksilverApp(suite.chainA)

	msg := types.MsgSubmitClaim{}
	tests := []struct {
		name     string
		malleate func()
		want     *types.MsgSubmitClaimResponse
		wantErr  string
	}{
		{
			"blank",
			func() {},
			nil,
			"a",
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
			"a",
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
			"a",
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
			"a",
		},
		{
			"valid_osmosis",
			func() {
				userAddress := utils.GenerateAccAddressForTest()
				osmoAddress := utils.ConvertAccAddressForTestUsingPrefix(userAddress, "osmo")
				locked := &osmolockup.PeriodLock{
					ID:       1,
					Owner:    osmoAddress,
					Duration: time.Hour * 72,
					Coins: sdk.NewCoins(
						sdk.NewCoin(
							"gamm/1",
							math.NewInt(10000000),
						),
					),
				}
				bz, err := locked.Marshal()
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
			"",
		},
		{
			"valid_liquid",
			func() {
				address := utils.GenerateAccAddressForTest()
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3"))

				cd := sdk.Coin{
					Denom:  "",
					Amount: math.NewInt(0),
				}
				bz, err := cd.Marshal()
				suite.Require().NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "osmosis-1",
					ClaimType:   cmtypes.ClaimTypeLiquidToken,
					Proofs: []*cmtypes.Proof{
						{
							Key:       key,
							Data:      bz,
							ProofOps:  &crypto.ProofOps{},
							Height:    0,
							ProofType: "bank",
						},
					},
				}
			},
			&types.MsgSubmitClaimResponse{},
			"",
		},
		{
			"valid_liquid",
			func() {
				address := utils.GenerateAccAddressForTest()
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3"))

				cd := sdk.Coin{
					Denom:  "",
					Amount: math.NewInt(0),
				}
				bz, err := cd.Marshal()
				suite.Require().NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "testchain1",
					ClaimType:   cmtypes.ClaimTypeLiquidToken,
					Proofs: []*cmtypes.Proof{
						{
							Key:       key,
							Data:      bz,
							ProofOps:  &crypto.ProofOps{},
							Height:    10,
							ProofType: "bank",
						},
					},
				}
			},
			&types.MsgSubmitClaimResponse{},
			"",
		},
	}
	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			tt.malleate()
			ctx := suite.chainA.GetContext()
			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			k := keeper.NewMsgServerImpl(appA.ParticipationRewardsKeeper)
			resp, err := k.SubmitClaim(sdk.WrapSDKContext(ctx), &msg)
			if tt.wantErr != "" {
				suite.Require().Errorf(err, tt.wantErr)
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

func (s *KeeperTestSuite) Test_msgServer_SubmitLocalClaim() {
	address := utils.GenerateAccAddressForTest()

	var msg *types.MsgSubmitClaim
	tests := []struct {
		name     string
		malleate func(ctx sdk.Context, appA *app.Quicksilver)
		generate func(ctx sdk.Context, appA *app.Quicksilver) *types.MsgSubmitClaim
		want     *types.MsgSubmitClaimResponse
		wantErr  string
		claims   []cmtypes.Claim
	}{
		{
			"local_callback_nil",
			func(ctx sdk.Context, appA *app.Quicksilver) {},
			func(ctx sdk.Context, appA *app.Quicksilver) *types.MsgSubmitClaim {
				address := utils.GenerateAccAddressForTest()
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3"))

				query := abci.RequestQuery{
					Data:   key,
					Path:   "/store/bank/key",
					Height: ctx.BlockHeight() - 2,
					Prove:  true,
				}

				resp := appA.BaseApp.Query(query)

				return &types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        s.chainB.ChainID,
					SrcZone:     s.chainA.ChainID,
					ClaimType:   cmtypes.ClaimTypeLiquidToken,
					Proofs: []*cmtypes.Proof{
						{
							Key:       resp.Key,
							Data:      resp.Value,
							ProofOps:  resp.ProofOps,
							Height:    resp.Height,
							ProofType: "bank",
						},
					},
				}
			},
			&types.MsgSubmitClaimResponse{},
			"",
			[]cmtypes.Claim{},
		},
		{
			"local_callback_value_invalid_denom",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))
				s.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", address, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))
			},
			func(ctx sdk.Context, appA *app.Quicksilver) *types.MsgSubmitClaim {
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("uqatom"))

				query := abci.RequestQuery{
					Data:   key,
					Path:   "/store/bank/key",
					Height: ctx.BlockHeight() - 2,
					Prove:  true,
				}

				resp := appA.BaseApp.Query(query)

				return &types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        s.chainB.ChainID,
					SrcZone:     s.chainA.ChainID,
					ClaimType:   cmtypes.ClaimTypeLiquidToken,
					Proofs: []*cmtypes.Proof{
						{
							Key:       resp.Key,
							Data:      resp.Value,
							ProofOps:  resp.ProofOps,
							Height:    resp.Height,
							ProofType: "bank",
						},
					},
				}
			},
			&types.MsgSubmitClaimResponse{},
			"",
			[]cmtypes.Claim{},
		},
		{
			"local_callback_value_valid_denom",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				s.Require().NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))
				s.Require().NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", address, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))

				// add uqatom to the list of allowed denoms for this zone
				rawPd := types.LiquidAllowedDenomProtocolData{
					ChainID:               s.chainA.ChainID,
					IbcDenom:              "uqatom",
					QAssetDenom:           "uqatom",
					RegisteredZoneChainID: s.chainB.ChainID,
				}

				blob, err := json.Marshal(rawPd)
				s.Require().NoError(err)
				pd := types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)

				appA.ParticipationRewardsKeeper.SetProtocolData(ctx, rawPd.GenerateKey(), pd)
			},
			func(ctx sdk.Context, appA *app.Quicksilver) *types.MsgSubmitClaim {
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("uqatom"))

				query := abci.RequestQuery{
					Data:   key,
					Path:   "/store/bank/key",
					Height: ctx.BlockHeight() - 2,
					Prove:  true,
				}

				resp := appA.BaseApp.Query(query)

				return &types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        s.chainB.ChainID,
					SrcZone:     s.chainA.ChainID,
					ClaimType:   cmtypes.ClaimTypeLiquidToken,
					Proofs: []*cmtypes.Proof{
						{
							Key:       resp.Key,
							Data:      resp.Value,
							ProofOps:  resp.ProofOps,
							Height:    resp.Height,
							ProofType: "bank",
						},
					},
				}
			},
			&types.MsgSubmitClaimResponse{},
			"",
			[]cmtypes.Claim{{
				UserAddress:   address.String(),
				ChainId:       s.chainB.ChainID,
				Module:        cmtypes.ClaimTypeLiquidToken,
				SourceChainId: s.chainA.ChainID,
				Amount:        100,
			}},
		},
	}
	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()

			appA := s.GetQuicksilverApp(s.chainA)
			// override disabled proof verification; lets test actual proofs :)
			appA.ParticipationRewardsKeeper.ValidateProofOps = utils.ValidateProofOps
			appA.ParticipationRewardsKeeper.ValidateSelfProofOps = utils.ValidateSelfProofOps

			s.coordinator.CommitNBlocks(s.chainA, 3)
			ctx := s.chainA.GetContext()
			tt.malleate(ctx, appA)
			s.coordinator.CommitNBlocks(s.chainA, 3)
			ctx = s.chainA.GetContext()
			s.Require().NoError(appA.ClaimsManagerKeeper.StoreSelfConsensusState(ctx, "epoch"))
			s.coordinator.CommitNBlocks(s.chainA, 1)

			ctx = s.chainA.GetContext()
			msg = tt.generate(ctx, appA)
			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)
			appA.ParticipationRewardsKeeper.CallbackHandler().RegisterCallbacks()

			k := keeper.NewMsgServerImpl(appA.ParticipationRewardsKeeper)
			resp, err := k.SubmitClaim(sdk.WrapSDKContext(ctx), msg)
			if tt.wantErr != "" {
				s.Require().Errorf(err, tt.wantErr)
				s.Require().Nil(resp)
				s.T().Logf("Error: %v", err)
				return
			}

			for _, expectedClaim := range tt.claims {
				actualClaim, found := appA.ClaimsManagerKeeper.GetClaim(ctx, expectedClaim.ChainId, expectedClaim.UserAddress, expectedClaim.Module, expectedClaim.SourceChainId)
				s.Require().True(found)
				s.Require().Equal(expectedClaim.Amount, actualClaim.Amount)

			}

			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.want, resp)
		})
	}
}
