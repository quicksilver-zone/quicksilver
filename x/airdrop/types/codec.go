package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgClaim{}, "quicksilver/MsgClaim", nil)
	cdc.RegisterConcrete(&MsgIncentivePoolSpend{}, "quicksilver/MsgIncentivePoolSpend", nil)
	cdc.RegisterConcrete(&RegisterZoneDropProposal{}, "quicksilver/RegisterZoneDropProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaim{},
		&MsgIncentivePoolSpend{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&RegisterZoneDropProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	govv1beta1.RegisterProposalType(ProposalTypeRegisterZoneDrop)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
