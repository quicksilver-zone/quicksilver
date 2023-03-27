package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
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

	fromAddr := (sdk.AccAddress)([]byte{0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc})
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
		(sdk.AccAddress)([]byte{0x84, 0xbf, 0xf8, 0x4c, 0x7d, 0xda, 0xd1, 0x1c, 0xb8, 0xc0, 0x73, 0x86, 0xe9, 0x19, 0x28, 0xc5, 0x67, 0x5c, 0xa4, 0xbc}))
	sigIntent.FromAddress = "abcdefghi"
	err = sigIntent.ValidateBasic()
	require.NotNil(t, err, "expecting a non-nil error")
	require.Contains(t, err.Error(), "invalid separator index")
	require.Contains(t, err.Error(), "undefined")
}

func TestMsgSignalIntent_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainId     string
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
				ChainId:     tt.fields.ChainId,
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
				DestinationAddress: utils.GenerateAccAddressForTest().String(),
				FromAddress:        utils.GenerateAccAddressForTest().String(),
			},
			true,
		},
		{
			"valid_zero",
			fields{
				Value: sdk.Coin{
					Denom:  "stake",
					Amount: math.ZeroInt(),
				},
				DestinationAddress: utils.GenerateAccAddressForTest().String(),
				FromAddress:        utils.GenerateAccAddressForTest().String(),
			},
			false,
		},
		{
			"valid_value",
			fields{
				Value: sdk.Coin{
					Denom:  "stake",
					Amount: math.ZeroInt(),
				},
				DestinationAddress: utils.GenerateAccAddressForTest().String(),
				FromAddress:        utils.GenerateAccAddressForTest().String(),
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
