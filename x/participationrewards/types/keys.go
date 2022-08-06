package types

const (
	// ModuleName defines the module name
	ModuleName = "participationrewards"
	// StoreKey defines the primary module store key
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the participationrewards store.
	QuerierRoute = StoreKey
	// RouterKey is the message route for participationrewards
	RouterKey = ModuleName

	ClaimTypeLiquidToken  = 0
	ClaimTypeOsmosisPool  = 1
	ClaimTypeCrescentPool = 2
	ClaimTypeSifchainPool = 3
)

var (
	KeyPrefixProtocolData = []byte{0x00}
	KeyPrefixClaim        = []byte{0x01}
)
