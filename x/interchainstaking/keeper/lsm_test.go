package keeper_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

const BondStatusUnbonded = "BOND_STATUS_UNBONDED"

func (suite *KeeperTestSuite) TestLsmSetGetDelete() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	caps, found := icsKeeper.GetLsmCaps(ctx, suite.chainB.ChainID)
	suite.False(found)
	suite.Nil(caps)

	allCaps := icsKeeper.AllLsmCaps(ctx)
	suite.Equal(0, len(allCaps))

	icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID, types.LsmCaps{ValidatorCap: sdk.NewDecWithPrec(50, 2), GlobalCap: sdk.NewDecWithPrec(25, 2), ValidatorBondCap: sdk.NewDec(500)})

	caps, found = icsKeeper.GetLsmCaps(ctx, suite.chainB.ChainID)
	suite.True(found)
	suite.Equal(caps.ValidatorBondCap, sdk.NewDec(500))

	allCaps = icsKeeper.AllLsmCaps(ctx)
	suite.Equal(1, len(allCaps))

	icsKeeper.DeleteLsmCaps(ctx, suite.chainB.ChainID)

	caps, found = icsKeeper.GetLsmCaps(ctx, suite.chainB.ChainID)
	suite.False(found)
	suite.Nil(caps)

	allCaps = icsKeeper.AllLsmCaps(ctx)
	suite.Equal(0, len(allCaps))
}

func (suite *KeeperTestSuite) TestGetTotalStakedSupply() {
	suite.SetupTest()
	suite.setupTestZones()
	tcs := []struct {
		Name     string
		Malleate func(icsKeeper *keeper.Keeper)
		Expect   math.Int
	}{
		{
			Name: "4x 1000000 VP bonded",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000000)
				validators[1].VotingPower = math.NewInt(1000000)
				validators[2].VotingPower = math.NewInt(1000000)
				validators[3].VotingPower = math.NewInt(1000000)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			Expect: sdk.NewInt(4000000),
		},
		{
			Name: "3x 1000000 VP bonded, 1x 1000000 unbonded",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000000)
				validators[1].VotingPower = math.NewInt(1000000)
				validators[2].VotingPower = math.NewInt(1000000)
				validators[3].VotingPower = math.NewInt(1000000)
				validators[3].Status = BondStatusUnbonded
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			Expect: sdk.NewInt(3000000),
		},
		{
			Name: "different vps, total 10000000",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(5000000)
				validators[1].VotingPower = math.NewInt(3000000)
				validators[2].VotingPower = math.NewInt(2000000)
				validators[3].VotingPower = math.NewInt(1000000)
				validators[3].Status = BondStatusUnbonded
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			Expect: sdk.NewInt(10000000),
		},
	}
	for _, t := range tcs {
		suite.Run(t.Name, func() {
			icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
			ctx := suite.chainA.GetContext()
			t.Malleate(icsKeeper)
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			suite.Equal(icsKeeper.GetTotalStakedSupply(ctx, &zone), t.Expect)
		})
	}
}

func (suite *KeeperTestSuite) TestGetLiquidStakedSupply() {
	tcs := []struct {
		Name     string
		Malleate func(icsKeeper *keeper.Keeper)
		Expect   sdk.Dec
	}{
		{
			Name: "4x 1000000 VP bonded, 0 liquid",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.ZeroDec()
				validators[2].LiquidShares = sdk.ZeroDec()
				validators[3].LiquidShares = sdk.ZeroDec()
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			Expect: sdk.ZeroDec(),
		},
		{
			Name: "3x 1000000 VP bonded, 1x 1000000 unbonded",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.NewDec(5000)
				validators[2].LiquidShares = sdk.NewDec(5000)
				validators[3].LiquidShares = sdk.ZeroDec()
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			Expect: sdk.NewDec(10000),
		},
		{
			Name: "different vps, total 10000000",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].LiquidShares = sdk.NewDec(1000)
				validators[1].LiquidShares = sdk.NewDec(2000)
				validators[2].LiquidShares = sdk.NewDec(3000)
				validators[3].LiquidShares = sdk.NewDec(5000)
				validators[3].Status = BondStatusUnbonded
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			Expect: sdk.NewDec(6000),
		},
	}
	for _, t := range tcs {
		suite.Run(t.Name, func() {
			suite.SetupTest()
			suite.setupTestZones()
			icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
			ctx := suite.chainA.GetContext()
			t.Malleate(icsKeeper)
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			suite.Equal(icsKeeper.GetLiquidStakedSupply(ctx, &zone), t.Expect)
		})
	}
}

func (suite *KeeperTestSuite) TestCheckExceedsGlobalCap() {
	tcs := []struct {
		Name     string
		Malleate func(icsKeeper *keeper.Keeper)
		Expect   bool
	}{
		{
			Name: "cap 5%, liquid 2% + 1; expect false",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[1].VotingPower = math.NewInt(1000)
				validators[2].VotingPower = math.NewInt(1000)
				validators[3].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.NewDec(80)
				validators[2].LiquidShares = sdk.ZeroDec()
				validators[3].LiquidShares = sdk.ZeroDec()
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
				icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID,
					types.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(500),
						GlobalCap:        sdk.NewDecWithPrec(5, 2),
					})
			},
			Expect: false,
		},
		{
			Name: "cap 5%, liquid 5% + 1; expect true",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[1].VotingPower = math.NewInt(1000)
				validators[2].VotingPower = math.NewInt(1000)
				validators[3].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.NewDec(60)
				validators[2].LiquidShares = sdk.NewDec(60)
				validators[3].LiquidShares = sdk.NewDec(80)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
				icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID,
					types.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(500),
						GlobalCap:        sdk.NewDecWithPrec(5, 2),
					})
			},
			Expect: true,
		},
		{
			Name: "no cap set, expect false",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[1].VotingPower = math.NewInt(1000)
				validators[2].VotingPower = math.NewInt(1000)
				validators[3].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.NewDec(20)
				validators[2].LiquidShares = sdk.NewDec(20)
				validators[3].LiquidShares = sdk.NewDec(10)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			Expect: false,
		},
	}
	for _, t := range tcs {
		suite.Run(t.Name, func() {
			suite.SetupTest()
			suite.setupTestZones()
			icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
			ctx := suite.chainA.GetContext()
			t.Malleate(icsKeeper)
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			suite.Equal(t.Expect, icsKeeper.CheckExceedsGlobalCap(ctx, &zone, sdk.NewInt(1)))
		})
	}
}

func (suite *KeeperTestSuite) TestCheckExceedsValidatorCap() {
	tcs := []struct {
		Name      string
		Malleate  func(icsKeeper *keeper.Keeper)
		ExpectErr bool
	}{
		{
			Name: "cap 50%, liquid 2% + 1; expect false",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[1].VotingPower = math.NewInt(1000)
				validators[2].VotingPower = math.NewInt(1000)
				validators[3].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.NewDec(20)
				validators[2].LiquidShares = sdk.ZeroDec()
				validators[3].LiquidShares = sdk.ZeroDec()
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))

				icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID,
					types.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(500),
						GlobalCap:        sdk.NewDecWithPrec(5, 2),
					})
			},
			ExpectErr: false,
		},
		{
			Name: "cap 50%, liquid 60% + 1; expect true",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[1].VotingPower = math.NewInt(1000)
				validators[2].VotingPower = math.NewInt(1000)
				validators[3].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.NewDec(600)
				validators[2].LiquidShares = sdk.NewDec(20)
				validators[3].LiquidShares = sdk.NewDec(10)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))

				icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID,
					types.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(500),
						GlobalCap:        sdk.NewDecWithPrec(5, 2),
					})
			},
			ExpectErr: true,
		},
		{
			Name: "no cap set, expect false",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[1].VotingPower = math.NewInt(1000)
				validators[2].VotingPower = math.NewInt(1000)
				validators[3].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.ZeroDec()
				validators[1].LiquidShares = sdk.NewDec(600)
				validators[2].LiquidShares = sdk.NewDec(20)
				validators[3].LiquidShares = sdk.NewDec(10)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			ExpectErr: false,
		},
	}
	for _, t := range tcs {
		suite.Run(t.Name, func() {
			suite.SetupTest()
			suite.setupTestZones()
			icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
			ctx := suite.chainA.GetContext()
			t.Malleate(icsKeeper)
			validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			if t.ExpectErr {
				suite.Error(icsKeeper.CheckExceedsValidatorCap(ctx, &zone, validators[1].ValoperAddress, sdk.NewInt(1)))
			} else {
				suite.NoError(icsKeeper.CheckExceedsValidatorCap(ctx, &zone, validators[1].ValoperAddress, sdk.NewInt(1)))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCheckExceedsValidatorBondCap() {
	tcs := []struct {
		Name      string
		Malleate  func(icsKeeper *keeper.Keeper)
		ExpectErr bool
	}{
		{
			Name: "valbond 5, multiplier 100, ls 400; expect false",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.NewDec(400)
				validators[0].ValidatorBondShares = sdk.NewDec(5)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))

				icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID,
					types.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(100),
						GlobalCap:        sdk.NewDecWithPrec(5, 2),
					})
			},
			ExpectErr: false,
		},
		{
			Name: "valbond 5, multiplier 100, ls 500; expect true",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.NewDec(500)
				validators[0].ValidatorBondShares = sdk.NewDec(5)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))

				icsKeeper.SetLsmCaps(ctx, suite.chainB.ChainID,
					types.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(100),
						GlobalCap:        sdk.NewDecWithPrec(5, 2),
					})
			},
			ExpectErr: true,
		},
		{
			Name: "no cap set, expect false",
			Malleate: func(icsKeeper *keeper.Keeper) {
				ctx := suite.chainA.GetContext()
				validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
				validators[0].VotingPower = math.NewInt(1000)
				validators[0].LiquidShares = sdk.NewDec(500)
				validators[0].ValidatorBondShares = sdk.NewDec(5)
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[0]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[1]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[2]))
				suite.NoError(icsKeeper.SetValidator(ctx, suite.chainB.ChainID, validators[3]))
			},
			ExpectErr: false,
		},
	}
	for _, t := range tcs {
		suite.Run(t.Name, func() {
			suite.SetupTest()
			suite.setupTestZones()
			icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
			ctx := suite.chainA.GetContext()
			t.Malleate(icsKeeper)
			zone, found := icsKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			validators := icsKeeper.GetValidators(ctx, suite.chainB.ChainID)
			if t.ExpectErr {
				suite.Error(icsKeeper.CheckExceedsValidatorBondCap(ctx, &zone, validators[0].ValoperAddress, math.NewInt(1)))
			} else {
				suite.NoError(icsKeeper.CheckExceedsValidatorBondCap(ctx, &zone, validators[0].ValoperAddress, math.NewInt(1)))
			}
		})
	}
}
