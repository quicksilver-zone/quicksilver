package types

import (
	"math"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ingenuity-build/quicksilver/utils"
)

// CalculateDeltas determines, for the current delegations, in delta between actual allocations and the target intent.
// Positive delta represents current allocation is below target, and vice versa.
func CalculateDeltas(currentAllocations map[string]sdkmath.Int, currentSum sdkmath.Int, targetAllocations ValidatorIntents) ValidatorIntents {
	deltas := make(ValidatorIntents, 0)

	targetValopers := func(in ValidatorIntents) []string {
		out := make([]string, 0, len(in))
		for _, i := range in {
			out = append(out, i.ValoperAddress)
		}
		return out
	}(targetAllocations)

	keySet := utils.Unique(append(targetValopers, utils.Keys(currentAllocations)...))
	sort.Strings(keySet)
	// for target allocations, raise the intent weight by the total delegated value to get target amount
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
		// positive == below target, negative == above target
		delta := targetAmount.Sub(current)
		deltas = append(deltas, &ValidatorIntent{Weight: sdk.NewDecFromInt(delta), ValoperAddress: valoper})
	}

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].ValoperAddress > deltas[j].ValoperAddress
	})

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].Weight.GT(deltas[j].Weight)
	})

	return deltas
}

// CalculateDeltas determines, for the current delegations, in delta between actual allocations and the target intent.
// Positive delta represents current allocation is below target, and vice versa.
func CalculateDeltasNew(currentAllocations map[string]sdkmath.Int, locked map[string]bool, currentSum sdkmath.Int, targetAllocations ValidatorIntents) (targets, sources Deltas) {
	targets = make(Deltas, 0)
	sources = make(Deltas, 0)

	targetValopers := func(in ValidatorIntents) []string {
		out := make([]string, 0, len(in))
		for _, i := range in {
			out = append(out, i.ValoperAddress)
		}
		return out
	}(targetAllocations)

	keySet := utils.Unique(append(targetValopers, utils.Keys(currentAllocations)...))
	sort.Strings(keySet)
	// for target allocations, raise the intent weight by the total delegated value to get target amount
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
		// positive == below target, negative == above target
		delta := targetAmount.Sub(current)

		if delta.IsPositive() {
			targets = append(targets, &Delta{Amount: delta, ValoperAddress: valoper})
		} else {
			if _, found := locked[valoper]; !found {
				// only append to sources if the delegation is not locked - i.e. it doesn't have an incoming redelegation.
				sources = append(sources, &Delta{Amount: delta.Abs(), ValoperAddress: valoper})
			}
		}
	}

	targets.Sort()
	sources.Sort()

	return
}

type Delta struct {
	ValoperAddress string
	Amount         sdkmath.Int
}

type Deltas []*Delta

func (d Deltas) Sort() {

	// filter zeros
	new := make(Deltas, 0)
	for _, delta := range d {
		if !delta.Amount.IsZero() {
			new = append(new, delta)
		}
	}
	d = new

	// sort keys by relative value of delta
	sort.SliceStable(d, func(i, j int) bool {
		// < sorts alphabetically.
		return d[i].ValoperAddress < d[j].ValoperAddress
	})

	// sort keys by relative value of delta
	sort.SliceStable(d, func(i, j int) bool {
		return d[i].Amount.GT(d[j].Amount)
	})
}

type RebalanceTarget struct {
	Amount sdkmath.Int
	Source string
	Target string
}

type RebalanceTargets []*RebalanceTarget

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

// DetermineAllocationsForRebalancing takes
func DetermineAllocationsForRebalancing(
	currentAllocations map[string]sdkmath.Int,
	currentLocked map[string]bool,
	currentSum sdkmath.Int,
	lockedSum sdkmath.Int,
	targetAllocations ValidatorIntents,
	logger log.Logger,
) RebalanceTargets {
	out := make(RebalanceTargets, 0)
	targets, sources := CalculateDeltasNew(currentAllocations, currentLocked, currentSum, targetAllocations)

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

// func (d Deltas) Render() (out string) {
// 	for _, delta := range d {
// 		out = out + fmt.Sprintf("%s:\t%d\n", delta.ValoperAddress, delta.Amount.Int64())
// 	}
// 	return
// }

// MinDeltas returns the lowest value in a slice of Deltas.
func MinDeltas(deltas ValidatorIntents) sdkmath.Int {
	minValue := sdk.NewInt(math.MaxInt64)
	for _, intent := range deltas {
		if minValue.GT(intent.Weight.TruncateInt()) {
			minValue = intent.Weight.TruncateInt()
		}
	}

	return minValue
}

// MaxDeltas returns the greatest value in a slice of Deltas.
func MaxDeltas(deltas ValidatorIntents) sdkmath.Int {
	maxValue := sdk.NewInt(math.MinInt64)
	for _, intent := range deltas {
		if maxValue.LT(intent.Weight.TruncateInt()) {
			maxValue = intent.Weight.TruncateInt()
		}
	}

	return maxValue
}
