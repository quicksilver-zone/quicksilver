package types

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (z Zone) SupportMultiSend() bool { return z.MultiSend }
func (z Zone) SupportLsm() bool       { return z.LiquidityModule }

func (z Zone) IsDelegateAddress(addr string) bool {
	return z.DelegationAddress.Address == addr
}

func (z *Zone) GetDelegationAccount() (*ICAAccount, error) {
	if z.DelegationAddress == nil {
		return nil, fmt.Errorf("no delegation account set: %v", z)
	}
	return z.DelegationAddress, nil
}

func (z *Zone) ValidateCoinsForZone(coins sdk.Coins, zoneVals []string) error {

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
func (z *Zone) UpdateIntentWithCoins(intent DelegatorIntent, multiplier sdk.Dec, inAmount sdk.Coins, vals []string) DelegatorIntent {
	// coinIntent is ordinal
	intent = intent.AddOrdinal(multiplier, z.ConvertCoinsToOrdinalIntents(inAmount, vals))
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

func (z *Zone) ConvertCoinsToOrdinalIntents(coins sdk.Coins, zoneVals []string) ValidatorIntents {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents, 0)
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
	panic("deprecated")
}

func (z Zone) GetValidatorsAddressesAsSlice() []string {
	panic("deprecated")
}
