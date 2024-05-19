package keeper_test

import (
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/airdrop/keeper"
	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	minttypes "github.com/quicksilver-zone/quicksilver/x/mint/types"
)

func (suite *KeeperTestSuite) Test_msgServer_Claim() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	userAddress := addressutils.GenerateAccAddressForTest().String()
	denom := "uatom" // same as test zone setup in keeper_test

	msg := types.MsgClaim{}
	tests := []struct {
		name     string
		malleate func()
		want     *types.MsgClaimResponse
		wantErr  bool
	}{
		{
			"blank",
			func() {},
			nil,
			true,
		},
		{
			"nodrop",
			func() {
				// no zone airdrop state

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionInitialClaim),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"noclaimrecord",
			func() {
				// set valid zone airdrop state
				suite.initTestZoneDrop()

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionInitialClaim),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			&types.MsgClaimResponse{
				Amount: 0,
			},
			true,
		},
		{
			"claim_initial",
			func() {
				// use existing state (from prev test)

				// add claim record
				cr := types.ClaimRecord{
					ChainId:          suite.chainB.ChainID,
					Address:          userAddress,
					ActionsCompleted: nil,
					MaxAllocation:    100000000,
					BaseValue:        10000000,
				}

				suite.setClaimRecord(cr)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionInitialClaim),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			&types.MsgClaimResponse{
				Amount: 15000000,
			},
			false,
		},
		{
			"claim_deposit_T5_insufficient",
			func() {
				// use existing state (from prev test)

				// add deposit
				rcpt := icstypes.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000), // 20% deposit
						),
					),
				}
				appA.InterchainstakingKeeper.SetReceipt(
					suite.chainA.GetContext(),
					rcpt,
				)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionDepositT5),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_deposit_T4_insufficient",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionDepositT4),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_deposit_T3",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionDepositT3),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			// expects 21% of MaxAllocation
			// - 8% for T3
			// - 7% for T2
			// - 6% for T1
			&types.MsgClaimResponse{
				Amount: 21000000,
			},
			false,
		},
		{
			"claim_deposit_T2_completed",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionDepositT2),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_deposit_T1_completed",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionDepositT1),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_deposit_T5",
			func() {
				// use existing state (from prev test)

				// add deposit
				rcpt := icstypes.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "T5_02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(1000000), // 10% deposit (sum 30%)
						),
					),
				}
				appA.InterchainstakingKeeper.SetReceipt(
					suite.chainA.GetContext(),
					rcpt,
				)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionDepositT5),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			// expects 19% of MaxAllocation
			// - 9% for T4
			// - 10% for T5
			// - 21% for T1-T3 already claimed
			&types.MsgClaimResponse{
				Amount: 19000000,
			},
			false,
		},
		{
			"claim_stake_no_bond",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionStakeQCK),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_stake_with_bond",
			func() {
				// use existing state (from prev test)

				// add staking delegation
				valAddress, err := sdk.ValAddressFromHex(suite.chainA.Vals.Validators[2].Address.String())
				suite.Require().NoError(err)

				del := staking.Delegation{
					DelegatorAddress: userAddress,
					ValidatorAddress: valAddress.String(),
					Shares:           sdk.MustNewDecFromStr("10.0"),
				}
				appA.StakingKeeper.SetDelegation(
					suite.chainA.GetContext(),
					del,
				)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionStakeQCK),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			&types.MsgClaimResponse{
				Amount: 15000000,
			},
			false,
		},
		{
			"claim_intent_not_set",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionSignalIntent),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_intent_is_set",
			func() {
				// use existing state (from prev test)

				// add intent
				valAddress, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[1].Address.String())
				suite.Require().NoError(err)

				zone, found := appA.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
				suite.Require().True(found)

				intent := icstypes.DelegatorIntent{
					Delegator: userAddress,
					Intents: []*icstypes.ValidatorIntent{
						{
							ValoperAddress: valAddress.String(),
							Weight:         sdk.OneDec(),
						},
					},
				}
				appA.InterchainstakingKeeper.SetDelegatorIntent(
					suite.chainA.GetContext(),
					&zone,
					intent,
					false,
				)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionSignalIntent),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			&types.MsgClaimResponse{
				Amount: 5000000,
			},
			false,
		},
		{
			"claim_gov_no_votes",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionQSGov),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_gov_voted",
			func() {
				// use existing state (from prev test)

				// add proposal and vote
				prop := govv1.Proposal{
					Id:     0,
					Status: govv1.StatusPassed,
				}
				appA.GovKeeper.SetProposal(suite.chainA.GetContext(), prop)

				vote := govv1.Vote{
					ProposalId: 0,
					Voter:      userAddress,
					Options: []*govv1.WeightedVoteOption{
						{
							Option: govv1.VoteOption_VOTE_OPTION_YES,
							Weight: "1.0",
						},
					},
				}
				appA.GovKeeper.SetVote(suite.chainA.GetContext(), vote)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionQSGov),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			&types.MsgClaimResponse{
				Amount: 10000000,
			},
			false,
		},
		{
			"claim_gbp_invalid",
			func() {
				// TODO: implement with GbP
			},
			nil,
			true,
		},
		{
			"claim_gbp_valid",
			func() {
				// TODO: implement with GbP
			},
			nil,
			true,
		},
		{
			"claim_osmosis_lp_nilproofs",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionOsmosis),
					Address: userAddress,
					Proofs:  nil,
				}
			},
			nil,
			true,
		},
		{
			"claim_osmosis_lp_zeroproof",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionOsmosis),
					Address: userAddress,
					Proofs: []*cmtypes.Proof{
						{},
					},
				}
			},
			nil,
			true,
		},
		{
			"claim_osmosis_lp_invalidproof",
			func() {
				// use existing state (from prev test)

				// add consensus state
				// consensusState := exported.ConsensusState{}

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  int64(types.ActionOsmosis),
					Address: userAddress,
					Proofs: []*cmtypes.Proof{
						{
							Key:  []byte(string("123")),
							Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
							ProofOps: &crypto.ProofOps{
								Ops: []crypto.ProofOp{
									{
										Type: "testtype",
										Key:  []byte(string("123")),
										Data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
									},
								},
							},
							ProofType: "lockup",
							Height:    0,
						},
					},
				}
			},
			nil,
			true,
		},
		{
			"Undefined action",
			func() {
				// use existing state (from prev test)

				msg = types.MsgClaim{
					ChainId: suite.chainB.ChainID,
					Action:  999,
					Address: userAddress,
					Proofs:  nil,
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

			k := keeper.NewMsgServerImpl(appA.AirdropKeeper)
			resp, err := k.Claim(sdk.WrapSDKContext(suite.chainA.GetContext()), &msg)
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

func (suite *KeeperTestSuite) Test_msgServer_IncentivePoolSpend() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	modAccAddr := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
	userAddress := addressutils.GenerateAccAddressForTest().String()
	denom := "uatom" // same as test zone setup in keeper_test
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromUint64(1000)))
	mintCoins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromUint64(100000000)))

	// set up mod acct with funds
	err := appA.BankKeeper.MintCoins(suite.chainA.GetContext(), minttypes.ModuleName, mintCoins)
	suite.Require().NoError(err)
	err = appA.BankKeeper.SendCoinsFromModuleToModule(suite.chainA.GetContext(), minttypes.ModuleName, types.ModuleName, mintCoins)
	suite.Require().NoError(err)

	msg := types.MsgIncentivePoolSpend{}
	tests := []struct {
		name     string
		malleate func()
		want     *types.MsgIncentivePoolSpendResponse
		wantErr  bool
	}{
		{
			name: "invalid authority",
			malleate: func() {
				msg = types.MsgIncentivePoolSpend{
					Authority:   "invalid",
					ToAddress:   userAddress,
					Amount:      coins,
					Title:       "Invalid Incentive Pool Spend Title",
					Description: "Invalid Incentive Pool Spend Description",
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid",
			malleate: func() {
				msg = types.MsgIncentivePoolSpend{
					Authority:   modAccAddr,
					ToAddress:   userAddress,
					Amount:      coins,
					Title:       "Valid Incentive Pool Spend Title",
					Description: "Valid Incentive Pool Spend Description",
				}
			},
			want:    &types.MsgIncentivePoolSpendResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			tt.malleate()

			k := keeper.NewMsgServerImpl(appA.AirdropKeeper)
			resp, err := k.IncentivePoolSpend(sdk.WrapSDKContext(suite.chainA.GetContext()), &msg)
			if tt.wantErr {
				suite.Require().Error(err)
				suite.Require().Nil(resp)
				suite.T().Logf("Error: %v", err)
				return
			}

			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.want, resp)

			// verify that balance has been properly transferred
			accAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
			suite.Require().NoError(err)
			balance := appA.BankKeeper.GetAllBalances(suite.chainA.GetContext(), accAddr)
			suite.Require().Equal(msg.Amount, balance)
		})
	}
}
