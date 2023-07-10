package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "lpfarm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

var (
	PositionKeyPrefix = []byte{0xd5}
)

func GetPositionKey(farmerAddr sdk.AccAddress, denom string) []byte {
	return append(append(PositionKeyPrefix, address.MustLengthPrefix(farmerAddr)...), denom...)
}
