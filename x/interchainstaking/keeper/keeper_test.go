package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	testsuite "github.com/stretchr/testify/suite"

	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v5/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	tmclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/utils/randomutils"
	ics "github.com/ingenuity-build/quicksilver/x/interchainstaking"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var testAddress = addressutils.GenerateAccAddressForTest().String()

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	testsuite.Run(t, new(KeeperTestSuite))
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

type KeeperTestSuite struct {
	testsuite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	path   *ibctesting.Path
}

func (suite *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
	quicksilver, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return quicksilver
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)         // initializes 2 test chains
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	suite.path = newQuicksilverPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)
}

func (suite *KeeperTestSuite) setupTestZones() {
	proposal := &icstypes.RegisterZoneProposal{
		Title:            "register zone A",
		Description:      "register zone A",
		ConnectionId:     suite.path.EndpointA.ConnectionID,
		LocalDenom:       "uqatom",
		BaseDenom:        "uatom",
		AccountPrefix:    "cosmos",
		ReturnToSender:   false,
		UnbondingEnabled: false,
		LiquidityModule:  true,
		DepositsEnabled:  true,
		Decimals:         6,
		Is_118:           true,
	}

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	err := quicksilver.InterchainstakingKeeper.HandleRegisterZoneProposal(ctx, proposal)
	suite.Require().NoError(err)

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)

	quicksilver.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", &tmclienttypes.ClientState{ChainId: suite.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
	quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: ctx.BlockTime()})
	quicksilver.IBCKeeper.ConnectionKeeper.SetConnection(ctx, suite.path.EndpointA.ConnectionID, connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0"})
	suite.Require().NoError(suite.setupChannelForICA(ctx, suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "deposit", zone.AccountPrefix))
	suite.Require().NoError(suite.setupChannelForICA(ctx, suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "withdrawal", zone.AccountPrefix))
	suite.Require().NoError(suite.setupChannelForICA(ctx, suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "performance", zone.AccountPrefix))
	suite.Require().NoError(suite.setupChannelForICA(ctx, suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "delegate", zone.AccountPrefix))
	zone, found = quicksilver.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext())
	for i := range vals {
		suite.Require().NoError(quicksilver.InterchainstakingKeeper.SetValidatorForZone(suite.chainA.GetContext(), &zone, app.DefaultConfig().Codec.MustMarshal(&vals[i])))
	}

	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)
}

func (suite *KeeperTestSuite) setupChannelForICA(ctx sdk.Context, chainID, connectionID, accountSuffix, remotePrefix string) error {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	ibcModule := ics.NewIBCModule(quicksilver.InterchainstakingKeeper)
	portID, err := icatypes.NewControllerPortID(chainID + "." + accountSuffix)
	if err != nil {
		return err
	}

	quicksilver.InterchainstakingKeeper.SetConnectionForPort(ctx, connectionID, portID)

	channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctx)
	quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icatypes.PortID, ChannelId: channelID}, ConnectionHops: []string{connectionID}})

	// channel, found := quicksilver.IBCKeeper.ChannelKeeper.GetChannel(ctx, portID, channelID)
	// suite.Require().True(found)
	// fmt.Printf("DEBUG: channel >>>\n%v\n<<<\n", channel)

	quicksilver.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, portID, channelID, 1)
	quicksilver.ICAControllerKeeper.SetActiveChannelID(ctx, connectionID, portID, channelID)
	key, err := quicksilver.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		ctx,
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}
	err = quicksilver.GetScopedIBCKeeper().ClaimCapability(
		ctx,
		key,
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}

	key, err = quicksilver.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		ctx,
		host.PortPath(portID),
	)
	if err != nil {
		return err
	}
	err = quicksilver.GetScopedIBCKeeper().ClaimCapability(
		ctx,
		key,
		host.PortPath(portID),
	)
	if err != nil {
		return err
	}

	addr, err := addressutils.EncodeAddressToBech32(remotePrefix, addressutils.GenerateAccAddressForTest())
	if err != nil {
		return err
	}
	quicksilver.ICAControllerKeeper.SetInterchainAccountAddress(ctx, connectionID, portID, addr)
	return ibcModule.OnChanOpenAck(ctx, portID, channelID, "", "")
}

func (suite *KeeperTestSuite) giveFunds(ctx sdk.Context, denom string, amount int64, address string) {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	balance := sdk.NewCoins(
		sdk.NewCoin(
			denom,
			math.NewInt(amount),
		),
	)
	err := quicksilver.MintKeeper.MintCoins(ctx, balance)
	suite.Require().NoError(err)
	addr, err := addressutils.AccAddressFromBech32(address, "")
	suite.Require().NoError(err)
	err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, balance)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestGetDelegatedAmount() {
	tc := []struct {
		name        string
		delegations func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation
		expected    math.Int
	}{
		{
			name: "empty delegations",
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				return out
			},
			expected: math.NewInt(0),
		},
		{
			name: "one delegation",
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DelegationAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			expected: math.NewInt(3000000),
		},
		{
			name: "multi delegation",
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out,
					icstypes.NewDelegation(zone.DelegationAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))),
					icstypes.NewDelegation(zone.DelegationAddress.Address, validators[1], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000000))),
					icstypes.NewDelegation(zone.DelegationAddress.Address, validators[2], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(20000000))),
				)
				return out
			},
			expected: math.NewInt(40000000),
		},
	}

	for _, tt := range tc {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			for _, delegation := range tt.delegations(ctx, quicksilver, zone) {
				icsKeeper.SetDelegation(ctx, &zone, delegation)
			}

			actual := icsKeeper.GetDelegatedAmount(ctx, &zone)
			suite.Require().Equal(tt.expected, actual.Amount)
			suite.Require().Equal(zone.BaseDenom, actual.Denom)
		})
	}
}

func (suite *KeeperTestSuite) TestGetUnbondingAmount() {
	tc := []struct {
		name     string
		records  func(zone icstypes.Zone) []icstypes.WithdrawalRecord
		expected math.Int
	}{
		{
			name: "no withdrawals",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				return out
			},
			expected: math.ZeroInt(),
		},
		{
			name: "one unbonding withdrawal",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)})
				return out
			},
			expected: math.NewInt(3000000),
		},
		{
			name: "one non-unbonding withdrawal",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusQueued, Txhash: randomutils.GenerateRandomHashAsHex(64)})
				return out
			},
			expected: math.ZeroInt(),
		},
		{
			name: "multi unbonding withdrawal",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out,
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
				)
				return out
			},
			expected: math.NewInt(14500000),
		},
		{
			name: "multi mixed withdrawal",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out,
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icstypes.WithdrawStatusQueued, Txhash: randomutils.GenerateRandomHashAsHex(64)},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icstypes.WithdrawStatusCompleted, Txhash: randomutils.GenerateRandomHashAsHex(64)},
				)
				return out
			},
			expected: math.NewInt(3000000),
		},
	}

	for _, tt := range tc {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			for _, record := range tt.records(zone) {
				icsKeeper.SetWithdrawalRecord(ctx, record)
			}

			actual := icsKeeper.GetUnbondingAmount(ctx, &zone)
			suite.Require().Equal(tt.expected, actual.Amount)
			suite.Require().Equal(zone.BaseDenom, actual.Denom)
		})
	}
}

func (suite *KeeperTestSuite) TestGetRatio() {
	tc := []struct {
		name        string
		records     func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord
		delegations func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation
		supply      math.Int
		expected    sdk.Dec
	}{
		{
			name: "no withdrawals, no unbonding, no supply",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				return out
			},
			supply:   math.ZeroInt(),
			expected: sdk.OneDec(),
		},
		{
			name: "no withdrawals, one delegation, expect 1.0",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(3000000),
			expected: sdk.OneDec(),
		},
		{
			name: "one withdrawal, no delegation, expect 1.0",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)})
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				return out
			},
			supply:   math.NewInt(3000000),
			expected: sdk.OneDec(),
		},
		{
			name: "one withdrawals, one delegation, one unbonding, expect 1.0",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)})
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(6000000),
			expected: sdk.OneDec(),
		},
		{
			name: "one non-unbond withdrawals, one delegation, one unbonding, expect 1.0",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusCompleted, Txhash: randomutils.GenerateRandomHashAsHex(64)})
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(3000000),
			expected: sdk.OneDec(),
		},
		{
			name: "multi unbonding withdrawal, delegation, expect 1.0",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out,
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
				)
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(17500000),
			expected: sdk.OneDec(),
		},
		{
			name: "multi unbonding withdrawal, delegation, sub 1.0",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out,
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)},
				)
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(10000000),
			expected: sdk.NewDecWithPrec(75, 2),
		},
		{
			name: "multi unbonding withdrawal, delegation, gt 1.0",
			records: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icstypes.WithdrawStatusUnbond, Txhash: randomutils.GenerateRandomHashAsHex(64)})
				return out
			},
			delegations: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) []icstypes.Delegation {
				validators := qs.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)
				out := make([]icstypes.Delegation, 0)
				out = append(out,
					icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))),
					icstypes.NewDelegation(zone.DepositAddress.Address, validators[1], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(16000000))),
					icstypes.NewDelegation(zone.DepositAddress.Address, validators[2], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(8000000))),
				)
				return out
			},
			supply:   math.NewInt(22500000),
			expected: sdk.NewDec(4).Quo(sdk.NewDec(3)),
		},
	}

	for _, tt := range tc {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.Require().True(found)

			for _, record := range tt.records(ctx, quicksilver, zone) {
				icsKeeper.SetWithdrawalRecord(ctx, record)
			}

			for _, delegation := range tt.delegations(ctx, quicksilver, zone) {
				icsKeeper.SetDelegation(ctx, &zone, delegation)
			}

			err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, tt.supply)))
			suite.Require().NoError(err)

			actual, isZero := icsKeeper.GetRatio(ctx, &zone, sdk.ZeroInt())
			suite.Require().Equal(tt.supply.IsZero(), isZero)
			suite.Require().Equal(tt.expected, actual)
		})
	}
}

func (suite *KeeperTestSuite) TestUpdateRedemptionRate() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	icsKeeper := quicksilver.InterchainstakingKeeper
	zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(3000))))
	suite.Require().NoError(err)

	// no change!
	suite.Require().Equal(sdk.OneDec(), zone.RedemptionRate)
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.ZeroInt())

	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	suite.Require().Equal(sdk.OneDec(), zone.RedemptionRate)

	// add 1%
	suite.Require().Equal(sdk.OneDec(), zone.RedemptionRate)
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.NewInt(30))
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(10)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(10)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(10)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewDecWithPrec(101, 2), zone.RedemptionRate)

	// add >2%; cap at 2%
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.NewInt(500))
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(166)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(167)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(167)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	// should be capped at 2% increase. (1.01*1.02 == 1.0302)
	suite.Require().Equal(sdk.NewDecWithPrec(10302, 4), zone.RedemptionRate)

	// add nothing, still cap at 2%
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.ZeroInt())
	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	// should be capped at 2% increase. (1.01*1.02*1.02 == 1.050804)
	suite.Require().Equal(sdk.NewDecWithPrec(1050804, 6), zone.RedemptionRate)

	delegationA.Amount.Amount = delegationA.Amount.Amount.SubRaw(500)
	delegationB.Amount.Amount = delegationB.Amount.Amount.SubRaw(500)
	delegationC.Amount.Amount = delegationC.Amount.Amount.SubRaw(500)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	// remove > 5%, cap at -5%
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.ZeroInt())
	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	suite.Require().Equal(sdk.NewDecWithPrec(9982638, 7), zone.RedemptionRate)
}

func (suite *KeeperTestSuite) TestOverrideRedemptionRateNoCap() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	icsKeeper := quicksilver.InterchainstakingKeeper
	zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(3000))))
	suite.Require().NoError(err)

	// no change!
	suite.Require().Equal(sdk.OneDec(), zone.RedemptionRate)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)

	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	suite.Require().Equal(sdk.OneDec(), zone.RedemptionRate)

	// add 1%
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(10)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(10)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(10)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)

	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewDecWithPrec(101, 2), zone.RedemptionRate)

	// add >2%; no cap
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(166)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(167)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(167)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)

	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewDecWithPrec(1176666666666666667, 18), zone.RedemptionRate)

	// add nothing, no change
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)
	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewDecWithPrec(1176666666666666667, 18), zone.RedemptionRate)

	delegationA.Amount.Amount = delegationA.Amount.Amount.SubRaw(500)
	delegationB.Amount.Amount = delegationB.Amount.Amount.SubRaw(500)
	delegationC.Amount.Amount = delegationC.Amount.Amount.SubRaw(500)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)
	zone, found = icsKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	suite.Require().Equal(sdk.NewDecWithPrec(676666666666666667, 18), zone.RedemptionRate)
}
