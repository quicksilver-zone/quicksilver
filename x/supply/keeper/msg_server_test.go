package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	minttypes "github.com/quicksilver-zone/quicksilver/x/mint/types"
	"github.com/quicksilver-zone/quicksilver/x/supply/keeper"
	"github.com/quicksilver-zone/quicksilver/x/supply/types"
)

func (suite *KeeperTestSuite) Test_msgServer_IncentivePoolSpend() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	modAccAddr := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
	userAddress := addressutils.GenerateAccAddressForTest().String()
	denom := "uatom" // same as test zone setup in keeper_test
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromUint64(1000)))
	mintCoins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromUint64(100000000)))

	// set up mod acct with funds
	err := appA.BankKeeper.MintCoins(suite.chainA.GetContext(), minttypes.ModuleName, mintCoins)
	suite.Require().NoError(err)
	err = appA.BankKeeper.SendCoinsFromModuleToModule(suite.chainA.GetContext(), minttypes.ModuleName, types.AirdropAccount, mintCoins)
	suite.Require().NoError(err)

	msg := types.MsgIncentivePoolSpend{}
	tests := []struct {
		name     string
		malleate func()
		want     *types.MsgIncentivePoolSpendResponse
		wantErr  bool
	}{
		{
			name: "invalid authority",
			malleate: func() {
				msg = types.MsgIncentivePoolSpend{
					Authority:   "invalid",
					ToAddress:   userAddress,
					Amount:      coins,
					Title:       "Invalid Incentive Pool Spend Title",
					Description: "Invalid Incentive Pool Spend Description",
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid",
			malleate: func() {
				msg = types.MsgIncentivePoolSpend{
					Authority:   modAccAddr,
					ToAddress:   userAddress,
					Amount:      coins,
					Title:       "Valid Incentive Pool Spend Title",
					Description: "Valid Incentive Pool Spend Description",
				}
			},
			want:    &types.MsgIncentivePoolSpendResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()

			k := keeper.NewMsgServerImpl(&appA.SupplyKeeper)
			resp, err := k.IncentivePoolSpend(sdk.WrapSDKContext(suite.chainA.GetContext()), &msg)
			if tt.wantErr {
				suite.Require().Error(err)
				suite.Require().Nil(resp)
				suite.T().Logf("Error: %v", err)
				return
			}

			suite.Require().NoError(err)
			suite.Require().NotNil(resp)
			suite.Require().Equal(tt.want, resp)

			// verify that balance has been properly transferred
			accAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
			suite.Require().NoError(err)
			balance := appA.BankKeeper.GetAllBalances(suite.chainA.GetContext(), accAddr)
			suite.Require().Equal(msg.Amount, balance)
		})
	}
}
