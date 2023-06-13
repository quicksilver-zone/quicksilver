package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "interchainstaking"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for interchainstaking.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
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
	KeyPrefixAddressZoneMapping          = []byte{0x0b}
	KeyPrefixValidatorsInfo              = []byte{0x0c}
	KeyPrefixRemoteAddress               = []byte{0x0d}
	KeyPrefixLocalAddress                = []byte{0x0e}

	// fill in missing 0d - 0f before adding 0x11!
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

// GetRemoteAddressKey gets the prefix for a remote address mapping.
func GetRemoteAddressKey(localAddress []byte, chainID string) []byte {
	return append(append(KeyPrefixRemoteAddress, localAddress...), []byte(chainID)...)
}

// GetLocalAddressKey gets the prefix for a local address mapping.
func GetLocalAddressKey(remoteAddress []byte, chainID string) []byte {
	return append(append(KeyPrefixLocalAddress, []byte(chainID)...), remoteAddress...)
}

// GetDelegationKey gets the key for delegator bond with validator.
// VALUE: staking/Delegation.
func GetDelegationKey(zone *Zone, delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(zone, delAddr), valAddr.Bytes()...)
}

// GetDelegationsKey gets the prefix for a delegator for all validators.
func GetDelegationsKey(zone *Zone, delAddr sdk.AccAddress) []byte {
	return append(append(KeyPrefixDelegation, []byte(zone.ChainId)...), delAddr.Bytes()...)
}

// GetPerformanceDelegationKey gets the key for delegator bond with validator.
// VALUE: staking/Delegation.
func GetPerformanceDelegationKey(zone *Zone, delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetPerformanceDelegationsKey(zone, delAddr), valAddr.Bytes()...)
}

// GetPerformanceDelegationsKey gets the prefix for a delegator for all validators.
func GetPerformanceDelegationsKey(zone *Zone, delAddr sdk.AccAddress) []byte {
	return append(append(KeyPrefixPerformanceDelegation, []byte(zone.ChainId)...), delAddr.Bytes()...)
}

func GetReceiptKey(chainID, txhash string) string {
	return fmt.Sprintf("%s/%s", chainID, strings.ToUpper(txhash))
}

// GetRedelegationKey gets the redelegation key.
// Unbondigng records are keyed by chainId, validator and epoch, as they must be unique with regard to this triple.
func GetRedelegationKey(chainID, source, destination string, epochNumber int64) []byte {
	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, uint64(epochNumber))
	return append(KeyPrefixRedelegationRecord, append(append([]byte(chainID), []byte(source+destination)...), epochBytes...)...)
}

func GetWithdrawalKey(chainID string, status int32) []byte {
	statusBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(statusBytes, uint64(status))
	key := KeyPrefixWithdrawalRecord
	key = append(key, append([]byte(chainID), statusBytes...)...)
	return key
}

// GetUnbondingKey gets the unbonding key.
// unbonding records are keyed by chainId, validator and epoch, as they must be unique with regard to this triple.
func GetUnbondingKey(chainID, validator string, epochNumber int64) []byte {
	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, uint64(epochNumber))
	return append(KeyPrefixUnbondingRecord, append(append([]byte(chainID), []byte(validator)...), epochBytes...)...)
}

// GetZoneValidatorsKey gets the validators key prefix for a given chain.
func GetZoneValidatorsKey(chainID string) []byte {
	return append(KeyPrefixValidatorsInfo, []byte(chainID)...)
}

// GetRemoteAddressPrefix gets the prefix for a remote address mapping.
func GetRemoteAddressPrefix(locaAddress []byte) []byte {
	return append(KeyPrefixRemoteAddress, locaAddress...)
}
