package keeper_test

import (
	"encoding/json"
	"time"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/keeper"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestCalcTokenValues() {
	type cases struct {
		name          string
		osmosisParams types.OsmosisParamsProtocolData
		osmosisPools  []types.OsmosisPoolProtocolData
		expectedTvs   keeper.TokenValues
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
					PoolID:         956,
					PoolName:       "qOsmo/osmo",
					LastUpdated:    time.Now().UTC(),
					PoolType:       "stableswap",
					PoolData:       json.RawMessage("{\"address\":\"osmo1q023e9m4d3ffvr96xwaeraa62yfvufkufkr7yf7lmacgkuspsuqsga4xp2\",\"id\":956,\"pool_params\":{\"swap_fee\":\"0.003000000000000000\",\"exit_fee\":\"0.000000000000000000\"},\"future_pool_governor\":\"168h\",\"total_shares\":{\"denom\":\"gamm/pool/956\",\"amount\":\"118922578939571354422559\"},\"pool_liquidity\":[{\"denom\":\"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC\",\"amount\":\"217240952822\"},{\"denom\":\"uosmo\",\"amount\":\"260096955062\"}],\"scaling_factors\":[\"1000000000\",\"1045466083\"],\"scaling_factor_controller\":\"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy\"}"),
					Denoms:         map[string]types.DenomWithZone{},
					IsIncentivized: true,
				},
				{
					PoolID:         944,
					PoolName:       "atom/qAtom",
					LastUpdated:    time.Now().UTC(),
					PoolType:       "stableswap",
					PoolData:       json.RawMessage("{\"address\":\"osmo1awr39mc2hrkt8gq8gt3882ru40ay45k8a3yg69nyypqe9g0ryycs66lhkh\",\"id\":944,\"pool_params\":{\"swap_fee\":\"0.003000000000000000\",\"exit_fee\":\"0.000000000000000000\"},\"future_pool_governor\":\"168h\",\"total_shares\":{\"denom\":\"gamm/pool/944\",\"amount\":\"9298235648962291280150\"},\"pool_liquidity\":[{\"denom\":\"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2\",\"amount\":\"73166403902\"},{\"denom\":\"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC\",\"amount\":\"98593533184\"}],\"scaling_factors\":[\"1071353717\",\"1000000000\"],\"scaling_factor_controller\":\"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy\"}"),
					Denoms:         map[string]types.DenomWithZone{},
					IsIncentivized: true,
				},
				{
					PoolID:      1,
					PoolName:    "Atom/Osmo",
					LastUpdated: time.Now().UTC(),
					PoolType:    "balancer",
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
					PoolType:    "balancer",
					PoolData:    json.RawMessage("{\"address\":\"osmo1k3j5wgcj8um2gnu8qxdm0mzzuh6x66p4p7gn6fraf3wnpfcvg9sq2zhx7j\",\"id\":\"952\",\"pool_params\":{\"swap_fee\":\"0.003000000000000000\",\"exit_fee\":\"0.000000000000000000\",\"smooth_weight_change_params\":null},\"future_pool_governor\":\"168h\",\"total_shares\":{\"denom\":\"gamm/pool/952\",\"amount\":\"281109110456689694028077\"},\"pool_assets\":[{\"token\":{\"denom\":\"ibc/635CB83EF1DFE598B10A3E90485306FD0D47D34217A4BE5FD9977FA010A5367D\",\"amount\":\"1036526700301\"},\"weight\":\"1073741824\"},{\"token\":{\"denom\":\"uosmo\",\"amount\":\"162265452817\"},\"weight\":\"1073741824\"}],\"total_weight\":\"2147483648\"}"),
					Denoms: map[string]types.DenomWithZone{
						"ibc/635CB83EF1DFE598B10A3E90485306FD0D47D34217A4BE5FD9977FA010A5367D": {Denom: "uqck", ChainID: "quicksilver-2"},
						"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
					},
					IsIncentivized: false,
				},
			},
			expectedTvs: keeper.TokenValues{
				"uatom": sdk.MustNewDecFromStr("18.680609802053228677"),
				"uosmo": sdk.MustNewDecFromStr("1.000000000000000000"),
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
			tvs, err := qs.ParticipationRewardsKeeper.CalcTokenValues(ctx)
			suite.NoError(err)
			suite.Equal(tt.expectedTvs, tvs)
		})
	}
}
