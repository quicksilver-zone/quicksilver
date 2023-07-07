package keeper_test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
	umeetypes "github.com/ingenuity-build/quicksilver/umee-types/leverage/types"
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
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)

	want := &types.OsmosisPoolProtocolData{
		PoolID:      1,
		PoolName:    "atom/osmo",
		LastUpdated: ctx.BlockTime(),
		PoolData:    expectedData,
		PoolType:    "balancer",
		Denoms: map[string]types.DenomWithZone{
			"ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3": {ChainID: "cosmoshub-4", Denom: "uatom"},
			"ibc/15E9C5CF5969080539DB395FA7D9C0868265217EFC528433671AAF9B1912D159": {ChainID: "osmosis-1", Denom: "uosmo"},
		},
	}

	pd, found := prk.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisPool, "1")
	suite.Require().True(found)

	ioppd, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, pd.Data)
	suite.Require().NoError(err)
	oppd := ioppd.(*types.OsmosisPoolProtocolData)
	suite.Require().Equal(want, oppd)
}

func (suite *KeeperTestSuite) executeValidatorSelectionRewardsCallback(performanceAddress string, valRewards map[string]sdk.Dec) {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: performanceAddress}
	bz := prk.GetCodec().MustMarshal(&rewardsQuery)

	qid := icqkeeper.GenerateQueryHash(
		suite.path.EndpointA.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
		bz,
		types.ModuleName,
	)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	var respJSON strings.Builder
	respJSON.Write([]byte(`{"rewards":[`))
	total := sdk.ZeroDec()
	i := 0
	for val, amount := range valRewards {
		if i > 0 {
			respJSON.Write([]byte(","))
		}
		respJSON.Write([]byte(fmt.Sprintf(`{"validator_address":%q,"reward":[{"denom":"uatom","amount":%q}]}`, val, amount.String())))
		total = total.Add(amount)
		i++
	}
	respJSON.Write([]byte(fmt.Sprintf(`],"total":[{"denom":"uatom","amount":%q}]}`, total.String())))
	qdtrResp := distrtypes.QueryDelegationTotalRewardsResponse{}
	err := json.Unmarshal([]byte(respJSON.String()), &qdtrResp)
	suite.Require().NoError(err)
	resp, err := qdtrResp.Marshal()
	suite.Require().NoError(err)

	err = keeper.ValidatorSelectionRewardsCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) executeSetEpochBlockCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	blockQuery := tmservice.GetLatestBlockRequest{}
	bz := prk.GetCodec().MustMarshal(&blockQuery)

	qid := icqkeeper.GenerateQueryHash(
		suite.path.EndpointA.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.base.tendermint.v1beta1.Service/GetLatestBlock",
		bz,
		types.ModuleName,
	)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)
	respJSON := `{"block_id":{"hash":"74pkkjg7u1eLtXxFCinnCmln3aVZVOqLCT3OnE3D+VA=","part_set_header":{"total":1,"hash":"UiLM70PpplmmZ85qC0ZKva5kYJmSZ2TTEZ4a7g9G92Q="}},"block":{"header":{"version":{"block":"11","app":"0"},"chain_id":"quickgaia-1","height":"90767","time":"2022-11-03T09:12:23.109926769Z","last_block_id":{"hash":"wCK5QmPuGJiRpn06Xu7ZjhxHwzBXVZrGqngMMzeRq8w=","part_set_header":{"total":1,"hash":"xYxXM7rX6Qcq/Yx3MpZeQA+FCeUbKSVr/FEuzfAFFQk="}},"last_commit_hash":"1Ev2iL1pTgyItBtSFbCRzxwdCtJfaCC1P+zWaDkJ/nU=","data_hash":"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=","validators_hash":"kQ9NNQ26Q3l5aXF2IwraaweoLzusIVDUA53AycOe1PI=","next_validators_hash":"kQ9NNQ26Q3l5aXF2IwraaweoLzusIVDUA53AycOe1PI=","consensus_hash":"BICRvH3cKD93v7+R1zxE2ljD34qcvIZ0Bdi389qtoi8=","app_hash":"j3x5PuEH14QVBqWUv6BhKitzXMTJ2w47h2Nj99JKtlI=","last_results_hash":"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=","evidence_hash":"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=","proposer_address":"YVEp2+U79qlwIpwaKr4rObvWbHo="},"data":{"txs":[]},"evidence":{"evidence":[]},"last_commit":{"height":"90766","round":0,"block_id":{"hash":"wCK5QmPuGJiRpn06Xu7ZjhxHwzBXVZrGqngMMzeRq8w=","part_set_header":{"total":1,"hash":"xYxXM7rX6Qcq/Yx3MpZeQA+FCeUbKSVr/FEuzfAFFQk="}},"signatures":[{"block_id_flag":"BLOCK_ID_FLAG_COMMIT","validator_address":"WRBCW5t/kdjOaTvYz9TaySfc8xU=","timestamp":"2022-11-03T09:12:23.109926769Z","signature":"LAbAFCM2MlT1QeNnhZoD8xPe6fO6GExtDNdwa8sokr9UZjurHWn3ad9U2BhTLFVKUF6j7r9G9ILKshljKn4/Aw=="},{"block_id_flag":"BLOCK_ID_FLAG_COMMIT","validator_address":"YVEp2+U79qlwIpwaKr4rObvWbHo=","timestamp":"2022-11-03T09:12:23.104507119Z","signature":"KxHLYnn97GG9prtLA+qurq5GvZogcoExpCvWmOOd8uS3m1Tug5qptSxZ2AObiUfyDwl23oqNNEhkp2XxsxcOCA=="},{"block_id_flag":"BLOCK_ID_FLAG_COMMIT","validator_address":"agVlYuY3F6RGyBUe5zgd0oFhoCM=","timestamp":"2022-11-03T09:12:23.109928409Z","signature":"vvL6yIdxyT4Eus/xw8/RWvFymUFNiOsJ+hHM/qwSJQt427hdUiIh/iH6+yZGz5bdpChW4/Y4bB1QnIA8q1SjBw=="}]}}}`
	glbrResp := tmservice.GetLatestBlockResponse{}
	err := suite.GetQuicksilverApp(suite.chainA).AppCodec().UnmarshalJSON([]byte(respJSON), &glbrResp)
	suite.Require().NoError(err)
	resp, err := glbrResp.Marshal()
	suite.Require().NoError(err)

	err = keeper.SetEpochBlockCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) executeUmeeReservesUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", umeetypes.KeyReserveAmount(umeeBaseDenom), types.ModuleName)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	data := sdk.NewInt(100000)
	resp, err := data.Marshal()
	suite.Require().NoError(err)
	expectedData, err := json.Marshal(data)
	suite.Require().NoError(err)

	// setup for expected

	err = keeper.UmeeReservesUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)

	want := &types.UmeeReservesProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	pd, found := prk.GetProtocolData(ctx, types.ProtocolDataTypeUmeeReserves, umeeBaseDenom)
	suite.Require().True(found)

	value, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeReserves, pd.Data)
	suite.Require().NoError(err)
	result := value.(*types.UmeeReservesProtocolData)
	suite.Require().Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeLeverageModuleBalanceUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	accountPrefix := banktypes.CreateAccountBalancesPrefix(authtypes.NewModuleAddress(umeetypes.LeverageModuleName))

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/bank/key", append(accountPrefix, []byte(umeeBaseDenom)...), types.ModuleName)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	data := sdk.NewInt(1400000)
	resp, err := data.Marshal()
	suite.Require().NoError(err)
	expectedData, err := json.Marshal(data)
	suite.Require().NoError(err)

	// setup for expected

	err = keeper.UmeeLeverageModuleBalanceUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)

	want := &types.UmeeLeverageModuleBalanceProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	pd, found := prk.GetProtocolData(ctx, types.ProtocolDataTypeUmeeLeverageModuleBalance, umeeBaseDenom)
	suite.Require().True(found)

	value, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeLeverageModuleBalance, pd.Data)
	suite.Require().NoError(err)
	result := value.(*types.UmeeLeverageModuleBalanceProtocolData)
	suite.Require().Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeUTokenSupplyUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", umeetypes.KeyUTokenSupply(umeetypes.UTokenPrefix+umeeBaseDenom), types.ModuleName)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	data := sdk.NewInt(100000)
	resp, err := data.Marshal()
	suite.Require().NoError(err)
	expectedData, err := json.Marshal(data)
	suite.Require().NoError(err)

	// setup for expected

	err = keeper.UmeeUTokenSupplyUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)

	want := &types.UmeeUTokenSupplyProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeetypes.UTokenPrefix + umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	pd, found := prk.GetProtocolData(ctx, types.ProtocolDataTypeUmeeUTokenSupply, umeetypes.UTokenPrefix+umeeBaseDenom)
	suite.Require().True(found)

	value, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeUTokenSupply, pd.Data)
	suite.Require().NoError(err)
	result := value.(*types.UmeeUTokenSupplyProtocolData)
	suite.Require().Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeTotalBorrowsUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", umeetypes.KeyAdjustedTotalBorrow(umeeBaseDenom), types.ModuleName)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	data := sdk.NewDec(150000)
	resp, err := data.Marshal()
	suite.Require().NoError(err)
	expectedData, err := json.Marshal(data)
	suite.Require().NoError(err)

	// setup for expected

	err = keeper.UmeeTotalBorrowsUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)

	want := &types.UmeeTotalBorrowsProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	pd, found := prk.GetProtocolData(ctx, types.ProtocolDataTypeUmeeTotalBorrows, umeeBaseDenom)
	suite.Require().True(found)

	value, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeTotalBorrows, pd.Data)
	suite.Require().NoError(err)
	result := value.(*types.UmeeTotalBorrowsProtocolData)
	suite.Require().Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeInterestScalarUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", umeetypes.KeyInterestScalar(umeeBaseDenom), types.ModuleName)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.Require().True(found, "qid: %s", qid)

	data := sdk.NewDec(1)
	resp, err := data.Marshal()
	suite.Require().NoError(err)
	expectedData, err := json.Marshal(data)
	suite.Require().NoError(err)

	// setup for expected

	err = keeper.UmeeInterestScalarUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.Require().NoError(err)

	want := &types.UmeeInterestScalarProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	pd, found := prk.GetProtocolData(ctx, types.ProtocolDataTypeUmeeInterestScalar, umeeBaseDenom)
	suite.Require().True(found)

	value, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeInterestScalar, pd.Data)
	suite.Require().NoError(err)
	result := value.(*types.UmeeInterestScalarProtocolData)
	suite.Require().Equal(want, result)
}
