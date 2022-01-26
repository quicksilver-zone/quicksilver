package types

const (
	// ModuleName defines the module name
	ModuleName = "interchainstaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// prefix bytes for the epoch persistent store
const (
	prefixZone = iota + 1
)

// KeyPrefixZone defines prefix key for storing registered zones
var KeyPrefixZone = []byte{prefixZone}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
