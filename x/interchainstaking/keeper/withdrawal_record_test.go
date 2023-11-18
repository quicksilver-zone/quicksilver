package keeper_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icskeeper "github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestUpdateWithdrawalRecordsForSlash() {
	tcs := []struct {
		Name            string
		InitialRecords  func(ctx sdk.Context, keeper *icskeeper.Keeper)
		ExpectedRecords func(ctx sdk.Context, keeper *icskeeper.Keeper) (out []types.WithdrawalRecord)
		Validator       func(validators []string) string
		Delta           sdk.Dec
		ExpectError     bool
	}{
		{
			Name: "single 5% slashing",
			InitialRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				withdrawalRecord := types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: zone.DelegationAddress.Address,
					Recipient: "cosmos1v4gek4mld0k5yhpe0fsln4takg558cdpyexv2rxr3dh45f2fqrgsw52m97",
					Distribution: []*types.Distribution{
						{Amount: 10000, Valoper: validators[0]},
						{Amount: 10000, Valoper: validators[1]},
						{Amount: 10000, Valoper: validators[2]},
						{Amount: 10000, Valoper: validators[3]},
					},
					Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(40000))),
					BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(32356)),
					Txhash:     "3BE21C1057ABBFBC44BE8993D2A4701C751507FF9901AA110B5993BA070C176B",
					Status:     types.WithdrawStatusUnbond,
				}
				keeper.SetWithdrawalRecord(ctx, withdrawalRecord)
			},
			ExpectedRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) (out []types.WithdrawalRecord) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				out = []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: zone.DelegationAddress.Address,
						Recipient: "cosmos1v4gek4mld0k5yhpe0fsln4takg558cdpyexv2rxr3dh45f2fqrgsw52m97",
						Distribution: []*types.Distribution{
							{Amount: 10000, Valoper: validators[0]},
							{Amount: 9500, Valoper: validators[1]},
							{Amount: 10000, Valoper: validators[2]},
							{Amount: 10000, Valoper: validators[3]},
						},
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(39500))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(32356)),
						Txhash:     "3BE21C1057ABBFBC44BE8993D2A4701C751507FF9901AA110B5993BA070C176B",
						Status:     types.WithdrawStatusUnbond,
					},
				}
				return out
			},
			Delta:     sdk.NewDecWithPrec(100, 2).Quo(sdk.NewDecWithPrec(95, 2)),
			Validator: func(validators []string) string { return validators[1] },

			ExpectError: false,
		},
		{
			Name: "multi record 5% slashing",
			InitialRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				withdrawalRecord := types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: zone.DelegationAddress.Address,
					Recipient: "cosmos1v4gek4mld0k5yhpe0fsln4takg558cdpyexv2rxr3dh45f2fqrgsw52m97",
					Distribution: []*types.Distribution{
						{Amount: 10000, Valoper: validators[0]},
						{Amount: 10000, Valoper: validators[1]},
						{Amount: 10000, Valoper: validators[2]},
						{Amount: 10000, Valoper: validators[3]},
					},
					Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(40000))),
					BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(32356)),
					Txhash:     "3BE21C1057ABBFBC44BE8993D2A4701C751507FF9901AA110B5993BA070C176B",
					Status:     types.WithdrawStatusUnbond,
				}
				keeper.SetWithdrawalRecord(ctx, withdrawalRecord)
				withdrawalRecord = types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: zone.DelegationAddress.Address,
					Recipient: "cosmos1nvkpj5n5mhy2ntvgn2cklntwx9ujvfvcacz5et",
					Distribution: []*types.Distribution{
						{Amount: 13000, Valoper: validators[0]},
						{Amount: 14000, Valoper: validators[1]},
						{Amount: 10000, Valoper: validators[2]},
						{Amount: 11000, Valoper: validators[3]},
					},
					Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(48000))),
					BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(36503)),
					Txhash:     "FB087A50A4836CBDFACA70D393AF110C28935276267B7BA2838BE3CEEA08F762",
					Status:     types.WithdrawStatusUnbond,
				}
				keeper.SetWithdrawalRecord(ctx, withdrawalRecord)
			},
			ExpectedRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) (out []types.WithdrawalRecord) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				out = []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: zone.DelegationAddress.Address,
						Recipient: "cosmos1v4gek4mld0k5yhpe0fsln4takg558cdpyexv2rxr3dh45f2fqrgsw52m97",
						Distribution: []*types.Distribution{
							{Amount: 10000, Valoper: validators[0]},
							{Amount: 9500, Valoper: validators[1]},
							{Amount: 10000, Valoper: validators[2]},
							{Amount: 10000, Valoper: validators[3]},
						},
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(39500))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(32356)),
						Txhash:     "3BE21C1057ABBFBC44BE8993D2A4701C751507FF9901AA110B5993BA070C176B",
						Status:     types.WithdrawStatusUnbond,
					},
					{
						ChainId:   zone.ChainId,
						Delegator: zone.DelegationAddress.Address,
						Recipient: "cosmos1nvkpj5n5mhy2ntvgn2cklntwx9ujvfvcacz5et",
						Distribution: []*types.Distribution{
							{Amount: 13000, Valoper: validators[0]},
							{Amount: 13300, Valoper: validators[1]},
							{Amount: 10000, Valoper: validators[2]},
							{Amount: 11000, Valoper: validators[3]},
						},
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(47300))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(36503)),
						Txhash:     "FB087A50A4836CBDFACA70D393AF110C28935276267B7BA2838BE3CEEA08F762",
						Status:     types.WithdrawStatusUnbond,
					},
				}
				return out
			},
			Delta:       sdk.NewDecWithPrec(100, 2).Quo(sdk.NewDecWithPrec(95, 2)),
			Validator:   func(validators []string) string { return validators[1] },
			ExpectError: false,
		},
		{
			Name: "overflow test",
			InitialRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				withdrawalRecord := types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: zone.DelegationAddress.Address,
					Recipient: "cosmos1v4gek4mld0k5yhpe0fsln4takg558cdpyexv2rxr3dh45f2fqrgsw52m97",
					Distribution: []*types.Distribution{
						{Amount: 9223372036854775807 + 1, Valoper: validators[1]}, // max int64 +1 - check for overflow
					},
					Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(40000))),
					BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(32356)),
					Txhash:     "3BE21C1057ABBFBC44BE8993D2A4701C751507FF9901AA110B5993BA070C176B",
					Status:     types.WithdrawStatusUnbond,
				}
				keeper.SetWithdrawalRecord(ctx, withdrawalRecord)
			},
			ExpectedRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) (out []types.WithdrawalRecord) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				out = []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: zone.DelegationAddress.Address,
						Recipient: "cosmos1v4gek4mld0k5yhpe0fsln4takg558cdpyexv2rxr3dh45f2fqrgsw52m97",
						Distribution: []*types.Distribution{
							{Amount: 9223372036854775807 + 1, Valoper: validators[1]},
						},
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(40000))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(32356)),
						Txhash:     "3BE21C1057ABBFBC44BE8993D2A4701C751507FF9901AA110B5993BA070C176B",
						Status:     types.WithdrawStatusUnbond,
					},
				}
				return out
			},
			Delta:       sdk.NewDecWithPrec(100, 2).Quo(sdk.NewDecWithPrec(95, 2)),
			Validator:   func(validators []string) string { return validators[1] },
			ExpectError: true,
		},

		{
			Name: "mismatch test",
			InitialRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				withdrawalRecord := types.WithdrawalRecord{
					ChainId:   zone.ChainId,
					Delegator: zone.DelegationAddress.Address,
					Recipient: "cosmos1nna7k5lywn99cd63elcfqm6p8c5c4qcug4aef5",
					Distribution: []*types.Distribution{
						{Amount: 100000000, Valoper: validators[1]}, // slashed amount exceeds total amount
					},
					Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10))),
					BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(10)),
					Txhash:     "8A698A142447087E5DC01F7BC3886EC1A6606D377D1FAC766FB279AD09F1407C",
					Status:     types.WithdrawStatusUnbond,
				}
				keeper.SetWithdrawalRecord(ctx, withdrawalRecord)
			},
			ExpectedRecords: func(ctx sdk.Context, keeper *icskeeper.Keeper) (out []types.WithdrawalRecord) {
				zone, _ := keeper.GetZone(ctx, suite.chainB.ChainID)
				validators := keeper.GetValidatorAddresses(ctx, zone.ChainId)
				out = []types.WithdrawalRecord{
					{
						ChainId:   zone.ChainId,
						Delegator: zone.DelegationAddress.Address,
						Recipient: "cosmos1nna7k5lywn99cd63elcfqm6p8c5c4qcug4aef5",
						Distribution: []*types.Distribution{
							{Amount: 100000000, Valoper: validators[1]}, // slashed amount exceeds total amount
						},
						Amount:     sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(10))),
						BurnAmount: sdk.NewCoin(zone.LocalDenom, math.NewInt(10)),
						Txhash:     "8A698A142447087E5DC01F7BC3886EC1A6606D377D1FAC766FB279AD09F1407C",
						Status:     types.WithdrawStatusUnbond,
					},
				}
				return out
			},
			Delta:       sdk.NewDecWithPrec(100, 2).Quo(sdk.NewDecWithPrec(95, 2)),
			Validator:   func(validators []string) string { return validators[1] },
			ExpectError: true,
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.Name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			app := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			zone, found := app.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			tc.InitialRecords(ctx, app.InterchainstakingKeeper)

			err := app.InterchainstakingKeeper.UpdateWithdrawalRecordsForSlash(ctx, &zone, tc.Validator(app.InterchainstakingKeeper.GetValidatorAddresses(ctx, zone.ChainId)), tc.Delta)
			if tc.ExpectError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
			ctx = suite.chainA.GetContext()
			for _, expected := range tc.ExpectedRecords(ctx, app.InterchainstakingKeeper) {
				wrd, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, expected.Txhash, types.WithdrawStatusUnbond)
				suite.True(found)
				suite.EqualValues(expected, wrd)
			}
		})
	}
}
