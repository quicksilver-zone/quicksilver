package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/airdrop/keeper"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (s *KeeperTestSuite) Test_msgServer_Claim() {
	appA := s.GetQuicksilverApp(s.chainA)

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
					ChainId: s.chainB.ChainID,
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
				s.initTestZoneDrop()

				msg = types.MsgClaim{
					ChainId: s.chainB.ChainID,
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
					ChainId:          s.chainB.ChainID,
					Address:          userAddress,
					ActionsCompleted: nil,
					MaxAllocation:    100000000,
					BaseValue:        10000000,
				}

				s.setClaimRecord(cr)

				msg = types.MsgClaim{
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					s.chainA.GetContext(),
					rcpt,
				)

				msg = types.MsgClaim{
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					s.chainA.GetContext(),
					rcpt,
				)

				msg = types.MsgClaim{
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
				valAddress, err := sdk.ValAddressFromHex(s.chainA.Vals.Validators[2].Address.String())
				s.Require().NoError(err)

				del := staking.Delegation{
					DelegatorAddress: userAddress,
					ValidatorAddress: valAddress.String(),
					Shares:           sdk.MustNewDecFromStr("10.0"),
				}
				appA.StakingKeeper.SetDelegation(
					s.chainA.GetContext(),
					del,
				)

				msg = types.MsgClaim{
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
				valAddress, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[1].Address.String())
				s.Require().NoError(err)

				zone, found := appA.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.Require().True(found)

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
					s.chainA.GetContext(),
					&zone,
					intent,
					false,
				)

				msg = types.MsgClaim{
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
				appA.GovKeeper.SetProposal(s.chainA.GetContext(), prop)

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
				appA.GovKeeper.SetVote(s.chainA.GetContext(), vote)

				msg = types.MsgClaim{
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
					ChainId: s.chainB.ChainID,
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
	}
	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			tt.malleate()

			k := keeper.NewMsgServerImpl(appA.AirdropKeeper)
			resp, err := k.Claim(sdk.WrapSDKContext(s.chainA.GetContext()), &msg)
			if tt.wantErr {
				s.Require().Error(err)
				s.Require().Nil(resp)
				s.T().Logf("Error: %v", err)
				return
			}

			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.want, resp)
		})
	}
}

func (s *KeeperTestSuite) Test_msgServer_IncentivePoolSpend() {
	appA := s.GetQuicksilverApp(s.chainA)

	modAccAddr := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
	userAddress := addressutils.GenerateAccAddressForTest().String()
	denom := "uatom" // same as test zone setup in keeper_test
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromUint64(1000)))
	mintCoins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromUint64(100000000)))

	// set up mod acct with funds
	err := appA.BankKeeper.MintCoins(s.chainA.GetContext(), minttypes.ModuleName, mintCoins)
	s.Require().NoError(err)
	err = appA.BankKeeper.SendCoinsFromModuleToModule(s.chainA.GetContext(), minttypes.ModuleName, types.ModuleName, mintCoins)
	s.Require().NoError(err)

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
					Authority: "invalid",
					ToAddress: userAddress,
					Amount:    coins,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid",
			malleate: func() {
				msg = types.MsgIncentivePoolSpend{
					Authority: modAccAddr,
					ToAddress: userAddress,
					Amount:    coins,
				}
			},
			want:    &types.MsgIncentivePoolSpendResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			tt.malleate()

			k := keeper.NewMsgServerImpl(appA.AirdropKeeper)
			resp, err := k.IncentivePoolSpend(sdk.WrapSDKContext(s.chainA.GetContext()), &msg)
			if tt.wantErr {
				s.Require().Error(err)
				s.Require().Nil(resp)
				s.T().Logf("Error: %v", err)
				return
			}

			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.want, resp)

			// verify that balance has been properly transferred
			accAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
			s.Require().NoError(err)
			balance := appA.BankKeeper.GetAllBalances(s.chainA.GetContext(), accAddr)
			s.Require().Equal(msg.Amount, balance)
		})
	}
}
