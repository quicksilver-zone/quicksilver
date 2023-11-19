package lsmtypes

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUnbondValidator{}
	_ sdk.Msg = &MsgTokenizeShares{}
	_ sdk.Msg = &MsgRedeemTokensForShares{}
	_ sdk.Msg = &MsgTransferTokenizeShareRecord{}
	_ sdk.Msg = &MsgDisableTokenizeShares{}
	_ sdk.Msg = &MsgEnableTokenizeShares{}
	_ sdk.Msg = &MsgValidatorBond{}
)

const (
	TypeMsgUndelegate                           = "begin_unbonding"
	TypeMsgUnbondValidator                      = "unbond_validator"
	TypeMsgEditValidator                        = "edit_validator"
	TypeMsgCreateValidator                      = "create_validator"
	TypeMsgDelegate                             = "delegate"
	TypeMsgBeginRedelegate                      = "begin_redelegate"
	TypeMsgCancelUnbondingDelegation            = "cancel_unbond"
	TypeMsgTokenizeShares                       = "tokenize_shares"
	TypeMsgRedeemTokensForShares                = "redeem_tokens_for_shares" //nolint:gosec
	TypeMsgTransferTokenizeShareRecord          = "transfer_tokenize_share_record"
	TypeMsgDisableTokenizeShares                = "disable_tokenize_shares"
	TypeMsgEnableTokenizeShares                 = "enable_tokenize_shares" //nolint:gosec
	TypeMsgValidatorBond                        = "validator_bond"
	TypeMsgWithdrawTokenizeShareRecordReward    = "withdraw_tokenize_share_record_reward"
	TypeMsgWithdrawAllTokenizeShareRecordReward = "withdraw_all_tokenize_share_record_reward"

	RouterKey  = "lsm"
	ModuleName = "lsm"
)

// NewMsgUnbondValidator creates a new MsgUnbondValidator instance.
//
//nolint:interfacer
func NewMsgUnbondValidator(valAddr sdk.ValAddress) *MsgUnbondValidator {
	return &MsgUnbondValidator{
		ValidatorAddress: valAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUnbondValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUnbondValidator) Type() string { return TypeMsgUnbondValidator }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUnbondValidator) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{valAddr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUnbondValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUnbondValidator) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	return nil
}

// NewMsgTokenizeShares creates a new MsgTokenizeShares instance.
//
//nolint:interfacer
func NewMsgTokenizeShares(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin, owner sdk.AccAddress) *MsgTokenizeShares {
	return &MsgTokenizeShares{
		DelegatorAddress:    delAddr.String(),
		ValidatorAddress:    valAddr.String(),
		Amount:              amount,
		TokenizedShareOwner: owner.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgTokenizeShares) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgTokenizeShares) Type() string { return TypeMsgTokenizeShares }

// GetSigners implements the sdk.Msg interface.
func (msg MsgTokenizeShares) GetSigners() []sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delegator}
}

// MsgTokenizeShares implements the sdk.Msg interface.
func (msg MsgTokenizeShares) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgTokenizeShares) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.TokenizedShareOwner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid tokenize share owner address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	return nil
}

// NewMsgRedeemTokensForShares creates a new MsgRedeemTokensForShares instance.
//
//nolint:interfacer
func NewMsgRedeemTokensForShares(delAddr sdk.AccAddress, amount sdk.Coin) *MsgRedeemTokensForShares {
	return &MsgRedeemTokensForShares{
		DelegatorAddress: delAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgRedeemTokensForShares) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRedeemTokensForShares) Type() string { return TypeMsgRedeemTokensForShares }

// GetSigners implements the sdk.Msg interface.
func (msg MsgRedeemTokensForShares) GetSigners() []sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgRedeemTokensForShares) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRedeemTokensForShares) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	return nil
}

// NewMsgTransferTokenizeShareRecord creates a new MsgTransferTokenizeShareRecord instance.
//
//nolint:interfacer
func NewMsgTransferTokenizeShareRecord(recordID uint64, sender, newOwner sdk.AccAddress) *MsgTransferTokenizeShareRecord {
	return &MsgTransferTokenizeShareRecord{
		TokenizeShareRecordId: recordID,
		Sender:                sender.String(),
		NewOwner:              newOwner.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgTransferTokenizeShareRecord) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgTransferTokenizeShareRecord) Type() string { return TypeMsgTransferTokenizeShareRecord }

// GetSigners implements the sdk.Msg interface.
func (msg MsgTransferTokenizeShareRecord) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgTransferTokenizeShareRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgTransferTokenizeShareRecord) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.NewOwner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid new owner address: %s", err)
	}

	return nil
}

// NewMsgDisableTokenizeShares creates a new MsgDisableTokenizeShares instance.
//
//nolint:interfacer
func NewMsgDisableTokenizeShares(delAddr sdk.AccAddress) *MsgDisableTokenizeShares {
	return &MsgDisableTokenizeShares{
		DelegatorAddress: delAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDisableTokenizeShares) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDisableTokenizeShares) Type() string { return TypeMsgDisableTokenizeShares }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDisableTokenizeShares) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDisableTokenizeShares) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDisableTokenizeShares) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	return nil
}

// NewMsgEnableTokenizeShares creates a new MsgEnableTokenizeShares instance.
//
//nolint:interfacer
func NewMsgEnableTokenizeShares(delAddr sdk.AccAddress) *MsgEnableTokenizeShares {
	return &MsgEnableTokenizeShares{
		DelegatorAddress: delAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEnableTokenizeShares) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEnableTokenizeShares) Type() string { return TypeMsgEnableTokenizeShares }

// GetSigners implements the sdk.Msg interface.
func (msg MsgEnableTokenizeShares) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgEnableTokenizeShares) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEnableTokenizeShares) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	return nil
}

// NewMsgValidatorBond creates a new MsgValidatorBond instance.
//
//nolint:interfacer
func NewMsgValidatorBond(delAddr sdk.AccAddress, valAddr sdk.ValAddress) *MsgValidatorBond {
	return &MsgValidatorBond{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgValidatorBond) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgValidatorBond) Type() string { return TypeMsgValidatorBond }

// GetSigners implements the sdk.Msg interface.
func (msg MsgValidatorBond) GetSigners() []sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgValidatorBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgValidatorBond) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	return nil
}

func NewMsgWithdrawTokenizeShareRecordReward(ownerAddr sdk.AccAddress, recordID uint64) *MsgWithdrawTokenizeShareRecordReward {
	return &MsgWithdrawTokenizeShareRecordReward{
		OwnerAddress: ownerAddr.String(),
		RecordId:     recordID,
	}
}

func (msg MsgWithdrawTokenizeShareRecordReward) Route() string { return ModuleName }
func (msg MsgWithdrawTokenizeShareRecordReward) Type() string {
	return TypeMsgWithdrawTokenizeShareRecordReward
}

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawTokenizeShareRecordReward) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawTokenizeShareRecordReward) GetSignBytes() []byte {
	bz := legacy.Cdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawTokenizeShareRecordReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.OwnerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}
	return nil
}

func NewMsgWithdrawAllTokenizeShareRecordReward(ownerAddr sdk.AccAddress) *MsgWithdrawAllTokenizeShareRecordReward {
	return &MsgWithdrawAllTokenizeShareRecordReward{
		OwnerAddress: ownerAddr.String(),
	}
}

func (msg MsgWithdrawAllTokenizeShareRecordReward) Route() string { return ModuleName }

func (msg MsgWithdrawAllTokenizeShareRecordReward) Type() string {
	return TypeMsgWithdrawAllTokenizeShareRecordReward
}

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawAllTokenizeShareRecordReward) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawAllTokenizeShareRecordReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawAllTokenizeShareRecordReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.OwnerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}
	return nil
}
