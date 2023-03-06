package simulation

import (
	"encoding/json"
	"fmt"
	appconfig "github.com/ingenuity-build/quicksilver/cmd/config"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ingenuity-build/quicksilver/x/mint/types"
)

// Simulation parameter constants.
const (
	epochProvisionsKey         = "genesis_epoch_provisions"
	reductionFactorKey         = "reduction_factor"
	reductionPeriodInEpochsKey = "reduction_period_in_epochs"
	distributionProportionsKey = "distribution_proportions"

	mintingRewardsDistributionStartEpochKey = "minting_rewards_distribution_start_epoch"

	epochIdentifier = "day"
	maxInt64        = int(^uint(0) >> 1)
)

// genDistributionProportions generates randomized DistributionProportions.
func genDistributionProportions(r *rand.Rand) types.DistributionProportions {
	staking := r.Int63n(99)
	left := int64(100) - staking
	poolIncentives := r.Int63n(left)
	left -= poolIncentives
	participationRewards := r.Int63n(left)
	communityPool := left - participationRewards

	return types.DistributionProportions{
		Staking:              sdk.NewDecWithPrec(staking, 2),
		PoolIncentives:       sdk.NewDecWithPrec(poolIncentives, 2),
		ParticipationRewards: sdk.NewDecWithPrec(participationRewards, 2),
		CommunityPool:        sdk.NewDecWithPrec(communityPool, 2),
	}
}

func genEpochProvisions(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(int64(r.Intn(maxInt64)))
}

func genReductionFactor(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(10)), 1)
}

func genReductionPeriodInEpochs(r *rand.Rand) int64 {
	return int64(r.Intn(maxInt64))
}

func genMintintRewardsDistributionStartEpoch(r *rand.Rand) int64 {
	return int64(r.Intn(maxInt64))
}

// RandomizedGenState generates a random GenesisState for mint.
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var epochProvisions sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, epochProvisionsKey, &epochProvisions, simState.Rand,
		func(r *rand.Rand) {
			epochProvisions = genEpochProvisions(r)
		},
	)

	var distributionProportions types.DistributionProportions
	simState.AppParams.GetOrGenerate(
		simState.Cdc, distributionProportionsKey, &distributionProportions, simState.Rand,
		func(r *rand.Rand) {
			distributionProportions = genDistributionProportions(r)
		},
	)

	var reductionFactor sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, reductionFactorKey, &reductionFactor, simState.Rand,
		func(r *rand.Rand) {
			reductionFactor = genReductionFactor(r)
		},
	)

	var reductionPeriodInEpochs int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, reductionPeriodInEpochsKey, &reductionPeriodInEpochs, simState.Rand,
		func(r *rand.Rand) { reductionPeriodInEpochs = genReductionPeriodInEpochs(r) },
	)

	var mintintRewardsDistributionStartEpoch int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintingRewardsDistributionStartEpochKey, &mintintRewardsDistributionStartEpoch, simState.Rand,
		func(r *rand.Rand) { mintintRewardsDistributionStartEpoch = genMintintRewardsDistributionStartEpoch(r) },
	)

	mintDenom := appconfig.DefaultBondDenom

	minter := types.Minter{EpochProvisions: epochProvisions}
	params := types.Params{
		MintDenom:                            mintDenom,
		GenesisEpochProvisions:               epochProvisions,
		EpochIdentifier:                      epochIdentifier,
		ReductionPeriodInEpochs:              reductionPeriodInEpochs,
		ReductionFactor:                      reductionFactor,
		DistributionProportions:              distributionProportions,
		MintingRewardsDistributionStartEpoch: mintintRewardsDistributionStartEpoch,
	}
	err := params.Validate()
	if err != nil {
		panic(err)
	}

	mintGenesis := types.GenesisState{
		Minter: minter,
		Params: params,
	}

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&mintGenesis)
}
