package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestRedelegationRecordSetGetIterate() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	testValidatorOne := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")
	testValidatorTwo := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")

	suite.SetupTest()

	records := quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(0, len(records))

	record := types.RedelegationRecord{
		ChainId:        "cosmoshub-4",
		EpochNumber:    1,
		Source:         testValidatorOne,
		Destination:    testValidatorTwo,
		Amount:         3000,
		CompletionTime: time.Now().Add(time.Hour).UTC(),
	}

	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	records = quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)

	suite.Equal(1, len(records))

	recordFetched, found := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, "cosmoshub-4", testValidatorOne, testValidatorTwo, 1)

	suite.True(found)
	suite.Equal(record, recordFetched)

	allRecords := quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(1, len(allRecords))
	allCosmosRecords := quicksilver.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "cosmoshub-4")
	suite.Equal(1, len(allCosmosRecords))
	allOtherChainRecords := quicksilver.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "elgafar-1")
	suite.Equal(0, len(allOtherChainRecords))

	quicksilver.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, "cosmoshub-4", testValidatorOne, testValidatorTwo, 1)

	allCosmosRecords = quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(0, len(allCosmosRecords))
}

func (suite *KeeperTestSuite) TestGCCompletedRedelegations() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	testValidatorOne := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")
	testValidatorTwo := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")
	testValidatorThree := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")

	suite.SetupTest()

	records := quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(0, len(records))

	currentTime := ctx.BlockTime()

	record := types.RedelegationRecord{
		ChainId:        "cosmoshub-4",
		EpochNumber:    1,
		Source:         testValidatorOne,
		Destination:    testValidatorTwo,
		Amount:         3000,
		CompletionTime: currentTime.Add(time.Hour).UTC(),
	}
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	record = types.RedelegationRecord{
		ChainId:        "cosmoshub-4",
		EpochNumber:    1,
		Source:         testValidatorOne,
		Destination:    testValidatorThree,
		Amount:         3000,
		CompletionTime: currentTime.Add(-time.Hour).UTC(),
	}
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)
	record = types.RedelegationRecord{
		ChainId:        "cosmoshub-4",
		EpochNumber:    1,
		Source:         testValidatorThree,
		Destination:    testValidatorTwo,
		Amount:         3000,
		CompletionTime: time.Time{},
	}
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	records = quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(3, len(records))

	err := quicksilver.InterchainstakingKeeper.GCCompletedRedelegations(ctx)
	suite.NoError(err)

	records = quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(2, len(records))

	_, found := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, "cosmoshub-4", testValidatorOne, testValidatorThree, 1)
	suite.False(found)
}

func (suite *KeeperTestSuite) TestDeleteRedelegationRecordByKey() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	testValidatorOne := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")
	testValidatorTwo := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")
	testValidatorThree := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")

	suite.SetupTest()

	// Currently there are 0 records
	records := quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(0, len(records))

	// Set 3 records
	currentTime := ctx.BlockTime()

	record := types.RedelegationRecord{
		ChainId:        "cosmoshub-4",
		EpochNumber:    1,
		Source:         testValidatorOne,
		Destination:    testValidatorTwo,
		Amount:         3000,
		CompletionTime: currentTime.Add(time.Hour).UTC(),
	}
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)

	record = types.RedelegationRecord{
		ChainId:        "cosmoshub-4",
		EpochNumber:    1,
		Source:         testValidatorOne,
		Destination:    testValidatorThree,
		Amount:         3000,
		CompletionTime: currentTime.Add(-time.Hour).UTC(),
	}
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)
	record = types.RedelegationRecord{
		ChainId:        "cosmoshub-4",
		EpochNumber:    1,
		Source:         testValidatorThree,
		Destination:    testValidatorTwo,
		Amount:         3000,
		CompletionTime: time.Time{},
	}
	quicksilver.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)
	// Check set 3 records
	records = quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(3, len(records))
	// Handle DeleteRedelegationRecordByKey for 3 records
	quicksilver.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(idx int64, key []byte, redelegation types.RedelegationRecord) bool {
		quicksilver.InterchainstakingKeeper.DeleteRedelegationRecordByKey(ctx, append(types.KeyPrefixRedelegationRecord, key...))
		return false
	})
	// Check DeleteRedelegationRecordByKey 3 records to 0 records
	records = quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Equal(0, len(records))
}

func (suite *KeeperTestSuite) TestGCCompletedUnbondings() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	if !found {
		suite.Fail("unable to retrieve zone for test")
	}
	records := quicksilver.InterchainstakingKeeper.AllWithdrawalRecords(ctx)
	suite.Equal(0, len(records))

	vals := quicksilver.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
	currentTime := ctx.BlockTime()

	record1 := types.WithdrawalRecord{
		ChainId:   suite.chainB.ChainID,
		Delegator: zone.DelegationAddress.Address,
		Distribution: []*types.Distribution{
			{
				Valoper: vals[0].ValoperAddress,
				Amount:  500,
			},
			{
				Valoper: vals[1].ValoperAddress,
				Amount:  500,
			},
		},
		Recipient:      user1.String(),
		Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(1000))),
		BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdkmath.NewInt(1000)),
		Txhash:         "1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F1191DF092BD",
		Status:         types.WithdrawStatusCompleted,
		CompletionTime: currentTime.Add(-25 * time.Hour).UTC(),
	}
	quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record1)

	record2 := types.WithdrawalRecord{
		ChainId:   suite.chainB.ChainID,
		Delegator: zone.DelegationAddress.Address,
		Distribution: []*types.Distribution{
			{
				Valoper: vals[0].ValoperAddress,
				Amount:  500,
			},
			{
				Valoper: vals[1].ValoperAddress,
				Amount:  500,
			},
		},
		Recipient:      user2.String(),
		Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(1000))),
		BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdkmath.NewInt(1000)),
		Txhash:         "91DF093BD1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F11",
		Status:         types.WithdrawStatusUnbond,
		CompletionTime: currentTime.Add(25 * time.Hour).UTC(),
	}
	quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record2)

	record3 := types.WithdrawalRecord{
		ChainId:   suite.chainB.ChainID,
		Delegator: zone.DelegationAddress.Address,
		Distribution: []*types.Distribution{
			{
				Valoper: vals[0].ValoperAddress,
				Amount:  500,
			},
			{
				Valoper: vals[1].ValoperAddress,
				Amount:  500,
			},
		},
		Recipient:      user2.String(),
		Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdkmath.NewInt(1000))),
		BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdkmath.NewInt(1000)),
		Txhash:         "2247B55EE89FB22FC68C62D61050B944F1191DF093BD1613D2E8FBF7C7294A4D",
		Status:         types.WithdrawStatusUnbond,
		CompletionTime: time.Time{},
	}
	quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record3)

	records = quicksilver.InterchainstakingKeeper.AllWithdrawalRecords(ctx)
	suite.Equal(3, len(records))

	err := quicksilver.InterchainstakingKeeper.GCCompletedUnbondings(ctx, &zone)
	suite.NoError(err)

	records = quicksilver.InterchainstakingKeeper.AllWithdrawalRecords(ctx)
	suite.Equal(2, len(records))

	_, found = quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, record1.ChainId, record1.Txhash, record1.Status)
	suite.False(found)
}
