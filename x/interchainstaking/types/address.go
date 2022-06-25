package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func AccAddressFromBech32(address string, checkHrp string) (addr sdk.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.AccAddress{}, errors.New("empty address string is not allowed")
	}

	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return nil, err
	}

	if checkHrp != "" {
		if checkHrp != hrp {
			return sdk.AccAddress{}, fmt.Errorf("unexpected hrp - got %s expected %s", hrp, checkHrp)
		}
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(bz), nil
}

func ValAddressFromBech32(address string, checkHrp string) (addr sdk.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.ValAddress{}, errors.New("empty address string is not allowed")
	}

	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return nil, err
	}

	if checkHrp != "" {
		if checkHrp != hrp {
			return sdk.ValAddress{}, fmt.Errorf("unexpected hrp - got %s expected %s", hrp, checkHrp)
		}
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.ValAddress(bz), nil
}
