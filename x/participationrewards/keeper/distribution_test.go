package keeper_test

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/keeper"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) TestCalcTokenValues() {
	type cases struct {
		name           string
		osmosisParams  types.OsmosisParamsProtocolData
		osmosisPools   []types.OsmosisPoolProtocolData
		osmosisClPools []types.OsmosisClPoolProtocolData
		expectedTvs    keeper.TokenValues
	}

	tests := []cases{
		{
			osmosisParams: types.OsmosisParamsProtocolData{
				ChainID:   "osmosis-1",
				BaseDenom: "uosmo",
				BaseChain: "osmosis-1",
			},
			osmosisPools: []types.OsmosisPoolProtocolData{
				{
					PoolID:      944,
					PoolName:    "atom/qAtom",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeStableSwap,
					PoolData:    json.RawMessage(`{"address":"osmo1awr39mc2hrkt8gq8gt3882ru40ay45k8a3yg69nyypqe9g0ryycs66lhkh","id":944,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/944","amount":"6108537302303463956540"},"pool_liquidity":[{"denom":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","amount":"42678069500"},{"denom":"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC","amount":"70488173547"}],"scaling_factors":[1202853876,1000000000],"scaling_factor_controller":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"}`),
					Denoms: map[string]types.DenomWithZone{
						"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {Denom: "uatom", ChainID: "cosmoshub-4"},
						"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC": {Denom: "uqatom", ChainID: "cosmoshub-4"},
					},
					IsIncentivized: true,
				},
				{
					PoolID:      1,
					PoolName:    "Atom/Osmo",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeBalancer,
					PoolData:    json.RawMessage("{\"address\":\"osmo1mw0ac6rwlp5r8wapwk3zs6g29h8fcscxqakdzw9emkne6c8wjp9q0t3v8t\",\"id\":1,\"pool_params\":{\"swap_fee\":\"0.002000000000000000\",\"exit_fee\":\"0.000000000000000000\",\"smooth_weight_change_params\":null},\"future_pool_governor\":\"24h\",\"total_shares\":{\"denom\":\"gamm/pool/1\",\"amount\":\"216987393856026889179749817\"},\"pool_assets\":[{\"token\":{\"denom\":\"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2\",\"amount\":\"1909639500022\"},\"weight\":\"536870912000000\"},{\"token\":{\"denom\":\"uosmo\",\"amount\":\"35673230362499\"},\"weight\":\"536870912000000\"}],\"total_weight\":\"1073741824000000\"}"),
					Denoms: map[string]types.DenomWithZone{
						"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {Denom: "uatom", ChainID: "cosmoshub-4"},
						"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
					},
					IsIncentivized: false,
				},
				{
					PoolID:      952,
					PoolName:    "qck/osmo",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeBalancer,
					PoolData:    json.RawMessage("{\"address\":\"osmo1k3j5wgcj8um2gnu8qxdm0mzzuh6x66p4p7gn6fraf3wnpfcvg9sq2zhx7j\",\"id\":952,\"pool_params\":{\"swap_fee\":\"0.003000000000000000\",\"exit_fee\":\"0.000000000000000000\",\"smooth_weight_change_params\":null},\"future_pool_governor\":\"168h\",\"total_shares\":{\"denom\":\"gamm/pool/952\",\"amount\":\"281109110456689694028077\"},\"pool_assets\":[{\"token\":{\"denom\":\"ibc/635CB83EF1DFE598B10A3E90485306FD0D47D34217A4BE5FD9977FA010A5367D\",\"amount\":\"1036526700301\"},\"weight\":\"1073741824\"},{\"token\":{\"denom\":\"uosmo\",\"amount\":\"162265452817\"},\"weight\":\"1073741824\"}],\"total_weight\":\"2147483648\"}"),
					Denoms: map[string]types.DenomWithZone{
						"ibc/635CB83EF1DFE598B10A3E90485306FD0D47D34217A4BE5FD9977FA010A5367D": {Denom: "uqck", ChainID: "quicksilver-2"},
						"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
					},
					IsIncentivized: false,
				},
				{
					PoolID:      956,
					PoolName:    "qOsmo/osmo",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeStableSwap,
					PoolData:    json.RawMessage("{\"address\":\"osmo1q023e9m4d3ffvr96xwaeraa62yfvufkufkr7yf7lmacgkuspsuqsga4xp2\",\"id\":956,\"pool_params\":{\"swap_fee\":\"0.003000000000000000\",\"exit_fee\":\"0.000000000000000000\"},\"future_pool_governor\":\"168h\",\"total_shares\":{\"denom\":\"gamm/pool/956\",\"amount\":\"118922578939571354422559\"},\"pool_liquidity\":[{\"denom\":\"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC\",\"amount\":\"217240952822\"},{\"denom\":\"uosmo\",\"amount\":\"260096955062\"}],\"scaling_factors\":[1000000000,1045466083],\"scaling_factor_controller\":\"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy\"}"),
					Denoms: map[string]types.DenomWithZone{
						"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC": {Denom: "uqosmo", ChainID: "osmosis-1"},
						"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
					},
					IsIncentivized: true,
				},
			},
			osmosisClPools: []types.OsmosisClPoolProtocolData{
				{
					PoolID:      1589,
					PoolName:    "qAtom/atom",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeCL,
					PoolData:    json.RawMessage(`{"@type":"/osmosis.concentratedliquidity.v1beta1.Pool","address":"osmo1qlyunnmszlve9z5c9g85u2psv30gsvdf7v9pry7pgtndmffej6fs2tkw04","incentives_address":"osmo12yf4vkpv5f7mpzp358cj9as9jma37errqg3mtv6tckehhgsvasvqu6pnlf","spread_rewards_address":"osmo1akvkltq8grxvmtdk034q9c64phncwf8772up73qv6hnyhgrw6c4sxwjs2s","id":1589,"current_tick_liquidity":"30095650.842241297950528495","token0":"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC","token1":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","current_sqrt_price":"1.077266816456319199106160167400821857","current_tick":160503,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.000500000000000000","last_liquidity_update":"2024-04-07T22:05:02.235723904Z"}`),
					Denoms: map[string]types.DenomWithZone{
						"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {Denom: "uatom", ChainID: "cosmoshub-4"},
						"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC": {Denom: "uqatom", ChainID: "cosmoshub-4"},
					},
					IsIncentivized: true,
				},
				{
					PoolID:      1590,
					PoolName:    "qOsmo/osmo",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeCL,
					PoolData:    json.RawMessage(`{"@type":"/osmosis.concentratedliquidity.v1beta1.Pool","address":"osmo1ww5567k7ya8g2rzfft30mgdqwy6jyfqfl60805kclwc3x44qzp6qgd3d3r","incentives_address":"osmo1e0wfxnfrzg50wkeyakg2qxjf4fw9lrtx88grlhuw9n2ajm2sth2sdpg9ky","spread_rewards_address":"osmo1e9grc4jp9jzwqqxy6cv6ypkkytyae50wygc9zsj3w5040a3yq4es6z09q7","id":1590,"current_tick_liquidity":"2087592960394.399195262984792926","token0":"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC","token1":"uosmo","current_sqrt_price":"1.052413845746709839272442045043131967","current_tick":107574,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.000500000000000000","last_liquidity_update":"2024-05-04T04:10:52.756088069Z"}`),
					Denoms: map[string]types.DenomWithZone{
						"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC": {Denom: "uqosmo", ChainID: "osmosis-1"},
						"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
					},
					IsIncentivized: true,
				},
				{
					PoolID:      1670,
					PoolName:    "dydx/osmo",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeCL,
					PoolData:    json.RawMessage(`{"@type":"/osmosis.concentratedliquidity.v1beta1.Pool","address":"osmo14nx06cfjs9xslz3ehd3tqs65pzad76479spnjfshdwcmxqnvp7rq2q7jcs","incentives_address":"osmo14nh0xcve9axjtu6sag3vkpa06vdet4pzz2arypx3e7r73wac4epsjlcku6","spread_rewards_address":"osmo1kf2p445mvuxve2j978xaxam2jsf2t5kjnn0fmpfasnxe3fmpergqs4m4gy","id":1670,"current_tick_liquidity":"276421909743.387332362512087187","token0":"ibc/094FB70C3006906F67F5D674073D2DAFAFB41537E7033098F5C752F211E7B6C2","token1":"uosmo","current_sqrt_price":"1.843760308332381016802562578162870388","current_tick":2399452,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.010000000000000000","last_liquidity_update":"2024-05-07T21:00:19.764413083Z"}`),
					Denoms: map[string]types.DenomWithZone{
						"ibc/094FB70C3006906F67F5D674073D2DAFAFB41537E7033098F5C752F211E7B6C2": {Denom: "usaga", ChainID: "ssc-1"},
						"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
					},
					IsIncentivized: false,
				},
				{
					PoolID:      1246,
					PoolName:    "dydx/usdc",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeCL,
					PoolData:    json.RawMessage(`{"@type":"/osmosis.concentratedliquidity.v1beta1.Pool","address":"osmo1trl50qymun2pk0cuquqk8m0pg5q8pl5dfudmjth7tn3crhvvkv9snhre37","incentives_address":"osmo1y25nwfch5rp3fyawkqjvvfpc8qe3tjywjs2hv6xfdg292umpl50qfk8aj7","spread_rewards_address":"osmo1u9ktz95vadadeve0l6qcfyj47mtd32v93xnz3ru4j33z326rjlhs4skgac","id":1246,"current_tick_liquidity":"292642113174175523.743603256873986902","token0":"ibc/831F0B1BBB1D08A2B75311892876D71565478C532967545476DF4C2D7492E48C","token1":"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4","current_sqrt_price":"0.000001473276708984120163314112246914","current_tick":-106829456,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.000500000000000000","last_liquidity_update":"2024-05-07T20:36:57.080586221Z"}`),
					Denoms: map[string]types.DenomWithZone{
						"ibc/831F0B1BBB1D08A2B75311892876D71565478C532967545476DF4C2D7492E48C": {Denom: "adydx", ChainID: "dydx-mainnet-1"},
						"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4": {Denom: "uusdc", ChainID: "noble-1"},
					},
					IsIncentivized: false,
				},
				{
					PoolID:      1263,
					PoolName:    "osmo/usdc",
					LastUpdated: time.Now().UTC(),
					PoolType:    types.PoolTypeCL,
					PoolData:    json.RawMessage(`{"@type":"/osmosis.concentratedliquidity.v1beta1.Pool","address":"osmo1hueh0egxjt6upzn6d80d3653k5vwjjluqz7lm8ffu4t84488p49s6vquyd","incentives_address":"osmo1glw62sxnarqlte4g4kxh7z5pje0x0llfwe53jfec7uydc9xygv7srmknr9","spread_rewards_address":"osmo15qle4t0azwg8etdn48n8lrvvwwut2wfe6lmkn5hk2faycs5gyalsf7eqgv","id":1263,"current_tick_liquidity":"5753612666799.012276383585267557","token0":"uosmo","token1":"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4","current_sqrt_price":"0.943068131935192879644524489130035036","current_tick":-1106225,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.000500000000000000","last_liquidity_update":"2024-05-07T21:37:57.090751219Z"}`),
					Denoms: map[string]types.DenomWithZone{
						"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4": {Denom: "uusdc", ChainID: "noble-1"},
						"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
					},
					IsIncentivized: false,
				},
			},

			expectedTvs: keeper.TokenValues{
				"uatom":  sdk.MustNewDecFromStr("18.680609802053228684"),
				"uosmo":  sdk.MustNewDecFromStr("1.000000000000000000"),
				"uqatom": sdk.MustNewDecFromStr("21.292547025239143631"),
				"uqck":   sdk.MustNewDecFromStr("0.156547296630061979"),
				"uqosmo": sdk.MustNewDecFromStr("1.076844826109830973"),
				"usaga":  sdk.MustNewDecFromStr("3.399452074581916723"),
				"adydx":  sdk.MustNewDecFromStr("0.000000000002440520"),
				"uusdc":  sdk.MustNewDecFromStr("1.124381939441023032"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()

			qs := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()
			osmoParamsJSON, err := json.Marshal(tt.osmosisParams)
			suite.NoError(err)
			data := types.ProtocolData{
				Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisParams)],
				Data: osmoParamsJSON,
			}
			qs.ParticipationRewardsKeeper.SetProtocolData(ctx, tt.osmosisParams.GenerateKey(), &data)

			for _, pool := range tt.osmosisPools {
				poolJSON, err := json.Marshal(pool)
				suite.NoError(err)
				data := types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)],
					Data: poolJSON,
				}
				qs.ParticipationRewardsKeeper.SetProtocolData(ctx, pool.GenerateKey(), &data)
			}
			for _, pool := range tt.osmosisClPools {
				poolJSON, err := json.Marshal(pool)
				suite.NoError(err)
				data := types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisCLPool)],
					Data: poolJSON,
				}
				qs.ParticipationRewardsKeeper.SetProtocolData(ctx, pool.GenerateKey(), &data)
			}
			tvs, err := qs.ParticipationRewardsKeeper.CalcTokenValues(ctx)
			suite.NoError(err)
			suite.Equal(tt.expectedTvs, tvs)
		})
	}
}

func (suite *KeeperTestSuite) TestCalcTokenValuesIncCLPools() {
	qs := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()
	osmoParams := types.OsmosisParamsProtocolData{
		ChainID:   "osmosis-1",
		BaseDenom: "uosmo",
		BaseChain: "osmosis-1",
	}
	osmoParamsJSON, err := json.Marshal(osmoParams)
	suite.NoError(err)
	data := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisParams)],
		Data: osmoParamsJSON,
	}
	qs.ParticipationRewardsKeeper.SetProtocolData(ctx, osmoParams.GenerateKey(), &data)

	pools := []types.OsmosisPoolProtocolData{}
	err = json.Unmarshal([]byte(prodPDString), &pools)
	suite.NoError(err)
	for _, pool := range pools {
		poolJSON, err := json.Marshal(pool)
		suite.NoError(err)
		data := types.ProtocolData{
			Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)],
			Data: poolJSON,
		}
		qs.ParticipationRewardsKeeper.SetProtocolData(ctx, pool.GenerateKey(), &data)
	}

	clpools := []types.OsmosisClPoolProtocolData{}
	err = json.Unmarshal([]byte(prodPDString2), &clpools)
	suite.NoError(err)
	for _, pool := range clpools {
		poolJSON, err := json.Marshal(pool)
		suite.NoError(err)
		data := types.ProtocolData{
			Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisCLPool)],
			Data: poolJSON,
		}
		qs.ParticipationRewardsKeeper.SetProtocolData(ctx, pool.GenerateKey(), &data)
	}

	tvs, err := qs.ParticipationRewardsKeeper.CalcTokenValues(ctx)
	suite.NoError(err)
	expected := map[string]sdk.Dec{
		"uatom":   sdk.MustNewDecFromStr("12.677691539450075775"),
		"ujuno":   sdk.MustNewDecFromStr("0.252676989792765540"),
		"uosmo":   sdk.MustNewDecFromStr("1.000000000000000000"),
		"uqatom":  sdk.MustNewDecFromStr("16.637863612013346262"),
		"uqosmo":  sdk.MustNewDecFromStr("1.201708298212189075"),
		"uqregen": sdk.MustNewDecFromStr("0.069536169040918330"),
		"uqsomm":  sdk.MustNewDecFromStr("0.073765486093994849"),
		"uqstars": sdk.MustNewDecFromStr("0.029935583655522436"),
		"uregen":  sdk.MustNewDecFromStr("0.052839431713034795"),
		"usomm":   sdk.MustNewDecFromStr("0.069792907126472964"),
		"ustars":  sdk.MustNewDecFromStr("0.021812145938948139"),
	}

	for denom, expectedValue := range expected {
		suite.Equal(tvs[denom], expectedValue)
	}

}

var prodPDString = `[
	  {"PoolID":1,"PoolName":"ATOM/OSMO","LastUpdated":"2024-06-29T17:00:50.109029446Z","PoolData":{"address":"osmo1mw0ac6rwlp5r8wapwk3zs6g29h8fcscxqakdzw9emkne6c8wjp9q0t3v8t","id":1,"pool_params":{"swap_fee":"0.002000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/1","amount":"51756990084398530307168264"},"pool_assets":[{"token":{"denom":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","amount":"575276315709"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"7293175680510"},"weight":"536870912000000"}],"total_weight":"1073741824000000"},"PoolType":"balancer","Denoms":{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2":{"Denom":"uatom","ChainID":"cosmoshub-4"},"uosmo":{"Denom":"uosmo","ChainID":"osmosis-1"}},"IsIncentivized":false},
	  {"PoolID":1087,"PoolName":"SOMM/qSOMM","LastUpdated":"2024-06-29T17:09:07.817925952Z","PoolData":{"address":"osmo1unwajz776rcsvaaehrq82qldwfw4zeqp7jgty09cw4lytuwfw3pqvs0cmt","id":1087,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/1087","amount":"1052495305723789511288418"},"pool_liquidity":[{"denom":"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E","amount":"105275083846"},{"denom":"ibc/EAF76AD1EEF7B16D167D87711FB26ABE881AC7D9F7E6D0CF313D5FA530417208","amount":"107870858010"}],"scaling_factors":[1057053811,1000000000],"scaling_factor_governor":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"},"PoolType":"stableswap","Denoms":{"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E":{"Denom":"usomm","ChainID":"sommelier-3"},"ibc/EAF76AD1EEF7B16D167D87711FB26ABE881AC7D9F7E6D0CF313D5FA530417208":{"Denom":"uqsomm","ChainID":"sommelier-3"}},"IsIncentivized":true},
	  {"PoolID":42,"PoolName":"REGEN/OSMO","LastUpdated":"2024-06-29T17:01:08.083152922Z","PoolData":{"address":"osmo1txawpctjs6phpqsnkx2r5qud7yvekw93394anhuzz4dquy5jggssgqtn0l","id":42,"pool_params":{"swap_fee":"0.002000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/42","amount":"26509813273140138926958"},"pool_assets":[{"token":{"denom":"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076","amount":"853135473387"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"45079193588"},"weight":"536870912000000"}],"total_weight":"1073741824000000"},"PoolType":"balancer","Denoms":{"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076":{"Denom":"uregen","ChainID":"regen-1"},"uosmo":{"Denom":"uosmo","ChainID":"osmosis-1"}},"IsIncentivized":false},
	  {"PoolID":497,"PoolName":"JUNO/OSMO","LastUpdated":"2024-06-29T17:01:37.013088205Z","PoolData":{"address":"osmo1h7yfu7x4qsv2urnkl4kzydgxegdfyjdry5ee4xzj98jwz0uh07rqdkmprr","id":497,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"total_shares":{"denom":"gamm/pool/497","amount":"148794601473627832073337"},"pool_assets":[{"token":{"denom":"ibc/46B44899322F3CD854D2D46DEEF881958467CDD4B3B10086DA49296BBED94BED","amount":"573672533941"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"144953849003"},"weight":"536870912000000"}],"total_weight":"1073741824000000"},"PoolType":"balancer","Denoms":{"ibc/46B44899322F3CD854D2D46DEEF881958467CDD4B3B10086DA49296BBED94BED":{"Denom":"ujuno","ChainID":"juno-1"},"uosmo":{"Denom":"uosmo","ChainID":"osmosis-1"}},"IsIncentivized":false},
	  {"PoolID":604,"PoolName":"STARS/OSMO","LastUpdated":"2024-06-29T17:02:07.675839861Z","PoolData":{"address":"osmo1thscstwxp87g0ygh7le3h92f9ff4sel9y9d2eysa25p43yf43rysk7jp93","id":604,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/604","amount":"27287206734846955128"},"pool_assets":[{"token":{"denom":"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4","amount":"6981654162834"},"weight":"21474836480"},{"token":{"denom":"uosmo","amount":"152284859495"},"weight":"21474836480"}],"total_weight":"42949672960"},"PoolType":"balancer","Denoms":{"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4":{"Denom":"ustars","ChainID":"stargaze-1"},"uosmo":{"Denom":"uosmo","ChainID":"osmosis-1"}},"IsIncentivized":false},
	  {"PoolID":627,"PoolName":"SOMM/OSMO","LastUpdated":"2024-06-29T17:01:37.013088205Z","PoolData":{"address":"osmo19qawwfrlkz9upglmpqj6akgz9ap7v2mnd05pxzgmxw3ywz58wnvqtet2mg","id":627,"pool_params":{"swap_fee":"0.002000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/627","amount":"58898324380432389355671"},"pool_assets":[{"token":{"denom":"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E","amount":"383959778025"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"26797669128"},"weight":"536870912000000"}],"total_weight":"1073741824000000"},"PoolType":"balancer","Denoms":{"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E":{"Denom":"usomm","ChainID":"sommelier-3"},"uosmo":{"Denom":"uosmo","ChainID":"osmosis-1"}},"IsIncentivized":false},
	  {"PoolID":903,"PoolName":"qSTARS/STARS","LastUpdated":"2024-06-29T17:00:50.109029446Z","PoolData":{"address":"osmo1cxlrfu8r0v3cyqj78fuvlsmhjdgna0r7tum8cpd0g3x7w7pte8fsfvcs84","id":903,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/903","amount":"33285682356451671868036"},"pool_liquidity":[{"denom":"ibc/46C83BB054E12E189882B5284542DB605D94C99827E367C9192CF0579CD5BC83","amount":"201086272390"},{"denom":"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4","amount":"672153817668"}],"scaling_factors":[1,1]},"PoolType":"stableswap","Denoms":{"ibc/46C83BB054E12E189882B5284542DB605D94C99827E367C9192CF0579CD5BC83":{"Denom":"uqstars","ChainID":"stargaze-1"},"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4":{"Denom":"ustars","ChainID":"stargaze-1"}},"IsIncentivized":true},
	  {"PoolID":944,"PoolName":"ATOM/qATOM","LastUpdated":"2024-06-29T17:01:08.083152922Z","PoolData":{"address":"osmo1awr39mc2hrkt8gq8gt3882ru40ay45k8a3yg69nyypqe9g0ryycs66lhkh","id":944,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/944","amount":"376682006004054987881"},"pool_liquidity":[{"denom":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","amount":"3590725451"},{"denom":"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC","amount":"3561329984"}],"scaling_factors":[1318616518,1000000000],"scaling_factor_governor":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"},"PoolType":"stableswap","Denoms":{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2":{"Denom":"uatom","ChainID":"cosmoshub-4"},"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC":{"Denom":"uqatom","ChainID":"cosmoshub-4"}},"IsIncentivized":true},
	  {"PoolID":948,"PoolName":"REGEN/qREGEN","LastUpdated":"2024-06-29T17:00:50.109029446Z","PoolData":{"address":"osmo1hylqy4uu5el36wykhzzhj786eh8rx4epyvg6nrtl503wjufz8z3sdptdzw","id":948,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/948","amount":"203791959210650860086725"},"pool_liquidity":[{"denom":"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076","amount":"247440093055"},{"denom":"ibc/79A676508A2ECA1021EDDC7BB9CF70CEEC9514C478DA526A5A8B3E78506C2206","amount":"177233580285"}],"scaling_factors":[1315922033,1000000000],"scaling_factor_governor":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"},"PoolType":"stableswap","Denoms":{"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076":{"Denom":"uregen","ChainID":"regen-1"},"ibc/79A676508A2ECA1021EDDC7BB9CF70CEEC9514C478DA526A5A8B3E78506C2206":{"Denom":"uqregen","ChainID":"regen-1"}},"IsIncentivized":true},
	  {"PoolID":956,"PoolName":"qOSMO/OSMO","LastUpdated":"2024-06-29T17:01:08.083152922Z","PoolData":{"address":"osmo1q023e9m4d3ffvr96xwaeraa62yfvufkufkr7yf7lmacgkuspsuqsga4xp2","id":956,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/956","amount":"2831883387499914752885"},"pool_liquidity":[{"denom":"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC","amount":"4678864416"},{"denom":"uosmo","amount":"7164359952"}],"scaling_factors":[1000000000,1182092750],"scaling_factor_governor":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"},"PoolType":"stableswap","Denoms":{"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC":{"Denom":"uqosmo","ChainID":"osmosis-1"},"uosmo":{"Denom":"uosmo","ChainID":"osmosis-1"}},"IsIncentivized":true}
	]`

var prodPDString2 = `[
	{"PoolID":1590,"PoolName":"qOSMO/OSMO","LastUpdated":"2024-06-29T17:01:08.083152922Z","PoolData":{"address":"osmo1ww5567k7ya8g2rzfft30mgdqwy6jyfqfl60805kclwc3x44qzp6qgd3d3r","incentives_address":"osmo1e0wfxnfrzg50wkeyakg2qxjf4fw9lrtx88grlhuw9n2ajm2sth2sdpg9ky","spread_rewards_address":"osmo1e9grc4jp9jzwqqxy6cv6ypkkytyae50wygc9zsj3w5040a3yq4es6z09q7","id":1590,"current_tick_liquidity":"987929067650.598626855339345890","token0":"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC","token1":"uosmo","current_sqrt_price":"1.102849323600434041469545632253046567","current_tick":216276,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.000500000000000000","last_liquidity_update":"2024-09-16T11:42:46.811999582Z"}, "Denoms":{"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC":{"Denom":"uqosmo","ChainID":"osmosis-1"},"uosmo":{"Denom":"uosmo","ChainID":"osmosis-1"}},"IsIncentivized":true}
]`
