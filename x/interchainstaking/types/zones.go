package types

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (z Zone) SupportMultiSend() bool { return z.MultiSend }
func (z Zone) SupportLsm() bool       { return z.LiquidityModule }

func (z Zone) IsDelegateAddress(addr string) bool {
	return z.DelegationAddress.Address == addr
}

func (z *Zone) GetValidatorByValoper(valoper string) (*Validator, bool) {
	for _, v := range z.GetValidatorsSorted() {
		if v.ValoperAddress == valoper {
			return v, true
		}
	}
	return nil, false
}

func (z *Zone) GetDelegationAccount() (*ICAAccount, error) {
	if z.DelegationAddress == nil {
		return nil, fmt.Errorf("no delegation account set: %v", z)
	}
	return z.DelegationAddress, nil
}

func (z *Zone) ValidateCoinsForZone(ctx sdk.Context, coins sdk.Coins) error {
	zoneVals := z.GetValidatorsAddressesAsSlice()

COINS:
	for _, coin := range coins.Sort() {
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

// this method exist to make testing easier!
func (z *Zone) UpdateIntentWithCoins(intent DelegatorIntent, multiplier sdk.Dec, inAmount sdk.Coins) DelegatorIntent {
	// coinIntent is ordinal
	intent = intent.AddOrdinal(multiplier, z.ConvertCoinsToOrdinalIntents(inAmount))
	return intent
}

// this method exist to make testing easier!
func (z *Zone) UpdateIntentWithMemo(intent DelegatorIntent, memo string, multiplier sdk.Dec, inAmount sdk.Coins) (DelegatorIntent, error) {
	// coinIntent is ordinal
	memoIntent, err := z.ConvertMemoToOrdinalIntents(inAmount, memo)
	if err != nil {
		return DelegatorIntent{}, err
	}
	intent = intent.AddOrdinal(multiplier, memoIntent)
	return intent, nil
}

func (z *Zone) ConvertCoinsToOrdinalIntents(coins sdk.Coins) ValidatorIntents {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents, 0)
	zoneVals := z.GetValidatorsAddressesAsSlice()
COINS:
	for _, coin := range coins {
		for _, v := range zoneVals {
			// if token share, add amount to
			if strings.HasPrefix(coin.Denom, v) {
				val, ok := out.GetForValoper(v)
				if !ok {
					val = &ValidatorIntent{ValoperAddress: v, Weight: sdk.ZeroDec()}
				}
				val.Weight = val.Weight.Add(sdk.NewDecFromInt(coin.Amount))
				out = out.SetForValoper(v, val)
				continue COINS
			}
		}
	}

	return out
}

func (z *Zone) ConvertMemoToOrdinalIntents(coins sdk.Coins, memo string) (ValidatorIntents, error) {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents, 0)

	if len(memo) == 0 {
		return out, errors.New("memo length unexpectedly zero")
	}

	memoBytes, err := base64.StdEncoding.DecodeString(memo)
	if err != nil {
		return out, fmt.Errorf("unable to determine intent from memo: Failed to decode base64 message: %s", err.Error())
	}

	if len(memoBytes)%21 != 0 { // memo must be one byte (1-200) weight then 20 byte valoperAddress
		return out, fmt.Errorf("unable to determine intent from memo: Message was incorrect length: %d", len(memoBytes))
	}

	for index := 0; index < len(memoBytes); {
		// truncate weight to 200
		rawWeight := int64(memoBytes[index])
		if rawWeight > 200 {
			return ValidatorIntents{}, fmt.Errorf("out of bounds value received in memo intent message; expected 0-200, got %d", rawWeight)
		}
		sdkWeight := sdk.NewDecFromInt(sdk.NewInt(rawWeight)).QuoInt(sdk.NewInt(200))
		coinWeight := sdkWeight.MulInt(coins.AmountOf(z.BaseDenom))
		index++
		address := memoBytes[index : index+20]
		index += 20
		valAddr, err := bech32.ConvertAndEncode(z.AccountPrefix+"valoper", address)
		if err != nil {
			return ValidatorIntents{}, err
		}
		val, ok := out.GetForValoper(valAddr)
		if !ok {
			val = &ValidatorIntent{ValoperAddress: valAddr, Weight: sdk.ZeroDec()}
		}
		val.Weight = val.Weight.Add(coinWeight)
		out = out.SetForValoper(valAddr, val)
	}
	return out, nil
}

func (z *Zone) GetValidatorsSorted() []*Validator {
	sort.SliceStable(z.Validators, func(i, j int) bool {
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

func (z Zone) GetBondedValidatorAddressesAsSlice() []string {
	l := make([]string, 0)
	for _, v := range z.Validators {
		if v.Status == "BOND_STATUS_BONDED" {
			l = append(l, v.ValoperAddress)
		}
	}

	sort.Strings(l)

	return l
}

func (z *Zone) GetAggregateIntentOrDefault() ValidatorIntents {
	var intents ValidatorIntents
	var filteredIntents ValidatorIntents

	if len(z.AggregateIntent) == 0 {
		intents = z.DefaultAggregateIntents()
	} else {
		intents = z.AggregateIntent
	}
	// filter intents here...
	// check validators for tombstoned
	for _, v := range intents {
		val, found := z.GetValidatorByValoper(v.ValoperAddress)
		// this case should not happen as we check the validity of a validator entry when intent is set.
		if !found {
			continue
		}
		// we should never let tombstoned validators into the list, even if they are explicitly selected
		if val.Tombstoned {
			continue
		}

		// we should never let denylist validators into the list, even if they are explicitly selected
		// if in deny list {
		// continue
		// }
		filteredIntents = append(filteredIntents, v)
	}

	return filteredIntents
}

// defaultAggregateIntents determines the default aggregate intent (for epoch 0)
func (z *Zone) DefaultAggregateIntents() ValidatorIntents {
	out := make(ValidatorIntents, 0)
	for _, val := range z.GetValidatorsSorted() {
		if val.CommissionRate.LTE(sdk.NewDecWithPrec(5, 1)) { // 50%; make this a param.
			if !val.Jailed && !val.Tombstoned && val.Status == stakingtypes.BondStatusBonded {
				out = append(out, &ValidatorIntent{ValoperAddress: val.GetValoperAddress(), Weight: sdk.OneDec()})
			}
		}
	}

	valCount := sdk.NewInt(int64(len(out)))

	// normalise the array (divide everything by length of intent list)
	for idx, intent := range out.Sort() {
		out[idx].Weight = intent.Weight.Quo(sdk.NewDecFromInt(valCount))
	}

	return out
}
