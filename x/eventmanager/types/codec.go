package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("quicksilver.eventmanager.v1.ConditionI", (*ConditionI)(nil))
	registry.RegisterImplementations((*ConditionI)(nil), &ConditionAll{}, &ConditionAnd{} /*&ConditionAny{},*/, &ConditionOr{})
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
