package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterZoneProposal_ValidateBasic(t *testing.T) {
	type fields struct {
		Title            string
		Description      string
		ConnectionId     string
		BaseDenom        string
		LocalDenom       string
		AccountPrefix    string
		ReturnToSender   bool
		UnbondingEnabled bool
		LiquidityModule  bool
		Decimals         int64
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
				ConnectionId:     "connection-0",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				LiquidityModule:  false,
				Decimals:         6,
			},
			wantErr: false,
		},
		{
			name: "invalid connection field",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
				ConnectionId:     "test",
				BaseDenom:        "uatom",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				LiquidityModule:  false,
				Decimals:         6,
			},
			wantErr: true,
		},
		{
			name: "invalid basdenom field",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
				ConnectionId:     "test",
				BaseDenom:        "0",
				LocalDenom:       "uqatom",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				LiquidityModule:  false,
				Decimals:         6,
			},
			wantErr: true,
		},
		{
			name: "invalid localdenom field",
			fields: fields{
				Title:            "Enable testzone-1",
				Description:      "onboard testzone-1",
				ConnectionId:     "test",
				BaseDenom:        "uatom",
				LocalDenom:       "0",
				AccountPrefix:    "cosmos",
				ReturnToSender:   false,
				UnbondingEnabled: false,
				LiquidityModule:  false,
				Decimals:         6,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := RegisterZoneProposal{
				Title:            tt.fields.Title,
				Description:      tt.fields.Description,
				ConnectionId:     tt.fields.ConnectionId,
				BaseDenom:        tt.fields.BaseDenom,
				LocalDenom:       tt.fields.LocalDenom,
				AccountPrefix:    tt.fields.AccountPrefix,
				ReturnToSender:   tt.fields.ReturnToSender,
				UnbondingEnabled: tt.fields.UnbondingEnabled,
				LiquidityModule:  tt.fields.LiquidityModule,
				Decimals:         tt.fields.Decimals,
			}

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
	uzp := &UpdateZoneProposal{
		Title:       "Testing right here",
		Description: "Testing description",
		ChainId:     "quicksilver",
		Changes: []*UpdateZoneValue{
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
	sink = (interface{})(nil)
}
