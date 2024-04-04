package types

const (
	// ModuleName defines the module name.
	ModuleName = "eventmanager"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for interchainquery.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

const (
	EventStatusUnspecified = 0
	EventStatusActive      = 1
	EventStatusPending     = 2
)

const (
	EventTypeUnspecified         = 0x00
	EventTypeICQQueryRewards     = 0x01
	EventTypeICQQueryDelegations = 0x02
	EventTypeICQAccountBalances  = 0x03
	EventTypeICQAccountBalance   = 0x04
	EventTypeICAWithdrawRewards  = 0x05
	EventTypeICADelegate         = 0x06
	EventTypeICAUnbond           = 0x07
)

var (
	KeyPrefixEvent = []byte{0x01}
)
