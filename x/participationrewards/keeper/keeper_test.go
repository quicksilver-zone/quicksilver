package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v5/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	tmclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	ics "github.com/ingenuity-build/quicksilver/x/interchainstaking"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
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

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	path *ibctesting.Path
}

func (s *KeeperTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *app.Quicksilver {
	app, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return app
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)         // initializes 2 test chains
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	suite.path = newQuicksilverPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)

	suite.coordinator.CurrentTime = time.Now().UTC()
	suite.coordinator.UpdateTime()

	suite.coreTest()
}

func (suite *KeeperTestSuite) coreTest() {
	qApp := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	// test zone
	performanceAddressTest := utils.GenerateAccAddressForTestWithPrefix("bcosmos")
	performanceAccountTest, err := icstypes.NewICAAccount(performanceAddressTest, fmt.Sprintf("%s.performance", suite.chainB.ChainID), "uatom")
	suite.Require().NoError(err)
	performanceAccountTest.WithdrawalAddress = utils.GenerateAccAddressForTestWithPrefix("bcosmos")

	zone := icstypes.Zone{
		ConnectionId:       suite.path.EndpointA.ConnectionID,
		ChainId:            suite.chainB.ChainID,
		AccountPrefix:      "bcosmos",
		LocalDenom:         "uqatom",
		BaseDenom:          "uatom",
		MultiSend:          true,
		LiquidityModule:    true,
		PerformanceAddress: performanceAccountTest,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)

	qApp.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", &tmclienttypes.ClientState{ChainId: suite.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
	qApp.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: ctx.BlockTime()})
	qApp.IBCKeeper.ConnectionKeeper.SetConnection(ctx, suite.path.EndpointA.ConnectionID, connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0"})
	suite.Require().NoError(setupChannelForICA(ctx, qApp, suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "deposit", zone.AccountPrefix))

	for _, val := range suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext()) {
		// refetch the zone for each validator, else we end up with an empty valset each time!
		zone, found := qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
		suite.Require().True(found)
		suite.Require().NoError(icskeeper.SetValidatorForZone(&qApp.InterchainstakingKeeper, suite.chainA.GetContext(), zone, app.DefaultConfig().Codec.MustMarshal(&val)))
	}

	// cosmos zone
	performanceAddressCosmos := utils.GenerateAccAddressForTestWithPrefix("cosmos")
	performanceAccountCosmos, err := icstypes.NewICAAccount(performanceAddressCosmos, "cosmoshub-4.performance", "uatom")
	suite.Require().NoError(err)
	performanceAccountCosmos.WithdrawalAddress = utils.GenerateAccAddressForTestWithPrefix("cosmos")

	zoneCosmos := icstypes.Zone{
		ConnectionId:       "connection-77001",
		ChainId:            "cosmoshub-4",
		AccountPrefix:      "cosmos",
		LocalDenom:         "uqatom",
		BaseDenom:          "uatom",
		MultiSend:          true,
		LiquidityModule:    true,
		PerformanceAddress: performanceAccountCosmos,
		Validators: []*icstypes.Validator{
			{
				ValoperAddress:  "cosmosvaloper1759teakrsvnx7rnur8ezc4qaq8669nhtgukm0x",
				CommissionRate:  sdk.MustNewDecFromStr("0.1"),
				DelegatorShares: sdk.NewDec(200032604739),
				VotingPower:     math.NewInt(200032604739),
				Score:           sdk.ZeroDec(),
			},
			{
				ValoperAddress:  "cosmosvaloper1jtjjyxtqk0fj85ud9cxk368gr8cjdsftvdt5jl",
				CommissionRate:  sdk.MustNewDecFromStr("0.1"),
				DelegatorShares: sdk.NewDec(200032604734),
				VotingPower:     math.NewInt(200032604734),
				Score:           sdk.ZeroDec(),
			},
			{
				ValoperAddress:  "cosmosvaloper1q86m0zq0p52h4puw5pg5xgc3c5e2mq52y6mth0",
				CommissionRate:  sdk.MustNewDecFromStr("0.1"),
				DelegatorShares: sdk.NewDec(200032604738),
				VotingPower:     math.NewInt(200032604738),
				Score:           sdk.ZeroDec(),
			},
		},
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneCosmos)

	// osmosis zone
	zoneOsmosis := icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       true,
		LiquidityModule: true,
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           utils.GenerateAccAddressForTestWithPrefix("osmo"),
			PortName:          "cosmoshub-4.performance",
			WithdrawalAddress: utils.GenerateAccAddressForTestWithPrefix("osmo"),
		},
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneOsmosis)

	// connection type
	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0),
		suite.chainB.ChainID,
	)
	// osmosis
	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisParams,
		fmt.Sprintf("{\"ChainID\": %q}", "osmosis-1"),
		"",
	)

	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", "connection-77002", "osmosis-1", 0),
		"osmosis-1",
	)

	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisPool,
		fmt.Sprintf(
			"{\"poolid\":%d,\"poolname\":%q,\"pooltype\":\"balancer\",\"zones\":{%q:%q,%q:%q}}",
			1,
			"atom/osmo",
			"cosmoshub-4",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"osmosis-1",
			"ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159",
		),
		"1",
	)

	suite.addProtocolData(
		types.ProtocolDataTypeLiquidToken,
		fmt.Sprintf(
			"{\"chainid\":%q,\"originchainid\":%q,\"denom\":%q,\"localdenom\":%q}",
			"osmosis-1",
			"cosmoshub-4",
			"uatom",
			"uqatom",
		),
		"osmosis-1/uqatom",
	)

	// advance the chains
	suite.coordinator.CommitNBlocks(suite.chainA, 1)
	suite.coordinator.CommitNBlocks(suite.chainB, 1)

	// callback test
	suite.executeOsmosisPoolUpdateCallback()

	// add some deposits
	suite.addReceipt(
		&zoneCosmos,
		testAddress,
		"testTxHash",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(120000000))),
	)

	// add some claims
	suite.addClaim(
		testAddress,
		"cosmoshub-4",
		cmtypes.ClaimTypeLiquidToken,
		"osmosis-1",
		40000000,
	)

	// Epoch boundary
	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.AfterEpochEnd(suite.chainA.GetContext(), "epoch", 3)

	suite.executeValidatorSelectionRewardsCallback(performanceAddressTest)
}

func (suite *KeeperTestSuite) addProtocolData(Type types.ProtocolDataType, Data string, Key string) {
	pd := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(Type)],
		Data: []byte(Data),
	}

	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.SetProtocolData(suite.chainA.GetContext(), Key, &pd)
}

func (suite *KeeperTestSuite) addReceipt(zone *icstypes.Zone, sender string, hash string, coins sdk.Coins) {
	receipt := icstypes.Receipt{
		ChainId: zone.ChainId,
		Sender:  sender,
		Txhash:  hash,
		Amount:  coins,
	}

	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetReceipt(suite.chainA.GetContext(), receipt)

	delegationAddress := utils.GenerateAccAddressForTestWithPrefix("cosmos")
	validatorAddress := utils.GenerateValAddressForTestWithPrefix("cosmos")
	delegation := icstypes.Delegation{
		DelegationAddress: delegationAddress,
		ValidatorAddress:  validatorAddress,
		Amount:            coins[0],
		Height:            1,
		RedelegationEnd:   101,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetDelegation(suite.chainA.GetContext(), zone, delegation)
}

func (suite *KeeperTestSuite) addClaim(address string, chainID string, claimType cmtypes.ClaimType, sourceChainID string, amount uint64) {
	claim := cmtypes.Claim{
		UserAddress:   address,
		ChainId:       chainID,
		Module:        claimType,
		SourceChainId: sourceChainID,
		Amount:        amount,
	}
	suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper.SetClaim(suite.chainA.GetContext(), &claim)
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
