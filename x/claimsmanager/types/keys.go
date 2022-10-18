package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "claimsmanager"
	// StoreKey defines the primary module store key
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the claimsmanager store.
	QuerierRoute = StoreKey
	// RouterKey is the message route for claimsmanager
	RouterKey = ModuleName
)

var (
	KeyPrefixClaim          = []byte{0x00}
	KeyPrefixLastEpochClaim = []byte{0x01}
)

// ClaimKey returns the key for storing a given claim.
func GetGenericKeyClaim(key []byte, chainID string, address string, module ClaimType, srcChainID string) []byte {
	typeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(typeBytes, uint32(module))
	key = append(key, []byte(chainID)...)
	key = append(key, []byte(address)...)
	key = append(key, typeBytes...)
	return append(key, []byte(srcChainID)...)
}

func GetKeyClaim(chainID string, address string, module ClaimType, srcChainID string) []byte {
	return GetGenericKeyClaim(KeyPrefixClaim, chainID, address, module, srcChainID)
}

func GetPrefixClaim(chainID string) []byte {
	return append(KeyPrefixClaim, []byte(chainID)...)
}

func GetPrefixUserClaim(chainID string, address string) []byte {
	return append(append(KeyPrefixClaim, []byte(chainID)...), []byte(address)...)
}

func GetKeyLastEpochClaim(chainID string, address string, module ClaimType, srcChainID string) []byte {
	return GetGenericKeyClaim(KeyPrefixLastEpochClaim, chainID, address, module, srcChainID)
}

func GetPrefixLastEpochClaim(chainID string) []byte {
	return append(KeyPrefixLastEpochClaim, []byte(chainID)...)
}

func GetPrefixLastEpochUserClaim(chainID string, address string) []byte {
	return append(append(KeyPrefixLastEpochClaim, []byte(chainID)...), []byte(address)...)
}
