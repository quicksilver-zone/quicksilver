package types

const (
	// ModuleName defines the module name.
	ModuleName = "epochs"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

// prefix bytes for the epoch persistent store.
const (
	prefixEpoch = 0x01
)

const (
	EpochIdentifierEpoch = "epoch"
	EpochIdentifierDay   = "day"
)

// KeyPrefixEpoch defines prefix key for storing epochs.
var KeyPrefixEpoch = []byte{prefixEpoch}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
