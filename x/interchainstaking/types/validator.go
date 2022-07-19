package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (v Validator) SharesToTokens(shares sdk.Dec) sdk.Int {
	return sdk.NewDecFromInt(v.VotingPower).Quo(v.DelegatorShares).TruncateInt()
}

func (di DelegatorIntent) AddOrdinal(multiplier sdk.Int, intents ValidatorIntents) DelegatorIntent {
	if len(intents) == 0 {
		return di
	}
	di.Ordinalize(multiplier)

OUTER:
	for _, idx := range intents.Keys() {
		if i, ok := intents[idx]; ok {
			for _, j := range di.Sorted() {
				if i.ValoperAddress == j.ValoperAddress {
					j.Weight = j.Weight.Add(i.Weight)
					continue OUTER
				}
			}
			di.Intents = append(di.Intents, i)
		}

	}

	return di.Normalize()
}

func (di DelegatorIntent) Normalize() DelegatorIntent {
	summedWeight := sdk.ZeroDec()
	for _, i := range di.Sorted() {
		summedWeight = summedWeight.Add(i.Weight)
	}
	for _, i := range di.Sorted() {
		i.Weight = i.Weight.QuoTruncate(summedWeight)
	}
	return di
}

func (di DelegatorIntent) Ordinalize(multiple sdk.Int) DelegatorIntent {
	for _, i := range di.Sorted() {
		i.Weight = i.Weight.MulInt(multiple)
	}
	return di
}

func (di DelegatorIntent) ToMap(multiple sdk.Int) map[string]sdk.Int {
	out := make(map[string]sdk.Int)
	di = di.Ordinalize(multiple)
	for _, i := range di.Sorted() {
		out[i.ValoperAddress] = i.Weight.TruncateInt()
	}
	return out
}

func (di DelegatorIntent) ToAllocations(multiple sdk.Int) Allocations {
	out := Allocations{}
	di = di.Ordinalize(multiple)
	for _, i := range di.Sorted() {
		out = out.Allocate(i.ValoperAddress, sdk.Coins{sdk.Coin{Denom: GenericToken, Amount: i.Weight.TruncateInt()}})
	}
	return out
}

func (di DelegatorIntent) ToValidatorIntents() ValidatorIntents {
	out := make(ValidatorIntents)
	for _, i := range di.Sorted() {
		out[i.ValoperAddress] = i
	}
	return out
}

func (d DelegatorIntent) Sorted() []*ValidatorIntent {
	sort.SliceStable(d.Intents, func(i, j int) bool {
		return d.Intents[i].ValoperAddress < d.Intents[j].ValoperAddress
	})
	return d.Intents
}
