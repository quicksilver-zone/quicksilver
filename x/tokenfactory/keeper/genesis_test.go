package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/tokenfactory/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
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

	suite.SetupTestForInitGenesis()
	app := suite.App

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.BankKeeper.SetDenomMetaData(suite.Ctx, banktypes.Metadata{Base: denom.GetDenom()})
		}
	}

	// check before initGenesis that the module account is nil
	tokenfactoryModuleAccount := app.AccountKeeper.GetAccount(suite.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().Nil(tokenfactoryModuleAccount)

	app.TokenFactoryKeeper.SetParams(suite.Ctx, types.Params{DenomCreationFee: sdk.Coins{sdk.NewInt64Coin("uqck", 100)}})
	app.TokenFactoryKeeper.InitGenesis(suite.Ctx, genesisState)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(suite.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().NotNil(tokenfactoryModuleAccount)

	exportedGenesis := app.TokenFactoryKeeper.ExportGenesis(suite.Ctx)
	suite.Require().NotNil(exportedGenesis)
	suite.Require().Equal(genesisState, *exportedGenesis)
}
