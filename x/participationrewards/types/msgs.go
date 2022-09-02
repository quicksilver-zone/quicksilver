package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

var _ sdk.Msg = &MsgSubmitClaim{}

// NewMsgSubmitClaim - construct a msg to submit a claim.
func NewMsgSubmitClaim(
	userAddress sdk.Address,
	zone string,
) *MsgSubmitClaim {
	return &MsgSubmitClaim{
		UserAddress: userAddress.String(),
		Zone:        zone,
	}
}

// GetSigners implements Msg.
func (msg MsgSubmitClaim) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.UserAddress)
	return []sdk.AccAddress{fromAddress}
}

// ValidateBasic implements Msg: stateless checks.
func (msg MsgSubmitClaim) ValidateBasic() error {
	errors := make(map[string]error)
	if _, err := sdk.AccAddressFromBech32(msg.UserAddress); err != nil {
		errors["UserAddress"] = err
	}

	if len(msg.Zone) == 0 {
		errors["Zone"] = ErrUndefinedAttribute
	}

	if len(msg.Key) == 0 {
		errors["Key"] = ErrUndefinedAttribute
	}

	if len(msg.Data) == 0 {
		errors["Data"] = ErrUndefinedAttribute
	}

	if len(msg.ProofOps) == 0 {
		errors["ProofOps"] = ErrUndefinedAttribute
	}

	if msg.Height == 0 {
		errors["ProofOps"] = ErrUndefinedAttribute
	}

	if len(msg.Data) != len(msg.Key) {
		errors["DataLength"] = ErrSliceLengthMismatch
	}

	if len(msg.ProofOps) != len(msg.Key) {
		errors["DataLength"] = ErrSliceLengthMismatch
	}

	// check for errors and return
	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
