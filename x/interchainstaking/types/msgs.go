package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// interchainstaking message types
const (
	TypeMsgRegisterZone      = "registerzone"
	TypeMsgRequestRedemption = "requestredemption"
	TypeMsgSignalIntent      = "signalintent"
)

var _ sdk.Msg = &MsgRegisterZone{}
var _ sdk.Msg = &MsgRequestRedemption{}
var _ sdk.Msg = &MsgSignalIntent{}

// NewMsgRegisterZone - construct a msg to register a new zone.
//nolint:interfacer
func NewMsgRegisterZone(
	identifier string,
	connection_id string,
	chain_id string,
	local_denom string,
	base_denom string,
	from_address sdk.Address,
	multi_send bool,
) *MsgRegisterZone {
	return &MsgRegisterZone{
		Identifier:   identifier,
		ConnectionId: connection_id,
		ChainId:      chain_id,
		LocalDenom:   local_denom,
		BaseDenom:    base_denom,
		FromAddress:  from_address.String(),
		MultiSend:    multi_send,
	}
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

//----------------------------------------------------------------

// NewMsgRequestRedemption - construct a msg to request redemption.
//nolint:interfacer
func NewMsgRequestRedemption(coin string, destination_address string, from_address sdk.Address) *MsgRequestRedemption {
	return &MsgRequestRedemption{Coin: coin, DestinationAddress: destination_address, FromAddress: from_address.String()}
}

// Route Implements Msg.
func (msg MsgRequestRedemption) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRequestRedemption) Type() string { return TypeMsgRegisterZone }

// ValidateBasic Implements Msg.
func (msg MsgRequestRedemption) ValidateBasic() error {
	// TODO: check from address

	// TODO: check for valid identifier

	// TODO: check for valid chain_id

	// TODO: check for valid denominations

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRequestRedemption) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgRequestRedemption) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

//----------------------------------------------------------------

func IntentsFromString(input string) ([]*ValidatorIntent, error) {
	out := []*ValidatorIntent{}
	parts := strings.Split(input, ";")
	for _, val := range parts {
		vparts := strings.SplitN(val, ",", 2)
		// validator should be a valoper address
		// weight should be a float
		// validate me please!
		weight, err := sdk.NewDecFromStr(vparts[1])
		if err != nil {
			return []*ValidatorIntent{}, err
		}
		v := ValidatorIntent{vparts[0], weight}
		out = append(out, &v)
	}
	return out, nil
}

// NewMsgRequestRedemption - construct a msg to request redemption.
//nolint:interfacer
func NewMsgSignalIntent(chain_id string, intents string, from_address sdk.Address) *MsgSignalIntent {
	intent_obj, err := IntentsFromString(intents)
	if err != nil {
		return nil
	}

	return &MsgSignalIntent{ChainId: chain_id, Intents: intent_obj, FromAddress: from_address.String()}
}

// Route Implements Msg.
func (msg MsgSignalIntent) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSignalIntent) Type() string { return TypeMsgRegisterZone }

// ValidateBasic Implements Msg.
func (msg MsgSignalIntent) ValidateBasic() error {
	// TODO: check from address

	// TODO: check for valid identifier

	// TODO: check for valid chain_id

	// TODO: check for valid denominations

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSignalIntent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSignalIntent) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}
