package types

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
	KeyPrefixProtocolData = []byte{0x00}
	KeyPrefixClaim        = []byte{0x01}
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
func GetKeyClaim(chainID string, address string) []byte {
	return append(append(KeyPrefixClaim, []byte(chainID)...), []byte(address)...)
}

func GetPrefixClaim(chainID string) []byte {
	return append(KeyPrefixClaim, []byte(chainID)...)
}
