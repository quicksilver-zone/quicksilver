package types_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func TestRegisterZoneProposal_ValidateBasic(t *testing.T) {
	type fields struct {
<<<<<<< HEAD
		Title           string
		Description     string
		ConnectionId    string
		BaseDenom       string
		LocalDenom      string
		AccountPrefix   string
		MultiSend       bool
		LiquidityModule bool
		MessagesPerTx   int64
=======
		Title            string
		Description      string
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
>>>>>>> origin/develop
	}

	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errorMsg string
	}{
		{
<<<<<<< HEAD
			"zero-length-title",
			fields{
				"",
				"description",
				"connection-0",
				"uatom",
				"uqatom",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"proposal title cannot be blank: invalid proposal content",
		},
		{
			"zero-length-desc",
			fields{
				"title",
				"",
				"connection-0",
				"uatom",
				"uqatom",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"proposal description cannot be blank: invalid proposal content",
		},
		{
			"zero-length-connection",
			fields{
				"title",
				"description",
				"",
				"uatom",
				"uqatom",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"invalid length connection string: ",
		},
		{
			"invalid-connection",
			fields{
				"title",
				"description",
				"abcdefghijklmnop",
				"uatom",
				"uqatom",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"invalid connection string: abcdefghijklmnop",
		},
		{
			"invalid-length-base-denom",
			fields{
				"title",
				"description",
				"connection-0",
				"",
				"uqatom",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"invalid denom: ",
		},
		{
			"invalid-base-denom",
			fields{
				"title",
				"description",
				"connection-0",
				"000",
				"uqatom",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"invalid denom: 000",
		},
		{
			"invalid-length-local-denom",
			fields{
				"title",
				"description",
				"connection-0",
				"uatom",
				"",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"invalid denom: ",
		},
		{
			"invalid-length-local-denom",
			fields{
				"title",
				"description",
				"connection-0",
				"uatom",
				"0000",
				"cosmos",
				true,
				false,
				5,
			},
			true,
			"invalid denom: 000",
		},
		{
			"invalid-length-prefix",
			fields{
				"title",
				"description",
				"connection-0",
				"uatom",
				"uqatom",
				"a",
				true,
				false,
				5,
			},
			true,
			"account prefix must be at least 2 characters",
		},
		{
			"invalid-messages-per-tx-0",
			fields{
				"title",
				"description",
				"connection-0",
				"uatom",
				"uqatom",
				"ki",
				true,
				false,
				0,
			},
			true,
			"messages_per_tx must be a positive non-zero integer",
		},
		{
			"invalid-messages-per-tx-negative",
			fields{
				"title",
				"description",
				"connection-0",
				"uatom",
				"uqatom",
				"cosmos",
				true,
				false,
				-1,
			},
			true,
			"messages_per_tx must be a positive non-zero integer",
		},
		{
			"invalid-messages-per-tx-negative",
			fields{
				"title",
				"description",
				"connection-0",
				"uatom",
				"uqatom",
				"cosmos",
				true,
				false,
				50,
			},
			false,
			"",
=======
			name: "test valid",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: false,
		},
		{
			name: "invalid gov content",
			fields: fields{
				Title:            "",
				Description:      "",
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
			},
			wantErr: true,
		},
		{
			name: "invalid connection field",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
		},
		{
			name: "invalid basedenom field",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
		},
		{
			name: "invalid localdenom field",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
		},
		{
			name: "invalid account prefix",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
		},
		{
			name: "liquidity",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
		},
		{
			name: "invalid decimals",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
		},
		{
			name: "invalid-messages-per-tx-0",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
		},
		{
			name: "invalid-messages-per-tx-negative",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
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
			},
			wantErr: true,
>>>>>>> origin/develop
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
<<<<<<< HEAD
			m := RegisterZoneProposal{
				Title:           tt.fields.Title,
				Description:     tt.fields.Description,
				ConnectionId:    tt.fields.ConnectionId,
				BaseDenom:       tt.fields.BaseDenom,
				LocalDenom:      tt.fields.LocalDenom,
				AccountPrefix:   tt.fields.AccountPrefix,
				MultiSend:       tt.fields.MultiSend,
				LiquidityModule: tt.fields.LiquidityModule,
				MessagesPerTx:   tt.fields.MessagesPerTx,
			}
=======
			m := types.NewRegisterZoneProposal(
				tt.fields.Title,
				tt.fields.Description,
				tt.fields.ConnectionID,
				tt.fields.BaseDenom,
				tt.fields.LocalDenom,
				tt.fields.AccountPrefix,
				tt.fields.ReturnToSender,
				tt.fields.UnbondingEnabled,
				tt.fields.Deposits,
				tt.fields.LiquidityModule,
				tt.fields.Decimals,
				tt.fields.MessagesPerTx,
			)
>>>>>>> origin/develop

			err := m.ValidateBasic()
			if tt.wantErr {
				require.ErrorContains(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

var sink interface{}

func BenchmarkUpdateZoneProposalString(b *testing.B) {
	uzp := &types.UpdateZoneProposal{
		Title:       "Testing right here",
		Description: "Testing description",
		ChainId:     "quicksilver",
		Changes: []*types.UpdateZoneValue{
			{Key: "K1", Value: "V1"},
			{Key: strings.Repeat("Ks", 100), Value: strings.Repeat("Vs", 128)},
			{Key: strings.Repeat("a", 64), Value: strings.Repeat("A", 28)},
		},
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		str := uzp.String()
		b.SetBytes(int64(len(str)))
		sink = str
	}

	if sink == nil {
		b.Fatal("Benchmark did not run")
	}
	sink = interface{}(nil)
}
