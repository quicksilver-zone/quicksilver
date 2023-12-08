package types_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func TestRegisterZoneProposal_ValidateBasic(t *testing.T) {
	type fields struct {
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
		Is_118           bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
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
				Is_118:           true,
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
				Is_118:           true,
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
				Is_118:           true,
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
				Is_118:           true,
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
				Is_118:           true,
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
				Is_118:           true,
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
				Is_118:           true,
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
				Is_118:           true,
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
				Is_118:           true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				tt.fields.Is_118,
			)

			err := m.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
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
