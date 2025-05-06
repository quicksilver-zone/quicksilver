package types

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ingenuity-build/multierror"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

// interchainstaking message types.
const (
	TypeMsgRequestRedemption          = "requestredemption"
	TypeMsgCancelRedemption           = "cancelredemption"
	TypeMsgRequeueRedemption          = "requeueredemption"
	TypeMsgUpdateRedemption           = "updateredemption"
	TypeMsgSignalIntent               = "signalintent"
	TypeMsgGovCloseChannel            = "govclosechannel"
	TypeMsgGovReopenChannel           = "govreopenchannel"
	TypeMsgGovSetLsmCaps              = "govsetlsmcaps"
	TypeMsgGovAddValidatorDenyList    = "govaddvalidatordenylist"
	TypeMsgGovRemoveValidatorDenyList = "govremovevalidatordenylist"
	TypeMsgGovExecuteICATx            = "govexecuteicatx"
)

var (
	_ sdk.Msg = &MsgRequestRedemption{}
	_ sdk.Msg = &MsgCancelRedemption{}
	_ sdk.Msg = &MsgRequeueRedemption{}
	_ sdk.Msg = &MsgSignalIntent{}
	_ sdk.Msg = &MsgGovCloseChannel{}
	_ sdk.Msg = &MsgGovReopenChannel{}
	_ sdk.Msg = &MsgGovSetLsmCaps{}
	_ sdk.Msg = &MsgGovAddValidatorDenyList{}
	_ sdk.Msg = &MsgGovRemoveValidatorDenyList{}
	_ sdk.Msg = &MsgGovExecuteICATx{}

	_ legacytx.LegacyMsg = &MsgRequestRedemption{}
	_ legacytx.LegacyMsg = &MsgCancelRedemption{}
	_ legacytx.LegacyMsg = &MsgRequeueRedemption{}
	_ legacytx.LegacyMsg = &MsgUpdateRedemption{}
	_ legacytx.LegacyMsg = &MsgSignalIntent{}
	_ legacytx.LegacyMsg = &MsgGovCloseChannel{}
	_ legacytx.LegacyMsg = &MsgGovReopenChannel{}
	_ legacytx.LegacyMsg = &MsgGovSetLsmCaps{}
	_ legacytx.LegacyMsg = &MsgGovAddValidatorDenyList{}
	_ legacytx.LegacyMsg = &MsgGovRemoveValidatorDenyList{}
	_ legacytx.LegacyMsg = &MsgGovExecuteICATx{}

	_ codectypes.UnpackInterfacesMessage = &MsgGovExecuteICATx{}
)

// NewMsgRequestRedemption - construct a msg to request redemption.
func NewMsgRequestRedemption(value sdk.Coin, destinationAddress string, fromAddress sdk.Address) *MsgRequestRedemption {
	return &MsgRequestRedemption{Value: value, DestinationAddress: destinationAddress, FromAddress: fromAddress.String()}
}

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

var hexpr = regexp.MustCompile("^[A-Fa-f0-9]{64}$")

// NewMsgCancelQueuedRedemption - construct a msg to cancel a requested redemption.
func NewMsgCancelRedemption(chainID string, hash string, fromAddress sdk.Address) *MsgCancelRedemption {
	return &MsgCancelRedemption{ChainId: chainID, Hash: hash, FromAddress: fromAddress.String()}
}

// ValidateBasic Implements Msg.
func (msg MsgCancelRedemption) ValidateBasic() error {
	errs := make(map[string]error)

	// check from address
	_, err := addressutils.AccAddressFromBech32(msg.FromAddress, "")
	if err != nil {
		errs["FromAddress"] = err
	}

	// check hash
	if !hexpr.MatchString(msg.Hash) {
		errs["Hash"] = fmt.Errorf("invalid sha256 hash - expecting 64 character hex string, got %s", msg.Hash)
	}

	// validate recipient address
	if msg.ChainId == "" {
		errs["ChainId"] = errors.New("chainId not provided")
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCancelRedemption) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgCancelRedemption) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

// ---------------------------------------------------------------

func NewMsgRequeueRedemption(chainID string, hash string, fromAddress sdk.Address) *MsgRequeueRedemption {
	return &MsgRequeueRedemption{ChainId: chainID, Hash: hash, FromAddress: fromAddress.String()}
}

// ValidateBasic Implements Msg.
func (msg MsgRequeueRedemption) ValidateBasic() error {
	errs := make(map[string]error)

	// check from address
	_, err := addressutils.AccAddressFromBech32(msg.FromAddress, "")
	if err != nil {
		errs["FromAddress"] = err
	}

	// check hash
	if !hexpr.MatchString(msg.Hash) {
		errs["Hash"] = fmt.Errorf("invalid sha256 hash - expecting 64 character hex string, got %s (%d)", msg.Hash, len(msg.Hash))
	}

	// validate recipient address
	if msg.ChainId == "" {
		errs["ChainId"] = errors.New("chainId not provided")
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRequeueRedemption) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgRequeueRedemption) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

// ---------------------------------------------------------------

func NewMsgUpdateRedemption(chainID string, hash string, newStatus int32, fromAddress sdk.Address) *MsgUpdateRedemption {
	return &MsgUpdateRedemption{ChainId: chainID, Hash: hash, NewStatus: newStatus, FromAddress: fromAddress.String()}
}

// ValidateBasic Implements Msg.
func (msg MsgUpdateRedemption) ValidateBasic() error {
	errs := make(map[string]error)

	// check from address
	_, err := addressutils.AccAddressFromBech32(msg.FromAddress, "")
	if err != nil {
		errs["FromAddress"] = err
	}

	// check hash
	if !hexpr.MatchString(msg.Hash) {
		errs["Hash"] = fmt.Errorf("invalid sha256 hash - expecting 64 character hex string, got %s (%d)", msg.Hash, len(msg.Hash))
	}

	// validate recipient address
	if msg.ChainId == "" {
		errs["ChainId"] = errors.New("chainId not provided")
	}

	switch msg.NewStatus {
	case WithdrawStatusTokenize: // intentionally removed as not currently supported, but included here for completeness.
		errs["NewStatus"] = errors.New("new status WithdrawStatusTokenize not supported")
	case WithdrawStatusQueued:
	case WithdrawStatusUnbond:
	case WithdrawStatusSend: // send is not a valid state for recovery, included here for completeness.
		errs["NewStatus"] = errors.New("new status WithdrawStatusSend not supported")
	case WithdrawStatusCompleted:
	default:
		errs["NewStatus"] = errors.New("new status not provided or invalid")
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateRedemption) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgUpdateRedemption) GetSigners() []sdk.AccAddress {
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

// ValidateBasic Implements Msg.
func (msg MsgSignalIntent) ValidateBasic() error {
	errm := make(map[string]error)
	if _, err := addressutils.AccAddressFromBech32(msg.FromAddress, ""); err != nil {
		errm["FromAddress"] = err
	}

	if msg.ChainId == "" {
		errm["ChainID"] = errors.New("chainId not provided")
	}

	intents, err := IntentsFromString(msg.Intents)
	if err != nil {
		errm["Intents"] = err
	} else {
		for i, intent := range intents {
			if _, _, err := bech32.DecodeAndConvert(intent.ValoperAddress); err != nil {
				istr := fmt.Sprintf("Intent_%02d_ValoperAddress", i)
				errm[istr] = err
			}
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
	if caps.GlobalCap.GT(sdk.OneDec()) || caps.GlobalCap.LT(sdk.ZeroDec()) {
		return errors.New("global cap must be between 0 and 1")
	}

	if caps.ValidatorCap.GT(sdk.OneDec()) || caps.ValidatorCap.LT(sdk.ZeroDec()) {
		return errors.New("validator cap must be between 0 and 1")
	}

	if caps.ValidatorBondCap.LTE(sdk.ZeroDec()) {
		return errors.New("validator bond cap must be greater than 0")
	}

	return nil
}

// MsgGovAddValidatorDenyList

// // ValidateBasic
func (msg MsgGovAddValidatorDenyList) ValidateBasic() error {
	_, err := addressutils.AccAddressFromBech32(msg.Authority, "")
	if err != nil {
		return err
	}

	if _, err := addressutils.ValAddressFromBech32(msg.OperatorAddress, ""); err != nil {
		return err
	}

	return nil
}

// // GetSignBytes Implements Msg.
func (msg MsgGovAddValidatorDenyList) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// // GetSigners Implements Msg.
func (msg MsgGovAddValidatorDenyList) GetSigners() []sdk.AccAddress {
	fromAddress, _ := addressutils.AccAddressFromBech32(msg.Authority, "")
	return []sdk.AccAddress{fromAddress}
}

// MsgGovRemoveValidatorDenyList

// // ValidateBasic

func (msg MsgGovRemoveValidatorDenyList) ValidateBasic() error {
	_, err := addressutils.AccAddressFromBech32(msg.Authority, "")
	if err != nil {
		return err
	}

	if _, err := addressutils.ValAddressFromBech32(msg.OperatorAddress, ""); err != nil {
		return err
	}

	return nil
}

// // GetSignBytes Implements Msg.
func (msg MsgGovRemoveValidatorDenyList) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// // GetSigners Implements Msg.
func (msg MsgGovRemoveValidatorDenyList) GetSigners() []sdk.AccAddress {
	fromAddress, _ := addressutils.AccAddressFromBech32(msg.Authority, "")
	return []sdk.AccAddress{fromAddress}
}

// MsgGovExecuteICATx

// GetSignBytes Implements Msg.
func (msg MsgGovExecuteICATx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgGovExecuteICATx) GetSigners() []sdk.AccAddress {
	fromAddress, _ := addressutils.AccAddressFromBech32(msg.Authority, "")
	return []sdk.AccAddress{fromAddress}
}

// ValidateBasic
func (msg MsgGovExecuteICATx) ValidateBasic() error {
	if _, err := addressutils.AccAddressFromBech32(msg.Authority, ""); err != nil {
		return err
	}

	if _, err := addressutils.AccAddressFromBech32(msg.Address, ""); err != nil {
		return err
	}

	if msg.ChainId == "" {
		return errors.New("invalid chain id")
	}

	if len(msg.Msgs) == 0 {
		return errors.New("no msgs provided")
	}

	if len(msg.Msgs) > 20 {
		return errors.New("max 20 msgs are supported")
	}

	return nil
}

func (msg MsgGovExecuteICATx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return tx.UnpackInterfaces(unpacker, msg.Msgs)
}

// legacytx.LegacyMsg implementations - remove with SDk v0.50

// MsgGovExecuteICATx
func (msg MsgGovExecuteICATx) Route() string {
	return RouterKey
}

func (msg MsgGovExecuteICATx) Type() string {
	return TypeMsgGovExecuteICATx
}

// MsgSignalIntent

func (msg MsgSignalIntent) Route() string {
	return RouterKey
}

func (msg MsgSignalIntent) Type() string {
	return TypeMsgSignalIntent
}

// MsgGovCloseChannel

func (msg MsgGovCloseChannel) Route() string {
	return RouterKey
}

func (msg MsgGovCloseChannel) Type() string {
	return TypeMsgGovCloseChannel
}

// MsgGovReopenChannel

func (msg MsgGovReopenChannel) Route() string {
	return RouterKey
}

func (msg MsgGovReopenChannel) Type() string {
	return TypeMsgGovReopenChannel
}

// MsgGovSetLsmCaps

func (msg MsgGovSetLsmCaps) Route() string {
	return RouterKey
}

func (msg MsgGovSetLsmCaps) Type() string {
	return TypeMsgGovSetLsmCaps
}

// MsgGovAddValidatorDenyList

func (msg MsgGovAddValidatorDenyList) Route() string {
	return RouterKey
}

func (msg MsgGovAddValidatorDenyList) Type() string {
	return TypeMsgGovAddValidatorDenyList
}

// MsgGovRemoveValidatorDenyList

func (msg MsgGovRemoveValidatorDenyList) Route() string {
	return RouterKey
}

func (msg MsgGovRemoveValidatorDenyList) Type() string {
	return TypeMsgGovRemoveValidatorDenyList
}

// MsgRequestRedemption

func (msg MsgRequestRedemption) Route() string {
	return RouterKey
}

func (msg MsgRequestRedemption) Type() string {
	return TypeMsgRequestRedemption
}

// MsgCancelRedemption

func (msg MsgCancelRedemption) Route() string {
	return RouterKey
}

func (msg MsgCancelRedemption) Type() string {
	return TypeMsgCancelRedemption
}

// MsgRequeueRedemption

func (msg MsgRequeueRedemption) Route() string {
	return RouterKey
}

func (msg MsgRequeueRedemption) Type() string {
	return TypeMsgRequeueRedemption
}

// MsgUpdateRedemption

func (msg MsgUpdateRedemption) Route() string {
	return RouterKey
}

func (msg MsgUpdateRedemption) Type() string {
	return TypeMsgUpdateRedemption
}
