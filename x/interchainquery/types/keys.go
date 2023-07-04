package types

const (
	// ModuleName defines the module name.
	ModuleName = "interchainquery"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for interchainquery.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

// prefix bytes for the interchainquery persistent store.
const (
	prefixData         = 0x01
	prefixQuery        = 0x02
	prefixLatestHeight = 0x03
)

var (
	KeyPrefixData         = []byte{prefixData}
	KeyPrefixQuery        = []byte{prefixQuery}
	KeyPrefixLatestHeight = []byte{prefixLatestHeight}
)
