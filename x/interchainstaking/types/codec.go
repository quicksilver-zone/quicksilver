package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/x/lsmtypes"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSignalIntent{}, "quicksilver/MsgSignalIntent", nil)
	cdc.RegisterConcrete(&MsgRequestRedemption{}, "quicksilver/MsgRequestRedemption", nil)
	cdc.RegisterConcrete(&RegisterZoneProposal{}, "quicksilver/RegisterZoneProposal", nil)
	cdc.RegisterConcrete(&UpdateZoneProposal{}, "quicksilver/UpdateZoneProposal", nil)
	lsmstakingtypes.RegisterLegacyAminoCodec(cdc)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSignalIntent{},
		&MsgRequestRedemption{},
		&MsgGovCloseChannel{},
		&MsgGovReopenChannel{},
		&MsgGovSetLsmCaps{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&UpdateZoneProposal{},
		&RegisterZoneProposal{},
	)

	lsmstakingtypes.RegisterInterfaces(registry)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	govv1beta1.RegisterProposalType(ProposalTypeRegisterZone)
	govv1beta1.RegisterProposalType(ProposalTypeUpdateZone)
	sdk.RegisterLegacyAminoCodec(amino)
}
