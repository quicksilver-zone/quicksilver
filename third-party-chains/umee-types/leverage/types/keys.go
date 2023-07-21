package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"

	"github.com/ingenuity-build/quicksilver/utils"
)

const (
	LeverageModuleName = "leverage"

	// StoreKey defines the primary module store key.
	StoreKey = LeverageModuleName
)

// KVStore key prefixes.
var (
	KeyPrefixCollateralAmount    = []byte{0x04}
	KeyPrefixReserveAmount       = []byte{0x05}
	KeyPrefixInterestScalar      = []byte{0x08}
	KeyPrefixAdjustedTotalBorrow = []byte{0x09}
	KeyPrefixUtokenSupply        = []byte{0x0A}
)

func KeyReserveAmount(tokenDenom string) []byte {
	// reserveamountprefix | denom | 0x00 for null-termination.
	return utils.ConcatBytes(1, KeyPrefixReserveAmount, []byte(tokenDenom))
}

// KeyAdjustedTotalBorrow returns a KVStore key for getting and setting the total ajdusted borrows for
// a given token.
func KeyAdjustedTotalBorrow(tokenDenom string) []byte {
	// totalBorrowedPrefix | denom | 0x00 for null-termination.
	return utils.ConcatBytes(1, KeyPrefixAdjustedTotalBorrow, []byte(tokenDenom))
}

// KeyInterestScalar returns a KVStore key for getting and setting the interest scalar for a
// given token.
func KeyInterestScalar(tokenDenom string) []byte {
	// interestScalarPrefix | denom | 0x00 for null-termination.
	return utils.ConcatBytes(1, KeyPrefixInterestScalar, []byte(tokenDenom))
}

// KeyUTokenSupply returns a KVStore key for getting and setting a utoken's total supply.
func KeyUTokenSupply(uTokenDenom string) []byte {
	// supplyprefix | denom | 0x00 for null-termination.
	return utils.ConcatBytes(1, KeyPrefixUtokenSupply, []byte(uTokenDenom))
}

// KeyCollateralAmount returns a KVStore key for getting and setting the amount of
// collateral stored for a user in a given denom.
func KeyCollateralAmount(addr sdk.AccAddress, uTokenDenom string) []byte {
	// collateralPrefix | lengthprefixed(addr) | denom | 0x00 for null-termination
	return utils.ConcatBytes(1, KeyCollateralAmountNoDenom(addr), []byte(uTokenDenom))
}

// KeyCollateralAmountNoDenom returns the common prefix used by all collateral associated
// with a given address.
func KeyCollateralAmountNoDenom(addr sdk.AccAddress) []byte {
	return utils.ConcatBytes(0, KeyPrefixCollateralAmount, address.MustLengthPrefix(addr))
}

// DenomFromKey extracts denom from a key with the form
// prefix | denom | 0x00.
func DenomFromKey(key, prefix []byte) string {
	return string(key[len(prefix) : len(key)-1])
}
