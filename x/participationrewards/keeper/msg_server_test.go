package keeper_test

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"

	"github.com/quicksilver-zone/quicksilver/app" //nolint:revive
	lockup "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/lockup/types"
	umeetypes "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/keeper"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
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
					UserAddress: addressutils.GenerateAccAddressForTest().String(),
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
				userAddress := addressutils.GenerateAccAddressForTest()
				osmoAddress := addressutils.GenerateAddressForTestWithPrefix("osmo")
				lockedResp := lockup.LockedResponse{
					Lock: &lockup.PeriodLock{
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
				suite.NoError(err)

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
				userAddress := addressutils.GenerateAccAddressForTest()
				osmoAddress := addressutils.MustEncodeAddressToBech32("osmo", userAddress)
				lockedResp := lockup.LockedResponse{
					Lock: &lockup.PeriodLock{
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
				suite.NoError(err)

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
			"invalid_umee_denom",
			func() {
				address := addressutils.GenerateAccAddressForTest()
				prefix := banktypes.CreateAccountBalancesPrefix(authtypes.NewModuleAddress(umeetypes.LeverageModuleName))
				key := banktypes.CreatePrefixedAccountStoreKey(prefix, []byte("u/ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3"))

				cd := sdk.Coin{
					Denom:  "u/random",
					Amount: math.NewInt(1000),
				}
				bz, err := cd.Marshal()
				suite.NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "testchain-1",
					ClaimType:   cmtypes.ClaimTypeUmeeToken,
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
			nil,
			"a",
		},
		{
			"valid_osmosis",
			func() {
				userAddress := addressutils.GenerateAccAddressForTest()
				osmoAddress := addressutils.MustEncodeAddressToBech32("osmo", userAddress)
				locked := &lockup.PeriodLock{
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
				suite.NoError(err)

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
			"valid_osmosis_unbonded_gamm",
			func() {
				userAddress := addressutils.GenerateAccAddressForTest()
				bankkey := banktypes.CreateAccountBalancesPrefix(userAddress)
				bankkey = append(bankkey, "gamm/1"...)

				cd := math.NewInt(10000000)
				bz, err := cd.Marshal()
				suite.Require().NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: userAddress.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "osmosis-1",
					ClaimType:   cmtypes.ClaimTypeOsmosisPool,
					Proofs: []*cmtypes.Proof{
						{
							Key:       bankkey,
							Data:      bz,
							ProofOps:  &crypto.ProofOps{},
							Height:    0,
							ProofType: types.ProofTypeBank,
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
				address := addressutils.GenerateAccAddressForTest()
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3"))

				cd := sdk.Coin{
					Denom:  "",
					Amount: math.NewInt(0),
				}
				bz, err := cd.Marshal()
				suite.NoError(err)

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
				address := addressutils.GenerateAccAddressForTest()
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3"))

				cd := sdk.Coin{
					Denom:  "",
					Amount: math.NewInt(0),
				}
				bz, err := cd.Marshal()
				suite.NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        "cosmoshub-4",
					SrcZone:     "testchain-1",
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
		{
			"valid_umee",
			func() {
				address := addressutils.GenerateAccAddressForTest()

				userAddress, _ := addressutils.EncodeAddressToBech32("umee", address)
				bankkey := banktypes.CreateAccountBalancesPrefix(address)
				bankkey = append(bankkey, "u/uumee"...)

				leveragekey := umeetypes.KeyCollateralAmount(address, "u/uumee")

				cd := math.NewInt(1000)
				bz, err := cd.Marshal()
				suite.NoError(err)

				msg = types.MsgSubmitClaim{
					UserAddress: userAddress,
					Zone:        "cosmoshub-4",
					SrcZone:     "testchain-1",
					ClaimType:   cmtypes.ClaimTypeUmeeToken,
					Proofs: []*cmtypes.Proof{
						{
							Key:       bankkey,
							Data:      bz,
							ProofOps:  &crypto.ProofOps{},
							Height:    10,
							ProofType: "bank",
						},
						{
							Key:       leveragekey,
							Data:      bz,
							ProofOps:  &crypto.ProofOps{},
							Height:    10,
							ProofType: "leverage",
						},
					},
				}
			},
			&types.MsgSubmitClaimResponse{},
			"",
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			ctx := suite.chainA.GetContext()
			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)

			k := keeper.NewMsgServerImpl(appA.ParticipationRewardsKeeper)
			resp, err := k.SubmitClaim(sdk.WrapSDKContext(ctx), &msg)
			if tt.wantErr != "" {
				suite.Errorf(err, tt.wantErr)
				suite.Nil(resp)
				suite.T().Logf("Error: %v", err)
				return
			}

			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.want, resp)
		})
	}
}

func (suite *KeeperTestSuite) Test_msgServer_SubmitLocalClaim() {
	address := addressutils.GenerateAccAddressForTest()

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
				address := addressutils.GenerateAccAddressForTest()
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3"))

				query := abci.RequestQuery{
					Data:   key,
					Path:   "/store/bank/key",
					Height: ctx.BlockHeight() - 2,
					Prove:  true,
				}

				resp := appA.Query(query)

				return &types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        suite.chainB.ChainID,
					SrcZone:     suite.chainA.ChainID,
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
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", address, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))
			},
			func(ctx sdk.Context, appA *app.Quicksilver) *types.MsgSubmitClaim {
				key := banktypes.CreatePrefixedAccountStoreKey(address, []byte("uqatom"))

				query := abci.RequestQuery{
					Data:   key,
					Path:   "/store/bank/key",
					Height: ctx.BlockHeight() - 2,
					Prove:  true,
				}

				resp := appA.Query(query)

				return &types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        suite.chainB.ChainID,
					SrcZone:     suite.chainA.ChainID,
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
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", address, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(100)))))

				// add uqatom to the list of allowed denoms for this zone
				rawPd := types.LiquidAllowedDenomProtocolData{
					ChainID:               suite.chainA.ChainID,
					IbcDenom:              "uqatom",
					QAssetDenom:           "uqatom",
					RegisteredZoneChainID: suite.chainB.ChainID,
				}

				blob, err := json.Marshal(rawPd)
				suite.NoError(err)
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

				resp := appA.Query(query)

				return &types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        suite.chainB.ChainID,
					SrcZone:     suite.chainA.ChainID,
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
				ChainId:       suite.chainB.ChainID,
				Module:        cmtypes.ClaimTypeLiquidToken,
				SourceChainId: suite.chainA.ChainID,
				Amount:        math.NewInt(100),
			}},
		},
		{
			"local_callback_value_valid_denom_overflow",
			func(ctx sdk.Context, appA *app.Quicksilver) {
				suite.NoError(appA.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(1000000000000000000).Mul(sdk.NewInt(1000))))))
				suite.NoError(appA.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", address, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(1000000000000000000).Mul(sdk.NewInt(1000))))))

				// add uqatom to the list of allowed denoms for this zone
				rawPd := types.LiquidAllowedDenomProtocolData{
					ChainID:               suite.chainA.ChainID,
					IbcDenom:              "uqatom",
					QAssetDenom:           "uqatom",
					RegisteredZoneChainID: suite.chainB.ChainID,
				}

				blob, err := json.Marshal(rawPd)
				suite.NoError(err)
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

				resp := appA.Query(query)

				return &types.MsgSubmitClaim{
					UserAddress: address.String(),
					Zone:        suite.chainB.ChainID,
					SrcZone:     suite.chainA.ChainID,
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
				ChainId:       suite.chainB.ChainID,
				Module:        cmtypes.ClaimTypeLiquidToken,
				SourceChainId: suite.chainA.ChainID,
				Amount:        sdk.NewInt(1000000000000000000).Mul(sdk.NewInt(1000)),
			}},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			appA := suite.GetQuicksilverApp(suite.chainA)
			// override disabled proof verification; lets test actual proofs :)
			appA.ParticipationRewardsKeeper.ValidateProofOps = utils.ValidateProofOps
			appA.ParticipationRewardsKeeper.ValidateSelfProofOps = utils.ValidateSelfProofOps

			suite.coordinator.CommitBlock(suite.chainA)
			ctx := suite.chainA.GetContext()
			tt.malleate(ctx, appA)
			suite.coordinator.CommitBlock(suite.chainA)
			suite.NoError(appA.ClaimsManagerKeeper.StoreSelfConsensusState(ctx, "epoch"))
			suite.coordinator.CommitBlock(suite.chainA)

			ctx = suite.chainA.GetContext() // update context to get latest block height
			msg = tt.generate(ctx, appA)
			params := appA.ParticipationRewardsKeeper.GetParams(ctx)
			params.ClaimsEnabled = true
			appA.ParticipationRewardsKeeper.SetParams(ctx, params)
			appA.ParticipationRewardsKeeper.CallbackHandler().RegisterCallbacks()

			k := keeper.NewMsgServerImpl(appA.ParticipationRewardsKeeper)
			resp, err := k.SubmitClaim(sdk.WrapSDKContext(ctx), msg)
			if tt.wantErr != "" {
				suite.Errorf(err, tt.wantErr)
				suite.Nil(resp)
				suite.T().Logf("Error: %v", err)
				return
			}

			for _, expectedClaim := range tt.claims {
				actualClaim, found := appA.ClaimsManagerKeeper.GetClaim(ctx, expectedClaim.ChainId, expectedClaim.UserAddress, expectedClaim.Module, expectedClaim.SourceChainId)
				suite.True(found)
				suite.Equal(expectedClaim.Amount, actualClaim.Amount)

			}

			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.want, resp)
		})
	}
}
