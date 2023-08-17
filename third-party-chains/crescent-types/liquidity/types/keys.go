package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "liquidity"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
)

var PoolKeyPrefix = []byte{0xab}

// GetPoolKey returns the store key to retrieve pool object from the pool id.
func GetPoolKey(poolID uint64) []byte {
	return append(PoolKeyPrefix, sdk.Uint64ToBigEndian(poolID)...)
}
