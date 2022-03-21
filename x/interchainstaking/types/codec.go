package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterZone{}, "cosmos-sdk/MsgRegisterZone", nil)
	cdc.RegisterConcrete(&MsgSignalIntent{}, "cosmos-sdk/MsgSignalIntent", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterZone{},
		&MsgSignalIntent{},
	)
	// registry.RegisterImplementations(
	// 	(*authz.Authorization)(nil),
	// 	&SendAuthorization{},
	// )

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
func init() {
	cryptocodec.RegisterCrypto(amino)
}
