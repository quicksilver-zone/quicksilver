package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestMsgRegisterZone_ValidateBasic(t *testing.T) {
	testAddress := addressutils.GenerateAccAddressForTest().String()

	type fields struct {
		Authority        string
		ConnectionID     string
		BaseDenom        string
		LocalDenom       string
		AccountPrefix    string
		ReturnToSender   bool
		UnbondingEnabled bool
		Deposits         bool
		LiquidityModule  bool
		Decimals         int64
		MessagesPerTx    int64
		Is118            bool
		SubzoneInfo      *types.SubzoneInfo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test valid",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         6,
				MessagesPerTx:    5,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: false,
		},
		{
			name: "invalid connection field",
			fields: fields{
				ConnectionID:     "test",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         6,
				MessagesPerTx:    5,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid basedenom field",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "0",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         6,
				MessagesPerTx:    5,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid localdenom field",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "0",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         6,
				MessagesPerTx:    5,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid account prefix",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "a",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         6,
				MessagesPerTx:    5,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
		{
			name: "liquidity",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  true,
				Decimals:         6,
				MessagesPerTx:    5,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid decimals",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         0,
				MessagesPerTx:    5,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid-messages-per-tx-0",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         0,
				MessagesPerTx:    0,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid-messages-per-tx-negative",
			fields: fields{
				ConnectionID:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				Deposits:         false,
				LiquidityModule:  false,
				Decimals:         0,
				MessagesPerTx:    -1,
				Is118:            true,
				Authority:        testAddress,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := types.MsgRegisterZone{
				ConnectionID:     tt.fields.ConnectionID,
				BaseDenom:        tt.fields.BaseDenom,
				LocalDenom:       tt.fields.LocalDenom,
				AccountPrefix:    tt.fields.AccountPrefix,
				LiquidityModule:  tt.fields.LiquidityModule,
				MessagesPerTx:    tt.fields.MessagesPerTx,
				ReturnToSender:   tt.fields.ReturnToSender,
				DepositsEnabled:  tt.fields.Deposits,
				UnbondingEnabled: tt.fields.UnbondingEnabled,
				Decimals:         tt.fields.Decimals,
				Is_118:           tt.fields.Is118,
				SubzoneInfo:      tt.fields.SubzoneInfo,
				Authority:        tt.fields.Authority,
			}

			err := m.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateZone_ValidateBasic(t *testing.T) {
	testAddress := addressutils.GenerateAccAddressForTest().String()
	validZoneID := "zone-1"

	type fields struct {
		Authority string
		ZoneID    string
		Changes   []*types.UpdateZoneValue
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test valid",
			fields: fields{
				Authority: testAddress,
				ZoneID:    validZoneID,
				Changes: []*types.UpdateZoneValue{
					{
						Key:   types.UpdateZoneKeyBaseDenom,
						Value: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid authority",
			fields: fields{
				Authority: "",
				ZoneID:    validZoneID,
				Changes: []*types.UpdateZoneValue{
					{
						Key:   types.UpdateZoneKeyBaseDenom,
						Value: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid zone ID empty",
			fields: fields{
				Authority: testAddress,
				ZoneID:    "",
				Changes: []*types.UpdateZoneValue{
					{
						Key:   types.UpdateZoneKeyBaseDenom,
						Value: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid update Key",
			fields: fields{
				Authority: testAddress,
				ZoneID:    validZoneID,
				Changes: []*types.UpdateZoneValue{
					{
						Key:   "invalid",
						Value: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid no updates",
			fields: fields{
				Authority: testAddress,
				ZoneID:    validZoneID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := types.MsgUpdateZone{
				Authority: tt.fields.Authority,
				ZoneID:    tt.fields.ZoneID,
				Changes:   tt.fields.Changes,
			}

			err := m.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

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

func TestMsgReopenIntent_ValidateBasic(t *testing.T) {
	validAddress := addressutils.GenerateAccAddressForTest()

	type fields struct {
		PortID       string
		ConnectionID string
		Authority    string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errorMsg string
	}{
		{
			"blank",
			fields{
				Authority: validAddress.String(),
			},
			true,
			"invalid port format",
		},
		{
			"invalid authority ",
			fields{
				Authority: "invalid",
			},
			true,
			"invalid authority address",
		},
		{
			"invalid port",
			fields{
				PortID:    "cat",
				Authority: validAddress.String(),
			},
			true,
			"invalid port format",
		},
		{
			"invalid account",
			fields{
				PortID:    "icacontroller-osmosis-4.bad",
				Authority: validAddress.String(),
			},
			true,
			"invalid port format; unexpected account",
		},
		{
			"invalid connection; too short",
			fields{
				PortID:       "icacontroller-osmosis-4.withdrawal",
				ConnectionID: "bad-1",
				Authority:    validAddress.String(),
			},
			true,
			"invalid connection string; too short",
		},
		{
			"invalid connection; too short",
			fields{
				PortID:       "icacontroller-osmosis-4.withdrawal",
				ConnectionID: "longenoughbutstillbad-1",
				Authority:    validAddress.String(),
			},
			true,
			"invalid connection string; incorrect prefix",
		},
		{
			"valid",
			fields{
				PortID:       "icacontroller-osmosis-4.withdrawal",
				ConnectionID: "connection-1",
				Authority:    validAddress.String(),
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
				Authority:    tt.fields.Authority,
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
