package stableswap

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	types "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/types"
	poolmanagertypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/poolmanager/types"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Pool{}, "osmosis/gamm/StableswapPool", nil)
	cdc.RegisterConcrete(&PoolParams{}, "osmosis/gamm/StableswapPoolParams", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"osmosis.poolmanager.v1beta1.PoolI",
		(*poolmanagertypes.PoolI)(nil),
		&Pool{},
	)
	registry.RegisterInterface(
		"osmosis.gamm.v1beta1.PoolI", // N.B.: the old proto-path is preserved for backwards-compatibility.
		(*types.CFMMPoolI)(nil),
		&Pool{},
	)

	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

const PoolTypeName string = "Stableswap"
