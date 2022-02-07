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

	PortID = ModuleName

	Version = "ics27-1"

	// this value defines the number of delegation accounts per zone. This can only ever increase.
	DelegationAccountCount = 10
)

// prefix bytes for the epoch persistent store
const (
	prefixZone        = iota + 1
	prefixPortMapping = iota + 1
)

var (
	KeyPrefixZone        = []byte{prefixZone}
	KeyPrefixPortMapping = []byte{prefixPortMapping}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
