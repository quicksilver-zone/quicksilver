package keeper_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/proto/tendermint/crypto"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) Test_OsmosisCL_ValidateClaim() {
	cases := []struct {
		name          string
		submitAddress string
		mappedAddress string
		expectErr     bool
	}{
		{
			"valid - submitter",
			"quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			"",
			false,
		},
		{
			"valid - mapped",
			"quick1th90qmcxwsr6wwjsna92l63u7j8qp3fmj47v5x",
			"osmo1jc24kwznud9m3mwqmcz3xw33ndjuufngltcdt6",
			false,
		},
		{
			"invalid - no mapped",
			"quick1th90qmcxwsr6wwjsna92l63u7j8qp3fmj47v5x",
			"",
			true,
		},
		{
			"invalid - mapped no match",
			"quick1th90qmcxwsr6wwjsna92l63u7j8qp3fmj47v5x",
			"osmo18e8drgypatsw0skt5ywzeqhlk365hlul3pnasr",
			true,
		},
	}

	for _, c := range cases {
		suite.Run(c.name, func() {
			suite.SetupTest()

			app := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			// create osmosis params protocol data
			osmosisParamPd := types.OsmosisParamsProtocolData{
				ChainID:   "osmosis-1",
				BaseDenom: "uosmo",
				BaseChain: "osmosis-1",
			}

			blob, err := json.Marshal(osmosisParamPd)
			suite.NoError(err)
			pd := types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisParams)], blob)
			app.ParticipationRewardsKeeper.SetProtocolData(ctx, osmosisParamPd.GenerateKey(), pd)

			// {"PoolID":1589,"PoolName":"qATOM/ATOM","LastUpdated":"2025-02-15T17:02:01.669154947Z","PoolData":,"PoolType":"concentrated-liquidity","Denoms":{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2":{"Denom":"uatom","ChainID":"cosmoshub-4"},"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC":{"Denom":"uqatom","ChainID":"cosmoshub-4"}},"IsIncentivized":true}
			poolPd := types.OsmosisClPoolProtocolData{
				PoolID:      1589,
				PoolName:    "qATOM/ATOM",
				LastUpdated: time.Now(),
				PoolData:    []byte(`{"address":"osmo1qlyunnmszlve9z5c9g85u2psv30gsvdf7v9pry7pgtndmffej6fs2tkw04","incentives_address":"osmo12yf4vkpv5f7mpzp358cj9as9jma37errqg3mtv6tckehhgsvasvqu6pnlf","spread_rewards_address":"osmo1akvkltq8grxvmtdk034q9c64phncwf8772up73qv6hnyhgrw6c4sxwjs2s","id":1589,"current_tick_liquidity":"12809281828.243176405811346618","token0":"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC","token1":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","current_sqrt_price":"1.188208396503728215435481479429664466","current_tick":411839,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.000500000000000000","last_liquidity_update":"2025-02-13T22:26:22.8002161Z"}`),
				PoolType:    "concentrated-liquidity",
				Denoms: map[string]types.DenomWithZone{
					"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {
						Denom:   "uatom",
						ChainID: "cosmoshub-4",
					},
					"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC": {
						Denom:   "uqatom",
						ChainID: "cosmoshub-4",
					},
				},
				IsIncentivized: true,
			}

			blob, err = json.Marshal(poolPd)
			suite.NoError(err)
			pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisCLPool)], blob)
			app.ParticipationRewardsKeeper.SetProtocolData(ctx, poolPd.GenerateKey(), pd)

			// remove the liquid token protocol data for uqatom on osmosis added in keeper_test.go
			app.ParticipationRewardsKeeper.DeleteProtocolData(ctx, types.GetProtocolDataKey(types.ProtocolDataTypeLiquidToken, []byte(fmt.Sprintf("%s_%s", "osmosis-1", cosmosIBCDenom))))

			liquidTokenPd := types.LiquidAllowedDenomProtocolData{
				ChainID:               "osmosis-1",
				RegisteredZoneChainID: "cosmoshub-4",
				QAssetDenom:           "uqatom",
				IbcDenom:              "ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC",
			}

			blob, err = json.Marshal(liquidTokenPd)
			suite.NoError(err)
			pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
			app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

			if c.mappedAddress != "" {
				app.InterchainstakingKeeper.SetRemoteAddressMap(ctx, addressutils.MustAccAddressFromBech32(c.submitAddress, ""), addressutils.MustAccAddressFromBech32(c.mappedAddress, ""), "osmosis-1")
			}
			// claim
			key, err := base64.StdEncoding.DecodeString(testCLKey)
			suite.NoError(err)
			data, err := base64.StdEncoding.DecodeString(testCLData)
			suite.NoError(err)
			proofOps := crypto.ProofOps{}
			err = json.Unmarshal([]byte(testCLProofOps), &proofOps)
			suite.NoError(err)

			msgClaim := types.MsgSubmitClaim{
				UserAddress: c.submitAddress,
				Zone:        "cosmoshub-4",
				SrcZone:     "osmosis-1",
				Proofs: []*cmtypes.Proof{
					{
						Key:       key,
						Data:      data,
						ProofOps:  &proofOps,
						Height:    27966503,
						ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeOsmosisCLPool)],
					},
				},
			}

			suite.NoError(msgClaim.ValidateBasic())

			out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeOsmosisCLPool].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
			if c.expectErr {
				suite.Error(err)
				suite.Equal(out, math.ZeroInt())
			} else {
				suite.NoError(err)
				suite.Equal(out, math.NewInt(923443))
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_OsmosisCL_ValidateClaim_IgnoreDoubleClaim() {
	suite.Run("double_claim", func() {
		suite.SetupTest()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		// create osmosis params protocol data
		osmosisParamPd := types.OsmosisParamsProtocolData{
			ChainID:   "osmosis-1",
			BaseDenom: "uosmo",
			BaseChain: "osmosis-1",
		}

		blob, err := json.Marshal(osmosisParamPd)
		suite.NoError(err)
		pd := types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisParams)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, osmosisParamPd.GenerateKey(), pd)

		// {"PoolID":1589,"PoolName":"qATOM/ATOM","LastUpdated":"2025-02-15T17:02:01.669154947Z","PoolData":,"PoolType":"concentrated-liquidity","Denoms":{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2":{"Denom":"uatom","ChainID":"cosmoshub-4"},"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC":{"Denom":"uqatom","ChainID":"cosmoshub-4"}},"IsIncentivized":true}
		poolPd := types.OsmosisClPoolProtocolData{
			PoolID:      1589,
			PoolName:    "qATOM/ATOM",
			LastUpdated: time.Now(),
			PoolData:    []byte(`{"address":"osmo1qlyunnmszlve9z5c9g85u2psv30gsvdf7v9pry7pgtndmffej6fs2tkw04","incentives_address":"osmo12yf4vkpv5f7mpzp358cj9as9jma37errqg3mtv6tckehhgsvasvqu6pnlf","spread_rewards_address":"osmo1akvkltq8grxvmtdk034q9c64phncwf8772up73qv6hnyhgrw6c4sxwjs2s","id":1589,"current_tick_liquidity":"12809281828.243176405811346618","token0":"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC","token1":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","current_sqrt_price":"1.188208396503728215435481479429664466","current_tick":411839,"tick_spacing":100,"exponent_at_price_one":-6,"spread_factor":"0.000500000000000000","last_liquidity_update":"2025-02-13T22:26:22.8002161Z"}`),
			PoolType:    "concentrated-liquidity",
			Denoms: map[string]types.DenomWithZone{
				"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {
					Denom:   "uatom",
					ChainID: "cosmoshub-4",
				},
				"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC": {
					Denom:   "uqatom",
					ChainID: "cosmoshub-4",
				},
			},
			IsIncentivized: true,
		}

		blob, err = json.Marshal(poolPd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisCLPool)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, poolPd.GenerateKey(), pd)

		// remove the liquid token protocol data for uqatom on osmosis added in keeper_test.go
		app.ParticipationRewardsKeeper.DeleteProtocolData(ctx, types.GetProtocolDataKey(types.ProtocolDataTypeLiquidToken, []byte(fmt.Sprintf("%s_%s", "osmosis-1", cosmosIBCDenom))))

		liquidTokenPd := types.LiquidAllowedDenomProtocolData{
			ChainID:               "osmosis-1",
			RegisteredZoneChainID: "cosmoshub-4",
			QAssetDenom:           "uqatom",
			IbcDenom:              "ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC",
		}

		blob, err = json.Marshal(liquidTokenPd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

		// claim
		key, err := base64.StdEncoding.DecodeString(testCLKey)
		suite.NoError(err)
		data, err := base64.StdEncoding.DecodeString(testCLData)
		suite.NoError(err)
		proofOps := crypto.ProofOps{}
		err = json.Unmarshal([]byte(testCLProofOps), &proofOps)
		suite.NoError(err)

		msgClaim := types.MsgSubmitClaim{
			UserAddress: "quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			Zone:        "cosmoshub-4",
			SrcZone:     "osmosis-1",
			Proofs: []*cmtypes.Proof{
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    27966503,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeOsmosisCLPool)],
				},
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    27966503,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeOsmosisCLPool)],
				},
			},
		}

		suite.NoError(msgClaim.ValidateBasic())

		out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeOsmosisCLPool].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
		suite.NoError(err)
		suite.Equal(out, math.NewInt(923443))
	})
}

func (suite *KeeperTestSuite) Test_OsmosisCL_ValidateClaim_FailUnclaimablePool() {
	suite.Run("fail_unclaimable_pool", func() {
		suite.SetupTest()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		// create osmosis params protocol data
		osmosisParamPd := types.OsmosisParamsProtocolData{
			ChainID:   "osmosis-1",
			BaseDenom: "uosmo",
			BaseChain: "osmosis-1",
		}

		blob, err := json.Marshal(osmosisParamPd)
		suite.NoError(err)
		pd := types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisParams)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, osmosisParamPd.GenerateKey(), pd)

		// remove the liquid token protocol data for uqatom on osmosis added in keeper_test.go
		app.ParticipationRewardsKeeper.DeleteProtocolData(ctx, types.GetProtocolDataKey(types.ProtocolDataTypeLiquidToken, []byte(fmt.Sprintf("%s_%s", "osmosis-1", cosmosIBCDenom))))

		liquidTokenPd := types.LiquidAllowedDenomProtocolData{
			ChainID:               "osmosis-1",
			RegisteredZoneChainID: "cosmoshub-4",
			QAssetDenom:           "uqatom",
			IbcDenom:              "ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC",
		}

		blob, err = json.Marshal(liquidTokenPd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

		// claim
		key, err := base64.StdEncoding.DecodeString(testCLKey)
		suite.NoError(err)
		data, err := base64.StdEncoding.DecodeString(testCLData)
		suite.NoError(err)
		proofOps := crypto.ProofOps{}
		err = json.Unmarshal([]byte(testCLProofOps), &proofOps)
		suite.NoError(err)

		msgClaim := types.MsgSubmitClaim{
			UserAddress: "quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			Zone:        "cosmoshub-4",
			SrcZone:     "osmosis-1",
			Proofs: []*cmtypes.Proof{
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    27966503,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeOsmosisCLPool)],
				},
			},
		}

		suite.NoError(msgClaim.ValidateBasic())

		out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeOsmosisCLPool].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
		suite.Error(err)
		suite.Equal(out, math.ZeroInt())
	})
}

const (
	testCLKey      = "CDQ0ODQzMDg="
	testCLData     = "CNTZkQISK29zbW8xamMyNGt3em51ZDltM213cW1jejN4dzMzbmRqdXVmbmdsdGNkdDYYtQwggJrAzP//////ASiAg4qjATIMCJ6C4a8GEO/llNcBOhkxMDk3MjQxOTk3MDA4ODY0MDM5MTIwMzk0"
	testCLProofOps = `{"ops":[{"type":"ics23:iavl","key":"CDQ0ODQzMDg=","data":"CpkKCggINDQ4NDMwOBJvCNTZkQISK29zbW8xamMyNGt3em51ZDltM213cW1jejN4dzMzbmRqdXVmbmdsdGNkdDYYtQwggJrAzP//////ASiAg4qjATIMCJ6C4a8GEO/llNcBOhkxMDk3MjQxOTk3MDA4ODY0MDM5MTIwMzk0Gg4IARgBIAEqBgACnIfaDSIsCAESKAIEvuHbDiBBTq/xC18ddO3mK7UhRyIvYPjIlKpRYxnasTjqXOPstCAiLAgBEigEBo7utBQgxMQbssXWidca/u5pvIlJniE7MohxGo6TpVEppmQOw58gIi4IARIHCBKO7rQUIBohILdn2gGk4Zg+OsQX5RFf1G1hhB0tjp+R/pyHNu+qieYnIi4IARIHCiTQx4oZIBohIJtqOUkahsBWGc4HprmsQHwsZmQYy5UQbwyBijru9WKYIi4IARIHDErC1YsZIBohIAvpXR1CPs1FiWcPS9UbP9i5nXhmO67FPoaJp8DaMTBIIiwIARIoDnDC1YsZIIURG42rZ3mzUv/1ujrPMow7MuSgyQyJhB6U1lfhlt9rICItCAESKRCAAorfkxkgvw8gVEZ3XMlfUXFzdR1rq+V1wl4MnoWlvjHiE7PBOuggIi8IARIIEqIE3PPFGSAaISDSATwq0Y8Yvf01PEO1/uK8cPX2bKZ/1j+1JiFpUm1nRSIvCAESCBaGCPLy+hkgGiEg+B7p/5wDLrO+xE7YLQggSX00f8YZsG1jc5+JPqWgbLQiLwgBEggYjg3OqMkaIBohIN9NF5ZK4w6J/0aTvu+Rtrbzy6wPmXZTHgQW81IjKUgYIi8IARIIGuwZzM7VGiAaISC0YFZDP0+5fhyh/PpCrfQRWCGVpudOnYhfrZPEuTWGgSIvCAESCBzOJszO1RogGiEgxn4MHr51Gf9gWElfXfUI44Qka6cMAhCkpl3Xtor1S5kiLwgBEgge+k7MztUaIBohIFamH1OXxNAu4L+6G1fuLUyLORSstoilhrBonKag6SMJIi0IARIpINB7zM7VGiBQ1I9ehggVa/prZeWfylOhri7fZ+7kZx1T59IGvZl2FyAiLggBEioimq8CzM7VGiDCmq4g+5FYu2tGEb9+RXs+skcsw/qykQOp6yEjBbYtWCAiMAgBEgkktooEmuTVGiAaISC40ZK80G/Comaj8Dafulslnanxyr+dAojTmUZJqyK69CIuCAESKibM9Aaa5NUaIPPMcsSUgOr6QyFwS/TNGnFU5UnG6Cdxpuf2vW3mKs/GICIuCAESKii2zgvM8NUaILw5IsMkkuoAm5ngIEQ1yE4s0yROYtBB8HS636KkUB+bICIuCAESKiqmnBnM8NUaIHeDfZKyPsji0lCeqhuQUGhATz6vcRo7SqjunMso+YhcICIwCAESCSya+DXM8NUaIBohII80wcbVBqyovgvSMFI8TzcB8nWSC7om0Iy/+4Xm0p4BIjEIARIKMLDxqwHM8NUaIBohIGRbX1NJW+cTtqiSR5dzUxSjXwsacDZvwZMQ22BvTvMdIjEIARIKMoKrrALM8NUaIBohIBBra7lkFMjztGFG/5QsbBsgKI3j7+NtWQialuznToO0IjEIARIKNo6mzg7M8NUaIBohINInbIzVSSg32wgPfu839mH9zbfxgGBXpwwpkAJ9pKOGIjEIARIKOPLEsSLM8NUaIBohINegwDlo4fhPy81H6TL0XmYmECAceqSmLQc3tbk+giV7"},{"type":"ics23:simple","key":"Y29uY2VudHJhdGVkbGlxdWlkaXR5","data":"CrYCChVjb25jZW50cmF0ZWRsaXF1aWRpdHkSIP8nWkBeKCaC9tR4v617t1TYpHBk5kHp2pD4dDLOTsIPGgkIARgBIAEqAQAiJwgBEgEBGiDZk8aockSx+3/iSkkoyTN5GHrFAAwRClbfkrtRPPTrliIlCAESIQGPMob8oHtwh4PiSJtyx1Atw7mathzNOOBhNVKjw/fgySIlCAESIQGDFiX29ktAYjAJ61JjEQ2brLVkNTN8kdtputXDLSm1PyInCAESAQEaIF3iNXEZHuGDV7kipO4La/jnyfhUilOqkQ+oLFqLLEfNIicIARIBARogbYmgAU0IfgpffCHi+ZOfhQFWSlseI8nsbUA8DZ0dHA8iJwgBEgEBGiANcgcCAPsnhwuGqF19EowWZq1Otf0RNBMxnT9upjDtHQ=="}]}`
)
