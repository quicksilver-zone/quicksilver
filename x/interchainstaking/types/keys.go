package types

import (
	"bytes"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

	GenericToken = "tokens"

	TxRetrieveCount = 100

	QueryParameters         = "params"
	QueryZones              = "zones"
	QueryZoneDepositAddress = "zones/deposit_address"

	ICASuffixDeposit     = "deposit"
	ICASuffixDelegate    = "delegate"
	ICASuffixWithdrawal  = "withdrawal"
	ICASuffixPerformance = "performance"

	BankStoreKey        = "store/bank/key"
	EscrowModuleAccount = "ics-escrow-account"
)

var (
	KeyPrefixZone                        = []byte{0x01}
	KeyPrefixIntent                      = []byte{0x02}
	KeyPrefixPortMapping                 = []byte{0x03}
	KeyPrefixReceipt                     = []byte{0x04}
	KeyPrefixWithdrawalRecord            = []byte{0x05}
	KeyPrefixUnbondingRecord             = []byte{0x06}
	KeyPrefixDelegation                  = []byte{0x07}
	KeyPrefixPerformanceDelegation       = []byte{0x08}
	KeyPrefixSnapshotIntent              = []byte{0x09}
	KeyPrefixRequeuedWithdrawalRecordSeq = []byte{0x0a}
	// fill in missing 0b - 0f before adding 0x11!
	KeyPrefixRedelegationRecord = []byte{0x10}
)

// ParseStakingDelegationKey parses the KV store key for a delegation from Cosmos x/staking module,
// as defined here: https://github.com/cosmos/cosmos-sdk/blob/v0.45.6/x/staking/types/keys.go#L180
func ParseStakingDelegationKey(key []byte) (sdk.AccAddress, sdk.ValAddress, error) {
	if len(key) < 1 {
		return nil, nil, errors.New("out of bounds reading byte 0")
	}
	if !bytes.Equal(key[0:1], []byte{0x31}) {
		return []byte{}, []byte{}, errors.New("not a valid delegation key")
	}
	if len(key) < 2 {
		return nil, nil, errors.New("out of bounds reading delegator address length")
	}
	delAddrLen := int(key[1])
	if len(key) < 2+delAddrLen {
		return nil, nil, errors.New("invalid delegator address length")
	}
	delAddr := key[2 : 2+delAddrLen]
	// use valAddrLen to validate the val address has not been truncated.
	valAddrLen := int(key[2+delAddrLen])
	if len(key) < 3+delAddrLen+valAddrLen {
		return nil, nil, errors.New("out of bounds reading validator address")
	}
	valAddr := key[3+delAddrLen : 3+delAddrLen+valAddrLen]
	return delAddr, valAddr, nil
}
