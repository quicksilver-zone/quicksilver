package keeper_test

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestHandleMsgTransferGood(t *testing.T) {
	quicksilver, ctx := app.GetAppWithContext(t, true)
	err := quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
	require.NoError(t, err)

	sender := addressutils.GenerateAccAddressForTest()
	senderAddr, err := sdk.Bech32ifyAddressBytes("cosmos", sender)
	require.NoError(t, err)

	txMacc := quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName)
	feeMacc := quicksilver.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	txMaccBalance := quicksilver.BankKeeper.GetAllBalances(ctx, txMacc)
	feeMaccBalance := quicksilver.BankKeeper.GetAllBalances(ctx, feeMacc)

	transferMsg := ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
		Sender:        senderAddr,
		Receiver:      quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName).String(),
	}
	require.NoError(t, quicksilver.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))

	txMaccBalance2 := quicksilver.BankKeeper.GetAllBalances(ctx, txMacc)
	feeMaccBalance2 := quicksilver.BankKeeper.GetAllBalances(ctx, feeMacc)

	// assert that ics module balance is now 100denom less than before HandleMsgTransfer()
	require.Equal(t, txMaccBalance.AmountOf("denom").Sub(txMaccBalance2.AmountOf("denom")), sdk.NewInt(100))
	// assert that fee collector module balance is now 100denom more than before HandleMsgTransfer()
	require.Equal(t, feeMaccBalance2.AmountOf("denom").Sub(feeMaccBalance.AmountOf("denom")), sdk.NewInt(100))
}

func TestHandleMsgTransferBadType(t *testing.T) {
	quicksilver, ctx := app.GetAppWithContext(t, true)
	err := quicksilver.BankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
	require.NoError(t, err)

	transferMsg := banktypes.MsgSend{}
	require.Error(t, quicksilver.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))
}

func TestHandleMsgTransferBadRecipient(t *testing.T) {
	recipient := addressutils.GenerateAccAddressForTest()
	quicksilver, ctx := app.GetAppWithContext(t, true)

	sender := addressutils.GenerateAccAddressForTest()
	senderAddr, err := sdk.Bech32ifyAddressBytes("cosmos", sender)
	require.NoError(t, err)

	transferMsg := ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
		Sender:        senderAddr,
		Receiver:      recipient.String(),
	}
	require.Error(t, quicksilver.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))
}

func (suite *KeeperTestSuite) TestHandleQueuedUnbondings() {
	tests := []struct {
		name             string
		records          func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.WithdrawalRecord
		delegations      func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.Delegation
		redelegations    func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.RedelegationRecord
		expectTransition []bool
		expectError      bool
	}{
		{
			name: "valid",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: vals[0], Amount: 1000000},
							{Valoper: vals[1], Amount: 1000000},
							{Valoper: vals[2], Amount: 1000000},
							{Valoper: vals[3], Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[2],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[3],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
				}
			},
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{}
			},
			expectTransition: []bool{true},
			expectError:      false,
		},
		{
			name: "valid - two",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: vals[0], Amount: 1000000},
							{Valoper: vals[1], Amount: 1000000},
							{Valoper: vals[2], Amount: 1000000},
							{Valoper: vals[3], Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: vals[0], Amount: 5000000},
							{Valoper: vals[1], Amount: 2500000},
							{Valoper: vals[2], Amount: 5000000},
							{Valoper: vals[3], Amount: 2500000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[2],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[3],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					},
				}
			},
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.RedelegationRecord {
				return []icstypes.RedelegationRecord{}
			},
			expectTransition: []bool{true, true},
			expectError:      false,
		},
		{
			name: "invalid - locked tokens",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: vals[0], Amount: 1000000},
							{Valoper: vals[1], Amount: 1000000},
							{Valoper: vals[2], Amount: 1000000},
							{Valoper: vals[3], Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[2],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[3],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000)),
					},
				}
			},
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.RedelegationRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.RedelegationRecord{
					{
						ChainId:        zone.ChainId,
						EpochNumber:    1,
						Source:         vals[3],
						Destination:    vals[0],
						Amount:         500000,
						CompletionTime: time.Now().Add(time.Hour),
					},
				}
			},
			expectTransition: []bool{false},
			expectError:      false,
		},
		{
			name: "mixed - locked tokens but both succeed (previously failed)",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: vals[0], Amount: 5000000},
							{Valoper: vals[1], Amount: 2500000},
							{Valoper: vals[2], Amount: 5000000},
							{Valoper: vals[3], Amount: 2500000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: vals[0], Amount: 1000000},
							{Valoper: vals[1], Amount: 1000000},
							{Valoper: vals[2], Amount: 1000000},
							{Valoper: vals[3], Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[2],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[3],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(6000000)),
					},
				}
			},
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *icstypes.Zone) []icstypes.RedelegationRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.RedelegationRecord{
					{
						ChainId:        zone.ChainId,
						EpochNumber:    1,
						Source:         vals[3],
						Destination:    vals[0],
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
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			if !found {
				suite.Fail("unable to retrieve zone for test")
			}
			records := test.records(ctx, quicksilver, &zone)
			delegations := test.delegations(ctx, quicksilver, &zone)
			redelegations := test.redelegations(ctx, quicksilver, &zone)

			// set up zones
			for _, record := range records {
				quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
			}

			for _, delegation := range delegations {
				quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
				valAddrBytes, err := addressutils.ValAddressFromBech32(delegation.ValidatorAddress, zone.GetValoperPrefix())
				suite.Require().NoError(err)
				val, _ := quicksilver.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
				val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
				val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))
			}

			for _, redelegation := range redelegations {
				quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, redelegation)
			}

			quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

			// trigger handler
			err := quicksilver.InterchainstakingKeeper.HandleQueuedUnbondings(ctx, &zone, 1)
			if test.expectError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			for idx, record := range records {
				// check record with old status is opposite to expectedTransition (if false, this record should exist in status 3)
				_, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusQueued)
				suite.Require().Equal(!test.expectTransition[idx], found)
				// check record with new status is as per expectedTransition (if false, this record should not exist in status 4)
				_, found = quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusUnbond)
				suite.Require().Equal(test.expectTransition[idx], found)

				if test.expectTransition[idx] {
					actualRecord, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusUnbond)
					suite.Require().True(found)
					for _, unbonding := range actualRecord.Distribution {
						r, found := quicksilver.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, unbonding.Valoper, 1)
						suite.Require().True(found)
						suite.Require().Contains(r.RelatedTxhash, record.Txhash)
					}
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleWithdrawForUser() {
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
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
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
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
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
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusSend,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 5000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1250000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 5000000},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: 1250000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
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
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			if !found {
				suite.Fail("unable to retrieve zone for test")
			}

			records := test.records(&zone)

			// set up zones
			for _, record := range records {
				quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				err := quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(record.BurnAmount))
				suite.Require().NoError(err)
				err = quicksilver.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount))
				suite.Require().NoError(err)
			}

			// trigger handler
			err := quicksilver.InterchainstakingKeeper.HandleWithdrawForUser(ctx, &zone, &test.message, test.memo)
			if test.err {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusSend, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash == test.memo {
					suite.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusCompleted, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash != test.memo {
					suite.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})
		})
	}
}

func (suite *KeeperTestSuite) TestHandleWithdrawForUserLSM() {
	v1 := addressutils.GenerateValAddressForTest().String()
	v2 := addressutils.GenerateValAddressForTest().String()
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
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: v1, Amount: 1000000},
							{Valoper: v2, Amount: 1000000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
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
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{Valoper: v1, Amount: 1000000},
							{Valoper: v2, Amount: 1500000},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
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
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			if !found {
				suite.Fail("unable to retrieve zone for test")
			}

			records := test.records(&zone)

			startBalance := quicksilver.BankKeeper.GetAllBalances(ctx, quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName))
			// set up zones
			for _, record := range records {
				quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				err := quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(record.BurnAmount))
				suite.Require().NoError(err)
				err = quicksilver.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount))
				suite.Require().NoError(err)
			}

			// trigger handler
			for i := range test.message {
				err := quicksilver.InterchainstakingKeeper.HandleWithdrawForUser(ctx, &zone, &test.message[i], test.memo)
				if test.err {
					suite.Require().Error(err)
				} else {
					suite.Require().NoError(err)
				}
			}

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusSend, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash == test.memo {
					suite.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icskeeper.WithdrawStatusCompleted, func(idx int64, withdrawal icstypes.WithdrawalRecord) bool {
				if withdrawal.Txhash != test.memo {
					suite.Require().Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			postBurnBalance := quicksilver.BankKeeper.GetAllBalances(ctx, quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName))
			suite.Require().Equal(startBalance, postBurnBalance)
		})
	}
}

func (suite *KeeperTestSuite) TestReceiveAckErrForBeginRedelegate() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	validators := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	// create redelegation record
	record := icstypes.RedelegationRecord{
		ChainId:     suite.chainB.ChainID,
		EpochNumber: 1,
		Source:      validators[0].ValoperAddress,
		Destination: validators[1].ValoperAddress,
		Amount:      1000,
	}

	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	redelegate := &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: validators[0].ValoperAddress, ValidatorDstAddress: validators[1].ValoperAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []sdk.Msg{redelegate})
	suite.Require().NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: fmt.Sprintf("rebalance/%d", 1),
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	ackBytes := []byte("{\"error\":\"ABCI code: 32: error handling packet on host chain: see events for details\"}")
	// call handler

	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, 1)
	suite.Require().True(found)

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes)
	suite.Require().NoError(err)

	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, 1)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestReceiveAckErrForBeginUndelegate() {
	hash1 := fmt.Sprintf("%x", sha256.Sum256([]byte{0x01}))
	hash2 := fmt.Sprintf("%x", sha256.Sum256([]byte{0x02}))
	hash3 := fmt.Sprintf("%x", sha256.Sum256([]byte{0x03}))
	delegator1 := addressutils.GenerateAccAddressForTest().String()
	delegator2 := addressutils.GenerateAccAddressForTest().String()

	tests := []struct {
		name                      string
		epoch                     int64
		withdrawalRecords         func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord
		unbondingRecords          func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.UnbondingRecord
		msgs                      func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []sdk.Msg
		expectedWithdrawalRecords func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord
	}{
		{
			name:  "1 wdr, 2 vals, 1k+1k, 1800 qasset",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  1000,
							},
							{
								Valoper: vals[1],
								Amount:  1000,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.UnbondingRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.UnbondingRecord{
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   1,
						Validator:     vals[0],
						RelatedTxhash: []string{hash1},
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[0],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[1],
								Amount:  1000,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:      suite.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
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
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  1000,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.UnbondingRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.UnbondingRecord{
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   1,
						Validator:     vals[0],
						RelatedTxhash: []string{hash1},
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[0],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:      suite.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
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
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  1000,
							},
							{
								Valoper: vals[1],
								Amount:  500,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1500))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1350)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  1000,
							},
							{
								Valoper: vals[1],
								Amount:  2000,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(2700)),
						Txhash:     hash2,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  600,
							},
							{
								Valoper: vals[1],
								Amount:  400,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash3,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.UnbondingRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.UnbondingRecord{
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   2,
						Validator:     vals[1],
						RelatedTxhash: []string{hash1, hash2, hash3},
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[1],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2900)),
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  1000,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  1000,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash2,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*icstypes.Distribution{
							{
								Valoper: vals[0],
								Amount:  600,
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(600))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(540)),
						Txhash:     hash3,
						Status:     icskeeper.WithdrawStatusUnbond,
					},
					{
						ChainId:      suite.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(500))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(450)),
						Txhash:       fmt.Sprintf("%064d", 1),
						Status:       icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:      suite.chainB.ChainID,
						Delegator:    delegator2,
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:       fmt.Sprintf("%064d", 2),
						Status:       icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:      suite.chainB.ChainID,
						Delegator:    delegator1,
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(400))),
						BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewInt(360)),
						Txhash:       fmt.Sprintf("%064d", 3),
						Status:       icskeeper.WithdrawStatusQueued,
					},
				}
			},
		},
		// TODO: fix this test
		// {
		//	name:  "2 wdr, random_rr, 1 vals, 1k; 2 vals; 123 + 456 ",
		//	epoch: 1,
		//	withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
		//		vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
		//		return []icstypes.WithdrawalRecord{
		//			{
		//				ChainId:   s.chainB.ChainID,
		//				Delegator: delegator1,
		//				Distribution: []*icstypes.Distribution{
		//					{
		//						Valoper: vals[0],
		//						Amount:  1000,
		//					},
		//				},
		//				Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
		//				Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
		//				BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewDec(1000).Quo(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
		//				Txhash:     hash1,
		//				Status:     icskeeper.WithdrawStatusUnbond,
		//			},
		//			{
		//				ChainId:   s.chainB.ChainID,
		//				Delegator: delegator2,
		//				Distribution: []*icstypes.Distribution{
		//					{
		//						Valoper: vals[1],
		//						Amount:  123,
		//					},
		//					{
		//						Valoper: vals[2],
		//						Amount:  456,
		//					},
		//				},
		//				Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
		//				Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(579))),
		//				BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewDec(579).Quo(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
		//				Txhash:     hash2,
		//				Status:     icskeeper.WithdrawStatusUnbond,
		//			},
		//		}
		//	},
		//	unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.UnbondingRecord {
		//		vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
		//		return []icstypes.UnbondingRecord{
		//			{
		//				ChainId:       s.chainB.ChainID,
		//				EpochNumber:   1,
		//				Validator:     vals[0],
		//				RelatedTxhash: []string{hash1},
		//			},
		//			{
		//				ChainId:       s.chainB.ChainID,
		//				EpochNumber:   1,
		//				Validator:     vals[1],
		//				RelatedTxhash: []string{hash2},
		//			},
		//			// {
		//			// 	ChainID:       s.chainB.ChainID,
		//			// 	EpochNumber:   1,
		//			// 	Validator:     vals[2],
		//			// 	RelatedTxhash: []string{hash2},
		//			// },
		//		}
		//	},
		//	msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []sdk.Msg {
		//		vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
		//		return []sdk.Msg{
		//			&stakingtypes.MsgUndelegate{
		//				DelegatorAddress: zone.DelegationAddress.Address,
		//				ValidatorAddress: vals[0],
		//				Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
		//			},
		//			&stakingtypes.MsgUndelegate{
		//				DelegatorAddress: zone.DelegationAddress.Address,
		//				ValidatorAddress: vals[1],
		//				Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(123)),
		//			},
		//		}
		//	},
		//	expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
		//		vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
		//		return []icstypes.WithdrawalRecord{
		//			{
		//				ChainId:      s.chainB.ChainID,
		//				Delegator:    delegator1,
		//				Distribution: nil,
		//				Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
		//				Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
		//				BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewDec(1000).Quo(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
		//				Txhash:       hash1,
		//				Status:       icskeeper.WithdrawStatusQueued,
		//			},
		//			{
		//				ChainId:      s.chainB.ChainID,
		//				Delegator:    delegator2,
		//				Distribution: nil,
		//				Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
		//				Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(123))),
		//				BurnAmount:   sdk.NewCoin(zone.LocalDenom, sdk.NewDec(123).Quo(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
		//				Txhash:       fmt.Sprintf("%064d", 1),
		//				Status:       icskeeper.WithdrawStatusQueued,
		//			},
		//			{
		//				ChainId:   s.chainB.ChainID,
		//				Delegator: delegator2,
		//				Distribution: []*icstypes.Distribution{
		//					{
		//						Valoper: vals[2],
		//						Amount:  456,
		//					},
		//				},
		//				Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
		//				Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(456))),
		//				BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewDec(456).Quo(sdk.MustNewDecFromStr(fmt.Sprintf("%f", randRr))).TruncateInt()),
		//				Txhash:     hash2,
		//				Status:     icskeeper.WithdrawStatusUnbond,
		//			},
		//		}
		//	},
		// },
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			if !found {
				suite.Fail("unable to retrieve zone for test")
			}

			for _, wdr := range test.withdrawalRecords(ctx, quicksilver, zone) {
				quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, wdr)
			}

			for _, ubr := range test.unbondingRecords(ctx, quicksilver, zone) {
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr)
			}

			data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), test.msgs(ctx, quicksilver, zone))
			suite.Require().NoError(err)

			// validate memo < 256 bytes
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: data,
				Memo: fmt.Sprintf("withdrawal/%d", test.epoch),
			}

			packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

			ackBytes := []byte("{\"error\":\"ABCI code: 32: error handling packet on host chain: see events for details\"}")
			// call handler

			for _, ubr := range test.unbondingRecords(ctx, quicksilver, zone) {
				_, found = quicksilver.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, ubr.Validator, test.epoch)
				suite.Require().True(found)
			}

			err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes)
			suite.Require().NoError(err)

			for _, ubr := range test.unbondingRecords(ctx, quicksilver, zone) {
				_, found = quicksilver.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, ubr.Validator, test.epoch)
				suite.Require().False(found)
			}

			for idx, ewdr := range test.expectedWithdrawalRecords(ctx, quicksilver, zone) {
				wdr, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, ewdr.Txhash, ewdr.Status)
				suite.Require().True(found)
				suite.Require().Equal(ewdr.Amount, wdr.Amount)
				suite.Require().Equal(ewdr.BurnAmount, wdr.BurnAmount)
				suite.Require().Equal(ewdr.Delegator, wdr.Delegator)
				suite.Require().Equal(ewdr.Distribution, wdr.Distribution, idx)
				suite.Require().Equal(ewdr.Status, wdr.Status)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRebalanceDueToIntentChange() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)

	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	vals := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	for _, val := range vals {
		valoper, _ := addressutils.ValAddressFromBech32(val.ValoperAddress, "cosmosvaloper")
		quicksilver.InterchainstakingKeeper.DeleteValidator(ctx, zone.ChainId, valoper)
	}

	val0 := icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val0)
	val1 := icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val1)
	val2 := icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val2)
	val3 := icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val3)

	vals = quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)

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
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[3].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
	}
	for _, delegation := range delegations {
		quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		addressBytes, err := addressutils.ValAddressFromBech32(delegation.ValidatorAddress, zone.GetValoperPrefix())
		suite.Require().NoError(err)
		val, _ := quicksilver.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, addressBytes)
		val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
		val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))
	}

	quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

	// trigger rebalance
	err := quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 1)
	suite.Require().NoError(err)

	// change intents to trigger redelegations from val[3]
	intents := icstypes.ValidatorIntents{
		{ValoperAddress: vals[0].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[1].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[2].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[3].ValoperAddress, Weight: sdk.NewDecWithPrec(1, 1)},
	}
	zone.AggregateIntent = intents

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 2)
	suite.Require().NoError(err)

	// mock ack for redelegations
	quicksilver.InterchainstakingKeeper.IteratePrefixedRedelegationRecords(ctx, []byte(zone.ChainId), func(idx int64, _ []byte, record icstypes.RedelegationRecord) (stop bool) {
		if record.EpochNumber == 2 {
			msg := stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    zone.DelegationAddress.Address,
				ValidatorSrcAddress: record.Source,
				ValidatorDstAddress: record.Destination,
				Amount:              sdk.NewCoin("uatom", sdkmath.NewInt(record.Amount)),
			}
			err := quicksilver.InterchainstakingKeeper.HandleBeginRedelegate(ctx, &msg, time.Now().Add(time.Hour*24*7), fmt.Sprintf("rebalance/%d", 2))
			if err != nil {
				return false
			}
		}
		return false
	})

	// check for redelegations
	_, present := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[0].ValoperAddress, 2)
	suite.Require().True(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[1].ValoperAddress, 2)
	suite.Require().True(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[2].ValoperAddress, 2)
	suite.Require().True(present)

	// change intents to trigger transitive redelegations which should fail rebalance
	zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	intents = icstypes.ValidatorIntents{
		{ValoperAddress: vals[0].ValoperAddress, Weight: sdk.NewDecWithPrec(1, 1)},
		{ValoperAddress: vals[1].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[2].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[3].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
	}
	zone.AggregateIntent = intents

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 3)
	suite.Require().NoError(err)

	// check for redelegations originating from val[0], they should not be present
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[1].ValoperAddress, 3)
	suite.Require().False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[2].ValoperAddress, 3)
	suite.Require().False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[3].ValoperAddress, 3)
	suite.Require().False(present)
}

func (suite *KeeperTestSuite) TestRebalanceDueToDelegationChange() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	vals := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	for _, val := range vals {
		valoper, _ := addressutils.ValAddressFromBech32(val.ValoperAddress, "cosmosvaloper")
		quicksilver.InterchainstakingKeeper.DeleteValidator(ctx, zone.ChainId, valoper)
	}

	val0 := icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val0)
	val1 := icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val1)
	val2 := icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val2)
	val3 := icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val3)
	delegations := []icstypes.Delegation{
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  val0.ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  val1.ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  val2.ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  val3.ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000)),
			RedelegationEnd:   0,
		},
	}
	for _, delegation := range delegations {
		quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		valAddrBytes, err := addressutils.ValAddressFromBech32(delegation.ValidatorAddress, zone.GetValoperPrefix())
		suite.Require().NoError(err)

		val, found := quicksilver.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		suite.Require().NoError(err)
		suite.Require().True(found)
		val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
		val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))

	}

	// trigger rebalance
	err := quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 1)
	suite.Require().NoError(err)

	quicksilver.InterchainstakingKeeper.IterateAllDelegations(ctx, &zone, func(delegation icstypes.Delegation) bool {
		if delegation.ValidatorAddress == val0.ValoperAddress {
			delegation.Amount = delegation.Amount.Add(sdk.NewInt64Coin("uatom", 4000))
			quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		}
		return false
	})

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 2)
	suite.Require().NoError(err)

	// mock ack for redelegations
	quicksilver.InterchainstakingKeeper.IteratePrefixedRedelegationRecords(ctx, []byte(zone.ChainId), func(idx int64, _ []byte, record icstypes.RedelegationRecord) (stop bool) {
		if record.EpochNumber == 2 {
			msg := stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    zone.DelegationAddress.Address,
				ValidatorSrcAddress: record.Source,
				ValidatorDstAddress: record.Destination,
				Amount:              sdk.NewCoin("uatom", sdkmath.NewInt(record.Amount)),
			}
			err := quicksilver.InterchainstakingKeeper.HandleBeginRedelegate(ctx, &msg, time.Now().Add(time.Hour*24*7), fmt.Sprintf("rebalance/%d", 2))
			if err != nil {
				return false
			}
		}
		return false
	})

	// check for redelegations
	_, present := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val0.ValoperAddress, val1.ValoperAddress, 2)
	suite.Require().False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val0.ValoperAddress, val2.ValoperAddress, 2)
	suite.Require().False(present)

	// change validator delegation to trigger transitive redelegations which should fail rebalance
	quicksilver.InterchainstakingKeeper.IterateAllDelegations(ctx, &zone, func(delegation icstypes.Delegation) bool {
		if delegation.ValidatorAddress == val0.ValoperAddress {
			delegation.Amount = delegation.Amount.Sub(sdk.NewInt64Coin("uatom", 4000))
			quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		}
		if delegation.ValidatorAddress == val1.ValoperAddress {
			delegation.Amount = delegation.Amount.Add(sdk.NewInt64Coin("uatom", 4000))
			quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
		}

		return false
	})

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 3)
	suite.Require().NoError(err)

	// check for redelegations originating from val[1], they should not be present
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val1.ValoperAddress, val0.ValoperAddress, 3)
	suite.Require().False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val1.ValoperAddress, val1.ValoperAddress, 3)
	suite.Require().False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val1.ValoperAddress, val3.ValoperAddress, 3)
	suite.Require().False(present)
}

func (suite *KeeperTestSuite) Test_v045Callback() {
	tests := []struct {
		name             string
		setStatements    func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]sdk.Msg, []byte)
		assertStatements func(ctx sdk.Context, quicksilver *app.Quicksilver) bool
	}{
		{
			name: "msg response with some data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]sdk.Msg, []byte) {
				err := quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
				suite.Require().NoError(err)
				sender := addressutils.GenerateAccAddressForTest()
				senderAddr, err := sdk.Bech32ifyAddressBytes("cosmos", sender)
				suite.Require().NoError(err)

				transferMsg := ibctransfertypes.MsgTransfer{
					SourcePort:    "transfer",
					SourceChannel: "channel-0",
					Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
					Sender:        senderAddr,
					Receiver:      quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName).String(),
				}
				response := ibctransfertypes.MsgTransferResponse{
					Sequence: 1,
				}

				respBytes := icatypes.ModuleCdc.MustMarshal(&response)
				return []sdk.Msg{&transferMsg}, respBytes
			},
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				txMacc := quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName)
				feeMacc := quicksilver.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
				txMaccBalance2 := quicksilver.BankKeeper.GetAllBalances(ctx, txMacc)
				feeMaccBalance2 := quicksilver.BankKeeper.GetAllBalances(ctx, feeMacc)

				// assert that ics module balance is now 100denom less than before HandleMsgTransfer()

				if txMaccBalance2.AmountOf("denom").Equal(sdk.ZeroInt()) && feeMaccBalance2.AmountOf("denom").Equal(sdk.NewInt(100)) {
					return true
				}
				return false
			},
		},
		{
			name: "msg response with nil data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]sdk.Msg, []byte) {
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				if !found {
					suite.Fail("unable to retrieve zone for test")
				}

				msgSetWithdrawAddress := distrtypes.MsgSetWithdrawAddress{
					DelegatorAddress: zone.PerformanceAddress.Address,
					WithdrawAddress:  zone.WithdrawalAddress.Address,
				}

				response := distrtypes.MsgSetWithdrawAddressResponse{}

				respBytes := icatypes.ModuleCdc.MustMarshal(&response)
				return []sdk.Msg{&msgSetWithdrawAddress}, respBytes
			},
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				if !found {
					suite.Fail("unable to retrieve zone for test")
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
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			msg, msgResponseBytes := test.setStatements(ctx, quicksilver)

			txMsgData := &sdk.TxMsgData{
				// we need to support this older deprecated type
				Data:         []*sdk.MsgData{{MsgType: "/bob", Data: msgResponseBytes}}, //nolint:staticcheck
				MsgResponses: []*codectypes.Any{},
			}

			ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

			acknowledgement := channeltypes.Acknowledgement{
				Response: &channeltypes.Acknowledgement_Result{
					Result: ackData,
				},
			}

			pdBytes, err := icatypes.SerializeCosmosTx(icatypes.ModuleCdc, msg)
			suite.Require().NoError(err)
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: pdBytes,
				Memo: "test_acknowledgement",
			}

			packetBytes, err := icatypes.ModuleCdc.MarshalJSON(&packetData)
			suite.Require().NoError(err)
			packet := channeltypes.Packet{
				Data: packetBytes,
			}

			suite.Require().NoError(quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, icatypes.ModuleCdc.MustMarshalJSON(&acknowledgement)))

			suite.Require().True(test.assertStatements(ctx, quicksilver))
		})
	}
}

func (suite *KeeperTestSuite) Test_v046Callback() {
	tests := []struct {
		name             string
		setStatements    func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]sdk.Msg, *codectypes.Any)
		assertStatements func(ctx sdk.Context, quicksilver *app.Quicksilver) bool
	}{
		{
			name: "msg response with some data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]sdk.Msg, *codectypes.Any) {
				err := quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
				suite.Require().NoError(err)
				sender := addressutils.GenerateAccAddressForTest()
				senderAddr, err := sdk.Bech32ifyAddressBytes("cosmos", sender)
				suite.Require().NoError(err)

				transferMsg := ibctransfertypes.MsgTransfer{
					SourcePort:    "transfer",
					SourceChannel: "channel-0",
					Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
					Sender:        senderAddr,
					Receiver:      quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName).String(),
				}
				response := ibctransfertypes.MsgTransferResponse{
					Sequence: 1,
				}

				anyResponse, err := codectypes.NewAnyWithValue(&response)
				suite.Require().NoError(err)
				return []sdk.Msg{&transferMsg}, anyResponse
			},
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				txMacc := quicksilver.AccountKeeper.GetModuleAddress(icstypes.ModuleName)
				feeMacc := quicksilver.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
				txMaccBalance2 := quicksilver.BankKeeper.GetAllBalances(ctx, txMacc)
				feeMaccBalance2 := quicksilver.BankKeeper.GetAllBalances(ctx, feeMacc)

				// assert that ics module balance is now 100denom less than before HandleMsgTransfer()

				if txMaccBalance2.AmountOf("denom").Equal(sdk.ZeroInt()) && feeMaccBalance2.AmountOf("denom").Equal(sdk.NewInt(100)) {
					return true
				}
				return false
			},
		},
		{
			name: "msg response with nil data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]sdk.Msg, *codectypes.Any) {
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				if !found {
					suite.Fail("unable to retrieve zone for test")
				}

				msgSetWithdrawAddress := distrtypes.MsgSetWithdrawAddress{
					DelegatorAddress: zone.PerformanceAddress.Address,
					WithdrawAddress:  zone.WithdrawalAddress.Address,
				}

				response := distrtypes.MsgSetWithdrawAddressResponse{}

				anyResponse, err := codectypes.NewAnyWithValue(&response)
				suite.Require().NoError(err)
				return []sdk.Msg{&msgSetWithdrawAddress}, anyResponse
			},
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				if !found {
					suite.Fail("unable to retrieve zone for test")
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
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			msg, anyResp := test.setStatements(ctx, quicksilver)

			txMsgData := &sdk.TxMsgData{
				Data:         []*sdk.MsgData{}, //nolint:staticcheck
				MsgResponses: []*codectypes.Any{anyResp},
			}

			ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

			acknowledgement := channeltypes.Acknowledgement{
				Response: &channeltypes.Acknowledgement_Result{
					Result: ackData,
				},
			}

			pdBytes, err := icatypes.SerializeCosmosTx(icatypes.ModuleCdc, msg)
			suite.Require().NoError(err)
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: pdBytes,
				Memo: "test_acknowledgement",
			}

			packetBytes, err := icatypes.ModuleCdc.MarshalJSON(&packetData)
			suite.Require().NoError(err)
			packet := channeltypes.Packet{
				Data: packetBytes,
			}

			suite.Require().NoError(quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, icatypes.ModuleCdc.MustMarshalJSON(&acknowledgement)))

			suite.Require().True(test.assertStatements(ctx, quicksilver))
		})
	}
}
