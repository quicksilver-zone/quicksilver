package types

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// NewDelegation creates a new delegation object
//nolint:interfacer
func NewDelegation(delegatorAddr string, validatorAddr string, amount sdk.Coin) Delegation {
	return Delegation{
		DelegationAddress: delegatorAddr,
		ValidatorAddress:  validatorAddr,
		Amount:            amount,
		Height:            0,
		RedelegationEnd:   0,
	}
}

// MustMarshalDelegation returns the delegation bytes. Panics if fails
func MustMarshalDelegation(cdc codec.BinaryCodec, delegation Delegation) []byte {
	return cdc.MustMarshal(&delegation)
}

// MustUnmarshalDelegation return the unmarshaled delegation from bytes.
// Panics if fails.
func MustUnmarshalDelegation(cdc codec.BinaryCodec, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}

	return delegation
}

// return the delegation
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation Delegation, err error) {
	err = cdc.Unmarshal(value, &delegation)
	return delegation, err
}

func (d Delegation) GetDelegatorAddr() sdk.AccAddress {
	_, delAddr, err := bech32.DecodeAndConvert(d.DelegationAddress)
	if err != nil {
		panic(err)
	}
	return delAddr
}

func (d Delegation) GetValidatorAddr() sdk.ValAddress {
	_, valAddr, err := bech32.DecodeAndConvert(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return valAddr
}

// -------------------------------------------------------------------------
// DelegationCandidates

type CandidateBins []CandidateBin

type CandidateBin struct {
	addr string
	val  sdk.Int
}

func (bins CandidateBins) GetSorted() CandidateBins {
	sort.Slice(bins, func(i, j int) bool {
		return bins[i].val.LT(bins[j].val)
	})
	return bins
}

type DelegationBin map[string]sdk.Int

func (bin DelegationBin) SumDelegations() sdk.Int {
	sum := sdk.ZeroInt()
	for _, delegation := range bin {
		sum = sum.Add(delegation)
	}
	return sum
}

func (bin DelegationBin) IsEmpty() bool {
	return len(bin) == 0
}

func (bin DelegationBin) HasValidator(valoperAddress string) bool {
	_, ok := bin[valoperAddress]
	return ok
}

type DelegationBins map[string]DelegationBin

func (bins DelegationBins) AddDelegation(valoperAddress string, delegationAddress string, amount sdk.Int) DelegationBins {
	if _, ok := bins[delegationAddress]; !ok {
		bins[delegationAddress] = DelegationBin{}
	}
	if bins[delegationAddress].HasValidator(valoperAddress) {
		bins[delegationAddress][valoperAddress] = bins[delegationAddress][valoperAddress].Add(amount)
	} else {
		bins[delegationAddress][valoperAddress] = amount
	}
	return bins
}

func (bins DelegationBins) Keys() []string {
	keys := []string{}
	for k := range bins {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (bins DelegationBins) GetSorted() CandidateBins {
	sortedBins := CandidateBins{}
	for binAddr, bin := range bins {
		sortedBins = append(sortedBins, CandidateBin{binAddr, bin.SumDelegations()})
	}
	sort.Slice(sortedBins, func(i, j int) bool {
		return sortedBins[i].val.LT(sortedBins[j].val)
	})
	return sortedBins
}

func (bins DelegationBins) DetermineThreshold() sdk.Int {

	return bins.GetSorted()[int(float64(0.33)*float64(len(bins)))].val
}

func (bins DelegationBins) SmallestBin() CandidateBin {
	return bins.GetSorted()[0]
}

func (bins DelegationBins) FindAccountForDelegation(validatorAddress string, coin sdk.Coin) (string, DelegationBins) {
	fmt.Println("Finding bin: ", coin)
	candidates := CandidateBins{}
	threshold := bins.DetermineThreshold()
	fmt.Println("Threshold is: ", threshold)

	for _, delAddr := range bins.Keys() {
		bin := bins[delAddr]
		binVal := bin.SumDelegations()
		if bin.HasValidator(validatorAddress) {
			// already contains
			if binVal.GTE(threshold) {
				fmt.Printf("Binval %s >= threshold %s adding as candidate\n", binVal, threshold)
				// oversubscribed :(
				candidates = append(candidates, CandidateBin{delAddr, binVal})
			} else {
				fmt.Printf("Binval %s < threshold %s; using\n", binVal, threshold)
				return delAddr, bins.AddDelegation(delAddr, validatorAddress, coin.Amount)
			}
		} else {
			// bin does not have this validator in...
			if bin.IsEmpty() {
				fmt.Println("Bin is empty; using")
				return delAddr, bins.AddDelegation(delAddr, validatorAddress, coin.Amount)
			}
		}
	}

	smallest := bins.SmallestBin()
	if len(candidates) > 0 {
		candidates = candidates.GetSorted()
		fmt.Println("Candidates: ", candidates)
		fmt.Println("Smallest: ", smallest)
		if smallest.val.LT(candidates[0].val.Quo(sdk.NewInt(3))) {
			return smallest.addr, bins.AddDelegation(smallest.addr, validatorAddress, coin.Amount)
		} else {
			return candidates[0].addr, bins.AddDelegation(candidates[0].addr, validatorAddress, coin.Amount)
		}
	} else {
		return smallest.addr, bins.AddDelegation(smallest.addr, validatorAddress, coin.Amount)
	}
}

func (bins DelegationBins) SumForValidator(valoper string) sdk.Int {
	out := sdk.ZeroInt()

	for _, bin := range bins {
		val, ok := bin[valoper]
		if ok {
			out = out.Add(val)
		}
	}

	return out
}

// --------------------------------------------------------
// DelegationPlans

type ValidatorIntents map[string]*ValidatorIntent

func (v ValidatorIntents) Keys() []string {
	out := []string{}

	for i := range v {
		out = append(out, i)
	}

	sort.Strings(out)

	return out
}

// MustMarshalDelegationPlan returns the delegation plan bytes. Panics if fails
func MustMarshalDelegationPlan(cdc codec.BinaryCodec, delegationPlan DelegationPlan) []byte {
	return cdc.MustMarshal(&delegationPlan)
}

// MustUnmarshalDelegationPlan return the unmarshaled delegation plan from bytes.
// Panics if fails.
func MustUnmarshalDelegationPlan(cdc codec.BinaryCodec, value []byte) DelegationPlan {
	delegationPlan, err := UnmarshalDelegationPlan(cdc, value)
	if err != nil {
		panic(err)
	}

	return delegationPlan
}

// return the delegation plan
func UnmarshalDelegationPlan(cdc codec.BinaryCodec, value []byte) (delegationPlan DelegationPlan, err error) {
	err = cdc.Unmarshal(value, &delegationPlan)
	return delegationPlan, err
}

func (d DelegationPlan) GetDelegatorAddr() sdk.AccAddress {
	_, delAddr, err := bech32.DecodeAndConvert(d.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return delAddr
}

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

func DelegationPlanFromUserIntent(zone RegisteredZone, coin sdk.Coin, intent ValidatorIntents) Allocations {

	out := Allocations{}

	for _, val := range intent.Keys() {
		out = out.Allocate(val, sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: sdk.Int(coin.Amount.ToDec().Mul(intent[val].Weight).TruncateInt())}})
	}
	fmt.Println("DelegationPlanFromUserIntent", out)

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
	return append(a, Allocation{Address: address, Amount: amount})
}

func (a Allocations) Sorted() Allocations {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Address < a[j].Address
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

type DelegationPlans []DelegationPlan
type Allocations []Allocation

func (d DelegationPlans) Add(valAddress string, delAddress string, amount sdk.Coins) DelegationPlans {
	for _, delPlan := range d {
		if delPlan.ValidatorAddress == valAddress && delPlan.DelegatorAddress == delAddress {
			delPlan.Value = delPlan.Value.Add(amount...)
			return d
		}
	}
	return append(d, DelegationPlan{ValidatorAddress: valAddress, DelegatorAddress: delAddress, Value: amount})
}

func DelegationPlanFromGlobalIntent(currentState DelegationBins, zone RegisteredZone, coin sdk.Coin, intent ValidatorIntents) (Allocations, error) {
	if coin.Denom != zone.BaseDenom {
		return nil, fmt.Errorf("expected base denom, got %s", coin.Denom)
	}

	type Diff struct {
		valoper string
		amount  sdk.Int
	}

	deltas := []Diff{}
	allocations := Allocations{}

	// fetch current state
	total := zone.GetDelegatedAmount().Amount

	for _, val := range intent.Keys() {
		current := currentState.SumForValidator(val)                                                  // fetch current delegations to validator
		percent := current.ToDec().Quo(total.Add(coin.Amount).ToDec())                                // what is this a percent of total + new
		deltaToIntent := intent[val].Weight.Sub(percent).MulInt(total.Add(coin.Amount)).TruncateInt() // what to we have to delegate to make it match intent?
		deltas = append(deltas, Diff{val, deltaToIntent})
	}

	// determinism baby!
	sort.Slice(deltas, func(i, j int) bool {
		return deltas[i].amount.LT(deltas[j].amount)
	})

	fmt.Println("deltas: ", deltas)

	distributableValue := coin.Amount

	for idx, delta := range deltas {
		if delta.amount.GT(sdk.ZeroInt()) {
			if delta.amount.GTE(distributableValue) {
				allocations = allocations.Allocate(delta.valoper, sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: distributableValue}})
				distributableValue = sdk.ZeroInt()
				break
			} else {
				allocations = allocations.Allocate(delta.valoper, sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: deltas[idx].amount}})
				distributableValue = distributableValue.Sub(deltas[idx].amount)
			}
		}
	}

	if distributableValue.GT(sdk.ZeroInt()) {
		for _, val := range intent.Keys() {
			valCoin := distributableValue.ToDec().Mul(intent[val].Weight).TruncateInt()
			distributableValue = distributableValue.Sub(valCoin)
			allocations = allocations.Allocate(val, sdk.Coins{sdk.NewCoin(zone.BaseDenom, valCoin)})
		}
	}

	if !allocations.Sum().IsEqual(sdk.Coins{coin}) {
		remainder := sdk.Coins{coin}.Sub(allocations.Sum())
		allocations = allocations.Allocate(deltas[len(deltas)-1].valoper, remainder)
	}

	fmt.Println("DelegationPlanFromGlobalIntent", allocations)

	return allocations, nil
}

func DelegationPlanFromCoins(zone RegisteredZone, coin sdk.Coin) Allocations {
	out := Allocations{}

	for _, val := range zone.GetValidatorsSorted() {
		if strings.HasPrefix(coin.Denom, val.ValoperAddress) {
			out = out.Allocate(val.ValoperAddress, sdk.NewCoins(coin))
			break
		}
	}

	fmt.Println("DelegationPlanFromCoins", out)

	return out
}
