package types

import (
	"encoding/json"
	"errors"
	"fmt"

	sdkioerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ingenuity-build/multierror"
)

// airdrop message types.

var (
	_ sdk.Msg = &MsgClaim{}
	_ sdk.Msg = &MsgIncentivePoolSpend{}
	_ sdk.Msg = &MsgRegisterZoneDrop{}
)

// NewMsgClaim constructs a msg to claim from a zone airdrop.
func NewMsgClaim(chainID string, action int64, fromAddress sdk.Address) *MsgClaim {
	return &MsgClaim{ChainId: chainID, Action: action, Address: fromAddress.String()}
}

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

//////////////////////////////////////////////////////////////////////////////////////////////////////

// NewMsgIncentivePoolSpend constructs a msg to claim from a zone airdrop.
func NewMsgIncentivePoolSpend(authority, toAddress sdk.Address, amt sdk.Coins) *MsgIncentivePoolSpend {
	return &MsgIncentivePoolSpend{
		Authority: authority.String(),
		ToAddress: toAddress.String(),
		Amount:    amt,
	}
}

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
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// ValidateBasic implements Msg.
func (msg MsgRegisterZoneDrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if msg.ZoneDrop == nil {
		return errors.New("proposal must contain a valid ZoneDrop")
	}

	if len(msg.ClaimRecords) == 0 {
		return errors.New("update must contain valid ClaimRecords")
	}

	if err := msg.ZoneDrop.ValidateBasic(); err != nil {
		return err
	}

	// decompress claim records
	crsb, err := Decompress(msg.ClaimRecords)
	if err != nil {
		return err
	}

	// unmarshal json
	var crs ClaimRecords
	if err := json.Unmarshal(crsb, &crs); err != nil {
		return err
	}

	sumMax := uint64(0)
	// validate ClaimRecords and process
	for i, cr := range crs {
		if err := cr.ValidateBasic(); err != nil {
			return fmt.Errorf("claim record %d, %w", i, err)
		}
		if len(cr.ActionsCompleted) != 0 {
			return fmt.Errorf("invalid zonedrop proposal claim record [%d]: contains completed actions", i)
		}

		if cr.ChainId != msg.ZoneDrop.ChainId {
			return fmt.Errorf("invalid zonedrop proposal claim record [%d]: chainID missmatch, expected %q got %q",
				i,
				msg.ZoneDrop.ChainId,
				cr.ChainId,
			)
		}

		sumMax += cr.MaxAllocation
	}

	// check allocations
	if sumMax > msg.ZoneDrop.Allocation {
		return fmt.Errorf("sum of claim records max allocations (%v) exceed zone airdrop allocation (%v)", sumMax, msg.ZoneDrop.Allocation)
	}

	return nil
}

// GetSignBytes implements Msg.
func (msg MsgRegisterZoneDrop) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgRegisterZoneDrop) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}
