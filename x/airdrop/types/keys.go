package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "airdrop"
	// StoreKey defines the primary module store key
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey
)

var (
	KeyPrefixZoneDrop    = []byte{0x01}
	KeyPrefixClaimRecord = []byte{0x02}
)

func GetKeyZoneDrop(chainId string) []byte {
	return append(KeyPrefixZoneDrop, []byte(chainId)...)
}

func GetKeyClaimRecord(chainId string, addr sdk.AccAddress) []byte {
	return append(append(KeyPrefixClaimRecord, []byte(chainId)...), addr...)
}
