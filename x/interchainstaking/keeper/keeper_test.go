package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	"github.com/stretchr/testify/suite"

	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	tmclienttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	ics "github.com/ingenuity-build/quicksilver/x/interchainstaking"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var testAddress = utils.GenerateAccAddressForTest().String()

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	path   *ibctesting.Path
}

func (s *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
	quicksilver, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return quicksilver
}

// SetupTest creates a coordinator with 2 test chains.
func (s *KeeperTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)         // initializes 2 test chains
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	s.path = newQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)
}

func (s *KeeperTestSuite) setupTestZones() {
	proposal := &icstypes.RegisterZoneProposal{
		Title:            "register zone A",
		Description:      "register zone A",
		ConnectionId:     s.path.EndpointA.ConnectionID,
		LocalDenom:       "uqatom",
		BaseDenom:        "uatom",
		AccountPrefix:    "cosmos",
		ReturnToSender:   false,
		UnbondingEnabled: false,
		LiquidityModule:  true,
		DepositsEnabled:  true,
		Decimals:         6,
	}

	quicksilver := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	err := quicksilver.InterchainstakingKeeper.HandleRegisterZoneProposal(ctx, proposal)
	s.Require().NoError(err)

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().True(found)

	quicksilver.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", &tmclienttypes.ClientState{ChainId: s.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
	quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: ctx.BlockTime()})
	quicksilver.IBCKeeper.ConnectionKeeper.SetConnection(ctx, s.path.EndpointA.ConnectionID, connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0"})
	s.Require().NoError(s.setupChannelForICA(ctx, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "deposit", zone.AccountPrefix))
	s.Require().NoError(s.setupChannelForICA(ctx, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "withdrawal", zone.AccountPrefix))
	s.Require().NoError(s.setupChannelForICA(ctx, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "performance", zone.AccountPrefix))
	s.Require().NoError(s.setupChannelForICA(ctx, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "delegate", zone.AccountPrefix))
	zone, found = quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext())
	for i := range vals {
		s.Require().NoError(quicksilver.InterchainstakingKeeper.SetValidatorForZone(s.chainA.GetContext(), &zone, app.DefaultConfig().Codec.MustMarshal(&vals[i])))
	}

	s.coordinator.CommitNBlocks(s.chainA, 2)
	s.coordinator.CommitNBlocks(s.chainB, 2)
}

func (s *KeeperTestSuite) setupChannelForICA(ctx sdk.Context, chainID, connectionID, accountSuffix, remotePrefix string) error {
	quicksilver := s.GetQuicksilverApp(s.chainA)

	ibcModule := ics.NewIBCModule(quicksilver.InterchainstakingKeeper)
	portID, err := icatypes.NewControllerPortID(chainID + "." + accountSuffix)
	if err != nil {
		return err
	}

	quicksilver.InterchainstakingKeeper.SetConnectionForPort(ctx, connectionID, portID)

	channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctx)
	quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: portID, ChannelId: channelID}, ConnectionHops: []string{connectionID}})

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

	addr, err := bech32.ConvertAndEncode(remotePrefix, utils.GenerateAccAddressForTest())
	if err != nil {
		return err
	}
	quicksilver.ICAControllerKeeper.SetInterchainAccountAddress(ctx, connectionID, portID, addr)
	return ibcModule.OnChanOpenAck(ctx, portID, channelID, "", "")
}

func (s *KeeperTestSuite) giveFunds(ctx sdk.Context, denom string, amount int64, address string) {
	quicksilver := s.GetQuicksilverApp(s.chainA)

	balance := sdk.NewCoins(
		sdk.NewCoin(
			denom,
			math.NewInt(amount),
		),
	)
	err := quicksilver.MintKeeper.MintCoins(ctx, balance)
	s.Require().NoError(err)
	addr, err := utils.AccAddressFromBech32(address, "")
	s.Require().NoError(err)
	err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, balance)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestGetDelegatedAmount() {
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
		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			quicksilver := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			for _, delegation := range tt.delegations(ctx, quicksilver, zone) {
				icsKeeper.SetDelegation(ctx, &zone, delegation)
			}

			actual := icsKeeper.GetDelegatedAmount(ctx, &zone)
			s.Require().Equal(tt.expected, actual.Amount)
			s.Require().Equal(zone.BaseDenom, actual.Denom)
		})
	}
}

func (s *KeeperTestSuite) TestGetUnbondingAmount() {
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
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			expected: math.NewInt(3000000),
		},
		{
			name: "one non-unbonding withdrawal",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusQueued, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			expected: math.ZeroInt(),
		},
		{
			name: "multi unbonding withdrawal",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out,
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
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
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icskeeper.WithdrawStatusQueued, Txhash: utils.GenerateRandomHashAsHex()},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusCompleted, Txhash: utils.GenerateRandomHashAsHex()},
				)
				return out
			},
			expected: math.NewInt(3000000),
		},
	}

	for _, tt := range tc {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			quicksilver := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			for _, record := range tt.records(zone) {
				icsKeeper.SetWithdrawalRecord(ctx, record)
			}

			actual := icsKeeper.GetUnbondingAmount(ctx, &zone)
			s.Require().Equal(tt.expected, actual.Amount)
			s.Require().Equal(zone.BaseDenom, actual.Denom)
		})
	}
}

func (s *KeeperTestSuite) TestGetRatio() {
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
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
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
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
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
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusCompleted, Txhash: utils.GenerateRandomHashAsHex()})
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
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
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
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
					icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()},
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
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
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
		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			quicksilver := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			for _, record := range tt.records(ctx, quicksilver, zone) {
				icsKeeper.SetWithdrawalRecord(ctx, record)
			}

			for _, delegation := range tt.delegations(ctx, quicksilver, zone) {
				icsKeeper.SetDelegation(ctx, &zone, delegation)
			}

			err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, tt.supply)))
			s.Require().NoError(err)

			actual, isZero := icsKeeper.GetRatio(ctx, &zone, sdk.ZeroInt())
			s.Require().Equal(tt.supply.IsZero(), isZero)
			s.Require().Equal(tt.expected, actual)
		})
	}
}

func (s *KeeperTestSuite) TestUpdateRedemptionRate() {
	s.SetupTest()
	s.setupTestZones()

	quicksilver := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()
	icsKeeper := quicksilver.InterchainstakingKeeper
	zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(3000))))
	s.Require().NoError(err)

	// no change!
	s.Require().Equal(sdk.OneDec(), zone.RedemptionRate)
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.ZeroInt())

	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(sdk.OneDec(), zone.RedemptionRate)

	// add 1%
	s.Require().Equal(sdk.OneDec(), zone.RedemptionRate)
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.NewInt(30))
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(10)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(10)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(10)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(sdk.NewDecWithPrec(101, 2), zone.RedemptionRate)

	// add >2%; cap at 2%
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.NewInt(500))
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(166)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(167)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(167)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	// should be capped at 2% increase. (1.01*1.02 == 1.0302)
	s.Require().Equal(sdk.NewDecWithPrec(10302, 4), zone.RedemptionRate)

	// add nothing, still cap at 2%
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.ZeroInt())
	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	// should be capped at 2% increase. (1.01*1.02*1.02 == 1.050804)
	s.Require().Equal(sdk.NewDecWithPrec(1050804, 6), zone.RedemptionRate)

	delegationA.Amount.Amount = delegationA.Amount.Amount.SubRaw(500)
	delegationB.Amount.Amount = delegationB.Amount.Amount.SubRaw(500)
	delegationC.Amount.Amount = delegationC.Amount.Amount.SubRaw(500)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	// remove > 5%, cap at -5%
	icsKeeper.UpdateRedemptionRate(ctx, &zone, sdk.ZeroInt())
	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	s.Require().Equal(sdk.NewDecWithPrec(9982638, 7), zone.RedemptionRate)
}

func (s *KeeperTestSuite) TestOverrideRedemptionRateNoCap() {
	s.SetupTest()
	s.setupTestZones()

	quicksilver := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()
	icsKeeper := quicksilver.InterchainstakingKeeper
	zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)

	err := quicksilver.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(3000))))
	s.Require().NoError(err)

	// no change!
	s.Require().Equal(sdk.OneDec(), zone.RedemptionRate)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)

	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(sdk.OneDec(), zone.RedemptionRate)

	// add 1%
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(10)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(10)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(10)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)

	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(sdk.NewDecWithPrec(101, 2), zone.RedemptionRate)

	// add >2%; no cap
	delegationA.Amount.Amount = delegationA.Amount.Amount.AddRaw(166)
	delegationB.Amount.Amount = delegationB.Amount.Amount.AddRaw(167)
	delegationC.Amount.Amount = delegationC.Amount.Amount.AddRaw(167)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)

	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(sdk.NewDecWithPrec(1176666666666666667, 18), zone.RedemptionRate)

	// add nothing, no change
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)
	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)
	s.Require().Equal(sdk.NewDecWithPrec(1176666666666666667, 18), zone.RedemptionRate)

	delegationA.Amount.Amount = delegationA.Amount.Amount.SubRaw(500)
	delegationB.Amount.Amount = delegationB.Amount.Amount.SubRaw(500)
	delegationC.Amount.Amount = delegationC.Amount.Amount.SubRaw(500)
	icsKeeper.SetDelegation(ctx, &zone, delegationA)
	icsKeeper.SetDelegation(ctx, &zone, delegationB)
	icsKeeper.SetDelegation(ctx, &zone, delegationC)
	icsKeeper.OverrideRedemptionRateNoCap(ctx, &zone)
	zone, found = icsKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	s.Require().Equal(sdk.NewDecWithPrec(676666666666666667, 18), zone.RedemptionRate)
}
