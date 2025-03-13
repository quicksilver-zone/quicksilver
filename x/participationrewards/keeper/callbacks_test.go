package keeper_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/model"
	gamm "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/types"
	leveragetypes "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icqkeeper "github.com/quicksilver-zone/quicksilver/x/interchainquery/keeper"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/keeper"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

var PoolCoinDenom = "pool1"

func (suite *KeeperTestSuite) TestOsmosisPoolUpdateCallback() {
	suite.SetupTest()

	// osmosis test pool
	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisPool,
		[]byte(fmt.Sprintf(
			"{\"poolid\":%d,\"poolname\":%q,\"pooltype\":\"stableswap\",\"denoms\":{%q:{\"chainid\": %q, \"denom\":%q}, %q:{\"chainid\": %q, \"denom\":%q}}}",
			944,
			"atom/qatom",
			"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			"cosmoshub-4",
			"uatom",
			"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC",
			"cosmoshub-4",
			"uqatom",
		)),
	)

	ctx := suite.chainA.GetContext()
	app := suite.GetQuicksilverApp(suite.chainA)
	prk := app.ParticipationRewardsKeeper

	prk.PrSubmodules[cmtypes.ClaimTypeOsmosisPool].Hooks(ctx, prk)

	osm := &keeper.OsmosisModule{}
	qid := icqkeeper.GenerateQueryHash("connection-77002", "osmosis-1", "store/gamm/key", osm.GetKeyPrefixPools(944), types.ModuleName, keeper.OsmosisPoolUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	resp, err := base64.StdEncoding.DecodeString("CjAvb3Ntb3Npcy5nYW1tLnBvb2xtb2RlbHMuc3RhYmxlc3dhcC52MWJldGExLlBvb2wS7QIKP29zbW8xYXdyMzltYzJocmt0OGdxOGd0Mzg4MnJ1NDBheTQ1azhhM3lnNjlueXlwcWU5ZzByeXljczY2bGhraBCwBxoVChAzMDAwMDAwMDAwMDAwMDAwEgEwIgQxNjhoKicKDWdhbW0vcG9vbC85NDQSFjMyNzgzMDY2NTQ2MjI2OTAzNDg3OTIyUwpEaWJjLzI3Mzk0RkIwOTJEMkVDQ0Q1NjEyM0M3NEYzNkU0QzFGOTI2MDAxQ0VBREE5Q0E5N0VBNjIyQjI1RjQxRTVFQjISCzI4MDUyMzMzNjEyMlMKRGliYy9GQTYwMjM2NEJFQzMwNUE2OTZDQkRGOTg3MDU4RTk5RDhCNDc5RjAzMThFNDczMTRDNDkxNzNFODgzOEM1QkFDEgszMzUyMjgzNzU2MjoK+NHZzgSAlOvcA0Irb3NtbzE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaG04bTN5eQ==")
	suite.NoError(err)

	var pdi gamm.CFMMPoolI
	err = prk.GetCodec().UnmarshalInterface(resp, &pdi)
	suite.NoError(err)

	err = keeper.OsmosisPoolUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)

	suite.NoError(err)

	_, pooldata, err := keeper.GetAndUnmarshalProtocolData[*types.OsmosisPoolProtocolData](ctx, prk, "944", types.ProtocolDataTypeOsmosisPool)
	suite.NoError(err)

	pool, err := pooldata.GetPool()
	suite.NoError(err)

	liq := pool.GetTotalPoolLiquidity(ctx)
	suite.Equal(liq.AmountOf("ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"), math.NewInt(28052333612))
	suite.Equal(liq.AmountOf("ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC"), math.NewInt(33522837562))
}

func (suite *KeeperTestSuite) TestOsmosisClPoolUpdateCallback() {
	suite.SetupTest()

	// osmosis test pool
	suite.addProtocolData(
		types.ProtocolDataTypeOsmosisCLPool,
		[]byte(fmt.Sprintf(
			"{\"poolid\":%d,\"poolname\":%q,\"pooltype\":\"concentrated-liquidity\",\"denoms\":{%q:{\"chainid\": %q, \"denom\":%q}, %q:{\"chainid\": %q, \"denom\":%q}}}",
			1089,
			"atom/qatom",
			"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			"cosmoshub-4",
			"uatom",
			"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC",
			"cosmoshub-4",
			"uqatom",
		)),
	)

	ctx := suite.chainA.GetContext()
	app := suite.GetQuicksilverApp(suite.chainA)
	prk := app.ParticipationRewardsKeeper

	prk.PrSubmodules[cmtypes.ClaimTypeOsmosisCLPool].Hooks(ctx, prk)

	osm := &keeper.OsmosisClModule{}
	qid := icqkeeper.GenerateQueryHash("connection-77002", "osmosis-1", "store/concentratedliquidity/key", osm.KeyPool(1089), types.ModuleName, keeper.OsmosisClPoolUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	resp, err := base64.StdEncoding.DecodeString("Cj9vc21vMXFseXVubm1zemx2ZTl6NWM5Zzg1dTJwc3YzMGdzdmRmN3Y5cHJ5N3BndG5kbWZmZWo2ZnMydGt3MDQSP29zbW8xMnlmNHZrcHY1ZjdtcHpwMzU4Y2o5YXM5am1hMzdlcnJxZzNtdHY2dGNrZWhoZ3N2YXN2cXU2cG5sZho/b3NtbzFha3ZrbHRxOGdyeHZtdGRrMDM0cTljNjRwaG5jd2Y4NzcydXA3M3F2NmhueWhncnc2YzRzeHdqczJzILUMKhozMDA5NTY1MDg0MjI0MTI5Nzk1MDUyODQ5NTJEaWJjL0ZBNjAyMzY0QkVDMzA1QTY5NkNCREY5ODcwNThFOTlEOEI0NzlGMDMxOEU0NzMxNEM0OTE3M0U4ODM4QzVCQUM6RGliYy8yNzM5NEZCMDkyRDJFQ0NENTYxMjNDNzRGMzZFNEMxRjkyNjAwMUNFQURBOUNBOTdFQTYyMkIyNUY0MUU1RUIyQiUxMDk4ODAwOTE3ODQxODM4MDUzMTIyMDE3MTU2NTE5NzI0NDE2SIPUDFBkWPr//////////wFiDzUwMDAwMDAwMDAwMDAwMGoLCI6xzLAGEIC5s3A=")
	suite.NoError(err)

	var pdi model.Pool
	err = prk.GetCodec().Unmarshal(resp, &pdi)
	suite.NoError(err)

	err = keeper.OsmosisClPoolUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)

	suite.NoError(err)

	_, pooldata, err := keeper.GetAndUnmarshalProtocolData[*types.OsmosisClPoolProtocolData](ctx, prk, "1089", types.ProtocolDataTypeOsmosisCLPool)
	suite.NoError(err)

	pool, err := pooldata.GetPool()
	suite.NoError(err)

	liq, err := pool.SpotPrice(ctx, "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", "ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC")
	suite.NoError(err)
	liq.Equal(osmomath.NewBigDec(1))
}

func (suite *KeeperTestSuite) executeOsmosisPoolUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	osm := &keeper.OsmosisModule{}
	qid := icqkeeper.GenerateQueryHash("connection-77002", "osmosis-1", "store/gamm/key", osm.GetKeyPrefixPools(1), types.ModuleName, keeper.OsmosisPoolUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	var err error
	// resp := []byte{`{"pool":{"@type":"/osmosis.concentratedliquidity.v1beta1.Pool","address":"osmo1x667pejfeygrp9x0725yuxwlg6cc83rv0ehswsg43wresaskxk2s0jwsej","current_sqrt_price":"1.132659465453955135474729197655764713","current_tick":"282917","current_tick_liquidity":"247410935830.646001002045074682","exponent_at_price_one":"-6","id":"1767","incentives_address":"osmo1mfnll4vtt8zk9afzupvm2ryh5xak00trlqshf05v2460sskqp97s452cem","last_liquidity_update":"2024-05-15T16:10:44.483047470Z","spread_factor":"0.000500000000000000","spread_rewards_address":"osmo14y9wf6hv9j5c5fxh54gx5fu43ute99n80nd5x2uzk6tspt7x2llq0h8pz4","tick_spacing":"100","token0":"ibc/79A676508A2ECA1021EDDC7BB9CF70CEEC9514C478DA526A5A8B3E78506C2206","token1":"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076"}}`}
	resp := []byte{10, 26, 47, 111, 115, 109, 111, 115, 105, 115, 46, 103, 97, 109, 109, 46, 118, 49, 98, 101, 116, 97, 49, 46, 80, 111, 111, 108, 18, 202, 2, 10, 63, 111, 115, 109, 111, 49, 109, 119, 48, 97, 99, 54, 114, 119, 108, 112, 53, 114, 56, 119, 97, 112, 119, 107, 51, 122, 115, 54, 103, 50, 57, 104, 56, 102, 99, 115, 99, 120, 113, 97, 107, 100, 122, 119, 57, 101, 109, 107, 110, 101, 54, 99, 56, 119, 106, 112, 57, 113, 48, 116, 51, 118, 56, 116, 16, 1, 26, 6, 10, 1, 48, 18, 1, 48, 34, 4, 49, 54, 56, 104, 42, 43, 10, 11, 103, 97, 109, 109, 47, 112, 111, 111, 108, 47, 49, 18, 28, 49, 48, 48, 48, 48, 48, 48, 50, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 48, 48, 50, 94, 10, 80, 10, 68, 105, 98, 99, 47, 49, 53, 69, 57, 67, 53, 67, 70, 53, 57, 54, 57, 48, 56, 48, 53, 51, 57, 68, 66, 51, 57, 53, 70, 65, 55, 68, 57, 67, 48, 56, 54, 56, 50, 54, 53, 50, 49, 55, 69, 70, 67, 53, 50, 56, 52, 51, 51, 54, 55, 49, 65, 65, 70, 57, 66, 49, 57, 49, 50, 68, 49, 53, 57, 18, 8, 49, 48, 48, 48, 48, 48, 48, 51, 18, 10, 49, 48, 55, 51, 55, 52, 49, 56, 50, 52, 50, 94, 10, 80, 10, 68, 105, 98, 99, 47, 51, 48, 50, 48, 57, 50, 50, 66, 55, 53, 55, 54, 70, 67, 55, 53, 66, 66, 69, 48, 53, 55, 65, 48, 50, 57, 48, 65, 57, 65, 69, 69, 70, 70, 52, 56, 57, 66, 66, 49, 49, 49, 51, 69, 54, 69, 51, 54, 53, 67, 69, 52, 55, 50, 68, 52, 66, 70, 66, 55, 70, 70, 65, 51, 18, 8, 49, 48, 48, 48, 48, 48, 48, 51, 18, 10, 49, 48, 55, 51, 55, 52, 49, 56, 50, 52, 58, 10, 50, 49, 52, 55, 52, 56, 51, 54, 52, 56}
	// respB64 := "Chovb3Ntb3Npcy5nYW1tLnYxYmV0YTEuUG9vbBLKAgo/b3NtbzFtdzBhYzZyd2xwNXI4d2Fwd2szenM2ZzI5aDhmY3NjeHFha2R6dzllbWtuZTZjOHdqcDlxMHQzdjh0EAEaBgoBMBIBMCIEMTY4aCorCgtnYW1tL3Bvb2wvMRIcMTAwMDAwMDI5OTk5OTk5OTk5OTk5OTk5OTkwMDJeClAKRGliYy8xNUU5QzVDRjU5NjkwODA1MzlEQjM5NUZBN0Q5QzA4NjgyNjUyMTdFRkM1Mjg0MzM2NzFBQUY5QjE5MTJEMTU5EggxMDAwMDAwMxIKMTA3Mzc0MTgyNDJeClAKRGliYy8zMDIwOTIyQjc1NzZGQzc1QkJFMDU3QTAyOTBBOUFFRUZGNDg5QkIxMTEzRTZFMzY1Q0U0NzJENEJGQjdGRkEzEggxMDAwMDAwMxIKMTA3Mzc0MTgyNDoKMjE0NzQ4MzY0OA=="
	// resp, err := base64.StdEncoding.DecodeString(respB64)
	// suite.NoError(err)

	// setup for expected
	var pdi gamm.CFMMPoolI
	err = prk.GetCodec().UnmarshalInterface(resp, &pdi)
	suite.NoError(err)
	expectedData, err := json.Marshal(pdi)
	suite.NoError(err)

	err = keeper.OsmosisPoolUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)

	want := &types.OsmosisPoolProtocolData{
		PoolID:      1,
		PoolName:    "atom/osmo",
		LastUpdated: ctx.BlockTime(),
		PoolData:    expectedData,
		PoolType:    "balancer",
		Denoms: map[string]types.DenomWithZone{
			cosmosIBCDenom:  {ChainID: "cosmoshub-4", Denom: "uatom"},
			osmosisIBCDenom: {ChainID: "osmosis-1", Denom: "uosmo"},
		},
	}

	_, oppd, err := keeper.GetAndUnmarshalProtocolData[*types.OsmosisPoolProtocolData](ctx, prk, "1", types.ProtocolDataTypeOsmosisPool)
	suite.NoError(err)
	suite.Equal(want, oppd)
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
		keeper.ValidatorSelectionRewardsCallbackID,
	)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

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
	suite.NoError(err)
	resp, err := qdtrResp.Marshal()
	suite.NoError(err)

	err = keeper.ValidatorSelectionRewardsCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)
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
		keeper.SetEpochBlockCallbackID,
	)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)
	respJSON := `{"block_id":{"hash":"74pkkjg7u1eLtXxFCinnCmln3aVZVOqLCT3OnE3D+VA=","part_set_header":{"total":1,"hash":"UiLM70PpplmmZ85qC0ZKva5kYJmSZ2TTEZ4a7g9G92Q="}},"block":{"header":{"version":{"block":"11","app":"0"},"chain_id":"quickgaia-1","height":"90767","time":"2022-11-03T09:12:23.109926769Z","last_block_id":{"hash":"wCK5QmPuGJiRpn06Xu7ZjhxHwzBXVZrGqngMMzeRq8w=","part_set_header":{"total":1,"hash":"xYxXM7rX6Qcq/Yx3MpZeQA+FCeUbKSVr/FEuzfAFFQk="}},"last_commit_hash":"1Ev2iL1pTgyItBtSFbCRzxwdCtJfaCC1P+zWaDkJ/nU=","data_hash":"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=","validators_hash":"kQ9NNQ26Q3l5aXF2IwraaweoLzusIVDUA53AycOe1PI=","next_validators_hash":"kQ9NNQ26Q3l5aXF2IwraaweoLzusIVDUA53AycOe1PI=","consensus_hash":"BICRvH3cKD93v7+R1zxE2ljD34qcvIZ0Bdi389qtoi8=","app_hash":"j3x5PuEH14QVBqWUv6BhKitzXMTJ2w47h2Nj99JKtlI=","last_results_hash":"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=","evidence_hash":"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=","proposer_address":"YVEp2+U79qlwIpwaKr4rObvWbHo="},"data":{"txs":[]},"evidence":{"evidence":[]},"last_commit":{"height":"90766","round":0,"block_id":{"hash":"wCK5QmPuGJiRpn06Xu7ZjhxHwzBXVZrGqngMMzeRq8w=","part_set_header":{"total":1,"hash":"xYxXM7rX6Qcq/Yx3MpZeQA+FCeUbKSVr/FEuzfAFFQk="}},"signatures":[{"block_id_flag":"BLOCK_ID_FLAG_COMMIT","validator_address":"WRBCW5t/kdjOaTvYz9TaySfc8xU=","timestamp":"2022-11-03T09:12:23.109926769Z","signature":"LAbAFCM2MlT1QeNnhZoD8xPe6fO6GExtDNdwa8sokr9UZjurHWn3ad9U2BhTLFVKUF6j7r9G9ILKshljKn4/Aw=="},{"block_id_flag":"BLOCK_ID_FLAG_COMMIT","validator_address":"YVEp2+U79qlwIpwaKr4rObvWbHo=","timestamp":"2022-11-03T09:12:23.104507119Z","signature":"KxHLYnn97GG9prtLA+qurq5GvZogcoExpCvWmOOd8uS3m1Tug5qptSxZ2AObiUfyDwl23oqNNEhkp2XxsxcOCA=="},{"block_id_flag":"BLOCK_ID_FLAG_COMMIT","validator_address":"agVlYuY3F6RGyBUe5zgd0oFhoCM=","timestamp":"2022-11-03T09:12:23.109928409Z","signature":"vvL6yIdxyT4Eus/xw8/RWvFymUFNiOsJ+hHM/qwSJQt427hdUiIh/iH6+yZGz5bdpChW4/Y4bB1QnIA8q1SjBw=="}]}}}`
	glbrResp := tmservice.GetLatestBlockResponse{}
	err := suite.GetQuicksilverApp(suite.chainA).AppCodec().UnmarshalJSON([]byte(respJSON), &glbrResp)
	suite.NoError(err)
	resp, err := glbrResp.Marshal()
	suite.NoError(err)

	err = keeper.SetEpochBlockCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) executeUmeeReservesUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", leveragetypes.KeyReserveAmount(umeeBaseDenom), types.ModuleName, keeper.UmeeReservesUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	data := sdk.NewInt(100000)
	resp, err := data.Marshal()
	suite.NoError(err)
	expectedData, err := json.Marshal(data)
	suite.NoError(err)

	// setup for expected

	err = keeper.UmeeReservesUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)

	want := &types.UmeeReservesProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	_, result, err := keeper.GetAndUnmarshalProtocolData[*types.UmeeReservesProtocolData](ctx, prk, umeeBaseDenom, types.ProtocolDataTypeUmeeReserves)
	suite.NoError(err)
	suite.Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeLeverageModuleBalanceUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	accountPrefix := banktypes.CreateAccountBalancesPrefix(authtypes.NewModuleAddress(leveragetypes.LeverageModuleName))

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/bank/key", append(accountPrefix, umeeBaseDenom...), types.ModuleName, keeper.UmeeLeverageModuleBalanceUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	data := sdk.NewInt(1400000)
	resp, err := data.Marshal()
	suite.NoError(err)
	expectedData, err := json.Marshal(data)
	suite.NoError(err)

	// setup for expected

	err = keeper.UmeeLeverageModuleBalanceUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)

	want := &types.UmeeLeverageModuleBalanceProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	_, result, err := keeper.GetAndUnmarshalProtocolData[*types.UmeeLeverageModuleBalanceProtocolData](ctx, prk, umeeBaseDenom, types.ProtocolDataTypeUmeeLeverageModuleBalance)
	suite.NoError(err)
	suite.Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeUTokenSupplyUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", leveragetypes.KeyUTokenSupply(leveragetypes.UTokenPrefix+umeeBaseDenom), types.ModuleName, keeper.UmeeUTokenSupplyUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	data := sdk.NewInt(100000)
	resp, err := data.Marshal()
	suite.NoError(err)
	expectedData, err := json.Marshal(data)
	suite.NoError(err)

	// setup for expected

	err = keeper.UmeeUTokenSupplyUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)

	want := &types.UmeeUTokenSupplyProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       leveragetypes.UTokenPrefix + umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	_, result, err := keeper.GetAndUnmarshalProtocolData[*types.UmeeUTokenSupplyProtocolData](ctx, prk, leveragetypes.UTokenPrefix+umeeBaseDenom, types.ProtocolDataTypeUmeeUTokenSupply)
	suite.NoError(err)
	suite.Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeTotalBorrowsUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", leveragetypes.KeyAdjustedTotalBorrow(umeeBaseDenom), types.ModuleName, keeper.UmeeTotalBorrowsUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	data := sdk.NewDec(150000)
	resp, err := data.Marshal()
	suite.NoError(err)
	expectedData, err := json.Marshal(data)
	suite.NoError(err)

	// setup for expected

	err = keeper.UmeeTotalBorrowsUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)

	want := &types.UmeeTotalBorrowsProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	_, result, err := keeper.GetAndUnmarshalProtocolData[*types.UmeeTotalBorrowsProtocolData](ctx, prk, umeeBaseDenom, types.ProtocolDataTypeUmeeTotalBorrows)
	suite.NoError(err)
	suite.Equal(want, result)
}

func (suite *KeeperTestSuite) executeUmeeInterestScalarUpdateCallback() {
	prk := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
	ctx := suite.chainA.GetContext()

	qid := icqkeeper.GenerateQueryHash(umeeTestConnection, umeeTestChain, "store/leverage/key", leveragetypes.KeyInterestScalar(umeeBaseDenom), types.ModuleName, keeper.UmeeInterestScalarUpdateCallbackID)

	query, found := prk.IcqKeeper.GetQuery(ctx, qid)
	suite.True(found, "qid: %s", qid)

	data := sdk.NewDec(1)
	resp, err := data.Marshal()
	suite.NoError(err)
	expectedData, err := json.Marshal(data)
	suite.NoError(err)

	// setup for expected

	err = keeper.UmeeInterestScalarUpdateCallback(
		ctx,
		prk,
		resp,
		query,
	)
	suite.NoError(err)

	want := &types.UmeeInterestScalarProtocolData{
		UmeeProtocolData: types.UmeeProtocolData{
			Denom:       umeeBaseDenom,
			LastUpdated: ctx.BlockTime(),
			Data:        expectedData,
		},
	}

	_, result, err := keeper.GetAndUnmarshalProtocolData[*types.UmeeInterestScalarProtocolData](ctx, prk, umeeBaseDenom, types.ProtocolDataTypeUmeeInterestScalar)
	suite.NoError(err)
	suite.Equal(want, result)
}
