package types

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/ingenuity-build/multierror"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	"github.com/quicksilver-zone/quicksilver/v7/utils/addressutils"
)

// interchainstaking message types.
const (
	TypeMsgRequestRedemption = "requestredemption"
	TypeMsgSignalIntent      = "signalintent"
)

var (
	_ sdk.Msg            = &MsgRequestRedemption{}
	_ sdk.Msg            = &MsgSignalIntent{}
	_ sdk.Msg            = &MsgGovCloseChannel{}
	_ sdk.Msg            = &MsgGovReopenChannel{}
	_ sdk.Msg            = &MsgGovSetLsmCaps{}
	_ legacytx.LegacyMsg = &MsgRequestRedemption{}
	_ legacytx.LegacyMsg = &MsgSignalIntent{}
)

// NewMsgRequestRedemption - construct a msg to request redemption.
func NewMsgRequestRedemption(value sdk.Coin, destinationAddress string, fromAddress sdk.Address) *MsgRequestRedemption {
	return &MsgRequestRedemption{Value: value, DestinationAddress: destinationAddress, FromAddress: fromAddress.String()}
}

// Route Implements Msg.
func (MsgRequestRedemption) Route() string { return RouterKey }

// Type Implements Msg.
func (MsgRequestRedemption) Type() string { return TypeMsgRequestRedemption }

// ValidateBasic Implements Msg.
func (msg MsgRequestRedemption) ValidateBasic() error {
	errs := make(map[string]error)

	// check from address
	_, err := addressutils.AccAddressFromBech32(msg.FromAddress, "")
	if err != nil {
		errs["FromAddress"] = err
	}

	// check coin
	if msg.Value.IsNil() || msg.Value.Amount.IsNil() {
		errs["Value"] = ErrCoinAmountNil
	} else if err = msg.Value.Validate(); err != nil {
		errs["Value"] = err
	} else if msg.Value.IsZero() {
		errs["Value"] = errors.New("cannot redeem zero-value coins")
	}

	// validate recipient address
	if msg.DestinationAddress == "" {
		errs["DestinationAddress"] = errors.New("recipient address not provided")
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRequestRedemption) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgRequestRedemption) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

// ----------------------------------------------------------------

var (
	// regexp defining what an intent looks like:
	// {value}{address}
	iexpr = regexp.MustCompile(`(\d\.\d+)(.+1\w+)`)
	// regexp defining what the intent string looks like:
	// {intent}(,{intent})...
	pexpr = regexp.MustCompile(fmt.Sprintf("^%s(,%s)*$", iexpr.String(), iexpr.String()))
)

// IntentsFromString parses and validates the given string into a slice
// containing pointers to ValidatorIntent.
//
// The combined weights must be 1.0 and the valoper addresses must be valid
// bech32 strings.
//
// Tokens are comma separated, e.g.
// "0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz".
func IntentsFromString(input string) ([]*ValidatorIntent, error) {
	if !pexpr.MatchString(input) {
		return nil, errors.New("invalid intents string")
	}

	out := []*ValidatorIntent{}

	wsum := sdkmath.LegacyZeroDec()
	istrs := strings.Split(input, ",")
	for i, istr := range istrs {
		wstr := iexpr.ReplaceAllString(istr, "$1")
		weight, err := sdkmath.LegacyNewDecFromStr(wstr)
		if err != nil {
			return nil, fmt.Errorf("intent token [%v]: %w", i, err)
		}

		if !weight.IsPositive() {
			return nil, fmt.Errorf("intent token [%v]: must not be negative nor zero", i)
		}

		wsum = wsum.Add(weight)

		v := &ValidatorIntent{
			iexpr.ReplaceAllString(istr, "$2"),
			weight,
		}
		out = append(out, v)
	}

	if !wsum.Equal(sdkmath.LegacyOneDec()) {
		return nil, errors.New("combined weight must be 1.0")
	}

	return out, nil
}

// NewMsgSignalIntent - construct a msg to update signalled intent.
func NewMsgSignalIntent(chainID, intents string, fromAddress sdk.Address) *MsgSignalIntent {
	return &MsgSignalIntent{ChainId: chainID, Intents: intents, FromAddress: fromAddress.String()}
}

// Route Implements Msg.
func (MsgSignalIntent) Route() string { return RouterKey }

// Type Implements Msg.
func (MsgSignalIntent) Type() string { return TypeMsgSignalIntent }

// ValidateBasic Implements Msg.
func (msg MsgSignalIntent) ValidateBasic() error {
	errm := make(map[string]error)
	if _, err := addressutils.AccAddressFromBech32(msg.FromAddress, ""); err != nil {
		errm["FromAddress"] = err
	}

	if msg.ChainId == "" {
		errm["ChainID"] = errors.New("undefined")
	}

	wantSum := sdkmath.LegacyOneDec()
	weightSum := sdkmath.LegacyNewDec(0)
	intents, err := IntentsFromString(msg.Intents)
	if err != nil {
		errm["Intents"] = err
	} else {
		for i, intent := range intents {
			if _, _, err := bech32.DecodeAndConvert(intent.ValoperAddress); err != nil {
				istr := fmt.Sprintf("Intent_%02d_ValoperAddress", i)
				errm[istr] = err
			}

			if intent.Weight.GT(wantSum) {
				istr := fmt.Sprintf("Intent_%02d_Weight", i)
				errm[istr] = fmt.Errorf("weight %d overruns maximum of %v", intent.Weight, wantSum)
			}
			weightSum = weightSum.Add(intent.Weight)
		}

		if !weightSum.Equal(wantSum) {
			errm["IntentWeights"] = fmt.Errorf("sum of weights is %v, not %v", weightSum, wantSum)
		}
	}
	if len(errm) > 0 {
		return multierror.New(errm)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSignalIntent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSignalIntent) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

// NewMsgGovCloseChannel

// GetSignBytes Implements Msg.
func (msg MsgGovCloseChannel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgGovCloseChannel) GetSigners() []sdk.AccAddress {
	fromAddress, _ := addressutils.AccAddressFromBech32(msg.Authority, "")
	return []sdk.AccAddress{fromAddress}
}

func (msg MsgGovCloseChannel) ValidateBasic() error {
	_, err := addressutils.AccAddressFromBech32(msg.Authority, "")
	if err != nil {
		return err
	}

	if err := ValidatePort(msg.PortId); err != nil {
		return err
	}

	return ValidateChannel(msg.ChannelId)
}

// NewMsgGovReopenChannel - construct a msg to update signalled intent.
func NewMsgGovReopenChannel(connectionID, portName string, fromAddress sdk.Address) *MsgGovReopenChannel {
	return &MsgGovReopenChannel{ConnectionId: connectionID, PortId: portName, Authority: fromAddress.String()}
}

// GetSignBytes Implements Msg.
func (msg MsgGovReopenChannel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgGovReopenChannel) GetSigners() []sdk.AccAddress {
	fromAddress, _ := addressutils.AccAddressFromBech32(msg.Authority, "")
	return []sdk.AccAddress{fromAddress}
}

// ValidateBasic
func (msg MsgGovReopenChannel) ValidateBasic() error {
	_, err := addressutils.AccAddressFromBech32(msg.Authority, "")
	if err != nil {
		return err
	}

	if err := ValidatePort(msg.PortId); err != nil {
		return err
	}

	return ValidateConnection(msg.ConnectionId)
}

// MsgGovSetLsmCaps

// GetSignBytes Implements Msg.
func (msg MsgGovSetLsmCaps) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgGovSetLsmCaps) GetSigners() []sdk.AccAddress {
	fromAddress, _ := addressutils.AccAddressFromBech32(msg.Authority, "")
	return []sdk.AccAddress{fromAddress}
}

// ValidateBasic
func (msg MsgGovSetLsmCaps) ValidateBasic() error {
	_, err := addressutils.AccAddressFromBech32(msg.Authority, "")
	if err != nil {
		return err
	}

	if len(msg.ChainId) == 0 || len(msg.ChainId) > 100 {
		return errors.New("invalid chain id")
	}

	return msg.Caps.Validate()
}

// Helpers
func ValidateConnection(connectionID string) error {
	if !strings.HasPrefix(connectionID, "connection-") {
		return errors.New("invalid connection")
	}

	_, err := strconv.ParseInt(connectionID[11:], 0, 64)

	return err
}

func ValidateChannel(channelID string) error {
	if !strings.HasPrefix(channelID, "channel-") {
		return errors.New("invalid channel")
	}

	_, err := strconv.ParseInt(channelID[8:], 0, 64)

	return err
}

func ValidatePort(portID string) error {
	// remove leading prefix icacontroller- if passed in msg
	portID = strings.ReplaceAll(portID, "icacontroller-", "")

	// validate the zone exists, and the format is valid (e.g. quickgaia-1.delegate)
	parts := strings.Split(portID, ".")

	if len(parts) != 2 {
		return errors.New("invalid port format")
	}

	if parts[1] != "delegate" && parts[1] != "deposit" && parts[1] != "performance" && parts[1] != "withdrawal" {
		return errors.New("invalid port format; unexpected account")
	}

	return nil
}

func (caps LsmCaps) Validate() error {
	if caps.GlobalCap.GT(sdkmath.LegacyOneDec()) || caps.GlobalCap.LT(sdkmath.LegacyZeroDec()) {
		return errors.New("global cap must be between 0 and 1")
	}

	if caps.ValidatorCap.GT(sdkmath.LegacyOneDec()) || caps.ValidatorCap.LT(sdkmath.LegacyZeroDec()) {
		return errors.New("validator cap must be between 0 and 1")
	}

	if caps.ValidatorBondCap.LTE(sdkmath.LegacyZeroDec()) {
		return errors.New("validator bond cap must be greater than 0")
	}

	return nil
}
