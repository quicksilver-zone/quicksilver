package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/airdrop/keeper"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

func (suite *KeeperTestSuite) TestHandleRegisterZoneDropProposal() {
	appA := suite.GetQuicksilverApp(suite.chainA)

	validZoneDrop := types.ZoneDrop{
		ChainId:    suite.chainB.ChainID,
		StartTime:  time.Now().Add(time.Hour),
		Duration:   time.Hour,
		Decay:      30 * time.Minute,
		Allocation: 1000000000,
		Actions: []sdk.Dec{
			0:  sdk.MustNewDecFromStr("0.15"), // 15%
			1:  sdk.MustNewDecFromStr("0.06"), // 21%
			2:  sdk.MustNewDecFromStr("0.07"), // 28%
			3:  sdk.MustNewDecFromStr("0.08"), // 36%
			4:  sdk.MustNewDecFromStr("0.09"), // 45%
			5:  sdk.MustNewDecFromStr("0.1"),  // 55%
			6:  sdk.MustNewDecFromStr("0.15"), // 70%
			7:  sdk.MustNewDecFromStr("0.05"), // 75%
			8:  sdk.MustNewDecFromStr("0.1"),  // 85%
			9:  sdk.MustNewDecFromStr("0.1"),  // 95%
			10: sdk.MustNewDecFromStr("0.05"), // 100%
		},
		IsConcluded: false,
	}
	userAddresses := []string{
		utils.GenerateAccAddressForTest().String(),
	}

	prop := types.RegisterZoneDropProposal{}
	tests := []struct {
		name     string
		malleate func()
		wantErr  bool
	}{
		{
			"blank",
			func() {},
			true,
		},
		{
			"invalid-zd-chainID",
			func() {
				zd := types.ZoneDrop{
					ChainId:    "test-01",
					StartTime:  time.Now().Add(-5 * time.Minute),
					Duration:   time.Hour,
					Decay:      30 * time.Minute,
					Allocation: 1000000000,
					Actions: []sdk.Dec{
						0:  sdk.MustNewDecFromStr("0.15"), // 15%
						1:  sdk.MustNewDecFromStr("0.06"), // 21%
						2:  sdk.MustNewDecFromStr("0.07"), // 28%
						3:  sdk.MustNewDecFromStr("0.08"), // 36%
						4:  sdk.MustNewDecFromStr("0.09"), // 45%
						5:  sdk.MustNewDecFromStr("0.1"),  // 55%
						6:  sdk.MustNewDecFromStr("0.15"), // 70%
						7:  sdk.MustNewDecFromStr("0.05"), // 75%
						8:  sdk.MustNewDecFromStr("0.1"),  // 85%
						9:  sdk.MustNewDecFromStr("0.1"),  // 95%
						10: sdk.MustNewDecFromStr("0.05"), // 100%
					},
					IsConcluded: false,
				}

				crs := make([]types.ClaimRecord, len(userAddresses))
				for i := range crs {
					crs[i] = types.ClaimRecord{
						ChainId:          suite.chainB.ChainID,
						Address:          userAddresses[i],
						ActionsCompleted: nil,
						MaxAllocation:    100000000,
						BaseValue:        10000000,
					}
				}

				prop = types.RegisterZoneDropProposal{
					Title:        "Test Zone Airdrop Proposal",
					Description:  "Adding this zone drop allows for automated testing",
					ZoneDrop:     &zd,
					ClaimRecords: suite.compressClaimRecords(crs),
				}
			},
			true,
		},
		{
			"invalid-zd-started",
			func() {
				zd := types.ZoneDrop{
					ChainId:    suite.chainB.ChainID,
					StartTime:  time.Now().Add(-5 * time.Minute),
					Duration:   time.Hour,
					Decay:      30 * time.Minute,
					Allocation: 1000000000,
					Actions: []sdk.Dec{
						0:  sdk.MustNewDecFromStr("0.15"), // 15%
						1:  sdk.MustNewDecFromStr("0.06"), // 21%
						2:  sdk.MustNewDecFromStr("0.07"), // 28%
						3:  sdk.MustNewDecFromStr("0.08"), // 36%
						4:  sdk.MustNewDecFromStr("0.09"), // 45%
						5:  sdk.MustNewDecFromStr("0.1"),  // 55%
						6:  sdk.MustNewDecFromStr("0.15"), // 70%
						7:  sdk.MustNewDecFromStr("0.05"), // 75%
						8:  sdk.MustNewDecFromStr("0.1"),  // 85%
						9:  sdk.MustNewDecFromStr("0.1"),  // 95%
						10: sdk.MustNewDecFromStr("0.05"), // 100%
					},
					IsConcluded: false,
				}

				crs := make([]types.ClaimRecord, len(userAddresses))
				for i := range crs {
					crs[i] = types.ClaimRecord{
						ChainId:          suite.chainB.ChainID,
						Address:          userAddresses[i],
						ActionsCompleted: nil,
						MaxAllocation:    100000000,
						BaseValue:        10000000,
					}
				}

				prop = types.RegisterZoneDropProposal{
					Title:        "Test Zone Airdrop Proposal",
					Description:  "Adding this zone drop allows for automated testing",
					ZoneDrop:     &zd,
					ClaimRecords: suite.compressClaimRecords(crs),
				}
			},
			true,
		},
		{
			"invalid-cr-chainID",
			func() {
				zd := validZoneDrop
				crs := make([]types.ClaimRecord, len(userAddresses))
				for i := range crs {
					crs[i] = types.ClaimRecord{
						ChainId:          "test-01",
						Address:          userAddresses[i],
						ActionsCompleted: nil,
						MaxAllocation:    100000000,
						BaseValue:        10000000,
					}
				}

				prop = types.RegisterZoneDropProposal{
					Title:        "Test Zone Airdrop Proposal",
					Description:  "Adding this zone drop allows for automated testing",
					ZoneDrop:     &zd,
					ClaimRecords: suite.compressClaimRecords(crs),
				}
			},
			true,
		},
		{
			"invalid-cr-completed-actions",
			func() {
				zd := validZoneDrop

				crs := make([]types.ClaimRecord, len(userAddresses))
				for i := range crs {
					crs[i] = types.ClaimRecord{
						ChainId: suite.chainB.ChainID,
						Address: userAddresses[i],
						ActionsCompleted: map[int32]*types.CompletedAction{
							1: {},
						},
						MaxAllocation: 100000000,
						BaseValue:     10000000,
					}
				}

				prop = types.RegisterZoneDropProposal{
					Title:        "Test Zone Airdrop Proposal",
					Description:  "Adding this zone drop allows for automated testing",
					ZoneDrop:     &zd,
					ClaimRecords: suite.compressClaimRecords(crs),
				}
			},
			true,
		},
		{
			"invalid-allocation-exceeded",
			func() {
				zd := validZoneDrop

				crs := make([]types.ClaimRecord, len(userAddresses))
				for i := range crs {
					crs[i] = types.ClaimRecord{
						ChainId:          suite.chainB.ChainID,
						Address:          userAddresses[i],
						ActionsCompleted: nil,
						MaxAllocation:    1000000001,
						BaseValue:        10000000,
					}
				}

				prop = types.RegisterZoneDropProposal{
					Title:        "Test Zone Airdrop Proposal",
					Description:  "Adding this zone drop allows for automated testing",
					ZoneDrop:     &zd,
					ClaimRecords: suite.compressClaimRecords(crs),
				}
			},
			true,
		},
		{
			"valid",
			func() {
				zd := validZoneDrop

				crs := make([]types.ClaimRecord, len(userAddresses))
				for i := range crs {
					crs[i] = types.ClaimRecord{
						ChainId:          suite.chainB.ChainID,
						Address:          userAddresses[i],
						ActionsCompleted: nil,
						MaxAllocation:    100000000,
						BaseValue:        10000000,
					}
				}

				prop = types.RegisterZoneDropProposal{
					Title:        "Test Zone Airdrop Proposal",
					Description:  "Adding this zone drop allows for automated testing",
					ZoneDrop:     &zd,
					ClaimRecords: suite.compressClaimRecords(crs),
				}
			},
			false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.malleate()

			k := appA.AirdropKeeper
			err := keeper.HandleRegisterZoneDropProposal(suite.chainA.GetContext(), k, &prop)
			if tt.wantErr {
				suite.Require().Error(err)
				suite.T().Logf("Error: %v", err)
				return
			}

			suite.Require().NoError(err)
		})
	}
}
