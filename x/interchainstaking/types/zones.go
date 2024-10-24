package types

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
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

func (z Zone) IsDepositAddress(addr string) bool {
	return z.DepositAddress != nil && z.DepositAddress.Address == addr
}

func (z Zone) IsDelegateAddress(addr string) bool {
	return z.DelegationAddress != nil && z.DelegationAddress.Address == addr
}

func (z Zone) IsWithdrawalAddress(addr string) bool {
	return z.WithdrawalAddress != nil && z.WithdrawalAddress.Address == addr
}

func (z *Zone) GetDelegationAccount() (*ICAAccount, error) {
	if z.DelegationAddress == nil {
		return nil, fmt.Errorf("no delegation account set: %v", z)
	}
	return z.DelegationAddress, nil
}

func (z *Zone) SetWithdrawalWaitgroup(logger log.Logger, num uint32, reason string) {
	logger.Info("setting withdrawal waitgroup", "zone", z.ChainId, "existing", z.WithdrawalWaitgroup, "new", num, "reason", reason)
	z.WithdrawalWaitgroup = num
}

func (z *Zone) DecrementWithdrawalWaitgroup(logger log.Logger, num uint32, reason string) error {
	if z.WithdrawalWaitgroup-num > z.WithdrawalWaitgroup { // underflow
		logger.Error("error decrementing withdrawal waitgroup: uint32 underflow ", "zone", z.ChainId, "existing", z.WithdrawalWaitgroup, "decrement", num, "reason", reason)
		return errors.New("unable to decrement the withdrawal waitgroup below 0")
	}
	logger.Info("decrementing withdrawal waitgroup", "zone", z.ChainId, "existing", z.WithdrawalWaitgroup, "decrement", num, "new", z.WithdrawalWaitgroup-num, "reason", reason)
	z.WithdrawalWaitgroup -= num
	return nil
}

func (z *Zone) IncrementWithdrawalWaitgroup(logger log.Logger, num uint32, reason string) error {
	if z.WithdrawalWaitgroup+num < z.WithdrawalWaitgroup { // overflow
		logger.Error("error incrementing withdrawal waitgroup: uint32 overflow ", "zone", z.ChainId, "existing", z.WithdrawalWaitgroup, "increment", num, "reason", reason)
		return errors.New("unable to increment the withdrawal waitgroup above 4294967295")
	}
	logger.Info("incrementing withdrawal waitgroup", "zone", z.ChainId, "existing", z.WithdrawalWaitgroup, "increment", num, "new", z.WithdrawalWaitgroup+num, "reason", reason)
	z.WithdrawalWaitgroup += num
	return nil
}

// ValidateCoinsForZone checks whether an inbound denomination is valid for this zone.
// valid coins comprise:
//   - non-lsm: staking denom only (z.BaseDenom)
//   - lsm: staking denom, and tokens prefixed with zone.AccountPrefix
//
// lsm: if valoper is not in zoneVals map, then we return true, false, and this tx will be revisited.
func (z *Zone) ValidateCoinsForZone(coins sdk.Coins, zoneVals map[string]bool) (valid bool, matchesVal bool) {
	for _, coin := range coins.Sort() {
		if coin.Denom == z.BaseDenom {
			continue
		}

		// if liquidity module enabled, check to see if this is a tokenized share.
		if z.LiquidityModule {
			coinParts := strings.Split(coin.Denom, "/")
			if len(coinParts) != 2 || !strings.HasPrefix(coinParts[0], z.AccountPrefix) {
				return false, false
			}

			if _, ok := zoneVals[coinParts[0]]; !ok {
				return true, false
			}

			continue
		}

		return false, false
	}
	return true, true
}

// memo functionality

// this method exist to make testing easier!
func (z *Zone) UpdateIntentWithCoins(intent DelegatorIntent, multiplier sdk.Dec, inAmount sdk.Coins, vals map[string]bool) DelegatorIntent {
	// coinIntent is ordinal
	return intent.AddOrdinal(multiplier, z.ConvertCoinsToOrdinalIntents(inAmount, vals))
}

func (*Zone) UpdateZoneIntentWithMemo(memoIntent ValidatorIntents, intent DelegatorIntent, multiplier sdk.Dec) DelegatorIntent {
	return intent.AddOrdinal(multiplier, memoIntent)
}

func (*Zone) ConvertCoinsToOrdinalIntents(coins sdk.Coins, zoneVals map[string]bool) ValidatorIntents {
	out := make(ValidatorIntents, 0, len(coins))
	for _, coin := range coins {
		coinParts := strings.Split(coin.Denom, "/")
		if len(coinParts) != 2 {
			continue
		}

		if _, ok := zoneVals[coinParts[0]]; !ok {
			continue
		}

		val, ok := out.GetForValoper(coinParts[0])
		if !ok {
			val = &ValidatorIntent{ValoperAddress: coinParts[0], Weight: sdk.ZeroDec()}
		}
		val.Weight = val.Weight.Add(sdk.NewDecFromInt(coin.Amount))
		out = out.SetForValoper(coinParts[0], val)
	}

	return out
}

const (
	FieldTypeAccountMap     int = 0x00
	FieldTypeReturnToSender int = 0x01
	FieldTypeIntent         int = 0x02
	FieldTypeAutoClaim      int = 0x03
)

type MemoField struct {
	ID   int
	Data []byte
}

type MemoFields map[int]MemoField

func (m MemoFields) RTS() bool {
	_, found := m[FieldTypeReturnToSender]
	return found
}

func (m MemoFields) AccountMap() ([]byte, bool) {
	field, found := m[FieldTypeAccountMap]
	return field.Data, found
}

func (m MemoFields) Intent(coins sdk.Coins, zone *Zone) (ValidatorIntents, bool) {
	field, found := m[FieldTypeIntent]
	if !found {
		return nil, false
	}

	validatorIntents, err := zone.validatorIntentsFromBytes(coins, field.Data)
	if err != nil {
		return validatorIntents, false
	}
	return validatorIntents, true
}

func (m MemoFields) AutoClaim() bool {
	_, found := m[FieldTypeAutoClaim]
	return found
}

func (m *MemoField) Validate() error {
	switch m.ID {
	case FieldTypeAccountMap:
		if len(m.Data) == 0 {
			return errors.New("invalid length for account map memo field 0")
		}
		// check if valid address
		_, err := sdk.Bech32ifyAddressBytes("test", m.Data)
		if err != nil {
			return fmt.Errorf("invalid address for account map memo field: address: %s", m.Data)
		}
	case FieldTypeReturnToSender:
		// do nothing - we ignore data if RTS
	case FieldTypeAutoClaim:
		// do nothing - we ignore data if AutoClaim
	case FieldTypeIntent:
		if len(m.Data)%21 != 0 { // memo must be one byte (1-200) weight then 20 byte valoperAddress
			return fmt.Errorf("invalid length for validator intent memo field %d", len(m.Data))
		}
	default:
		return fmt.Errorf("invalid field type %d", m.ID)
	}

	return nil
}

func (*Zone) DecodeMemo(memo string) (memoFields MemoFields, err error) {
	if memo == "" {
		return memoFields, nil
	}

	memoBytes, err := base64.StdEncoding.DecodeString(memo)
	if err != nil {
		return memoFields, fmt.Errorf("failed to decode base64 message: %w", err)
	}

	memoFields, err = ParseMemoFields(memoBytes)
	if err != nil {
		return memoFields, fmt.Errorf("unable to decode memo field: %w", err)
	}

	return memoFields, err
}

func (z *Zone) validatorIntentsFromBytes(coins sdk.Coins, weightBytes []byte) (validatorIntents ValidatorIntents, err error) {
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
		valAddr, err := addressutils.EncodeAddressToBech32(z.AccountPrefix+"valoper", sdk.ValAddress(address))
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
	if len(fieldBytes) < 2 {
		return MemoFields{}, errors.New("invalid field bytes format")
	}

	memoFields := make(MemoFields)

	idx := 0
	for idx < len(fieldBytes) {
		// prevent out of bounds
		if len(fieldBytes[idx:]) < 2 {
			return memoFields, errors.New("invalid field bytes format")
		}

		fieldID := int(fieldBytes[idx])
		idx++
		fieldLength := int(fieldBytes[idx])
		idx++

		var data []byte
		switch {
		case fieldLength == 0:
			data = nil
		case len(fieldBytes[idx:]) < fieldLength:
			return memoFields, errors.New("invalid field length for memo field")
		default:
			data = fieldBytes[idx : idx+fieldLength]
		}

		memoField := MemoField{
			ID:   fieldID,
			Data: data,
		}
		err := memoField.Validate()
		if err != nil {
			return memoFields, fmt.Errorf("invalid memo field: %w", err)
		}

		if _, found := memoFields[fieldID]; found {
			return memoFields, fmt.Errorf("duplicate field ID found in memo: fieldID: %d", fieldID)
		}

		memoFields[fieldID] = memoField

		idx += fieldLength
	}

	// secondary sanity check
	if idx != len(fieldBytes) {
		return memoFields, errors.New("error parsing multiple fields")
	}

	return memoFields, nil
}
