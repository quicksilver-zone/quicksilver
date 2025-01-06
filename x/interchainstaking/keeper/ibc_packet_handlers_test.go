package keeper_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/utils/randomutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/x/lsmtypes"
)

var TestChannel = channeltypes.Channel{
	State:          channeltypes.OPEN,
	Ordering:       channeltypes.UNORDERED,
	Counterparty:   channeltypes.NewCounterparty("transfer", "channel-0"),
	ConnectionHops: []string{"connection-0"},
}

const queryAllBalancesPath = "cosmos.bank.v1beta1.Query/AllBalances"

func (suite *KeeperTestSuite) TestHandleMsgTransferGood() {
	user1 := addressutils.GenerateAddressForTestWithPrefix("quick")
	user2 := addressutils.GenerateAddressForTestWithPrefix("quick")
	user3 := addressutils.GenerateAddressForTestWithPrefix("quick")

	tcs := []struct {
		name                 string
		amount               sdk.Coin
		claims               []cmtypes.Claim
		expectedFeeCollector math.Int
		supply               sdk.Coin
		expectedBalances     map[string]math.Int
		feeAmount            sdk.Dec
	}{
		{
			name:                 "staking denom - 100% to fc",
			amount:               sdk.NewCoin("uatom", math.NewInt(100)),
			supply:               sdk.NewCoin("uqatom", math.NewInt(1000)),
			expectedFeeCollector: math.NewInt(100),
		},
		{
			name:                 "non staking denom, no claims - 100% to fc",
			amount:               sdk.NewCoin("transfer/channel-569/untrn", math.NewInt(100)),
			supply:               sdk.NewCoin("uqatom", math.NewInt(1000)),
			expectedFeeCollector: math.NewInt(100),
		},
		{
			name:   "non staking denom, claims 100%, default (2.5%) to fc, remainder proportional to users",
			amount: sdk.NewCoin("ujuno", math.NewInt(100)),
			claims: []cmtypes.Claim{
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "osmosis-1", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "quicksilver-2", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "juno-1", Amount: math.NewInt(2000)},
			},
			supply:               sdk.NewCoin("uqatom", math.NewInt(4000)),
			expectedFeeCollector: math.NewInt(3),
			expectedBalances: map[string]math.Int{
				user1: math.NewInt(97),
			},
		},
		{
			name:   "non staking denom, claims 100%, non-default (9%) to fc, remainder proportional to users",
			amount: sdk.NewCoin("ujuno", math.NewInt(1000)),
			claims: []cmtypes.Claim{
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "osmosis-1", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "quicksilver-2", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "juno-1", Amount: math.NewInt(2000)},
				{UserAddress: user2, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "juno-1", Amount: math.NewInt(2000)},
			},
			supply:               sdk.NewCoin("uqatom", math.NewInt(6000)),
			expectedFeeCollector: math.NewInt(92),
			expectedBalances: map[string]math.Int{
				user1: math.NewInt(605),
				user2: math.NewInt(303),
			},
			feeAmount: sdk.NewDecWithPrec(9, 2),
		},
		{
			name:   "non staking denom, claims >100% - default (2.5%) to fc, remainder proportional to users",
			amount: sdk.NewCoin("ujuno", math.NewInt(1000)),
			claims: []cmtypes.Claim{
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "osmosis-1", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "quicksilver-2", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "juno-1", Amount: math.NewInt(2000)},
				{UserAddress: user2, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "juno-1", Amount: math.NewInt(2000)},
				{UserAddress: user3, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "quicksilver-2", Amount: math.NewInt(1000)},
			},
			supply:               sdk.NewCoin("uqatom", math.NewInt(6000)),
			expectedFeeCollector: math.NewInt(27),
			expectedBalances: map[string]math.Int{
				user1: math.NewInt(556),
				user2: math.NewInt(278),
				user3: math.NewInt(139),
			},
		},
		{
			name:   "non staking denom, claims <100% - default (2.5%) to fc, remainder proportional to users",
			amount: sdk.NewCoin("ujuno", math.NewInt(1000)),
			claims: []cmtypes.Claim{
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "osmosis-1", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "quicksilver-2", Amount: math.NewInt(1000)},
				{UserAddress: user1, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "juno-1", Amount: math.NewInt(1000)},
				{UserAddress: user2, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "juno-1", Amount: math.NewInt(1000)},
				{UserAddress: user3, ChainId: suite.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "quicksilver-2", Amount: math.NewInt(1000)},
			},
			supply:               sdk.NewCoin("uqatom", math.NewInt(6000)),
			expectedFeeCollector: math.NewInt(190),
			expectedBalances: map[string]math.Int{
				user1: math.NewInt(486),
				user2: math.NewInt(162),
				user3: math.NewInt(162),
			},
		},
	}
	for _, tc := range tcs {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			quicksilver.InterchainstakingKeeper.IBCKeeper.ChannelKeeper.SetChannel(ctx, "transfer", "channel-0", TestChannel)
			channel, cfound := quicksilver.InterchainstakingKeeper.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", "channel-0")
			suite.True(cfound)

			ibcDenom := utils.DeriveIbcDenom("transfer", "channel-0", channel.Counterparty.PortId, channel.Counterparty.ChannelId, tc.amount.Denom)

			err := quicksilver.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ibcDenom, tc.amount.Amount)))
			suite.NoError(err)

			if !tc.feeAmount.IsNil() {
				params := quicksilver.InterchainstakingKeeper.GetParams(ctx)
				params.CommissionRate.Set(tc.feeAmount)
				quicksilver.InterchainstakingKeeper.SetParams(ctx, params)
			}

			for _, claim := range tc.claims {
				quicksilver.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &claim)
			}

			err = quicksilver.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(tc.supply))
			suite.NoError(err)

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			sender := zone.WithdrawalAddress.Address
			suite.NoError(err)

			txMacc := quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName)
			feeMacc := quicksilver.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)

			transferPacket := ibctransfertypes.FungibleTokenPacketData{
				Amount:   tc.amount.Amount.String(),
				Denom:    tc.amount.Denom,
				Sender:   sender,
				Receiver: quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName).String(),
			}

			suite.NoError(quicksilver.InterchainstakingKeeper.HandleMsgTransfer(ctx, transferPacket, utils.DeriveIbcDenom("transfer", "channel-0", channel.Counterparty.PortId, channel.Counterparty.ChannelId, tc.amount.Denom)))

			txMaccBalance := quicksilver.BankKeeper.GetAllBalances(ctx, txMacc)
			feeMaccBalance := quicksilver.BankKeeper.GetAllBalances(ctx, feeMacc)
			suite.NoError(err)

			// assert that ics module balance is nil
			suite.Equal(sdk.Coins{}, txMaccBalance)

			// assert that fee collector module balance is the expected value
			suite.Equal(feeMaccBalance.AmountOf(ibcDenom), tc.expectedFeeCollector)

			for address, balance := range tc.expectedBalances {
				suite.Equal(quicksilver.BankKeeper.GetBalance(ctx, addressutils.MustAccAddressFromBech32(address, ""), ibcDenom), sdk.NewCoin(ibcDenom, balance))
			}
		})
	}
}

func TestHandleMsgTransferBadRecipient(t *testing.T) {
	recipient := addressutils.GenerateAccAddressForTest()
	quicksilver, ctx := app.GetAppWithContext(t, true)

	sender := addressutils.GenerateAccAddressForTest()
	senderAddr, err := sdk.Bech32ifyAddressBytes("cosmos", sender)
	require.NoError(t, err)

	transferMsg := ibctransfertypes.FungibleTokenPacketData{
		Denom:    "denom",
		Amount:   "100",
		Sender:   senderAddr,
		Receiver: recipient.String(),
	}
	require.Error(t, quicksilver.InterchainstakingKeeper.HandleMsgTransfer(ctx, transferMsg, "raa"))
}

func (suite *KeeperTestSuite) TestHandleQueuedUnbondings() {
	tests := []struct {
		name             string
		records          func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.WithdrawalRecord
		delegations      func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.Delegation
		redelegations    func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.RedelegationRecord
		expectTransition []bool
		expectError      bool
	}{
		{
			name: "valid",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:       "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:       types.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
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
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.RedelegationRecord {
				return []types.RedelegationRecord{}
			},
			expectTransition: []bool{true},
			expectError:      false,
		},
		{
			name: "valid - int64 overflow",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(4000000000000).Mul(sdk.NewInt(4000000000000))),
						Txhash:       "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:       types.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000000000).Mul(sdk.NewInt(4000000000000))),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000000000).Mul(sdk.NewInt(4000000000000))),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[2],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000000000).Mul(sdk.NewInt(4000000000000))),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[3],
						Amount:            sdk.NewCoin("uatom", sdk.NewInt(1000000000000).Mul(sdk.NewInt(4000000000000))),
					},
				}
			},
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.RedelegationRecord {
				return []types.RedelegationRecord{}
			},
			expectTransition: []bool{true},
			expectError:      false,
		},
		{
			name: "valid - two",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:       "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:       types.WithdrawStatusQueued,
					},
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:       "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:       types.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
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
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.RedelegationRecord {
				return []types.RedelegationRecord{}
			},
			expectTransition: []bool{true, true},
			expectError:      false,
		},
		{
			name: "invalid - locked tokens",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:       "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:       types.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
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
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.RedelegationRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.RedelegationRecord{
					{
						ChainId:        zone.ChainId,
						EpochNumber:    1,
						Source:         vals[3],
						Destination:    vals[0],
						Amount:         math.NewInt(500000),
						CompletionTime: time.Now().Add(time.Hour),
					},
				}
			},
			expectTransition: []bool{false},
			expectError:      false,
		},
		{
			name: "mixed - locked tokens but both succeed (previously failed)",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:       "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:       types.WithdrawStatusQueued,
					},
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:       "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:       types.WithdrawStatusQueued,
					},
				}
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
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
			redelegations: func(ctx sdk.Context, qs *app.Quicksilver, zone *types.Zone) []types.RedelegationRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.RedelegationRecord{
					{
						ChainId:        zone.ChainId,
						EpochNumber:    1,
						Source:         vals[3],
						Destination:    vals[0],
						Amount:         math.NewInt(1000001),
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
				err := quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				suite.NoError(err)
			}

			for _, delegation := range delegations {
				quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
				valAddrBytes, err := addressutils.ValAddressFromBech32(delegation.ValidatorAddress, zone.GetValoperPrefix())
				suite.NoError(err)
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
				suite.Error(err)
			} else {
				suite.NoError(err)
			}

			for idx, record := range records {
				// check record with old status is opposite to expectedTransition (if false, this record should exist in status 3)
				_, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, types.WithdrawStatusQueued)
				suite.Equal(!test.expectTransition[idx], found)
				// check record with new status is as per expectedTransition (if false, this record should not exist in status 4)
				_, found = quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, types.WithdrawStatusUnbond)
				suite.Equal(test.expectTransition[idx], found)

				if test.expectTransition[idx] {
					actualRecord, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, types.WithdrawStatusUnbond)
					suite.True(found)
					for _, unbonding := range actualRecord.Distribution {
						r, found := quicksilver.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, unbonding.Valoper, 1)
						suite.True(found)
						suite.Contains(r.RelatedTxhash, record.Txhash)
					}
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleWithdrawForUser() {
	tests := []struct {
		name    string
		records func(zone *types.Zone) []types.WithdrawalRecord
		message banktypes.MsgSend
		memo    string
		err     bool
	}{
		{
			name: "invalid - no matching record",
			records: func(zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:      zone.ChainId,
						Delegator:    addressutils.GenerateAccAddressForTest().String(),
						Distribution: nil,
						Recipient:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:       nil,
						BurnAmount:   sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:       "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:       types.WithdrawStatusQueued,
					},
				}
			},
			message: banktypes.MsgSend{},
			memo:    "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:     true,
		},
		{
			name: "valid",
			records: func(zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*types.Distribution{
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     types.WithdrawStatusSend,
					},
				}
			},
			message: banktypes.MsgSend{
				Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
			},
			memo: "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:  false,
		},
		{
			name: "valid - two",
			records: func(zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*types.Distribution{
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1000000)},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     types.WithdrawStatusSend,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*types.Distribution{
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(5000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1250000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(5000000)},
							{Valoper: addressutils.GenerateValAddressForTest().String(), Amount: math.NewInt(1250000)},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     types.WithdrawStatusSend,
					},
				}
			},
			message: banktypes.MsgSend{
				Amount: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
			},
			memo: "unbondSend/d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
			err:  false,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test := test
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
				err := quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				suite.NoError(err)
				err = quicksilver.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(record.BurnAmount))
				suite.NoError(err)
				err = quicksilver.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount))
				suite.NoError(err)
			}

			// trigger handler
			err := quicksilver.InterchainstakingKeeper.HandleWithdrawForUser(ctx, &zone, &test.message, test.memo)
			if test.err {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}

			hash, err := types.ParseTxMsgMemo(test.memo, types.MsgTypeUnbondSend)
			suite.NoError(err)

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusSend, func(idx int64, withdrawal types.WithdrawalRecord) bool {
				if withdrawal.Txhash == hash {
					suite.Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusCompleted, func(idx int64, withdrawal types.WithdrawalRecord) bool {
				if withdrawal.Txhash != hash {
					suite.Fail("unexpected withdrawal record; status should be Completed.")
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
		records func(zone *types.Zone) []types.WithdrawalRecord
		message []banktypes.MsgSend
		memo    string
		err     bool
	}{
		{
			name: "valid",
			records: func(zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*types.Distribution{
							{Valoper: v1, Amount: math.NewInt(1000000)},
							{Valoper: v2, Amount: math.NewInt(1000000)},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     types.WithdrawStatusSend,
					},
				}
			},
			message: []banktypes.MsgSend{
				{Amount: sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000)))},
				{Amount: sdk.NewCoins(sdk.NewCoin(v2+"2", sdk.NewInt(1000000)))},
			},
			memo: "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:  false,
		},
		{
			name: "valid - unequal",
			records: func(zone *types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: addressutils.GenerateAccAddressForTest().String(),
						Distribution: []*types.Distribution{
							{Valoper: v1, Amount: math.NewInt(1000000)},
							{Valoper: v2, Amount: math.NewInt(1500000)},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2500000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2500000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     types.WithdrawStatusSend,
					},
				}
			},
			message: []banktypes.MsgSend{
				{Amount: sdk.NewCoins(sdk.NewCoin(v2+"1", sdk.NewInt(1500000)))},
				{Amount: sdk.NewCoins(sdk.NewCoin(v1+"2", sdk.NewInt(1000000)))},
			},
			memo: "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
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

			startBalance := quicksilver.BankKeeper.GetAllBalances(ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
			// set up zones
			for _, record := range records {
				err := quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				suite.NoError(err)
				err = quicksilver.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(record.BurnAmount))
				suite.NoError(err)
				err = quicksilver.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount))
				suite.NoError(err)
			}

			// trigger handler
			for i := range test.message {
				err := quicksilver.InterchainstakingKeeper.HandleWithdrawForUser(ctx, &zone, &test.message[i], test.memo)
				if test.err {
					suite.Error(err)
				} else {
					suite.NoError(err)
				}
			}

			hash, err := types.ParseTxMsgMemo(test.memo, types.MsgTypeUnbondSend)
			suite.NoError(err)

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusSend, func(idx int64, withdrawal types.WithdrawalRecord) bool {
				if withdrawal.Txhash == hash {
					suite.Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			quicksilver.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusCompleted, func(idx int64, withdrawal types.WithdrawalRecord) bool {
				if withdrawal.Txhash != hash {
					suite.Fail("unexpected withdrawal record; status should be Completed.")
				}
				return false
			})

			postBurnBalance := quicksilver.BankKeeper.GetAllBalances(ctx, quicksilver.AccountKeeper.GetModuleAddress(types.ModuleName))
			suite.Equal(startBalance, postBurnBalance)
		})
	}
}

func (suite *KeeperTestSuite) TestHandleWithdrawRewards() {
	val := addressutils.GenerateValAddressForTest().String()
	user := addressutils.GenerateAddressForTestWithPrefix("quick")
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone)
		msg       func(zone *types.Zone) sdk.Msg
		checks    func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone)
		triggered bool
		err       bool
	}{
		{
			name:   "wrong msg",
			setup:  func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {},
			checks: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {},
			msg: func(zone *types.Zone) sdk.Msg {
				return &distrtypes.MsgWithdrawValidatorCommission{
					ValidatorAddress: val,
				}
			},
			triggered: false,
			err:       true,
		},
		{
			name: "wrong context",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {
				zone.ConnectionId = ""
				quicksilver.InterchainstakingKeeper.SetZone(ctx, zone)
			},
			checks: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {},
			msg: func(zone *types.Zone) sdk.Msg {
				return &distrtypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: user,
					ValidatorAddress: val,
				}
			},
			triggered: false,
			err:       true,
		},
		// try to decrement when waitgroup = 0
		{
			name: "try to decrement when waitgroup = 0",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {
				zone.SetWithdrawalWaitgroup(quicksilver.Logger(), 0, "init")
				quicksilver.InterchainstakingKeeper.SetZone(ctx, zone)
			},
			checks: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {},
			msg: func(zone *types.Zone) sdk.Msg {
				return &distrtypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: user,
					ValidatorAddress: val,
				}
			},
			triggered: false,
			err:       false, // was true but we don't fail on this case now in case historic messages needs to be delivered after wg zeroed.
		},
		{
			name: "valid case with balances != 0",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {
				zone.SetWithdrawalWaitgroup(quicksilver.Logger(), 1, "init")

				balances := sdk.NewCoins(
					sdk.NewCoin(
						zone.BaseDenom,
						math.NewInt(10_000_000),
					),
				)
				zone.DelegationAddress.Balance = balances
				quicksilver.InterchainstakingKeeper.SetZone(ctx, zone)
			},
			checks: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {},
			msg: func(zone *types.Zone) sdk.Msg {
				return &distrtypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: user,
					ValidatorAddress: val,
				}
			},
			triggered: true,
			err:       false,
		},
		{
			name: "valid case trigger redemption rate and check if delegatorAddress == performanceAddress",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {
				zone.SetWithdrawalWaitgroup(quicksilver.Logger(), 1, "init")
				quicksilver.InterchainstakingKeeper.SetZone(ctx, zone)
			},
			checks: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {
				newZone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, zone.ChainId)
				suite.True(found)
				suite.Zero(newZone.GetWithdrawalWaitgroup())
			},
			msg: func(zone *types.Zone) sdk.Msg {
				return &distrtypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: user,
					ValidatorAddress: val,
				}
			},
			triggered: true,
			err:       false,
		},
		{
			name: "valid case trigger redemption rate and without checking if delegatorAddress == performanceAddress",
			setup: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {
				zone.SetWithdrawalWaitgroup(quicksilver.Logger(), 0, "init")
				quicksilver.InterchainstakingKeeper.SetZone(ctx, zone)
			},
			checks: func(ctx sdk.Context, quicksilver *app.Quicksilver, zone *types.Zone) {},
			msg: func(zone *types.Zone) sdk.Msg {
				return &distrtypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: zone.PerformanceAddress.Address,
					ValidatorAddress: val,
				}
			},
			triggered: true,
			err:       false,
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

			test.setup(ctx, quicksilver, &zone)
			prevAllBalancesQueryCnt := 0
			for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
				if query.QueryType == queryAllBalancesPath {
					prevAllBalancesQueryCnt++
				}
			}

			err := quicksilver.InterchainstakingKeeper.HandleWithdrawRewards(ctx, test.msg(&zone), zone.ConnectionId)
			if test.err {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}

			allBalancesQueryCnt := 0
			for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
				if query.QueryType == queryAllBalancesPath {
					allBalancesQueryCnt++
				}
			}
			if test.triggered {
				suite.Equal(prevAllBalancesQueryCnt+1, allBalancesQueryCnt)
			} else {
				suite.Equal(prevAllBalancesQueryCnt, allBalancesQueryCnt)
			}

			test.checks(ctx, quicksilver, &zone)
		})
	}
}

type txAck struct {
	msgs    []sdk.Msg
	memo    string
	success bool
}

func (suite *KeeperTestSuite) TestHandleFailedUndelegate() {
	user := addressutils.GenerateAddressForTestWithPrefix("quick")
	user2 := addressutils.GenerateAddressForTestWithPrefix("quick")
	beneficiary := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	beneficiary2 := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	hash := randomutils.GenerateRandomHashAsHex(32)
	hash2 := randomutils.GenerateRandomHashAsHex(32)

	tests := []struct {
		name     string
		malleate func(ctx sdk.Context, quicksilver *app.Quicksilver)
		txs      func(ctx sdk.Context, quicksilver *app.Quicksilver) []txAck
		expected func(ctx sdk.Context, quicksilver *app.Quicksilver) []types.WithdrawalRecord
		err      bool
	}{
		{
			name: "failed unbond - single wdr, single dist",
			malleate: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)

				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(
					ctx, types.WithdrawalRecord{
						ChainId:   suite.chainB.ChainID,
						Delegator: user,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(312000000)},
						},
						Recipient:      beneficiary,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(312000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(300000000)),
						Txhash:         hash,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       false,
						EpochNumber:    1,
					},
				)
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   1,
					Validator:     vals[0],
					RelatedTxhash: []string{hash},
					Amount:        sdk.NewCoin("uatom", math.NewInt(312000000)),
				})
			},
			txs: func(ctx sdk.Context, quicksilver *app.Quicksilver) []txAck {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)
				return []txAck{
					{
						msgs: []sdk.Msg{
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", math.NewInt(312000000))},
						},
						memo:    types.EpochWithdrawalMemo(1),
						success: false,
					},
				}
			},
			expected: func(ctx sdk.Context, quicksilver *app.Quicksilver) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:        suite.chainB.ChainID,
						Delegator:      user,
						Distribution:   nil,
						Recipient:      beneficiary,
						Amount:         nil,
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(299999999)),
						Txhash:         "0000000000000000000000000000000000000000000000000000000000000000",
						Status:         types.WithdrawStatusQueued,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       true,
						EpochNumber:    1,
					},
				}
			},
			err: false,
		},
		{
			name: "failed unbond - multi related wdr, single dist",
			malleate: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)

				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(
					ctx, types.WithdrawalRecord{
						ChainId:   suite.chainB.ChainID,
						Delegator: user,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(312000000)},
						},
						Recipient:      beneficiary,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(312000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(300000000)),
						Txhash:         hash,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       false,
						EpochNumber:    1,
					},
				)
				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(
					ctx, types.WithdrawalRecord{
						ChainId:   suite.chainB.ChainID,
						Delegator: user2,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(624000000)},
						},
						Recipient:      beneficiary2,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(624000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(600000000)),
						Txhash:         hash2,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       false,
						EpochNumber:    1,
					},
				)
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   1,
					Validator:     vals[0],
					RelatedTxhash: []string{hash, hash2},
					Amount:        sdk.NewCoin("uatom", math.NewInt(936000000)),
				})
			},
			txs: func(ctx sdk.Context, quicksilver *app.Quicksilver) []txAck {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)
				return []txAck{
					{
						msgs: []sdk.Msg{
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", math.NewInt(936000000))},
						},
						memo:    types.EpochWithdrawalMemo(1),
						success: false,
					},
				}
			},
			expected: func(ctx sdk.Context, quicksilver *app.Quicksilver) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:        suite.chainB.ChainID,
						Delegator:      user,
						Distribution:   nil,
						Recipient:      beneficiary,
						Amount:         nil,
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(299999999)),
						Txhash:         "0000000000000000000000000000000000000000000000000000000000000000",
						Status:         types.WithdrawStatusQueued,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       true,
						EpochNumber:    1,
					},
					{
						ChainId:        suite.chainB.ChainID,
						Delegator:      user2,
						Distribution:   nil,
						Recipient:      beneficiary2,
						Amount:         nil,
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(599999999)),
						Txhash:         "0000000000000000000000000000000000000000000000000000000000000001",
						Status:         types.WithdrawStatusQueued,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       true,
						EpochNumber:    1,
					},
				}
			},
			err: false,
		},
		{
			name: "failed unbond - multi related wdr, multi dist",
			malleate: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)

				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(
					ctx, types.WithdrawalRecord{
						ChainId:   suite.chainB.ChainID,
						Delegator: user,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(156000000)},
							{Valoper: vals[1], Amount: math.NewInt(156000000)},
						},
						Recipient:      beneficiary,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(312000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(300000000)),
						Txhash:         hash,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       false,
						EpochNumber:    2,
					},
				)
				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(
					ctx, types.WithdrawalRecord{
						ChainId:   suite.chainB.ChainID,
						Delegator: user2,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(312000000)},
							{Valoper: vals[2], Amount: math.NewInt(312000000)},
						},
						Recipient:      beneficiary2,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(624000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(600000000)),
						Txhash:         hash2,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       false,
						EpochNumber:    2,
					},
				)
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   2,
					Validator:     vals[0],
					RelatedTxhash: []string{hash, hash2},
					Amount:        sdk.NewCoin("uatom", math.NewInt(468000000)),
				})
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   2,
					Validator:     vals[1],
					RelatedTxhash: []string{hash},
					Amount:        sdk.NewCoin("uatom", math.NewInt(156000000)),
				})
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   2,
					Validator:     vals[2],
					RelatedTxhash: []string{hash2},
					Amount:        sdk.NewCoin("uatom", math.NewInt(312000000)),
				})
			},
			txs: func(ctx sdk.Context, quicksilver *app.Quicksilver) []txAck {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)
				return []txAck{
					{
						msgs: []sdk.Msg{
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", math.NewInt(468000000))},
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[1], Amount: sdk.NewCoin("uatom", math.NewInt(156000000))},
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[2], Amount: sdk.NewCoin("uatom", math.NewInt(312000000))},
						},
						memo:    types.EpochWithdrawalMemo(2),
						success: false,
					},
				}
			},
			expected: func(ctx sdk.Context, quicksilver *app.Quicksilver) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{
					{
						ChainId:        suite.chainB.ChainID,
						Delegator:      user,
						Distribution:   nil,
						Recipient:      beneficiary,
						Amount:         nil,
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(300000000)),
						Txhash:         "0000000000000000000000000000000000000000000000000000000000000000",
						Status:         types.WithdrawStatusQueued,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       true,
						EpochNumber:    2,
					},
					{
						ChainId:        suite.chainB.ChainID,
						Delegator:      user2,
						Distribution:   nil,
						Recipient:      beneficiary2,
						Amount:         nil,
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(600000000)),
						Txhash:         "0000000000000000000000000000000000000000000000000000000000000001",
						Status:         types.WithdrawStatusQueued,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       true,
						EpochNumber:    2,
					},
				}
			},
			err: false,
		},
		{
			name: "failed unbond - multi related wdr, multi dist, partial success",
			malleate: func(ctx sdk.Context, quicksilver *app.Quicksilver) {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)

				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(
					ctx, types.WithdrawalRecord{
						ChainId:   suite.chainB.ChainID,
						Delegator: user,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(156000000)},
							{Valoper: vals[1], Amount: math.NewInt(156000000)},
						},
						Recipient:      beneficiary,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(312000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(300000000)),
						Txhash:         hash,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       false,
						EpochNumber:    1,
					},
				)
				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(
					ctx, types.WithdrawalRecord{
						ChainId:   suite.chainB.ChainID,
						Delegator: user2,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(312000000)},
							{Valoper: vals[2], Amount: math.NewInt(312000000)},
						},
						Recipient:      beneficiary2,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(624000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(600000000)),
						Txhash:         hash2,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       false,
						EpochNumber:    1,
					},
				)
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   1,
					Validator:     vals[0],
					RelatedTxhash: []string{hash, hash2},
					Amount:        sdk.NewCoin("uatom", math.NewInt(468000000)),
				})
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   1,
					Validator:     vals[1],
					RelatedTxhash: []string{hash},
					Amount:        sdk.NewCoin("uatom", math.NewInt(156000000)),
				})
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, types.UnbondingRecord{
					ChainId:       suite.chainB.ChainID,
					EpochNumber:   1,
					Validator:     vals[2],
					RelatedTxhash: []string{hash2},
					Amount:        sdk.NewCoin("uatom", math.NewInt(312000000)),
				})
			},
			txs: func(ctx sdk.Context, quicksilver *app.Quicksilver) []txAck {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)
				return []txAck{
					{
						msgs: []sdk.Msg{
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", math.NewInt(468000000))},
						},
						memo:    types.EpochWithdrawalMemo(1),
						success: true,
					},
					{
						msgs: []sdk.Msg{
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[1], Amount: sdk.NewCoin("uatom", math.NewInt(156000000))},
							&stakingtypes.MsgUndelegate{DelegatorAddress: z.DelegationAddress.Address, ValidatorAddress: vals[2], Amount: sdk.NewCoin("uatom", math.NewInt(312000000))},
						},
						memo:    types.EpochWithdrawalMemo(1),
						success: false,
					},
				}
			},
			expected: func(ctx sdk.Context, quicksilver *app.Quicksilver) []types.WithdrawalRecord {
				z, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, z.ChainId)
				suite.True(found)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: user,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(156000000)},
						},
						Recipient:      beneficiary,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(156000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(150000001)),
						Txhash:         hash,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(6 * time.Hour),
						Acknowledged:   true,
						Requeued:       false,
						EpochNumber:    1,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: user2,
						Distribution: []*types.Distribution{
							{Valoper: vals[0], Amount: math.NewInt(312000000)},
						},
						Recipient:      beneficiary2,
						Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(312000000))),
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(300000001)),
						Txhash:         hash2,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(6 * time.Hour),
						Acknowledged:   true,
						Requeued:       false,
						EpochNumber:    1,
					},
					{
						ChainId:        suite.chainB.ChainID,
						Delegator:      user,
						Distribution:   nil,
						Recipient:      beneficiary,
						Amount:         nil,
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(149999999)),
						Txhash:         fmt.Sprintf("%064d", 0),
						Status:         types.WithdrawStatusQueued,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       true,
						EpochNumber:    1,
					},
					{
						ChainId:        suite.chainB.ChainID,
						Delegator:      user2,
						Distribution:   nil,
						Recipient:      beneficiary2,
						Amount:         nil,
						BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(299999999)),
						Txhash:         fmt.Sprintf("%064d", 1),
						Status:         types.WithdrawStatusQueued,
						CompletionTime: time.Time{},
						Acknowledged:   false,
						Requeued:       true,
						EpochNumber:    1,
					},
				}
			},
			err: false,
		},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			test.malleate(ctx, quicksilver)

			txs := test.txs(ctx, quicksilver)
			for _, tx := range txs {
				ctx := suite.chainA.GetContext()
				packet, err := makePacketFromMsgs(quicksilver.AppCodec(), tx.msgs, tx.memo)
				suite.NoError(err)
				ack, err := makeAckForMsgs(ctx, quicksilver.AppCodec(), tx.msgs, tx.success)
				suite.NoError(err)
				bz, err := quicksilver.AppCodec().MarshalJSON(&ack)
				suite.NoError(err)
				err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, bz, suite.chainB.ChainID)
				suite.NoError(err)
			}

			ctx = suite.chainA.GetContext()
			suite.ElementsMatch(quicksilver.InterchainstakingKeeper.AllZoneWithdrawalRecords(ctx, suite.chainB.ChainID), test.expected(ctx, quicksilver))
		})
	}
}

func makeAckForMsgs(ctx sdk.Context, cdc codec.Codec, msgs []sdk.Msg, success bool) (channeltypes.Acknowledgement, error) {
	// If the operation was not successful, return an error acknowledgement
	if !success {
		return channeltypes.NewErrorAcknowledgement(fmt.Errorf("an error happened")), nil
	}

	// Initialize an empty TxMsgData object to hold the responses
	msgData := sdk.TxMsgData{}

	// Iterate through the messages to process their responses
	for _, msg := range msgs {
		// Use UnsafePackAny to avoid the need for a concrete type, which allows for flexibility with msg types
		actualMsg := codectypes.UnsafePackAny(msg)

		// Check if the message type is MsgUndelegate and process accordingly
		if actualMsg.TypeUrl == "/cosmos.staking.v1beta1.MsgUndelegate" {
			// Calculate the completion time for undelegation, set to 6 hours from the current block time
			t := ctx.BlockTime().Add(time.Hour * 6)
			// Create a response for the MsgUndelegate
			resp := stakingtypes.MsgUndelegateResponse{CompletionTime: t}
			// Pack the response into an Any type
			respAny, err := codectypes.NewAnyWithValue(&resp)
			if err != nil {
				return channeltypes.Acknowledgement{}, err
			}
			// Append the response to the MsgData responses
			msgData.MsgResponses = append(msgData.MsgResponses, respAny)
		}
	}

	// Marshal the msgData into bytes for the acknowledgement payload
	bz, err := cdc.Marshal(&msgData)
	if err != nil {
		return channeltypes.Acknowledgement{}, err
	}

	// Return a result acknowledgement with the marshaled msgData
	return channeltypes.NewResultAcknowledgement(bz), nil
}

func makePacketFromMsgs(cdc codec.Codec, msgs []sdk.Msg, memo string) (channeltypes.Packet, error) {
	protoMsgs := make([]proto.Message, len(msgs))
	for i, msg := range msgs {
		protoMsgs[i] = msg.(proto.Message)
	}
	data, err := icatypes.SerializeCosmosTx(cdc, protoMsgs)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: memo,
	}

	return channeltypes.NewPacket(packetData.GetBytes(), 1, "icahost", "channel-1", "icacontroller-%s.delegate", "channel-1", clienttypes.Height{RevisionNumber: 0, RevisionHeight: 1000000}, 0), nil
}

func (suite *KeeperTestSuite) TestHandleFailedUnbondSend() {
	v1 := addressutils.GenerateValAddressForTest().String()
	v2 := addressutils.GenerateValAddressForTest().String()
	user := addressutils.GenerateAddressForTestWithPrefix("quick")
	tests := []struct {
		name    string
		record  func(zone *types.Zone) types.WithdrawalRecord
		message []banktypes.MsgSend
		memo    string
		err     []bool
		check   bool
	}{
		{
			name:    "invalid - unable to parse tx hash",
			message: []banktypes.MsgSend{},
			memo:    "withdrawal/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:     []bool{true},
			check:   false,
		},
		{
			name: "invalid - no matching record",
			record: func(zone *types.Zone) types.WithdrawalRecord {
				return types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: addressutils.GenerateAccAddressForTest().String(),
					Distribution: []*types.Distribution{
						{Valoper: v1, Amount: math.NewInt(1000000)},
						{Valoper: v2, Amount: math.NewInt(1000000)},
					},
					Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
					Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
					BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
					Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
					Status:     types.WithdrawStatusQueued,
				}
			},
			message: []banktypes.MsgSend{},
			memo:    "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:     []bool{true},
			check:   false,
		},
		{
			name: "invalid - try msg send 2 times with one txHash",
			record: func(zone *types.Zone) types.WithdrawalRecord {
				return types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: addressutils.GenerateAccAddressForTest().String(),
					Distribution: []*types.Distribution{
						{Valoper: v1, Amount: math.NewInt(1000000)},
						{Valoper: v2, Amount: math.NewInt(1000000)},
					},
					Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
					Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
					BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2000000)),
					Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
					Status:     types.WithdrawStatusSend,
				}
			},
			message: []banktypes.MsgSend{
				{
					FromAddress: user,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				},
				{
					FromAddress: user,
					Amount:      sdk.NewCoins(sdk.NewCoin(v2+"2", sdk.NewInt(1000000))),
				},
			},
			memo:  "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   []bool{false, true},
			check: true,
		},
		{
			name: "valid",
			record: func(zone *types.Zone) types.WithdrawalRecord {
				return types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: addressutils.GenerateAccAddressForTest().String(),
					Distribution: []*types.Distribution{
						{Valoper: v1, Amount: math.NewInt(1000000)},
						{Valoper: v2, Amount: math.NewInt(1000000)},
					},
					Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
					Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
					BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2000000)),
					Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
					Status:     types.WithdrawStatusSend,
				}
			},
			message: []banktypes.MsgSend{
				{
					FromAddress: user,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				},
			},
			memo:  "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   []bool{false},
			check: true,
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

			var record types.WithdrawalRecord
			if test.record != nil {
				// set up zones
				record = test.record(&zone)
				err := quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				suite.NoError(err)
			}

			// set address for zone mapping
			quicksilver.InterchainstakingKeeper.SetAddressZoneMapping(ctx, user, zone.ChainId)

			// trigger handler
			for i := range test.message {
				err := quicksilver.InterchainstakingKeeper.HandleFailedUnbondSend(ctx, &test.message[i], test.memo)
				if test.err[i] {
					suite.Error(err)
				} else {
					suite.NoError(err)
				}
			}

			if test.check {
				newRecord, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, types.WithdrawStatusUnbond)
				if !found {
					suite.Fail("unable to retrieve new withdrawal record for test")
				}
				suite.Equal(ctx.BlockTime().Add(types.DefaultWithdrawalRequeueDelay), newRecord.CompletionTime)
				suite.Equal(newRecord.Status, types.WithdrawStatusUnbond)
			}
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
	record := types.RedelegationRecord{
		ChainId:     suite.chainB.ChainID,
		EpochNumber: 1,
		Source:      validators[0].ValoperAddress,
		Destination: validators[1].ValoperAddress,
		Amount:      math.NewInt(1000),
	}

	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	redelegate := &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: validators[0].ValoperAddress, ValidatorDstAddress: validators[1].ValoperAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{redelegate})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: types.EpochRebalanceMemo(1),
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	ackBytes := []byte("{\"error\":\"ABCI code: 32: error handling packet on host chain: see events for details\"}")
	// call handler

	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, 1)
	suite.True(found)

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, 1)
	suite.False(found)
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

	val0 := types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2_000_000_000), Status: stakingtypes.BondStatusBonded}
	err := quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val0)
	suite.NoError(err)

	val1 := types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2_000_000_000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val1)
	suite.NoError(err)

	val2 := types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2_000_000_000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val2)
	suite.NoError(err)

	val3 := types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2_000_000_000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val3)
	suite.NoError(err)

	vals = quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)

	delegations := []types.Delegation{
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[0].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1_000_000_000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[1].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1_000_000_000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[2].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1_000_000_000)),
			RedelegationEnd:   0,
		},
		{
			DelegationAddress: zone.DelegationAddress.Address,
			ValidatorAddress:  vals[3].ValoperAddress,
			Amount:            sdk.NewCoin("uatom", sdk.NewInt(1_000_000_000)),
			RedelegationEnd:   0,
		},
	}
	for _, delegation := range delegations {
		quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
		addressBytes, err := addressutils.ValAddressFromBech32(delegation.ValidatorAddress, zone.GetValoperPrefix())
		suite.NoError(err)
		val, _ := quicksilver.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, addressBytes)
		val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
		val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))
	}

	quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 1)
	suite.NoError(err)

	// change intents to trigger redelegations from val[3]
	intents := types.ValidatorIntents{
		{ValoperAddress: vals[0].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[1].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[2].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[3].ValoperAddress, Weight: sdk.NewDecWithPrec(1, 1)},
	}
	zone.AggregateIntent = intents

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 2)
	suite.NoError(err)

	// mock ack for redelegations
	quicksilver.InterchainstakingKeeper.IteratePrefixedRedelegationRecords(ctx, []byte(zone.ChainId), func(idx int64, _ []byte, record types.RedelegationRecord) (stop bool) {
		if record.EpochNumber == 2 {
			msg := stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    zone.DelegationAddress.Address,
				ValidatorSrcAddress: record.Source,
				ValidatorDstAddress: record.Destination,
				Amount:              sdk.NewCoin("uatom", record.Amount),
			}
			err := quicksilver.InterchainstakingKeeper.HandleBeginRedelegate(ctx, &msg, time.Now().Add(time.Hour*24*7), types.EpochRebalanceMemo(2))
			if err != nil {
				return false
			}
		}
		return false
	})

	// check for redelegations
	_, present := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[0].ValoperAddress, 2)
	suite.True(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[1].ValoperAddress, 2)
	suite.True(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[3].ValoperAddress, vals[2].ValoperAddress, 2)
	suite.True(present)

	// change intents to trigger transitive redelegations which should fail rebalance
	zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	intents = types.ValidatorIntents{
		{ValoperAddress: vals[0].ValoperAddress, Weight: sdk.NewDecWithPrec(1, 1)},
		{ValoperAddress: vals[1].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[2].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: vals[3].ValoperAddress, Weight: sdk.NewDecWithPrec(3, 1)},
	}
	zone.AggregateIntent = intents

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 3)
	suite.NoError(err)

	// check for redelegations originating from val[0], they should not be present
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[1].ValoperAddress, 3)
	suite.False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[2].ValoperAddress, 3)
	suite.False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, vals[0].ValoperAddress, vals[3].ValoperAddress, 3)
	suite.False(present)
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

	val0 := types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err := quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val0)
	suite.NoError(err)

	val1 := types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val1)
	suite.NoError(err)

	val2 := types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val2)
	suite.NoError(err)

	val3 := types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("1"), VotingPower: sdk.NewInt(2000), Status: stakingtypes.BondStatusBonded}
	err = quicksilver.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, val3)
	suite.NoError(err)

	delegations := []types.Delegation{
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
		quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
		valAddrBytes, err := addressutils.ValAddressFromBech32(delegation.ValidatorAddress, zone.GetValoperPrefix())
		suite.NoError(err)

		val, found := quicksilver.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		suite.NoError(err)
		suite.True(found)
		val.VotingPower = val.VotingPower.Add(delegation.Amount.Amount)
		val.DelegatorShares = val.DelegatorShares.Add(sdk.NewDecFromInt(delegation.Amount.Amount))

	}

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 1)
	suite.NoError(err)

	quicksilver.InterchainstakingKeeper.IterateAllDelegations(ctx, zone.ChainId, func(delegation types.Delegation) bool {
		if delegation.ValidatorAddress == val0.ValoperAddress {
			delegation.Amount = delegation.Amount.Add(sdk.NewInt64Coin("uatom", 4000))
			quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
		}
		return false
	})

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 2)
	suite.NoError(err)

	// mock ack for redelegations
	quicksilver.InterchainstakingKeeper.IteratePrefixedRedelegationRecords(ctx, []byte(zone.ChainId), func(idx int64, _ []byte, record types.RedelegationRecord) (stop bool) {
		if record.EpochNumber == 2 {
			msg := stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    zone.DelegationAddress.Address,
				ValidatorSrcAddress: record.Source,
				ValidatorDstAddress: record.Destination,
				Amount:              sdk.NewCoin("uatom", record.Amount),
			}
			err := quicksilver.InterchainstakingKeeper.HandleBeginRedelegate(ctx, &msg, time.Now().Add(time.Hour*24*7), types.EpochRebalanceMemo(2))
			if err != nil {
				return false
			}
		}
		return false
	})

	// check for redelegations
	_, present := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val0.ValoperAddress, val1.ValoperAddress, 2)
	suite.False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val0.ValoperAddress, val2.ValoperAddress, 2)
	suite.False(present)

	// change validator delegation to trigger transitive redelegations which should fail rebalance
	quicksilver.InterchainstakingKeeper.IterateAllDelegations(ctx, zone.ChainId, func(delegation types.Delegation) bool {
		if delegation.ValidatorAddress == val0.ValoperAddress {
			delegation.Amount = delegation.Amount.Sub(sdk.NewInt64Coin("uatom", 4000))
			quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
		}
		if delegation.ValidatorAddress == val1.ValoperAddress {
			delegation.Amount = delegation.Amount.Add(sdk.NewInt64Coin("uatom", 4000))
			quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
		}

		return false
	})

	// trigger rebalance
	err = quicksilver.InterchainstakingKeeper.Rebalance(ctx, &zone, 3)
	suite.NoError(err)

	// check for redelegations originating from val[1], they should not be present
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val1.ValoperAddress, val0.ValoperAddress, 3)
	suite.False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val1.ValoperAddress, val1.ValoperAddress, 3)
	suite.False(present)
	_, present = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, val1.ValoperAddress, val3.ValoperAddress, 3)
	suite.False(present)
}

func (suite *KeeperTestSuite) Test_v045Callback() {
	tests := []struct {
		name             string
		setStatements    func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]proto.Message, []byte)
		assertStatements func(ctx sdk.Context, quicksilver *app.Quicksilver) bool
	}{
		{
			name: "msg response with some data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]proto.Message, []byte) {
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				if !found {
					suite.Fail("unable to retrieve zone for test")
				}

				val := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0]

				sendMsg := stakingtypes.MsgDelegate{
					Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(100)),
					DelegatorAddress: zone.PerformanceAddress.Address,
					ValidatorAddress: val,
				}
				response := stakingtypes.MsgDelegateResponse{}

				respBytes := icatypes.ModuleCdc.MustMarshal(&response)
				return []proto.Message{&sendMsg}, respBytes
			},
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				return true
			},
		},
		{
			name: "msg response with nil data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]proto.Message, []byte) {
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
				return []proto.Message{&msgSetWithdrawAddress}, respBytes
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
				Data:         []*sdk.MsgData{{MsgType: "/bob", Data: msgResponseBytes}}, // nolint:staticcheck
				MsgResponses: []*codectypes.Any{},
			}

			ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

			acknowledgement := channeltypes.Acknowledgement{
				Response: &channeltypes.Acknowledgement_Result{
					Result: ackData,
				},
			}

			pdBytes, err := icatypes.SerializeCosmosTx(icatypes.ModuleCdc, msg)
			suite.NoError(err)
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: pdBytes,
				Memo: "test_acknowledgement",
			}

			packetBytes, err := icatypes.ModuleCdc.MarshalJSON(&packetData)
			suite.NoError(err)
			packet := channeltypes.Packet{
				Data: packetBytes,
			}
			ctx = suite.chainA.GetContext()
			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			suite.NoError(quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, icatypes.ModuleCdc.MustMarshalJSON(&acknowledgement), zone.ConnectionId))

			suite.True(test.assertStatements(ctx, quicksilver))
		})
	}
}

func (suite *KeeperTestSuite) Test_v046Callback() {
	tests := []struct {
		name             string
		setStatements    func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]proto.Message, *codectypes.Any)
		assertStatements func(ctx sdk.Context, quicksilver *app.Quicksilver) bool
	}{
		{
			name: "msg response with some data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]proto.Message, *codectypes.Any) {
				zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
				if !found {
					suite.Fail("unable to retrieve zone for test")
				}

				val := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)[0]

				sendMsg := stakingtypes.MsgDelegate{
					Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(100)),
					DelegatorAddress: zone.PerformanceAddress.Address,
					ValidatorAddress: val,
				}
				response := stakingtypes.MsgDelegateResponse{}

				anyResponse, err := codectypes.NewAnyWithValue(&response)
				suite.NoError(err)
				return []proto.Message{&sendMsg}, anyResponse
			},
			assertStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) bool {
				return true
			},
		},
		{
			name: "msg response with nil data",
			setStatements: func(ctx sdk.Context, quicksilver *app.Quicksilver) ([]proto.Message, *codectypes.Any) {
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
				suite.NoError(err)
				return []proto.Message{&msgSetWithdrawAddress}, anyResponse
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
				Data:         []*sdk.MsgData{}, // nolint:staticcheck
				MsgResponses: []*codectypes.Any{anyResp},
			}

			ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

			acknowledgement := channeltypes.Acknowledgement{
				Response: &channeltypes.Acknowledgement_Result{
					Result: ackData,
				},
			}

			pdBytes, err := icatypes.SerializeCosmosTx(icatypes.ModuleCdc, msg)
			suite.NoError(err)
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: pdBytes,
				Memo: "test_acknowledgement",
			}

			packetBytes, err := icatypes.ModuleCdc.MarshalJSON(&packetData)
			suite.NoError(err)
			packet := channeltypes.Packet{
				Data: packetBytes,
			}

			suite.NoError(quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, icatypes.ModuleCdc.MustMarshalJSON(&acknowledgement), "connection-0"))

			suite.True(test.assertStatements(ctx, quicksilver))
		})
	}
}

func (suite *KeeperTestSuite) TestReceiveAckForBeginUndelegate() {
	hash1 := randomutils.GenerateRandomHashAsHex(32)
	hash2 := randomutils.GenerateRandomHashAsHex(32)
	hash3 := randomutils.GenerateRandomHashAsHex(32)
	delegator1 := addressutils.GenerateAddressForTestWithPrefix("quick")
	delegator2 := addressutils.GenerateAddressForTestWithPrefix("quick")
	oneMonth := time.Now().AddDate(0, 1, 0).UTC()
	nilTime := time.Time{}

	tests := []struct {
		name                      string
		epoch                     int64
		withdrawalRecords         func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord
		unbondingRecords          func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.UnbondingRecord
		msgs                      func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []proto.Message
		completionTime            time.Time
		expectedWithdrawalRecords func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord
	}{
		{
			name:  "1 wdr, 2 vals, 1k+1k, 1800 qasset",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:     hash1,
						Status:     types.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.UnbondingRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.UnbondingRecord{
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   1,
						Validator:     vals[0],
						RelatedTxhash: []string{hash1},
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []proto.Message {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []proto.Message{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[0],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			completionTime: oneMonth,
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: oneMonth,
					},
				}
			},
		},
		{
			name:  "1 wdr, 1 vals, 1k, 900 qasset",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     types.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.UnbondingRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.UnbondingRecord{
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   1,
						Validator:     vals[0],
						RelatedTxhash: []string{hash1},
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []proto.Message {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []proto.Message{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[0],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			completionTime: oneMonth,
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: oneMonth,
					},
				}
			},
		},
		{
			name:  "3 wdr, 2 vals, 1k+0.5k, 1350 qasset; 1k+2k, 2700 qasset; 600+400, 900qasset",
			epoch: 2,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(500),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1500))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1350)),
						Txhash:     hash1,
						Status:     types.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(2000),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(2700)),
						Txhash:     hash2,
						Status:     types.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(600),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(400),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash3,
						Status:     types.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.UnbondingRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.UnbondingRecord{
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   2,
						Validator:     vals[1],
						RelatedTxhash: []string{hash1, hash2, hash3},
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []proto.Message {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []proto.Message{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[1],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2900)),
					},
				}
			},
			completionTime: nilTime,
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(500),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1500))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1350)),
						Txhash:     hash1,
						Status:     types.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(2000),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(2700)),
						Txhash:     hash2,
						Status:     types.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(600),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(400),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash3,
						Status:     types.WithdrawStatusUnbond,
					},
				}
			},
		},
		{
			name:  "2 wdr, 1 vals, 1k; 2 vals; 123 + 456 ",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:     hash1,
						Status:     types.WithdrawStatusUnbond,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[1],
								Amount:  math.NewInt(123),
							},
							{
								Valoper: vals[2],
								Amount:  math.NewInt(456),
							},
						},
						Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(579))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, sdk.NewInt(521)),
						Txhash:     hash2,
						Status:     types.WithdrawStatusUnbond,
					},
				}
			},
			unbondingRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.UnbondingRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.UnbondingRecord{
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   1,
						Validator:     vals[0],
						RelatedTxhash: []string{hash1},
					},
					{
						ChainId:       suite.chainB.ChainID,
						EpochNumber:   1,
						Validator:     vals[1],
						RelatedTxhash: []string{hash2},
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []proto.Message {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []proto.Message{
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[0],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
					&stakingtypes.MsgUndelegate{
						DelegatorAddress: zone.DelegationAddress.Address,
						ValidatorAddress: vals[1],
						Amount:           sdk.NewCoin(zone.BaseDenom, sdk.NewInt(123)),
					},
				}
			},
			completionTime: nilTime,
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{}
			},
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

			for _, wdr := range test.withdrawalRecords(ctx, quicksilver, zone) {
				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, wdr)
			}

			for _, ubr := range test.unbondingRecords(ctx, quicksilver, zone) {
				quicksilver.InterchainstakingKeeper.SetUnbondingRecord(ctx, ubr)
			}

			msgs := test.msgs(ctx, quicksilver, zone)
			data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), msgs)
			suite.NoError(err)

			// validate memo < 256 bytes
			packetData := icatypes.InterchainAccountPacketData{
				Type: icatypes.EXECUTE_TX,
				Data: data,
				Memo: types.EpochWithdrawalMemo(test.epoch),
			}

			packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

			responses := make([]*codectypes.Any, 0)

			for range msgs {
				response := stakingtypes.MsgUndelegateResponse{
					CompletionTime: test.completionTime,
				}

				anyResponse, err := codectypes.NewAnyWithValue(&response)
				suite.NoError(err)
				responses = append(responses, anyResponse)
			}

			txMsgData := &sdk.TxMsgData{
				MsgResponses: responses,
			}

			ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

			acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
			ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
			suite.NoError(err)

			// call handler

			for _, ubr := range test.unbondingRecords(ctx, quicksilver, zone) {
				_, found = quicksilver.InterchainstakingKeeper.GetUnbondingRecord(ctx, zone.ChainId, ubr.Validator, test.epoch)
				suite.True(found)
			}

			err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
			suite.NoError(err)

			for idx, ewdr := range test.expectedWithdrawalRecords(ctx, quicksilver, zone) {
				wdr, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, ewdr.Txhash, ewdr.Status)
				suite.True(found)
				suite.Equal(ewdr.Amount, wdr.Amount)
				suite.Equal(ewdr.BurnAmount, wdr.BurnAmount)
				suite.Equal(ewdr.Delegator, wdr.Delegator)
				suite.Equal(ewdr.Distribution, wdr.Distribution, idx)
				suite.Equal(ewdr.Status, wdr.Status)
				suite.Equal(ewdr.CompletionTime, wdr.CompletionTime)
				suite.True(wdr.Acknowledged)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestReceiveAckForBeginRedelegateNonNilCompletion() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	complete := time.Now().UTC().AddDate(0, 0, 21)

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	validators := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	// create redelegation record
	record := types.RedelegationRecord{
		ChainId:     suite.chainB.ChainID,
		EpochNumber: 1,
		Source:      validators[0].ValoperAddress,
		Destination: validators[1].ValoperAddress,
		Amount:      math.NewInt(1000),
	}

	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	beforeSource := types.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  validators[0].ValoperAddress,
		Amount:            sdk.NewCoin(zone.BaseDenom, math.NewInt(2000)),
	}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, beforeSource)

	redelegate := &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: validators[0].ValoperAddress, ValidatorDstAddress: validators[1].ValoperAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{redelegate})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: types.EpochRebalanceMemo(1),
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := stakingtypes.MsgUndelegateResponse{
		CompletionTime: complete,
	}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	// call handler

	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, 1)
	suite.True(found)

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	afterRedelegation, found := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, 1)
	suite.True(found)
	suite.Equal(complete, afterRedelegation.CompletionTime)

	afterSource, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, validators[1].ValoperAddress)
	suite.True(found)
	suite.Equal(beforeSource.Amount.Sub(redelegate.Amount), afterSource.Amount)

	afterTarget, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, validators[1].ValoperAddress)
	suite.True(found)
	suite.Equal(complete.Unix(), afterTarget.RedelegationEnd)
	// target did not exist before redelegation
	suite.Equal(redelegate.Amount, afterTarget.Amount)
}

func (suite *KeeperTestSuite) TestReceiveAckForBeginRedelegateNilCompletion() {
	suite.SetupTest()
	suite.setupTestZones()

	epoch := int64(2)

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	complete := time.Time{}

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	validators := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	// create redelegation record
	record := types.RedelegationRecord{
		ChainId:     suite.chainB.ChainID,
		EpochNumber: epoch,
		Source:      validators[0].ValoperAddress,
		Destination: validators[1].ValoperAddress,
		Amount:      sdk.NewInt(1000),
	}

	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	beforeTarget := types.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  validators[1].ValoperAddress,
		Amount:            sdk.NewCoin(zone.BaseDenom, math.NewInt(2000)),
	}

	beforeSource := types.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  validators[0].ValoperAddress,
		Amount:            sdk.NewCoin(zone.BaseDenom, math.NewInt(1001)),
	}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, beforeTarget)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, beforeSource)

	redelegate := &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: validators[0].ValoperAddress, ValidatorDstAddress: validators[1].ValoperAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{redelegate})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: types.EpochRebalanceMemo(epoch),
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := stakingtypes.MsgUndelegateResponse{
		CompletionTime: complete,
	}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	// call handler

	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, epoch)
	suite.True(found)

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	_, found = quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, epoch)
	suite.False(found) // redelegation record should have been removed.

	afterSource, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, validators[0].ValoperAddress)
	suite.True(found)
	suite.Equal(beforeSource.Amount.Sub(redelegate.Amount), afterSource.Amount)

	afterTarget, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, validators[1].ValoperAddress)
	suite.True(found)
	suite.Equal(complete.Unix(), afterTarget.RedelegationEnd)
	suite.Equal(beforeTarget.Amount.Add(redelegate.Amount), afterTarget.Amount)
}

func (suite *KeeperTestSuite) TestReceiveAckForBeginRedelegateNoExistingRecord() {
	suite.SetupTest()
	suite.setupTestZones()

	epoch := int64(2)

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	complete := ctx.BlockTime().Add(time.Hour * 72)

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	validators := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)

	beforeTarget := types.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  validators[1].ValoperAddress,
		Amount:            sdk.NewCoin(zone.BaseDenom, math.NewInt(2000)),
	}

	beforeSource := types.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  validators[0].ValoperAddress,
		Amount:            sdk.NewCoin(zone.BaseDenom, math.NewInt(1001)),
	}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, beforeTarget)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, beforeSource)

	redelegate := &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: validators[0].ValoperAddress, ValidatorDstAddress: validators[1].ValoperAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{redelegate})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: types.EpochRebalanceMemo(epoch),
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := stakingtypes.MsgUndelegateResponse{
		CompletionTime: complete,
	}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	// call handler

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	createdRecord, found := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, zone.ChainId, validators[0].ValoperAddress, validators[1].ValoperAddress, epoch)
	suite.True(found) // redelegation record should have been removed.

	suite.Equal(redelegate.Amount.Amount, createdRecord.Amount)
	suite.Equal(redelegate.ValidatorDstAddress, createdRecord.Destination)
	suite.Equal(redelegate.ValidatorSrcAddress, createdRecord.Source)
	suite.Equal(epoch, createdRecord.EpochNumber)
	suite.Equal(complete, createdRecord.CompletionTime)

	afterSource, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, validators[0].ValoperAddress)
	suite.True(found)
	suite.Equal(beforeSource.Amount.Sub(redelegate.Amount), afterSource.Amount)

	afterTarget, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, validators[1].ValoperAddress)
	suite.True(found)
	suite.Equal(complete.Unix(), afterTarget.RedelegationEnd)
	suite.Equal(beforeTarget.Amount.Add(redelegate.Amount), afterTarget.Amount)
}

func (suite *KeeperTestSuite) TestReceiveAckForWithdrawReward() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	val := addressutils.GenerateValAddressForTest().String()
	user := addressutils.GenerateAddressForTestWithPrefix("quick")

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	zone.SetWithdrawalWaitgroup(quicksilver.Logger(), 1, "init")
	quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

	withdrawReward := &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: user,
		ValidatorAddress: val,
	}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{withdrawReward})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := distrtypes.MsgWithdrawDelegatorRewardResponse{}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	// Success case
	prevAllBalancesQueryCnt := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.QueryType == queryAllBalancesPath {
			prevAllBalancesQueryCnt++
		}
	}

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	allBalancesQueryCnt := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.QueryType == queryAllBalancesPath {
			allBalancesQueryCnt++
		}
	}

	suite.Equal(prevAllBalancesQueryCnt+1, allBalancesQueryCnt)
}

func (suite *KeeperTestSuite) TestReceiveAckForRedeemTokens() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}

	vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
	delegationRecord := types.Delegation{
		DelegationAddress: zone.DelegationAddress.Address,
		ValidatorAddress:  vals[0],
		Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
	}
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegationRecord)
	txHash := randomutils.GenerateRandomHashAsHex(32)
	t := ctx.BlockTime().Add(-time.Minute * 5)
	quicksilver.InterchainstakingKeeper.SetReceipt(ctx, types.Receipt{
		ChainId:   zone.ChainId,
		Sender:    addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
		Txhash:    txHash,
		Amount:    sdk.NewCoins(sdk.NewCoin(vals[0]+"/1", sdk.NewInt(100))),
		FirstSeen: &t,
	})

	redeemTokens := &lsmstakingtypes.MsgRedeemTokensForShares{
		DelegatorAddress: zone.DelegationAddress.Address,
		Amount:           sdk.NewCoin(vals[0]+"/1", sdk.NewInt(100)),
	}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{redeemTokens})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: txHash,
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := lsmstakingtypes.MsgRedeemTokensForSharesResponse{
		Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(100)),
	}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)

	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("connectionID"), suite.path.EndpointA.ConnectionID))
	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	delegationRecord, found = quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, vals[0])
	suite.True(found)
	suite.Equal(delegationRecord.Amount, sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1100)))
}

func (suite *KeeperTestSuite) TestReceiveAckForTokenizedShares() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}

	vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
	user := addressutils.GenerateAddressForTestWithPrefix("quick")
	txHash := randomutils.GenerateRandomHashAsHex(32)

	withdrawalRecord := types.WithdrawalRecord{
		ChainId:   suite.chainB.ChainID,
		Delegator: user,
		Distribution: []*types.Distribution{
			{
				Valoper: vals[0],
				Amount:  math.NewInt(1000),
			},
		},
		Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
		Amount:         sdk.Coins{},
		BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(800)),
		Txhash:         txHash,
		Status:         types.WithdrawStatusTokenize,
		CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
		Acknowledged:   false,
	}
	_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, withdrawalRecord)
	_, found = quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, txHash, types.WithdrawStatusTokenize)
	suite.True(found)

	tokenizeShares := &lsmstakingtypes.MsgTokenizeShares{
		DelegatorAddress:    zone.DelegationAddress.Address,
		ValidatorAddress:    vals[0],
		Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
		TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
	}
	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{tokenizeShares})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: txHash,
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := lsmstakingtypes.MsgTokenizeSharesResponse{
		Amount: sdk.NewCoin(vals[0]+"/1", sdk.NewInt(1000)),
	}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)
	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	_, found = quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, txHash, types.WithdrawStatusTokenize)
	suite.False(found)

	wr, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, txHash, types.WithdrawStatusSend)
	suite.True(found)

	suite.Equal(wr.Amount[0], response.Amount)
}

func (suite *KeeperTestSuite) TestReceiveAckForDelegate() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}

	vals := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
	user := addressutils.GenerateAddressForTestWithPrefix("quick")
	txHash := randomutils.GenerateRandomHashAsHex(32)

	firstSeen := ctx.BlockTime().Add(-10 * time.Hour)
	completed := ctx.BlockTime().Add(-1 * time.Hour)
	receipt := types.Receipt{
		ChainId:   zone.ChainId,
		Sender:    user,
		Txhash:    txHash,
		Amount:    sdk.Coins{sdk.NewCoin("uatom", sdk.NewInt(1000))},
		FirstSeen: &firstSeen,
		Completed: &completed,
	}
	quicksilver.InterchainstakingKeeper.SetReceipt(ctx, receipt)

	withdrawReward := &stakingtypes.MsgDelegate{
		DelegatorAddress: zone.DelegationAddress.Address,
		ValidatorAddress: vals[0],
		Amount:           sdk.NewCoin("uatom", sdk.NewInt(1000)),
	}

	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{withdrawReward})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: txHash,
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := stakingtypes.MsgDelegateResponse{}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)
	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	newCompleted := ctx.BlockTime()
	newReceipt, found := quicksilver.InterchainstakingKeeper.GetReceipt(ctx, zone.ChainId, txHash)
	suite.True(found)
	suite.Equal(newReceipt.ChainId, receipt.ChainId)
	suite.Equal(newReceipt.Sender, receipt.Sender)
	suite.Equal(newReceipt.Amount, receipt.Amount)
	suite.Equal(newReceipt.FirstSeen, receipt.FirstSeen)
	suite.Equal(newReceipt.Completed, &newCompleted)

	delegation, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, zone.DelegationAddress.Address, vals[0])
	suite.True(found)
	suite.Equal(delegation.Amount, sdk.NewCoin("uatom", sdk.NewInt(1000)))
}

func (suite *KeeperTestSuite) TestReceiveAckForBankSend() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}

	ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("connectionID"), suite.path.EndpointA.ConnectionID))
	quicksilver.InterchainstakingKeeper.IBCKeeper.ChannelKeeper.SetChannel(ctx, "transfer", "channel-0", TestChannel)

	withdrawReward := &banktypes.MsgSend{
		FromAddress: zone.DepositAddress.Address,
		ToAddress:   zone.DelegationAddress.Address,
		Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1_000_000))),
	}

	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{withdrawReward})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	response := banktypes.MsgSendResponse{}

	anyResponse, err := codectypes.NewAnyWithValue(&response)
	suite.NoError(err)

	txMsgData := &sdk.TxMsgData{
		MsgResponses: []*codectypes.Any{anyResponse},
	}

	ackData := icatypes.ModuleCdc.MustMarshal(txMsgData)
	acknowledgement := channeltypes.NewResultAcknowledgement(ackData)
	ackBytes, err := icatypes.ModuleCdc.MarshalJSON(&acknowledgement)
	suite.NoError(err)

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestReceiveAckErrForBankSend() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}

	v1 := addressutils.GenerateValAddressForTest().String()
	v2 := addressutils.GenerateValAddressForTest().String()
	user := addressutils.GenerateAddressForTestWithPrefix("quick")

	withdrawalRecord := types.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: addressutils.GenerateAccAddressForTest().String(),
		Distribution: []*types.Distribution{
			{Valoper: v1, Amount: math.NewInt(1000000)},
			{Valoper: v2, Amount: math.NewInt(1000000)},
		},
		Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
		Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
		BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2000000)),
		Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
		Status:     types.WithdrawStatusSend,
	}
	_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, withdrawalRecord)
	quicksilver.InterchainstakingKeeper.SetAddressZoneMapping(ctx, user, zone.ChainId)

	send := &banktypes.MsgSend{
		FromAddress: zone.DelegationAddress.Address,
		Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
	}

	data, err := icatypes.SerializeCosmosTx(quicksilver.InterchainstakingKeeper.GetCodec(), []proto.Message{send})
	suite.NoError(err)

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
	}

	packet := channeltypes.Packet{Data: quicksilver.InterchainstakingKeeper.GetCodec().MustMarshalJSON(&packetData)}

	ackBytes := []byte("{\"error\":\"ABCI code: 32: error handling packet on host chain: see events for details\"}")

	err = quicksilver.InterchainstakingKeeper.HandleAcknowledgement(ctx, packet, ackBytes, zone.ConnectionId)
	suite.NoError(err)

	newRecord, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D", types.WithdrawStatusUnbond)
	if !found {
		suite.Fail("unable to retrieve new withdrawal record for test")
	}
	suite.Equal(ctx.BlockTime().Add(types.DefaultWithdrawalRequeueDelay), newRecord.CompletionTime)
}

func (suite *KeeperTestSuite) TestHandleMaturedUbondings() {
	hash1 := randomutils.GenerateRandomHashAsHex(32)
	hash2 := randomutils.GenerateRandomHashAsHex(32)
	delegator1 := addressutils.GenerateAddressForTestWithPrefix("quick")
	delegator2 := addressutils.GenerateAddressForTestWithPrefix("quick")

	tests := []struct {
		name                      string
		epoch                     int64
		withdrawalRecords         func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord
		completionTime            time.Time
		expectedWithdrawalRecords func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord
	}{
		{
			name:  "1 wdr, valid completion time, acknowledged ",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusSend,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
		},
		{
			name:  "1 wdr, invalid completion time, acknowledged ",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
		},
		{
			name:  "1 wdr, invalid completion time, unacknowledged ",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(1 * time.Hour),
						Acknowledged:   false,
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(1 * time.Hour),
						Acknowledged:   false,
					},
				}
			},
		},

		{
			name:  "valid completion time, Unacknowledged ",
			epoch: 1,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   false,
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(900)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   false,
					},
				}
			},
		},
		{
			name:  "valid completion time, 1 acknowledged and 1 unacknowledged ",
			epoch: 2,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(500),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1500))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1350)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   false,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(2000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(2700)),
						Txhash:         hash2,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator1,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(500),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1500))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1350)),
						Txhash:         hash1,
						Status:         types.WithdrawStatusUnbond,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   false,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator2,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(2000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000))),
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(2700)),
						Txhash:         hash2,
						Status:         types.WithdrawStatusSend,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
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

			for _, wdr := range test.withdrawalRecords(ctx, quicksilver, zone) {
				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, wdr)
			}

			err := quicksilver.InterchainstakingKeeper.HandleMaturedWithdrawals(ctx, &zone)
			suite.NoError(err)

			for idx, ewdr := range test.expectedWithdrawalRecords(ctx, quicksilver, zone) {
				wdr, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, ewdr.Txhash, ewdr.Status)
				suite.True(found)
				suite.Equal(ewdr.Amount, wdr.Amount)
				suite.Equal(ewdr.BurnAmount, wdr.BurnAmount)
				suite.Equal(ewdr.Delegator, wdr.Delegator)
				suite.Equal(ewdr.Distribution, wdr.Distribution, idx)
				suite.Equal(ewdr.Status, wdr.Status)
				suite.Equal(ewdr.CompletionTime, wdr.CompletionTime)
				suite.Equal(ewdr.Acknowledged, wdr.Acknowledged)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleTokenizedShares() {
	txHash := randomutils.GenerateRandomHashAsHex(32)
	txHash1 := randomutils.GenerateRandomHashAsHex(32)
	delegator := addressutils.GenerateAddressForTestWithPrefix("quick")

	tests := []struct {
		name                      string
		epoch                     int64
		txHash                    string
		withdrawalRecords         func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord
		msgs                      func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg
		sharesAmount              func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coins
		expectedWithdrawalRecords func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord
	}{
		{
			name:   "1 wdr, 2 distributions, 2 msgs, withdraw success",
			epoch:  1,
			txHash: txHash,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.Coins{},
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[0],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(500)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[1],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(500)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
				}
			},
			sharesAmount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coins {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return sdk.NewCoins(
					sdk.NewCoin(vals[0]+"/1", sdk.NewInt(1000)),
					sdk.NewCoin(vals[1]+"/2", sdk.NewInt(1000)),
				)
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				return []types.WithdrawalRecord{}
			},
		},
		{
			name:   "1 wdr, 2 distributions, 1 msgs, withdraw half",
			epoch:  1,
			txHash: txHash,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.Coins{},
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[0],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
				}
			},
			sharesAmount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coins {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return sdk.NewCoins(
					sdk.NewCoin(vals[0]+"/1", sdk.NewInt(1000)),
				)
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.Coins{sdk.NewCoin(vals[0]+"/1", sdk.NewInt(1000))},
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
		},
		{
			name:   "1 wdr, 2 distributions, 1 msgs, not match amount",
			epoch:  1,
			txHash: txHash,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.Coins{},
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[0],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(500)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
				}
			},
			sharesAmount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coins {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return sdk.NewCoins(
					sdk.NewCoin(vals[0]+"/1", sdk.NewInt(500)),
				)
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         nil,
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
		},
		{
			name:   "1 wdr, 2 distributions, 2 msgs, not match denom",
			epoch:  1,
			txHash: txHash,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.Coins{},
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[0],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(500)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[0],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(500)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
				}
			},
			sharesAmount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coins {
				return sdk.NewCoins(
					sdk.NewCoin("not_match_denom_0x", sdk.NewInt(1000)),
					sdk.NewCoin("not_match_denom_1x", sdk.NewInt(1000)),
				)
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         nil,
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
		},
		{
			name:   "2 wdr, 2 distributions, 2 msgs, not match denom",
			epoch:  1,
			txHash: txHash,
			withdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.Coins{},
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         sdk.Coins{},
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash1,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[0],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
					&lsmstakingtypes.MsgTokenizeShares{
						DelegatorAddress:    zone.DelegationAddress.Address,
						ValidatorAddress:    vals[0],
						Amount:              sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
						TokenizedShareOwner: addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
					},
				}
			},
			sharesAmount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coins {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return sdk.NewCoins(
					sdk.NewCoin(vals[0]+"/1", sdk.NewInt(1000)),
					sdk.NewCoin(vals[1]+"/2", sdk.NewInt(1000)),
				)
			},
			expectedWithdrawalRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)

				return []types.WithdrawalRecord{
					{
						ChainId:   suite.chainB.ChainID,
						Delegator: delegator,
						Distribution: []*types.Distribution{
							{
								Valoper: vals[0],
								Amount:  math.NewInt(1000),
							},
							{
								Valoper: vals[1],
								Amount:  math.NewInt(1000),
							},
						},
						Recipient:      addressutils.GenerateAddressForTestWithPrefix(zone.GetAccountPrefix()),
						Amount:         nil,
						BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1800)),
						Txhash:         txHash1,
						Status:         types.WithdrawStatusTokenize,
						CompletionTime: ctx.BlockTime().Add(-1 * time.Hour),
						Acknowledged:   true,
					},
				}
			},
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

			shareAmount := test.sharesAmount(ctx, quicksilver, zone)
			wdrs := test.withdrawalRecords(ctx, quicksilver, zone)
			for _, wdr := range wdrs {
				_ = quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, wdr)
			}

			for index, msg := range test.msgs(ctx, quicksilver, zone) {
				err := quicksilver.InterchainstakingKeeper.HandleTokenizedShares(ctx, msg, shareAmount[index], test.txHash)
				suite.NoError(err)
			}

			ewdrs := test.expectedWithdrawalRecords(ctx, quicksilver, zone)

			if len(ewdrs) == 0 {
				_, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, test.txHash, types.WithdrawStatusTokenize)
				suite.False(found)
			}

			for idx, ewdr := range ewdrs {
				wdr, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, ewdr.Txhash, ewdr.Status)
				suite.True(found)
				suite.Equal(ewdr.Amount, wdr.Amount)
				suite.Equal(ewdr.BurnAmount, wdr.BurnAmount)
				suite.Equal(ewdr.Delegator, wdr.Delegator)
				suite.Equal(ewdr.Distribution, wdr.Distribution, idx)
				suite.Equal(ewdr.Status, wdr.Status)
				suite.Equal(ewdr.CompletionTime, wdr.CompletionTime)
				suite.Equal(ewdr.Acknowledged, wdr.Acknowledged)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestTriggerRedemptionRate() {
	suite.Run("trigger redemption rate", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)

		prevAllBalancesQueryCnt := 0
		for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
			if query.QueryType == queryAllBalancesPath {
				prevAllBalancesQueryCnt++
			}
		}

		err := quicksilver.InterchainstakingKeeper.TriggerRedemptionRate(ctx, &zone)
		suite.NoError(err)

		allBalancesQueryCnt := 0
		for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
			if query.QueryType == queryAllBalancesPath {
				allBalancesQueryCnt++
			}
		}

		suite.Equal(prevAllBalancesQueryCnt+1, allBalancesQueryCnt)
	})
}

func (suite *KeeperTestSuite) TestGetValidatorForToken() {
	tests := []struct {
		name            string
		err             bool
		setupConnection bool
		amount          func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coin
		expectVal       func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) string
	}{
		{
			name:            "Found validator",
			err:             false,
			setupConnection: true,
			amount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coin {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return sdk.NewCoin(vals[0]+"/1", sdk.NewInt(100))
			},
			expectVal: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) string {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return vals[0]
			},
		},
		{
			name:            "Not found validator",
			err:             true,
			setupConnection: true,
			amount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coin {
				return sdk.NewCoin("hello", sdk.NewInt(100))
			},
			expectVal: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) string {
				return ""
			},
		},
		{
			name:            "Not setup connection",
			err:             true,
			setupConnection: false,
			amount: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) sdk.Coin {
				return sdk.NewCoin("hello", sdk.NewInt(100))
			},
			expectVal: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) string {
				return ""
			},
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
			amount := test.amount(ctx, quicksilver, zone)
			resVal, err := quicksilver.InterchainstakingKeeper.GetValidatorForToken(ctx, amount, suite.path.EndpointA.ConnectionID)

			if test.err {
				suite.Error(err)
			} else {
				suite.NoError(err)
				expVal := test.expectVal(ctx, quicksilver, zone)
				suite.Equal(expVal, resVal)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleCompleteSend() {
	testCases := []struct {
		name          string
		message       func(zone *types.Zone) sdk.Msg
		memo          string
		expectedError error
	}{
		{
			name: "unexpected completed send",
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: "",
					ToAddress:   "",
					Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1_000_000))),
				}
			},
			expectedError: errors.New("unexpected completed send (2) from  to  (amount: 1000000uatom)"),
		},
		{
			name: "From WithdrawalAddress",
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.WithdrawalAddress.Address,
					ToAddress:   zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1_000_000))),
				}
			},
			expectedError: nil,
		},
		{
			name: "From DepositAddress to DelegateAddress",
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DepositAddress.Address,
					ToAddress:   zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1_000_000))),
				}
			},
			memo:          "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			expectedError: nil,
		},
		{
			name: "From DepositAddress",
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DelegationAddress.Address,
					ToAddress:   "",
					Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1_000_000))),
				}
			},
			memo:          "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			expectedError: errors.New("no matching withdrawal record found"),
		},
		{
			name: "From DepositAddress to Withdrawal Address",
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DepositAddress.Address,
					ToAddress:   zone.WithdrawalAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1_000_000))),
				}
			},
			memo:          "",
			expectedError: nil,
		},
		{
			name: "From DelegateAddress to Withdrawal Address",
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DelegationAddress.Address,
					ToAddress:   zone.WithdrawalAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1_000_000))),
				}
			},
			memo:          "",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("connectionID"), suite.path.EndpointA.ConnectionID))
			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			if !found {
				suite.Fail("unable to retrieve zone for test")
			}
			quicksilver.InterchainstakingKeeper.IBCKeeper.ChannelKeeper.SetChannel(ctx, "transfer", "channel-0", TestChannel)

			msg := tc.message(&zone)

			err := quicksilver.InterchainstakingKeeper.HandleCompleteSend(ctx, msg, tc.memo, zone.ConnectionId)
			if tc.expectedError != nil {
				suite.Equal(tc.expectedError, err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleFailedBankSend() {
	v1 := addressutils.GenerateValAddressForTest().String()
	v2 := addressutils.GenerateValAddressForTest().String()
	user := addressutils.GenerateAddressForTestWithPrefix("quick")
	tests := []struct {
		name            string
		record          func(zone *types.Zone) types.WithdrawalRecord
		setupConnection bool
		message         func(zone *types.Zone) sdk.Msg
		memo            string
		err             bool
		check           bool
	}{
		{
			name:            "invalid - can not cast to MsgSend",
			setupConnection: false,
			message: func(zone *types.Zone) sdk.Msg {
				return &types.MsgRequestRedemption{}
			},
			memo:  "withdrawal/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   true,
			check: false,
		},
		{
			name:            "invalid - not has connection",
			setupConnection: false,
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				}
			},
			memo:  "withdrawal/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   true,
			check: false,
		},
		{
			name:            "Send from DelegateAddress then HandleFailedUnbondSend, invalid - unable to parse tx hash",
			setupConnection: true,
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				}
			},
			memo:  "withdrawal/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   true,
			check: false,
		},
		{
			name:            "Send from DelegateAddress then HandleFailedUnbondSend, invalid - no matching record",
			setupConnection: true,
			record: func(zone *types.Zone) types.WithdrawalRecord {
				return types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: addressutils.GenerateAccAddressForTest().String(),
					Distribution: []*types.Distribution{
						{Valoper: v1, Amount: math.NewInt(1000000)},
						{Valoper: v2, Amount: math.NewInt(1000000)},
					},
					Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
					Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
					BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
					Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
					Status:     types.WithdrawStatusQueued,
				}
			},
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				}
			},
			memo:  "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   true,
			check: false,
		},
		{
			name:            "Send from DelegateAddress then HandleFailedUnbondSend, invalid - try msg send 2 times with one txHash",
			setupConnection: true,
			record: func(zone *types.Zone) types.WithdrawalRecord {
				return types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: addressutils.GenerateAccAddressForTest().String(),
					Distribution: []*types.Distribution{
						{Valoper: v1, Amount: math.NewInt(1000000)},
						{Valoper: v2, Amount: math.NewInt(1000000)},
					},
					Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
					Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
					BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2000000)),
					Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
					Status:     types.WithdrawStatusSend,
				}
			},
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				}
			},
			memo:  "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   false,
			check: true,
		},
		{
			name:            "Send from DelegateAddress then HandleFailedUnbondSend, valid",
			setupConnection: true,
			record: func(zone *types.Zone) types.WithdrawalRecord {
				return types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: addressutils.GenerateAccAddressForTest().String(),
					Distribution: []*types.Distribution{
						{Valoper: v1, Amount: math.NewInt(1000000)},
						{Valoper: v2, Amount: math.NewInt(1000000)},
					},
					Recipient:  addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix),
					Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000))),
					BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(2000000)),
					Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
					Status:     types.WithdrawStatusSend,
				}
			},
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				}
			},
			memo:  "unbondSend/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   false,
			check: true,
		},
		{
			name:            "Send from WithdrawlAddress, valid - but nothing change",
			setupConnection: true,
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.WithdrawalAddress.WithdrawalAddress,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				}
			},
			memo:  "withdrawal/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   false,
			check: false,
		},
		{
			name:            "Send from DepositAddress to DelegationAddress, valid - but nothing change",
			setupConnection: true,
			message: func(zone *types.Zone) sdk.Msg {
				return &banktypes.MsgSend{
					FromAddress: zone.DepositAddress.Address,
					ToAddress:   zone.DelegationAddress.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(v1+"1", sdk.NewInt(1000000))),
				}
			},
			memo:  "withdrawal/7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
			err:   false,
			check: false,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			if test.setupConnection {
				ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("connectionID"), suite.path.EndpointA.ConnectionID))
			}
			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			if !found {
				suite.Fail("unable to retrieve zone for test")
			}

			var record types.WithdrawalRecord
			if test.record != nil {
				record = test.record(&zone)
				err := quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
				suite.NoError(err)
			}

			// set address for zone mapping
			quicksilver.InterchainstakingKeeper.SetAddressZoneMapping(ctx, user, zone.ChainId)
			msg := test.message(&zone)
			err := quicksilver.InterchainstakingKeeper.HandleFailedBankSend(ctx, msg, test.memo, zone.ConnectionId)

			if test.err {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}

			if test.check {
				newRecord, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, types.WithdrawStatusUnbond)
				if !found {
					suite.Fail("unable to retrieve new withdrawal record for test")
				}
				suite.Equal(ctx.BlockTime().Add(types.DefaultWithdrawalRequeueDelay), newRecord.CompletionTime)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleRedeemTokens() {
	tests := []struct {
		name                      string
		errs                      []bool
		msgs                      func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg
		delegationRecords         func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation
		expectedDelegationRecords func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation
	}{
		{
			name: "1 record, 1 msgs, redeem success",
			errs: []bool{false},
			delegationRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgRedeemTokensForShares{
						DelegatorAddress: zone.DelegationAddress.Address,
						Amount:           sdk.NewCoin(vals[0]+"/1", sdk.NewInt(200)),
					},
				}
			},
			expectedDelegationRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1200)),
					},
				}
			},
		},
		{
			name: "2 record, 2 msgs, redeem success",
			errs: []bool{false, false},
			delegationRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgRedeemTokensForShares{
						DelegatorAddress: zone.DelegationAddress.Address,
						Amount:           sdk.NewCoin(vals[0]+"/1", sdk.NewInt(100)),
					},
					&lsmstakingtypes.MsgRedeemTokensForShares{
						DelegatorAddress: zone.DelegationAddress.Address,
						Amount:           sdk.NewCoin(vals[1]+"/2", sdk.NewInt(200)),
					},
				}
			},
			expectedDelegationRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1100)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1200)),
					},
				}
			},
		},
		{
			name: "2 record, 2 msgs, redeem failed first msg",
			errs: []bool{true, false},
			delegationRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
				}
			},
			msgs: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []sdk.Msg {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []sdk.Msg{
					&lsmstakingtypes.MsgRedeemTokensForShares{
						DelegatorAddress: zone.DelegationAddress.Address,
						Amount:           sdk.NewCoin("hello", sdk.NewInt(100)),
					},
					&lsmstakingtypes.MsgRedeemTokensForShares{
						DelegatorAddress: zone.DelegationAddress.Address,
						Amount:           sdk.NewCoin(vals[1]+"/2", sdk.NewInt(200)),
					},
				}
			},
			expectedDelegationRecords: func(ctx sdk.Context, qs *app.Quicksilver, zone types.Zone) []types.Delegation {
				vals := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				return []types.Delegation{
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[0],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000)),
					},
					{
						DelegationAddress: zone.DelegationAddress.Address,
						ValidatorAddress:  vals[1],
						Amount:            sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1200)),
					},
				}
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("connectionID"), suite.path.EndpointA.ConnectionID))

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			if !found {
				suite.Fail("unable to retrieve zone for test")
			}

			for _, dr := range test.delegationRecords(ctx, quicksilver, zone) {
				quicksilver.InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, dr)
			}

			t := ctx.BlockTime().Add(-time.Minute * 2)

			for idx, msg := range test.msgs(ctx, quicksilver, zone) {
				txHash := randomutils.GenerateRandomHashAsHex(32)
				lsmMsg, ok := msg.(*lsmstakingtypes.MsgRedeemTokensForShares)
				suite.True(ok)
				quicksilver.InterchainstakingKeeper.SetReceipt(ctx, types.Receipt{
					ChainId:   suite.chainB.ChainID,
					Sender:    lsmMsg.DelegatorAddress,
					Txhash:    txHash,
					Amount:    sdk.NewCoins(lsmMsg.Amount),
					FirstSeen: &t,
				})

				err := quicksilver.InterchainstakingKeeper.HandleRedeemTokens(ctx, msg, sdk.NewCoin(zone.BaseDenom, lsmMsg.Amount.Amount), txHash, zone.ConnectionId)
				if test.errs[idx] {
					suite.Error(err)
				} else {
					suite.NoError(err)
				}
			}

			for _, edr := range test.expectedDelegationRecords(ctx, quicksilver, zone) {
				dr, found := quicksilver.InterchainstakingKeeper.GetDelegation(ctx, zone.ChainId, edr.DelegationAddress, edr.ValidatorAddress)
				suite.True(found)
				suite.Equal(dr.Amount, edr.Amount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleFailedDelegate_Batch_OK() {
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	zone.SetWithdrawalWaitgroup(app.Logger(), 100, "init")
	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	vals := app.InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
	msg := stakingtypes.MsgDelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}
	var msgMsg sdk.Msg = &msg
	err := app.InterchainstakingKeeper.HandleFailedDelegate(ctx, msgMsg, "batch/12345678")
	suite.NoError(err)

	zone, found = app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)
	suite.Equal(uint32(99), zone.GetWithdrawalWaitgroup())
}

func (suite *KeeperTestSuite) TestHandleFailedDelegate_PerfAddress_OK() {
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	zone.SetWithdrawalWaitgroup(app.Logger(), 100, "init")
	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	vals := app.InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
	msg := stakingtypes.MsgDelegate{DelegatorAddress: zone.PerformanceAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}
	var msgMsg sdk.Msg = &msg
	err := app.InterchainstakingKeeper.HandleFailedDelegate(ctx, msgMsg, "batch/12345678")
	suite.NoError(err)

	zone, found = app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)
	// delegator was perf address, no change in waitgroup
	suite.Equal(uint32(100), zone.GetWithdrawalWaitgroup())
}

func (suite *KeeperTestSuite) TestHandleFailedDelegate_NotBatch_OK() {
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	zone.SetWithdrawalWaitgroup(app.Logger(), 100, "init")
	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	vals := app.InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
	msg := stakingtypes.MsgDelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}
	var msgMsg sdk.Msg = &msg
	err := app.InterchainstakingKeeper.HandleFailedDelegate(ctx, msgMsg, randomutils.GenerateRandomHashAsHex(32))
	suite.NoError(err)

	zone, found = app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)
	// memo was not a batch id, so don't decrement withdrawal wg
	suite.Equal(uint32(100), zone.GetWithdrawalWaitgroup())
}

func (suite *KeeperTestSuite) TestHandleFailedDelegate_BatchTriggerRR_OK() {
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	zone.SetWithdrawalWaitgroup(app.Logger(), 1, "init")
	app.InterchainstakingKeeper.SetZone(ctx, &zone)
	preQueries := app.InterchainQueryKeeper.AllQueries(ctx)

	vals := app.InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
	msg := stakingtypes.MsgDelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}
	var msgMsg sdk.Msg = &msg
	err := app.InterchainstakingKeeper.HandleFailedDelegate(ctx, msgMsg, "batch/12345678")
	suite.NoError(err)

	zone, found = app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)
	// memo was not a batch id, so don't decrement withdrawal wg
	suite.Equal(uint32(0), zone.GetWithdrawalWaitgroup())

	postQueries := app.InterchainQueryKeeper.AllQueries(ctx)

	// we should have exactly one additional query
	suite.Equal(len(postQueries), len(preQueries)+1)

	distributeRewardsPreQueryCount := 0
	distributeRewardsPostQueryCount := 0
	for _, q := range preQueries {
		if q.CallbackId == "distributerewards" {
			distributeRewardsPreQueryCount++
		}
	}

	for _, q := range postQueries {
		if q.CallbackId == "distributerewards" {
			distributeRewardsPostQueryCount++
		}
	}

	suite.Equal(distributeRewardsPostQueryCount, distributeRewardsPreQueryCount+1)
}

func (suite *KeeperTestSuite) TestHandleFailedDelegate_BadAddr_Fail() {
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	zone.SetWithdrawalWaitgroup(app.Logger(), 100, "init")
	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	vals := app.InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
	msg := stakingtypes.MsgDelegate{DelegatorAddress: addressutils.GenerateAddressForTestWithPrefix("cosmos"), ValidatorAddress: vals[0], Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}
	var msgMsg sdk.Msg = &msg
	err := app.InterchainstakingKeeper.HandleFailedDelegate(ctx, msgMsg, randomutils.GenerateRandomHashAsHex(32))
	suite.ErrorContains(err, "unable to find zone for address")
}

func (suite *KeeperTestSuite) TestHandleFailedDelegate_BadMsg_Fail() {
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	zone.SetWithdrawalWaitgroup(app.Logger(), 100, "init")
	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	vals := app.InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
	msg := stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: vals[0], ValidatorDstAddress: vals[1], Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}
	var msgMsg sdk.Msg = &msg
	err := app.InterchainstakingKeeper.HandleFailedDelegate(ctx, msgMsg, "batch/12345678")
	suite.ErrorContains(err, "unable to cast source message to MsgDelegate")
}
