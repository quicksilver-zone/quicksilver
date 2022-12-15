package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func AccAddressFromBech32(address string, checkHRP string) (addr sdk.AccAddress, err error) {
	if strings.TrimSpace(address) == "" {
		return sdk.AccAddress{}, errors.New("empty address string is not allowed")
	}

	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return nil, err
	}

	if checkHRP != "" {
		if checkHRP != hrp {
			return sdk.AccAddress{}, fmt.Errorf("unexpected hrp - got %s expected %s", hrp, checkHRP)
		}
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(bz), nil
}

func ValAddressFromBech32(address string, checkHRP string) (addr sdk.ValAddress, err error) {
	if strings.TrimSpace(address) == "" {
		return sdk.ValAddress{}, errors.New("empty address string is not allowed")
	}

	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return nil, err
	}

	if checkHRP != "" {
		if checkHRP != hrp {
			return sdk.ValAddress{}, fmt.Errorf("unexpected hrp - got %s expected %s", hrp, checkHRP)
		}
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.ValAddress(bz), nil
}

func GenerateAccAddressForTest() sdk.AccAddress {
	size := 32 // change the length of the generated random string here

	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		panic(err)
	}

	return sdk.AccAddress(rb)
}

func GenerateValAddressForTest() sdk.ValAddress {
	size := 32 // change the length of the generated random string here

	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		panic(err)
	}

	return sdk.ValAddress(rb)
}

func GenerateValAddressForTestWithPrefix(hrp string) string {
	addr, err := bech32.ConvertAndEncode(hrp, GenerateValAddressForTest())
	if err != nil {
		panic(err)
	}
	return addr
}

func GenerateAccAddressForTestWithPrefix(hrp string) string {
	addr, err := bech32.ConvertAndEncode(hrp, GenerateAccAddressForTest())
	if err != nil {
		panic(err)
	}
	return addr
}

func ConvertAccAddressForTestUsingPrefix(address sdk.AccAddress, prefix string) string {
	addr, err := bech32.ConvertAndEncode(prefix, address)
	if err != nil {
		panic(err)
	}
	return addr
}
