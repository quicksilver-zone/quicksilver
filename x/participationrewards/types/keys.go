package types

import (
	"encoding/binary"
)

type ProtocolDataType int32

const (
	// ModuleName defines the module name
	ModuleName = "participationrewards"
	// StoreKey defines the primary module store key
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the participationrewards store.
	QuerierRoute = StoreKey
	// RouterKey is the message route for participationrewards
	RouterKey = ModuleName

	ProtocolDataConnection   ProtocolDataType = 0
	ProtocolDataLiquidToken  ProtocolDataType = 1
	ProtocolDataOsmosisPool  ProtocolDataType = 2
	ProtocolDataCrescentPool ProtocolDataType = 3
	ProtocolDataSifchainPool ProtocolDataType = 4
)

var (
	KeyPrefixProtocolData   = []byte{0x00}
	KeyPrefixClaim          = []byte{0x01}
	KeyPrefixLastEpochClaim = []byte{0x02}
)

var ProtocolDataType_name = map[ProtocolDataType]string{ //nolint:revive,stylecheck // conform with protobuf3 enum
	ProtocolDataConnection:   "connection",
	ProtocolDataLiquidToken:  "liquidtoken",
	ProtocolDataOsmosisPool:  "osmosispool",
	ProtocolDataCrescentPool: "crescentpool",
	ProtocolDataSifchainPool: "sifchainpool",
}

var ProtocolDataType_value = map[string]ProtocolDataType{ //nolint:revive,stylecheck // conform with protobuf3 enum
	"connection":   ProtocolDataConnection,
	"liquidtoken":  ProtocolDataLiquidToken,
	"osmosispool":  ProtocolDataOsmosisPool,
	"crescentpool": ProtocolDataCrescentPool,
	"sifchainpool": ProtocolDataSifchainPool,
}

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
