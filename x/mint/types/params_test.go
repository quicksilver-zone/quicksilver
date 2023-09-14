package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/mint/types"
)

func TestParamsValidate(t *testing.T) {
	tests := []struct {
		name    string
		params  types.Params
		isValid bool
	}{
		{
			name:    "valid genesis",
			params:  types.DefaultParams(),
			isValid: true,
		},
		{
			name: "invalid mint denom",
			params: types.Params{
				MintDenom:               "", // empty string
				GenesisEpochProvisions:  sdk.NewDec(200000000 / 122),
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1), // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1), // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1), // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1), // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid genesis epoch provisions",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1), // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1), // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1), // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1), // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid epoch reduction period",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: -1,                        // negative
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1), // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1), // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1), // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1), // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid reduction factor 1",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),
				EpochIdentifier:         "day",                      // 1 day
				ReductionPeriodInEpochs: 365,                        // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(-75, 2), // negative
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1), // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1), // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1), // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1), // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid reduction factor 2",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 1), // greater than 1
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1), // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1), // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1), // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1), // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid distribution proportions 1",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(-3, 1), // -0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1),  // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1),  // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1),  // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid distribution proportions 2",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1),  // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(-3, 1), // -0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1),  // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1),  // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid distribution proportions 3",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1),  // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1),  // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(-3, 1), // -0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1),  // 0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid distribution proportions 4",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1),  // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1),  // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1),  // 0.3
					CommunityPool:        sdk.NewDecWithPrec(-1, 1), // -0.1
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid distribution proportions 5",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1), // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1), // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1), // 0.3
					CommunityPool:        sdk.NewDecWithPrec(3, 1), // 0.3
				},
				MintingRewardsDistributionStartEpoch: 0,
			},
			isValid: false,
		},
		{
			name: "invalid distribution start epoch",
			params: types.Params{
				MintDenom:               sdk.DefaultBondDenom,
				GenesisEpochProvisions:  sdk.NewDec(-1),            // negative
				EpochIdentifier:         "day",                     // 1 day
				ReductionPeriodInEpochs: 365,                       // 1 year
				ReductionFactor:         sdk.NewDecWithPrec(75, 2), // 0.75
				DistributionProportions: types.DistributionProportions{
					Staking:              sdk.NewDecWithPrec(3, 1), // 0.3
					PoolIncentives:       sdk.NewDecWithPrec(3, 1), // 0.3
					ParticipationRewards: sdk.NewDecWithPrec(3, 1), // 0.3
					CommunityPool:        sdk.NewDecWithPrec(1, 1), // 0.3
				},
				MintingRewardsDistributionStartEpoch: -1, // negative
			},
			isValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if !tc.isValid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// TestGetDeveloperVestingProportion sanity checks that participation
// rewards proportion equals to the value set by
// parameter for participation rewards.
func TestGetDistributionProportions(t *testing.T) {
	expected := sdk.NewDecWithPrec(4, 1)

	params := types.Params{
		DistributionProportions: types.DistributionProportions{
			ParticipationRewards: expected,
		},
	}

	actual := params.GetDistributionProportions().ParticipationRewards
	require.Equal(t, expected, actual)
}
