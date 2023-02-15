package keeper_test

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_ZoneInfos() {
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	tests := []struct {
		name         string
		malleate     func()
		req          *types.QueryZonesInfoRequest
		wantErr      bool
		expectLength int
	}{
		{
			"ZoneInfos_No_State",
			func() {},
			&types.QueryZonesInfoRequest{},
			false,
			0,
		},
		{
			"ZoneInfos_Nil_Request",
			func() {},
			nil,
			true,
			0,
		},
		{
			"ZoneInfos_Valid_Request",
			func() {
				// setup zones
				suite.setupTestZones()
			},
			&types.QueryZonesInfoRequest{},
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
								ValoperAddress: zone.GetValidatorsAddressesAsSlice()[0],
								Weight:         sdk.OneDec(),
							},
						},
					},
				}
				for _, intent := range intents {
					icsKeeper.SetIntent(ctx, zone, intent, false)
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
				// TODO: set standardized delegatyions for keeper_test package
				delegation := types.Delegation{
					DelegationAddress: testAddress,
					ValidatorAddress:  zone.GetValidatorsAddressesAsSlice()[0],
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
					zone,
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
				DelegatorAddress: "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
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
						Valoper: zone.GetValidatorsAddressesAsSlice()[0],
						Amount:  10000000,
					},
					{
						Valoper: zone.GetValidatorsAddressesAsSlice()[1],
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
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
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
						Valoper: zone.GetValidatorsAddressesAsSlice()[0],
						Amount:  10000000,
					},
					{
						Valoper: zone.GetValidatorsAddressesAsSlice()[1],
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
						Validator:     zone.GetValidatorsAddressesAsSlice()[0],
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
						Source:      zone.GetValidatorsAddressesAsSlice()[1],
						Destination: zone.GetValidatorsAddressesAsSlice()[0],
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
