package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/balancer"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/stableswap"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	gamm.RegisterInterfaces(registry)
	balancer.RegisterInterfaces(registry)
	stableswap.RegisterInterfaces(registry)

	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
