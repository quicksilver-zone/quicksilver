package types

import (
	"encoding/base64"
	fmt "fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (z Zone) SupportMultiSend() bool { return z.MultiSend }
func (z Zone) SupportLsm() bool       { return z.LiquidityModule }

func (z Zone) IsDelegateAddress(addr string) bool {
	for _, acc := range z.DelegationAddresses {
		if acc.Address == addr {
			return true
		}
	}
	return false
}

func (z *Zone) GetValidatorByValoper(valoper string) (*Validator, error) {
	for _, v := range z.GetValidatorsSorted() {
		if v.ValoperAddress == valoper {
			return v, nil
		}
	}
	return nil, fmt.Errorf("invalid validator -> %s", valoper)
}

func (z *Zone) GetDelegationAccountByAddress(address string) (*ICAAccount, error) {
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

func (z *Zone) ValidateCoinsForZone(ctx sdk.Context, coins sdk.Coins) error {
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

func (z *Zone) ConvertCoinsToOrdinalIntents(coins sdk.Coins) ValidatorIntents {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents)
	zoneVals := z.GetValidatorsAddressesAsSlice()
COINS:
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
				continue COINS
			}
		}
	}

	return out
}

func (z *Zone) ConvertMemoToOrdinalIntents(coins sdk.Coins, memo string) ValidatorIntents {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents)

	if len(memo) == 0 {
		return out
	}

	memoBytes, err := base64.StdEncoding.DecodeString(memo)
	if err != nil {
		fmt.Println("unable to determine intent from memo: Failed to decode base64 message", err)
		return out
	}

	if len(memoBytes)%21 != 0 { // memo must be one byte (1-200) weight then 20 byte valoperAddress
		fmt.Println("unable to determine intent from memo: Message was incorrect length", len(memoBytes))
		return out
	}

	for index := 0; index < len(memoBytes); {
		sdkWeight := sdk.NewDecFromInt(sdk.NewInt(int64(memoBytes[index]))).QuoInt(sdk.NewInt(200))
		coinWeight := sdkWeight.MulInt(coins.AmountOf(z.BaseDenom))
		index++
		address := memoBytes[index : index+20]
		index += 20
		valAddr, _ := bech32.ConvertAndEncode(z.AccountPrefix+"valoper", address)

		val, ok := out[valAddr]
		if !ok {
			val = &ValidatorIntent{ValoperAddress: valAddr, Weight: sdk.ZeroDec()}
		}
		val.Weight = val.Weight.Add(coinWeight)
		out[valAddr] = val
	}
	return out
}

func (z *Zone) GetValidatorsSorted() []*Validator {
	sort.Slice(z.Validators, func(i, j int) bool {
		return z.Validators[i].ValoperAddress < z.Validators[j].ValoperAddress
	})
	return z.Validators
}

func (z Zone) GetValidatorsAddressesAsSlice() []string {
	l := make([]string, 0)
	for _, v := range z.Validators {
		l = append(l, v.ValoperAddress)
	}

	sort.Strings(l)

	return l
}

func (z *Zone) GetDelegatedAmount() sdk.Coin {
	out := sdk.NewCoin(z.BaseDenom, sdk.ZeroInt())
	for _, da := range z.DelegationAddresses {
		out = out.Add(da.DelegatedBalance)
	}
	return out
}

func (z *Zone) GetDelegationAccounts() []*ICAAccount {
	delegationAccounts := z.DelegationAddresses
	sort.Slice(delegationAccounts, func(i, j int) bool {
		return delegationAccounts[i].Address < delegationAccounts[j].Address
	})
	return delegationAccounts
}

func (z *Zone) GetAggregateIntentOrDefault() ValidatorIntents {
	if len(z.AggregateIntent) == 0 {
		return z.DefaultAggregateIntents()
	}
	return z.AggregateIntent
}

// defaultAggregateIntents determines the default aggregate intent (for epoch 0)
func (z *Zone) DefaultAggregateIntents() ValidatorIntents {
	out := make(ValidatorIntents)
	for _, val := range z.GetValidatorsSorted() {
		if val.CommissionRate.LTE(sdk.NewDecWithPrec(5, 1)) { // 50%; make this a param.
			out[val.GetValoperAddress()] = &ValidatorIntent{ValoperAddress: val.GetValoperAddress(), Weight: sdk.OneDec()}
		}
	}

	valCount := sdk.NewInt(int64(len(out)))

	// normalise the array (divide everything by length of intent list)
	for _, key := range out.Keys() {
		if val, ok := out[key]; ok {
			val.Weight = val.Weight.Quo(sdk.NewDecFromInt(valCount))
			out[key] = val
		}
	}

	return out
}
