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

	suite.setupTestZones()

	// test ProtocolData
	akpd := qApp.ParticipationRewardsKeeper.AllKeyedProtocolDatas(suite.chainA.GetContext())
	// initially we expect none
	suite.Require().Equal(0, len(akpd))

	suite.setupTestProtocolData()

	akpd = qApp.ParticipationRewardsKeeper.AllKeyedProtocolDatas(suite.chainA.GetContext())
	// added 5 in setupTestProtocolData
	suite.Require().Equal(5, len(akpd))

	// advance the chains
	suite.coordinator.CommitNBlocks(suite.chainA, 1)
	suite.coordinator.CommitNBlocks(suite.chainB, 1)

	// callback test
	suite.executeSetEpochBlockCallback()
	suite.executeOsmosisPoolUpdateCallback()

	suite.setupTestDeposits()
	suite.setupTestIntents()

	err := qApp.ParticipationRewardsKeeper.AfterEpochEnd(suite.chainA.GetContext(), "epoch", 1)
	suite.Require().NoError(err)

	suite.setupTestClaims()

	err = qApp.ParticipationRewardsKeeper.AfterEpochEnd(suite.chainA.GetContext(), "epoch", 2)
	suite.Require().NoError(err)

	// Epoch boundary
	err = qApp.ParticipationRewardsKeeper.AfterEpochEnd(suite.chainA.GetContext(), "epoch", 3)
	suite.Require().NoError(err)

	_, found := qApp.ClaimsManagerKeeper.GetLastEpochClaim(suite.chainA.GetContext(), "cosmoshub-4", "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure", cmtypes.ClaimTypeLiquidToken, "osmosis-1")
	suite.Require().True(found)

	// zone for remote chain
	zone, found := qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)

	valRewards := make(map[string]sdk.Dec)
	for _, val := range zone.Validators {
		valRewards[val.ValoperAddress] = sdk.NewDec(100000000)
	}

	suite.executeValidatorSelectionRewardsCallback(zone.PerformanceAddress.Address, valRewards)
}

func (suite *KeeperTestSuite) setupTestZones() {
	qApp := suite.GetQuicksilverApp(suite.chainA)

	// test zone
	testzone := icstypes.Zone{
		ConnectionId:    suite.path.EndpointA.ConnectionID,
		ChainId:         suite.chainB.ChainID,
		AccountPrefix:   "bcosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       true,
		LiquidityModule: true,
	}
	qApp.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &testzone)

	qApp.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), "07-tendermint-0", &tmclienttypes.ClientState{ChainId: suite.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
	qApp.IBCKeeper.ClientKeeper.SetClientConsensusState(suite.chainA.GetContext(), "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: suite.chainA.GetContext().BlockTime()})
	suite.Require().NoError(suite.setupChannelForICA(suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "performance", testzone.AccountPrefix))

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
	qApp.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneCosmos)

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
	qApp.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneOsmosis)
}

func (suite *KeeperTestSuite) setupChannelForICA(chainID string, connectionID string, accountSuffix string, remotePrefix string) error {
	qApp := suite.GetQuicksilverApp(suite.chainA)

	ibcModule := ics.NewIBCModule(qApp.InterchainstakingKeeper)
	portID, err := icatypes.NewControllerPortID(chainID + "." + accountSuffix)
	if err != nil {
		return err
	}

	qApp.InterchainstakingKeeper.SetConnectionForPort(suite.chainA.GetContext(), connectionID, portID)

	channelID := qApp.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(suite.chainA.GetContext())
	qApp.IBCKeeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), portID, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icatypes.PortID, ChannelId: channelID}, ConnectionHops: []string{connectionID}})

	// channel, found := qApp.IBCKeeper.ChannelKeeper.GetChannel(suite.chainA.GetContext(), portID, channelID)
	// suite.Require().True(found)
	// fmt.Printf("DEBUG: channel >>>\n%v\n<<<\n", channel)

	qApp.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.chainA.GetContext(), portID, channelID, 1)
	qApp.ICAControllerKeeper.SetActiveChannelID(suite.chainA.GetContext(), connectionID, portID, channelID)
	key, err := qApp.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		suite.chainA.GetContext(),
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}
	err = qApp.GetScopedIBCKeeper().ClaimCapability(
		suite.chainA.GetContext(),
		key,
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}

	key, err = qApp.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		suite.chainA.GetContext(),
		host.PortPath(portID),
	)
	if err != nil {
		return err
	}
	err = qApp.GetScopedIBCKeeper().ClaimCapability(
		suite.chainA.GetContext(),
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
	qApp.ICAControllerKeeper.SetInterchainAccountAddress(suite.chainA.GetContext(), connectionID, portID, addr)
	return ibcModule.OnChanOpenAck(suite.chainA.GetContext(), portID, channelID, "", "")
}

func (suite *KeeperTestSuite) setupTestProtocolData() {
	// connection type for ibc testsuite chainB
	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0),
		suite.chainB.ChainID,
	)
	// osmosis params
	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisParams,
		fmt.Sprintf("{\"ChainID\": %q}", "osmosis-1"),
		types.OsmosisParamsKey,
	)
	// osmosis test chain
	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", "connection-77002", "osmosis-1", 0),
		"osmosis-1",
	)
	// osmosis test pool
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
	// atom (cosmoshub) on osmosis
	suite.addProtocolData(
		types.ProtocolDataTypeLiquidToken,
		fmt.Sprintf(
			"{\"chainid\":%q,\"registeredzonechainid\":%q,\"ibcdenom\":%q,\"qassetdenom\":%q}",
			"osmosis-1",
			"cosmoshub-4",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"uqatom",
		),
		"osmosis-1/ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
	)
}

func (suite *KeeperTestSuite) addProtocolData(Type types.ProtocolDataType, Data string, Key string) {
	pd := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(Type)],
		Data: []byte(Data),
	}

	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.SetProtocolData(suite.chainA.GetContext(), Key, &pd)
}

func (suite *KeeperTestSuite) setupTestDeposits() {
	qApp := suite.GetQuicksilverApp(suite.chainA)

	// add deposit to chainB zone
	zone, found := qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)

	suite.addReceipt(
		&zone,
		testAddress,
		"testTxHash03",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(150000000))),
	)

	// add deposit to cosmos zone
	zone, found = qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), "cosmoshub-4")
	suite.Require().True(found)

	suite.addReceipt(
		&zone,
		testAddress,
		"testTxHash01",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(120000000))),
	)

	// add deposit to osmosis zone
	zone, found = qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), "osmosis-1")
	suite.Require().True(found)

	suite.addReceipt(
		&zone,
		testAddress,
		"testTxHash02",
		sdk.NewCoins(sdk.NewCoin("uosmo", math.NewInt(100000000))),
	)
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

func (suite *KeeperTestSuite) setupTestIntents() {
	qApp := suite.GetQuicksilverApp(suite.chainA)

	// chainB
	zone, found := qApp.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)

	suite.addIntent(
		testAddress,
		zone,
		icstypes.ValidatorIntents{
			{
				ValoperAddress: zone.Validators[0].ValoperAddress,
				Weight:         sdk.MustNewDecFromStr("0.3"),
			},
			{
				ValoperAddress: zone.Validators[1].ValoperAddress,
				Weight:         sdk.MustNewDecFromStr("0.4"),
			},
			{
				ValoperAddress: zone.Validators[2].ValoperAddress,
				Weight:         sdk.MustNewDecFromStr("0.3"),
			},
		},
	)
}

func (suite *KeeperTestSuite) addIntent(address string, zone icstypes.Zone, intents icstypes.ValidatorIntents) {
	intent := icstypes.DelegatorIntent{
		Delegator: address,
		Intents:   intents,
	}
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetIntent(suite.chainA.GetContext(), zone, intent, false)
}

func (suite *KeeperTestSuite) setupTestClaims() {
	// add some claims
	suite.addClaim(
		testAddress,
		"cosmoshub-4",
		cmtypes.ClaimTypeLiquidToken,
		"osmosis-1",
		40000000,
	)

	suite.addClaim(
		"quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
		"cosmoshub-4",
		cmtypes.ClaimTypeLiquidToken,
		"osmosis-1",
		0,
	)
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

/*func (suite *KeeperTestSuite) testReopenChannel(zone ) {
	qApp := suite.GetQuicksilverApp(suite.chainA)

	// connection
	connectionID, err := qApp.InterchainstakingKeeper.GetConnectionForPort(suite.chainA.GetContext(), zone.PerformanceAddress.PortName)
	suite.Require().NoError(err)
	fmt.Printf("DEBUG: connectionID %q\n", connectionID)

	// channel
	channelID, found := qApp.ICAControllerKeeper.GetActiveChannelID(suite.chainA.GetContext(), connectionID, zone.PerformanceAddress.PortName)
	suite.Require().True(found)
	fmt.Printf("DEBUG: channelID %q\n", channelID)
	channel, found := qApp.IBCKeeper.ChannelKeeper.GetChannel(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID)
	suite.Require().True(found)
	fmt.Printf("Channel: %v\n", channel)

	// close channel
	channelCap := suite.chainA.GetChannelCapability(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID)
	err = qApp.IBCKeeper.ChannelKeeper.ChanCloseInit(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID, channelCap)
	suite.Require().True(found)

	// qApp.IBCKeeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID, channeltypes.Channel{State: channeltypes.CLOSED, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: "icahost", ChannelId: channelID}, ConnectionHops: []string{connectionID}})
	channel, found = qApp.IBCKeeper.ChannelKeeper.GetChannel(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID)
	suite.Require().True(found)
	fmt.Printf("Channel: %v\n", channel)

	// // attempt to reopen channel here
	// conn, found := qApp.IBCKeeper.ConnectionKeeper.GetConnection(suite.chainA.GetContext(), connectionID)
	// suite.Require().True(found)
	// fmt.Printf("Connection: %v\n", conn)

	// fmt.Println("Channel Closed...")
	// msg := channeltypes.NewMsgChannelOpenInit(zone.PerformanceAddress.PortName, icatypes.Version, channeltypes.ORDERED, []string{connectionID}, icatypes.PortID, icatypes.ModuleName)
	// handler := qApp.MsgServiceRouter().Handler(msg)
	// _, err = handler(suite.chainA.GetContext(), msg)
	// suite.Require().NoError(err)
}*/
