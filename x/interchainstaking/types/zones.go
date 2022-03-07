package types

import (
	fmt "fmt"
	math "math"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (z RegisteredZone) GetDelegationAccountsByLowestBalance(qty int64) []*ICAAccount {
	delegationAccounts := z.DelegationAddresses
	sort.Slice(delegationAccounts, func(i, j int) bool {
		return delegationAccounts[i].DelegatedBalance.Amount.GT(delegationAccounts[j].DelegatedBalance.Amount)
	})
	if qty > 0 {
		return delegationAccounts[:int(math.Min(float64(len(delegationAccounts)-1), float64(qty)))]
	}
	return delegationAccounts
}

func (z RegisteredZone) SupportMultiSend() bool { return z.MultiSend }

func (z RegisteredZone) GetValidatorByValoper(valoper string) (*Validator, error) {
	for _, v := range z.Validators {
		if v.ValoperAddress == valoper {
			return v, nil
		}
	}
	return nil, fmt.Errorf("invalid validator %s", valoper)
}

func (z RegisteredZone) GetDelegationsForDelegator(delegator string) []*Delegation {
	delegations := []*Delegation{}
	for _, v := range z.Validators {
		delegation, err := v.GetDelegationForDelegator(delegator)
		if err != nil {
			continue
		}
		delegations = append(delegations, delegation)
	}
	return delegations
}

func (z *RegisteredZone) ValidateCoinsForZone(ctx sdk.Context, coins sdk.Coins) error {

	zoneVals := z.GetValidatorsAsSlice()
COINS:
	for _, coin := range coins {
		if coin.Denom == z.BaseDenom {
			continue
		}

		for _, v := range zoneVals {
			if strings.HasPrefix(coin.Denom, v) {
				// continue 2 levels
				continue COINS
			}
		}
		return fmt.Errorf("invalid denom for zone: %s", coin.Denom)

	}
	return nil
}

func (z *RegisteredZone) ConvertCoinsToOrdinalIntents(ctx sdk.Context, coins sdk.Coins) map[string]*ValidatorIntent {
	// should we be return DelegatorIntent here?
	out := make(map[string]*ValidatorIntent)
	zoneVals := z.GetValidatorsAsSlice()
	for _, coin := range coins {
		for _, v := range zoneVals {
			// if token share, add amount to
			if strings.HasPrefix(coin.Denom, v) {
				val, ok := out[v]
				if !ok {
					val = &ValidatorIntent{ValoperAddress: v, Weight: sdk.ZeroDec()}
				}
				val.Weight = val.Weight.Add(sdk.NewDecFromInt(coin.Amount))
				out[v] = val
			}
		}
	}

	return out
}

func (z RegisteredZone) GetValidatorsAsSlice() []string {
	l := make([]string, 0)
	for _, v := range z.Validators {
		l = append(l, v.ValoperAddress)
	}
	return l
}

func (z RegisteredZone) DetermineStateIntentDiff(intents []DelegatorIntent) map[string]sdk.Dec {
	aggregateIntent := make(map[string]sdk.Dec)
	totalAggregateIntent := sdk.ZeroDec()
	currentState := make(map[string]sdk.Dec)
	totalDelegations := sdk.ZeroDec()
	diff := make(map[string]sdk.Dec)

	for _, intent := range intents {
		for _, vIntent := range intent.Intents {
			vStake, found := aggregateIntent[vIntent.ValoperAddress]
			if !found {
				vStake = sdk.ZeroDec()
			}
			vStake = vStake.Add(vIntent.Weight)
			aggregateIntent[vIntent.ValoperAddress] = vStake
		}
	}

	// sum total aggregate intent
	for _, val := range aggregateIntent {
		totalAggregateIntent = totalAggregateIntent.Add(val)
	}

	for _, i := range z.Validators {
		stake := sdk.ZeroDec()
		for _, delegation := range i.GetDelegations() {
			stake = stake.Add(delegation.Amount)
		}
		currentState[i.ValoperAddress] = stake
		totalDelegations = totalDelegations.Add(stake)
	}

	ratio := totalDelegations.Quo(totalAggregateIntent) // will always be >= 1.0

	for _, i := range z.Validators {
		current, found := currentState[i.ValoperAddress]
		if !found {
			panic("this shouldn't happen...")
		}
		desired, found := aggregateIntent[i.ValoperAddress]
		if !found {
			desired = sdk.ZeroDec() // this is okay! just means nobody wants this validator anymore!
		}
		thisDiff := desired.Mul(ratio).Sub(current)
		if !thisDiff.Equal(sdk.ZeroDec()) {
			diff[i.ValoperAddress] = thisDiff
		}
	}

	return diff
}
