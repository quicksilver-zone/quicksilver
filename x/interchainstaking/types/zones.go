package types

import (
	"encoding/hex"
	fmt "fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ingenuity-build/quicksilver/utils"
)

func (z RegisteredZone) GetDelegationAccountsByLowestBalance(qty uint64) []*ICAAccount {
	delegationAccounts := z.GetDelegationAccounts()
	sort.SliceStable(delegationAccounts, func(i, j int) bool {
		return delegationAccounts[i].DelegatedBalance.Amount.GT(delegationAccounts[j].DelegatedBalance.Amount)
	})
	if qty > 0 {
		return delegationAccounts[:int(utils.MinU64(append([]uint64{}, uint64(len(delegationAccounts)), qty)))]
	}
	return delegationAccounts
}

func (z RegisteredZone) SupportMultiSend() bool { return z.MultiSend }

func (z *RegisteredZone) GetValidatorByValoper(valoper string) (*Validator, error) {
	for _, v := range z.GetValidatorsSorted() {
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
	for _, account := range z.GetDelegationAccounts() {
		if account.GetAddress() == address {
			return account, nil
		}
	}
	return nil, fmt.Errorf("unable to find delegation account: %s", address)
}

func (z *RegisteredZone) ValidateCoinsForZone(ctx sdk.Context, coins sdk.Coins) error {

	zoneVals := z.GetValidatorsAddressesAsSlice()
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
	zoneVals := z.GetValidatorsAddressesAsSlice()
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

func (z *RegisteredZone) ConvertMemoToOrdinalIntents(ctx sdk.Context, coins sdk.Coins, memo string) map[string]*ValidatorIntent {
	// should we be return DelegatorIntent here?
	out := make(map[string]*ValidatorIntent)

	if len(memo) == 0 {
		return out
	}

	memoBytes, err := hex.DecodeString(memo)
	if err != nil {
		return out
	}

	if len(memoBytes)%33 != 0 { // memo must be one byte (1-200) weight then 32 byte valoperAddress
		return out
	}

	for remaining := len(memoBytes); remaining > 0; {
		sdkWeight := sdk.NewDecFromInt(sdk.NewInt(int64(memoBytes[0]))).QuoInt(sdk.NewInt(2)).MulInt(coins.AmountOf(z.BaseDenom))
		address := memoBytes[1:33]
		valAddr, _ := bech32.ConvertAndEncode(z.AccountPrefix+"valoper", address)

		val, ok := out[valAddr]
		if !ok {
			val = &ValidatorIntent{ValoperAddress: valAddr, Weight: sdk.ZeroDec()}
		}
		val.Weight = val.Weight.Add(sdkWeight)
		out[valAddr] = val

		memoBytes = memoBytes[33:]
	}

	return out
}

func (z *RegisteredZone) GetValidatorsSorted() []*Validator {
	vals := z.Validators
	sort.Slice(vals, func(i, j int) bool {
		return vals[i].ValoperAddress < vals[j].ValoperAddress
	})
	return vals
}

func (z RegisteredZone) GetValidatorsAddressesAsSlice() []string {
	l := make([]string, 0)
	for _, v := range z.Validators {
		l = append(l, v.ValoperAddress)
	}

	sort.Strings(l)

	return l
}

// func (z RegisteredZone) ApplyDiffsToDistribution(distribution map[string]sdk.Coin, diffs map[string]sdk.Int) (map[string]sdk.Coin, sdk.Int) {
// 	remaining := sdk.ZeroInt()
// 	// sort map to ordered slice
// 	for _, val := range sortMapToSlice(diffs) {
// 		thisAmount, ok := distribution[val.str]
// 		if !ok {
// 			// no allocation to this val from intents, so skip.
// 			// TODO: should we _add_ a new distribution here? We could easily, we just need to know the denom.
// 			continue
// 		}

// 		if val.i.GT(sdk.ZeroInt()) {
// 			if thisAmount.Amount.LTE(val.i) { // if the new additional value is LTE the positive diff, remove it all and assign all values to remaining.
// 				delete(distribution, val.str)
// 				remaining = remaining.Add(thisAmount.Amount)
// 			} else { // GT
// 				distribution[val.str] = distribution[val.str].SubAmount(val.i)
// 				remaining = remaining.Add(val.i)
// 			}
// 		} else {
// 			// increase new amounts by diff from remaining
// 			if val.i.Abs().GTE(remaining) {
// 				distribution[val.str] = distribution[val.str].AddAmount(remaining) // negative addition :(
// 				remaining = sdk.ZeroInt()
// 				break
// 			} else {
// 				distribution[val.str] = distribution[val.str].SubAmount(val.i) // negative addition :(
// 				remaining = remaining.Add(val.i)
// 			}
// 		}
// 	}

// 	return distribution, remaining
// }

// type sortableStringInt struct {
// 	str string
// 	i   sdk.Int
// }

// func sortMapToSlice(numbers map[string]sdk.Int) []sortableStringInt {
// 	out := []sortableStringInt{}
// 	for str, int := range numbers {
// 		if !int.IsZero() {
// 			out = append(out, sortableStringInt{str: str, i: int})
// 		}
// 	}
// 	// sort
// 	sort.SliceStable(out, func(i, j int) bool {
// 		return out[i].i.GT(out[j].i)
// 	})
// 	return out
// }

func (z *RegisteredZone) GetDelegatedAmount() sdk.Coin {
	out := sdk.NewCoin(z.BaseDenom, sdk.ZeroInt())
	for _, da := range z.GetDelegationAccounts() {
		out = out.Add(da.DelegatedBalance)
	}
	return out
}

func (z *RegisteredZone) GetDelegationAccounts() []*ICAAccount {
	delegationAccounts := z.DelegationAddresses
	sort.Slice(delegationAccounts, func(i, j int) bool {
		return delegationAccounts[i].Address < delegationAccounts[j].Address
	})
	return delegationAccounts
}
