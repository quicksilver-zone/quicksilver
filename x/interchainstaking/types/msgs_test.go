package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func TestIntentsFromString(t *testing.T) {
	// 1. Ensure we can properly parse intents with their weights.
	intents := "0.3cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0,0.3cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf,0.4cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll"

	wantIntents := []*types.ValidatorIntent{
		{
			ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			Weight:         sdk.MustNewDecFromStr("0.3"),
		},
		{
			ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			Weight:         sdk.MustNewDecFromStr("0.3"),
		},
		{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			Weight:         sdk.MustNewDecFromStr("0.4"),
		},
	}
	intentsSlice, err := types.IntentsFromString(intents)
	require.NoError(t, err)
	require.Equal(t, wantIntents, intentsSlice, "intents mismatch")

	// 2. Ensure that if the weights don't add up to 1.0 that it fails.
	// 2.1. Greater than 1.0
	malIntents, err := types.IntentsFromString("1.3cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0,2.3cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf,3.4cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll")
	require.Nil(t, malIntents, "expecting nil intents")
	require.NotNil(t, err, "expecting a non-nil error")
	require.Contains(t, err.Error(), "combined weight must be 1.0")

	// 2.2. Less than 1.0
	malIntents, err = types.IntentsFromString("0.03cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0,0.3cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf,0.2cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll")
	require.Nil(t, malIntents, "expecting nil intents")
	require.NotNil(t, err, "expecting a non-nil error")
	require.Contains(t, err.Error(), "combined weight must be 1.0")

	// 3. Invalid intent strings not matching: (\d.\d)(.+1\w+) are rejected
	malIntents, err = types.IntentsFromString("foo,bar,few")
	require.Nil(t, malIntents, "expecting nil intents")
	require.NotNil(t, err, "expecting a non-nil error")
	require.Contains(t, err.Error(), "invalid intents string")

	fromAddr := sdk.AccAddress([]byte{0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc})
	sigIntent := types.NewMsgSignalIntent("quicksilver", intents, fromAddr)
	err = sigIntent.ValidateBasic()
	if err != nil {
		t.Fatal(err.Error())
	}
	require.Nil(t, err)

	// Check the router key.
	gotRoute := sigIntent.Route()
	wantRoute := "interchainstaking"
	require.Equal(t, wantRoute, gotRoute, "mismatch in route")

	// Check the type.
	gotType := sigIntent.Type()
	wantType := "signalintent"
	require.Equal(t, wantType, gotType, "mismatch in type")

	// Check the signBytes.
	signBytes := sigIntent.GetSignBytes()
	require.True(t, len(signBytes) != 0, "expecting signBytes to be produced")

	// Signers should return the from address.
	gotSigners := sigIntent.GetSigners()
	wantSigners := []sdk.AccAddress{fromAddr}
	require.Equal(t, wantSigners, gotSigners, "mismatch in signers")
}

func TestIntentsFromStringInvalidValoperAddressesFailsOnValidate(t *testing.T) {
	negativeIntents, err := types.IntentsFromString("0.5cosmosvaloper1sjllsnramtg7ewxqwwrwjxfgc4n4ef9u2lcnp0,-0.5cosmosvaloper156g8f9837p7d4c46p8yt3rlals9c5vuurfrrzf")
	require.Nil(t, negativeIntents, "expecting nil intents")
	require.NotNil(t, err, "expecting a non-nil error")
	require.Contains(t, err.Error(), "must not be negative")

	// The valoper addresses have invalid checksums, but they'll only be caught on invoking .ValidateBasic()
	intents := "1.7cosmosvaloper1sjllsnramtg7ewxqwwrwjxfgc4n4ef9u2lcnp0,-0.5cosmosvaloper156g8f9837p7d4c46p8yt3rlals9c5vuurfrrzf"

	sigIntent := types.NewMsgSignalIntent("", intents,
		sdk.AccAddress([]byte{0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc}))
	sigIntent.FromAddress = "abcdefghi"
	err = sigIntent.ValidateBasic()
	require.NotNil(t, err, "expecting a non-nil error")
	require.Contains(t, err.Error(), "invalid separator index")
	require.Contains(t, err.Error(), "undefined")
}

func TestMsgSignalIntent_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainID     string
		Intents     string
		FromAddress string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"blank",
			fields{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.MsgSignalIntent{
				ChainId:     tt.fields.ChainID,
				Intents:     tt.fields.Intents,
				FromAddress: tt.fields.FromAddress,
			}
			err := msg.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgRequestRedemption_ValidateBasic(t *testing.T) {
	type fields struct {
		Value              sdk.Coin
		DestinationAddress string
		FromAddress        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"nil",
			fields{},
			true,
		},
		{
			"empty_coin",
			fields{
				Value: sdk.Coin{
					Denom: "stake",
				},
				DestinationAddress: addressutils.GenerateAccAddressForTest().String(),
				FromAddress:        addressutils.GenerateAccAddressForTest().String(),
			},
			true,
		},
		{
			"invalid_nil_destination_address",
			fields{
				Value: sdk.Coin{
					Denom:  "stake",
					Amount: sdkmath.OneInt(),
				},
				DestinationAddress: "",
				FromAddress:        addressutils.GenerateAccAddressForTest().String(),
			},
			true,
		},
		{
			"invalid_nil_from_address",
			fields{
				Value: sdk.Coin{
					Denom:  "stake",
					Amount: sdkmath.OneInt(),
				},
				DestinationAddress: addressutils.GenerateAccAddressForTest().String(),
				FromAddress:        "",
			},
			true,
		},
		{
			"invalid_zero",
			fields{
				Value: sdk.Coin{
					Denom:  "stake",
					Amount: sdkmath.ZeroInt(),
				},
				DestinationAddress: addressutils.GenerateAccAddressForTest().String(),
				FromAddress:        addressutils.GenerateAccAddressForTest().String(),
			},
			true,
		},
		{
			"invalid_negative",
			fields{
				Value: sdk.Coin{
					Denom:  "stake",
					Amount: sdkmath.NewInt(-1),
				},
				DestinationAddress: addressutils.GenerateAccAddressForTest().String(),
				FromAddress:        addressutils.GenerateAccAddressForTest().String(),
			},
			true,
		},
		{
			"valid_value",
			fields{
				Value: sdk.Coin{
					Denom:  "stake",
					Amount: sdkmath.OneInt(),
				},
				DestinationAddress: addressutils.GenerateAccAddressForTest().String(),
				FromAddress:        addressutils.GenerateAccAddressForTest().String(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.MsgRequestRedemption{
				Value:              tt.fields.Value,
				DestinationAddress: tt.fields.DestinationAddress,
				FromAddress:        tt.fields.FromAddress,
			}
			err := msg.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgReopen_ValidateBasic(t *testing.T) {
	type fields struct {
		PortID       string
		ConnectionID string
		FromAddress  string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errorMsg string
	}{
		{
			"blank",
			fields{},
			true,
			"empty address string is not allowed",
		},
		{
			"invalid port",
			fields{
				FromAddress: addressutils.GenerateAddressForTestWithPrefix("cosmos"),
				PortID:      "cat",
			},
			true,
			"invalid port format",
		},
		{
			"invalid account",
			fields{
				FromAddress: addressutils.GenerateAddressForTestWithPrefix("cosmos"),
				PortID:      "icacontroller-osmosis-4.bad",
			},
			true,
			"invalid port format; unexpected account",
		},
		{
			"invalid connection; too short",
			fields{
				FromAddress:  addressutils.GenerateAddressForTestWithPrefix("cosmos"),
				PortID:       "icacontroller-osmosis-4.withdrawal",
				ConnectionID: "bad-1",
			},
			true,
			"invalid connection",
		},
		{
			"invalid connection; too short",
			fields{
				FromAddress:  addressutils.GenerateAddressForTestWithPrefix("cosmos"),
				PortID:       "icacontroller-osmosis-4.withdrawal",
				ConnectionID: "longenoughbutstillbad-1",
			},
			true,
			"invalid connection",
		},
		{
			"valid",
			fields{
				FromAddress:  addressutils.GenerateAddressForTestWithPrefix("cosmos"),
				PortID:       "icacontroller-osmosis-4.withdrawal",
				ConnectionID: "connection-1",
			},
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.MsgGovReopenChannel{
				PortId:       tt.fields.PortID,
				ConnectionId: tt.fields.ConnectionID,
				Authority:    tt.fields.FromAddress,
			}
			err := msg.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				require.ErrorContains(t, err, tt.errorMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

// Given a valid connectionID starting with 'connection-', the function should return nil
func TestValidateConnection_ValidConnectionID_ReturnsNil(t *testing.T) {
	connectionID := "connection-123"

	err := types.ValidateConnection(connectionID)

	require.NoError(t, err)
}

// Given a connectionID with the minimum length of 12 characters, the function should return nil
func TestValidateConnection_MinimumLengthConnectionID_ReturnsNil(t *testing.T) {
	connectionID := "connection-1"

	err := types.ValidateConnection(connectionID)

	require.NoError(t, err)
}

// Given an empty string, the function should return an error
func TestValidateConnection_EmptyString_ReturnsError(t *testing.T) {
	connectionID := ""

	err := types.ValidateConnection(connectionID)

	require.Error(t, err)
}

// Given a connectionID starting with 'connection-' and ending with a non-numeric character, the function should return an error
func TestValidateConnection_ConnectionIDWithNonNumericCharacter_ReturnsError(t *testing.T) {
	connectionID := "connection-abc"

	err := types.ValidateConnection(connectionID)

	require.Error(t, err)
}

// Given a connectionID starting with 'connection' (without the hyphen), the function should return an error
func TestValidateConnection_ConnectionIDWithoutHyphen_ReturnsError(t *testing.T) {
	connectionID := "connection123"

	err := types.ValidateConnection(connectionID)

	require.Error(t, err)
}

// Should return nil for a valid portID in the format 'zone.account'
func TestValidatePort_ValidPortID_ReturnsNil(t *testing.T) {
	portID := "quickgaia-1.delegate"

	err := types.ValidatePort(portID)

	require.Nil(t, err)
}

// Should remove the 'icacontroller-' prefix from the portID before validation
func TestValidatePort_RemovePrefix_ReturnsNil(t *testing.T) {
	portID := "icacontroller-quickgaia-1.delegate"

	err := types.ValidatePort(portID)

	require.Nil(t, err)
}

// Should accept 'delegate', 'deposit', 'performance', and 'withdrawal' as valid account types
func TestValidatePort_ValidAccountTypes_ReturnsNil(t *testing.T) {
	accountTypes := []string{"delegate", "deposit", "performance", "withdrawal"}

	for _, accountType := range accountTypes {
		portID := fmt.Sprintf("quickgaia-1.%s", accountType)

		err := types.ValidatePort(portID)

		require.Nil(t, err)
	}
}

// Should return an error for an invalid portID format
func TestValidatePort_InvalidFormat_ReturnsError(t *testing.T) {
	portID := "invalidformat"

	err := types.ValidatePort(portID)

	require.Error(t, err)
}

// Should return an error for a portID with more than one dot separator
func TestValidatePort_MultipleDotSeparators_ReturnsError(t *testing.T) {
	portID := "quickgaia-1.delegate.extra"

	err := types.ValidatePort(portID)

	require.Error(t, err)
}

// Should return an error for a portID with an unexpected account type
func TestValidatePort_UnexpectedAccountType_ReturnsError(t *testing.T) {
	portID := "quickgaia-1.unexpected"

	err := types.ValidatePort(portID)

	require.Error(t, err)
}

// Given a valid channel ID starting with "channel-", the function should return no error.
func TestValidateChannel_ValidChannelID_ReturnsNoError(t *testing.T) {
	channelID := "channel-123"
	err := types.ValidateChannel(channelID)
	require.NoError(t, err)
}

// Given a valid channel ID starting with "channel-" and followed by a number, the function should return no error.
func TestValidateChannel_ValidChannelIDWithNumber_ReturnsNoError(t *testing.T) {
	channelID := "channel-456"
	err := types.ValidateChannel(channelID)
	require.NoError(t, err)
}

// Given a channel ID starting with "channel-" and followed by a number greater than 0, the function should return no error.
func TestValidateChannel_ValidChannelIDWithPositiveNumber_ReturnsNoError(t *testing.T) {
	channelID := "channel-789"
	err := types.ValidateChannel(channelID)
	require.NoError(t, err)
}

// Given an empty string as channel ID, the function should return an error.
func TestValidateChannel_EmptyChannelID_ReturnsError(t *testing.T) {
	channelID := ""
	err := types.ValidateChannel(channelID)
	require.Error(t, err)
}

// Given a channel ID that does not start with "channel-", the function should return an error.
func TestValidateChannel_ChannelIDWithoutPrefix_ReturnsError(t *testing.T) {
	channelID := "invalid-channel"
	err := types.ValidateChannel(channelID)
	require.Error(t, err)
}

// Given a channel ID that starts with "channel" (without the trailing "-"), the function should return an error.
func TestValidateChannel_ChannelIDWithoutTrailingDash_ReturnsError(t *testing.T) {
	channelID := "channel123"
	err := types.ValidateChannel(channelID)
	require.Error(t, err)
}

// test the string "channel-"
func TestValidateChannel_OnlyPrefix_ReturnsError(t *testing.T) {
	channelID := "channel-"
	err := types.ValidateChannel(channelID)
	require.Error(t, err)
}

func TestCloseChannelValidateBasic(t *testing.T) {
	cases := []struct {
		Name string
		Msg  types.MsgGovCloseChannel
		Err  string
	}{
		{
			Name: "valid",
			Msg:  types.MsgGovCloseChannel{Title: "test", Description: "test", ChannelId: "channel-1", PortId: "icacontroller-juno-1.delegate", Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "",
		},
		{
			Name: "invalid channel",
			Msg:  types.MsgGovCloseChannel{Title: "test", Description: "test", ChannelId: "cat", PortId: "juno-1.delegate", Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "invalid channel",
		},
		{
			Name: "invalid port",
			Msg:  types.MsgGovCloseChannel{Title: "test", Description: "test", ChannelId: "channel-1", PortId: "icacontroller-juno-1.bad", Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "invalid port",
		},
		{
			Name: "invalid authority",
			Msg:  types.MsgGovCloseChannel{Title: "test", Description: "test", ChannelId: "channel-1", PortId: "icacontroller-juno-1.delegate", Authority: "aaa"},
			Err:  "decoding bech32 failed: invalid bech32 string length 3",
		},
	}

	for _, c := range cases {
		err := c.Msg.ValidateBasic()
		if c.Err == "" { // happy
			require.NoError(t, err, c.Name)
		} else {
			require.ErrorContains(t, err, c.Err, c.Name)
		}
	}
}

func TestReopenChannelValidateBasic(t *testing.T) {
	cases := []struct {
		Name string
		Msg  types.MsgGovReopenChannel
		Err  string
	}{
		{
			Name: "valid",
			Msg:  types.MsgGovReopenChannel{Title: "test", Description: "test", ConnectionId: "connection-1", PortId: "icacontroller-juno-1.delegate", Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "",
		},
		{
			Name: "invalid connection",
			Msg:  types.MsgGovReopenChannel{Title: "test", Description: "test", ConnectionId: "cat", PortId: "juno-1.delegate", Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "invalid connection",
		},
		{
			Name: "invalid port",
			Msg:  types.MsgGovReopenChannel{Title: "test", Description: "test", ConnectionId: "connection-1", PortId: "icacontroller-juno-1.bad", Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "invalid port",
		},
		{
			Name: "invalid authority",
			Msg:  types.MsgGovReopenChannel{Title: "test", Description: "test", ConnectionId: "connection-1", PortId: "icacontroller-juno-1.delegate", Authority: "aaa"},
			Err:  "decoding bech32 failed: invalid bech32 string length 3",
		},
	}

	for _, c := range cases {
		err := c.Msg.ValidateBasic()
		if c.Err == "" { // happy
			require.NoError(t, err, c.Name)
		} else {
			require.ErrorContains(t, err, c.Err, c.Name)
		}
	}
}

func TestGovSetLsmCaps(t *testing.T) {
	cases := []struct {
		Name string
		Msg  types.MsgGovSetLsmCaps
		Err  string
	}{
		{
			Name: "valid",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", ChainId: "chain-1", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec(), ValidatorBondCap: sdk.NewDec(250), GlobalCap: sdk.NewDecWithPrec(50, 2)}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "",
		},
		{
			Name: "invalid empty chain id",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec(), ValidatorBondCap: sdk.NewDec(250), GlobalCap: sdk.NewDecWithPrec(50, 2)}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "invalid chain id",
		},
		{
			Name: "invalid bad authority",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec(), ValidatorBondCap: sdk.NewDec(250), GlobalCap: sdk.NewDecWithPrec(50, 2)}, Authority: "raa"},
			Err:  "decoding bech32 failed: invalid bech32 string length 3",
		},
		{
			Name: "invalid validator cap < 0",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", ChainId: "chain-1", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec().Neg(), ValidatorBondCap: sdk.NewDec(250), GlobalCap: sdk.NewDecWithPrec(50, 2)}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "validator cap must be between 0 and 1",
		},
		{
			Name: "invalid validator cap > 1",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", ChainId: "chain-1", Caps: &types.LsmCaps{ValidatorCap: sdk.NewDec(250), ValidatorBondCap: sdk.NewDec(250), GlobalCap: sdk.NewDecWithPrec(50, 2)}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "validator cap must be between 0 and 1",
		},
		{
			Name: "invalid negative bond cap",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", ChainId: "chain-1", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec(), ValidatorBondCap: sdk.NewDec(250).Neg(), GlobalCap: sdk.NewDecWithPrec(50, 2)}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "validator bond cap must be greater than 0",
		},
		{
			Name: "invalid zero bond cap",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", ChainId: "chain-1", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec(), ValidatorBondCap: sdk.ZeroDec(), GlobalCap: sdk.NewDecWithPrec(50, 2)}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "validator bond cap must be greater than 0",
		},
		{
			Name: "invalid - global cap > 1",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", ChainId: "chain-1", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec(), ValidatorBondCap: sdk.NewDec(250), GlobalCap: sdk.NewDecWithPrec(50, 1)}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "global cap must be between 0 and 1",
		},
		{
			Name: "invalid - global cap < 0",
			Msg:  types.MsgGovSetLsmCaps{Title: "test", Description: "test", ChainId: "chain-1", Caps: &types.LsmCaps{ValidatorCap: sdkmath.LegacyOneDec(), ValidatorBondCap: sdk.NewDec(250), GlobalCap: sdk.NewDecWithPrec(50, 2).Neg()}, Authority: addressutils.GenerateAddressForTestWithPrefix("quick")},
			Err:  "global cap must be between 0 and 1",
		},
	}

	for _, c := range cases {
		err := c.Msg.ValidateBasic()
		if c.Err == "" { // happy
			require.NoError(t, err, c.Name)
		} else {
			require.ErrorContains(t, err, c.Err, c.Name)
		}
	}
}
