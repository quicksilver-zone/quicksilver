package types

import (
	"fmt"
	"math"
	"sort"

	"github.com/tendermint/tendermint/libs/log"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils"
)

// CalculateAllocationDeltas determines, for the current delegations, in delta between actual allocations and the target intent.
// Returns a slice of deltas for each of target allocations (underallocated) and source allocations (overallocated).
func CalculateAllocationDeltas(
	currentAllocations map[string]sdkmath.Int,
	locked map[string]bool,
	currentSum sdkmath.Int,
	targetAllocations ValidatorIntents,
	maxCanAllocate map[string]sdkmath.Int,
) (targets, sources AllocationDeltas) {
	targets = make(AllocationDeltas, 0)
	sources = make(AllocationDeltas, 0)

	// reduce ValidatorIntents to slice of Valoper addresses.
	targetValopers := func(in ValidatorIntents) []string {
		out := make([]string, 0, len(in))
		for _, i := range in {
			out = append(out, i.ValoperAddress)
		}
		return out
	}(targetAllocations)

	// create a slide of unique valopers across current and target allocations.
	keySet := utils.Unique(append(targetValopers, utils.Keys(currentAllocations)...))
	sort.Strings(keySet)

	// for target allocations, raise the intent weight by the total delegated value to get target amount.
	for _, valoper := range keySet {
		current, ok := currentAllocations[valoper]
		if !ok {
			current = sdk.ZeroInt()
		}

		target, ok := targetAllocations.GetForValoper(valoper)
		if !ok {
			target = &ValidatorIntent{ValoperAddress: valoper, Weight: sdk.ZeroDec()}
		}
		targetAmount := target.Weight.MulInt(currentSum).TruncateInt()

		// diff between target and current allocations
		// positive == below target (target), negative == above target (source)
		delta := targetAmount.Sub(current)
		max, ok := maxCanAllocate[valoper]
		if !ok {
			max = delta
		}
		if max.LT(delta) {
			delta = max
		}

		if delta.IsPositive() {
			targets = append(targets, &AllocationDelta{Amount: delta, ValoperAddress: valoper})
		} else {
			if _, found := locked[valoper]; !found {
				// only append to sources if the delegation is not locked - i.e. it doesn't have an incoming redelegation.
				// TODO: this needs to be locked amounts for unbonding purposes. Redelegations do not care about amounts, but unbondings do.
				sources = append(sources, &AllocationDelta{Amount: delta.Abs(), ValoperAddress: valoper})
			}
		}
	}

	// sort for determinism.
	targets.Sort()
	sources.Sort()

	return targets, sources
}

type AllocationDelta struct {
	ValoperAddress string
	Amount         sdkmath.Int
}

type AllocationDeltas []*AllocationDelta

func (deltas AllocationDeltas) Sort() {
	// filter zeros
	newAllocationDeltas := make(AllocationDeltas, 0)
	for _, delta := range deltas {
		if !delta.Amount.IsZero() {
			newAllocationDeltas = append(newAllocationDeltas, delta)
		}
	}
	deltas = newAllocationDeltas

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		// < sorts alphabetically.
		return deltas[i].ValoperAddress < deltas[j].ValoperAddress
	})

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].Amount.GT(deltas[j].Amount)
	})
}

func (deltas AllocationDeltas) Sum() (sum sdkmath.Int) {
	sum = sdkmath.ZeroInt()
	for _, delta := range deltas {
		sum = sum.Add(delta.Amount)
	}
	return sum
}

func (deltas AllocationDeltas) GetForValoper(valoper string) (out *AllocationDelta, found bool) {
	for _, delta := range deltas {
		if delta.ValoperAddress == valoper {
			return delta, true
		}
	}
	return out, false
}

// Render AllocationDeltas as a string.
func (deltas AllocationDeltas) String() (out string) {
	for _, delta := range deltas {
		out = fmt.Sprintf("%s%s:\t%d\n", out, delta.ValoperAddress, delta.Amount.Int64())
	}
	return out
}

// MinDelta returns the lowest value in a slice of AllocationDeltas.
func (deltas AllocationDeltas) MinDelta() sdkmath.Int {
	minValue := sdk.NewInt(math.MaxInt64)
	for _, delta := range deltas {
		if minValue.GT(delta.Amount) {
			minValue = delta.Amount
		}
	}

	return minValue
}

// MaxDelta returns the greatest value in a slice of AllocationDeltas.
func (deltas AllocationDeltas) MaxDelta() sdkmath.Int {
	maxValue := sdk.NewInt(math.MinInt64)
	for _, delta := range deltas {
		if maxValue.LT(delta.Amount) {
			maxValue = delta.Amount
		}
	}

	return maxValue
}

// Negate the values of all the AllocationDeltas.
func (deltas *AllocationDeltas) Negate() {
	for _, delta := range *deltas {
		delta.Amount = delta.Amount.Neg()
	}
}

type RebalanceTarget struct {
	Amount sdkmath.Int
	Source string
	Target string
}

type RebalanceTargets []*RebalanceTarget

// Sort RebalanceTargets deterministically.
func (t RebalanceTargets) Sort() {
	// sort keys by relative value of delta
	sort.SliceStable(t, func(i, j int) bool {
		// < sorts alphabetically.
		return t[i].Source < t[j].Source
	})

	// sort keys by relative value of delta
	sort.SliceStable(t, func(i, j int) bool {
		// < sorts alphabetically.
		return t[i].Target < t[j].Target
	})

	// sort keys by relative value of delta
	sort.SliceStable(t, func(i, j int) bool {
		return t[i].Amount.LT(t[j].Amount)
	})
}

func (rb RebalanceTargets) RemoveDuplicates() RebalanceTargets {
	encountered := make(map[string]bool)
	result := make(RebalanceTargets, 0)

	for _, r := range rb {
		if r.Amount.IsZero() {
			continue
		}
		key := fmt.Sprintf("%v-%s-%s", r.Amount.String(), r.Source, r.Target)
		if !encountered[key] {
			encountered[key] = true
			result = append(result, r)
		}
	}
	return result
}

// DetermineAllocationsForRebalancing takes maps of current and locked delegations, and based upon the target allocations,
// attempts to satisfy the target allocations in the fewest number of transformations. It returns a slice of RebalanceTargets.
func DetermineAllocationsForRebalancing(
	currentAllocations map[string]sdkmath.Int,
	currentLocked map[string]bool,
	currentSum sdkmath.Int,
	lockedSum sdkmath.Int,
	targetAllocations ValidatorIntents,
	maxCanAllocate map[string]sdkmath.Int,
	logger log.Logger,
) RebalanceTargets {
	out := make(RebalanceTargets, 0)
	targets, sources := CalculateAllocationDeltas(currentAllocations, currentLocked, currentSum, targetAllocations, maxCanAllocate)

	// rebalanceBudget = (total_delegations - locked)/2 == 50% of (total_delegations - locked)
	// TODO: make this 2 (max_redelegation_factor) a param.
	rebalanceBudget := currentSum.Sub(lockedSum).Quo(sdk.NewInt(2))

	if logger != nil {
		logger.Debug("Rebalancing", "total", currentSum, "totalLocked", lockedSum, "rebalanceBudget", rebalanceBudget)
	}

TARGET:
	// targets are validators with a delegation deficit, sorted in descending order.
	// that is, those at the top should be satisfied first to maximise progress toward goal.
	for _, target := range targets {
		// amount is amount we should try to redelegate toward target. This may be constrained by the remaining redelegateBudget.
		// if it is zero (i.e. we hit the redelegation budget) break out of the loop.
		amount := sdkmath.MinInt(target.Amount, rebalanceBudget)
		if amount.IsZero() {
			break
		}
		sources.Sort()
		// sources are validators with available balance to redelegate, sorted in desc order.
		for _, source := range sources {
			switch {
			case source.Amount.IsZero():
				// if source is zero, skip.
				continue
			case source.Amount.GTE(amount):
				// if source >= amount, fully satisfy target.
				out = append(out, &RebalanceTarget{Amount: amount, Target: target.ValoperAddress, Source: source.ValoperAddress})
				source.Amount = source.Amount.Sub(amount)
				target.Amount = target.Amount.Sub(amount)
				rebalanceBudget = rebalanceBudget.Sub(amount)
				continue TARGET
			case source.Amount.LT(amount):
				// if source < amount, partially satisfy amount.
				out = append(out, &RebalanceTarget{Amount: source.Amount, Target: target.ValoperAddress, Source: source.ValoperAddress})
				amount = amount.Sub(source.Amount)
				target.Amount = target.Amount.Sub(source.Amount)
				rebalanceBudget = rebalanceBudget.Sub(source.Amount)
				source.Amount = source.Amount.Sub(source.Amount)
				if amount.IsZero() || rebalanceBudget.IsZero() {
					// if the amount is fully satisfied or the rebalanceBudget is zero, skip to next target.
					continue TARGET
				}
				// otherwise, try next source.
			}
		}
		// we only get here if we are unable to satisfy targets due to rebalanceBudget depletion.
		if logger != nil {
			logger.Info("unable to satisfy targets with available sources.")
		}
	}

	out.Sort()

	return out
}
