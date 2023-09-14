package keeper_test

import (
	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (s *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
				},
			},
			{
				Denom: "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
				},
			},
			{
				Denom: "factory/quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "quick1ve2nremzdnu7e55khlrt2282qhh98dh4708ppf",
				},
			},
		},
	}

	s.SetupTestForInitGenesis()
	app := s.App

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.BankKeeper.SetDenomMetaData(s.Ctx, banktypes.Metadata{Base: denom.GetDenom()})
		}
	}

	// check before initGenesis that the module account is nil
	tokenfactoryModuleAccount := app.AccountKeeper.GetAccount(s.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	app.TokenFactoryKeeper.SetParams(s.Ctx, types.Params{DenomCreationFee: sdk.Coins{sdk.NewInt64Coin("uqck", 100)}})
	app.TokenFactoryKeeper.InitGenesis(s.Ctx, genesisState)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(s.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)

	exportedGenesis := app.TokenFactoryKeeper.ExportGenesis(s.Ctx)
	s.Require().NotNil(exportedGenesis)
	s.Require().Equal(genesisState, *exportedGenesis)
}
