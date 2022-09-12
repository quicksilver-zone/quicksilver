package types

import (
	"fmt"

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

	if len(msg.Proofs) == 0 {
		errors["Proofs"] = ErrUndefinedAttribute
	}

	if len(msg.Proofs) > 0 {
		for i, p := range msg.Proofs {
			pLabel := fmt.Sprintf("Proof [%d]:", i)
			if err := p.ValidateBasic(); err != nil {
				errors[pLabel] = err
			}
		}
	}

	// check for errors and return
	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

func (p Proof) ValidateBasic() error {
	errors := make(map[string]error)

	if len(p.Key) == 0 {
		errors["Key"] = ErrUndefinedAttribute
	}

	if len(p.Data) == 0 {
		errors["Data"] = ErrUndefinedAttribute
	}

	if p.ProofOps == nil {
		errors["ProofOps"] = ErrUndefinedAttribute
	}

	if p.Height < 0 {
		errors["Height"] = ErrNegativeAttribute
	}

	// check for errors and return
	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
