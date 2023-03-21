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

type RebalanceTarget struct {
	Amount sdkmath.Int
	Source string
	Target string
}

func DetermineAllocationsForRebalancing(
	currentAllocations map[string]sdkmath.Int,
	currentLocked map[string]bool,
	currentSum sdkmath.Int,
	targetAllocations ValidatorIntents,
	existingRedelegations []RedelegationRecord,
	log log.Logger,
) []RebalanceTarget {
	out := make([]RebalanceTarget, 0)
	deltas := CalculateDeltas(currentAllocations, currentSum, targetAllocations)

	wantToRebalance := sdk.ZeroInt()
	canRebalanceFrom := sdk.ZeroInt()

	totalLocked := int64(0)
	lockedPerValidator := map[string]int64{}
	for _, redelegation := range existingRedelegations {
		totalLocked += redelegation.Amount
		thisLocked, found := lockedPerValidator[redelegation.Destination]
		if !found {
			thisLocked = 0
		}
		lockedPerValidator[redelegation.Destination] = thisLocked + redelegation.Amount
	}
	for _, valoper := range utils.Keys(currentAllocations) {
		// if validator already has a redelegation _to_ it, we can no longer redelegate _from_ it (transitive redelegations)
		// remove _locked_ amount from lpv and total locked for purposes of rebalancing.
		if currentLocked[valoper] {
			thisLocked, found := lockedPerValidator[valoper]
			if !found {
				thisLocked = 0
			}
			totalLocked = totalLocked - thisLocked + currentAllocations[valoper].Int64()
			lockedPerValidator[valoper] = currentAllocations[valoper].Int64()
		}
	}

	// TODO: make these params
	maxCanRebalanceTotal := currentSum.Sub(sdkmath.NewInt(totalLocked)).Quo(sdk.NewInt(2))
	maxCanRebalance := sdkmath.MinInt(maxCanRebalanceTotal, currentSum.Quo(sdk.NewInt(7)))
	if log != nil {
		log.Debug("Rebalancing", "totalLocked", totalLocked, "lockedPerValidator", lockedPerValidator, "canRebalanceTotal", maxCanRebalanceTotal, "canRebalanceEpoch", maxCanRebalance)
	}

	// deltas are sorted in CalculateDeltas; don't re-sort.
	for _, delta := range deltas {
		switch {
		case delta.Weight.IsZero():
			// do nothing
		case delta.Weight.IsPositive():
			// if delta > current value - locked value, truncate, as we cannot rebalance locked tokens.
			wantToRebalance = wantToRebalance.Add(delta.Weight.TruncateInt())
		case delta.Weight.IsNegative():
			if delta.Weight.Abs().GT(sdk.NewDecFromInt(currentAllocations[delta.ValoperAddress].Sub(sdkmath.NewInt(lockedPerValidator[delta.ValoperAddress])))) {
				delta.Weight = sdk.NewDecFromInt(currentAllocations[delta.ValoperAddress].Sub(sdkmath.NewInt(lockedPerValidator[delta.ValoperAddress]))).Neg()
				if log != nil {
					log.Debug("Truncated delta due to locked tokens", "valoper", delta.ValoperAddress, "delta", delta.Weight.Abs())
				}
			}
			canRebalanceFrom = canRebalanceFrom.Add(delta.Weight.Abs().TruncateInt())
		}
	}

	toRebalance := sdk.MinInt(sdk.MinInt(wantToRebalance, canRebalanceFrom), maxCanRebalance)

	if toRebalance.Equal(sdkmath.ZeroInt()) {
		if log != nil {
			log.Debug("No rebalancing this epoch")
		}
		return []RebalanceTarget{}
	}
	if log != nil {
		log.Debug("Will rebalance this epoch", "amount", toRebalance)
	}

	tgtIdx := 0
	srcIdx := len(deltas) - 1
	for i := 0; toRebalance.GT(sdk.ZeroInt()); {
		i++
		if i > 20 {
			break
		}
		src := deltas[srcIdx]
		tgt := deltas[tgtIdx]
		if src.ValoperAddress == tgt.ValoperAddress {
			break
		}
		var amount sdkmath.Int
		if src.Weight.Abs().TruncateInt().IsZero() { //nolint:gocritic
			srcIdx--
			continue
		} else if src.Weight.Abs().TruncateInt().GT(toRebalance) { // amount == rebalance
			amount = toRebalance
		} else {
			amount = src.Weight.Abs().TruncateInt()
		}

		if tgt.Weight.Abs().TruncateInt().IsZero() {
			tgtIdx++
			continue
		} else if tgt.Weight.Abs().TruncateInt().LTE(toRebalance) {
			amount = sdk.MinInt(amount, tgt.Weight.Abs().TruncateInt())
		}

		out = append(out, RebalanceTarget{Amount: amount, Target: tgt.ValoperAddress, Source: src.ValoperAddress})
		deltas[srcIdx].Weight = src.Weight.Add(sdk.NewDecFromInt(amount))
		deltas[tgtIdx].Weight = tgt.Weight.Sub(sdk.NewDecFromInt(amount))
		toRebalance = toRebalance.Sub(amount)

	}

	// sort keys by relative value of delta
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Source < out[j].Source
	})

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Target < out[j].Target
	})

	// sort keys by relative value of delta
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Amount.GT(out[j].Amount)
	})

	return out
}

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
