package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (v Validator) SharesToTokens(shares sdk.Dec) sdk.Int {
	return v.VotingPower.ToDec().Quo(v.DelegatorShares).TruncateInt()
}

func (di DelegatorIntent) AddOrdinal(multiplier sdk.Int, intents map[string]*ValidatorIntent) DelegatorIntent {
	if len(intents) == 0 {
		return di
	}
	di.Ordinalize(multiplier)

OUTER:
	for _, i := range intents {
		for _, j := range di.Intents {
			if i.ValoperAddress == j.ValoperAddress {
				j.Weight = j.Weight.Add(i.Weight)
				continue OUTER
			}
		}
		di.Intents = append(di.Intents, i)
	}

	return di.Normalize()
}

func (di DelegatorIntent) Normalize() DelegatorIntent {
	summedWeight := sdk.ZeroDec()
	for _, i := range di.Intents {
		summedWeight = summedWeight.Add(i.Weight)
	}
	for _, i := range di.Intents {
		i.Weight = i.Weight.QuoTruncate(summedWeight)
	}
	return di
}

func (di DelegatorIntent) Ordinalize(multiple sdk.Int) DelegatorIntent {
	for _, i := range di.Intents {
		i.Weight = i.Weight.MulInt(multiple)
	}
	return di
}

func (di DelegatorIntent) ToMap(multiple sdk.Int) map[string]sdk.Int {
	out := make(map[string]sdk.Int)
	di = di.Ordinalize(multiple)
	for _, i := range di.Intents {
		out[i.ValoperAddress] = i.Weight.TruncateInt()
	}
	return out
}
