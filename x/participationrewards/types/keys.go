package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "participationrewards"
	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the participationrewards store.
	QuerierRoute = StoreKey
	// RouterKey is the message route for participationrewards.
	RouterKey = ModuleName

	OsmosisParamsKey = "osmosisparams"
)

var KeyPrefixProtocolData = []byte{0x00}

func GetProtocolDataKey(pdType ProtocolDataType, key string) []byte {
	return append(sdk.Uint64ToBigEndian(uint64(pdType)), []byte(key)...)
}

func GetPrefixProtocolDataKey(pdType ProtocolDataType) []byte {
	return sdk.Uint64ToBigEndian(uint64(pdType))
}
