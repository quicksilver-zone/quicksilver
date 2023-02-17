package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgTransferGood(t *testing.T) {
	app, ctx := app.GetAppWithContext(t, true)
	app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))

	sender := utils.GenerateAccAddressForTest()
	senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)

	txMacc := app.AccountKeeper.GetModuleAddress(icstypes.ModuleName)
	feeMacc := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	txMaccBalance := app.BankKeeper.GetAllBalances(ctx, txMacc)
	feeMaccBalance := app.BankKeeper.GetAllBalances(ctx, feeMacc)

	transferMsg := ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
		Sender:        senderAddr,
		Receiver:      app.AccountKeeper.GetModuleAddress(icstypes.ModuleName).String(),
	}
	require.NoError(t, app.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))

	txMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, txMacc)
	feeMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, feeMacc)

	// assert that ics module balance is now 100denom less than before HandleMsgTransfer()
	require.Equal(t, txMaccBalance.AmountOf("denom").Sub(txMaccBalance2.AmountOf("denom")), sdk.NewInt(100))
	// assert that fee collector module balance is now 100denom more than before HandleMsgTransfer()
	require.Equal(t, feeMaccBalance2.AmountOf("denom").Sub(feeMaccBalance.AmountOf("denom")), sdk.NewInt(100))
}

func TestHandleMsgTransferBadType(t *testing.T) {
	app, ctx := app.GetAppWithContext(t, true)
	app.BankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))

	transferMsg := banktypes.MsgSend{}
	require.Error(t, app.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))
}

func TestHandleMsgTransferBadRecipient(t *testing.T) {
	recipient := utils.GenerateAccAddressForTest()
	app, ctx := app.GetAppWithContext(t, true)

	sender := utils.GenerateAccAddressForTest()
	senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)

	transferMsg := ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
		Sender:        senderAddr,
		Receiver:      recipient.String(),
	}
	require.Error(t, app.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))
}

// TODO: add test cases for send.
// func (s *KeeperTestSuite) TestHandleSendToDelegate() {
// 	tests := []struct {
// 		name string
// 	}{
// 		{
// 			name: "valid",
// 		},
// 	}

// 	for _, test := range tests {
// 		s.Run(test.name, func() {

// 			s.SetupTest()
// 			s.setupTestZones()

// 			recipient := utils.GenerateAccAddressForTest()
// 			app := s.GetQuicksilverApp(s.chainA)
// 			ctx := s.chainA.GetContext()
// 			ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("connectionID"), s.path.EndpointA.ConnectionID))

// 			sender := utils.GenerateAccAddressForTest()
// 			senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)

// 			sendMsg := banktypes.MsgSend{
// 				FromAddress: senderAddr,
// 				ToAddress:   recipient.String(),
// 				Amount:      sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))),
// 			}
// 			s.Require().NoError(app.InterchainstakingKeeper.HandleCompleteSend(ctx, &sendMsg, ""))
// 		})
// 	}
// }

func mustGetTestBech32Address(hrp string) string {
	outAddr, err := bech32.ConvertAndEncode(hrp, utils.GenerateAccAddressForTest())
	if err != nil {
		panic(err)
	}
	return outAddr
}

func (s *KeeperTestSuite) TestHandleQueuedUnbondings() {
	tests := []struct {
		name             string
		records          func(zone *icstypes.Zone) []icstypes.WithdrawalRecord
		delegations      func(zone *icstypes.Zone) []icstypes.Delegation
		redelegations    func(zone *icstypes.Zone) []icstypes.RedelegationRecord
		expectTransition []bool
		expectError      bool
	}{
		{
			name: "valid",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: zone.Validators[0].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[1].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[2].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[3].ValoperAddress, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(zone *icstypes.Zone) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[0].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[1].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[2].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[3].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
				}
			},
			redelegations: func(zone *icstypes.Zone) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{}
			},
			expectTransition: []bool{true},
			expectError:      false,
		},
		{
			name: "valid - two",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: zone.Validators[0].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[1].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[2].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[3].ValoperAddress, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: zone.Validators[0].ValoperAddress, Amount: 5000000},
							{Valoper: zone.Validators[1].ValoperAddress, Amount: 2500000},
							{Valoper: zone.Validators[2].ValoperAddress, Amount: 5000000},
							{Valoper: zone.Validators[3].ValoperAddress, Amount: 2500000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(zone *icstypes.Zone) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[0].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[1].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[2].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[3].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
				}
			},
			redelegations: func(zone *icstypes.Zone) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{}
			},
			expectTransition: []bool{true, true},
			expectError:      false,
		},
		{
			name: "invalid - locked tokens",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: zone.Validators[0].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[1].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[2].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[3].ValoperAddress, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(zone *icstypes.Zone) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[0].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[1].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[2].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[3].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
				}
			},
			redelegations: func(zone *icstypes.Zone) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{
					{
						ChainId:        zone.ChainId,
						EpochNumber:    1,
						Source:         zone.Validators[3].ValoperAddress,
						Destination:    zone.Validators[0].ValoperAddress,
						Amount:         50000,
						CompletionTime: time.Now().Add(time.Hour),
					},
				}
			},
			expectTransition: []bool{false},
			expectError:      true,
		},
		{
			name: "mixed - locked tokens but both succeed (previously failed)",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: zone.Validators[0].ValoperAddress, Amount: 5000000},
							{Valoper: zone.Validators[1].ValoperAddress, Amount: 2500000},
							{Valoper: zone.Validators[2].ValoperAddress, Amount: 5000000},
							{Valoper: zone.Validators[3].ValoperAddress, Amount: 2500000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: zone.Validators[0].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[1].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[2].ValoperAddress, Amount: 1000000},
							{Valoper: zone.Validators[3].ValoperAddress, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(zone *icstypes.Zone) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[0].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[1].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[2].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  zone.Validators[3].ValoperAddress,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
				}
			},
			redelegations: func(zone *icstypes.Zone) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{
					{
						ChainId:        zone.ChainId,
						EpochNumber:    1,
						Source:         zone.Validators[3].ValoperAddress,
						Destination:    zone.Validators[0].ValoperAddress,
						Amount:         1000001,
						CompletionTime: time.Now().Add(time.Hour),
					},
				}
			},
			expectTransition: []bool{true, true},
			expectError:      false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			if !found {
				s.Fail("unable to retrieve zone for test")
			}
			records := test.records(&zone)
			delegations := test.delegations(&zone)
			redelegations := test.redelegations(&zone)

			// set up zones
			for _, record := range records {
				app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
			}

			for _, delegation := range delegations {
				app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
				val, _ := zone.GetValidatorByValoper(delegation.ValidatorAddress)
				val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
				val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))
			}

			for _, redelegation := range redelegations {
				app.InterchainstakingKeeper.SetRedelegationRecord(ctx, redelegation)
			}

			app.InterchainstakingKeeper.SetZone(ctx, &zone)

			// trigger handler
			err := app.InterchainstakingKeeper.HandleQueuedUnbondings(ctx, &zone, 1)
			if test.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

			for idx, record := range records {
				// check record with old status is opposite to expectedTransition (if false, this record should exist in status 3)
				_, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusQueued)
				s.Require().Equal(!test.expectTransition[idx], found)
				// check record with new status is as per expectedTransition (if false, this record should not exist in status 4)
				_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusUnbond)
				s.Require().Equal(test.expectTransition[idx], found)

				if test.expectTransition[idx] {
					actualRecord, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusUnbond)
					s.Require().True(found)
					for _, unbonding := range actualRecord.Distribution {
						r, found := app.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, unbonding.Valoper, 1)
						fmt.Println(r)
						s.Require().True(found)
						s.Require().Contains(r.RelatedTxhash, record.Txhash)
					}
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestHandleWithdrawForUser() {
	tests := []struct {
		name    string
		records func(zone *icstypes.Zone) []icstypes.WithdrawalRecord
		message banktypes.MsgSend
		memo    string
		err     bool
	}{
		{
			name: "invalid - no matching record",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			message: banktypes.MsgSend{},
			memo:    "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:     true,
		},
		{
			name: "valid",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusSend,
					},
				}
			},
			message: banktypes.MsgSend{
				Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
			},
			memo: "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:  false,
		},
		{
			name: "valid - two",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusSend,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 5000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1250000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 5000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1250000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusSend,
					},
				}
			},
			message: banktypes.MsgSend{
				Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
			},
			memo: "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
			err:  false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			if !found {
				s.Fail("unable to retrieve zone for test")
			}

			records := test.records(&zone)

			// set up zones
			for _, record := range records {
				app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				err := app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(record.BurnAmount))
				s.Require().NoError(err)
				err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount))
				s.Require().NoError(err)
			}

			// trigger handler
			err := app.InterchainstakingKeeper.HandleWithdrawForUser(ctx, &zone, &test.message, test.memo)
			if test.err {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

			app.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusSend, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash == test.memo {
					s.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			app.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusCompleted, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash != test.memo {
					s.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})
		})
	}
}

func (s *KeeperTestSuite) TestHandleWithdrawForUserLSM() {
	v1 := utils.GenerateValAddressForTest().String()
	v2 := utils.GenerateValAddressForTest().String()
	tests := []struct {
		name    string
		records func(zone *icstypes.Zone) []icstypes.WithdrawalRecord
		message []banktypes.MsgSend
		memo    string
		err     bool
	}{
		{
			name: "valid",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: v1, Amount: 1000000},
							{Valoper: v2, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusSend,
					},
				}
			},
			message: []banktypes.MsgSend{
				{Amount: sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000)))},
				{Amount: sdk.NewCoins(sdk.NewCoin(v2+"2", sdk.NewInt(1000000)))},
			},
			memo: "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:  false,
		},
		{
			name: "valid - unequal",
			records: func(zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: v1, Amount: 1000000},
							{Valoper: v2, Amount: 1500000},
						},
						Recipient:  mustGetTestBech32Address(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2500000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2500000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusSend,
					},
				}
			},
			message: []banktypes.MsgSend{
				{Amount: sdk.NewCoins(sdk.NewCoin(v2+"1", sdk.NewInt(1500000)))},
				{Amount: sdk.NewCoins(sdk.NewCoin(v1+"2", sdk.NewInt(1000000)))},
			},
			memo: "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:  false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			if !found {
				s.Fail("unable to retrieve zone for test")
			}

			records := test.records(&zone)

			startBalance := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(icstypes.ModuleName))
			// set up zones
			for _, record := range records {
				app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				err := app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(record.BurnAmount))
				s.Require().NoError(err)
				err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount))
				s.Require().NoError(err)
			}

			// trigger handler
			for _, msg := range test.message {
				err := app.InterchainstakingKeeper.HandleWithdrawForUser(ctx, &zone, &msg, test.memo)
				if test.err {
					s.Require().Error(err)
				} else {
					s.Require().NoError(err)
				}
			}

			app.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusSend, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash == test.memo {
					s.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			app.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusCompleted, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash != test.memo {
					s.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			postBurnBalance := app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(icstypes.ModuleName))
			s.Require().Equal(startBalance, postBurnBalance)
		})
	}
}

func (s *KeeperTestSuite) TestReceiveAckErrForBeginRedelegate() {
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	if !found {
		s.Fail("unable to retrieve zone for test")
	}

	// create redelegation record
	record := icstypes.RedelegationRecord{
		ChainId:     s.chainB.ChainID,
		EpochNumber: 1,
		Source:      zone.Validators[0].ValoperAddress,
		Destination: zone.Validators[1].ValoperAddress,
		Amount:      1000,
	}

	app.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	redelegate := &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: zone.Validators[0].ValoperAddress, ValidatorDstAddress: zone.Validators[1].ValoperAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	data, err := icatypes.SerializeCosmosTx(app.InterchainstakingKeeper.GetCodec(), []sdk.Msg{redelegate})
	s.Require().NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: fmt.Sprintf("rebalance/%d", 1),
	}

	packet := channeltypes.Packet{Data: app.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	ackBytes := []byte("{\"error\":\"ABCI code: 32: error handling packet on host chain: see events for details\"}")
	// call handler

	_, found = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, zone.Validators[0].ValoperAddress, zone.Validators[1].ValoperAddress, 1)
	s.Require().True(found)

	err = app.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes)
	s.Require().NoError(err)

	_, found = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, zone.Validators[0].ValoperAddress, zone.Validators[1].ValoperAddress, 1)
	s.Require().False(found)
}

func (s *KeeperTestSuite) TestReceiveAckErrForBeginUndelegate() {
	hash1 := fmt.Sprintf("%x", sha256.Sum256([]byte{0x01}))
	hash2 := fmt.Sprintf("%x", sha256.Sum256([]byte{0x02}))
	hash3 := fmt.Sprintf("%x", sha256.Sum256([]byte{0x03}))
	delegator1 := utils.GenerateAccAddressForTest().String()
	delegator2 := utils.GenerateAccAddressForTest().String()

	randRr := rand.Float32() + 1.0

	tests := []struct {
		name                      string
		epoch                     int64
		withdrawalRecords         func(zone icstypes.Zone) []icstypes.WithdrawalRecord
		unbondingRecords          func(zone icstypes.Zone) []icstypes.UnbondingRecord
		msgs                      func(zone icstypes.Zone) []sdk.Msg
		expectedWithdrawalRecords func(zone icstypes.Zone) []icstypes.WithdrawalRecord
	}{
		{
			name:  "1 wdr, 2 vals, 1k+1k, 1800 qasset",
			epoch: 1,
			withdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  1000,
							},
							{
								Valoper: zone.Validators[1].ValoperAddress,
								Amount:  1000,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(zone icstypes.Zone) []icstypes.UnbondingRecord {
				return []icstypes.UnbondingRecord{
					{
						ChainId:       s.chainB.ChainID,
						EpochNumber:   1,
						Validator:     zone.Validators[0].ValoperAddress,
						RelatedTxhash: []string{hash1},
					},
				}
			},
			msgs: func(zone icstypes.Zone) []sdk.Msg {
				return []sdk.Msg{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: zone.Validators[0].ValoperAddress,
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			expectedWithdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[1].ValoperAddress,
								Amount:  1000,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:      s.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:       fmt.Sprintf("%064d", 1),
						Status:       icskeeper.WithdrawStatusQueued,
					},
				}
			},
		},
		{
			name:  "1 wdr, 1 vals, 1k, 900 qasset",
			epoch: 1,
			withdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  1000,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(zone icstypes.Zone) []icstypes.UnbondingRecord {
				return []icstypes.UnbondingRecord{
					{
						ChainId:       s.chainB.ChainID,
						EpochNumber:   1,
						Validator:     zone.Validators[0].ValoperAddress,
						RelatedTxhash: []string{hash1},
					},
				}
			},
			msgs: func(zone icstypes.Zone) []sdk.Msg {
				return []sdk.Msg{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: zone.Validators[0].ValoperAddress,
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			expectedWithdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:      s.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:       hash1,
						Status:       icskeeper.WithdrawStatusQueued,
					},
				}
			},
		},
		{
			name:  "3 wdr, 2 vals, 1k+0.5k, 1350 qasset; 1k+2k, 2700 qasset; 600+400, 900qasset",
			epoch: 2,
			withdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  1000,
							},
							{
								Valoper: zone.Validators[1].ValoperAddress,
								Amount:  500,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1500))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1350)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  1000,
							},
							{
								Valoper: zone.Validators[1].ValoperAddress,
								Amount:  2000,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(2700)),
						Txhash:     hash2,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  600,
							},
							{
								Valoper: zone.Validators[1].ValoperAddress,
								Amount:  400,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash3,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(zone icstypes.Zone) []icstypes.UnbondingRecord {
				return []icstypes.UnbondingRecord{
					{
						ChainId:       s.chainB.ChainID,
						EpochNumber:   2,
						Validator:     zone.Validators[1].ValoperAddress,
						RelatedTxhash: []string{hash1, hash2, hash3},
					},
				}
			},
			msgs: func(zone icstypes.Zone) []sdk.Msg {
				return []sdk.Msg{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: zone.Validators[1].ValoperAddress,
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2900)),
					},
				}
			},
			expectedWithdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  1000,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  1000,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash2,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  600,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(600))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(540)),
						Txhash:     hash3,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:      s.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(500))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(450)),
						Txhash:       fmt.Sprintf("%064d", 1),
						Status:       icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:      s.chainB.ChainID,
						Delegator:    delegator2,
						Distribution: nil,
						Recipient:    mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:       fmt.Sprintf("%064d", 2),
						Status:       icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:      s.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(400))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(360)),
						Txhash:       fmt.Sprintf("%064d", 3),
						Status:       icskeeper.WithdrawStatusQueued,
					},
				}
			},
		},
		{
			name:  "2 wdr, random_rr, 1 vals, 1k; 2 vals; 123 + 456 ",
			epoch: 1,
			withdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[0].ValoperAddress,
								Amount:  1000,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewDec(1000).QuoTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[1].ValoperAddress,
								Amount:  123,
							},
							{
								Valoper: zone.Validators[2].ValoperAddress,
								Amount:  456,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(579))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewDec(579).QuoTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
						Txhash:     hash2,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(zone icstypes.Zone) []icstypes.UnbondingRecord {
				return []icstypes.UnbondingRecord{
					{
						ChainId:       s.chainB.ChainID,
						EpochNumber:   1,
						Validator:     zone.Validators[0].ValoperAddress,
						RelatedTxhash: []string{hash1},
					},
					{
						ChainId:       s.chainB.ChainID,
						EpochNumber:   1,
						Validator:     zone.Validators[1].ValoperAddress,
						RelatedTxhash: []string{hash2},
					},
					// {
					// 	ChainId:       s.chainB.ChainID,
					// 	EpochNumber:   1,
					// 	Validator:     zone.Validators[2].ValoperAddress,
					// 	RelatedTxhash: []string{hash2},
					// },
				}
			},
			msgs: func(zone icstypes.Zone) []sdk.Msg {
				return []sdk.Msg{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: zone.Validators[0].ValoperAddress,
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: zone.Validators[1].ValoperAddress,
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(123)),
					},
				}
			},
			expectedWithdrawalRecords: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:      s.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewDec(1000).QuoTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
						Txhash:       hash1,
						Status:       icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:      s.chainB.ChainID,
						Delegator:    delegator2,
						Distribution: nil,
						Recipient:    mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(123))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewDec(123).QuoTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
						Txhash:       fmt.Sprintf("%064d", 1),
						Status:       icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   s.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: zone.Validators[2].ValoperAddress,
								Amount:  456,
							},
						},
						Recipient:  mustGetTestBech32Address(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(456))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewDec(456).QuoTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
						Txhash:     hash2,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			if !found {
				s.Fail("unable to retrieve zone for test")
			}

			for _, wdr := range test.withdrawalRecords(zone) {
				app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, wdr)
			}

			for _, ubr := range test.unbondingRecords(zone) {
				app.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr)
			}

			data, err := icatypes.SerializeCosmosTx(app.InterchainstakingKeeper.GetCodec(), test.msgs(zone))
			s.Require().NoError(err)

			// validate memo < 256 bytes
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: data,
				Memo: fmt.Sprintf("withdrawal/%d", test.epoch),
			}

			packet := channeltypes.Packet{Data: app.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

			ackBytes := []byte("{\"error\":\"ABCI code: 32: error handling packet on host chain: see events for details\"}")
			// call handler

			for _, ubr := range test.unbondingRecords(zone) {
				_, found = app.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, ubr.Validator, test.epoch)
				s.Require().True(found)
			}

			err = app.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes)
			s.Require().NoError(err)

			for _, ubr := range test.unbondingRecords(zone) {
				_, found = app.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, ubr.Validator, test.epoch)
				s.Require().False(found)
			}

			for idx, ewdr := range test.expectedWithdrawalRecords(zone) {
				wdr, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, ewdr.Txhash, ewdr.Status)
				s.Require().True(found)
				s.Require().Equal(ewdr.Amount, wdr.Amount)
				s.Require().Equal(ewdr.BurnAmount, wdr.BurnAmount)
				s.Require().Equal(ewdr.Delegator, wdr.Delegator)
				s.Require().Equal(ewdr.Distribution, wdr.Distribution, idx)
				s.Require().Equal(ewdr.Status, wdr.Status)
			}
		})
	}
}
func (s *KeeperTestSuite) TestRebalanceDueToIntentChange() {

	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	if !found {
		s.Fail("unable to retrieve zone for test")
	}
	vals := zone.Validators

	delegations := []icstypes.Delegation{
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[0].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[1].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[2].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		}, {
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[3].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
	}
	for _, delegation := range delegations {
		app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		val, _ := zone.GetValidatorByValoper(delegation.ValidatorAddress)
		val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
		val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))
	}

	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	// trigger rebalance
	err := app.InterchainstakingKeeper.Rebalance(ctx, zone, 1)
	s.Require().NoError(err)

	//change intents to trigger redelegations from val[3]
	intents := icstypes.ValidatorIntents{
		{ValoperAddress: vals[0].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[1].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[2].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[3].ValoperAddress, Weight: sdk.NewDecWithPrec(1, 1)},
	}
	zone.AggregateIntent = intents

	// trigger rebalance
	err = app.InterchainstakingKeeper.Rebalance(ctx, zone, 2)

	// mock ack for redelegations
	app.InterchainstakingKeeper.IteratePrefixedRedelegationRecords(ctx, []byte(zone.ChainId), func(idx int64, _ []byte, record icstypes.RedelegationRecord) (stop bool) {
		if record.EpochNumber == 2 {
			msg := stakingtypes.MsgBeginRedelegate{zone.DelegationAddress.Address,
				record.Source,
				record.Destination,
				sdk.NewCoin("uatom", sdkmath.NewInt(record.Amount))}
			err := app.InterchainstakingKeeper.HandleBeginRedelegate(ctx, &msg, time.Now().Add(time.Hour*24*7), fmt.Sprintf("rebalance/%d", 2))
			if err != nil {
				return false
			}
		}
		return false
	})

	// check for redelegations
	_, present := app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[0].ValoperAddress, 2)
	s.Require().True(present)
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[1].ValoperAddress, 2)
	s.Require().True(present)
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[2].ValoperAddress, 2)
	s.Require().True(present)

	// change intents to trigger transitive redelegations which should fail rebalance
	zone, _ = app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	intents = icstypes.ValidatorIntents{
		{ValoperAddress: vals[0].ValoperAddress, Weight: sdk.NewDecWithPrec(1, 1)},
		{ValoperAddress: vals[1].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[2].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[3].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
	}
	zone.AggregateIntent = intents

	// trigger rebalance
	err = app.InterchainstakingKeeper.Rebalance(ctx, zone, 3)

	//check for redelegations originating from val[0], they should not be present
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[1].ValoperAddress, 3)
	s.Require().False(present)
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[2].ValoperAddress, 3)
	s.Require().False(present)
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[3].ValoperAddress, 3)
	s.Require().False(present)
}

func (s *KeeperTestSuite) TestRebalanceDueToDelegationChange() {

	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	if !found {
		s.Fail("unable to retrieve zone for test")
	}
	vals := zone.Validators

	delegations := []icstypes.Delegation{
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[0].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[1].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[2].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		}, {
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[3].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
	}
	for _, delegation := range delegations {
		app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		val, _ := zone.GetValidatorByValoper(delegation.ValidatorAddress)
		val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
		val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))
	}

	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	// trigger rebalance
	err := app.InterchainstakingKeeper.Rebalance(ctx, zone, 1)
	s.Require().NoError(err)

	app.InterchainstakingKeeper.IterateAllDelegations(ctx, &zone, func(delegation icstypes.Delegation) bool {
		if delegation.ValidatorAddress == vals[0].ValoperAddress {
			delegation.Amount = delegation.Amount.Add(sdk.NewInt64Coin("uatom", 4000))
			app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		}
		return false
	})

	// trigger rebalance
	err = app.InterchainstakingKeeper.Rebalance(ctx, zone, 2)

	// mock ack for redelegations
	app.InterchainstakingKeeper.IteratePrefixedRedelegationRecords(ctx, []byte(zone.ChainId), func(idx int64, _ []byte, record icstypes.RedelegationRecord) (stop bool) {
		if record.EpochNumber == 2 {
			msg := stakingtypes.MsgBeginRedelegate{zone.DelegationAddress.Address,
				record.Source,
				record.Destination,
				sdk.NewCoin("uatom", sdkmath.NewInt(record.Amount))}
			err := app.InterchainstakingKeeper.HandleBeginRedelegate(ctx, &msg, time.Now().Add(time.Hour*24*7), fmt.Sprintf("rebalance/%d", 2))
			if err != nil {
				return false
			}
		}
		return false
	})

	// check for redelegations
	_, present := app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[2].ValoperAddress, 2)
	s.Require().True(present)
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[3].ValoperAddress, 2)
	s.Require().True(present)

	// change validator delegation to trigger transitive redelegations which should fail rebalance
	app.InterchainstakingKeeper.IterateAllDelegations(ctx, &zone, func(delegation icstypes.Delegation) bool {
		if delegation.ValidatorAddress == vals[0].ValoperAddress {
			delegation.Amount = delegation.Amount.Sub(sdk.NewInt64Coin("uatom", 4000))
			app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		}
		if delegation.ValidatorAddress == vals[2].ValoperAddress {
			delegation.Amount = delegation.Amount.Add(sdk.NewInt64Coin("uatom", 4000))
			app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		}

		return false
	})

	// trigger rebalance
	err = app.InterchainstakingKeeper.Rebalance(ctx, zone, 3)

	//check for redelegations originating from val[1], they should not be present
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[2].ValoperAddress, vals[0].ValoperAddress, 3)
	s.Require().False(present)
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[2].ValoperAddress, vals[1].ValoperAddress, 3)
	s.Require().False(present)
	_, present = app.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[2].ValoperAddress, vals[3].ValoperAddress, 3)
	s.Require().False(present)
}
