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

	candidates := CandidateBins{}
	threshold := bins.DetermineThreshold()
	for _, delAddr := range bins.Keys() {
		bin := bins[delAddr]
		binVal := bin.SumDelegations()
		if bin.HasValidator(validatorAddress) {
			// already contains
			if binVal.GTE(threshold) {
				// oversubscribed :(
				candidates = append(candidates, CandidateBin{delAddr, binVal})
			} else {
				return delAddr, bins.AddDelegation(delAddr, validatorAddress, coin.Amount)
			}
		} else {
			// bin does not have this validator in...
			if bin.IsEmpty() {
				return delAddr, bins.AddDelegation(delAddr, validatorAddress, coin.Amount)
			}
		}
	}

	smallest := bins.SmallestBin()
	if len(candidates) > 0 {
		candidates = candidates.GetSorted()

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

type SendPlan map[string]sdk.Coins

func (s SendPlan) Keys() []string {
	out := []string{}

	for i := range s {
		out = append(out, i)
	}

	sort.Strings(out)

	return out
}

func (distrPlan DistributionPlan) ToSendPlan() SendPlan {
	out := make(SendPlan)
	for _, delegateAccount := range distrPlan.Keys() {
		if plan, ok := distrPlan.Value[delegateAccount]; ok {
			for _, planKey := range plan.Keys() {
				if _, ok := out[delegateAccount]; !ok {
					out[delegateAccount] = plan.Value[planKey].Value
				} else {
					out[delegateAccount] = out[delegateAccount].Add(plan.Value[planKey].Value...)
				}
			}
		}
	}
	return out
}

func (distrPlan DistributionPlan) Keys() []string {
	out := []string{}

	for i := range distrPlan.GetValue() {
		out = append(out, i)
	}

	sort.Strings(out)

	return out
}

func (delPlan DelegationPlan) Keys() []string {
	out := []string{}

	for i := range delPlan.GetValue() {
		out = append(out, i)
	}

	sort.Strings(out)

	return out
}

func (distrPlan DistributionPlan) Add(delegationAddress string, plan *DelegationPlan) DistributionPlan {
	if distrPlan.Value == nil {
		distrPlan.Value = make(map[string]*DelegationPlan)
	}
	if _, ok := distrPlan.Value[delegationAddress]; !ok {
		distrPlan.Value[delegationAddress] = NewEmptyDelegationPlan()
	}
	distrPlan.Value[delegationAddress] = distrPlan.Value[delegationAddress].Merge(plan)
	return distrPlan
}

func (distrPlan DistributionPlan) Merge(b DistributionPlan) DistributionPlan {
	for _, delegationAddress := range b.Keys() {
		if plan, ok := distrPlan.Value[delegationAddress]; ok {
			distrPlan.Value[delegationAddress] = plan.Merge(b.GetValue()[delegationAddress])
		} else {
			distrPlan.Value[delegationAddress] = b.GetValue()[delegationAddress]
		}
	}
	return distrPlan
}

func (delegationPlan DelegationPlan) Sum() sdk.Coins {
	out := sdk.Coins{}
	for _, planItem := range delegationPlan.Value {
		out = out.Add(planItem.Value...)
	}
	return out

}

func NewEmptyDelegationPlan() *DelegationPlan {
	return &DelegationPlan{Value: map[string]*DelegationPlan_DelegationPlanItem{}}
}

func NewSingleDelegationPlan(validator string, amount sdk.Coins) *DelegationPlan {
	out := NewEmptyDelegationPlan()
	out.Value[validator] = &DelegationPlan_DelegationPlanItem{Value: amount}
	return out
}

func NewPlanItem(coins sdk.Coins) *DelegationPlan_DelegationPlanItem {
	return &DelegationPlan_DelegationPlanItem{Value: coins}
}

func (delegationPlan DelegationPlan) Merge(b *DelegationPlan) *DelegationPlan {
	for _, validatorAddress := range b.Keys() {
		if existingCoins, ok := delegationPlan.Value[validatorAddress]; ok {
			delegationPlan.Value[validatorAddress] = NewPlanItem(existingCoins.Value.Add(b.GetValue()[validatorAddress].GetValue()...))
		} else {
			delegationPlan.Value[validatorAddress] = b.GetValue()[validatorAddress]
		}
	}
	return &delegationPlan
}

func DelegationPlanFromUserIntent(zone RegisteredZone, coin sdk.Coin, intent ValidatorIntents) (*DelegationPlan, error) {

	out := NewEmptyDelegationPlan()

	for _, val := range intent.Keys() {
		out.Value[val] = NewPlanItem(sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: sdk.Int(coin.Amount.ToDec().Mul(intent[val].Weight).TruncateInt())}})
	}

	return out, nil
}

func DelegationPlanFromGlobalIntent(currentState DelegationBins, zone RegisteredZone, coin sdk.Coin, intent ValidatorIntents) (*DelegationPlan, error) {
	if coin.Denom != zone.BaseDenom {
		return nil, fmt.Errorf("expected base denom, got %s", coin.Denom)
	}

	type Diff struct {
		valoper string
		amount  sdk.Int
	}

	deltas := []Diff{}

	out := NewEmptyDelegationPlan()

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

	distributableValue := coin.Amount

	for idx, delta := range deltas {
		if delta.amount.GT(sdk.ZeroInt()) {
			if delta.amount.GTE(distributableValue) {
				out.Value[delta.valoper] = NewPlanItem(sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: distributableValue}})
				distributableValue = sdk.ZeroInt()
				break
			} else {
				distributableValue = distributableValue.Sub(deltas[idx].amount)
				out.Value[delta.valoper] = NewPlanItem(sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: deltas[idx].amount}})
			}
		}
	}

	if distributableValue.GT(sdk.ZeroInt()) {
		for _, val := range intent.Keys() {
			valCoin := distributableValue.ToDec().Mul(intent[val].Weight).TruncateInt()
			distributableValue = distributableValue.Sub(valCoin)
			out.Value[val] = NewPlanItem(out.Value[val].Value.Add(sdk.NewCoin(zone.BaseDenom, valCoin)))
		}
	}

	if !out.Sum().IsEqual(sdk.Coins{coin}) {
		remainder := sdk.Coins{coin}.Sub(out.Sum())
		out.Value[deltas[len(deltas)-1].valoper] = NewPlanItem(out.Value[deltas[len(deltas)-1].valoper].Value.Add(remainder...))
	}

	return out, nil
}

func DelegationPlanFromCoins(zone RegisteredZone, coin sdk.Coin) *DelegationPlan {
	out := NewEmptyDelegationPlan()

	for _, val := range zone.GetValidatorsSorted() {
		if strings.HasPrefix(coin.Denom, val.ValoperAddress) {
			_, ok := out.Value[val.ValoperAddress]
			if !ok {
				out.Value[val.ValoperAddress] = NewPlanItem(sdk.NewCoins(coin))
			} else {
				out.Value[val.ValoperAddress] = NewPlanItem(out.Value[val.ValoperAddress].Value.Add(coin))
			}
			break
		}
	}

	return out
}
