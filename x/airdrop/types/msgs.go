package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

// airdrop message types
const (
	TypeMsgClaim = "claim"
)

var _ sdk.Msg = &MsgClaim{}

// NewMsgClaim constructs a msg to claim from a zone airdrop.
func NewMsgClaim(chainID string, action int64, fromAddress sdk.Address) *MsgClaim {
	return &MsgClaim{ChainId: chainID, Action: action, Address: fromAddress.String()}
}

// Route implements Msg.
func (msg MsgClaim) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgClaim) Type() string { return TypeMsgClaim }

// ValidateBasic implements Msg.
func (msg MsgClaim) ValidateBasic() error {
	errors := make(map[string]error)

	if len(msg.ChainId) == 0 {
		errors["ChainId"] = ErrUndefinedAttribute
	}

	action := int(msg.Action)
	if action < 0 || action >= len(Action_value) {
		errors["Action"] = fmt.Errorf("%w, got %d", ErrActionOutOfBounds, msg.Action)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		errors["Address"] = err
	}

	for i, p := range msg.Proofs {
		pLabel := fmt.Sprintf("Proof [%d]:", i)
		if len(p.Key) == 0 {
			errors[pLabel+" Key"] = ErrUndefinedAttribute
		}

		if len(p.Data) == 0 {
			errors[pLabel+" Data"] = ErrUndefinedAttribute
		}

		if p.ProofOps == nil {
			errors[pLabel+" ProofOps"] = ErrUndefinedAttribute
		}

		if p.Height < 0 {
			errors[pLabel+" Height"] = ErrNegativeAttribute
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
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
