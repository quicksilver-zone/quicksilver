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
