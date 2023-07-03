package keeper_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	umeetypes "github.com/ingenuity-build/quicksilver/umee/leverage/types"

	testsuite "github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	tmclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	epochtypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	ics "github.com/ingenuity-build/quicksilver/x/interchainstaking"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

var (
	testAddress        = addressutils.GenerateAddressForTestWithPrefix("cosmos")
	umeeTestConnection = "connection-77003"
	umeeTestChain      = "umee-1"
	umeeBaseDenom      = "uumee"
)

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

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	path *ibctesting.Path
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

	suite.coordinator.CurrentTime = time.Now().UTC()
	suite.coordinator.UpdateTime()

	suite.coreTest()
}

func (suite *KeeperTestSuite) coreTest() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	suite.setupTestZones()

	// test ProtocolData
	akpd := quicksilver.ParticipationRewardsKeeper.AllKeyedProtocolDatas(suite.chainA.GetContext())
	// initially we expect one - the 'local' chain
	suite.Require().Equal(1, len(akpd))

	suite.setupTestProtocolData()

	akpd = quicksilver.ParticipationRewardsKeeper.AllKeyedProtocolDatas(suite.chainA.GetContext())
	// added 6 in setupTestProtocolData
	suite.Require().Equal(14, len(akpd))

	// advance the chains
	suite.coordinator.CommitNBlocks(suite.chainA, 1)
	suite.coordinator.CommitNBlocks(suite.chainB, 1)

	// callback test
	suite.executeSetEpochBlockCallback()
	suite.executeOsmosisPoolUpdateCallback()
	suite.executeUmeeReservesUpdateCallback()
	suite.executeUmeeTotalBorrowsUpdateCallback()
	suite.executeUmeeInterestScalarUpdateCallback()
	suite.executeUmeeLeverageModuleBalanceUpdateCallback()
	suite.executeUmeeUTokenSupplyUpdateCallback()

	suite.setupTestDeposits()
	suite.setupTestIntents()

	quicksilver.EpochsKeeper.AfterEpochEnd(suite.chainA.GetContext(), epochtypes.EpochIdentifierEpoch, 1)

	suite.setupTestClaims()

	quicksilver.EpochsKeeper.AfterEpochEnd(suite.chainA.GetContext(), epochtypes.EpochIdentifierEpoch, 2)
	// Epoch boundary
	ctx := suite.chainA.GetContext()

	quicksilver.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		suite.Require().NoError(quicksilver.BankKeeper.MintCoins(ctx, "mint", sdk.NewCoins(sdk.NewCoin(quicksilver.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
		suite.Require().NoError(quicksilver.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, sdk.NewCoins(sdk.NewCoin(quicksilver.StakingKeeper.BondDenom(ctx), sdk.NewIntFromUint64(zone.HoldingsAllocation)))))
		return false
	})

	_, found := quicksilver.ClaimsManagerKeeper.GetLastEpochClaim(ctx, "cosmoshub-4", "quick16pxh2v4hr28h2gkntgfk8qgh47pfmjfhzgeure", cmtypes.ClaimTypeLiquidToken, "osmosis-1")
	suite.Require().True(found)

	quicksilver.EpochsKeeper.AfterEpochEnd(suite.chainA.GetContext(), epochtypes.EpochIdentifierEpoch, 3)

	// zone for remote chain
	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	valRewards := make(map[string]sdk.Dec)
	for _, val := range quicksilver.InterchainstakingKeeper.GetValidators(suite.chainA.GetContext(), suite.chainB.ChainID) {
		valRewards[val.ValoperAddress] = sdk.NewDec(100000000)
	}

	suite.executeValidatorSelectionRewardsCallback(zone.PerformanceAddress.Address, valRewards)
}

func (suite *KeeperTestSuite) setupTestZones() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	withdrawalAddress1 := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	withdrawalAddress2 := addressutils.GenerateAddressForTestWithPrefix("osmo")

	// test zone
	testzone := icstypes.Zone{
		ConnectionId:     suite.path.EndpointA.ConnectionID,
		ChainId:          suite.chainB.ChainID,
		AccountPrefix:    "cosmos",
		LocalDenom:       "uqatom",
		BaseDenom:        "uatom",
		ReturnToSender:   false,
		LiquidityModule:  true,
		DepositsEnabled:  true,
		UnbondingEnabled: false,
		Is_118:           true,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           withdrawalAddress1,
			PortName:          suite.chainB.ChainID + ".withrawal",
			WithdrawalAddress: withdrawalAddress1,
		},
	}
	selftestzone := icstypes.Zone{
		ConnectionId:     suite.path.EndpointB.ConnectionID,
		ChainId:          suite.chainA.ChainID,
		AccountPrefix:    "osmo",
		LocalDenom:       "uqosmo",
		BaseDenom:        "uosmo",
		ReturnToSender:   false,
		LiquidityModule:  true,
		DepositsEnabled:  true,
		UnbondingEnabled: false,
		Is_118:           true,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           withdrawalAddress2,
			PortName:          suite.chainA.ChainID + ".withrawal",
			WithdrawalAddress: withdrawalAddress2,
		},
	}

	quicksilver.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &selftestzone)
	quicksilver.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &testzone)

	quicksilver.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), "07-tendermint-0", &tmclienttypes.ClientState{ChainId: suite.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})

	quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(suite.chainA.GetContext(), "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &tmclienttypes.ConsensusState{Timestamp: suite.chainA.GetContext().BlockTime()})
	suite.Require().NoError(suite.setupChannelForICA(suite.chainB.ChainID, suite.path.EndpointA.ConnectionID, "performance", testzone.AccountPrefix))

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext())
	zone, found := quicksilver.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)

	for i := range vals {
		suite.Require().NoError(quicksilver.InterchainstakingKeeper.SetValidatorForZone(suite.chainA.GetContext(), &zone, app.DefaultConfig().Codec.MustMarshal(&vals[i])))
	}

	// self zone
	performanceAddressOsmo := addressutils.GenerateAddressForTestWithPrefix("osmo")
	performanceAccountOsmo, err := icstypes.NewICAAccount(performanceAddressOsmo, "testchain1.performance")
	suite.Require().NoError(err)
	withdrawalAddressOsmo := addressutils.GenerateAddressForTestWithPrefix("osmo")
	withdrawalAccountOsmo, err := icstypes.NewICAAccount(withdrawalAddressOsmo, "testchain1.withdrawal")
	suite.Require().NoError(err)
	performanceAccountOsmo.WithdrawalAddress = withdrawalAddressOsmo

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
		Is_118:             true,
		Decimals:           6,
		PerformanceAddress: performanceAccountOsmo,
		WithdrawalAddress:  withdrawalAccountOsmo,
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
	quicksilver.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneSelf)

	// cosmos zone
	performanceAddressCosmos := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	performanceAccountCosmos, err := icstypes.NewICAAccount(performanceAddressCosmos, "cosmoshub-4.performance")
	suite.Require().NoError(err)
	performanceAccountCosmos.WithdrawalAddress = addressutils.GenerateAddressForTestWithPrefix("cosmos")

	withdrawalAddressCosmos := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	withdrawalAccountCosmos, err := icstypes.NewICAAccount(withdrawalAddressCosmos, "cosmoshub-4.withdrawal")
	suite.Require().NoError(err)
	performanceAccountOsmo.WithdrawalAddress = withdrawalAddressCosmos

	zoneCosmos := icstypes.Zone{
		ConnectionId:       "connection-77001",
		ChainId:            "cosmoshub-4",
		AccountPrefix:      "cosmos",
		LocalDenom:         "uqatom",
		BaseDenom:          "uatom",
		ReturnToSender:     false,
		LiquidityModule:    true,
		PerformanceAddress: performanceAccountCosmos,
		Is_118:             true,
		WithdrawalAddress:  withdrawalAccountCosmos,
	}
	quicksilver.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneCosmos)
	cosmosVals := []icstypes.Validator{
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
	}
	for _, cosmosVal := range cosmosVals {
		quicksilver.InterchainstakingKeeper.SetValidator(suite.chainA.GetContext(), zoneCosmos.ChainId, cosmosVal)
	}

	withdrawalAddress := addressutils.GenerateAddressForTestWithPrefix("osmo")

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
			Address:           addressutils.GenerateAddressForTestWithPrefix("osmo"),
			PortName:          "osmosis-1.performance",
			WithdrawalAddress: withdrawalAddress,
		},
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           withdrawalAddress,
			PortName:          "osmosis-1.withrawal",
			WithdrawalAddress: withdrawalAddress,
		},
		Is_118: true,
	}
	quicksilver.InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zoneOsmosis)
}

func (suite *KeeperTestSuite) setupChannelForICA(chainID, connectionID, accountSuffix, remotePrefix string) error {
	suite.T().Helper()
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	ibcModule := ics.NewIBCModule(quicksilver.InterchainstakingKeeper)
	portID, err := icatypes.NewControllerPortID(chainID + "." + accountSuffix)
	if err != nil {
		return err
	}

	quicksilver.InterchainstakingKeeper.SetConnectionForPort(suite.chainA.GetContext(), connectionID, portID)

	channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(suite.chainA.GetContext())
	quicksilver.IBCKeeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), portID, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icatypes.PortID, ChannelId: channelID}, ConnectionHops: []string{connectionID}})

	quicksilver.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.chainA.GetContext(), portID, channelID, 1)
	quicksilver.ICAControllerKeeper.SetActiveChannelID(suite.chainA.GetContext(), connectionID, portID, channelID)
	key, err := quicksilver.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		suite.chainA.GetContext(),
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}
	err = quicksilver.GetScopedIBCKeeper().ClaimCapability(
		suite.chainA.GetContext(),
		key,
		host.ChannelCapabilityPath(portID, channelID),
	)
	if err != nil {
		return err
	}

	key, err = quicksilver.InterchainstakingKeeper.ScopedKeeper().NewCapability(
		suite.chainA.GetContext(),
		host.PortPath(portID),
	)
	if err != nil {
		return err
	}
	err = quicksilver.GetScopedIBCKeeper().ClaimCapability(
		suite.chainA.GetContext(),
		key,
		host.PortPath(portID),
	)
	if err != nil {
		return err
	}

	addr := addressutils.GenerateAddressForTestWithPrefix(remotePrefix)
	quicksilver.ICAControllerKeeper.SetInterchainAccountAddress(suite.chainA.GetContext(), connectionID, portID, addr)
	return ibcModule.OnChanOpenAck(suite.chainA.GetContext(), portID, channelID, "", "")
}

func (suite *KeeperTestSuite) setupTestProtocolData() {
	// connection type for ibc testsuite chainB
	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		[]byte(fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, 0)),
	)
	// umee params
	suite.addProtocolData(
		types.ProtocolDataTypeUmeeParams,
		[]byte(fmt.Sprintf("{\"ChainID\": %q}", umeeTestChain)),
	)
	// umee test chain
	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		[]byte(fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", umeeTestConnection, umeeTestChain, 0)),
	)
	// umee test reserves
	upd, _ := json.Marshal(types.UmeeReservesProtocolData{UmeeProtocolData: types.UmeeProtocolData{Denom: umeeBaseDenom}})
	suite.addProtocolData(
		types.ProtocolDataTypeUmeeReserves,
		upd,
	)
	// umee test leverage module balance
	upd, _ = json.Marshal(types.UmeeLeverageModuleBalanceProtocolData{UmeeProtocolData: types.UmeeProtocolData{Denom: umeeBaseDenom}})
	suite.addProtocolData(
		types.ProtocolDataTypeUmeeLeverageModuleBalance,
		upd,
	)
	// umee test borrows
	upd, _ = json.Marshal(types.UmeeTotalBorrowsProtocolData{UmeeProtocolData: types.UmeeProtocolData{Denom: umeeBaseDenom}})
	suite.addProtocolData(
		types.ProtocolDataTypeUmeeTotalBorrows,
		upd,
	)
	// umee test interest scalar
	upd, _ = json.Marshal(types.UmeeInterestScalarProtocolData{UmeeProtocolData: types.UmeeProtocolData{Denom: umeeBaseDenom}})
	suite.addProtocolData(
		types.ProtocolDataTypeUmeeInterestScalar,
		upd,
	)
	// umee test utoken supply
	upd, _ = json.Marshal(types.UmeeInterestScalarProtocolData{UmeeProtocolData: types.UmeeProtocolData{Denom: umeetypes.UTokenPrefix + umeeBaseDenom}})
	suite.addProtocolData(
		types.ProtocolDataTypeUmeeUTokenSupply,
		upd,
	)
	// osmosis params
	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisParams,
		[]byte(fmt.Sprintf("{\"ChainID\": %q, \"BaseDenom\": %q, \"BaseChain\": %q}", "osmosis-1", "uosmo", "osmosis-1")),
	)
	// osmosis test chain
	suite.addProtocolData(
		types.ProtocolDataTypeConnection,
		[]byte(fmt.Sprintf("{\"connectionid\": %q,\"chainid\": %q,\"lastepoch\": %d}", "connection-77002", "osmosis-1", 0)),
	)
	// osmosis test pool
	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisPool,
		[]byte(fmt.Sprintf(
			"{\"poolid\":%d,\"poolname\":%q,\"pooltype\":\"balancer\",\"denoms\":{%q:{\"chainid\": %q, \"denom\":%q}, %q:{\"chainid\": %q, \"denom\":%q}}}",
			1,
			"atom/osmo",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"cosmoshub-4",
			"uatom",
			"ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159",
			"osmosis-1",
			"uosmo",
		)),
	)

	// atom (cosmoshub) on osmosis
	suite.addProtocolData(
		types.ProtocolDataTypeLiquidToken,
		[]byte(fmt.Sprintf(
			"{\"chainid\":%q,\"registeredzonechainid\":%q,\"ibcdenom\":%q,\"qassetdenom\":%q}",
			"osmosis-1",
			"cosmoshub-4",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"uqatom",
		)),
	)
	// atom (cosmoshub) on local chain
	suite.addProtocolData(types.ProtocolDataTypeLiquidToken,
		[]byte(fmt.Sprintf(
			"{\"chainid\":%q,\"registeredzonechainid\":%q,\"ibcdenom\":%q,\"qassetdenom\":%q}",
			"testchain1",
			"cosmoshub-4",
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"uqatom",
		)),
	)
}

func (suite *KeeperTestSuite) addProtocolData(dataType types.ProtocolDataType, data []byte) {
	suite.T().Helper()

	pd := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(dataType)],
		Data: data,
	}

	upd, err := types.UnmarshalProtocolData(dataType, pd.Data)
	if err != nil {
		panic(err)
	}

	suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper.SetProtocolData(suite.chainA.GetContext(), upd.GenerateKey(), &pd)
}

func (suite *KeeperTestSuite) setupTestDeposits() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	// add deposit to chainB zone
	zone, found := quicksilver.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)

	suite.addReceipt(
		&zone,
		testAddress,
		"testTxHash03",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(150000000))),
	)

	// add deposit to cosmos zone
	zone, found = quicksilver.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), "cosmoshub-4")
	suite.Require().True(found)

	suite.addReceipt(
		&zone,
		testAddress,
		"testTxHash01",
		sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(120000000))),
	)

	// add deposit to osmosis zone
	zone, found = quicksilver.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), "osmosis-1")
	suite.Require().True(found)

	suite.addReceipt(
		&zone,
		testAddress,
		"testTxHash02",
		sdk.NewCoins(sdk.NewCoin("uosmo", math.NewInt(100000000))),
	)
}

func (suite *KeeperTestSuite) addReceipt(zone *icstypes.Zone, sender, hash string, coins sdk.Coins) {
	t := time.Now().Add(-time.Hour)
	t2 := time.Now().Add(-5 * time.Minute)
	receipt := icstypes.Receipt{
		ChainId:   zone.ChainId,
		Sender:    sender,
		Txhash:    hash,
		Amount:    coins,
		FirstSeen: &t,
		Completed: &t2,
	}

	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetReceipt(suite.chainA.GetContext(), receipt)

	delegationAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	validatorAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")
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
	quicksilver := suite.GetQuicksilverApp(suite.chainA)

	// chainB
	zone, found := quicksilver.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.Require().True(found)
	vals := quicksilver.InterchainstakingKeeper.GetValidators(suite.chainA.GetContext(), suite.chainB.ChainID)

	suite.addIntent(
		testAddress,
		zone,
		icstypes.ValidatorIntents{
			{
				ValoperAddress: vals[0].ValoperAddress,
				Weight:         sdk.MustNewDecFromStr("0.3"),
			},
			{
				ValoperAddress: vals[1].ValoperAddress,
				Weight:         sdk.MustNewDecFromStr("0.4"),
			},
			{
				ValoperAddress: vals[2].ValoperAddress,
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
	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetDelegatorIntent(suite.chainA.GetContext(), &zone, intent, false)
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
		1000,
	)
}

func (suite *KeeperTestSuite) addClaim(address, chainID string, claimType cmtypes.ClaimType, sourceChainID string, amount uint64) {
	claim := cmtypes.Claim{
		UserAddress:   address,
		ChainId:       chainID,
		Module:        claimType,
		SourceChainId: sourceChainID,
		Amount:        amount,
	}
	suite.GetQuicksilverApp(suite.chainA).ClaimsManagerKeeper.SetClaim(suite.chainA.GetContext(), &claim)
}
