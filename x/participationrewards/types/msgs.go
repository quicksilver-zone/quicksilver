package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSubmitClaim{}

// NewMsgSubmitClaim - construct a msg to submit a claim.
//nolint:interfacer
func NewMsgSubmitClaim(
	user_address sdk.Address,
	zone string,
	asset_type string,
) *MsgSubmitClaim {
	return &MsgSubmitClaim{
		UserAddress: user_address.String(),
		Zone:        zone,
		AssetType:   asset_type,
	}
}

// GetSigners Implements Msg.
func (msg MsgSubmitClaim) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.UserAddress)
	return []sdk.AccAddress{fromAddress}
}

// ValidateBasic Implements Msg.
func (msg MsgSubmitClaim) ValidateBasic() error {
	errors := make(map[string]error)
	if _, err := sdk.AccAddressFromBech32(msg.UserAddress); err != nil {
		errors["UserAddress"] = err
	}

	// TODO: check for valid zone (chain_id)

	// TODO: check for valid asset type (sdk.Coin)

	return nil
}
