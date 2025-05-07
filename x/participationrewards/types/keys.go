package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "participationrewards"
	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
	// QuerierRoute is the querier route for the participationrewards store.
	QuerierRoute = StoreKey
	// RouterKey is the message route for participationrewards.
	RouterKey = ModuleName

	OsmosisParamsKey  = "osmosisparams"
	MembraneParamsKey = "membraneparams"
	UmeeParamsKey     = "umeeparams"
	CrescentParamsKey = "crescentparams"
	ProofTypeBank     = "bank"
	ProofTypeLeverage = "leverage"
	ProofTypeLPFarm   = "lpfarm"
)

var KeyPrefixProtocolData = []byte{0x00}

func GetProtocolDataKey(pdType ProtocolDataType, key []byte) []byte {
	if pdType < 1 {
		panic(fmt.Sprintf("protocol data type is negative or undefined: %d", pdType))
	}
	return append(sdk.Uint64ToBigEndian(uint64(pdType)), key...) //nolint:gosec
}

func GetPrefixProtocolDataKey(pdType ProtocolDataType) []byte {
	if pdType < 1 {
		panic(fmt.Sprintf("protocol data type is negative or undefined: %d", pdType))
	}
	return sdk.Uint64ToBigEndian(uint64(pdType)) //nolint:gosec
}
