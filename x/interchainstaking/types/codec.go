package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSignalIntent{},
		&MsgRequestRedemption{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateZoneProposal{},
		&RegisterZoneProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	govtypes.RegisterProposalType(ProposalTypeRegisterZone)
	govtypes.RegisterProposalTypeCodec(&RegisterZoneProposal{}, "quicksilver/RegisterZoneProposal")

	govtypes.RegisterProposalType(ProposalTypeUpdateZone)
	govtypes.RegisterProposalTypeCodec(&UpdateZoneProposal{}, "quicksilver/UpdateZoneProposal")
}
