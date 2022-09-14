package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/keeper"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

func (suite *KeeperTestSuite) Test_msgServer_Claim() {
	appA := suite.GetQuicksilverApp(suite.chainA)
	// appB := suite.GetQuicksilverApp(suite.chainB)

	denom := "stake"
	userAddress := "osmo1pgfzn0zhxjjgte7hprwtnqyhrn534lqka2dkuu"

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
					Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
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
					Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
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
					Txhash:  "T5_01",
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
					Intents: map[string]*icstypes.ValidatorIntent{
						valAddress.String(): {
							ValoperAddress: valAddress.String(),
							Weight:         sdk.OneDec(),
						},
					},
				}
				appA.InterchainstakingKeeper.SetIntent(
					suite.chainA.GetContext(),
					zone,
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
				prop := gov.Proposal{
					ProposalId: 0,
					Status:     gov.StatusPassed,
				}
				appA.GovKeeper.SetProposal(suite.chainA.GetContext(), prop)

				vote := gov.Vote{
					ProposalId: 0,
					Voter:      userAddress,
					Option:     gov.OptionYes,
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
					Proofs: []*types.Proof{
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
					Proofs: []*types.Proof{
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
							Height: 0,
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
			suite.Require().Equal(resp, tt.want)
		})
	}
}
