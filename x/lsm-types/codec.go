package lsm_types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers the necessary x/distribution interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgWithdrawTokenizeShareRecordReward{}, "cosmos-sdk/MsgWithdrawTokenizeShareRecordReward", nil)
	cdc.RegisterConcrete(&MsgWithdrawAllTokenizeShareRecordReward{}, "cosmos-sdk/MsgWithdrawAllTokenizeShareRecordReward", nil)
	cdc.RegisterConcrete(&MsgCancelUnbondingDelegation{}, "cosmos-sdk/MsgCancelUnbondingDelegation", nil)
	cdc.RegisterConcrete(&MsgValidatorBond{}, "cosmos-sdk/MsgValidatorBond", nil)
	cdc.RegisterConcrete(&MsgUnbondValidator{}, "cosmos-sdk/MsgUnbondValidator", nil)
	cdc.RegisterConcrete(&MsgTokenizeShares{}, "cosmos-sdk/MsgTokenizeShares", nil)
	cdc.RegisterConcrete(&MsgRedeemTokensForShares{}, "cosmos-sdk/MsgRedeemTokensForShares", nil)
	cdc.RegisterConcrete(&MsgTransferTokenizeShareRecord{}, "cosmos-sdk/MsgTransferTokenizeRecord", nil)
	cdc.RegisterConcrete(&MsgDisableTokenizeShares{}, "cosmos-sdk/MsgDisableTokenizeShares", nil)
	cdc.RegisterConcrete(&MsgEnableTokenizeShares{}, "cosmos-sdk/MsgEnableTokenizeShares", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgWithdrawTokenizeShareRecordReward{},
		&MsgWithdrawAllTokenizeShareRecordReward{},
		&MsgCancelUnbondingDelegation{},
		&MsgValidatorBond{},
		&MsgUnbondValidator{},
		&MsgTokenizeShares{},
		&MsgRedeemTokensForShares{},
		&MsgTransferTokenizeShareRecord{},
		&MsgDisableTokenizeShares{},
		&MsgEnableTokenizeShares{},
	)
}

var (
	amino = codec.NewLegacyAmino()
	// ModuleCdc references the global x/staking module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
}
