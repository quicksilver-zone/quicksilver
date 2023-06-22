package addressutils

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/ingenuity-build/quicksilver/utils/randomutils"
)

// AddressFromBech32 decodes a bech32 encoded address into a byte-slice, and validates the prefix (hrp).
// An empty prefix param skips the checking.
// Returns an error if address is zero-length, invalid or the prefix does not match.
func AddressFromBech32(address, prefix string) (addr []byte, err error) {
	if strings.TrimSpace(address) == "" {
		return nil, errors.New("empty address string is not allowed")
	}

	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return nil, err
	}

	if prefix != "" {
		if prefix != hrp {
			return nil, fmt.Errorf("unexpected prefix - got %s expected %s", hrp, prefix)
		}
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// AccAddressFromBech32 decodes a bech32 encoded address into an sdk.AccAddress, and validates the prefix (hrp).
// An empty prefix param skips the checking.
// Returns an error if address is zero-length, invalid or the prefix does not match.
func AccAddressFromBech32(address, prefix string) (sdk.AccAddress, error) {
	addr, err := AddressFromBech32(address, prefix)
	return sdk.AccAddress(addr), err
}

// MustAccAddressFromBech32 decodes a bech32 encoded address into an sdk.AccAddress, and validates the prefix (hrp).
// An empty prefix param skips the checking.
// Panics if address is zero-length, invalid or the prefix does not match.
func MustAccAddressFromBech32(address, prefix string) sdk.AccAddress {
	accAddress, err := AccAddressFromBech32(address, prefix)
	if err != nil {
		panic(err)
	}
	return accAddress
}

// ValAddressFromBech32 decodes a bech32 encoded address into an sdk.ValAddress, and validates the prefix (hrp).
// An empty prefix param skips the checking.
// Returns an error if address is zero-length, invalid or the prefix does not match.
func ValAddressFromBech32(address, prefix string) (sdk.ValAddress, error) {
	addr, err := AddressFromBech32(address, prefix)
	return sdk.ValAddress(addr), err
}

// MustValAddressFromBech32 decodes a bech32 encoded address into an sdk.ValAddress, and validates the prefix (hrp).
// An empty prefix param skips the checking.
// Panics if address is zero-length, invalid or the prefix does not match.
func MustValAddressFromBech32(address, prefix string) sdk.ValAddress {
	valAddress, err := ValAddressFromBech32(address, prefix)
	if err != nil {
		panic(err)
	}
	return valAddress
}

// GenerateAccAddressForTest generates a random sdk.AccAddress for test purposes.
func GenerateAccAddressForTest() sdk.AccAddress {
	return sdk.AccAddress(randomutils.GenerateRandomBytes(32))
}

// GenerateValAddressForTest generates a random sdk.ValAddress for test purposes.
func GenerateValAddressForTest() sdk.ValAddress {
	return sdk.ValAddress(randomutils.GenerateRandomBytes(32))
}

// GenerateAddressForTestWithPrefix generates a random bech32 address with the specified prefix for test purposes.
func GenerateAddressForTestWithPrefix(prefix string) string {
	// AccAddress and ValAddress are simple a byte slice, so it doesn't matter this is AccAddress below.
	return MustEncodeAddressToBech32(prefix, GenerateAccAddressForTest())
}

// EncodeAddressToBech32 encodes an sdk.Address interface with the specified prefix.
// Identical behaviour to bech32.ConvertAndDecode(); added to addressutils for consistency.
// Error is thrown if encoding fails.
func EncodeAddressToBech32(prefix string, address sdk.Address) (string, error) {
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}

// EncodeAddressToBech32 encodes an sdk.Address interface with the specified prefix.
// Identical behaviour to bech32.ConvertAndDecode(); added to addressutils for consistency.
// Panics if encoding fails.
func MustEncodeAddressToBech32(prefix string, address sdk.Address) string {
	addr, err := EncodeAddressToBech32(prefix, address)
	if err != nil {
		panic(err)
	}
	return addr
}
