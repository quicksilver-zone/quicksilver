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
	EventStatusUnspecified = int32(0)
	EventStatusActive      = int32(1)
	EventStatusPending     = int32(2)

	EventTypeUnspecified            = int32(0x00)
	EventTypeICQQueryRewards        = int32(0x01)
	EventTypeICQQueryDelegations    = int32(0x02)
	EventTypeICQQueryDelegation     = int32(0x03)
	EventTypeICQAccountBalances     = int32(0x04)
	EventTypeICQAccountBalance      = int32(0x05)
	EventTypeICAWithdrawRewards     = int32(0x06)
	EventTypeICADelegate            = int32(0x07)
	EventTypeICAUnbond              = int32(0x08)
	EventTypeICQGetLatestBlock      = int32(0x09)
	EventTypeICQOsmosisPool         = int32(0x0a)
	EventTypeSubmodules             = int32(0x0b)
	EventTypeCalculateTvls          = int32(0x0c)
	EventTypeDistributeRewards      = int32(0x0d)
	EventTypeICQUmeeUTokenSupply    = int32(0x0e)
	EventTypeICQUmeeInterestScalar  = int32(0x0f)
	EventTypeICQUmeeLeverageBalance = int32(0x10)
	EventTypeICQUmeeReserves        = int32(0x11)
	EventTypeICQUmeeTotalBorrows    = int32(0x12)

	FieldEventType   = "eventtype"
	FieldModule      = "module"
	FieldEventStatus = "eventstatus"
	FieldChainID     = "chainid"
	FieldIdentifier  = "identifier"
	FieldCallback    = "callback"
)

var KeyPrefixEvent = []byte{0x01}