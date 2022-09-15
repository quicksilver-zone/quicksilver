package types

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr string, validatorAddr string, amount sdk.Coin) Delegation {
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

// return the delegation
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation Delegation, err error) {
	if bytes.Equal(value, []byte("")) {
		return Delegation{}, fmt.Errorf("unable to unmarshal zero-length byte slice")
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

// -------------------------------------------------------------------------
// DelegationCandidates

func (a Allocations) DetermineThreshold() sdk.Int {
	return a.SortedByAmount()[int(float64(0.33)*float64(len(a)))].SumAll()
}

func (a Allocations) SmallestBin() Allocation {
	return *a.SortedByAmount()[0]
}

func (a Allocations) FindAccountForDelegation(validatorAddress string, coin sdk.Coin) (string, Allocations) {
	candidates := Allocations{}
	threshold := a.DetermineThreshold()

	for _, bin := range a.SortedByAmount() {
		binVal := bin.SumAll()
		if bin.Amount.AmountOf(validatorAddress).GT(sdk.ZeroInt()) { // does this allocation contain any valoper coins?
			// already contains
			if binVal.GTE(threshold) {
				// oversubscribed :(
				candidates = candidates.Allocate(bin.Address, bin.Amount)
			} else {
				return bin.Address, a.Allocate(bin.Address, sdk.Coins{sdk.Coin{Denom: validatorAddress, Amount: coin.Amount}})
			}
		} else {
			// bin does not have this validator in...
			if bin.Amount.IsZero() {
				return bin.Address, a.Allocate(bin.Address, sdk.Coins{sdk.Coin{Denom: validatorAddress, Amount: coin.Amount}})
			}
		}
	}

	smallest := a.SmallestBin()
	if len(candidates) > 0 {
		candidates = candidates.SortedByAmount()
		if smallest.SumAll().LT(candidates[0].SumAll().Quo(sdk.NewInt(3))) {
			return smallest.Address, a.Allocate(smallest.Address, sdk.Coins{sdk.Coin{Denom: validatorAddress, Amount: coin.Amount}})
		}
		return candidates[0].Address, a.Allocate(candidates[0].Address, sdk.Coins{sdk.Coin{Denom: validatorAddress, Amount: coin.Amount}})
	}
	return smallest.Address, a.Allocate(smallest.Address, sdk.Coins{sdk.Coin{Denom: validatorAddress, Amount: coin.Amount}})
}

// --------------------------------------------------------
// DelegationPlans

type ValidatorIntents map[string]*ValidatorIntent

func (v ValidatorIntents) Keys() []string {
	keys := make([]string, len(v))
	i := 0
	for key := range v {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	return keys
}

// MustMarshalDelegationPlan returns the delegation plan bytes.
// This function will panic on failure.
func MustMarshalDelegationPlan(cdc codec.BinaryCodec, delegationPlan DelegationPlan) []byte {
	return cdc.MustMarshal(&delegationPlan)
}

// MustUnmarshalDelegationPlan return the unmarshaled delegation plan from bytes.
// This function will panic on failure.
func MustUnmarshalDelegationPlan(cdc codec.BinaryCodec, value []byte) DelegationPlan {
	delegationPlan, err := UnmarshalDelegationPlan(cdc, value)
	if err != nil {
		panic(err)
	}

	return delegationPlan
}

// return the delegation plan
func UnmarshalDelegationPlan(cdc codec.BinaryCodec, value []byte) (delegationPlan DelegationPlan, err error) {
	if bytes.Equal(value, []byte("")) {
		return delegationPlan, fmt.Errorf("unable to unmarshal zero length byte slice")
	}
	err = cdc.Unmarshal(value, &delegationPlan)
	return delegationPlan, err
}

// This function will panic on failure.
func (d DelegationPlan) GetDelegatorAddr() sdk.AccAddress {
	_, delAddr, err := bech32.DecodeAndConvert(d.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return delAddr
}

// This function will panic on failure.
func (d DelegationPlan) GetValidatorAddr() sdk.ValAddress {
	_, valAddr, err := bech32.DecodeAndConvert(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return valAddr
}

func NewDelegationPlan(delAddr, valAddr string, amount sdk.Coins) DelegationPlan {
	return DelegationPlan{DelegatorAddress: delAddr, ValidatorAddress: valAddr, Value: amount}
}

func DelegationPlanFromUserIntent(zone Zone, coin sdk.Coin, intent ValidatorIntents) Allocations {
	out := Allocations{}

	for _, val := range intent.Keys() {
		out = out.Allocate(val, sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: sdk.NewDecFromInt(coin.Amount).Mul(intent[val].Weight).TruncateInt()}})
	}
	return out
}

type Allocation struct {
	Address string
	Amount  sdk.Coins
}

func (a Allocations) Allocate(address string, amount sdk.Coins) Allocations {
	for _, allocation := range a {
		if allocation.Address == address {
			allocation.Amount = allocation.Amount.Add(amount...)
			return a
		}
	}
	return append(a, &Allocation{Address: address, Amount: amount})
}

func (a Allocations) Get(address string) *Allocation {
	for _, allocation := range a {
		if allocation.Address == address {
			return allocation
		}
	}
	return nil
}

func (a Allocations) Sorted() Allocations {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Address < a[j].Address
	})

	return a
}

func (a Allocations) SortedByAmount() Allocations {
	a = a.Sorted() // sort by address first so that sorting on amount is deterministic.
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].SumAll().LT(a[j].SumAll())
	})

	return a
}

func (a Allocations) Sum() sdk.Coins {
	out := sdk.Coins{}
	for _, allocation := range a {
		out = out.Add(allocation.Amount...)
	}
	return out
}

// remove amount from address. Return the amount that could not be subtracted.
func (a Allocations) Sub(amount sdk.Coins, address string) (Allocations, sdk.Coins) {
	if allocation := a.Get(address); allocation != nil {
		subAmount := allocation.Amount
		for _, coin := range amount {
			var amountToSub sdk.Coins
			if subAmount.AmountOf(coin.Denom).GTE(coin.Amount) {
				amountToSub = sdk.Coins{coin}
			} else {
				amountToSub = sdk.Coins{sdk.NewCoin(coin.Denom, subAmount.AmountOf(coin.Denom))}
			}
			subAmount = subAmount.Sub(amountToSub...)
			amount = amount.Sub(amountToSub...)
		}
		allocation.Amount = subAmount
	}
	return a, amount
}

func (a Allocations) SumForDenom(denom string) sdk.Int {
	out := sdk.ZeroInt()
	for _, allocation := range a {
		out = out.Add(allocation.Amount.AmountOf(denom))
	}
	return out
}

func (a Allocation) SumAll() sdk.Int {
	// warning: this treats all denoms as fungible. It might not be what you want to do!
	out := sdk.ZeroInt()
	for _, coin := range a.Amount {
		out = out.Add(coin.Amount)
	}
	return out
}

func (a Allocations) SumAll() sdk.Int {
	// warning: this treats all denoms as fungible. It might not be what you want to do!
	out := sdk.ZeroInt()
	for _, allocation := range a {
		for _, coin := range allocation.Amount {
			out = out.Add(coin.Amount)
		}
	}
	return out
}

type (
	Allocations []*Allocation
	Diffs       []*Diff
)

func (a Diffs) Sorted() Diffs {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Valoper < a[j].Valoper
	})

	return a
}

func (a Diffs) SortedByAmount() Diffs {
	a = a.Sorted() // sort by address first so that sorting on amount is deterministic.
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Amount.LT(a[j].Amount)
	})

	return a
}

func DetermineIntentDelta(currentState Allocations, total sdk.Int, intent ValidatorIntents) Diffs {
	deltas := Diffs{}

	if total.IsZero() {
		return deltas
	}

	for _, val := range intent.Keys() {
		current := currentState.SumForDenom(val)                                     // fetch current delegations to validator
		percent := sdk.NewDecFromInt(current).Quo(sdk.NewDecFromInt(total))          // what is this a percent of total + new
		deltaToIntent := intent[val].Weight.Sub(percent).MulInt(total).TruncateInt() // what to we have to delegate to make it match intent?
		deltas = append(deltas, &Diff{val, deltaToIntent})
	}

	return deltas.SortedByAmount()
}

type Diff struct {
	Valoper string
	Amount  sdk.Int
}

func DelegationPlanFromGlobalIntent(currentTotal sdk.Coin, currentState Allocations, coin sdk.Coin, intent ValidatorIntents) (Allocations, error) {
	if coin.Denom != currentTotal.Denom {
		return nil, fmt.Errorf("expected base denom, got %s", coin.Denom)
	}

	allocations := Allocations{}

	deltas := DetermineIntentDelta(currentState, currentTotal.Amount.Add(coin.Amount), intent)

	distributableValue := coin.Amount

	for idx, delta := range deltas {
		if delta.Amount.GT(sdk.ZeroInt()) {
			if delta.Amount.GTE(distributableValue) {
				allocations = allocations.Allocate(delta.Valoper, sdk.Coins{sdk.Coin{Denom: currentTotal.Denom, Amount: distributableValue}})
				distributableValue = sdk.ZeroInt()
				break
			} else {
				allocations = allocations.Allocate(delta.Valoper, sdk.Coins{sdk.Coin{Denom: currentTotal.Denom, Amount: deltas[idx].Amount}})
				distributableValue = distributableValue.Sub(deltas[idx].Amount)
			}
		}
	}

	if distributableValue.GT(sdk.ZeroInt()) {
		for _, val := range intent.Keys() {
			valCoin := sdk.NewDecFromInt(distributableValue).Mul(intent[val].Weight).TruncateInt()
			distributableValue = distributableValue.Sub(valCoin)
			allocations = allocations.Allocate(val, sdk.Coins{sdk.NewCoin(currentTotal.Denom, valCoin)})
		}
	}

	if !allocations.Sum().IsEqual(sdk.Coins{coin}) {
		remainder := sdk.Coins{coin}.Sub(allocations.Sum()...)
		allocations = allocations.Allocate(deltas[len(deltas)-1].Valoper, remainder)
	}
	return allocations, nil
}

func DelegationPlanFromCoins(zone Zone, coin sdk.Coin) Allocations {
	out := Allocations{}

	for _, val := range zone.GetValidatorsSorted() {
		if strings.HasPrefix(coin.Denom, val.ValoperAddress) {
			out = out.Allocate(val.ValoperAddress, sdk.NewCoins(coin))
			break
		}
	}

	return out
}
