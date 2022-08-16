package types

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

// interchainstaking message types
const (
	TypeMsgRegisterZone      = "registerzone"
	TypeMsgRequestRedemption = "requestredemption"
	TypeMsgSignalIntent      = "signalintent"
)

// NewMsgRequestRedemption - construct a msg to request redemption.
func NewMsgRequestRedemption(coin string, destinationAddress string, fromAddress sdk.Address) *MsgRequestRedemption {
	return &MsgRequestRedemption{Coin: coin, DestinationAddress: destinationAddress, FromAddress: fromAddress.String()}
}

// Route Implements Msg.
func (msg MsgRequestRedemption) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRequestRedemption) Type() string { return TypeMsgRegisterZone }

// ValidateBasic Implements Msg.
func (msg MsgRequestRedemption) ValidateBasic() error {
	// TODO: check from address
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return err
	}

	// check coin
	coin, err := sdk.ParseCoinNormalized(msg.Coin)
	if err != nil {
		return err
	}

	if !coin.IsPositive() {
		return fmt.Errorf("expected positive value, got %v", msg.Coin)
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

// IntentsFromString validates and parses the given string into a slice
// containing pointers to ValidatorIntent.
//
// The combined weights must be 1.0 and the valoper addresses must be valid
// bech32 strings. (what about zero weights?)
//
// Tokens are comma separated, e.g.
// "0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz".
func IntentsFromString(input string) ([]*ValidatorIntent, error) {
	iexpr := regexp.MustCompile(`(\d.\d+)(.+1\w+)`)
	pexpr := regexp.MustCompile(fmt.Sprintf("^%s(,%s)*$", iexpr.String(), iexpr.String()))
	if !pexpr.MatchString(input) {
		return nil, fmt.Errorf("invalid intents string")
	}

	out := []*ValidatorIntent{}

	istrs := strings.Split(input, ",")
	for i, istr := range istrs {
		wstr := iexpr.ReplaceAllString(istr, "$1")
		weight, err := sdk.NewDecFromStr(wstr)
		if err != nil {
			return nil, fmt.Errorf("intent token [%v]: %w", i, err)
		}

		v := &ValidatorIntent{
			iexpr.ReplaceAllString(istr, "$2"),
			weight,
		}
		out = append(out, v)
	}

	return out, nil
}

// NewMsgRequestRedemption - construct a msg to request redemption.
func NewMsgSignalIntent(chainID string, intents []*ValidatorIntent, fromAddress sdk.Address) *MsgSignalIntent {
	return &MsgSignalIntent{ChainId: chainID, Intents: intents, FromAddress: fromAddress.String()}
}

// Route Implements Msg.
func (msg MsgSignalIntent) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSignalIntent) Type() string { return TypeMsgSignalIntent }

// ValidateBasic Implements Msg.
func (msg MsgSignalIntent) ValidateBasic() error {
	errors := make(map[string]error)
	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		errors["FromAddress"] = err
	}

	// TODO: check for valid chainID

	wantSum := sdk.MustNewDecFromStr("1.0")
	weightSum := sdk.NewDec(0)
	for i, intent := range msg.Intents {
		if _, _, err := bech32.DecodeAndConvert(intent.ValoperAddress); err != nil {
			istr := fmt.Sprintf("Intent_%02d_ValoperAddress", i)
			errors[istr] = err
		}

		if intent.Weight.GT(wantSum) {
			istr := fmt.Sprintf("Intent_%02d_Weight", i)
			errors[istr] = fmt.Errorf("weight %d overruns maximum of %v", intent.Weight, wantSum)
		}
		weightSum = weightSum.Add(intent.Weight)
	}

	if !weightSum.Equal(wantSum) {
		errors["IntentWeights"] = fmt.Errorf("sum of weights is %v, not %v", weightSum, wantSum)
	}

	if len(errors) > 0 {
		return multierror.New(errors)
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
