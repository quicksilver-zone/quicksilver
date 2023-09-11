package types

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ingenuity-build/multierror"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

// interchainstaking message types.
const (
	TypeMsgRequestRedemption = "requestredemption"
	TypeMsgSignalIntent      = "signalintent"
)

var (
	_ sdk.Msg            = &MsgRequestRedemption{}
	_ sdk.Msg            = &MsgSignalIntent{}
	_ legacytx.LegacyMsg = &MsgRequestRedemption{}
	_ legacytx.LegacyMsg = &MsgSignalIntent{}
)

// NewMsgRequestRedemption - construct a msg to request redemption.
func NewMsgRequestRedemption(value sdk.Coin, destinationAddress string, fromAddress sdk.Address) *MsgRequestRedemption {
	return &MsgRequestRedemption{Value: value, DestinationAddress: destinationAddress, FromAddress: fromAddress.String()}
}

// Route Implements Msg.
func (msg MsgRequestRedemption) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRequestRedemption) Type() string { return TypeMsgRequestRedemption }

// ValidateBasic Implements Msg.
func (msg MsgRequestRedemption) ValidateBasic() error {
	errs := make(map[string]error)

	// check from address
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
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

//----------------------------------------------------------------

// IntentsFromString parses and validates the given string into a slice
// containing pointers to ValidatorIntent.
//
// The combined weights must be 1.0 and the valoper addresses must be valid
// bech32 strings.
//
// Tokens are comma separated, e.g.
// "0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz".
func IntentsFromString(input string) ([]*ValidatorIntent, error) {
	// regexp defining what an intent looks like:
	// {value}{address}
	iexpr := regexp.MustCompile(`(\d\.\d+)(.+1\w+)`)
	// regexp defining what the intent string looks like:
	// {intent}(,{intent})...
	pexpr := regexp.MustCompile(fmt.Sprintf("^%s(,%s)*$", iexpr.String(), iexpr.String()))
	if !pexpr.MatchString(input) {
		return nil, errors.New("invalid intents string")
	}

	out := []*ValidatorIntent{}

	wsum := sdk.ZeroDec()
	istrs := strings.Split(input, ",")
	for i, istr := range istrs {
		wstr := iexpr.ReplaceAllString(istr, "$1")
		weight, err := sdk.NewDecFromStr(wstr)
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

	if !wsum.Equal(sdk.OneDec()) {
		return nil, errors.New("combined weight must be 1.0")
	}

	return out, nil
}

// NewMsgSignalIntent - construct a msg to update signalled intent.
func NewMsgSignalIntent(chainID, intents string, fromAddress sdk.Address) *MsgSignalIntent {
	return &MsgSignalIntent{ChainId: chainID, Intents: intents, FromAddress: fromAddress.String()}
}

// Route Implements Msg.
func (msg MsgSignalIntent) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSignalIntent) Type() string { return TypeMsgSignalIntent }

// ValidateBasic Implements Msg.
func (msg MsgSignalIntent) ValidateBasic() error {
	errm := make(map[string]error)
	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		errm["FromAddress"] = err
	}

	if msg.ChainId == "" {
		errm["ChainID"] = errors.New("undefined")
	}

	wantSum := sdk.OneDec()
	weightSum := sdk.NewDec(0)
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

// NewMsgGovCloseChannel - construct a msg to update signalled intent.
func NewMsgGovCloseChannel(channelID, portName string, fromAddress sdk.Address) *MsgGovCloseChannel {
	return &MsgGovCloseChannel{ChannelId: channelID, PortId: portName, Authority: fromAddress.String()}
}

// GetSignBytes Implements Msg.
func (msg MsgGovCloseChannel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgGovCloseChannel) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

// check channel id is correct format. validate port name?
func (msg MsgGovCloseChannel) ValidateBasic() error { return nil }

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
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

// check channel id is correct format. validate port name?
func (msg MsgGovReopenChannel) ValidateBasic() error {
	// validate the zone exists, and the format is valid (e.g. quickgaia-1.delegate)
	parts := strings.Split(msg.PortId, ".")
	if len(parts) != 2 {
		return errors.New("invalid port format")
	}

	if parts[1] != "delegate" && parts[1] != "deposit" && parts[1] != "performance" && parts[1] != "withdrawal" {
		return errors.New("invalid port format; unexpected account")
	}

	if len(msg.ConnectionId) < 12 {
		return errors.New("invalid connection string; too short")
	}

	if msg.ConnectionId[0:11] != "connection-" {
		return errors.New("invalid connection string; incorrect prefix")
	}

	return nil
}
