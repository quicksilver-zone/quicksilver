package types

import (
	"bytes"
	"github.com/ingenuity-build/quicksilver/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "lpfarm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

var (
	LastBlockTimeKey           = []byte{0xd0}
	LastPlanIdKey              = []byte{0xd1}
	NumPrivatePlansKey         = []byte{0xd2}
	PlanKeyPrefix              = []byte{0xd3}
	FarmKeyPrefix              = []byte{0xd4}
	PositionKeyPrefix          = []byte{0xd5}
	HistoricalRewardsKeyPrefix = []byte{0xd6}
)

func GetPlanKey(id uint64) []byte {
	return append(PlanKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetFarmKey(denom string) []byte {
	return append(FarmKeyPrefix, denom...)
}

func GetPositionKey(farmerAddr sdk.AccAddress, denom string) []byte {
	return append(append(PositionKeyPrefix, address.MustLengthPrefix(farmerAddr)...), denom...)
}

// GetPositionsByFarmerKeyPrefix returns a key prefix for iterating through
// all the positions owned by a farmer.
func GetPositionsByFarmerKeyPrefix(farmerAddr sdk.AccAddress) []byte {
	return append(PositionKeyPrefix, address.MustLengthPrefix(farmerAddr)...)
}

func GetHistoricalRewardsKey(denom string, period uint64) []byte {
	return append(append(HistoricalRewardsKeyPrefix, utils.LengthPrefixString(denom)...), sdk.Uint64ToBigEndian(period)...)
}

// GetHistoricalRewardsByDenomKeyPrefix returns a key prefix for iterating
// through all the historical rewards belong to a denom.
func GetHistoricalRewardsByDenomKeyPrefix(denom string) []byte {
	return append(HistoricalRewardsKeyPrefix, utils.LengthPrefixString(denom)...)
}

func ParseFarmKey(key []byte) (denom string) {
	if !bytes.HasPrefix(key, FarmKeyPrefix) {
		panic("key does not have proper prefix")
	}
	denom = string(key[1:])
	return
}

func ParsePositionKey(key []byte) (farmerAddr sdk.AccAddress, denom string) {
	if !bytes.HasPrefix(key, PositionKeyPrefix) {
		panic("key does not have proper prefix")
	}
	farmerAddrLen := key[1]
	farmerAddr = key[2 : 2+farmerAddrLen]
	denom = string(key[2+farmerAddrLen:])
	return
}

func ParseHistoricalRewardsKey(key []byte) (denom string, period uint64) {
	if !bytes.HasPrefix(key, HistoricalRewardsKeyPrefix) {
		panic("key does not have proper prefix")
	}
	denomLen := key[1]
	denom = string(key[2 : 2+denomLen])
	period = sdk.BigEndianToUint64(key[2+denomLen:])
	return
}
