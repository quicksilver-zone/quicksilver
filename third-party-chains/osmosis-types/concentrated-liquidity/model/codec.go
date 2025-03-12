package model

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	types "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/types"
	poolmanagertypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/poolmanager/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Pool{}, "osmosis/cl-pool", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"osmosis.swaprouter.v1beta1.PoolI",
		(*poolmanagertypes.PoolI)(nil),
		&Pool{},
	)

	registry.RegisterInterface(
		"osmosis.concentratedliquidity.v1beta1.ConcentratedPoolExtension",
		(*types.ConcentratedPoolExtension)(nil),
		&Pool{},
	)

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
