package types

import (
	sdkioerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/multierror"
)

var _ sdk.Msg = &MsgUpdateParams{}

// ValidateBasic performs stateless validation for Proof.
func (p *Proof) ValidateBasic() error {
	errs := make(map[string]error)

	if len(p.Key) == 0 {
		errs["Key"] = ErrUndefinedAttribute
	}

	if len(p.Data) == 0 {
		errs["Data"] = ErrUndefinedAttribute
	}

	if p.ProofOps == nil {
		errs["ProofOps"] = ErrUndefinedAttribute
	}

	if p.Height < 0 {
		errs["Height"] = ErrNegativeAttribute
	}

	if p.ProofType == "" {
		errs["ProofType"] = ErrUndefinedAttribute
	}

	// check for errors and return
	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// GetSignBytes implements the LegacyMsg interface.
func (m *MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkioerrors.Wrap(err, "invalid authority address")
	}

	return m.Params.Validate()
}
