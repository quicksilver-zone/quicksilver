package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/simapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
			ValoperAddress:      validator.String(),
			CommissionRate:      sdk.NewDec(5.0),
			DelegatorShares:     sdk.NewDec(1000.0),
			VotingPower:         sdkmath.NewInt(1000),
			Status:              stakingtypes.BondStatusBonded,
			Score:               sdk.NewDec(0),
			LiquidShares:        sdk.ZeroDec(),
			ValidatorBondShares: sdk.ZeroDec(),
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

		newValidator := stakingtypes.Validator{
			OperatorAddress: validator.String(),
			ConsensusPubkey: pkAny,
		}
		consAddr, err := newValidator.GetConsAddr()
		suite.NoError(err)

		app.InterchainstakingKeeper.SetValidatorAddrByConsAddr(ctx, zone.ChainId, newValidator.OperatorAddress, consAddr)

		_, found = app.InterchainstakingKeeper.GetValidatorAddrByConsAddr(ctx, zone.ChainId, sdk.ConsAddress(PKs[0].Address()))
		suite.True(found)

		app.InterchainstakingKeeper.DeleteValidatorAddrByConsAddr(ctx, zone.ChainId, sdk.ConsAddress(PKs[0].Address()))
		_, found = app.InterchainstakingKeeper.GetValidatorAddrByConsAddr(ctx, zone.ChainId, sdk.ConsAddress(PKs[0].Address()))
		suite.False(found)
	})
}

func (suite *KeeperTestSuite) TestGetActiveValidators() {
	suite.Run("active validators", func() {
		suite.SetupTest()
		suite.setupTestZones()

		app := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)

		validators := app.InterchainstakingKeeper.GetActiveValidators(ctx, "not a chain id")
		suite.Len(validators, 0)

		validators = app.InterchainstakingKeeper.GetActiveValidators(ctx, zone.ChainId)
		count := len(validators)

		validator1 := addressutils.GenerateValAddressForTest()
		validator2 := addressutils.GenerateValAddressForTest()

		newValidator1 := types.Validator{
			ValoperAddress:  validator1.String(),
			CommissionRate:  sdk.NewDec(5.0),
			DelegatorShares: sdk.NewDec(1000.0),
			VotingPower:     sdkmath.NewInt(1000),
			Status:          stakingtypes.BondStatusBonded,
			Score:           sdk.NewDec(0),
		}
		newValidator2 := newValidator1
		newValidator2.ValoperAddress = validator2.String()
		newValidator2.Status = stakingtypes.BondStatusUnbonded

		err := app.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, newValidator1)
		suite.NoError(err)

		err = app.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, newValidator2)
		suite.NoError(err)

		validators = app.InterchainstakingKeeper.GetActiveValidators(ctx, zone.ChainId)
		count2 := len(validators)

		suite.Equal(count+1, count2)
	})
}
