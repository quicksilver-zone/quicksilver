package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "airdrop"
	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the module.
	QuerierRoute = StoreKey
	// RouterKey is the message route for the module.
	RouterKey = ModuleName
)

var (
	KeyPrefixZoneDrop    = []byte{0x01}
	KeyPrefixClaimRecord = []byte{0x02}
)

func GetKeyZoneDrop(chainID string) []byte {
	return append(KeyPrefixZoneDrop, []byte(chainID)...)
}

func GetKeyClaimRecord(chainID string, addr sdk.AccAddress) []byte {
	return append(append(KeyPrefixClaimRecord, []byte(chainID)...), addr...)
}

func GetPrefixClaimRecord(chainID string) []byte {
	return append(KeyPrefixClaimRecord, []byte(chainID)...)
}
