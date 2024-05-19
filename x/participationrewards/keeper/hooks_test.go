package keeper_test

import (
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
		pd, found = k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, zone.ChainId+"_"+zone.BaseDenom)
		suite.True(found)
		upd, err = types.UnmarshalProtocolData(types.ProtocolDataTypeLiquidToken, pd.Data)
		suite.NoError(err)
		lpd = upd.(*types.LiquidAllowedDenomProtocolData)
		suite.Equal(zone.ChainId, lpd.RegisteredZoneChainID)
		suite.Equal(zone.LocalDenom, lpd.QAssetDenom)
		suite.Equal(zone.BaseDenom, lpd.IbcDenom)

		// TODO: add tests for osmosis zone and umee zone, need setup connection data for testing
	})
}
