package keeper_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

var PKs = simapp.CreateTestPubKeys(10)

func (suite *KeeperTestSuite) TestStoreGetDeleteValidator() {
	suite.Run("validator - store / get / delete", func() {
		suite.SetupTest()
		suite.setupTestZones()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)

		validator := addressutils.GenerateValAddressForTest()

		valAddrBytes, err := addressutils.ValAddressFromBech32(validator.String(), zone.GetValoperPrefix())
		suite.NoError(err)
		_, found = app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		suite.False(found)

		count := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))

		newValidator := types.Validator{
			ValoperAddress:  validator.String(),
			CommissionRate:  sdk.NewDec(5.0),
			DelegatorShares: sdk.NewDec(1000.0),
			VotingPower:     sdk.NewInt(1000),
			Status:          stakingtypes.BondStatusBonded,
			Score:           sdk.NewDec(0),
		}
		err = app.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, newValidator)
		suite.NoError(err)

		count2 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))

		suite.Equal(count+1, count2)

		fetchedValidator, found := app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		suite.True(found)
		suite.Equal(newValidator, fetchedValidator)

		app.InterchainstakingKeeper.DeleteValidator(ctx, zone.ChainId, valAddrBytes)

		count3 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))
		suite.Equal(count, count3)
	})
}

func (suite *KeeperTestSuite) TestStoreGetDeleteValidatorByConsAddr() {
	suite.Run("validator - store / get / delete by consensus address", func() {
		suite.SetupTest()
		suite.setupTestZones()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)

		validator := addressutils.GenerateValAddressForTest()

		_, found = app.InterchainstakingKeeper.GetValidatorAddrByConsAddr(ctx, zone.ChainId, sdk.ConsAddress(PKs[0].Address()))
		suite.False(found)

		pkAny, err := codectypes.NewAnyWithValue(PKs[0])
		suite.Require().NoError(err)

		newValidator := types.Validator{
			ValoperAddress:  validator.String(),
			CommissionRate:  sdk.NewDec(5.0),
			DelegatorShares: sdk.NewDec(1000.0),
			VotingPower:     sdk.NewInt(1000),
			Status:          stakingtypes.BondStatusBonded,
			Score:           sdk.NewDec(0),
			ConsensusPubkey: pkAny,
		}

		err = app.InterchainstakingKeeper.SetValidatorAddrByConsAddr(ctx, zone.ChainId, newValidator)
		suite.NoError(err)

		_, found = app.InterchainstakingKeeper.GetValidatorAddrByConsAddr(ctx, zone.ChainId, sdk.ConsAddress(PKs[0].Address()))
		suite.True(found)

		app.InterchainstakingKeeper.DeleteValidatorAddrByConsAddr(ctx, zone.ChainId, sdk.ConsAddress(PKs[0].Address()))
		_, found = app.InterchainstakingKeeper.GetValidatorAddrByConsAddr(ctx, zone.ChainId, sdk.ConsAddress(PKs[0].Address()))
		suite.False(found)
	})
}
