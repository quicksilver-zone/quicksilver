package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm/pool-models/balancer"
	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm/pool-models/stableswap"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitClaim{}, "quicksilver/MsgSubmitClaim", nil)
	cdc.RegisterConcrete(&AddProtocolDataProposal{}, "quicksilver/AddProtocolDataProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	gamm.RegisterInterfaces(registry)
	balancer.RegisterInterfaces(registry)
	stableswap.RegisterInterfaces(registry)

	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSubmitClaim{},
		&MsgGovRemoveProtocolData{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&AddProtocolDataProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	govv1beta1.RegisterProposalType(ProposalTypeAddProtocolData)
	sdk.RegisterLegacyAminoCodec(amino)
}
