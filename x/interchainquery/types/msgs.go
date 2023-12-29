package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// interchainquery message types.
const (
	TypeMsgSubmitQueryResponse = "submitqueryresponse"
)

var _ sdk.Msg = &MsgSubmitQueryResponse{}

// Route Implements Msg.
func (MsgSubmitQueryResponse) Route() string { return RouterKey }

// Type Implements Msg.
func (MsgSubmitQueryResponse) Type() string { return TypeMsgSubmitQueryResponse }

// ValidateBasic Implements Msg.
func (msg MsgSubmitQueryResponse) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return err
	}

	if msg.Height < 0 {
		return errors.New("height must be non-negative")
	}

	if len(msg.ChainId) < 2 {
		return errors.New("invalid chain id")
	}

	if len(msg.QueryId) != 64 {
		return errors.New("invalid query id")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSubmitQueryResponse) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSubmitQueryResponse) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}
