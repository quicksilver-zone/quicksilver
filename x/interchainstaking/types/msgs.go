package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// bank message types
const (
	TypeMsgRegisterZone = "registerzone"
)

var _ sdk.Msg = &MsgRegisterZone{}

// NewMsgRegisterZone - construct a msg to send coins from one account to another.
//nolint:interfacer
func NewMsgRegisterZone(identifier string, chain_id string, local_denom string, remote_denom string, from_address sdk.Address) *MsgRegisterZone {
	return &MsgRegisterZone{Identifier: identifier, ChainId: chain_id, LocalDenom: local_denom, RemoteDenom: remote_denom, FromAddress: from_address.String()}
}

// Route Implements Msg.
func (msg MsgRegisterZone) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRegisterZone) Type() string { return TypeMsgRegisterZone }

// ValidateBasic Implements Msg.
func (msg MsgRegisterZone) ValidateBasic() error {
	// TODO: check from address

	// TODO: check for valid identifier

	// TODO: check for valid chain_id

	// TODO: check for valid denominations

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRegisterZone) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgRegisterZone) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}
