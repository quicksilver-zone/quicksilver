package types

import (
	"errors"
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

func (vi ValidatorIntents) GetForValoper(valoper string) (*ValidatorIntent, bool) {
	for _, i := range vi {
		if i.ValoperAddress == valoper {
			return i, true
		}
	}
	return nil, false
}

func (vi ValidatorIntents) SetForValoper(valoper string, intent *ValidatorIntent) ValidatorIntents {
	for idx, i := range vi {
		if i.ValoperAddress == valoper {
			vi[idx] = vi[len(vi)-1]
			vi = vi[:len(vi)-1]
			break
		}
	}
	vi = append(vi, intent)

	return vi.Sort()
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
		total = total.AddMut(i.Weight)
	}

	out := make(ValidatorIntents, 0)
	for _, i := range vi {
		out = append(out, &ValidatorIntent{ValoperAddress: i.ValoperAddress, Weight: i.Weight.Quo(total)})
	}
	return out.Sort()
}

func DetermineAllocationsForDelegation(currentAllocations map[string]sdkmath.Int, currentSum sdkmath.Int, targetAllocations ValidatorIntents, amount sdk.Coins) map[string]sdkmath.Int {
	input := amount[0].Amount
	deltas, sources := CalculateAllocationDeltas(currentAllocations, map[string]bool{}, currentSum, targetAllocations)
	// take targets and sources, and flip that shit.
	// sources -> negate -> join -> sort.
	largestSource := sources.MaxDelta()

	// negate all values in sources.
	sources.Negate()
	deltas = append(deltas, sources...)

	sum := sdk.ZeroInt()

	// raise all deltas such that the minimum value is zero.
	for idx := range deltas {
		deltas[idx].Amount = deltas[idx].Amount.Add(largestSource)
		// sum here instead of calling Sum() later to save looping over slice again.
		sum = sum.Add(deltas[idx].Amount)
	}

	// unequalSplit is the portion of input that should be distributed in attempt to make targets == 0
	unequalSplit := sdk.MinInt(sum, input)

	if !unequalSplit.IsZero() {
		for idx := range deltas {
			deltas[idx].Amount = sdk.NewDecFromInt(deltas[idx].Amount).QuoInt(sum).MulInt(unequalSplit).TruncateInt()
		}
	}

	// equalSplit is the portion of input that should be distributed equally across all validators, once targets are zero.
	equalSplit := sdk.NewDecFromInt(input.Sub(unequalSplit))

	// replace this portion with allocation proportional to targetAllocations!
	if !equalSplit.IsZero() {
		for _, targetAllocation := range targetAllocations.Sort() {
			delta, found := deltas.GetForValoper(targetAllocation.GetValoperAddress())
			if found {
				delta.Amount = delta.Amount.Add(equalSplit.Mul(targetAllocation.Weight).TruncateInt())
			} else {
				delta = &AllocationDelta{ValoperAddress: targetAllocation.GetValoperAddress(), Amount: equalSplit.Mul(targetAllocation.Weight).TruncateInt()}
				deltas = append(deltas, delta)
			}
		}
	}

	// dust is the portion of the input that was truncated in previous calculations; add this to the first validator in the list,
	// once sorted alphabetically. This will always be a small amount, and will count toward the delta calculations on the next run.

	outSum := sdk.ZeroInt()
	outWeights := make(map[string]sdkmath.Int)
	for _, delta := range deltas {
		outWeights[delta.ValoperAddress] = delta.Amount
		outSum = outSum.Add(delta.Amount)
	}
	dust := input.Sub(outSum)
	outWeights[deltas[0].ValoperAddress] = outWeights[deltas[0].ValoperAddress].Add(dust)

	return outWeights
}
