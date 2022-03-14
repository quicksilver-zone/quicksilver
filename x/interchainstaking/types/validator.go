package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (v Validator) GetDelegationForDelegator(delegator string) (*Delegation, error) {
	for _, d := range v.Delegations {
		if d.DelegationAddress == delegator {
			return d, nil
		}
	}
	return nil, fmt.Errorf("no delegation for for delegator %s", delegator)
}

func (di DelegatorIntent) AddOrdinal(multiplier sdk.Int, intents map[string]*ValidatorIntent) DelegatorIntent {
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
	fmt.Println("Intents (not normalised)", "intents", di.Intents)

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
