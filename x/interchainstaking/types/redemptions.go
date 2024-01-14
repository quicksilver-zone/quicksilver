package types

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DetermineAllocationsForUndelegation(currentAllocations map[string]math.Int, lockedAllocations map[string]bool, currentSum math.Int, targetAllocations ValidatorIntents, availablePerValidator map[string]math.Int, amount sdk.Coins) map[string]math.Int {
	// this is brooooken
	input := amount[0].Amount
	underAllocated, overAllocated := CalculateAllocationDeltas(currentAllocations, lockedAllocations, currentSum /* .Sub(input) */, targetAllocations, make(map[string]math.Int))
	outSum := sdkmath.ZeroInt()
	outWeights := make(map[string]math.Int)

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
			outWeights[overAllocated[idx].ValoperAddress] = sdkmath.LegacyNewDecFromInt(overAllocated[idx].Amount.Add(math.OneInt())).Quo(sdkmath.LegacyNewDecFromInt(sum)).Mul(sdkmath.LegacyNewDecFromInt(overAllocationSplit)).TruncateInt()
			if outWeights[overAllocated[idx].ValoperAddress].GT(availablePerValidator[overAllocated[idx].ValoperAddress]) {
				// use up all of overAllocated[idx] and set available to zero.
				outWeights[overAllocated[idx].ValoperAddress] = availablePerValidator[overAllocated[idx].ValoperAddress]
				availablePerValidator[overAllocated[idx].ValoperAddress] = sdkmath.ZeroInt()
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
		return outWeights
	}

	// negate all values in underallocated.
	underAllocated.Negate()
	// append the two slices
	// nolint:gocritic
	deltas := append(overAllocated, underAllocated...)
	deltas.Sort()

	maxValue := deltas.MaxDelta()
	sum = sdkmath.ZeroInt()

	// drop all deltas such that the maximum value is zero, and invert.
	for idx := range deltas {
		deltas[idx].Amount = deltas[idx].Amount.Add(maxValue)
		// sum here instead of calling Sum() later to save looping over slice again.
		sum = sum.Add(deltas[idx].Amount.Abs())
	}

	// unequalSplit is the portion of input that should be distributed in attempt to make deltas == 0
	unequalSplit := sdkmath.MinInt(sum, input)

	if !unequalSplit.IsZero() {
		for idx := range deltas {
			allocation := sdkmath.LegacyNewDecFromInt(deltas[idx].Amount).Quo(sdkmath.LegacyNewDecFromInt(sum)).Mul(sdkmath.LegacyNewDecFromInt(unequalSplit)).TruncateInt()
			_, ok := availablePerValidator[deltas[idx].ValoperAddress]
			if !ok {
				availablePerValidator[deltas[idx].ValoperAddress] = sdkmath.ZeroInt()
			}

			if allocation.GT(availablePerValidator[deltas[idx].ValoperAddress]) {
				allocation = availablePerValidator[deltas[idx].ValoperAddress]
				availablePerValidator[deltas[idx].ValoperAddress] = sdkmath.ZeroInt()
			} else {
				availablePerValidator[deltas[idx].ValoperAddress] = availablePerValidator[deltas[idx].ValoperAddress].Sub(allocation)
			}

			deltas[idx].Amount = deltas[idx].Amount.Sub(allocation)
			value, ok := outWeights[deltas[idx].ValoperAddress]
			if !ok {
				value = sdkmath.ZeroInt()
			}

			outWeights[deltas[idx].ValoperAddress] = value.Add(allocation)
			outSum = outSum.Add(allocation)
			input = input.Sub(allocation)
		}
	}

	// equalSplit is the portion of input that should be distributed across all validators proportion to intent, once targets are met.
	deltas.Sort()
	if !outSum.Equal(amount[0].Amount) {
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
			perValidatorAmount := sdkmath.LegacyNewDecFromInt(origin).Mul(weights[idx].Weight).TruncateInt()
			value, ok := outWeights[weights[idx].ValoperAddress]
			if !ok {
				value = sdkmath.ZeroInt()
			}
			if perValidatorAmount.GT(availablePerValidator[weights[idx].ValoperAddress]) {
				perValidatorAmount = availablePerValidator[weights[idx].ValoperAddress]
				availablePerValidator[weights[idx].ValoperAddress] = sdkmath.ZeroInt()
			} else {
				availablePerValidator[weights[idx].ValoperAddress] = availablePerValidator[weights[idx].ValoperAddress].Sub(perValidatorAmount)
			}
			outWeights[weights[idx].ValoperAddress] = value.Add(perValidatorAmount)
			outSum = outSum.Add(perValidatorAmount)
			input = input.Sub(perValidatorAmount)
		}
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

	return outWeights
}
