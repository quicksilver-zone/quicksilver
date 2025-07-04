package keeper_test

import (
	"encoding/base64"
	"encoding/json"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/proto/tendermint/crypto"

	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) Test_LiquidToken_ValidateClaim() {
	cases := []struct {
		name          string
		submitAddress string
		expectErr     bool
	}{
		{
			"valid - submitter",
			"quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			false,
		},
		{
			"invalid - bad user",
			"quick1th90qmcxwsr6wwjsna92l63u7j8qp3fmj47v5x",
			true,
		},
	}

	for _, c := range cases {
		suite.Run(c.name, func() {
			suite.SetupTest()

			app := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			liquidTokenPd := types.LiquidAllowedDenomProtocolData{
				ChainID:               "quicksilver-2",
				RegisteredZoneChainID: "cosmoshub-4",
				QAssetDenom:           "uqatom",
				IbcDenom:              "uqatom",
			}

			blob, err := json.Marshal(liquidTokenPd)
			suite.NoError(err)
			pd := types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
			app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

			// claim
			key, err := base64.StdEncoding.DecodeString(testLiquidKey)
			suite.NoError(err)
			data, err := base64.StdEncoding.DecodeString(testLiquidData)
			suite.NoError(err)
			proofOps := crypto.ProofOps{}
			err = json.Unmarshal([]byte(testLiquidProofOps), &proofOps)
			suite.NoError(err)

			msgClaim := types.MsgSubmitClaim{
				UserAddress: c.submitAddress,
				Zone:        "cosmoshub-4",
				SrcZone:     "quicksilver-2",
				Proofs: []*cmtypes.Proof{
					{
						Key:       key,
						Data:      data,
						ProofOps:  &proofOps,
						Height:    11012784,
						ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeLiquidToken)],
					},
				},
			}

			suite.NoError(msgClaim.ValidateBasic())

			out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeLiquidToken].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
			if c.expectErr {
				suite.Error(err)
				suite.Equal(out, math.ZeroInt())
			} else {
				suite.NoError(err)
				suite.Equal(out, math.NewInt(24307))
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_LiquidToken_ValidateClaim_IgnoreDoubleClaim() {
	suite.Run("double_claim", func() {
		suite.SetupTest()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		liquidTokenPd := types.LiquidAllowedDenomProtocolData{
			ChainID:               "quicksilver-2",
			RegisteredZoneChainID: "cosmoshub-4",
			QAssetDenom:           "uqatom",
			IbcDenom:              "uqatom",
		}

		blob, err := json.Marshal(liquidTokenPd)
		suite.NoError(err)
		pd := types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

		// claim
		key, err := base64.StdEncoding.DecodeString(testLiquidKey)
		suite.NoError(err)
		data, err := base64.StdEncoding.DecodeString(testLiquidData)
		suite.NoError(err)
		proofOps := crypto.ProofOps{}
		err = json.Unmarshal([]byte(testLiquidProofOps), &proofOps)
		suite.NoError(err)

		msgClaim := types.MsgSubmitClaim{
			UserAddress: "quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			Zone:        "cosmoshub-4",
			SrcZone:     "quicksilver-2",
			Proofs: []*cmtypes.Proof{
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    11012784,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeLiquidToken)],
				},
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    11012784,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeLiquidToken)],
				},
			},
		}

		suite.NoError(msgClaim.ValidateBasic())

		out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeLiquidToken].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
		suite.NoError(err)
		suite.Equal(out, math.NewInt(24307))
	})
}

func (suite *KeeperTestSuite) Test_LiquidToken_ValidateClaim_IgnoreUnclaimableDenom() {
	suite.Run("ignore_unclaimable_denom", func() {
		suite.SetupTest()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		// claim
		key, err := base64.StdEncoding.DecodeString(testLiquidKey)
		suite.NoError(err)
		data, err := base64.StdEncoding.DecodeString(testLiquidData)
		suite.NoError(err)
		proofOps := crypto.ProofOps{}
		err = json.Unmarshal([]byte(testLiquidProofOps), &proofOps)
		suite.NoError(err)

		msgClaim := types.MsgSubmitClaim{
			UserAddress: "quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			Zone:        "cosmoshub-4",
			SrcZone:     "quicksilver-2",
			Proofs: []*cmtypes.Proof{
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    11012784,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeLiquidToken)],
				},
			},
		}

		suite.NoError(msgClaim.ValidateBasic())

		out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeLiquidToken].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
		suite.NoError(err)
		suite.Equal(out, math.ZeroInt())
	})
}

const (
	testLiquidKey      = "AhSWFVs4U+NLuO3A3gUTOjGbZc4maHVxYXRvbQ=="
	testLiquidData     = "MjQzMDc="
	testLiquidProofOps = `{"ops":[{"type":"ics23:iavl","key":"AhSWFVs4U+NLuO3A3gUTOjGbZc4maHVxYXRvbQ==","data":"CsYHChwCFJYVWzhT40u47cDeBRM6MZtlziZodXFhdG9tEgUyNDMwNxoOCAEYASABKgYAAszj7QYiLggBEgcCBNTz6wcgGiEgmgMjlAM00s2a2lenOIYuvAbMV+H9QCCqtj/hZMgA2soiLggBEgcEBvy3vAogGiEgqAWIJspAQtpxIiFrFTUMIJzOBDYVJJv9w13nZxOFIPgiLggBEgcGCvy3vAogGiEg5jmaiZEF+qDULa/ueYmD7Ab6iBr9oxH4xhxV/BxvNS4iLAgBEigKIvy3vAog1+g+9cSMxJZ/iPlJl/OmnuE23gt3Ft8BB20xt2YxklQgIi4IARIHDDL8t7wKIBohIFTr+vKs01Dc7510FZvfyb7AuTLwhEfKJBoekCmtQcR0IiwIARIoDnz8t7wKIG4SojIoWKVyzDmbkfXm1SxYih/w40xno1NBPHXSsvErICItCAESKRDWAfy3vAogSac61jxugtgqmFF7IL7TN22Iq53GtLtsIMerCTYrb5wgIi8IARIIEuAC/Le8CiAaISD1TPe5zCQ/dqTWgYX8tc9fDIN5AvQvJJIfDV1p1iCE4SItCAESKRSaBd6NwAogBJtfyw4wOcVmh7dCsJV+hg5XFbTGCZ7GisZQynJt1IkgIi8IARIIFuQK3o3ACiAaISCV2dejXG5G7yKFztK5tttizBRVwGLXmlI+JM4brOkMniItCAESKRjcEoSqwAog/e0Erfd6Ec6aYvvrmgUpZYBJuwu1CcPKcvaqmgA7xNggIi0IARIpGpwqhKrACiDl/8Iv0JSXb2OqAF5EPlQcZwj7dmT3fre7jGfqp3wInyAiLQgBEike9lbeqsAKIBSD+0MRBlI8Qy83zyOpPj/z6Z5hoANmlsvsOieaDJRMICIwCAESCSDEtAHeqsAKIBohIBkVO8zaaGddtVYeNMtYfMvsq9/a8qI8VpAsNHtcST8LIi4IARIqIqj1At6qwAogPKAJpaaoQ5eDHF9Zx3n3/DRohKwcDJENHj+Yd0UkNFwgIjAIARIJJKbKBd6qwAogGiEgGZr0fmisUw+WDXL4vyDIA48YQkYabjg8SYR1buDO5BsiMAgBEgkm9ooL4KrACiAaISDfRdB5aSVbwZ5I/9TF3bMKsflv+TpW0DFKHGtN3lHJ7iIuCAESKirw1hbgqsAKIGQ0rQv3DodvNZLrWVS5gQekxTi4BFGFN3vL333k1HcIICIwCAESCSyuoS7gqsAKIBohIK+Ocslu3BA5Q5XCDQt9orU8ZidCdT685eEG8jMF5yw5"},{"type":"ics23:simple","key":"YmFuaw==","data":"CvwBCgRiYW5rEiCZCW+QMlvWAlZAEOp/YTiIFCZzOSLlhAlUDQgcIt46UxoJCAEYASABKgEAIiUIARIhAWTNKVDOHlvI/MUZLTRSuaYC5kFh271dzwHT1hcT6axwIiUIARIhAY976trxkzU4XVP5v5ZUimhYDaOzhf8iRGUsUCuW0QC6IicIARIBARogcHyvc3OA+QXw6TFfaQnn8RqOpvf4XjzWrjoEx3MIoJQiJwgBEgEBGiA2auJ+PsYuKKCH9ZTpnfRn1MmFpBvoE16VBQvBJJbZ4SInCAESAQEaINnKOLCptf6IjVE2FPYsmvbaQMZvxG1zXvURHgKh/tHr"}]}`
)
