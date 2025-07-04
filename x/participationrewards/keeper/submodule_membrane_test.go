package keeper_test

import (
	"encoding/base64"
	"encoding/json"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/proto/tendermint/crypto"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) Test_Membrane_ValidateClaim() {
	cases := []struct {
		name          string
		submitAddress string
		mappedAddress string
		expectErr     bool
	}{
		{
			"valid - submitter",
			"quick16qqhmsqcs4j6mfa92flnz4n8tj2s53jwdhy7an",
			"",
			false,
		},
		{
			"valid - mapped",
			"quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			"osmo16qqhmsqcs4j6mfa92flnz4n8tj2s53jwwg8ujn",
			false,
		},
		{
			"invalid - no mapped",
			"quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
			"",
			true,
		},
		{
			"invalid - mapped no match",
			"quick1jc24kwznud9m3mwqmcz3xw33ndjuufngu5m0y6",
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

			membranePd := types.MembraneProtocolData{
				ContractAddress: "osmo1gy5gpqqlth0jpm9ydxlmff6g5mpnfvrfxd3mfc8dhyt03waumtzqt8exxr",
			}

			blob, err = json.Marshal(membranePd)
			suite.NoError(err)
			pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeMembraneParams)], blob)
			app.ParticipationRewardsKeeper.SetProtocolData(ctx, membranePd.GenerateKey(), pd)

			liquidTokenPd := types.LiquidAllowedDenomProtocolData{
				ChainID:               "osmosis-1",
				RegisteredZoneChainID: "cosmoshub-4",
				QAssetDenom:           "uqatom",
				IbcDenom:              "ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC",
			}

			blob, err = json.Marshal(liquidTokenPd)
			suite.NoError(err)
			pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
			app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

			if c.mappedAddress != "" {
				app.InterchainstakingKeeper.SetRemoteAddressMap(ctx, addressutils.MustAccAddressFromBech32(c.submitAddress, ""), addressutils.MustAccAddressFromBech32(c.mappedAddress, ""), "osmosis-1")
			}

			// claim
			key, err := base64.StdEncoding.DecodeString(testMembraneKey)
			suite.NoError(err)
			data, err := base64.StdEncoding.DecodeString(testMembraneData)
			suite.NoError(err)
			proofOps := crypto.ProofOps{}
			err = json.Unmarshal([]byte(testMembraneProofOps), &proofOps)
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
						Height:    38973687,
						ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeMembrane)],
					},
				},
			}

			suite.NoError(msgClaim.ValidateBasic())

			out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeMembrane].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
			if c.expectErr {
				suite.Error(err)
				suite.Equal(out, math.ZeroInt())
			} else {
				suite.NoError(err)
				suite.Equal(out, math.NewInt(103200))
			}
		})
	}
}

func (suite *KeeperTestSuite) Test_Membrane_ValidateClaim_IgnoreDoubleClaim() {
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

		membranePd := types.MembraneProtocolData{
			ContractAddress: "osmo1gy5gpqqlth0jpm9ydxlmff6g5mpnfvrfxd3mfc8dhyt03waumtzqt8exxr",
		}

		blob, err = json.Marshal(membranePd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeMembraneParams)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, membranePd.GenerateKey(), pd)

		liquidTokenPd := types.LiquidAllowedDenomProtocolData{
			ChainID:               "osmosis-1",
			RegisteredZoneChainID: "cosmoshub-4",
			QAssetDenom:           "uqatom",
			IbcDenom:              "ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC",
		}

		blob, err = json.Marshal(liquidTokenPd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

		// claim
		key, err := base64.StdEncoding.DecodeString(testMembraneKey)
		suite.NoError(err)
		data, err := base64.StdEncoding.DecodeString(testMembraneData)
		suite.NoError(err)
		proofOps := crypto.ProofOps{}
		err = json.Unmarshal([]byte(testMembraneProofOps), &proofOps)
		suite.NoError(err)

		msgClaim := types.MsgSubmitClaim{
			UserAddress: "quick16qqhmsqcs4j6mfa92flnz4n8tj2s53jwdhy7an",
			Zone:        "cosmoshub-4",
			SrcZone:     "osmosis-1",
			Proofs: []*cmtypes.Proof{
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    38973687,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeMembrane)],
				},
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    38973687,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeMembrane)],
				},
			},
		}

		suite.NoError(msgClaim.ValidateBasic())

		out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeMembrane].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
		suite.NoError(err)
		suite.Equal(out, math.NewInt(103200))
	})
}

func (suite *KeeperTestSuite) Test_Membrane_ValidateClaim_FailContractMismatch() {
	suite.Run("fail_contract_mismatch", func() {
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

		membranePd := types.MembraneProtocolData{
			ContractAddress: "osmo1gy5gpqqlth0jpm9ydxlmff6g5mpnfvrfd3mfc8dhyt03waumtzqt8exxr",
		}

		blob, err = json.Marshal(membranePd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeMembraneParams)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, membranePd.GenerateKey(), pd)

		liquidTokenPd := types.LiquidAllowedDenomProtocolData{
			ChainID:               "osmosis-1",
			RegisteredZoneChainID: "cosmoshub-4",
			QAssetDenom:           "uqatom",
			IbcDenom:              "ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC",
		}

		blob, err = json.Marshal(liquidTokenPd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeLiquidToken)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, liquidTokenPd.GenerateKey(), pd)

		// claim
		key, err := base64.StdEncoding.DecodeString(testMembraneKey)
		suite.NoError(err)
		data, err := base64.StdEncoding.DecodeString(testMembraneData)
		suite.NoError(err)
		proofOps := crypto.ProofOps{}
		err = json.Unmarshal([]byte(testMembraneProofOps), &proofOps)
		suite.NoError(err)

		msgClaim := types.MsgSubmitClaim{
			UserAddress: "quick16qqhmsqcs4j6mfa92flnz4n8tj2s53jwdhy7an",
			Zone:        "cosmoshub-4",
			SrcZone:     "osmosis-1",
			Proofs: []*cmtypes.Proof{
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    38973687,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeMembrane)],
				},
			},
		}

		suite.NoError(msgClaim.ValidateBasic())

		out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeMembrane].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
		suite.Error(err)
		suite.Equal(out, math.ZeroInt())
	})
}

func (suite *KeeperTestSuite) Test_Membrane_ValidateClaim_IgnoreUnclaimableDenom() {
	suite.Run("ignore_unclaimable_denom", func() {
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

		membranePd := types.MembraneProtocolData{
			ContractAddress: "osmo1gy5gpqqlth0jpm9ydxlmff6g5mpnfvrfxd3mfc8dhyt03waumtzqt8exxr",
		}

		blob, err = json.Marshal(membranePd)
		suite.NoError(err)
		pd = types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeMembraneParams)], blob)
		app.ParticipationRewardsKeeper.SetProtocolData(ctx, membranePd.GenerateKey(), pd)

		// claim
		key, err := base64.StdEncoding.DecodeString(testMembraneKey)
		suite.NoError(err)
		data, err := base64.StdEncoding.DecodeString(testMembraneData)
		suite.NoError(err)
		proofOps := crypto.ProofOps{}
		err = json.Unmarshal([]byte(testMembraneProofOps), &proofOps)
		suite.NoError(err)

		msgClaim := types.MsgSubmitClaim{
			UserAddress: "quick16qqhmsqcs4j6mfa92flnz4n8tj2s53jwdhy7an",
			Zone:        "cosmoshub-4",
			SrcZone:     "osmosis-1",
			Proofs: []*cmtypes.Proof{
				{
					Key:       key,
					Data:      data,
					ProofOps:  &proofOps,
					Height:    38973687,
					ProofType: cmtypes.ClaimType_name[int32(cmtypes.ClaimTypeMembrane)],
				},
			},
		}

		suite.NoError(msgClaim.ValidateBasic())

		out, err := app.ParticipationRewardsKeeper.PrSubmodules[cmtypes.ClaimTypeMembrane].ValidateClaim(ctx, app.ParticipationRewardsKeeper, &msgClaim)
		suite.NoError(err)
		suite.Equal(out, math.ZeroInt())
	})
}

const (
	testMembraneKey      = "A0EogIAfXd8g7KRpv7SnSKbDNLBpM2O04O25Fvi7vNrEAAlwb3NpdGlvbnNvc21vMTZxcWhtc3FjczRqNm1mYTkyZmxuejRuOHRqMnM1M2p3d2c4dWpu"
	testMembraneData     = "W3sicG9zaXRpb25faWQiOiI1MDQiLCJjb2xsYXRlcmFsX2Fzc2V0cyI6W3siYXNzZXQiOnsiaW5mbyI6eyJuYXRpdmVfdG9rZW4iOnsiZGVub20iOiJpYmMvNDJEMjQ4NzlENDU2OUNFNjQ3N0I3RTg4MjA2QURCRkU0N0MyMjJDNkNBRDUxQTU0MDgzRTRBNzI1OTQyNjlGQyJ9fSwiYW1vdW50IjoiMTAzMjAwIn0sIm1heF9ib3Jyb3dfTFRWIjoiMC40NSIsIm1heF9MVFYiOiIwLjUiLCJyYXRlX2luZGV4IjoiMS4wMTE3OTkyODQ4OTMxODcwNjkiLCJwb29sX2luZm8iOm51bGwsImhpa2VfcmF0ZXMiOm51bGx9XSwiY3JlZGl0X2Ftb3VudCI6IjAifV0="
	testMembraneProofOps = `{"ops":[{"type":"ics23:iavl","key":"A0EogIAfXd8g7KRpv7SnSKbDNLBpM2O04O25Fvi7vNrEAAlwb3NpdGlvbnNvc21vMTZxcWhtc3FjczRqNm1mYTkyZmxuejRuOHRqMnM1M2p3d2c4dWpu","data":"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"},{"type":"ics23:simple","key":"d2FzbQ==","data":"Cs8BCgR3YXNtEiDo+wTkcVnz1m+B8dis/42DRh9fCF9A+fjjTAOSnLbTBRoJCAEYASABKgEAIiUIARIhAVirLdYOjvUtrVftQ+mUFW2K8Ncclf50prIrUKy3zYwSIiUIARIhAXJTfzGWZ1g2udZGlNuyhSij8Tc6Bmri7i8Y/c5JCtvnIiUIARIhAcpV7I8vwXOou+H55ABrhUK38yG6kaEHwu0mKt+DyTL5IiUIARIhAfENd7bGDqzjX8LTIIW8dlG/2rUlozsyfONY+yVzGiJs"}]}`
)
