package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

// participationrewars message types
const (
	TypeMsgSubmitClaim = "submitclaim"
)

var (
	_ sdk.Msg            = &MsgSubmitClaim{}
	_ legacytx.LegacyMsg = &MsgSubmitClaim{}
)

// NewMsgSubmitClaim - construct a msg to submit a claim.
func NewMsgSubmitClaim(
	userAddress sdk.Address,
	srcZone string,
	zone string,
	claimType cmtypes.ClaimType,
	proofs []*cmtypes.Proof,
) *MsgSubmitClaim {
	return &MsgSubmitClaim{
		UserAddress: userAddress.String(),
		SrcZone:     srcZone,
		Zone:        zone,
		ClaimType:   claimType,
		Proofs:      proofs,
	}
}

// GetSignBytes implements LegacyMsg.
func (msg MsgSubmitClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// Route implements LegacyMsg.
func (msg MsgSubmitClaim) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSubmitClaim) Type() string { return TypeMsgSubmitClaim }

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

	if len(msg.SrcZone) == 0 {
		errors["SrcZone"] = ErrUndefinedAttribute
	}

	ct := int(msg.ClaimType)
	if ct < 1 || ct >= len(cmtypes.ClaimType_value) {
		errors["Action"] = fmt.Errorf("%w, got %d", cmtypes.ErrClaimTypeOutOfBounds, msg.ClaimType)
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
