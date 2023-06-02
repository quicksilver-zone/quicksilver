package keeper_test

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (s *KeeperTestSuite) TestKeeper_Zones() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryZonesRequest
		wantErr      bool
		expectLength int
	}{
		{
			"Zones_No_State",
			func() {},
			&types.QueryZonesRequest{},
			false,
			0,
		},
		{
			"Zones_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"Zones_Valid_Request",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryZonesRequest{},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.Zones(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Zones))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_ZoneValidators() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryZoneValidatorsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"ZoneValidatorsInfo_No_State",
			func() {},
			&types.QueryZoneValidatorsRequest{},
			false,
			0,
		},
		{
			"ZoneValidatorsInfo_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"ZoneValidatorsInfo_Valid_Request",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryZoneValidatorsRequest{},
			false,
			4,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.ZoneValidators(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Validators))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_DepositAccount() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name     string
		malleate func()
		req      *types.QueryDepositAccountForChainRequest
		wantErr  bool
	}{
		{
			"DepositAccount_No_State",
			func() {},
			&types.QueryDepositAccountForChainRequest{},
			true,
		},
		{
			"DepositAccount_Nil_Request",
			func() {},
			nil,
			true,
		},
		{
			"DepositAccount_Invalid_Request",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryDepositAccountForChainRequest{},
			true,
		},
		{
			"DepositAccount_Valid_Request",
			func() {
				// use state set from previous tests
			},
			&types.QueryDepositAccountForChainRequest{
				ChainId: s.chainB.ChainID,
			},
			false,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.DepositAccount(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_DelegatorIntent() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryDelegatorIntentRequest
		wantErr      bool
		expectLength int
	}{
		{
			"DelegatorIntent_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"DelegatorIntent_Invalid_Zone",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryDelegatorIntentRequest{
				ChainId:          "boguschain",
				DelegatorAddress: testAddress,
			},
			true,
			0,
		},
		{
			"DelegatorIntent_No_Zone_Intents",
			func() {},
			&types.QueryDelegatorIntentRequest{
				ChainId:          s.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			0,
		},
		{
			"DelegatorIntent_No_Delegator_Intents",
			func() {},
			&types.QueryDelegatorIntentRequest{
				ChainId:          s.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			0,
		},
		{
			"DelegatorIntent_Valid_Intents",
			func() {
				zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)
				// give funds
				s.giveFunds(ctx, zone.LocalDenom, 5000000, testAddress)
				// set intents
				// TODO: set standardized intents for keeper_test package
				intents := []types.DelegatorIntent{
					{
						Delegator: testAddress,
						Intents: types.ValidatorIntents{
							&types.ValidatorIntent{
								ValoperAddress: icsKeeper.GetValidators(ctx, s.chainB.ChainID)[0].ValoperAddress,
								Weight:         sdk.OneDec(),
							},
						},
					},
				}
				for _, intent := range intents {
					icsKeeper.SetDelegatorIntent(ctx, &zone, intent, false)
				}
			},
			&types.QueryDelegatorIntentRequest{
				ChainId:          s.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.DelegatorIntent(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Intent.Intents))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_Delegations() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryDelegationsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"Delegations_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"Delegations_Invalid_Zone",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryDelegationsRequest{
				ChainId: "boguschain",
			},
			true,
			0,
		},
		{
			"Delegations_No_Zone_Delegations",
			func() {},
			&types.QueryDelegationsRequest{
				ChainId: s.chainB.ChainID,
			},
			false,
			0,
		},
		{
			"Delegations_Valid_Delegations",
			func() {
				zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)

				// set delegation
				// TODO: set standardized delegations for keeper_test package
				delegation := types.Delegation{
					DelegationAddress: testAddress,
					ValidatorAddress:  icsKeeper.GetValidators(ctx, s.chainB.ChainID)[0].ValoperAddress,
					Amount:            sdk.NewCoin("denom", sdk.NewInt(15000)),
				}
				icsKeeper.SetDelegation(ctx, &zone, delegation)
			},
			&types.QueryDelegationsRequest{
				ChainId: s.chainB.ChainID,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.Delegations(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Delegations))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_Receipts() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryReceiptsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"Receipts_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"Receipts_Invalid_Zone",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryReceiptsRequest{
				ChainId: "boguschain",
			},
			true,
			0,
		},
		{
			"Receipts_No_Zone_Receipts",
			func() {},
			&types.QueryReceiptsRequest{
				ChainId: s.chainB.ChainID,
			},
			false,
			0,
		},
		{
			"Receipts_Valid_Receipts",
			func() {
				zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)

				// set receipts
				receipt := icsKeeper.NewReceipt(
					ctx,
					&zone,
					testAddress,
					"testReceiptHash#01",
					sdk.NewCoins(
						sdk.NewCoin(zone.BaseDenom, math.NewInt(50000000)),
					),
				)
				icsKeeper.SetReceipt(ctx, *receipt)
			},
			&types.QueryReceiptsRequest{
				ChainId: s.chainB.ChainID,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.Receipts(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Receipts))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_ZoneWithdrawalRecords() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryWithdrawalRecordsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"ZoneWithdrawalRecords_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"ZoneWithdrawalRecords_Invalid_Zone",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryWithdrawalRecordsRequest{
				ChainId: "boguschain",
			},
			true,
			0,
		},
		{
			"ZoneWithdrawalRecords_No_Zone_Records",
			func() {},
			&types.QueryWithdrawalRecordsRequest{
				ChainId:          s.chainB.ChainID,
				DelegatorAddress: "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
			},
			false,
			0,
		},
		{
			"ZoneWithdrawalRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)

				distribution := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, s.chainB.ChainID)[0].ValoperAddress,
						Amount:  10000000,
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, s.chainB.ChainID)[1].ValoperAddress,
						Amount:  20000000,
					},
				}

				// set records
				icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					"quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
					distribution,
					testAddress,
					sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(15000000))),
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					icskeeper.WithdrawStatusQueued,
					time.Time{},
				)
			},
			&types.QueryWithdrawalRecordsRequest{
				ChainId:          s.chainB.ChainID,
				DelegatorAddress: "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.ZoneWithdrawalRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_WithdrawalRecords() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryWithdrawalRecordsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"WithdrawalRecords_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"WithdrawalRecords_No_Zone_Records",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryWithdrawalRecordsRequest{},
			false,
			0,
		},
		{
			"WithdrawalRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)

				distribution := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, s.chainB.ChainID)[0].ValoperAddress,
						Amount:  10000000,
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, s.chainB.ChainID)[1].ValoperAddress,
						Amount:  20000000,
					},
				}

				// set records
				icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					"quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
					distribution,
					testAddress,
					sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(15000000))),
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					icskeeper.WithdrawStatusQueued,
					time.Time{},
				)
			},
			&types.QueryWithdrawalRecordsRequest{},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.WithdrawalRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_UnbondingRecords() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryUnbondingRecordsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"UnbondingRecords_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"UnbondingRecords_No_Zone_Records",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryUnbondingRecordsRequest{},
			false,
			0,
		},
		{
			"UnbondingRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)

				icsKeeper.SetUnbondingRecord(
					ctx,
					types.UnbondingRecord{
						ChainId:       zone.ChainId,
						EpochNumber:   1,
						Validator:     icsKeeper.GetValidators(ctx, s.chainB.ChainID)[0].ValoperAddress,
						RelatedTxhash: []string{"ABC012"},
					},
				)
			},
			&types.QueryUnbondingRecordsRequest{},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.UnbondingRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Unbondings))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_RedelegationRecords() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryRedelegationRecordsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"RedelegationRecords_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"RedelegationRecords_No_Zone_Records",
			func() {
				// setup zones
				s.setupTestZones()
			},
			&types.QueryRedelegationRecordsRequest{},
			false,
			0,
		},
		{
			"RedelegationRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
				s.Require().True(found)

				icsKeeper.SetRedelegationRecord(
					ctx,
					types.RedelegationRecord{
						ChainId:     zone.ChainId,
						EpochNumber: 1,
						Source:      icsKeeper.GetValidators(ctx, s.chainB.ChainID)[1].ValoperAddress,
						Destination: icsKeeper.GetValidators(ctx, s.chainB.ChainID)[0].ValoperAddress,
						Amount:      10000000,
					})
			},
			&types.QueryRedelegationRecordsRequest{},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.RedelegationRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.Redelegations))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (s *KeeperTestSuite) TestKeeper_MappedAccounts() {
	icsKeeper := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper
	usrAddress1, _ := utils.AccAddressFromBech32("quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq", "quick")
	ctx := s.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryMappedAccountsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"MappedAccounts_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"MappedAccounts_NoRecords_Request",
			func() {
				// setup zones
				zone := types.Zone{
					ConnectionId:    "connection-77001",
					ChainId:         "evmos_9001-1",
					AccountPrefix:   "evmos",
					LocalDenom:      "uqevmos",
					BaseDenom:       "uevmos",
					MultiSend:       false,
					LiquidityModule: false,
					Is_118:          false,
				}
				icsKeeper.SetZone(ctx, &zone)
			},
			&types.QueryMappedAccountsRequest{Address: "quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq"},
			false,
			0,
		},
		{
			"MappedAccounts_ValidRecord_Request",
			func() {
				// setup zones
				s.setupTestZones()
				zone := types.Zone{
					ConnectionId:    "connection-77881",
					ChainId:         "evmos_9001-1",
					AccountPrefix:   "evmos",
					LocalDenom:      "uqevmos",
					BaseDenom:       "uevmos",
					MultiSend:       false,
					LiquidityModule: false,
					Is_118:          false,
				}
				icsKeeper.SetZone(ctx, &zone)

				icsKeeper.SetRemoteAddressMap(ctx, usrAddress1, utils.GenerateRandomHash(), zone.ChainId)
			},
			&types.QueryMappedAccountsRequest{Address: "quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq"},
			false,
			1,
		},

		{
			"MappedAccounts_ValidMultipleRecord_Request",
			func() {
				// setup zones
				zone := types.Zone{
					ConnectionId:    "connection-77881",
					ChainId:         "evmos_9001-1",
					AccountPrefix:   "evmos",
					LocalDenom:      "uqevmos",
					BaseDenom:       "uevmos",
					MultiSend:       false,
					LiquidityModule: false,
					Is_118:          false,
				}
				icsKeeper.SetZone(ctx, &zone)

				icsKeeper.SetRemoteAddressMap(ctx, usrAddress1, utils.GenerateRandomHash(), zone.ChainId)

				zone2 := types.Zone{
					ConnectionId:    "connection-77891",
					ChainId:         "injective-1",
					AccountPrefix:   "injective",
					LocalDenom:      "uqinj",
					BaseDenom:       "uinj",
					MultiSend:       false,
					LiquidityModule: false,
					Is_118:          false,
				}
				icsKeeper.SetZone(ctx, &zone2)

				icsKeeper.SetRemoteAddressMap(ctx, usrAddress1, utils.GenerateRandomHash(), zone2.ChainId)
			},
			&types.QueryMappedAccountsRequest{Address: "quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq"},
			false,
			2,
		},
	}

	// run tests:
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.MappedAccounts(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				s.T().Logf("Error:\n%v\n", err)
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tt.expectLength, len(resp.MappedAccounts))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			s.Require().NoError(err)

			s.T().Logf("Response:\n%s\n", vstr)
		})
	}
}
