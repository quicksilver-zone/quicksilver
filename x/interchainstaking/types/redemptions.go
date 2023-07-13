package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DetermineAllocationsForUndelegation(currentAllocations map[string]math.Int, currentSum math.Int, targetAllocations ValidatorIntents, availablePerValidator map[string]math.Int, amount sdk.Coins) map[string]math.Int {
	input := amount[0].Amount
	deltas := CalculateDeltas(currentAllocations, currentSum.Sub(input), targetAllocations)
	sum := sdk.ZeroInt()
	outSum := sdk.ZeroInt()
	outWeights := make(map[string]math.Int)

	// deltas: +ve is below target; -ve is above target.

	// q1: can we satisfy this unbonding using _just_ above target allocations.
	// example:
	// we have v1: 5000, v2: 1800; v3: 1200; v4: 1000 and targets of 50%, 20%, 15% and 5% respectively.
	// deltas == 500, 0, -150, 550
	// an unbonding of 300 should come from v1, v4 (as they has an excess of > unbond amount) _before_ touching anything else.

	for idx := range deltas {
		if deltas[idx].Weight.IsNegative() {
			sum = sum.Add(deltas[idx].Weight.TruncateInt().Abs())
		}
	}

	overAllocationSplit := sdk.MinInt(sum, input)
	if !overAllocationSplit.IsZero() {
		for idx := range deltas {
			if !deltas[idx].Weight.IsNegative() {
				continue
			}
			outWeights[deltas[idx].ValoperAddress] = deltas[idx].Weight.Quo(sdk.NewDecFromInt(sum)).Mul(sdk.NewDecFromInt(overAllocationSplit)).TruncateInt().Abs()
			if outWeights[deltas[idx].ValoperAddress].GT(availablePerValidator[deltas[idx].ValoperAddress]) {
				outWeights[deltas[idx].ValoperAddress] = availablePerValidator[deltas[idx].ValoperAddress]
				availablePerValidator[deltas[idx].ValoperAddress] = sdk.ZeroInt()
			} else {
				availablePerValidator[deltas[idx].ValoperAddress] = availablePerValidator[deltas[idx].ValoperAddress].Sub(outWeights[deltas[idx].ValoperAddress])
			}
			deltas[idx].Weight = deltas[idx].Weight.Add(sdk.NewDecFromInt(outWeights[deltas[idx].ValoperAddress]))
			outSum = outSum.Add(outWeights[deltas[idx].ValoperAddress])

		}
	}
	input = input.Sub(outSum)
	if input.IsZero() {
		return outWeights
	}

	maxValue := MaxDelta(deltas)
	sum = sdk.ZeroInt()

	// drop all deltas such that the maximum value is zero, and invert.
	for idx := range deltas {
		deltas[idx].Weight = deltas[idx].Weight.Sub(sdk.NewDecFromInt(maxValue)).Abs()
		sum = sum.Add(deltas[idx].Weight.TruncateInt().Abs())

	}

	// unequalSplit is the portion of input that should be distributed in attempt to make targets == 0
	unequalSplit := sdk.MinInt(sum, input)

	if !unequalSplit.IsZero() {
		for idx := range deltas {
			allocation := deltas[idx].Weight.Quo(sdk.NewDecFromInt(sum)).Mul(sdk.NewDecFromInt(unequalSplit))
			_, ok := availablePerValidator[deltas[idx].ValoperAddress]
			if !ok {
				availablePerValidator[deltas[idx].ValoperAddress] = sdk.ZeroInt()
			}

			if allocation.TruncateInt().GT(availablePerValidator[deltas[idx].ValoperAddress]) {
				allocation = sdk.NewDecFromInt(availablePerValidator[deltas[idx].ValoperAddress])
				availablePerValidator[deltas[idx].ValoperAddress] = sdk.ZeroInt()
			} else {
				availablePerValidator[deltas[idx].ValoperAddress] = availablePerValidator[deltas[idx].ValoperAddress].Sub(allocation.TruncateInt())
			}

			deltas[idx].Weight = deltas[idx].Weight.Sub(allocation)
			value, ok := outWeights[deltas[idx].ValoperAddress]
			if !ok {
				value = sdk.ZeroInt()
			}

			outWeights[deltas[idx].ValoperAddress] = value.Add(allocation.TruncateInt())
			outSum = outSum.Add(allocation.TruncateInt())
			input = input.Sub(allocation.TruncateInt())
		}
	}

	// equalSplit is the portion of input that should be distributed equally across all validators, once targets are met.

	if !outSum.Equal(amount[0].Amount) {
		each := sdk.NewDecFromInt(input).Quo(sdk.NewDec(int64(len(deltas))))
		for idx := range deltas {
			value, ok := outWeights[deltas[idx].ValoperAddress]
			if !ok {
				value = sdk.ZeroInt()
			}
			if each.TruncateInt().GT(availablePerValidator[deltas[idx].ValoperAddress]) {
				each = sdk.NewDecFromInt(availablePerValidator[deltas[idx].ValoperAddress])
				availablePerValidator[deltas[idx].ValoperAddress] = sdk.ZeroInt()
			} else {
				availablePerValidator[deltas[idx].ValoperAddress] = availablePerValidator[deltas[idx].ValoperAddress].Sub(each.TruncateInt())
			}
			outWeights[deltas[idx].ValoperAddress] = value.Add(each.TruncateInt())
			outSum = outSum.Add(each.TruncateInt())
			input = input.Sub(each.TruncateInt())
		}
	}

	// dust is the portion of the input that was truncated in previous calculations; add this to the last validator in the list,
	// which should be the biggest source. This will always be a small amount, and will count toward the delta calculations on the next run.
	dust := amount[0].Amount.Sub(outSum)
	for idx := len(deltas) - 1; idx >= 0; idx-- {
		if dust.LTE(availablePerValidator[deltas[idx].ValoperAddress]) {
			outWeights[deltas[idx].ValoperAddress] = outWeights[deltas[idx].ValoperAddress].Add(dust)
			break
		}
	}

	return outWeights
}
