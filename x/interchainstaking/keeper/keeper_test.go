package keeper_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	"github.com/stretchr/testify/suite"

	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v5/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	tmclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	ics "github.com/ingenuity-build/quicksilver/x/interchainstaking"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var TestOwnerAddress = utils.GenerateAccAddressForTest().String()

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	path   *ibctesting.Path
}

func (s *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
	app, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return app
}

func (s *KeeperTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1))
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2))

	s.path = newQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)
}

func setupChannelForICA(ctx sdk.Context, qApp *app.Quicksilver, chainID string, connectionID string, accountSuffix string, remotePrefix string) error {
	ibcModule := ics.NewIBCModule(qApp.InterchainstakingKeeper)
	portID, err := icatypes.NewControllerPortID(chainID + "." + accountSuffix)
	if err != nil {
		return err
	}
	channelID := qApp.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctx)
	qApp.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, channelID, channeltypes.Channel{ConnectionHops: []string{connectionID}, State: channeltypes.OPEN, Counterparty: channeltypes.Counterparty{PortId: "icahost", ChannelId: channelID}})
	qApp.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, portID, channelID, 1)
	qApp.ICAControllerKeeper.SetActiveChannelID(ctx, connectionID, portID, channelID)
	key, err := qApp.InterchainstakingKeeper.ScopedKeeper().NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return err
	}
	// err = qApp.InterchainstakingKeeper.ClaimCapability(ctx, key, host.ChannelCapabilityPath(portID, channelID))
	// if err != nil {
	// 	return err
	// }
	// key, err = qApp.GetScopedIBCKeeper().NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	// if err != nil {
	// 	return err
	// }
	err = qApp.GetScopedIBCKeeper().ClaimCapability(ctx, key, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return err
	}
	addr, err := bech32.ConvertAndEncode(remotePrefix, utils.GenerateAccAddressForTest())
	if err != nil {
		return err
	}
	qApp.ICAControllerKeeper.SetInterchainAccountAddress(ctx, connectionID, portID, addr)
	return ibcModule.OnChanOpenAck(ctx, portID, channelID, "", "")
}

func (s *KeeperTestSuite) SetupZones() {
	proposal := &icstypes.RegisterZoneProposal{
		Title:           "register zone A",
		Description:     "register zone A",
		ConnectionId:    s.path.EndpointA.ConnectionID,
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		AccountPrefix:   "cosmos",
		MultiSend:       true,
		LiquidityModule: true,
	}

	ctx := s.chainA.GetContext()

	// Set special testing context (e.g. for test / debug output)
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("TEST"), "TEST"))

	err := icskeeper.HandleRegisterZoneProposal(ctx, s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper, proposal)
	s.Require().NoError(err)

	// Simulate "cosmos.staking.v1beta1.Query/Validators" response

	qApp := s.GetQuicksilverApp(s.chainA)
	zone, _ := qApp.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)

	qApp.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", &tmclienttypes.ClientState{ChainId: s.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
	qApp.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: ctx.BlockTime()})
	qApp.IBCKeeper.ConnectionKeeper.SetConnection(ctx, s.path.EndpointA.ConnectionID, connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0"})
	s.Require().NoError(setupChannelForICA(ctx, qApp, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "deposit", zone.AccountPrefix))
	s.Require().NoError(setupChannelForICA(ctx, qApp, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "withdrawal", zone.AccountPrefix))
	s.Require().NoError(setupChannelForICA(ctx, qApp, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "performance", zone.AccountPrefix))
	s.Require().NoError(setupChannelForICA(ctx, qApp, s.chainB.ChainID, s.path.EndpointA.ConnectionID, "delegate", zone.AccountPrefix))

	for _, val := range s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext()) {
		// refetch the zone for each validator, else we end up with an empty valset each time!
		zone, found := qApp.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
		s.Require().True(found)
		s.Require().NoError(icskeeper.SetValidatorForZone(&qApp.InterchainstakingKeeper, s.chainA.GetContext(), zone, app.DefaultConfig().Codec.MustMarshal(&val)))
	}

	s.coordinator.CommitNBlocks(s.chainA, 2)
	s.coordinator.CommitNBlocks(s.chainB, 2)
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

func GetICSKeeper(t *testing.T) (*icskeeper.Keeper, sdk.Context) {
	app := app.Setup(t, false)
	keeper := app.InterchainstakingKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "mercury-1", Time: time.Now().UTC()})

	return &keeper, ctx
}

func (s *KeeperTestSuite) TestGetDelegatedAmount() {
	tc := []struct {
		name        string
		delegations func(zone icstypes.Zone) []icstypes.Delegation
		expected    math.Int
	}{
		{
			name: "empty delegations",
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				return out
			},
			expected: math.NewInt(0),
		},
		{
			name: "one delegation",
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			expected: math.NewInt(3000000),
		},
		{
			name: "multi delegation",
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[1], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(17000000))))
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[2], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(20000000))))
				return out
			},
			expected: math.NewInt(40000000),
		},
	}

	for _, tt := range tc {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupZones()

			qapp := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			for _, delegation := range tt.delegations(zone) {
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
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			expected: math.NewInt(14500000),
		},
		{
			name: "multi mixed withdrawal",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icskeeper.WithdrawStatusQueued, Txhash: utils.GenerateRandomHashAsHex()})
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusCompleted, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			expected: math.NewInt(3000000),
		},
	}

	for _, tt := range tc {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupZones()

			qapp := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := qapp.InterchainstakingKeeper
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
		records     func(zone icstypes.Zone) []icstypes.WithdrawalRecord
		delegations func(zone icstypes.Zone) []icstypes.Delegation
		supply      math.Int
		expected    sdk.Dec
	}{
		{
			name: "no withdrawals, no unbonding, no supply",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				return out
			},
			supply:   math.ZeroInt(),
			expected: sdk.OneDec(),
		},
		{
			name: "no withdrawals, one delegation, expect 1.0",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(3000000),
			expected: sdk.OneDec(),
		},
		{
			name: "one withdrawal, no delegation, expect 1.0",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				return out
			},
			supply:   math.NewInt(3000000),
			expected: sdk.OneDec(),
		},
		{
			name: "one withdrawals, one delegation, one unbonding, expect 1.0",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(6000000),
			expected: sdk.OneDec(),
		},
		{
			name: "one non-unbond withdrawals, one delegation, one unbonding, expect 1.0",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusCompleted, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(3000000),
			expected: sdk.OneDec(),
		},
		{
			name: "multi unbonding withdrawal, delegation, expect 1.0",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(17500000),
			expected: sdk.OneDec(),
		},
		{
			name: "multi unbonding withdrawal, delegation, sub 1.0",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1500000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				return out
			},
			supply:   math.NewInt(10000000),
			expected: sdk.NewDecWithPrec(75, 2),
		},
		{
			name: "multi unbonding withdrawal, delegation, gt 1.0",
			records: func(zone icstypes.Zone) []icstypes.WithdrawalRecord {
				out := make([]icstypes.WithdrawalRecord, 0)
				out = append(out, icstypes.WithdrawalRecord{ChainId: zone.ChainId, Delegator: zone.DelegationAddress.Address, Recipient: utils.GenerateAccAddressForTestWithPrefix(zone.AccountPrefix), Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(3000000))), Status: icskeeper.WithdrawStatusUnbond, Txhash: utils.GenerateRandomHashAsHex()})
				return out
			},
			delegations: func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[0], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[1], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(16000000))))
				out = append(out, icstypes.NewDelegation(zone.DepositAddress.Address, validators[2], sdk.NewCoin(zone.BaseDenom, sdk.NewInt(8000000))))
				return out
			},
			supply:   math.NewInt(22500000),
			expected: sdk.NewDec(4).Quo(sdk.NewDec(3)),
		},
	}

	for _, tt := range tc {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupZones()

			qapp := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)

			for _, record := range tt.records(zone) {
				icsKeeper.SetWithdrawalRecord(ctx, record)
			}

			for _, delegation := range tt.delegations(zone) {
				icsKeeper.SetDelegation(ctx, &zone, delegation)
			}

			qapp.MintKeeper.MintCoins(ctx, sdk.NewCoins(sdk.NewCoin(zone.LocalDenom, tt.supply)))

			actual := icsKeeper.GetRatio(ctx, zone, sdk.ZeroInt())
			s.Require().Equal(tt.expected, actual)
		})
	}
}
