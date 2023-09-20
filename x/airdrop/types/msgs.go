package types

import (
	"fmt"

	"github.com/ingenuity-build/multierror"

	sdkioerrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

// airdrop message types.

const (
	TypeMsgClaim              = "claim"
	TypeMsgIncentivePoolSpend = "incentive-pool-spend"
)

var (
	_ sdk.Msg            = &MsgClaim{}
	_ legacytx.LegacyMsg = &MsgClaim{}
)

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
	errs := make(map[string]error)

	if msg.ChainId == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	action := int(msg.Action)
	if action < 1 || action >= len(Action_value) {
		errs["Action"] = fmt.Errorf("%w, got %d", ErrActionOutOfBounds, msg.Action)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		errs["Address"] = err
	}

	for i, p := range msg.Proofs {
		pLabel := fmt.Sprintf("Proof [%d]:", i)
		if len(p.Key) == 0 {
			errs[pLabel+" Key"] = ErrUndefinedAttribute
		}

		if len(p.Data) == 0 {
			errs[pLabel+" Data"] = ErrUndefinedAttribute
		}

		if p.ProofOps == nil {
			errs[pLabel+" ProofOps"] = ErrUndefinedAttribute
		}

		if p.Height < 0 {
			errs[pLabel+" Height"] = ErrNegativeAttribute
		}
	}

	if len(errs) > 0 {
		return multierror.New(errs)
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

// NewMsgIncentivePoolSpend constructs a msg to claim from a zone airdrop.
func NewMsgIncentivePoolSpend(authority, toAddress sdk.Address, amt sdk.Coins) *MsgIncentivePoolSpend {
	return &MsgIncentivePoolSpend{
		Authority: authority.String(),
		ToAddress: toAddress.String(),
		Amount:    amt,
	}
}

// Route implements Msg.
func (msg MsgIncentivePoolSpend) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgIncentivePoolSpend) Type() string { return TypeMsgClaim }

// ValidateBasic implements Msg.
func (msg MsgIncentivePoolSpend) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	if from.Equals(to) {
		return sdkerrors.ErrInvalidAddress.Wrapf("to and from addresses equal: %s", err)
	}

	if !msg.Amount.IsValid() {
		return sdkioerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkioerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes implements Msg.
func (msg MsgIncentivePoolSpend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgIncentivePoolSpend) GetSigners() []sdk.AccAddress {
	address, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{address}
}
