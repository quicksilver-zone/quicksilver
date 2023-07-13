package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_DelegationStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	zoneValidatorAddresses := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

	performanceDelegations := icsKeeper.GetAllPerformanceDelegations(ctx, &zone)
	suite.Require().Len(performanceDelegations, 4)

	performanceDelegationPointers := icsKeeper.GetAllPerformanceDelegationsAsPointer(ctx, &zone)
	for i, pdp := range performanceDelegationPointers {
		suite.Require().Equal(performanceDelegations[i], *pdp)
	}

	// update performance delegation
	updateDelegation, found := icsKeeper.GetPerformanceDelegation(ctx, &zone, zoneValidatorAddresses[0])
	suite.Require().True(found)
	suite.Require().Equal(uint64(0), updateDelegation.Amount.Amount.Uint64())

	updateDelegation.Amount.Amount = sdkmath.NewInt(10000)
	icsKeeper.SetPerformanceDelegation(ctx, &zone, updateDelegation)

	updatedDelegation, found := icsKeeper.GetPerformanceDelegation(ctx, &zone, zoneValidatorAddresses[0])
	suite.Require().True(found)
	suite.Require().Equal(updateDelegation, updatedDelegation)

	// check that there are no delegations
	delegations := icsKeeper.GetAllDelegations(ctx, &zone)
	suite.Require().Len(delegations, 0)

	// set delegations
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		types.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[0],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000)),
		),
	)
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		types.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[1],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000000)),
		),
	)
	icsKeeper.SetDelegation(
		ctx,
		&zone,
		types.NewDelegation(
			zone.DelegationAddress.Address,
			zoneValidatorAddresses[2],
			sdk.NewCoin(zone.BaseDenom, sdk.NewInt(20000000)),
		),
	)

	// check for delegations set above
	delegations = icsKeeper.GetAllDelegations(ctx, &zone)
	suite.Require().Len(delegations, 3)

	// load and match pointers
	delegationPointers := icsKeeper.GetAllDelegationsAsPointer(ctx, &zone)
	for i, dp := range delegationPointers {
		suite.Require().Equal(delegations[i], *dp)
	}

	// get delegations for delegation address and match
	addr, err := sdk.AccAddressFromBech32(zone.DelegationAddress.GetAddress())
	suite.Require().NoError(err)
	dds := icsKeeper.GetDelegatorDelegations(ctx, &zone, addr)
	suite.Require().Len(dds, 3)
	suite.Require().Equal(delegations, dds)
}

type delegationUpdate struct {
	delegation types.Delegation
	absolute   bool
}

func (suite *KeeperTestSuite) TestUpdateDelegation() {
	del1 := addressutils.GenerateAccAddressForTest()

	val1 := addressutils.GenerateValAddressForTest()
	val2 := addressutils.GenerateValAddressForTest()
	val3 := addressutils.GenerateValAddressForTest()
	val4 := addressutils.GenerateValAddressForTest()
	val5 := addressutils.GenerateValAddressForTest()
	val6 := addressutils.GenerateValAddressForTest()

	tests := []struct {
		name       string
		delegation *types.Delegation
		updates    []delegationUpdate
		expected   types.Delegation
	}{
		{
			"single update, relative increase +3000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val1.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(6000))},
		},
		{
			"single update, relative increase +3000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   true,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val2.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
		},
		{
			"multi update, relative increase +3000, +2000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
					absolute:   false,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val3.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(8000))},
		},
		{
			"multi update, relative +3000, absolute +2000",
			&types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(3000))},
					absolute:   false,
				},
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
					absolute:   true,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val4.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(2000))},
		},
		{
			"new delegation, relative increase +10000",
			nil,
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val5.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(10000))},
					absolute:   false,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val5.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(10000))},
		},
		{
			"new delegation, absolute increase +15000",
			nil,
			[]delegationUpdate{
				{
					delegation: types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val6.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(15000))},
					absolute:   true,
				},
			},
			types.Delegation{DelegationAddress: del1.String(), ValidatorAddress: val6.String(), Amount: sdk.NewCoin("denom", sdk.NewInt(15000))},
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			qApp := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			zone, found := qApp.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			if tt.delegation != nil {
				qApp.InterchainstakingKeeper.SetDelegation(ctx, &zone, *tt.delegation)
			}

			for _, update := range tt.updates {
				err := qApp.InterchainstakingKeeper.UpdateDelegationRecordForAddress(ctx, update.delegation.DelegationAddress, update.delegation.ValidatorAddress, update.delegation.Amount, &zone, update.absolute)
				suite.Require().NoError(err)
			}

			actual, found := qApp.InterchainstakingKeeper.GetDelegation(ctx, &zone, tt.expected.DelegationAddress, tt.expected.ValidatorAddress)
			suite.Require().True(found)
			suite.Require().Equal(tt.expected, actual)
		})
	}
}

func (suite *KeeperTestSuite) TestStoreGetDeleteDelegation() {
	suite.Run("delegation - store / get / delete", func() {
		suite.SetupTest()
		suite.setupTestZones()

		qApp := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := qApp.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.Require().True(found)

		delegator := addressutils.GenerateAccAddressForTest()
		validator := addressutils.GenerateValAddressForTest()

		_, found = qApp.InterchainstakingKeeper.GetDelegation(ctx, &zone, delegator.String(), validator.String())
		suite.Require().False(found)

		newDelegation := types.NewDelegation(delegator.String(), validator.String(), sdk.NewCoin("uatom", sdk.NewInt(5000)))
		qApp.InterchainstakingKeeper.SetDelegation(ctx, &zone, newDelegation)

		fetchedDelegation, found := qApp.InterchainstakingKeeper.GetDelegation(ctx, &zone, delegator.String(), validator.String())
		suite.Require().True(found)
		suite.Require().Equal(newDelegation, fetchedDelegation)

		allDelegations := qApp.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)
		suite.Require().Len(allDelegations, 1)

		err := qApp.InterchainstakingKeeper.RemoveDelegation(ctx, &zone, newDelegation)
		suite.Require().NoError(err)

		allDelegations2 := qApp.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)
		suite.Require().Len(allDelegations2, 0)
	})
}

func (suite *KeeperTestSuite) TestFlushOutstandingDelegations() {
	userAddress := addressutils.GenerateAccAddressForTest().String()
	denom := "uatom"
	tests := []struct {
		name               string
		setStatements      func(ctx sdk.Context, quicksilver *app.Quicksilver)
		delAddrBalance     sdk.Coin
		mockAck            bool
		expectedDelegation sdk.Coins
		assertStatements   func(ctx sdk.Context, quicksilver *app.Quicksilver) bool
	}{
		{
			name:           "case 0: zero delegation balance, no pending receipts, no excluded receipts",
			setStatements:  func(ctx sdk.Context, quicksilver *app.Quicksilver) {},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.ZeroInt()),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				return true
			},
		},
		{
			name: "case 1: zero delegation balance, 2 pending receipts and no excluded receipts",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				receiptOneTime := cutOffTime.Add(-2 * time.Hour)
				receiptTwoTime := cutOffTime.Add(-3 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000),
						),
					),
					FirstSeen: &receiptOneTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100),
						),
					),
					FirstSeen: &receiptTwoTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(0)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}

					return false
				})

				suite.Require().Equal(0, count)
				return true
			},
		},
		{
			name: "case 2: zero delegation balance, 1 pending receipt and 1 excluded receipt",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				receiptOneTime := cutOffTime.Add(-2 * time.Hour)
				receiptTwoTime := cutOffTime.Add(2 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000),
						),
					),
					FirstSeen: &receiptOneTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100),
						),
					),
					FirstSeen: &receiptTwoTime,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(100)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				suite.Require().Equal(1, count)
				return true
			},
		},
		{
			name: "case 3: non-zero delegation balance, 1 pending receipts and 1 excluded receipts ",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)  // -24h
				receiptOneTime := cutOffTime.Add(-2 * time.Hour) // -26h
				receiptTwoTime := cutOffTime.Add(2 * time.Hour)  // -22h

				rcpt1 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000),
						),
					),
					FirstSeen: &receiptOneTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100),
						),
					),
					FirstSeen: &receiptTwoTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(2000100)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				suite.Require().Equal(1, count)
				return true
			},
			mockAck:            true,
			expectedDelegation: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
		},
		{
			name: "case 4: non-zero delegation balance, 2 pending receipts",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				receiptOneTime := cutOffTime.Add(-2 * time.Hour)
				receiptTwoTime := cutOffTime.Add(-3 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000),
						),
					),
					FirstSeen: &receiptOneTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100),
						),
					),
					FirstSeen: &receiptTwoTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(2000100)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				suite.Require().Equal(0, count)
				return true
			},
			mockAck:            true,
			expectedDelegation: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
		},
		{
			name: "case 5: zero delegation balance, 1 pending receipt, 1 excluded receipt",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				receiptOneTime := cutOffTime.Add(-2 * time.Hour)
				receiptTwoTime := cutOffTime.Add(2 * time.Hour)

				rcpt1 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000),
						),
					),
					FirstSeen: &receiptOneTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100),
						),
					),
					FirstSeen: &receiptTwoTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(0)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				suite.Require().Equal(1, count)
				return true
			},
			// zero delegation balance must mean that we cannot delegate anything.
			mockAck: false,
		},
		{
			name: "case 6: low delegation account balance",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				cutOffTime := ctx.BlockTime().AddDate(0, 0, -1)
				receiptOneTime := cutOffTime.Add(-2 * time.Hour)
				rcpt1 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit01",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(2000000),
						),
					),
					FirstSeen: &receiptOneTime,
					Completed: nil,
				}

				rcpt2 := types.Receipt{
					ChainId: suite.chainB.ChainID,
					Sender:  userAddress,
					Txhash:  "TestDeposit02",
					Amount: sdk.NewCoins(
						sdk.NewCoin(
							denom,
							sdk.NewIntFromUint64(100),
						),
					),
					FirstSeen: &receiptOneTime,
					Completed: nil,
				}
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)
			},
			delAddrBalance: sdk.NewCoin("uatom", sdkmath.NewInt(100)),
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				count := 0
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				suite.Require().True(found)
				quicksilver.InterchainstakingKeeper.IterateZoneReceipts(ctx, &zone, func(index int64, receiptInfo types.Receipt) (stop bool) {
					if receiptInfo.Completed == nil {
						count++
					}
					return false
				})
				suite.Require().Equal(0, count)
				return true
			},
			// delegation balance == 100, which equals the value of the second receipt.
			mockAck:            true,
			expectedDelegation: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100))),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			test.setStatements(ctx, quicksilver)
			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)
			err := quicksilver.InterchainstakingKeeper.FlushOutstandingDelegations(ctx, &zone, test.delAddrBalance)
			// refetch zone after FlushOutstandingDelegations setZone().
			ctx = suite.chainA.GetContext()
			if test.mockAck {
				var msgs []sdk.Msg
				allocations, err := quicksilver.InterchainstakingKeeper.DeterminePlanForDelegation(ctx, &zone, test.expectedDelegation)
				suite.Require().NoError(err)
				msgs = append(msgs, quicksilver.InterchainstakingKeeper.PrepareDelegationMessagesForCoins(&zone, allocations)...)
				for _, msg := range msgs {
					err := quicksilver.InterchainstakingKeeper.HandleDelegate(ctx, msg, "batch/1577836910")
					suite.Require().NoError(err)
				}
			}
			// fmt.Println(before, after)
			suite.Require().NoError(err)
			isCorrect := test.assertStatements(ctx, quicksilver)
			suite.Require().True(isCorrect)
		})
	}
}
