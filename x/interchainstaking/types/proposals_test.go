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
		MessagesPerTx   int64
	}

	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errorMsg string
	}{
		{
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
		},
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
				MessagesPerTx:   tt.fields.MessagesPerTx,
			}

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
