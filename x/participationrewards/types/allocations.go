package types

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RewardsAllocation struct {
	ValidatorSelection math.Int
	Holdings           math.Int
	Lockup             math.Int
}

// GetRewardsAllocations returns an instance of rewardsAllocation with values
// set according to the given moduleBalance and distribution proportions.
func GetRewardsAllocations(moduleBalance math.Int, proportions DistributionProportions) (*RewardsAllocation, error) {
	if moduleBalance.IsNil() || moduleBalance.IsZero() {
		return nil, ErrNothingToAllocate
	}

	if sum := proportions.Total(); !sum.Equal(sdk.OneDec()) {
		return nil, fmt.Errorf("%w: got %v", ErrInvalidTotalProportions, sum)
	}

	var allocation RewardsAllocation

	// split participation rewards allocations
	allocation.ValidatorSelection = sdk.NewDecFromInt(moduleBalance).Mul(proportions.ValidatorSelectionAllocation).TruncateInt()
	allocation.Holdings = sdk.NewDecFromInt(moduleBalance).Mul(proportions.HoldingsAllocation).TruncateInt()
	allocation.Lockup = sdk.NewDecFromInt(moduleBalance).Mul(proportions.LockupAllocation).TruncateInt()

	// use sum to check total distribution to collect and allocate dust
	sum := allocation.Lockup.Add(allocation.ValidatorSelection).Add(allocation.Holdings)
	dust := moduleBalance.Sub(sum)

	// Add dust to validator choice allocation (favors decentralization)
	allocation.ValidatorSelection = allocation.ValidatorSelection.Add(dust)

	return &allocation, nil
}
