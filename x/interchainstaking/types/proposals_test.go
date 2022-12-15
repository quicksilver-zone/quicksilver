package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterZoneProposal_ValidateBasic(t *testing.T) {
	type fields struct {
		Title           string
		Description     string
		ConnectionId    string
		BaseDenom       string
		LocalDenom      string
		AccountPrefix   string
		MultiSend       bool
		LiquidityModule bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := RegisterZoneProposal{
				Title:           tt.fields.Title,
				Description:     tt.fields.Description,
				ConnectionId:    tt.fields.ConnectionId,
				BaseDenom:       tt.fields.BaseDenom,
				LocalDenom:      tt.fields.LocalDenom,
				AccountPrefix:   tt.fields.AccountPrefix,
				MultiSend:       tt.fields.MultiSend,
				LiquidityModule: tt.fields.LiquidityModule,
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
