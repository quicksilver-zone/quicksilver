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
	srcZone string,
	zone string,
	claimType ClaimType,
	proofs []*Proof,
) *MsgSubmitClaim {

	return &MsgSubmitClaim{
		UserAddress: userAddress.String(),
		SrcZone:     srcZone,
		Zone:        zone,
		ClaimType:   claimType,
		Proofs:      proofs,
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

	ct := int(msg.ClaimType)
	if ct < 1 || ct >= len(ClaimType_value) {
		errors["Action"] = fmt.Errorf("%w, got %d", ErrClaimTypeOutOfBounds, msg.ClaimType)
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

	if len(p.ProofType) == 0 {
		errors["ProofType"] = ErrUndefinedAttribute
	}

	// check for errors and return
	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
