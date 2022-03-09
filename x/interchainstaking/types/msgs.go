package types

import (
	"fmt"
	"regexp"
	"strconv"
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

// IntentsFromString validates and parses the given string into a slice
// containing pointers to ValidatorIntent. The combined weights must be 1.0 and
// the valoper addresses must be of valid Bech32 string format.
//
// Tokens are comma separated, e.g.
// "0.3cosmos1xxxxxxxxx,0.3cosmos1yyyyyyyyy,0.4cosmos1zzzzzzzzz".
func IntentsFromString(input string) ([]*ValidatorIntent, error) {
	iexpr := regexp.MustCompile(`(\d.\d+)(\w+1\w+)`)
	pexpr := regexp.MustCompile(fmt.Sprintf("%s(,%s)*", iexpr.String(), iexpr.String()))
	if !pexpr.MatchString(input) {
		return nil, fmt.Errorf("invalid input string")
	}

	out := []*ValidatorIntent{}
	wsum := 0.0

	istrs := strings.Split(input, ",")
	for i, istr := range istrs {
		wstr := iexpr.ReplaceAllString(istr, "$1")
		wf, _ := strconv.ParseFloat(wstr, 32)
		weight, err := sdk.NewDecFromStr(wstr)
		if err != nil {
			return nil, err
		}
		wsum += wf

		valoper, err := sdk.AccAddressFromBech32(iexpr.ReplaceAllString(istr, "$2"))
		if err != nil {
			return nil, fmt.Errorf("invalid address for token [%v]", i)
		}

		v := &ValidatorIntent{valoper.String(), weight}
		out = append(out, v)
	}

	if wsum != 1.0 {
		return nil, fmt.Errorf("incomplete intent, sum of weights must be 1.0 (got %.2f)", wsum)
	}

	return out, nil
}

// NewMsgRequestRedemption - construct a msg to request redemption.
//nolint:interfacer
func NewMsgSignalIntent(chain_id string, intents []*ValidatorIntent, from_address sdk.Address) *MsgSignalIntent {
	return &MsgSignalIntent{ChainId: chain_id, Intents: intents, FromAddress: from_address.String()}
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
