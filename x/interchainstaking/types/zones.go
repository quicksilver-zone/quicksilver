package types

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (z Zone) SupportReturnToSender() bool { return z.ReturnToSender }
func (z Zone) IsUnbondingEnabled() bool    { return z.UnbondingEnabled }
func (z Zone) SupportLsm() bool            { return z.LiquidityModule }

func (z *Zone) GetValoperPrefix() string {
	if z != nil {
		return z.AccountPrefix + "valoper"
	}
	return ""
}

func (z Zone) IsDelegateAddress(addr string) bool {
	return z.DelegationAddress != nil && z.DelegationAddress.Address == addr
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

// memo functionality

// this method exist to make testing easier!
func (z *Zone) UpdateIntentWithCoins(intent DelegatorIntent, multiplier sdk.Dec, inAmount sdk.Coins, vals []string) DelegatorIntent {
	// coinIntent is ordinal
	intent = intent.AddOrdinal(multiplier, z.ConvertCoinsToOrdinalIntents(inAmount, vals))
	return intent
}

type ZoneMemoUpdate func() error

func (z *Zone) UpdateZoneWithMemo(intent DelegatorIntent, memo string, multiplier sdk.Dec, inAmount sdk.Coins) (DelegatorIntent, []ZoneMemoUpdate, error) {
	memoIntent, memoFields, err := z.DecodeMemo(inAmount, memo)
	if err != nil {
		return DelegatorIntent{}, nil, fmt.Errorf("error decoding memo: %w", err)
	}

	intent = intent.AddOrdinal(multiplier, memoIntent)

	updates, err := z.memoUpdatesFromFields(memoFields)
	if err != nil {
		return DelegatorIntent{}, nil, fmt.Errorf("error generating memo field updates: %w", err)
	}

	return intent, updates, nil
}

func (z *Zone) memoUpdatesFromFields(memoFields MemoFields) ([]ZoneMemoUpdate, error) {
	updates := make([]ZoneMemoUpdate, 0)

	_ = memoFields
	// TODO ADD LOGIC FOR UPDATING BASED ON FIELDS

	return updates, nil
}

func (z *Zone) ConvertCoinsToOrdinalIntents(coins sdk.Coins, zoneVals []string) ValidatorIntents {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents, 0)
COINS:
	for _, coin := range coins {
		for _, v := range zoneVals {
			// if token share, add amount to
			if !strings.HasPrefix(coin.Denom, v) {
				continue
			}
			val, ok := out.GetForValoper(v)
			if !ok {
				val = &ValidatorIntent{ValoperAddress: v, Weight: sdk.ZeroDec()}
			}
			val.Weight = val.Weight.Add(sdk.NewDecFromInt(coin.Amount))
			out = out.SetForValoper(v, val)
			continue COINS
		}
	}

	return out
}

func (z *Zone) ConvertMemoToOrdinalIntents(coins sdk.Coins, memo string) (ValidatorIntents, error) {
	// should we be return DelegatorIntent here?

	validatorIntents, _, err := z.DecodeMemo(coins, memo)
	if err != nil {
		return ValidatorIntents{}, fmt.Errorf("error decoding memo: %w", err)
	}

	return validatorIntents, nil
}

// decode memo
// if zone.Is_118:
//		decode as we have ( up to 8 validators)
// 		return
//
// decode as we have (up to 6 validators)
// look for separator 0xFF
// field_id = [
//  0x00 = map
//  0x01 = rts
// ] // will scale to future fields

var separator = []byte{byte(255)}

type fieldType int

const (
	AccountMap fieldType = iota
	ReturnToSender
	// add more here.
)

type MemoField struct {
	ID   int
	Data []byte
}

type MemoFields []MemoField

func (z *Zone) DecodeMemo(coins sdk.Coins, memo string) (validatorIntents ValidatorIntents, memoFields MemoFields, err error) {
	if memo == "" {
		return validatorIntents, memoFields, errors.New("memo length unexpectedly zero")
	}

	memoBytes, err := base64.StdEncoding.DecodeString(memo)
	if err != nil {
		return validatorIntents, memoFields, fmt.Errorf("failed to decode base64 message: %w", err)
	}

	parts := bytes.Split(memoBytes, separator)
	valWeightsBytes := parts[0]
	if len(valWeightsBytes)%21 != 0 { // memo must be one byte (1-200) weight then 20 byte valoperAddress
		return validatorIntents, memoFields, fmt.Errorf("unable to determine intent from memo: Message was incorrect length: %d", len(memoBytes))
	}

	switch {
	case len(parts) == 0:
		return validatorIntents, memoFields, errors.New("invalid memo format")

	case len(parts) == 1:
		if len(valWeightsBytes)/21 > 8 {
			return validatorIntents, memoFields, errors.New("memo format not currently supported")
		}

	default:
		// iterate through all non-validator weights parts of the memo
		memoFields, err = ParseMemoFields(parts[1])
		if err != nil {
			return validatorIntents, memoFields, fmt.Errorf("unable to decode memo field: %w", err)
		}
	}

	validatorIntents, err = z.validatorIntentsFromBytes(coins, valWeightsBytes)

	return validatorIntents, memoFields, err
}

func (z *Zone) validatorIntentsFromBytes(coins sdk.Coins, weightBytes []byte) (ValidatorIntents, error) {
	validatorIntents := make(ValidatorIntents, 0)

	for index := 0; index < len(weightBytes); {
		// truncate weight to 200
		rawWeight := int64(weightBytes[index])
		if rawWeight > 200 {
			return validatorIntents, fmt.Errorf("out of bounds value received in memo intent message; expected 0-200, got %d", rawWeight)
		}
		sdkWeight := sdk.NewDecFromInt(sdk.NewInt(rawWeight)).QuoInt(sdk.NewInt(200))
		coinWeight := sdkWeight.MulInt(coins.AmountOf(z.BaseDenom))
		index++
		address := weightBytes[index : index+20]
		index += 20
		valAddr, err := bech32.ConvertAndEncode(z.AccountPrefix+"valoper", address)
		if err != nil {
			return validatorIntents, err
		}
		val, ok := validatorIntents.GetForValoper(valAddr)
		if !ok {
			val = &ValidatorIntent{ValoperAddress: valAddr, Weight: sdk.ZeroDec()}
		}
		val.Weight = val.Weight.Add(coinWeight)
		validatorIntents = validatorIntents.SetForValoper(valAddr, val)
	}

	return validatorIntents, nil
}

func ParseMemoFields(fieldBytes []byte) (MemoFields, error) {
	if len(fieldBytes) < 3 {
		return MemoFields{}, errors.New("invalid field bytes format")
	}

	memoFields := make([]MemoField, 0)

	idx := 0
	for idx < len(fieldBytes) {
		// prevent out of bounds
		if len(fieldBytes[idx:]) < 3 {
			return MemoFields{}, errors.New("invalid field bytes format")
		}

		fieldID := int(fieldBytes[idx])
		idx++
		fieldLength := int(fieldBytes[idx])
		idx++

		if len(fieldBytes[idx:]) < fieldLength {
			return MemoFields{}, errors.New("invalid field length for memo field")
		}

		memoFields = append(memoFields, MemoField{
			ID:   fieldID,
			Data: fieldBytes[idx : idx+fieldLength], // add one for half-open range
		})
		idx += fieldLength
	}

	// secondary sanity check
	if idx != len(fieldBytes) {
		return MemoFields{}, errors.New("error parsing multiple fields")
	}

	return memoFields, nil
}
