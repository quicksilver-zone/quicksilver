package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (v Validator) GetAddressBytes() ([]byte, error) {
	_, addr, err := bech32.DecodeAndConvert(v.ValoperAddress)
	return addr, err
}

// check to see if two validator instances are equal. Used in testing.
// func (v Validator) IsEqual(other Validator) bool {
//	if v.ValoperAddress != other.ValoperAddress {
//		return false
//	}
//
//	if !v.CommissionRate.Equal(other.CommissionRate) {
//		return false
//	}
//
//	if !v.DelegatorShares.Equal(other.DelegatorShares) {
//		return false
//	}
//
//	if !v.VotingPower.Equal(other.VotingPower) {
//		return false
//	}
//	return true
// }

func (v Validator) SharesToTokens(shares sdk.Dec) sdkmath.Int {
	if v.DelegatorShares.IsNil() || v.DelegatorShares.IsZero() {
		return sdk.ZeroInt()
	}

	return sdk.NewDecFromInt(v.VotingPower).Quo(v.DelegatorShares).Mul(shares).TruncateInt()
}

func (di DelegatorIntent) AddOrdinal(multiplier sdk.Dec, intents ValidatorIntents) DelegatorIntent {
	if len(intents) == 0 {
		return di
	}

	if len(di.Intents) == 0 {
		di.Intents = make(ValidatorIntents, 0)
	}

	ordinalized := di.Ordinalize(multiplier)

OUTER:
	for _, i := range intents.Sort() {
		for jdx, j := range ordinalized.SortedIntents() {
			if i.ValoperAddress == j.ValoperAddress {
				ordinalized.Intents[jdx].Weight = j.Weight.Add(i.Weight)
				continue OUTER
			}
		}
		ordinalized.Intents = append(ordinalized.Intents, i)
	}

	// we may have appended above, so resort intents.
	ordinalized.SortedIntents()

	return ordinalized.Normalize()
}

func (di DelegatorIntent) IntentForValoper(valoper string) (*ValidatorIntent, bool) {
	for _, intent := range di.Intents {
		if intent.ValoperAddress == valoper {
			return intent, true
		}
	}
	return nil, false
}

func (di DelegatorIntent) MustIntentForValoper(valoper string) *ValidatorIntent {
	intent, found := di.IntentForValoper(valoper)
	if !found {
		panic("intent not found")
	}
	return intent
}

func (di DelegatorIntent) Normalize() DelegatorIntent {
	summedWeight := sdk.ZeroDec()
	// cached sorted intents as we don't modify in the first iteration.
	sortedIntents := di.SortedIntents()
	for _, i := range sortedIntents {
		summedWeight = summedWeight.Add(i.Weight)
	}

	// zero summed weight, we should panic here, something is very wrong...
	if summedWeight.IsZero() {
		return di
	}

	for idx, i := range sortedIntents {
		di.Intents[idx].Weight = i.Weight.QuoTruncate(summedWeight)
	}

	return di
}

func (di DelegatorIntent) Ordinalize(multiple sdk.Dec) DelegatorIntent {
	for idx, i := range di.SortedIntents() {
		di.Intents[idx].Weight = i.Weight.Mul(multiple)
	}

	return di
}

func (di *DelegatorIntent) SortedIntents() ValidatorIntents {
	di.Intents = di.Intents.Sort()
	return di.Intents
}
