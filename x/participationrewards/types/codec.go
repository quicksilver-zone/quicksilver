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
	cdc.RegisterConcrete(&MsgSubmitClaim{}, "cosmos-sdk/MsgSubmitClaim", nil)
	cdc.RegisterConcrete(&AddProtocolDataProposal{}, "cosmos-sdk/AddProtocolDataProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSubmitClaim{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddProtocolDataProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	cryptocodec.RegisterCrypto(amino)
	govtypes.RegisterProposalType(ProposalTypeAddProtocolData)
	govtypes.RegisterProposalTypeCodec(&AddProtocolDataProposal{}, "cosmos-sdk/AddProtocolDataProposal")
}
