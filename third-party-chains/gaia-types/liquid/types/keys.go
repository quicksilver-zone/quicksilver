package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

var (
	LiquidValidatorPrefix = []byte{0x8} // key for liquid validator prefix
)

func GetLiquidValidatorKey(operatorAddress sdk.ValAddress) []byte {
	return append(LiquidValidatorPrefix, addressutils.MustLengthPrefix(operatorAddress)...)
}
