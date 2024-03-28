package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name.
	ModuleName = "claimsmanager"
	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the claimsmanager store.
	QuerierRoute = StoreKey
)

var (
	KeyPrefixClaim          = []byte{0x00}
	KeyPrefixLastEpochClaim = []byte{0x01}
	KeySelfConsensusState   = []byte{0x02}
)

// GetGenericKeyClaim returns the key for storing a given claim.
func GetGenericKeyClaim(key []byte, chainID, address string, module ClaimType, srcChainID string) []byte {
	typeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(typeBytes, uint32(module))
	key = append(key, chainID...)
	key = append(key, 0x00)
	key = append(key, address...)
	key = append(key, typeBytes...)
	return append(key, srcChainID...)
}

func GetKeyClaim(chainID, address string, module ClaimType, srcChainID string) []byte {
	return GetGenericKeyClaim(KeyPrefixClaim, chainID, address, module, srcChainID)
}

func GetPrefixClaim(chainID string) []byte {
	return append(KeyPrefixClaim, chainID...)
}

func GetPrefixUserClaim(chainID, address string) []byte {
	key := KeyPrefixClaim
	key = append(key, chainID...)
	key = append(key, 0x00)
	key = append(key, address...)
	return key
}

func GetKeyLastEpochClaim(chainID, address string, module ClaimType, srcChainID string) []byte {
	return GetGenericKeyClaim(KeyPrefixLastEpochClaim, chainID, address, module, srcChainID)
}

func GetPrefixLastEpochClaim(chainID string) []byte {
	return append(KeyPrefixLastEpochClaim, chainID...)
}

func GetPrefixLastEpochUserClaim(chainID, address string) []byte {
	key := KeyPrefixLastEpochClaim
	key = append(key, chainID...)
	key = append(key, 0x00)
	key = append(key, address...)
	return key
}
