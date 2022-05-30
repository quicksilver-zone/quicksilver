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
	// this value defines the number of delegation accounts a given deposit should be shared amongst
	DelegationAccountSplit = 9

	// beginblocker intervals
	DepositInterval             = 5
	DelegateInterval            = 25
	DelegateDelegationsInterval = 100 // probably wants to be somewhere in the region of 1000 (c. 3h) in prod with 7s blocks.
	ValidatorSetInterval        = 25

	QueryParameters                   = "params"
	QueryRegisteredZonesInfo          = "zones"
	QueryRegisteredZoneDepositAddress = "zones/deposit_address"
)

// prefix bytes for the epoch persistent store
const (
	prefixZone             = iota + 1
	prefixIntent           = iota + 1
	prefixPortMapping      = iota + 1
	prefixReceipt          = iota + 1
	prefixWithdrawalRecord = iota + 1
)

var (
	KeyPrefixZone             = []byte{prefixZone}
	KeyPrefixIntent           = []byte{prefixIntent}
	KeyPrefixPortMapping      = []byte{prefixPortMapping}
	KeyPrefixReceipt          = []byte{prefixReceipt}
	KeyPrefixWithdrawalRecord = []byte{prefixWithdrawalRecord}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
