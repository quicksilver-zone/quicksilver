package keeper_test

import (
	"encoding/json"
	"strings"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/utils/randomutils"
	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Zones))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Validators))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Intent.Intents))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
					suite.Equal(intent.ChainId, suite.chainB.ChainID)
					suite.Equal(len(intent.Intent.Intents), 0)
				}
			},
		},
		{
			"DelegatorIntent_Valid_Intents across multiple zones",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)
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
				suite.Equal(len(intents), 2)
				suite.Equal(intents[0].ChainId, "cosmoshub-4")
				suite.Equal(intents[1].ChainId, suite.chainB.ChainID)
				for _, intent := range intents {
					suite.Equal(intent.Intent.Delegator, testAddress)
					suite.Equal(len(intent.Intent.Intents), 1)
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			tt.verify(resp.Intents)

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)

				// set delegation
				// TODO: set standardized delegations for keeper_test package
				delegation := types.Delegation{
					DelegationAddress: testAddress,
					ValidatorAddress:  icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
					Amount:            sdk.NewCoin("denom", sdk.NewInt(15000)),
				}
				icsKeeper.SetDelegation(ctx, zone.ChainId, delegation)
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Delegations))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)

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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Receipts))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
	suite.True(found)

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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.EqualValues(tt.want, resp)

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)

				distributions := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  math.NewInt(10000000),
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  math.NewInt(20000000),
					},
				}

				// set records
				err := icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distributions,
					testAddress,
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
					time.Time{},
					icsKeeper.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch,
				)
				suite.NoError(err)
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)

				distributions := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  math.NewInt(10000000),
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  math.NewInt(20000000),
					},
				}

				// set records
				err := icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distributions,
					testAddress,
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
					time.Time{},
					icsKeeper.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch,
				)
				suite.NoError(err)
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)

				distributions := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  math.NewInt(10000000),
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  math.NewInt(20000000),
					},
				}

				// set records
				err := icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distributions,
					testAddress,
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
					time.Time{},
					icsKeeper.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch,
				)
				suite.NoError(err)
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Withdrawals))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)

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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Unbondings))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.True(found)

				icsKeeper.SetRedelegationRecord(
					ctx,
					types.RedelegationRecord{
						ChainId:     zone.ChainId,
						EpochNumber: 1,
						Source:      icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Destination: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:      math.NewInt(10000000),
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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.Redelegations))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

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
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.expectLength, len(resp.RemoteAddressMap))

			vstr, err := json.MarshalIndent(resp, "", "\t")
			suite.NoError(err)

			suite.T().Logf("Response:\n%s\n", vstr)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_Zone() {
	testCases := []struct {
		name     string
		malleate func()
		req      *types.QueryZoneRequest
		wantErr  bool
	}{
		{
			name:     "empty request",
			malleate: func() {},
			req:      nil,
			wantErr:  true,
		},
		{
			name:     "zone not found",
			malleate: func() {},
			req:      &types.QueryZoneRequest{ChainId: suite.chainB.ChainID},
			wantErr:  true,
		},
		{
			name: "zone valid request",
			malleate: func() {
				suite.SetupTest()
				suite.setupTestZones()
			},
			req:     &types.QueryZoneRequest{ChainId: suite.chainB.ChainID},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()
			icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
			ctx := suite.chainA.GetContext()

			resp, err := icsKeeper.Zone(ctx, tc.req)
			if tc.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.NotNil(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ZoneValidatorDenyList() {
	testCases := []struct {
		name           string
		req            *types.QueryDenyListRequest
		wantErr        bool
		expectedLength int
	}{
		{
			name:           "empty request",
			req:            nil,
			wantErr:        true,
			expectedLength: 0,
		},
		{
			name:           "zone not found",
			req:            &types.QueryDenyListRequest{ChainId: "abcd"},
			wantErr:        false,
			expectedLength: 0,
		},
		{
			name:           "zone valid request",
			req:            &types.QueryDenyListRequest{ChainId: suite.chainB.ChainID},
			wantErr:        false,
			expectedLength: 2,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper

			// Set 2 validators to deny list
			validator1 := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0]
			valAddr, err := sdk.ValAddressFromBech32(validator1.ValoperAddress)
			suite.NoError(err)
			err = icsKeeper.SetZoneValidatorToDenyList(ctx, suite.chainB.ChainID, valAddr)
			suite.NoError(err)

			validator2 := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1]
			valAddr, err = sdk.ValAddressFromBech32(validator2.ValoperAddress)
			suite.NoError(err)
			err = icsKeeper.SetZoneValidatorToDenyList(ctx, suite.chainB.ChainID, valAddr)
			suite.NoError(err)
			denyList, err := icsKeeper.ValidatorDenyList(ctx, tc.req)
			if tc.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Error(err)
				suite.Empty(denyList)
			} else {
				suite.NotNil(denyList)
				if tc.expectedLength == 2 {
					suite.Equal(&types.QueryDenyListResponse{Validators: []string{validator1.ValoperAddress, validator2.ValoperAddress}}, denyList)
				} else {
					suite.Empty(denyList.Validators)
				}

			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_UserZoneWithdrawalRecords() {
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
			"UserZoneWithdrawalRecords_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"UserZoneWithdrawalRecords_Invalid_Address",
			func() {
				// setup zones
				suite.setupTestZones()
			},
			&types.QueryWithdrawalRecordsRequest{
				ChainId:          "boguschain",
				DelegatorAddress: "incorrect address",
			},
			true,
			0,
		},
		{
			"UserZoneWithdrawalRecords_No_Withdrawal_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				tempAddr := addressutils.GenerateAccAddressForTest().String()
				suite.True(found)

				distributions := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  sdk.NewInt(10000000),
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  sdk.NewInt(20000000),
					},
				}

				// set records
				err := icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distributions,
					tempAddr,
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
					time.Time{},
					icsKeeper.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch,
				)
				suite.NoError(err)
			},
			&types.QueryWithdrawalRecordsRequest{
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			0,
		},
		{
			"UserZoneWithdrawalRecords_Valid_Records",
			func() {
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)

				distributions := []*types.Distribution{
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[0].ValoperAddress,
						Amount:  sdk.NewInt(10000000),
					},
					{
						Valoper: icsKeeper.GetValidators(ctx, suite.chainB.ChainID)[1].ValoperAddress,
						Amount:  sdk.NewInt(20000000),
					},
				}

				// set records
				err := icsKeeper.AddWithdrawalRecord(
					ctx,
					zone.ChainId,
					delegatorAddress,
					distributions,
					testAddress,
					sdk.NewCoin(zone.LocalDenom, math.NewInt(15000000)),
					"ABC012",
					types.WithdrawStatusQueued,
					time.Time{},
					icsKeeper.EpochsKeeper.GetEpochInfo(ctx, epochstypes.EpochIdentifierEpoch).CurrentEpoch,
				)
				suite.NoError(err)
			},
			&types.QueryWithdrawalRecordsRequest{
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: testAddress,
			},
			false,
			1,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			tc.malleate()
			icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
			ctx := suite.chainA.GetContext()

			resp, err := icsKeeper.UserZoneWithdrawalRecords(ctx, tc.req)
			if tc.wantErr {
				suite.T().Logf("Error:\n%v\n", err)
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.NotNil(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ClaimedPercentage() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	tests := []struct {
		name     string
		malleate func()
		req      *types.QueryClaimedPercentageRequest
		wantErr  string
		resp     *types.QueryClaimedPercentageResponse
	}{
		{
			"ClaimedPercentage_Nil_Request",
			func() {},
			nil,
			"empty request",
			nil,
		},
		{
			"ClaimedPercentage_Invalid_Zone",
			func() {
				// setup zones
				suite.setupTestZones()
			},
			&types.QueryClaimedPercentageRequest{
				ChainId: "boguschain",
			},
			"no zone found",
			nil,
		},
		{
			"ClaimedPercentage_Valid_Claims",
			func() {
				addr1, addr2, addr3 := addressutils.GenerateAccAddressForTest().String(), addressutils.GenerateAccAddressForTest().String(), addressutils.GenerateAccAddressForTest().String()
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)
				claims := []claimsmanagertypes.Claim{}
				claims = append(claims, claimsmanagertypes.NewClaim(addr1, zone.ChainId, claimsmanagertypes.ClaimTypeOsmosisPool, "", math.NewInt(1000)))
				claims = append(claims, claimsmanagertypes.NewClaim(addr2, zone.ChainId, claimsmanagertypes.ClaimTypeOsmosisPool, "", math.NewInt(2000)))
				claims = append(claims, claimsmanagertypes.NewClaim(addr3, zone.ChainId, claimsmanagertypes.ClaimTypeOsmosisPool, "", math.NewInt(3000)))
				for _, claim := range claims {
					icsKeeper.ClaimsManagerKeeper.SetClaim(ctx, &claim) // #nosec G601
					err := suite.GetQuicksilverApp(suite.chainA).MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, claim.Amount)))
					suite.NoError(err)
				}
			},
			&types.QueryClaimedPercentageRequest{
				ChainId: suite.chainB.ChainID,
			},
			"",
			&types.QueryClaimedPercentageResponse{
				Percentage: sdk.NewDec(1),
			},
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.ClaimedPercentage(
				ctx,
				tt.req,
			)
			if tt.wantErr != "" {
				suite.Error(err)
				if g, w := err.Error(), tt.wantErr; !strings.Contains(g, w) {
					suite.T().Fatalf("Error mismatch:\n\t%q\n\tdoes not contain\n\t%q", g, w)
				}
				return
			}
			suite.NoError(err)
			suite.Equal(tt.resp, resp)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ClaimedPercentageByClaimType() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	tests := []struct {
		name     string
		malleate func()
		req      *types.QueryClaimedPercentageRequest
		wantErr  string
		resp     *types.QueryClaimedPercentageResponse
	}{
		{
			"ClaimedPercentage_Nil_Request",
			func() {},
			nil,
			"empty request",
			nil,
		},
		{
			"ClaimedPercentage_Invalid_Zone",
			func() {
				// setup zones
				suite.setupTestZones()
			},
			&types.QueryClaimedPercentageRequest{
				ChainId: "boguschain",
			},
			"no zone found",
			nil,
		},
		{
			"ClaimedPercentage_Invalid_ClaimType",
			func() {
			},
			&types.QueryClaimedPercentageRequest{
				ChainId:   suite.chainB.ChainID,
				ClaimType: 10000,
			},
			"claim type must be a valid number",
			nil,
		},
		{
			"ClaimedPercentage_No_Zone_Records",
			func() {},
			&types.QueryClaimedPercentageRequest{
				ChainId:   suite.chainB.ChainID,
				ClaimType: claimsmanagertypes.ClaimTypeOsmosisPool,
			},
			"",
			&types.QueryClaimedPercentageResponse{
				Percentage: sdk.ZeroDec(),
			},
		},
		{
			"ClaimedPercentage_Valid_Claims_ClaimTypeOsmosisPool",
			func() {
				addr1, addr2, addr3 := addressutils.GenerateAccAddressForTest().String(), addressutils.GenerateAccAddressForTest().String(), addressutils.GenerateAccAddressForTest().String()
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)
				claims := []claimsmanagertypes.Claim{}
				claims = append(claims, claimsmanagertypes.NewClaim(addr1, zone.ChainId, claimsmanagertypes.ClaimTypeLiquidToken, "", math.NewInt(1000)))
				claims = append(claims, claimsmanagertypes.NewClaim(addr2, zone.ChainId, claimsmanagertypes.ClaimTypeOsmosisPool, "", math.NewInt(3000)))
				claims = append(claims, claimsmanagertypes.NewClaim(addr3, zone.ChainId, claimsmanagertypes.ClaimTypeOsmosisPool, "", math.NewInt(6000)))
				for _, claim := range claims {
					icsKeeper.ClaimsManagerKeeper.SetClaim(ctx, &claim) // #nosec G601
					err := suite.GetQuicksilverApp(suite.chainA).MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, claim.Amount)))
					suite.NoError(err)
				}
			},
			&types.QueryClaimedPercentageRequest{
				ChainId:   suite.chainB.ChainID,
				ClaimType: claimsmanagertypes.ClaimTypeOsmosisPool,
			},
			"",
			&types.QueryClaimedPercentageResponse{
				Percentage: sdk.MustNewDecFromStr("0.9"),
			},
		},
		{
			"ClaimedPercentage_Valid_Claims_ClaimTypeLiquid",
			func() {
				addr1, addr2, addr3 := addressutils.GenerateAccAddressForTest().String(), addressutils.GenerateAccAddressForTest().String(), addressutils.GenerateAccAddressForTest().String()
				zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.True(found)
				claims := []claimsmanagertypes.Claim{}
				claims = append(claims, claimsmanagertypes.NewClaim(addr1, zone.ChainId, claimsmanagertypes.ClaimTypeLiquidToken, "", math.NewInt(1000)))
				claims = append(claims, claimsmanagertypes.NewClaim(addr2, zone.ChainId, claimsmanagertypes.ClaimTypeOsmosisPool, "", math.NewInt(3000)))
				claims = append(claims, claimsmanagertypes.NewClaim(addr3, zone.ChainId, claimsmanagertypes.ClaimTypeOsmosisPool, "", math.NewInt(6000)))
				for _, claim := range claims {
					icsKeeper.ClaimsManagerKeeper.SetClaim(ctx, &claim) // #nosec G601
					err := suite.GetQuicksilverApp(suite.chainA).MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, claim.Amount)))
					suite.NoError(err)
				}
			},
			&types.QueryClaimedPercentageRequest{
				ChainId:   suite.chainB.ChainID,
				ClaimType: claimsmanagertypes.ClaimTypeLiquidToken,
			},
			"",
			&types.QueryClaimedPercentageResponse{
				Percentage: sdk.MustNewDecFromStr("0.1"),
			},
		},
	}

	// run tests:
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()
			resp, err := icsKeeper.ClaimedPercentageByClaimType(
				ctx,
				tt.req,
			)
			if tt.wantErr != "" {
				suite.Error(err)
				if g, w := err.Error(), tt.wantErr; !strings.Contains(g, w) {
					suite.T().Fatalf("Error mismatch:\n\t%q\n\tdoes not contain\n\t%q", g, w)
				}
				return
			}
			suite.NoError(err)
			suite.NotNil(resp)
			suite.Equal(tt.resp, resp)
		})
	}
}
