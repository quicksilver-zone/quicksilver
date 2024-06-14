package types

import (
	"errors"
	"fmt"
	"sort"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// NewDelegation creates a new delegation object.
func NewDelegation(delegatorAddr, validatorAddr string, amount sdk.Coin) Delegation {
	return Delegation{
		DelegationAddress: delegatorAddr,
		ValidatorAddress:  validatorAddr,
		Amount:            amount,
		Height:            0,
		RedelegationEnd:   0,
	}
}

// MustMarshalDelegation returns the delegation bytes.
// This function will panic on failure.
func MustMarshalDelegation(cdc codec.BinaryCodec, delegation Delegation) []byte {
	return cdc.MustMarshal(&delegation)
}

// MustUnmarshalDelegation return the unmarshaled delegation from bytes.
// This function will panic on failure.
func MustUnmarshalDelegation(cdc codec.BinaryCodec, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}

	return delegation
}

// UnmarshalDelegation return the delegation.
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation Delegation, err error) {
	if len(value) == 0 {
		return Delegation{}, errors.New("unable to unmarshal zero-length byte slice")
	}
	err = cdc.Unmarshal(value, &delegation)
	return delegation, err
}

// This function will panic on failure.
func (d Delegation) GetDelegatorAddr() sdk.AccAddress {
	_, delAddr, err := bech32.DecodeAndConvert(d.DelegationAddress)
	if err != nil {
		panic(err)
	}
	return delAddr
}

// This function will panic on failure.
func (d Delegation) GetValidatorAddr() sdk.ValAddress {
	_, valAddr, err := bech32.DecodeAndConvert(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return valAddr
}

type ValidatorIntents []*ValidatorIntent

func (vi ValidatorIntents) Sort() ValidatorIntents {
	sort.SliceStable(vi, func(i, j int) bool {
		return vi[i].ValoperAddress < vi[j].ValoperAddress
	})
	return vi
}

func (vi ValidatorIntents) Remove(valoper string) ValidatorIntents {
	for i, v := range vi {
		if v.ValoperAddress == valoper {
			vi[i] = vi[len(vi)-1]
			return vi[:len(vi)-1]
		}
	}
	return vi
}

func (vi ValidatorIntents) GetForValoper(valoper string) (*ValidatorIntent, bool) {
	for _, i := range vi {
		if i.ValoperAddress == valoper {
			return i, true
		}
	}
	return nil, false
}

func (vi ValidatorIntents) SetForValoper(valoper string, intent *ValidatorIntent) ValidatorIntents {
	idx := -1 // the index of the valoper if found
	for i, v := range vi {
		// Search for the valoper.
		if v.ValoperAddress == valoper {
			idx = i
			break
		}
	}

	if idx >= 0 { // We found the valoper so just replace it
		vi[idx] = intent
	} else {
		vi = append(vi, intent)
		return vi.Sort()
	}
	return vi
}

func (vi ValidatorIntents) MustGetForValoper(valoper string) *ValidatorIntent {
	intent, found := vi.GetForValoper(valoper)
	if !found || intent == nil {
		return &ValidatorIntent{ValoperAddress: valoper, Weight: sdk.ZeroDec()}
	}
	return intent
}

func (vi ValidatorIntents) Normalize() ValidatorIntents {
	total := sdk.ZeroDec()
	for _, i := range vi {
		if !i.Weight.IsNil() {
			total = total.AddMut(i.Weight)
		}
	}

	out := make(ValidatorIntents, 0)

	if total.IsZero() {
		return out
	}

	for _, i := range vi {
		if !i.Weight.IsNil() {
			out = append(out, &ValidatorIntent{ValoperAddress: i.ValoperAddress, Weight: i.Weight.Quo(total)})
		}
	}
	return out.Sort()
}

func DetermineAllocationsForDelegation(currentAllocations map[string]sdkmath.Int, currentSum sdkmath.Int, targetAllocations ValidatorIntents, amount sdk.Coins, maxCanAllocate map[string]sdkmath.Int) (map[string]sdkmath.Int, error) {
	if amount.IsZero() {
		return make(map[string]sdkmath.Int, 0), fmt.Errorf("unable to delegate zero amount")
	}
	if len(targetAllocations) == 0 {
		return make(map[string]sdkmath.Int, 0), fmt.Errorf("unable to process nil delegation targets")
	}
	input := amount[0].Amount
	deltas, _ := CalculateAllocationDeltas(currentAllocations, map[string]bool{}, currentSum.Add(amount[0].Amount), targetAllocations, maxCanAllocate)
	sum := deltas.Sum()

	// unequalAllocation is the portion of input that should be distributed in attempt to make targets == 0 (that is, in line with intent).
	unequalAllocation := sdk.MinInt(sum, input)

	if !unequalAllocation.IsZero() {
		for idx := range deltas {
			deltas[idx].Amount = sdk.NewDecFromInt(deltas[idx].Amount).QuoInt(sum).MulInt(unequalAllocation).TruncateInt()
		}
	}

	// proportionalAllocation is the portion of input that should be distributed proportionally to intent,  once targets are zero, respecting caps.
	proportionalAllocation := sdk.NewDecFromInt(input.Sub(unequalAllocation))

	rounds := 0
	// set maximum number of rounds, in case we get stuck in a weird loop we cannot resolve. If we exit the after this point, the remainder will be treated as dust.
	maxRounds := 10
	for ok := proportionalAllocation.IsPositive(); ok; ok = proportionalAllocation.IsPositive() && rounds < maxRounds {
		// normalise targetAllocations, so maxed caps are handled nicely.
		targetAllocations = targetAllocations.Normalize().Sort()
		// initialise roundAllocation
		roundAllocation := sdk.ZeroInt()
		// for each target
		for _, targetAllocation := range targetAllocations {
			// does this target validator have a cap?
			max, hasMax := maxCanAllocate[targetAllocation.ValoperAddress]
			// does it have an existing allocation?
			delta, found := deltas.GetForValoper(targetAllocation.GetValoperAddress())
			if !found {
				// no existing delta, create new delta with zero
				delta = &AllocationDelta{ValoperAddress: targetAllocation.GetValoperAddress(), Amount: sdk.ZeroInt()}
				deltas = append(deltas, delta)
			}
			// allocate to this validator based on weight
			thisAllocation := proportionalAllocation.Mul(targetAllocation.Weight).TruncateInt()
			// if there is a cap...
			if hasMax {
				// belt and braces.
				if max.LT(sdk.ZeroInt()) {
					return nil, errors.New("maxCanAllocate underflow")
				}
				// determine if cap is breached
				if delta.Amount.Add(thisAllocation).GTE(max) {
					// if so, truncate and remove from target allocations for next round
					thisAllocation = max.Sub(delta.Amount)
					delta.Amount = max
					targetAllocations = targetAllocations.Remove(delta.ValoperAddress)
				} else {
					// if not, increase delta
					delta.Amount = delta.Amount.Add(thisAllocation)
				}
			}
			// track round allocations to deduct from running total
			roundAllocation = roundAllocation.Add(thisAllocation)
		}
		// deduct from running total
		proportionalAllocation = proportionalAllocation.Sub(sdk.NewDecFromInt(roundAllocation))
		// bail after N rounds
		rounds++
	}

	// dust is the portion of the input that was truncated in previous calculations; add this to the first validator in the list,
	// once sorted alphabetically. This will always be a small amount, and will count toward the delta calculations on the next run.

	outSum := sdk.ZeroInt()
	outWeights := make(map[string]sdkmath.Int)
	for _, delta := range deltas {
		if !delta.Amount.IsZero() {
			outWeights[delta.ValoperAddress] = delta.Amount
			outSum = outSum.Add(delta.Amount)
		}
	}
	if outSum.GT(input) {
		return nil, errors.New("outSum overflow; cannot be greater than input amount")
	}

	// dust := input.Sub(outSum)
	// if !dust.IsZero() {
	// 	outWeights[deltas[0].ValoperAddress] = outWeights[deltas[0].ValoperAddress].Add(dust)
	// }

	return outWeights, nil
}
