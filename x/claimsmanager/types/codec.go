package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	cltypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity"
	clpool "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/model"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/balancer"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/stableswap"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/poolmanager"
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
	poolmanager.RegisterInterfaces(registry)
	cltypes.RegisterInterfaces(registry)
	clpool.RegisterInterfaces(registry)

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
