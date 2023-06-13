package keeper_test

import (
	"time"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestRedelegationRecordSetGetIterate() {
	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	testValidatorOne := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")
	testValidatorTwo := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")

	suite.SetupTest()

	records := quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Require().Equal(0, len(records))

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

	suite.Require().Equal(1, len(records))

	recordFetched, found := quicksilver.InterchainstakingKeeper.GetRedelegationRecord(ctx, "cosmoshub-4", testValidatorOne, testValidatorTwo, 1)

	suite.Require().True(found)
	suite.Require().Equal(record, recordFetched)

	allRecords := quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Require().Equal(1, len(allRecords))
	allCosmosRecords := quicksilver.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "cosmoshub-4")
	suite.Require().Equal(1, len(allCosmosRecords))
	allOtherChainRecords := quicksilver.InterchainstakingKeeper.ZoneRedelegationRecords(ctx, "elgafar-1")
	suite.Require().Equal(0, len(allOtherChainRecords))

	quicksilver.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, "cosmoshub-4", testValidatorOne, testValidatorTwo, 1)

	allCosmosRecords = quicksilver.InterchainstakingKeeper.AllRedelegationRecords(ctx)
	suite.Require().Equal(0, len(allCosmosRecords))
}
