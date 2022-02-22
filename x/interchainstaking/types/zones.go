package types

import (
	fmt "fmt"
	"sort"
)

func (z RegisteredZone) GetDelegationAccountsByLowestBalance(qty int64) []*ICAAccount {
	delegationAccounts := z.DelegationAddresses
	sort.Slice(delegationAccounts, func(i, j int) bool {
		return delegationAccounts[i].DelegatedBalance.Amount.GT(delegationAccounts[j].DelegatedBalance.Amount)
	})
	if qty > 0 {
		return delegationAccounts[:qty]
	}
	return delegationAccounts
}

func (z RegisteredZone) SupportMultiSend() bool { return false } // this should become part of the constructor/changable by governance

func (z RegisteredZone) GetValidatorByValoper(valoper string) (*Validator, error) {
	for _, v := range z.Validators {
		if v.ValoperAddress == valoper {
			return v, nil
		}
	}
	return nil, fmt.Errorf("invalid validator %s", valoper)
}

func (v Validator) GetDelegationForDelegator(delegator string) (*Delegation, error) {
	for _, d := range v.Delegations {
		if d.DelegationAddress == delegator {
			return d, nil
		}
	}
	return nil, fmt.Errorf("no delegation for for delegator %s", delegator)
}
