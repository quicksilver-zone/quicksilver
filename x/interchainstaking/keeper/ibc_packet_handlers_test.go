package keeper_test

import (
	"fmt"
	"testing"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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
	val1 := utils.GenerateValAddressForTest().String()
	val2 := utils.GenerateValAddressForTest().String()
	val3 := utils.GenerateValAddressForTest().String()
	val4 := utils.GenerateValAddressForTest().String()

	tests := []struct {
		name             string
		records          func(chainID string, hrp string) []icstypes.WithdrawalRecord
		delegations      func(chainID string, delegateAddress string, hrp string) []icstypes.Delegation
		redelegations    func(chainID string, delegateAddress string, hrp string) []icstypes.RedelegationRecord
		expectTransition []bool
	}{
		{
			name: "valid",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: val1, Amount: 1000000},
							{Valoper: val2, Amount: 1000000},
							{Valoper: val3, Amount: 1000000},
							{Valoper: val4, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(chainID string, delegateAddress string, hrp string) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val1,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val2,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val3,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val4,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
				}
			},
			redelegations: func(chainID string, delegateAddress string, hrp string) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{}
			},
			expectTransition: []bool{true},
		},
		{
			name: "valid - two",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: val1, Amount: 1000000},
							{Valoper: val2, Amount: 1000000},
							{Valoper: val3, Amount: 1000000},
							{Valoper: val4, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: val1, Amount: 5000000},
							{Valoper: val2, Amount: 2500000},
							{Valoper: val3, Amount: 5000000},
							{Valoper: val4, Amount: 2500000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(chainID string, delegateAddress string, hrp string) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val1,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val2,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val3,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val4,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
				}
			},
			redelegations: func(chainID string, delegateAddress string, hrp string) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{}
			},
			expectTransition: []bool{true, true},
		},
		{
			name: "invalid - locked tokens",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: val1, Amount: 1000000},
							{Valoper: val2, Amount: 1000000},
							{Valoper: val3, Amount: 1000000},
							{Valoper: val4, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(chainID string, delegateAddress string, hrp string) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val1,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val2,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val3,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val4,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
				}
			},
			redelegations: func(chainID string, delegateAddress string, hrp string) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{
					{
						ChainId:        chainID,
						EpochNumber:    1,
						Source:         val4,
						Destination:    val1,
						Amount:         50000,
						CompletionTime: time.Now().Add(time.Hour),
					},
				}
			},
			expectTransition: []bool{false},
		},
		{
			name: "mixed - locked tokens prohibit first unbond, but second permitted",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: val1, Amount: 5000000},
							{Valoper: val2, Amount: 2500000},
							{Valoper: val3, Amount: 5000000},
							{Valoper: val4, Amount: 2500000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: val1, Amount: 1000000},
							{Valoper: val2, Amount: 1000000},
							{Valoper: val3, Amount: 1000000},
							{Valoper: val4, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(chainID string, delegateAddress string, hrp string) []icstypes.Delegation {
				return []icstypes.Delegation{
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val1,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val2,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val3,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: delegateAddress,
						ValidatorAddress:  val4,
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
				}
			},
			redelegations: func(chainID string, delegateAddress string, hrp string) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{
					{
						ChainId:        chainID,
						EpochNumber:    1,
						Source:         val4,
						Destination:    val1,
						Amount:         1000001,
						CompletionTime: time.Now().Add(time.Hour),
					},
				}
			},
			expectTransition: []bool{false, true},
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
			zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: val1, VotingPower: sdk.ZeroInt(), DelegatorShares: sdk.ZeroDec()})
			zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: val2, VotingPower: sdk.ZeroInt(), DelegatorShares: sdk.ZeroDec()})
			zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: val3, VotingPower: sdk.ZeroInt(), DelegatorShares: sdk.ZeroDec()})
			zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: val4, VotingPower: sdk.ZeroInt(), DelegatorShares: sdk.ZeroDec()})

			records := test.records(s.chainB.ChainID, zone.AccountPrefix)
			delegations := test.delegations(s.chainB.ChainID, zone.DelegationAddress.Address, zone.AccountPrefix)
			redelegations := test.redelegations(s.chainB.ChainID, zone.DelegationAddress.Address, zone.AccountPrefix)

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
			s.Require().NoError(err)

			for idx, record := range records {
				// check record with old status is opposite to expectedTransition (if false, this record should exist in status 3)
				_, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusQueued)
				s.Require().Equal(!test.expectTransition[idx], found)
				// check record with new status is as per expectedTransition (if false, this record should not exist in status 4)
				_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusUnbond)
				s.Require().Equal(test.expectTransition[idx], found)

				if test.expectTransition[idx] {
					for _, unbonding := range record.Distribution {
						r, found := app.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, unbonding.Valoper, 1)
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
		records func(chainID string, hrp string) []icstypes.WithdrawalRecord
		message banktypes.MsgSend
		memo    string
		err     bool
	}{
		{
			name: "invalid - no matching record",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
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
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
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
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusSend,
					},
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 5000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1250000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 5000000},
							{Valoper: utils.GenerateValAddressForTest().String(), Amount: 1250000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
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

			records := test.records(s.chainB.ChainID, zone.AccountPrefix)

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
		records func(chainID string, hrp string) []icstypes.WithdrawalRecord
		message []banktypes.MsgSend
		memo    string
		err     bool
	}{
		{
			name: "valid",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: v1, Amount: 1000000},
							{Valoper: v2, Amount: 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
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
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: v1, Amount: 1000000},
							{Valoper: v2, Amount: 1500000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
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

			records := test.records(s.chainB.ChainID, zone.AccountPrefix)

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

func (s *KeeperTestSuite) Test_v045Callback() {
	tests := []struct {
		name             string
		setStatements    func(ctx sdk.Context, app *app.Quicksilver) ([]sdk.Msg, []byte)
		assertStatements func(ctx sdk.Context, app *app.Quicksilver) bool
	}{
		{
			name: "msg response with some data",
			setStatements: func(ctx sdk.Context, app *app.Quicksilver) ([]sdk.Msg, []byte) {
				app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
				sender := utils.GenerateAccAddressForTest()
				senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)
				transferMsg := ibctransfertypes.MsgTransfer{
					SourcePort:    "transfer",
					SourceChannel: "channel-0",
					Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
					Sender:        senderAddr,
					Receiver:      app.AccountKeeper.GetModuleAddress(icstypes.ModuleName).String(),
				}
				response := ibctransfertypes.MsgTransferResponse{
					Sequence: 1,
				}

				respBytes := icatypes.ModuleCdc.MustMarshal(&response)
				return []sdk.Msg{&transferMsg}, respBytes
			},
			assertStatements: func(ctx sdk.Context, app *app.Quicksilver) bool {
				txMacc := app.AccountKeeper.GetModuleAddress(icstypes.ModuleName)
				feeMacc := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
				txMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, txMacc)
				feeMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, feeMacc)

				// assert that ics module balance is now 100denom less than before HandleMsgTransfer()

				if txMaccBalance2.AmountOf("denom").Equal(sdk.ZeroInt()) && feeMaccBalance2.AmountOf("denom").Equal(sdk.NewInt(100)) {
					return true
				}
				return false
			},
		},
		{
			name: "msg response with nil data",
			setStatements: func(ctx sdk.Context, app *app.Quicksilver) ([]sdk.Msg, []byte) {
				zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				if !found {
					s.Fail("unable to retrieve zone for test")
				}

				msgSetWithdrawAddress := distrtypes.MsgSetWithdrawAddress{
					DelegatorAddress: zone.PerformanceAddress.Address,
					WithdrawAddress:  zone.WithdrawalAddress.Address,
				}

				response := distrtypes.MsgSetWithdrawAddressResponse{}

				respBytes := icatypes.ModuleCdc.MustMarshal(&response)
				return []sdk.Msg{&msgSetWithdrawAddress}, respBytes
			},
			assertStatements: func(ctx sdk.Context, app *app.Quicksilver) bool {
				zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				if !found {
					s.Fail("unable to retrieve zone for test")
				}
				// assert that withdraw address is set
				if zone.WithdrawalAddress.Address == zone.PerformanceAddress.WithdrawalAddress {
					return true
				}
				return false
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			msg, msgResponseBytes := test.setStatements(ctx, app)

			txMsgData := &sdk.TxMsgData{
				Data:         []*sdk.MsgData{{MsgType: "/bob", Data: msgResponseBytes}},
				MsgResponses: []*codectypes.Any{},
			}

			ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

			acknowledgement := channeltypes.Acknowledgement{
				Response: &channeltypes.Acknowledgement_Result{
					Result: ackData,
				},
			}

			pdBytes, err := icatypes.SerializeCosmosTx(icatypes.ModuleCdc, msg)
			s.Require().NoError(err)
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: pdBytes,
				Memo: "test_acknowledgement",
			}

			packetBytes, err := icatypes.ModuleCdc.MarshalJSON(&packetData)
			s.Require().NoError(err)
			packet := channeltypes.Packet{
				Data: packetBytes,
			}

			s.Require().NoError(app.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, icatypes.ModuleCdc.MustMarshalJSON(&acknowledgement)))

			s.Require().True(test.assertStatements(ctx, app))
		})
	}
}

func (s *KeeperTestSuite) Test_v046Callback() {
	tests := []struct {
		name             string
		setStatements    func(ctx sdk.Context, app *app.Quicksilver) ([]sdk.Msg, *codectypes.Any)
		assertStatements func(ctx sdk.Context, app *app.Quicksilver) bool
	}{
		{
			name: "msg response with some data",
			setStatements: func(ctx sdk.Context, app *app.Quicksilver) ([]sdk.Msg, *codectypes.Any) {
				app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
				sender := utils.GenerateAccAddressForTest()
				senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)
				transferMsg := ibctransfertypes.MsgTransfer{
					SourcePort:    "transfer",
					SourceChannel: "channel-0",
					Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
					Sender:        senderAddr,
					Receiver:      app.AccountKeeper.GetModuleAddress(icstypes.ModuleName).String(),
				}
				response := ibctransfertypes.MsgTransferResponse{
					Sequence: 1,
				}

				anyresponse, _ := codectypes.NewAnyWithValue(&response)
				return []sdk.Msg{&transferMsg}, anyresponse
			},
			assertStatements: func(ctx sdk.Context, app *app.Quicksilver) bool {
				txMacc := app.AccountKeeper.GetModuleAddress(icstypes.ModuleName)
				feeMacc := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
				txMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, txMacc)
				feeMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, feeMacc)

				// assert that ics module balance is now 100denom less than before HandleMsgTransfer()

				if txMaccBalance2.AmountOf("denom").Equal(sdk.ZeroInt()) && feeMaccBalance2.AmountOf("denom").Equal(sdk.NewInt(100)) {
					return true
				}
				return false
			},
		},
		{
			name: "msg response with nil data",
			setStatements: func(ctx sdk.Context, app *app.Quicksilver) ([]sdk.Msg, *codectypes.Any) {
				zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				if !found {
					s.Fail("unable to retrieve zone for test")
				}

				msgSetWithdrawAddress := distrtypes.MsgSetWithdrawAddress{
					DelegatorAddress: zone.PerformanceAddress.Address,
					WithdrawAddress:  zone.WithdrawalAddress.Address,
				}

				response := distrtypes.MsgSetWithdrawAddressResponse{}

				anyresponse, _ := codectypes.NewAnyWithValue(&response)
				return []sdk.Msg{&msgSetWithdrawAddress}, anyresponse
			},
			assertStatements: func(ctx sdk.Context, app *app.Quicksilver) bool {
				zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
				if !found {
					s.Fail("unable to retrieve zone for test")
				}
				// assert that withdraw address is set
				if zone.WithdrawalAddress.Address == zone.PerformanceAddress.WithdrawalAddress {
					return true
				}
				return false
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			msg, anyResp := test.setStatements(ctx, app)

			txMsgData := &sdk.TxMsgData{
				Data:         []*sdk.MsgData{},
				MsgResponses: []*codectypes.Any{anyResp},
			}

			ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

			acknowledgement := channeltypes.Acknowledgement{
				Response: &channeltypes.Acknowledgement_Result{
					Result: ackData,
				},
			}

			pdBytes, err := icatypes.SerializeCosmosTx(icatypes.ModuleCdc, msg)
			s.Require().NoError(err)
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: pdBytes,
				Memo: "test_acknowledgement",
			}

			packetBytes, err := icatypes.ModuleCdc.MarshalJSON(&packetData)
			s.Require().NoError(err)
			packet := channeltypes.Packet{
				Data: packetBytes,
			}

			s.Require().NoError(app.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, icatypes.ModuleCdc.MustMarshalJSON(&acknowledgement)))

			s.Require().True(test.assertStatements(ctx, app))
		})
	}
}
