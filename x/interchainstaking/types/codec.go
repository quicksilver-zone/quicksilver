package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/third-party-chains/gaia-types/liquid/types"
	"github.com/quicksilver-zone/quicksilver/utils/proofs"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSignalIntent{}, "quicksilver/MsgSignalIntent", nil)
	cdc.RegisterConcrete(&MsgRequestRedemption{}, "quicksilver/MsgRequestRedemption", nil)
	cdc.RegisterConcrete(&MsgCancelRedemption{}, "quicksilver/MsgCancelRedemption", nil)
	cdc.RegisterConcrete(&MsgRequeueRedemption{}, "quicksilver/MsgRequeueRedemption", nil)
	cdc.RegisterConcrete(&RegisterZoneProposal{}, "quicksilver/RegisterZoneProposal", nil)
	cdc.RegisterConcrete(&UpdateZoneProposal{}, "quicksilver/UpdateZoneProposal", nil)
	cdc.RegisterConcrete(&MsgGovExecuteICATx{}, "quicksilver/MsgGovExecuteICATx", nil)
	cdc.RegisterConcrete(&MsgGovSetZoneOffboarding{}, "quicksilver/MsgGovSetZoneOffboarding", nil)
	cdc.RegisterConcrete(&MsgGovCancelAllPendingRedemptions{}, "quicksilver/MsgGovCancelAllPendingRedemptions", nil)
	cdc.RegisterConcrete(&MsgGovForceUnbondAllDelegations{}, "quicksilver/MsgGovForceUnbondAllDelegations", nil)
	lsmstakingtypes.RegisterLegacyAminoCodec(cdc)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSignalIntent{},
		&MsgRequestRedemption{},
		&MsgCancelRedemption{},
		&MsgRequeueRedemption{},
		&MsgGovCloseChannel{},
		&MsgGovReopenChannel{},
		&MsgGovSetLsmCaps{},
		&MsgGovAddValidatorDenyList{},
		&MsgGovRemoveValidatorDenyList{},
		&MsgGovExecuteICATx{},
		&MsgGovSetZoneOffboarding{},
		&MsgGovCancelAllPendingRedemptions{},
		&MsgGovForceUnbondAllDelegations{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&UpdateZoneProposal{},
		&RegisterZoneProposal{},
	)

	proofs.RegisterInterfaces(registry)

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
