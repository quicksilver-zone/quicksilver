package types

import (
	fmt "fmt"
	math "math"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (z RegisteredZone) GetDelegationAccountsByLowestBalance(qty uint64) []*ICAAccount {
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

func (z *RegisteredZone) GetValidatorByValoper(valoper string) (*Validator, error) {
	for _, v := range z.Validators {
		if v.ValoperAddress == valoper {
			return v, nil
		}
	}
	return nil, fmt.Errorf("invalid validator %s", valoper)
}

func (z *RegisteredZone) GetDelegationAccountByAddress(address string) (*ICAAccount, error) {
	if z.DelegationAddresses == nil {
		return nil, fmt.Errorf("no delegation accounts set: %v", z)
	}
	for _, account := range z.DelegationAddresses {
		if account.GetAddress() == address {
			return account, nil
		}
	}
	return nil, fmt.Errorf("unable to find delegation account: %s", address)
}

func (z *RegisteredZone) GetDelegationsForDelegator(delegator string) []*Delegation {
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

func (z RegisteredZone) DetermineStateIntentDiff(aggregateIntent map[string]*ValidatorIntent) map[string]sdk.Int {
	totalAggregateIntent := sdk.ZeroDec()
	currentState := make(map[string]sdk.Int)
	totalDelegations := sdk.ZeroInt()
	diff := make(map[string]sdk.Int)

	// sum total aggregate intent
	for _, val := range aggregateIntent {
		totalAggregateIntent = totalAggregateIntent.Add(val.Weight)

	}

	if totalAggregateIntent.IsZero() {
		// if totalAggregateIntent is zero (that is, we have no intent set - which can happen
		// if we have only ever have native tokens staked and nbody has signalled intent) give
		// every validator an equal intent artificially.

		// this can be removed when we cache intent.
		if aggregateIntent == nil {
			aggregateIntent = make(map[string]*ValidatorIntent)
		}

		for _, val := range z.Validators {
			aggregateIntent[val.ValoperAddress] = &ValidatorIntent{ValoperAddress: val.ValoperAddress, Weight: sdk.OneDec()}
			totalAggregateIntent = totalAggregateIntent.Add(sdk.OneDec())
		}
	}

	for _, i := range z.Validators {
		stake := sdk.ZeroInt()
		for _, delegation := range i.GetDelegations() {
			stake = stake.Add(delegation.Amount.TruncateInt())
		}
		currentState[i.ValoperAddress] = stake
		totalDelegations = totalDelegations.Add(stake)
	}
	ratio := totalDelegations.ToDec().Quo(totalAggregateIntent) // will always be >= 1.0

	for _, i := range z.Validators {
		current, found := currentState[i.ValoperAddress]
		if !found {
			// this probably can happen if we have intent for a validator not in the set
			// (although we _should_ have all validators, current and past in the set).
			panic("this shouldn't happen...")
		}
		desired, found := aggregateIntent[i.ValoperAddress]
		if !found {
			desired = &ValidatorIntent{ValoperAddress: i.ValoperAddress, Weight: sdk.ZeroDec()} // this is okay! just means nobody wants this validator anymore!
		}
		thisDiff := desired.Weight.Mul(ratio).TruncateInt().Sub(current)
		if !thisDiff.Equal(sdk.ZeroInt()) {
			diff[i.ValoperAddress] = thisDiff
		}
	}
	return diff
}

func (z RegisteredZone) ApplyDiffsToDistribution(distribution map[string]sdk.Coin, diffs map[string]sdk.Int) (map[string]sdk.Coin, sdk.Int) {
	remaining := sdk.ZeroInt()
	// sort map to ordered slice
	for _, val := range sortMapToSlice(diffs) {
		thisAmount, ok := distribution[val.str]
		if !ok {
			// no allocation to this val from intents, so skip.
			// TODO: should we _add_ a new distribution here? We could easily, we just need to know the denom.
			continue
		}

		if val.i.GT(sdk.ZeroInt()) {
			if thisAmount.Amount.LTE(val.i) { // if the new additional value is LTE the positive diff, remove it all and assign all values to remaining.
				delete(distribution, val.str)
				remaining = remaining.Add(thisAmount.Amount)
			} else { // GT
				distribution[val.str] = distribution[val.str].SubAmount(val.i)
				remaining = remaining.Add(val.i)
			}
		} else {
			// increase new amounts by diff from remaining
			if val.i.Abs().GTE(remaining) {
				distribution[val.str] = distribution[val.str].AddAmount(remaining) // negative addition :(
				remaining = sdk.ZeroInt()
				break
			} else {
				distribution[val.str] = distribution[val.str].SubAmount(val.i) // negative addition :(
				remaining = remaining.Add(val.i)
			}
		}
	}

	return distribution, remaining
}

type sortableStringInt struct {
	str string
	i   sdk.Int
}

func sortMapToSlice(numbers map[string]sdk.Int) []sortableStringInt {
	out := []sortableStringInt{}
	for str, int := range numbers {
		if !int.IsZero() {
			out = append(out, sortableStringInt{str: str, i: int})
		}
	}
	// sort
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].i.GT(out[j].i)
	})
	return out
}

func (z *RegisteredZone) GetRedemptionTargets(requests map[string]sdk.Int, denom string) map[string]map[string]sdk.Coin {
	out := make(map[string]map[string]sdk.Coin)
	for valoper, tokens := range requests {

		remainingTokens := tokens
		// TODO: order delegations from highest to lowest, as a reference. We wish to even these out as much as possible.
		// return a map of delegation bucket deviation from median.

		validator, err := z.GetValidatorByValoper(valoper)
		if err != nil {
			continue
		}

		for _, i := range validator.Delegations {
			if i.Amount.TruncateInt().GTE(remainingTokens) {
				if out[i.DelegationAddress] == nil {
					out[i.DelegationAddress] = make(map[string]sdk.Coin)
				}
				out[i.DelegationAddress][i.ValidatorAddress] = sdk.NewCoin(denom, remainingTokens)
				break
			} else {
				val := i.Amount.TruncateInt()
				remainingTokens = remainingTokens.Sub(val)
				if out[i.DelegationAddress] == nil {
					out[i.DelegationAddress] = make(map[string]sdk.Coin)
				}
				out[i.DelegationAddress][i.ValidatorAddress] = sdk.NewCoin(denom, val)
			}
		}

	}
	return out
}

// func (z *RegisteredZone) UpdateDelegatedAmount() {

// 	sum := map[string]sdk.Dec{}
// 	for _, validator := range z.Validators {
// 		for _, delegation := range validator.Delegations {
// 			_, ok := sum[delegation.DelegationAddress]
// 			if !ok {
// 				sum[delegation.DelegationAddress] = delegation.Amount
// 			} else {
// 				sum[delegation.DelegationAddress] = sum[delegation.DelegationAddress].Add(delegation.Amount)
// 			}
// 		}
// 	}

// 	out := sdk.NewCoin(z.BaseDenom, sdk.ZeroInt())
// 	for _, da := range z.DelegationAddresses {
// 		val, ok := sum[da.Address]
// 		if ok {
// 			delCoin := sdk.NewCoin(z.BaseDenom, val.TruncateInt())
// 			if da.DelegatedBalance.IsNil() || da.DelegatedBalance.IsZero() || !da.DelegatedBalance.Equal(delCoin) {
// 				da.DelegatedBalance = delCoin
// 			}
// 			out = out.Add(da.DelegatedBalance)
// 		}
// 	}
// }

func (z *RegisteredZone) GetDelegatedAmount() sdk.Coin {
	out := sdk.NewCoin(z.BaseDenom, sdk.ZeroInt())
	for _, da := range z.DelegationAddresses {
		out = out.Add(da.DelegatedBalance)
	}
	return out
}
