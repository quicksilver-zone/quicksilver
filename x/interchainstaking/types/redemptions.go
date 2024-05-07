package types

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils"
)

// remove zero and potential negative values
func filter(in map[string]math.Int) map[string]math.Int {
	out := make(map[string]math.Int)
	for _, v := range utils.Keys[math.Int](in) {
		if !in[v].IsNil() && in[v].IsPositive() {
			out[v] = in[v]
		}
	}
	return out
}

func DetermineAllocationsForUndelegation(currentAllocations map[string]math.Int, lockedAllocations map[string]bool, currentSum math.Int, targetAllocations ValidatorIntents, availablePerValidator map[string]math.Int, amount sdk.Coins) (map[string]math.Int, error) {
	outWeights := make(map[string]math.Int)
	if len(amount) != 1 {
		return outWeights, fmt.Errorf("amount was invalid, expected sdk.Coins of length 1, got length %d", len(amount))
	}

	if !amount[0].Amount.IsPositive() {
		return outWeights, fmt.Errorf("amount was invalid, expected positive value, got %s", amount[0].Amount.String())
	}
	input := amount[0].Amount
	underAllocated, overAllocated := CalculateAllocationDeltas(currentAllocations, lockedAllocations, currentSum /* .Sub(input) */, targetAllocations, make(map[string]math.Int))

	outSum := sdk.ZeroInt()

	// deltas: +ve is below target; -ve is above target.

	// q1: can we satisfy this unbonding using _just_ above target allocations.
	// example:
	// we have v1: 5000, v2: 1800; v3: 1200; v4: 1000 and targets of 50%, 20%, 15% and 5% respectively.
	// deltas == -500, 0, 150, -550
	// an unbonding of 300 should come from v1, v4 (as they has an excess of > unbond amount) _before_ touching anything else.

	sum := overAllocated.Sum()

	overAllocationSplit := sdk.MinInt(sum, input)

	// if the sum of 'overallocated' validators > 0 (else div/nil), try to use these to satisfy the unbonding first.
	if !overAllocationSplit.IsZero() {
		for idx := range overAllocated {
			// use Amount+1 in the line below to avoid 1 remaining where truncation leaves 1 remaining - e.g. 1000 => 333/333/333 + 1.
			outWeights[overAllocated[idx].ValoperAddress] = sdk.NewDecFromInt(overAllocated[idx].Amount).Quo(sdk.NewDecFromInt(sum)).Mul(sdk.NewDecFromInt(overAllocationSplit)).TruncateInt()
			if outWeights[overAllocated[idx].ValoperAddress].GT(availablePerValidator[overAllocated[idx].ValoperAddress]) {
				// use up all of overAllocated[idx] and set available to zero.
				outWeights[overAllocated[idx].ValoperAddress] = availablePerValidator[overAllocated[idx].ValoperAddress]
				availablePerValidator[overAllocated[idx].ValoperAddress] = sdk.ZeroInt()
			} else {
				// or don't, and reduce available as appropriate.
				availablePerValidator[overAllocated[idx].ValoperAddress] = availablePerValidator[overAllocated[idx].ValoperAddress].Sub(outWeights[overAllocated[idx].ValoperAddress])
			}
			overAllocated[idx].Amount = overAllocated[idx].Amount.Sub(outWeights[overAllocated[idx].ValoperAddress])
			outSum = outSum.Add(outWeights[overAllocated[idx].ValoperAddress])
		}
	}

	// if remaining amount to distribute is zero, shortcut exit here.
	input = input.Sub(outSum)
	if input.IsZero() {
		return filter(outWeights), nil
	}

	if input.IsNegative() {
		return map[string]math.Int{}, fmt.Errorf("input is unexpectedly negative (1), aborting")
	}

	// negate all values in underallocated.
	underAllocated.Negate()
	// append the two slices
	// nolint:gocritic
	deltas := append(overAllocated, underAllocated...)
	deltas.Sort()

	maxValue := deltas.MaxDelta()
	sum = sdk.ZeroInt()

	// drop all deltas such that the maximum value is zero, and invert.
	for idx := range deltas {
		deltas[idx].Amount = deltas[idx].Amount.Add(maxValue).Abs()
		// sum here instead of calling Sum() later to save looping over slice again.
		sum = sum.Add(deltas[idx].Amount.Abs())
	}

	// unequalSplit is the portion of input that should be distributed in attempt to make deltas == 0
	unequalSplit := sdk.MinInt(sum, input)

	if unequalSplit.IsPositive() {
		for idx := range deltas {
			allocation := sdk.NewDecFromInt(deltas[idx].Amount).Quo(sdk.NewDecFromInt(sum)).Mul(sdk.NewDecFromInt(unequalSplit)).TruncateInt()
			_, ok := availablePerValidator[deltas[idx].ValoperAddress]
			if !ok {
				availablePerValidator[deltas[idx].ValoperAddress] = sdk.ZeroInt()
			}

			if allocation.GT(availablePerValidator[deltas[idx].ValoperAddress]) {
				allocation = availablePerValidator[deltas[idx].ValoperAddress]
				availablePerValidator[deltas[idx].ValoperAddress] = sdk.ZeroInt()
			} else {
				availablePerValidator[deltas[idx].ValoperAddress] = availablePerValidator[deltas[idx].ValoperAddress].Sub(allocation)
			}

			deltas[idx].Amount = deltas[idx].Amount.Sub(allocation)
			value, ok := outWeights[deltas[idx].ValoperAddress]
			if !ok {
				value = sdk.ZeroInt()
			}

			outWeights[deltas[idx].ValoperAddress] = value.Add(allocation)
			outSum = outSum.Add(allocation)
			input = input.Sub(allocation)
		}
	}

	if input.IsNegative() {
		return map[string]math.Int{}, fmt.Errorf("input is unexpectedly negative (2), aborting")
	}

	// equalSplit is the portion of input that should be distributed across all validators proportion to intent, once targets are met.
	deltas.Sort()
	if outSum.LT(amount[0].Amount) {
		// remove validators with no remaining balance from intents, and split remaining amount proportionally.
		newTargetAllocations := make(ValidatorIntents, 0, len(targetAllocations))
		for idx := range targetAllocations.Sort() {
			if !availablePerValidator[targetAllocations[idx].ValoperAddress].IsZero() {
				newTargetAllocations = append(newTargetAllocations, targetAllocations[idx])
			}
		}

		weights := newTargetAllocations.Normalize().Sort()

		origin := input
		for idx := range weights {
			perValidatorAmount := sdk.NewDecFromInt(origin).Mul(weights[idx].Weight).TruncateInt()
			value, ok := outWeights[weights[idx].ValoperAddress]
			if !ok {
				value = sdk.ZeroInt()
			}
			if perValidatorAmount.GT(availablePerValidator[weights[idx].ValoperAddress]) {
				perValidatorAmount = availablePerValidator[weights[idx].ValoperAddress]
				availablePerValidator[weights[idx].ValoperAddress] = sdk.ZeroInt()
			} else {
				availablePerValidator[weights[idx].ValoperAddress] = availablePerValidator[weights[idx].ValoperAddress].Sub(perValidatorAmount)
			}
			outWeights[weights[idx].ValoperAddress] = value.Add(perValidatorAmount)
			outSum = outSum.Add(perValidatorAmount)
			input = input.Sub(perValidatorAmount)
		}
	}

	if !outSum.LTE(amount[0].Amount) {
		return map[string]math.Int{}, fmt.Errorf("outSum (%s) is unexpectedly greater than the input amount (%s), aborting", outSum.String(), amount[0].Amount.String())
	}
	if input.IsNegative() {
		return map[string]math.Int{}, fmt.Errorf("input is unexpectedly negative (2), aborting")
	}

	// dust is the portion of the input that was truncated in previous calculations; add this to the first validator in the list,
	// available balance permitting. this should be the biggest source. This will usually be a small amount, and will negated by
	// the delta calculations on the next run.
	dust := amount[0].Amount.Sub(outSum)
	for idx := 0; idx <= len(deltas)-1; idx++ {
		if dust.LTE(availablePerValidator[deltas[idx].ValoperAddress]) {
			outWeights[deltas[idx].ValoperAddress] = outWeights[deltas[idx].ValoperAddress].Add(dust)
			break
		}
	}

	return filter(outWeights), nil
}
