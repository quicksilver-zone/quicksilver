package keeper_test

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/utils/randomutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var delegatorAddress = "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure"

func (suite *KeeperTestSuite) TestKeeper_Zones() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
			},
			&types.QueryZonesRequest{},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.Zones(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Zones))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ZoneValidators() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
			},
			&types.QueryZoneValidatorsRequest{},
			false,
			4,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.ZoneValidators(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Validators))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_DepositAccount() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
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
				ChainId: suite.chainB.ChainID,
			},
			false,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.DepositAccount(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_DelegatorIntent() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
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
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			0,
		},
		{
			"DelegatorIntent_No_Delegator_Intents",
			func() {},
			&types.QueryDelegatorIntentRequest{
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			0,
		},
		{
			"DelegatorIntent_Valid_Intents",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				// give funds
				suite.giveFunds(ctx, zone.LocalDenom, 5000000, testAddress)
				// set intents
				// TODO: set standardized intents for keeper_test package
				intents := []types.DelegatorIntent{
					{
						Delegator: testAddress,
						Intents: types.ValidatorIntents{
							&types.ValidatorIntent{
								ValoperAddress: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
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
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.DelegatorIntent(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Intent.Intents))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_DelegatorIntents() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	tests := []struct {
		name     string
		malleate func()
		req      *types.QueryDelegatorIntentsRequest
		wantErr  bool
		verify   func(delegation []*types.DelegatorIntentsResponse)
	}{
		{
			name:     "DelegatorIntent_Nil_Request",
			malleate: func() {},
			req:      nil,
			wantErr:  true,
			verify: func([]*types.DelegatorIntentsResponse) {
			},
		},
		{
			"DelegatorIntent_No_Delegator_Intents",
			func() {
			},
			&types.QueryDelegatorIntentsRequest{
				DelegatorAddress: testAddress,
			},
			false,
			func(intents []*types.DelegatorIntentsResponse) {
				for _, intent := range intents {
					suite.Require().Equal(intent.ChainId, suite.chainB.ChainID)
					suite.Require().Equal(len(intent.Intent.Intents), 0)
				}
			},
		},
		{
			"DelegatorIntent_Valid_Intents across multiple zones",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				// give funds
				suite.giveFunds(ctx, zone.LocalDenom, 5000000, testAddress)
				// set intents
				intents := []types.DelegatorIntent{
					{
						Delegator: testAddress,
						Intents: types.ValidatorIntents{
							&types.ValidatorIntent{
								ValoperAddress: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
								Weight:         sdk.OneDec(),
							},
						},
					},
				}
				for _, intent := range intents {
					icsKeeper.SetDelegatorIntent(ctx, &zone, intent, false)
				}

				// cosmos zone
				zone = types.Zone{
					ConnectionId:    "connection-77001",
					ChainId:         "cosmoshub-4",
					AccountPrefix:   "cosmos",
					LocalDenom:      "uqatom",
					BaseDenom:       "uatom",
					MultiSend:       false,
					LiquidityModule: false,
					Is_118:          true,
				}
				icsKeeper.SetZone(ctx, &zone)
				// give funds
				suite.giveFunds(ctx, zone.LocalDenom, 5000000, testAddress)
				// set intents
				intents = []types.DelegatorIntent{
					{
						Delegator: testAddress,
						Intents: types.ValidatorIntents{
							&types.ValidatorIntent{
								ValoperAddress: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
								Weight:         sdk.OneDec(),
							},
						},
					},
				}
				for _, intent := range intents {
					icsKeeper.SetDelegatorIntent(ctx, &zone, intent, false)
				}
			},
			&types.QueryDelegatorIntentsRequest{
				DelegatorAddress: testAddress,
			},
			false,
			func(intents []*types.DelegatorIntentsResponse) {
				suite.Require().Equal(len(intents), 2)
				suite.Require().Equal(intents[0].ChainId, "cosmoshub-4")
				suite.Require().Equal(intents[1].ChainId, suite.chainB.ChainID)
				for _, intent := range intents {
					suite.Require().Equal(intent.Intent.Delegator, testAddress)
					suite.Require().Equal(len(intent.Intent.Intents), 1)
				}
			},
		},
	}

	// run tests:
	suite.setupTestZones()

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.DelegatorIntents(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			tt.verify(resp.Intents)

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_Delegations() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
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
				ChainId: suite.chainB.ChainID,
			},
			false,
			0,
		},
		{
			"Delegations_Valid_Delegations",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)

				// set delegation
				// TODO: set standardized delegations for keeper_test package
				delegation := types.Delegation{
					DelegationAddress: testAddress,
					ValidatorAddress:  icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
					Amount:            sdk.NewCoin("denom", sdk.NewInt(15000)),
				}
				icsKeeper.SetDelegation(ctx, &zone, delegation)
			},
			&types.QueryDelegationsRequest{
				ChainId: suite.chainB.ChainID,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.Delegations(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Delegations))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_Receipts() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
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
				ChainId: suite.chainB.ChainID,
			},
			false,
			0,
		},
		{
			"Receipts_Valid_Receipts",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)

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
				ChainId: suite.chainB.ChainID,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.Receipts(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Receipts))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_TxStatus() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()
	suite.setupTestZones()

	testReceiptHash := "testReceiptHash#01"

	zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	testReceipt := icsKeeper.NewReceipt(
		ctx,
		&zone,
		testAddress,
		testReceiptHash,
		sdk.NewCoins(
			sdk.NewCoin(zone.BaseDenom, math.NewInt(50000000)),
		),
	)

	tests := []struct {
		name     string
		malleate func()
		req      *types.QueryTxStatusRequest
		want     *types.QueryTxStatusResponse
		wantErr  bool
	}{
		{
			"Nil_Request",
			func() {},
			nil,
			nil,
			true,
		},
		{
			"empty_TxHash",
			func() {},
			&types.QueryTxStatusRequest{
				ChainId: suite.chainB.ChainID,
				TxHash:  "",
			},
			nil,
			true,
		},
		{
			"Invalid_Zone",
			func() {},
			&types.QueryTxStatusRequest{
				ChainId: "boguschain",
				TxHash:  "unimportant",
			},
			nil,
			true,
		},
		{
			name:     "Receipts_No_Zone_Receipts",
			malleate: func() {},
			req: &types.QueryTxStatusRequest{
				ChainId: suite.chainB.ChainID,
				TxHash:  "randomhash",
			},
			want:    nil,
			wantErr: true,
		},
		{
			"Receipts_Valid_Receipts",
			func() {
				icsKeeper.SetReceipt(ctx, *testReceipt)
			},
			&types.QueryTxStatusRequest{
				ChainId: suite.chainB.ChainID,
				TxHash:  testReceiptHash,
			},
			&types.QueryTxStatusResponse{Receipt: testReceipt},
			false,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.TxStatus(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().EqualValues(tt.want, resp)

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ZoneWithdrawalRecords() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
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
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: delegatorAddress,
			},
			false,
			0,
		},
		{
			"ZoneWithdrawalRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)

				distribution := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  10000000,
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  20000000,
					},
				}

				// set records
				icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distribution,
					testAddress,
					sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(15000000))),
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
					time.Time{},
				)
			},
			&types.QueryWithdrawalRecordsRequest{
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: delegatorAddress,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.ZoneWithdrawalRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_UserWithdrawalRecords() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryUserWithdrawalRecordsRequest
		wantErr      bool
		expectLength int
	}{
		{
			"UserWithdrawalRecords_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"UserWithdrawalRecords_Invalid_Address",
			func() {
				// setup zones
				suite.setupTestZones()
			},
			&types.QueryUserWithdrawalRecordsRequest{
				UserAddress: "incorrect address",
			},
			true,
			0,
		},
		{
			"UserWithdrawalRecords_No_Withdrawal_Records",
			func() {},
			&types.QueryUserWithdrawalRecordsRequest{
				UserAddress: testAddress,
			},
			false,
			0,
		},
		{
			"UserWithdrawalRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)

				distribution := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  10000000,
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  20000000,
					},
				}

				// set records
				icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distribution,
					testAddress,
					sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(15000000))),
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
					time.Time{},
				)
			},
			&types.QueryUserWithdrawalRecordsRequest{
				UserAddress: delegatorAddress,
			},
			false,
			1,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.UserWithdrawalRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_WithdrawalRecords() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
			},
			&types.QueryWithdrawalRecordsRequest{},
			false,
			0,
		},
		{
			"WithdrawalRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)

				distribution := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  10000000,
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  20000000,
					},
				}

				// set records
				icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distribution,
					testAddress,
					sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(15000000))),
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
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
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.WithdrawalRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_UnbondingRecords() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
			},
			&types.QueryUnbondingRecordsRequest{},
			false,
			0,
		},
		{
			"UnbondingRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)

				icsKeeper.SetUnbondingRecord(
					ctx,
					types.UnbondingRecord{
						ChainId:       zone.ChainId,
						EpochNumber:   1,
						Validator:     icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
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
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.UnbondingRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Unbondings))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_RedelegationRecords() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

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
				suite.setupTestZones()
			},
			&types.QueryRedelegationRecordsRequest{},
			false,
			0,
		},
		{
			"RedelegationRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)

				icsKeeper.SetRedelegationRecord(
					ctx,
					types.RedelegationRecord{
						ChainId:     zone.ChainId,
						EpochNumber: 1,
						Source:      icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Destination: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
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
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.RedelegationRecords(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.Redelegations))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_MappedAccounts() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	usrAddress1, _ := addressutils.AccAddressFromBech32("cosmos1vwh8mkgefn73vpsv7td68l3tynayck07engahn", "cosmos")
	ctx := suite.chainA.GetContext()

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
			&types.QueryMappedAccountsRequest{Address: "cosmos1vwh8mkgefn73vpsv7td68l3tynayck07engahn"},
			false,
			0,
		},
		{
			"MappedAccounts_ValidRecord_Request",
			func() {
				// setup zones
				suite.setupTestZones()
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

				icsKeeper.SetRemoteAddressMap(ctx, usrAddress1, randomutils.GenerateRandomBytes(32), zone.ChainId)
			},
			&types.QueryMappedAccountsRequest{Address: "cosmos1vwh8mkgefn73vpsv7td68l3tynayck07engahn"},
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

				icsKeeper.SetRemoteAddressMap(ctx, usrAddress1, randomutils.GenerateRandomBytes(32), zone.ChainId)

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

				icsKeeper.SetRemoteAddressMap(ctx, usrAddress1, randomutils.GenerateRandomBytes(32), zone2.ChainId)
			},
			&types.QueryMappedAccountsRequest{Address: "cosmos1vwh8mkgefn73vpsv7td68l3tynayck07engahn"},
			false,
			2,
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.MappedAccounts(
				ctx,
				tt.req,
			)
			if tt.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.expectLength, len(resp.RemoteAddressMap))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.Require().NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}
