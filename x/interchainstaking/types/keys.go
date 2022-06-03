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

var (
	KeyPrefixZone             = []byte{0x01}
	KeyPrefixIntent           = []byte{0x02}
	KeyPrefixPortMapping      = []byte{0x03}
	KeyPrefixReceipt          = []byte{0x04}
	KeyPrefixWithdrawalRecord = []byte{0x05}
	KeyPrefixDelegation       = []byte{0x06}

//	KeyPrefixDelegationByValidator = []byte{0x08}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
