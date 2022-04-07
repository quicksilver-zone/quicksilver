package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func AccAddressFromBech32(address string, check_hrp string) (addr sdk.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.AccAddress{}, errors.New("empty address string is not allowed")
	}

	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return nil, err
	}

	if check_hrp != "" {
		if check_hrp != hrp {
			return sdk.AccAddress{}, fmt.Errorf("unexpected hrp - got %s expected %s", hrp, check_hrp)
		}
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(bz), nil
}

func ValAddressFromBech32(address string, check_hrp string) (addr sdk.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.ValAddress{}, errors.New("empty address string is not allowed")
	}

	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return nil, err
	}

	if check_hrp != "" {
		if check_hrp != hrp {
			return sdk.ValAddress{}, fmt.Errorf("unexpected hrp - got %s expected %s", hrp, check_hrp)
		}
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.ValAddress(bz), nil
}
