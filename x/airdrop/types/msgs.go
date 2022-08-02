package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// airdrop message types
const (
	TypeMsgClaim = "claim"
)

// NewMsgClaim constructs a msg to claim from a zone airdrop.
//nolint:interfacer
func NewMsgClaim(chainID string, action int32, fromAddress sdk.Address) *MsgClaim {
	return &MsgClaim{ChainId: chainID, Action: action, Address: fromAddress.String()}
}

// Route implements Msg.
func (msg MsgClaim) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgClaim) Type() string { return TypeMsgClaim }

// ValidateBasic implements Msg.
func (msg MsgClaim) ValidateBasic() error {
	errors := make(map[string]error)

	if msg.ChainId == "" {
		errors["ChainId"] = fmt.Errorf("invalid ChainID, empty string")
	}

	if int(msg.Action) >= len(Action_value) {
		errors["Action"] = fmt.Errorf("invalid action, expects range [0-%d), got %d", len(Action_value), msg.Action)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		errors["Address"] = err
	}

	if len(errors) > 0 {
		return NewMultiError(errors)
	}

	return nil
}

// GetSignBytes implements Msg.
func (msg MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgClaim) GetSigners() []sdk.AccAddress {
	address, _ := sdk.AccAddressFromBech32(msg.Address)
	return []sdk.AccAddress{address}
}
