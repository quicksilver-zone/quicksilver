package app

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"

	"github.com/quicksilver-zone/quicksilver/app/upgrades"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prtypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func init() {
	ibctesting.DefaultTestingAppInit = SetupTestingApp
}

// TestKeeperTestSuite runs all the tests within this package.
func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func newQuicksilverPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

type AppTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	path *ibctesting.Path
}

func (*AppTestSuite) GetQuicksilverApp(chain *ibctesting.TestChain) *Quicksilver {
	app, ok := chain.App.(*Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return app
}

// SetupTest creates a coordinator with 2 test chains.
func (s *AppTestSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)         // initializes 2 test chains
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1)) // convenience and readability
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2)) // convenience and readability

	s.path = newQuicksilverPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(s.path)

	s.coordinator.CurrentTime = time.Now().UTC()
	s.coordinator.UpdateTime()
}

func addVestingAccount(ctx sdk.Context, ak *authkeeper.AccountKeeper, address string, numPeriods int64, periodLength int64, total int64) {
	start := int64(1704240000)
	duration := numPeriods * periodLength
	perPeriod := total / numPeriods
	dust := total - (perPeriod * numPeriods)

	periods := make(vestingtypes.Periods, 0, numPeriods)
	for i := numPeriods; i > 0; i-- {
		periods = append(periods, vestingtypes.Period{Length: periodLength, Amount: sdk.NewCoins(sdk.NewCoin("uqck", math.NewInt(perPeriod)))})
	}
	periods[0].Amount.Add(sdk.NewCoin("uqck", math.NewInt(dust)))
	vest := vestingtypes.NewPeriodicVestingAccountRaw(
		vestingtypes.NewBaseVestingAccount(
			authtypes.NewBaseAccountWithAddress(addressutils.MustAccAddressFromBech32(address, "")),
			sdk.NewCoins(sdk.NewCoin("uqck", math.NewInt(total))),
			start+duration,
		),
		start,
		periods,
	)
	ak.SetAccount(ctx, vest)
}

func (s *AppTestSuite) InitV150TestZones() {
	// zone to match prod
	zone := icstypes.Zone{ConnectionId: "connection-1", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-0", ChainId: "stargaze-1", AccountPrefix: "stars", LocalDenom: "uqstars", BaseDenom: "ustars"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-50", ChainId: "juno-1", AccountPrefix: "juno", LocalDenom: "uqjuno", BaseDenom: "ujuno"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-2", ChainId: "osmosis-1", AccountPrefix: "osmo", LocalDenom: "uqosmo", BaseDenom: "uosmo"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-9", ChainId: "regen-1", AccountPrefix: "regen", LocalDenom: "uqregen", BaseDenom: "uregen"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	zone = icstypes.Zone{ConnectionId: "connection-54", ChainId: "sommelier-3", AccountPrefix: "somm", LocalDenom: "uqsomm", BaseDenom: "usomm"}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l", 10, 864000, 5000000000)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f", 10, 864000, 5000000000)

	// set counterparty channels to match prod so we can assert denoms
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-0", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-124")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-1", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-467")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-2", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-522")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-52", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-65")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-49", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-53")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-101", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-59")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-86", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-272")})
	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-17", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-62")})
}

func (s *AppTestSuite) TestV010500UpgradeHandler() {
	s.InitV150TestZones()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()
	validators := app.StakingKeeper.GetAllValidators(ctx)
	// Setting up
	accountA := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l", ""))
	err := app.BankKeeper.SendCoins(ctx, s.chainA.SenderAccount.GetAddress(), accountA.GetAddress(), sdk.Coins{
		sdk.NewInt64Coin("stake", 100),
	})
	s.NoError(err)

	accountB := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f", ""))
	err = app.BankKeeper.SendCoins(ctx, s.chainA.SenderAccount.GetAddress(), accountB.GetAddress(), sdk.Coins{
		sdk.NewInt64Coin("stake", 300),
	})
	s.NoError(err)

	// Stake old account
	amountA, _ := sdk.NewIntFromString("100")
	_, err = app.StakingKeeper.Delegate(ctx, accountA.GetAddress(), amountA, stakingtypes.Unbonded, validators[0], true)
	s.NoError(err)

	amountB, _ := sdk.NewIntFromString("100")
	_, err = app.StakingKeeper.Delegate(ctx, accountB.GetAddress(), amountB, stakingtypes.Unbonded, validators[0], true)
	s.NoError(err)

	amountB, _ = sdk.NewIntFromString("200")
	_, err = app.StakingKeeper.Delegate(ctx, accountB.GetAddress(), amountB, stakingtypes.Unbonded, validators[1], true)
	s.NoError(err)

	// set withdrawal records
	zone, found := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.True(found)

	user1 := addressutils.GenerateAddressForTestWithPrefix("quick")
	user2 := addressutils.GenerateAddressForTestWithPrefix("quick")
	recipient1 := addressutils.GenerateAddressForTestWithPrefix("quick")
	recipient2 := addressutils.GenerateAddressForTestWithPrefix("quick")

	// queued
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:      zone.ChainId,
		Delegator:    user1,
		Recipient:    recipient1,
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(3000)),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 1),
		Acknowledged: false,
		Status:       icstypes.WithdrawStatusQueued,
		EpochNumber:  1,
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:      zone.ChainId,
		Delegator:    user1,
		Recipient:    recipient1,
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(3000)),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 2),
		Acknowledged: false,
		Status:       icstypes.WithdrawStatusQueued,
		EpochNumber:  1,
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:      zone.ChainId,
		Delegator:    user1,
		Recipient:    recipient1,
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(3000)),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 3),
		Acknowledged: false,
		Status:       icstypes.WithdrawStatusQueued,
		EpochNumber:  1,
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:      zone.ChainId,
		Delegator:    user1,
		Recipient:    recipient2,
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(4000)),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 4),
		Acknowledged: false,
		Status:       icstypes.WithdrawStatusQueued,
		EpochNumber:  1,
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:      zone.ChainId,
		Delegator:    user2,
		Recipient:    recipient2,
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(5000)),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 5),
		Acknowledged: false,
		Status:       icstypes.WithdrawStatusQueued,
		EpochNumber:  1,
	})
	s.NoError(err)

	// unbonding

	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: user1,
		Recipient: recipient1,
		Distribution: []*icstypes.Distribution{
			{
				Valoper: "cosmosvaloper100000000000000000000000000000000000",
				Amount:  math.NewInt(600),
			},
			{
				Valoper: "cosmosvaloper111111111111111111111111111111111111",
				Amount:  math.NewInt(400),
			},
			{
				Valoper: "cosmosvaloper122222222222222222222222222222222222",
				Amount:  math.NewInt(200),
			},
		},
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(1000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(1200))),
		Requeued:       true,
		Txhash:         fmt.Sprintf("%064d", 6),
		Acknowledged:   true,
		Status:         icstypes.WithdrawStatusUnbond,
		EpochNumber:    1,
		CompletionTime: ctx.BlockTime().Add(180 * time.Minute),
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: user1,
		Recipient: recipient1,
		Distribution: []*icstypes.Distribution{
			{
				Valoper: "cosmosvaloper100000000000000000000000000000000000",
				Amount:  math.NewInt(600),
			},
			{
				Valoper: "cosmosvaloper111111111111111111111111111111111111",
				Amount:  math.NewInt(800),
			},
			{
				Valoper: "cosmosvaloper122222222222222222222222222222222222",
				Amount:  math.NewInt(800),
			},
			{
				Valoper: "cosmosvaloper133333333333333333333333333333333333",
				Amount:  math.NewInt(800),
			},
		},
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(2000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(2400))),
		Requeued:       true,
		Txhash:         fmt.Sprintf("%064d", 7),
		Acknowledged:   true,
		Status:         icstypes.WithdrawStatusUnbond,
		EpochNumber:    1,
		CompletionTime: ctx.BlockTime().Add(182 * time.Minute),
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: user2,
		Recipient: recipient2,
		Distribution: []*icstypes.Distribution{
			{
				Valoper: "cosmosvaloper100000000000000000000000000000000000",
				Amount:  math.NewInt(600),
			},
			{
				Valoper: "cosmosvaloper122222222222222222222222222222222222",
				Amount:  math.NewInt(200),
			},
			{
				Valoper: "cosmosvaloper133333333333333333333333333333333333",
				Amount:  math.NewInt(800),
			},
			{
				Valoper: "cosmosvaloper144444444444444444444444444444444444",
				Amount:  math.NewInt(2000),
			},
		},
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(3000)),
		Amount:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(3600))),
		Requeued:       true,
		Txhash:         fmt.Sprintf("%064d", 8),
		Acknowledged:   true,
		Status:         icstypes.WithdrawStatusUnbond,
		EpochNumber:    1,
		CompletionTime: ctx.BlockTime().Add(182 * time.Minute),
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: user2,
		Recipient: recipient2,
		Distribution: []*icstypes.Distribution{
			{
				Valoper: "cosmosvaloper100000000000000000000000000000000000",
				Amount:  math.NewInt(800),
			},
			{
				Valoper: "cosmosvaloper122222222222222222222222222222222222",
				Amount:  math.NewInt(400),
			},
			{
				Valoper: "cosmosvaloper133333333333333333333333333333333333",
				Amount:  math.NewInt(1200),
			},
			{
				Valoper: "cosmosvaloper144444444444444444444444444444444444",
				Amount:  math.NewInt(2200),
			},
		},
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(4000)),
		Amount:       sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(4800))),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 9),
		Acknowledged: true,
		Status:       icstypes.WithdrawStatusUnbond,
		EpochNumber:  1,
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: user2,
		Recipient: recipient1,
		Distribution: []*icstypes.Distribution{
			{
				Valoper: "cosmosvaloper100000000000000000000000000000000000",
				Amount:  math.NewInt(1000),
			},
			{
				Valoper: "cosmosvaloper122222222222222222222222222222222222",
				Amount:  math.NewInt(1200),
			},
			{
				Valoper: "cosmosvaloper133333333333333333333333333333333333",
				Amount:  math.NewInt(1200),
			},
			{
				Valoper: "cosmosvaloper144444444444444444444444444444444444",
				Amount:  math.NewInt(1600),
			},
		},
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(5000)),
		Amount:       sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(6000))),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 10),
		Acknowledged: false,
		Status:       icstypes.WithdrawStatusUnbond,
		EpochNumber:  1,
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: user1,
		Recipient: recipient1,
		Distribution: []*icstypes.Distribution{
			{
				Valoper: "cosmosvaloper100000000000000000000000000000000000",
				Amount:  math.NewInt(1500),
			},
			{
				Valoper: "cosmosvaloper122222222222222222222222222222222222",
				Amount:  math.NewInt(1500),
			},
			{
				Valoper: "cosmosvaloper133333333333333333333333333333333333",
				Amount:  math.NewInt(1500),
			},
			{
				Valoper: "cosmosvaloper144444444444444444444444444444444444",
				Amount:  math.NewInt(1500),
			},
			{
				Valoper: "cosmosvaloper155555555555555555555555555555555555",
				Amount:  math.NewInt(1500),
			},
		},
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(6000)),
		Amount:       sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(7500))),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 11),
		Acknowledged: true,
		Status:       icstypes.WithdrawStatusUnbond,
		EpochNumber:  2,
	})
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(s.chainA.GetContext(), icstypes.WithdrawalRecord{
		ChainId:   zone.ChainId,
		Delegator: user2,
		Recipient: recipient2,
		Distribution: []*icstypes.Distribution{
			{
				Valoper: "cosmosvaloper100000000000000000000000000000000000",
				Amount:  math.NewInt(1750),
			},
			{
				Valoper: "cosmosvaloper122222222222222222222222222222222222",
				Amount:  math.NewInt(1750),
			},
			{
				Valoper: "cosmosvaloper133333333333333333333333333333333333",
				Amount:  math.NewInt(1750),
			},
			{
				Valoper: "cosmosvaloper144444444444444444444444444444444444",
				Amount:  math.NewInt(1750),
			},
			{
				Valoper: "cosmosvaloper155555555555555555555555555555555555",
				Amount:  math.NewInt(1750),
			},
		},
		BurnAmount:   sdk.NewCoin("uqatom", math.NewInt(7000)),
		Amount:       sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(8750))),
		Requeued:     true,
		Txhash:       fmt.Sprintf("%064d", 12),
		Acknowledged: true,
		Status:       icstypes.WithdrawStatusUnbond,
		EpochNumber:  2,
	})
	s.NoError(err)
	handler := upgrades.V010500UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)

	_, err = handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	_, found = app.StakingKeeper.GetDelegation(ctx, accountA.GetAddress(), validators[0].GetOperator())
	s.False(found)

	migratedA := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1h0sqndv2y4xty6uk0sv4vckgyc5aa7n5at7fll", ""))
	stakeBalanceA := app.BankKeeper.GetBalance(ctx, migratedA.GetAddress(), "stake")
	s.Equal(int64(100), stakeBalanceA.Amount.Int64())

	migratedB := app.AccountKeeper.GetAccount(ctx, addressutils.MustAccAddressFromBech32("quick1n4g6037cjm0e0v2nvwj2ngau7pk758wtwk6lwq", ""))
	stakeBalanceB := app.BankKeeper.GetBalance(ctx, migratedB.GetAddress(), "stake")
	s.Equal(int64(300), stakeBalanceB.Amount.Int64())

	// Check the vest period of new account
	vestMigratedA, ok := migratedA.(*vestingtypes.PeriodicVestingAccount)
	s.True(ok)
	s.Equal(int64(5000000000), vestMigratedA.OriginalVesting.AmountOf("uqck").Int64())
	s.Equal(float64(864000), vestMigratedA.VestingPeriods[0].Duration().Seconds())

	vestMigratedB, ok := migratedB.(*vestingtypes.PeriodicVestingAccount)
	s.True(ok)
	s.Equal(int64(5000000000), vestMigratedB.OriginalVesting.AmountOf("uqck").Int64())
	s.Equal(float64(864000), vestMigratedB.VestingPeriods[0].Duration().Seconds())

	z, existed := app.InterchainstakingKeeper.GetLocalDenomZoneMapping(ctx, "uqatom")
	s.True(existed)
	s.Equal(z.ChainId, "cosmoshub-4")
	s.Equal(z.ConnectionId, "connection-1")

	// 512 should be the sum of 01, 02, 03
	wdr, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, z.ChainId, fmt.Sprintf("%064d", 512), icstypes.WithdrawStatusQueued)
	s.True(found)
	s.Equal(wdr.BurnAmount, sdk.NewCoin("uqatom", math.NewInt(9000)))
	s.Equal(wdr.CompletionTime, time.Time{})
	s.True(wdr.Requeued)
	s.Nil(wdr.Amount)
	s.Nil(wdr.Distribution)

	// 513 and 514 should be 04 and 05 requeued respectively (due to differing recipient)
	wdr, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, z.ChainId, fmt.Sprintf("%064d", 513), icstypes.WithdrawStatusQueued)
	s.True(found)
	s.Equal(wdr.BurnAmount, sdk.NewCoin("uqatom", math.NewInt(4000)))
	s.Equal(wdr.CompletionTime, time.Time{})
	s.True(wdr.Requeued)
	s.Nil(wdr.Amount)
	s.Nil(wdr.Distribution)

	wdr, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, z.ChainId, fmt.Sprintf("%064d", 514), icstypes.WithdrawStatusQueued)
	s.True(found)
	s.Equal(wdr.BurnAmount, sdk.NewCoin("uqatom", math.NewInt(5000)))
	s.Equal(wdr.CompletionTime, time.Time{})
	s.True(wdr.Requeued)
	s.Nil(wdr.Amount)
	s.Nil(wdr.Distribution)

	// 010 shouldn't be touched; it was not acknowledged and will be requeued.
	wdr, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, z.ChainId, fmt.Sprintf("%064d", 10), icstypes.WithdrawStatusUnbond)
	s.True(found)
	s.Equal(wdr.BurnAmount, sdk.NewCoin("uqatom", math.NewInt(5000)))
	s.Equal(wdr.CompletionTime, time.Time{})
	s.True(wdr.Requeued)
	s.False(wdr.Acknowledged)
	s.Equal(wdr.Amount, sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(6000))))

	// 06 + 07 are a pair; distribution should be merged.
	wdr, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, z.ChainId, fmt.Sprintf("%064d", 515), icstypes.WithdrawStatusUnbond)
	s.True(found)
	s.Equal(wdr.BurnAmount, sdk.NewCoin("uqatom", math.NewInt(3000)))
	// use the latest completion time (+182 mins)
	s.Equal(wdr.CompletionTime, ctx.BlockTime().Add(182*time.Minute))
	s.True(wdr.Requeued)
	s.True(wdr.Acknowledged)
	s.Equal(wdr.Amount, sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(3600))))
	s.ElementsMatch(wdr.Distribution, []*icstypes.Distribution{{Valoper: "cosmosvaloper100000000000000000000000000000000000", Amount: math.NewInt(1200)}, {Valoper: "cosmosvaloper111111111111111111111111111111111111", Amount: math.NewInt(1200)}, {Valoper: "cosmosvaloper122222222222222222222222222222222222", Amount: math.NewInt(1000)}, {Valoper: "cosmosvaloper133333333333333333333333333333333333", Amount: math.NewInt(800)}})

	wdrs := app.InterchainstakingKeeper.AllWithdrawalRecords(ctx)
	s.Equal(35, len(wdrs)) // 8 from requeue collation, 27 new records from restituion

	// test protocol data

	tvs, err := app.ParticipationRewardsKeeper.CalcTokenValues(ctx)
	s.NoError(err)
	expectedTvs := map[string]sdk.Dec{ // relative prices between assets as of 2024-03-09T11:00
		"uosmo":  sdk.MustNewDecFromStr("1.000000000000000000"),
		"uatom":  sdk.MustNewDecFromStr("8.312793554467208145"),
		"ustars": sdk.MustNewDecFromStr("0.024508540336823926"),
		"uregen": sdk.MustNewDecFromStr("0.034894445954581256"),
		"usomm":  sdk.MustNewDecFromStr("0.108532538179923692"),
		"ujuno":  sdk.MustNewDecFromStr("0.256801530018076838"),
	}
	for denom, value := range expectedTvs {
		s.Equal(tvs[denom], value)
	}
}

func (s *AppTestSuite) TestV010501UpgradeHandler() {
	s.InitV150TestZones()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	handler := upgrades.V010501UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	osmosisQatom, found := app.ParticipationRewardsKeeper.GetProtocolData(ctx, prtypes.ProtocolDataTypeLiquidToken, "osmosis-1_ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC")
	s.True(found)
	prdata, err := prtypes.UnmarshalProtocolData(prtypes.ProtocolDataTypeLiquidToken, osmosisQatom.Data)
	s.NoError(err)
	lpd, ok := prdata.(*prtypes.LiquidAllowedDenomProtocolData)
	s.True(ok)
	s.Equal("osmosis-1", lpd.ChainID)
	s.Equal("cosmoshub-4", lpd.RegisteredZoneChainID)
	s.Equal("uqatom", lpd.QAssetDenom)
}

func (s *AppTestSuite) TestV010503UpgradeHandler() {
	s.InitV150TestZones()
	app := s.GetQuicksilverApp(s.chainA)
	ctx := s.chainA.GetContext()

	user1 := addressutils.GenerateAddressForTestWithPrefix("quick")
	recipient1 := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	val1 := addressutils.GenerateAddressForTestWithPrefix("cosmovaloper")
	val2 := addressutils.GenerateAddressForTestWithPrefix("cosmovaloper")

	wdr1 := icstypes.WithdrawalRecord{
		ChainId:    s.chainB.ChainID,
		BurnAmount: sdk.NewInt64Coin("uqatom", 300),
		Distribution: []*icstypes.Distribution{
			{Valoper: val1, XAmount: 110},
			{Valoper: val2, XAmount: 220},
		},
		Amount:         sdk.NewCoins(sdk.NewInt64Coin("uatom", 330)),
		Txhash:         fmt.Sprintf("%064d", 1),
		Status:         icstypes.WithdrawStatusUnbond,
		Delegator:      user1,
		Recipient:      recipient1,
		EpochNumber:    2,
		CompletionTime: ctx.BlockTime().Add(3 * 24 * time.Hour),
		Acknowledged:   true,
	}

	wdr2 := icstypes.WithdrawalRecord{
		ChainId:      s.chainB.ChainID,
		BurnAmount:   sdk.NewInt64Coin("uqatom", 300),
		Distribution: nil,
		Txhash:       fmt.Sprintf("%064d", 1),
		Status:       icstypes.WithdrawStatusQueued,
		Delegator:    user1,
		Recipient:    recipient1,
		EpochNumber:  2,
	}

	err := app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, wdr1)
	s.NoError(err)
	err = app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, wdr2)
	s.NoError(err)

	app.ClaimsManagerKeeper.SetClaim(ctx, &cmtypes.Claim{UserAddress: user1, ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "osmosis-1", XAmount: 3000})
	app.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &cmtypes.Claim{UserAddress: user1, ChainId: s.chainB.ChainID, Module: cmtypes.ClaimTypeLiquidToken, SourceChainId: "osmosis-1", XAmount: 2900})

	handler := upgrades.V010503UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)

	_, err = handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	wdr2actual, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, wdr2.ChainId, wdr2.Txhash, wdr2.Status)
	s.True(found)
	s.Equal(wdr2actual, wdr2)

	wdr1actual, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, wdr1.ChainId, wdr1.Txhash, wdr1.Status)
	s.True(found)
	s.Contains(wdr1actual.Distribution, &icstypes.Distribution{Valoper: val1, Amount: sdk.NewInt(110)})
	s.Contains(wdr1actual.Distribution, &icstypes.Distribution{Valoper: val2, Amount: sdk.NewInt(220)})

	claims := app.ClaimsManagerKeeper.AllZoneUserClaims(ctx, s.chainB.ChainID, user1)
	s.Equal(claims[0].Amount, math.NewInt(3000))

	leclaims := app.ClaimsManagerKeeper.AllZoneLastEpochUserClaims(ctx, s.chainB.ChainID, user1)
	s.Equal(leclaims[0].Amount, math.NewInt(2900))
}

// Init a zone with some zero burnAmount withdrawal records

func (s *AppTestSuite) InitV160rc0TestZone() {
	cosmosWithdrawal := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosPerformance := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDeposit := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDelegate := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	// cosmos zone
	zone := icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           cosmosWithdrawal,
			PortName:          "icacontroller-cosmoshub-4.withdrawal",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           cosmosDelegate,
			PortName:          "icacontroller-cosmoshub-4.delegate",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           cosmosDeposit,
			PortName:          "icacontroller-cosmoshub-4.deposit",
			WithdrawalAddress: cosmosWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           cosmosPerformance,
			PortName:          "icacontroller-cosmoshub-4.performance",
			WithdrawalAddress: cosmosWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
}

// UncheckedSetWithdrawalRecord store the withdrawal record without checking the burnAmount.
// WARNING: This function is intended for testing purposes only and should not be used in production code.
func (s *AppTestSuite) UncheckedSetWithdrawalRecord(ctx sdk.Context, app *Quicksilver, record icstypes.WithdrawalRecord) {
	key, err := hex.DecodeString(record.Txhash)
	if err != nil {
		panic(err)
	}

	store := prefix.NewStore(ctx.KVStore(app.GetKey(icstypes.StoreKey)), icstypes.GetWithdrawalKey(record.ChainId, record.Status))
	bz := app.InterchainstakingKeeper.GetCodec().MustMarshal(&record)
	store.Set(key, bz)
}

func (s *AppTestSuite) InitV160TestZones() {
	cosmosWithdrawal := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosPerformance := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDeposit := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	cosmosDelegate := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	// cosmos zone
	zone := icstypes.Zone{
		ConnectionId:    "connection-77001",
		ChainId:         "cosmoshub-4",
		AccountPrefix:   "cosmos",
		LocalDenom:      "uqatom",
		BaseDenom:       "uatom",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           cosmosWithdrawal,
			PortName:          "icacontroller-cosmoshub-4.withdrawal",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           cosmosDelegate,
			PortName:          "icacontroller-cosmoshub-4.delegate",
			WithdrawalAddress: cosmosWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           cosmosDeposit,
			PortName:          "icacontroller-cosmoshub-4.deposit",
			WithdrawalAddress: cosmosWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           cosmosPerformance,
			PortName:          "icacontroller-cosmoshub-4.performance",
			WithdrawalAddress: cosmosWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

	osmoWithdrawal := addressutils.GenerateAddressForTestWithPrefix("osmo")
	osmoPerformance := addressutils.GenerateAddressForTestWithPrefix("osmo")
	osmoDeposit := addressutils.GenerateAddressForTestWithPrefix("osmo")
	osmoDelegate := addressutils.GenerateAddressForTestWithPrefix("osmo")
	// osmosis zone
	zone = icstypes.Zone{
		ConnectionId:    "connection-77002",
		ChainId:         "osmosis-1",
		AccountPrefix:   "osmo",
		LocalDenom:      "uqosmo",
		BaseDenom:       "uosmo",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           osmoWithdrawal,
			PortName:          "icacontroller-osmosis-1.withdrawal",
			WithdrawalAddress: osmoWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           osmoDelegate,
			PortName:          "icacontroller-osmosis-1.delegate",
			WithdrawalAddress: osmoWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           osmoDeposit,
			PortName:          "icacontroller-osmosis-1.deposit",
			WithdrawalAddress: osmoWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           osmoPerformance,
			PortName:          "icacontroller-osmosis-1.performance",
			WithdrawalAddress: osmoWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	// uni-5 zone

	junoWithdrawal := addressutils.GenerateAddressForTestWithPrefix("juno")
	junoPerformance := addressutils.GenerateAddressForTestWithPrefix("juno")
	junoDeposit := addressutils.GenerateAddressForTestWithPrefix("juno")
	junoDelegate := addressutils.GenerateAddressForTestWithPrefix("juno")

	zone = icstypes.Zone{
		ConnectionId:    "connection-77003",
		ChainId:         "juno-1",
		AccountPrefix:   "juno",
		LocalDenom:      "uqjuno",
		BaseDenom:       "ujuno",
		MultiSend:       false,
		LiquidityModule: false,
		WithdrawalAddress: &icstypes.ICAAccount{
			Address:           junoWithdrawal,
			PortName:          "icacontroller-juno-1.withdrawal",
			WithdrawalAddress: junoWithdrawal,
		},
		DelegationAddress: &icstypes.ICAAccount{
			Address:           junoDelegate,
			PortName:          "icacontroller-juno-1.delegate",
			WithdrawalAddress: junoWithdrawal,
		},
		DepositAddress: &icstypes.ICAAccount{
			Address:           junoDeposit,
			PortName:          "icacontroller-juno-1.deposit",
			WithdrawalAddress: junoWithdrawal,
		},
		PerformanceAddress: &icstypes.ICAAccount{
			Address:           junoPerformance,
			PortName:          "icacontroller-juno-1.performance",
			WithdrawalAddress: junoWithdrawal,
		},
	}
	s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
	addVestingAccount(s.chainA.GetContext(), &s.GetQuicksilverApp(s.chainA).AccountKeeper, "quick1qfyntnmlvznvrkk9xqppmcxqcluv7wd74nmyus", 10, 864000, 5000000000)

	// set withdrawal records
	invalidWithdrawal := icstypes.WithdrawalRecord{
		ChainId:        zone.ChainId,
		Delegator:      cosmosDelegate,
		Recipient:      cosmosWithdrawal,
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(0)),
		Requeued:       true,
		Txhash:         fmt.Sprintf("%064d", 1),
		Acknowledged:   false,
		Status:         icstypes.WithdrawStatusQueued,
		EpochNumber:    1,
		CompletionTime: time.Time{},
	}
	s.UncheckedSetWithdrawalRecord(s.chainA.GetContext(), s.GetQuicksilverApp(s.chainA), invalidWithdrawal)
	validWithdrawal := icstypes.WithdrawalRecord{
		ChainId:        zone.ChainId,
		Delegator:      cosmosDelegate,
		Recipient:      cosmosWithdrawal,
		BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(1000)),
		Requeued:       true,
		Txhash:         fmt.Sprintf("%064d", 2),
		Acknowledged:   false,
		Status:         icstypes.WithdrawStatusQueued,
		EpochNumber:    1,
		CompletionTime: time.Time{},
	}
	s.UncheckedSetWithdrawalRecord(s.chainA.GetContext(), s.GetQuicksilverApp(s.chainA), validWithdrawal)

	s.GetQuicksilverApp(s.chainA).IBCKeeper.ChannelKeeper.SetChannel(s.chainA.GetContext(), "transfer", "channel-2", channeltypes.Channel{Counterparty: channeltypes.NewCounterparty("transfer", "channel-522")})
}

func (s *AppTestSuite) TestV010505UpgradeHandler() {
	s.InitV160TestZones()
	app := s.GetQuicksilverApp(s.chainA)

	ctx := s.chainA.GetContext()

	handler := upgrades.V010505UpgradeHandler(app.mm,
		app.configurator, &app.AppKeepers)

	_, err := handler(ctx, types.Plan{}, app.mm.GetVersionMap())
	s.NoError(err)

	osmoZone, ok := app.InterchainstakingKeeper.GetZone(ctx, "osmosis-1")
	s.True(ok)
	s.Equal(math.NewInt(1_000_000), osmoZone.DustThreshold)

	cosmosZone, ok := app.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
	s.True(ok)
	s.Equal(math.NewInt(500_000), cosmosZone.DustThreshold)

	junoZone, ok := app.InterchainstakingKeeper.GetZone(ctx, "juno-1")
	s.True(ok)
	s.Equal(math.NewInt(2_000_000), junoZone.DustThreshold)

	// check if the invalid withdrawal record is removed
	_, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "juno-1", fmt.Sprintf("%064d", 1), icstypes.WithdrawStatusQueued)
	s.False(found)
	// check if the valid withdrawal record is still there
	_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "juno-1", fmt.Sprintf("%064d", 2), icstypes.WithdrawStatusQueued)
	s.True(found)
}
