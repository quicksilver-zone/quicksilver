package keeper_test

import (
	"encoding/json"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
	icqkeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) executeOsmosisPoolUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	osm := &keeper.OsmosisModule{}
	qid := icqkeeper.GenerateQueryHash("connection-77002", "osmosis-1", "store/gamm/key", osm.GetKeyPrefixPools(1), types.ModuleName)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	var err error
	resp := []byte{10, 26, 47, 111, 115, 109, 111, 115, 105, 115, 46, 103, 97, 109, 109, 46, 118, 49, 98, 101, 116, 97, 49, 46, 80, 111, 111, 108, 18, 202, 2, 10, 63, 111, 115, 109, 111, 49, 109, 119, 48, 97, 99, 54, 114, 119, 108, 112, 53, 114, 56, 119, 97, 112, 119, 107, 51, 122, 115, 54, 103, 50, 57, 104, 56, 102, 99, 115, 99, 120, 113, 97, 107, 100, 122, 119, 57, 101, 109, 107, 110, 101, 54, 99, 56, 119, 106, 112, 57, 113, 48, 116, 51, 118, 56, 116, 16, 1, 26, 6, 10, 1, 48, 18, 1, 48, 34, 4, 49, 54, 56, 104, 42, 43, 10, 11, 103, 97, 109, 109, 47, 112, 111, 111, 108, 47, 49, 18, 28, 49, 48, 48, 48, 48, 48, 48, 50, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 48, 48, 50, 94, 10, 80, 10, 68, 105, 98, 99, 47, 49, 53, 69, 57, 67, 53, 67, 70, 53, 57, 54, 57, 48, 56, 48, 53, 51, 57, 68, 66, 51, 57, 53, 70, 65, 55, 68, 57, 67, 48, 56, 54, 56, 50, 54, 53, 50, 49, 55, 69, 70, 67, 53, 50, 56, 52, 51, 51, 54, 55, 49, 65, 65, 70, 57, 66, 49, 57, 49, 50, 68, 49, 53, 57, 18, 8, 49, 48, 48, 48, 48, 48, 48, 51, 18, 10, 49, 48, 55, 51, 55, 52, 49, 56, 50, 52, 50, 94, 10, 80, 10, 68, 105, 98, 99, 47, 51, 48, 50, 48, 57, 50, 50, 66, 55, 53, 55, 54, 70, 67, 55, 53, 66, 66, 69, 48, 53, 55, 65, 48, 50, 57, 48, 65, 57, 65, 69, 69, 70, 70, 52, 56, 57, 66, 66, 49, 49, 49, 51, 69, 54, 69, 51, 54, 53, 67, 69, 52, 55, 50, 68, 52, 66, 70, 66, 55, 70, 70, 65, 51, 18, 8, 49, 48, 48, 48, 48, 48, 48, 51, 18, 10, 49, 48, 55, 51, 55, 52, 49, 56, 50, 52, 58, 10, 50, 49, 52, 55, 52, 56, 51, 54, 52, 56}
	// respB64 := "Chovb3Ntb3Npcy5nYW1tLnYxYmV0YTEuUG9vbBLKAgo/b3NtbzFtdzBhYzZyd2xwNXI4d2Fwd2szenM2ZzI5aDhmY3NjeHFha2R6dzllbWtuZTZjOHdqcDlxMHQzdjh0EAEaBgoBMBIBMCIEMTY4aCorCgtnYW1tL3Bvb2wvMRIcMTAwMDAwMDI5OTk5OTk5OTk5OTk5OTk5OTkwMDJeClAKRGliYy8xNUU5QzVDRjU5NjkwODA1MzlEQjM5NUZBN0Q5QzA4NjgyNjUyMTdFRkM1Mjg0MzM2NzFBQUY5QjE5MTJEMTU5EggxMDAwMDAwMxIKMTA3Mzc0MTgyNDJeClAKRGliYy8zMDIwOTIyQjc1NzZGQzc1QkJFMDU3QTAyOTBBOUFFRUZGNDg5QkIxMTEzRTZFMzY1Q0U0NzJENEJGQjdGRkEzEggxMDAwMDAwMxIKMTA3Mzc0MTgyNDoKMjE0NzQ4MzY0OA=="
	// resp, err := base64.StdEncoding.DecodeString(respB64)
	// suite.Require().NoError(err)

	// setup for expected
	var pdi gamm.PoolI
	err = prk.GetCodec().UnmarshalInterface(resp, &pdi)
	suite.Require().NoError(err)
	expectedData, err := json.Marshal(pdi)
	suite.Require().NoError(err)

	err = keeper.OsmosisPoolUpdateCallback(
		prk,
		ctx,
		resp,
		query,
	)
	suite.Require().NoError(err)

	want := types.OsmosisPoolProtocolData{
		PoolID:      1,
		PoolName:    "atom/osmo",
		LastUpdated: ctx.BlockTime(),
		PoolData:    expectedData,
		PoolType:    "balancer",
		Zones: map[string]string{
			"cosmoshub-4": "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
			"osmosis-1":   "ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159",
		},
	}

	pd, found := prk.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisPool, "1")
	suite.Require().True(found)

	ioppd, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, pd.Data)
	suite.Require().NoError(err)
	oppd := ioppd.(types.OsmosisPoolProtocolData)
	suite.Require().Equal(want, oppd)
}

func (suite *KeeperTestSuite) executeValidatorSelectionRewardsCallback(performanceAddress string) {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: performanceAddress}
	bz := prk.GetCodec().MustMarshal(&rewardsQuery)

	qid := icqkeeper.GenerateQueryHash(
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
		bz,
		types.ModuleName,
	)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	respJSON := `{"rewards":[{"validator_address":"cosmosvaloper1q86m0zq0p52h4puw5pg5xgc3c5e2mq52y6mth0","reward":[{"denom":"uatom","amount":"519569.212352312102820680"}]},{"validator_address":"cosmosvaloper1jtjjyxtqk0fj85ud9cxk368gr8cjdsftvdt5jl","reward":[{"denom":"uatom","amount":"519569.719818287718073104"}]},{"validator_address":"cosmosvaloper1759teakrsvnx7rnur8ezc4qaq8669nhtgukm0x","reward":[{"denom":"uatom","amount":"519569.180489558573776457"}]}],"total":[{"denom":"uatom","amount":"1558708.112660158394670241"}]}`
	qdtrResp := distrtypes.QueryDelegationTotalRewardsResponse{}
	err := json.Unmarshal([]byte(respJSON), &qdtrResp)
	suite.Require().NoError(err)
	resp, err := qdtrResp.Marshal()

	err = keeper.ValidatorSelectionRewardsCallback(
		prk,
		ctx,
		resp,
		query,
	)
	suite.Require().NoError(err)
}
