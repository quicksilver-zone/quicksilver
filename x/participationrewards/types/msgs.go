package types

import (
	"encoding/hex"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/multierror"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

var (
	_ sdk.Msg = &MsgSubmitClaim{}
	_ sdk.Msg = &MsgGovRemoveProtocolData{}
	_ sdk.Msg = &MsgAddProtocolData{}
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

// GetSigners implements Msg.
func (msg MsgSubmitClaim) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.UserAddress)
	return []sdk.AccAddress{fromAddress}
}

// ValidateBasic implements Msg: stateless checks.
func (msg MsgSubmitClaim) ValidateBasic() error {
	errs := make(map[string]error)
	if _, err := sdk.AccAddressFromBech32(msg.UserAddress); err != nil {
		errs["UserAddress"] = err
	}

	if msg.Zone == "" {
		errs["Zone"] = ErrUndefinedAttribute
	}

	if msg.SrcZone == "" {
		errs["SrcZone"] = ErrUndefinedAttribute
	}

	ct := int(msg.ClaimType)
	if ct < 1 || ct >= len(cmtypes.ClaimType_value) {
		errs["Action"] = fmt.Errorf("%w, got %d", cmtypes.ErrClaimTypeOutOfBounds, msg.ClaimType)
	}

	if len(msg.Proofs) == 0 {
		errs["Proofs"] = ErrUndefinedAttribute
	}

	if len(msg.Proofs) > 0 {
		for i, p := range msg.Proofs {
			err := p.ValidateBasic()
			if err == nil {
				return nil
			}

			pLabel := fmt.Sprintf("Proof [%s]", hex.EncodeToString(p.Key))
			if _, ok := errs[pLabel]; ok {
				pLabel += fmt.Sprintf("-%d", i)
			}
			errs[pLabel+":"] = err
		}
	}

	// check for errors and return
	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

// NewMsgGovRemoveProtocolData - construct a governance proposal msg to remove protocoldata by key.
func NewMsgGovRemoveProtocolData(key string, fromAddress sdk.Address) *MsgGovRemoveProtocolData {
	return &MsgGovRemoveProtocolData{Key: key, Authority: fromAddress.String()}
}

// GetSignBytes Implements Msg.
func (msg MsgGovRemoveProtocolData) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgGovRemoveProtocolData) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

// Validate.
func (msg MsgGovRemoveProtocolData) ValidateBasic() error {
	// check title is non-empty
	if msg.Title == "" {
		return errors.New("title must not be empty")
	}

	// check description is non-empty
	if msg.Description == "" {
		return errors.New("description must not be empty")
	}

	// check key is non-empty
	if msg.Key == "" {
		return errors.New("key must not be empty")
	}

	// check authority is non-empty
	if msg.Authority == "" {
		return errors.New("authority must not be empty")
	}

	// check authority bech32 is valid
	_, err := addressutils.AddressFromBech32(msg.Authority, "")
	return err
}

// GetSigners Implements Msg.
func (msg MsgAddProtocolData) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

// Validate Implements Msg.
func (msg MsgAddProtocolData) ValidateBasic() error {
	pdtv, exists := ProtocolDataType_value[msg.Type]
	if !exists {
		return ErrUnknownProtocolDataType
	}

	pd, err := UnmarshalProtocolData(ProtocolDataType(pdtv), msg.Data)
	if err != nil {
		return err
	}

	if err := pd.ValidateBasic(); err != nil {
		return err
	}

	_, err = addressutils.AddressFromBech32(msg.Authority, "")
	return err
}
