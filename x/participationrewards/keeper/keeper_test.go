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

	s.coordinator.CurrentTime = time.Now().UTC()
	s.coordinator.UpdateTime()

	s.coreTest()
}

func (s *KeeperTestSuite) coreTest() {
	quicksilver := s.GetQuicksilverApp(s.chainA)

	s.setupTestZones()

	// test ProtocolData
	akpd := quicksilver.ParticipationRewardsKeeper.AllKeyedProtocolDatas(s.chainA.GetContext())
	// initially we expect one - the 'local' chain
	s.Require().Equal(1, len(akpd))

	s.setupTestProtocolData()

	akpd = quicksilver.ParticipationRewardsKeeper.AllKeyedProtocolDatas(s.chainA.GetContext())
	// added 6 in setupTestProtocolData
	s.Require().Equal(7, len(akpd))

	// advance the chains
	s.coordinator.CommitNBlocks(s.chainA, 1)
	s.coordinator.CommitNBlocks(s.chainB, 1)

	// callback test
	s.executeSetEpochBlockCallback()
	s.executeOsmosisPoolUpdateCallback()

	s.setupTestDeposits()
	s.setupTestIntents()

	err := quicksilver.ParticipationRewardsKeeper.AfterEpochEnd(s.chainA.GetContext(), "epoch", 1)
	s.Require().NoError(err)

	s.setupTestClaims()

	err = quicksilver.ParticipationRewardsKeeper.AfterEpochEnd(s.chainA.GetContext(), "epoch", 2)
	s.Require().NoError(err)

	// Epoch boundary
	err = quicksilver.ParticipationRewardsKeeper.AfterEpochEnd(s.chainA.GetContext(), "epoch", 3)
	s.Require().NoError(err)

	_, found := quicksilver.ClaimsManagerKeeper.GetLastEpochClaim(s.chainA.GetContext(), "cosmoshub-4", "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure", cmtypes.ClaimTypeLiquidToken, "osmosis-1")
	s.Require().True(found)

	// zone for remote chain
	zone, found := quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().True(found)

	valRewards := make(map[string]sdk.Dec)
	for _, val := range zone.Validators {
		valRewards[val.ValoperAddress] = sdk.NewDec(100000000)
	}

	s.executeValidatorSelectionRewardsCallback(zone.PerformanceAddress.Address, valRewards)
}

func (s *KeeperTestSuite) setupTestZones() {
	quicksilver := s.GetQuicksilverApp(s.chainA)

	// test zone
	testzone := icstypes.Zone{
		ConnectionId:     s.path.EndpointA.ConnectionID,
		ChainId:          s.chainB.ChainID,
		AccountPrefix:    "bcosmos",
		LocalDenom:       "uqatom",
		BaseDenom:        "uatom",
		ReturnToSender:   false,
		LiquidityModule:  true,
		DepositsEnabled:  true,
		UnbondingEnabled: false,
	}
	selftestzone := icstypes.Zone{
		ConnectionId:     s.path.EndpointB.ConnectionID,
		ChainId:          s.chainA.ChainID,
		AccountPrefix:    "osmo",
		LocalDenom:       "uqosmo",
		BaseDenom:        "uosmo",
		ReturnToSender:   false,
		LiquidityModule:  true,
		DepositsEnabled:  true,
		UnbondingEnabled: false,
	}

	quicksilver.InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &selftestzone)
	quicksilver.InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &testzone)

	quicksilver.IBCKeeper.ClientKeeper.SetClientState(s.chainA.GetContext(), "07-tendermint-0", &tmclienttypes.ClientState{ChainId: s.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})

	quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(s.chainA.GetContext(), "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: s.chainA.GetContext().BlockTime()})
	s.Require().NoError(s.setupChannelForICA(s.chainB.ChainID, s.path.EndpointA.ConnectionID, "performance", testzone.AccountPrefix))

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext())
	for i := range vals {
		// refetch the zone for each validator, else we end up with an empty valset each time!
		zone, found := quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
		s.Require().True(found)
		s.Require().NoError(quicksilver.InterchainstakingKeeper.SetValidatorForZone(s.chainA.GetContext(), &zone, app.DefaultConfig().Codec.MustMarshal(&vals[i])))
	}

	vals = s.GetQuicksilverApp(s.chainA).StakingKeeper.GetBondedValidatorsByPower(s.chainA.GetContext())
	for i := range vals {
		zone, found := quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainA.ChainID)
		s.Require().True(found)
		s.Require().NoError(quicksilver.InterchainstakingKeeper.SetValidatorForZone(s.chainA.GetContext(), &zone, app.DefaultConfig().Codec.MustMarshal(&vals[i])))
	}

	// self zone
	performanceAddressOsmo := utils.GenerateAccAddressForTestWithPrefix("osmo")
	performanceAccountOsmo, err := icstypes.NewICAAccount(performanceAddressOsmo, "self")
	s.Require().NoError(err)
	performanceAccountOsmo.WithdrawalAddress = utils.GenerateAccAddressForTestWithPrefix("osmo")

	zoneSelf := icstypes.Zone{
		ConnectionId:       "connection-77004",
		ChainId:            "testchain1",
		AccountPrefix:      "osmo",
		LocalDenom:         "uqosmo",
		BaseDenom:          "uosmo",
		ReturnToSender:     false,
		UnbondingEnabled:   false,
		LiquidityModule:    true,
		DepositsEnabled:    true,
		Decimals:           6,
		PerformanceAddress: performanceAccountOsmo,
		Validators: []*icstypes.Validator{
			{
				ValoperAddress:  "osmovaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep88n0y4",
				CommissionRate:  sdk.MustNewDecFromStr("0.1"),
				DelegatorShares: sdk.NewDec(200032604739),
				VotingPower:     math.NewInt(200032604739),
				Score:           sdk.ZeroDec(),
			},
			{
				ValoperAddress:  "osmovaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpf6t4agt",
				CommissionRate:  sdk.MustNewDecFromStr("0.1"),
				DelegatorShares: sdk.NewDec(200032604734),
				VotingPower:     math.NewInt(200032604734),
				Score:           sdk.ZeroDec(),
			},
			{
				ValoperAddress:  "osmovaloper15urq2dtp9qce4fyc85m6upwm9xul3049wh9czc",
				CommissionRate:  sdk.MustNewDecFromStr("0.1"),
				DelegatorShares: sdk.NewDec(200032604738),
				VotingPower:     math.NewInt(200032604738),
				Score:           sdk.ZeroDec(),
			},
		},
	}
	quicksilver.InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zoneSelf)

	// cosmos zone
	performanceAddressCosmos := utils.GenerateAccAddressForTestWithPrefix("cosmos")
	performanceAccountCosmos, err := icstypes.NewICAAccount(performanceAddressCosmos, "cosmoshub-4.performance")
	s.Require().NoError(err)
	performanceAccountCosmos.WithdrawalAddress = utils.GenerateAccAddressForTestWithPrefix("cosmos")

	zoneCosmos := icstypes.Zone{
		ConnectionId:       "connection-77001",
		ChainId:            "cosmoshub-4",
		AccountPrefix:      "cosmos",
		LocalDenom:         "uqatom",
		BaseDenom:          "uatom",
		ReturnToSender:     false,
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
	quicksilver.InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zoneCosmos)

	// osmosis zone
	zoneOsmosis := icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		ReturnToSender:  false,
		LiquidityModule: true,
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           utils.GenerateAccAddressForTestWithPrefix("osmo"),
			PortName:          "cosmoshub-4.performance",
			WithdrawalAddress: utils.GenerateAccAddressForTestWithPrefix("osmo"),
		},
	}
	quicksilver.InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zoneOsmosis)
}

func (s *KeeperTestSuite) setupChannelForICA(chainID, connectionID, accountSuffix, remotePrefix string) error {
	s.T().Helper()
	quicksilver := s.GetQuicksilverApp(s.chainA)

	ibcModule := ics.NewIBCModule(quicksilver.InterchainstakingKeeper)
	portID, err := icatypes.NewControllerPortID(chainID + "." + accountSuffix)
	if err != nil {
		return err
	}

	quicksilver.InterchainstakingKeeper.SetConnectionForPort(s.chainA.GetContext(), connectionID, portID)

	channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(s.chainA.GetContext())
	quicksilver.IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), portID, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icatypes.PortID, ChannelId: channelID}, ConnectionHops: []string{connectionID}})

	// channel, found := quicksilver.IBCKeeper.ChannelKeeper.GetChannel(suite.chainA.GetContext(), portID, channelID)
	// suite.Require().True(found)
	// fmt.Printf("DEBUG: channel >>>\n%v\n<<<\n", channel)

	quicksilver.IBCKeeper.ChannelKeeper.SetNextSequenceSend(s.chainA.GetContext(), portID, channelID, 1)
	quicksilver.ICAControllerKeeper.SetActiveChannelID(s.chainA.GetContext(), connectionID, portID, channelID)
	key, err := quicksilver.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		s.chainA.GetContext(),
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}
	err = quicksilver.GetScopedIBCKeeper().ClaimCapability(
		s.chainA.GetContext(),
		key,
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}

	key, err = quicksilver.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		s.chainA.GetContext(),
		host.PortPath(portID),
	)
	if err != nil {
		return err
	}
	err = quicksilver.GetScopedIBCKeeper().ClaimCapability(
		s.chainA.GetContext(),
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
	quicksilver.ICAControllerKeeper.SetInterchainAccountAddress(s.chainA.GetContext(), connectionID, portID, addr)
	return ibcModule.OnChanOpenAck(s.chainA.GetContext(), portID, channelID, "", "")
}

func (s *KeeperTestSuite) setupTestProtocolData() {
	// connection type for ibc testsuite chainB
	s.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", s.path.EndpointB.ConnectionID, s.chainB.ChainID, 0),
		s.chainB.ChainID,
	)
	// osmosis params
	s.addProtocolData(
		types.ProtocolDataTypeOsmosisParams,
		fmt.Sprintf("{\"ChainID\": %q}", "osmosis-1"),
		types.OsmosisParamsKey,
	)
	// osmosis test chain
	s.addProtocolData(
		types.ProtocolDataTypeConnection,
		fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", "connection-77002", "osmosis-1", 0),
		"osmosis-1",
	)
	// osmosis test pool
	s.addProtocolData(
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
	s.addProtocolData(
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
	// atom (cosmoshub) on local chain
	s.addProtocolData(types.ProtocolDataTypeLiquidToken,
		fmt.Sprintf(
			"{\"chainid\":%q,\"registeredzonechainid\":%q,\"ibcdenom\":%q,\"qassetdenom\":%q}",
			"testchain1",
			"cosmoshub-4",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"uqatom",
		),
		"testchain1/ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3")
}

func (s *KeeperTestSuite) addProtocolData(dataType types.ProtocolDataType, data, key string) {
	s.T().Helper()

	pd := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(dataType)],
		Data: []byte(data),
	}

	s.GetQuicksilverApp(s.chainA).ParticipationRewardsKeeper.SetProtocolData(s.chainA.GetContext(), key, &pd)
}

func (s *KeeperTestSuite) setupTestDeposits() {
	quicksilver := s.GetQuicksilverApp(s.chainA)

	// add deposit to chainB zone
	zone, found := quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().True(found)

	s.addReceipt(
		&zone,
		testAddress,
		"testTxHash03",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(150000000))),
	)

	// add deposit to cosmos zone
	zone, found = quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), "cosmoshub-4")
	s.Require().True(found)

	s.addReceipt(
		&zone,
		testAddress,
		"testTxHash01",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(120000000))),
	)

	// add deposit to osmosis zone
	zone, found = quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), "osmosis-1")
	s.Require().True(found)

	s.addReceipt(
		&zone,
		testAddress,
		"testTxHash02",
		sdk.NewCoins(sdk.NewCoin("uosmo", math.NewInt(100000000))),
	)
}

func (s *KeeperTestSuite) addReceipt(zone *icstypes.Zone, sender, hash string, coins sdk.Coins) {
	receipt := icstypes.Receipt{
		ChainId: zone.ChainId,
		Sender:  sender,
		Txhash:  hash,
		Amount:  coins,
	}

	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetReceipt(s.chainA.GetContext(), receipt)

	delegationAddress := utils.GenerateAccAddressForTestWithPrefix("cosmos")
	validatorAddress := utils.GenerateValAddressForTestWithPrefix("cosmos")
	delegation := icstypes.Delegation{
		DelegationAddress: delegationAddress,
		ValidatorAddress:  validatorAddress,
		Amount:            coins[0],
		Height:            1,
		RedelegationEnd:   101,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetDelegation(s.chainA.GetContext(), zone, delegation)
}

func (s *KeeperTestSuite) setupTestIntents() {
	quicksilver := s.GetQuicksilverApp(s.chainA)

	// chainB
	zone, found := quicksilver.InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
	s.Require().True(found)

	s.addIntent(
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

func (s *KeeperTestSuite) addIntent(address string, zone icstypes.Zone, intents icstypes.ValidatorIntents) {
	intent := icstypes.DelegatorIntent{
		Delegator: address,
		Intents:   intents,
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetDelegatorIntent(s.chainA.GetContext(), &zone, intent, false)
}

func (s *KeeperTestSuite) setupTestClaims() {
	// add some claims
	s.addClaim(
		testAddress,
		"cosmoshub-4",
		cmtypes.ClaimTypeLiquidToken,
		"osmosis-1",
		40000000,
	)

	s.addClaim(
		"quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure",
		"cosmoshub-4",
		cmtypes.ClaimTypeLiquidToken,
		"osmosis-1",
		0,
	)
}

func (s *KeeperTestSuite) addClaim(address, chainID string, claimType cmtypes.ClaimType, sourceChainID string, amount uint64) {
	claim := cmtypes.Claim{
		UserAddress:   address,
		ChainId:       chainID,
		Module:        claimType,
		SourceChainId: sourceChainID,
		Amount:        amount,
	}
	s.GetQuicksilverApp(s.chainA).ClaimsManagerKeeper.SetClaim(s.chainA.GetContext(), &claim)
}

/*func (suite *KeeperTestSuite) testReopenChannel(zone ) {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	// connection
	connectionID, err := quicksilver.InterchainstakingKeeper.GetConnectionForPort(suite.chainA.GetContext(), zone.PerformanceAddress.PortName)
	suite.Require().NoError(err)
	fmt.Printf("DEBUG: connectionID %q\n", connectionID)

	// channel
	channelID, found := quicksilver.ICAControllerKeeper.GetActiveChannelID(suite.chainA.GetContext(), connectionID, zone.PerformanceAddress.PortName)
	suite.Require().True(found)
	fmt.Printf("DEBUG: channelID %q\n", channelID)
	channel, found := quicksilver.IBCKeeper.ChannelKeeper.GetChannel(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID)
	suite.Require().True(found)
	fmt.Printf("Channel: %v\n", channel)

	// close channel
	channelCap := suite.chainA.GetChannelCapability(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID)
	err = quicksilver.IBCKeeper.ChannelKeeper.ChanCloseInit(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID, channelCap)
	suite.Require().True(found)

	// quicksilver.IBCKeeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID, channeltypes.Channel{State: channeltypes.CLOSED, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: "icahost", ChannelId: channelID}, ConnectionHops: []string{connectionID}})
	channel, found = quicksilver.IBCKeeper.ChannelKeeper.GetChannel(suite.chainA.GetContext(), zone.PerformanceAddress.PortName, channelID)
	suite.Require().True(found)
	fmt.Printf("Channel: %v\n", channel)

	// // attempt to reopen channel here
	// conn, found := quicksilver.IBCKeeper.ConnectionKeeper.GetConnection(suite.chainA.GetContext(), connectionID)
	// suite.Require().True(found)
	// fmt.Printf("Connection: %v\n", conn)

	// fmt.Println("Channel Closed...")
	// msg := channeltypes.NewMsgChannelOpenInit(zone.PerformanceAddress.PortName, icatypes.Version, channeltypes.ORDERED, []string{connectionID}, icatypes.PortID, icatypes.ModuleName)
	// handler := quicksilver.MsgServiceRouter().Handler(msg)
	// _, err = handler(suite.chainA.GetContext(), msg)
	// suite.Require().NoError(err)
}*/
