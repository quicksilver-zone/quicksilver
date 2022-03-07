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
	for _, i := range intents {
		fmt.Println("Intent to add", "intent", i)

		for _, j := range di.Intents {
			if i.ValoperAddress == j.ValoperAddress {
				fmt.Println("Validator found", "valoper", j.ValoperAddress)
				fmt.Println("Adding intent", "weight", i.Weight)

				j.Weight = j.Weight.Add(i.Weight)
				continue
			}
		}

		// we don't have an intent for this validator yet!
		di.Intents = append(di.Intents, i)
	}
	fmt.Println("Intents (not normalised)", "intents", di.Intents)

	return di.Normalize()
}

func (di DelegatorIntent) Normalize() DelegatorIntent {
	summedWeight := sdk.ZeroDec()
	for _, i := range di.Intents {
		fmt.Println("Intent found", "intent", i)
		summedWeight = summedWeight.Add(i.Weight)
	}
	fmt.Println("Summed weight", "weight", summedWeight)

	validateWeight := sdk.ZeroDec()
	for _, i := range di.Intents {
		i.Weight = i.Weight.QuoTruncate(summedWeight)
		validateWeight = validateWeight.Add(i.Weight)
	}
	fmt.Println("Validate weight", "weight", validateWeight)
	if !validateWeight.LTE(sdk.OneDec()) && len(di.Intents) > 0 {
		panic("Normalize should equal 1")
	}

	return di
}

func (di DelegatorIntent) Ordinalize(multiple sdk.Int) DelegatorIntent {
	for _, i := range di.Intents {
		i.Weight = i.Weight.MulInt(multiple)
	}
	return di
}
