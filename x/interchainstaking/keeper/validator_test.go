package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (s *KeeperTestSuite) TestStoreGetDeleteValidator() {
	s.Run("validator - store / get / delete", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		ctx := s.chainA.GetContext()

		zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().True(found)

		validator := utils.GenerateValAddressForTest()

		valAddrBytes, err := utils.ValAddressFromBech32(validator.String(), zone.GetValoperPrefix())
		s.Require().NoError(err)
		_, found = app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		s.Require().False(found)

		count := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))

		newValidator := types.Validator{
			ValoperAddress:  validator.String(),
			CommissionRate:  sdk.NewDec(5.0),
			DelegatorShares: sdk.NewDec(1000.0),
			VotingPower:     sdk.NewInt(1000),
			Status:          stakingtypes.BondStatusBonded,
			Score:           sdk.NewDec(0),
		}
		app.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, newValidator)

		count2 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))

		s.Require().Equal(count+1, count2)

		fetchedValidator, found := app.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
		s.Require().True(found)
		s.Require().Equal(newValidator, fetchedValidator)

		app.InterchainstakingKeeper.DeleteValidator(ctx, zone.ChainId, valAddrBytes)

		count3 := len(app.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId))
		s.Require().Equal(count, count3)
	})
}
