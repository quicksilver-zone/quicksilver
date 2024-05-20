package keeper_test

import (
	"github.com/quicksilver-zone/quicksilver/utils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (suite *KeeperTestSuite) TestAfterZoneCreated() {
	zone := &icstypes.Zone{
		ChainId:         "testzone-1",
		ConnectionId:    "connection-0",
		AccountPrefix:   "test",
		LocalDenom:      "uqtst",
		BaseDenom:       "utst",
		TransferChannel: "channel-1",
	}
	suite.Run("ProtocolData", func() {
		suite.setupChannelForHookTest()
		suite.setupTestProtocolData()
		k := suite.GetQuicksilverApp(suite.chainA).ParticipationRewardsKeeper
		ctx := suite.chainA.GetContext()
		suite.NoError(k.Hooks().AfterZoneCreated(ctx, zone))
		ctx = suite.chainA.GetContext()
		// we want to fetch this for the local (quicksilver) zone; not the host chain.
		pd, found := k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, ctx.ChainID()+"_"+zone.LocalDenom)
		suite.True(found)
		upd, err := types.UnmarshalProtocolData(types.ProtocolDataTypeLiquidToken, pd.Data)
		suite.NoError(err)

		lpd := upd.(*types.LiquidAllowedDenomProtocolData)
		suite.Equal(zone.ChainId, lpd.RegisteredZoneChainID)
		suite.Equal(zone.LocalDenom, lpd.QAssetDenom)
		suite.Equal(ctx.ChainID(), lpd.ChainID)

		// check host chain
		hostChainIBCDenom := utils.DeriveIbcDenom("transfer", counterpartyTestzoneChannel, "", "", "uqtst")
		pd, found = k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, zone.ChainId+"_"+hostChainIBCDenom)
		suite.True(found)
		suite.NotNil(pd)
		upd, err = types.UnmarshalProtocolData(types.ProtocolDataTypeLiquidToken, pd.Data)
		suite.NoError(err)
		lpd = upd.(*types.LiquidAllowedDenomProtocolData)
		suite.Equal(zone.ChainId, lpd.RegisteredZoneChainID)
		suite.Equal(zone.LocalDenom, lpd.QAssetDenom)
		suite.Equal(hostChainIBCDenom, lpd.IbcDenom)

		// check osmosis
		osmosisIBCDenom := utils.DeriveIbcDenom("transfer", counterpartyOsmosisChannel, "", "", "uqtst")
		pd, found = k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, "osmosis-1_"+osmosisIBCDenom)
		suite.True(found)
		upd, err = types.UnmarshalProtocolData(types.ProtocolDataTypeLiquidToken, pd.Data)
		suite.NoError(err)
		lpd = upd.(*types.LiquidAllowedDenomProtocolData)
		suite.Equal(zone.ChainId, lpd.RegisteredZoneChainID)
		suite.Equal(zone.LocalDenom, lpd.QAssetDenom)
		suite.Equal(osmosisIBCDenom, lpd.IbcDenom)

		// check umee
		umeeIBCDenom := utils.DeriveIbcDenom("transfer", counterpartyUmeeChannel, "", "", "uqtst")
		pd, found = k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, "umee-types-1_"+umeeIBCDenom)
		suite.True(found)
		upd, err = types.UnmarshalProtocolData(types.ProtocolDataTypeLiquidToken, pd.Data)
		suite.NoError(err)
		lpd = upd.(*types.LiquidAllowedDenomProtocolData)
		suite.Equal(zone.ChainId, lpd.RegisteredZoneChainID)
		suite.Equal(zone.LocalDenom, lpd.QAssetDenom)
		suite.Equal(umeeIBCDenom, lpd.IbcDenom)
	})
}
