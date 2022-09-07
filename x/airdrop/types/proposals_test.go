package types

import (
	"testing"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRegisterZoneDropProposal_ValidateBasic(t *testing.T) {
	type fields struct {
		Title        string
		Description  string
		ZoneDrop     *ZoneDrop
		ClaimRecords []byte
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
		{
			"invalid-00",
			fields{
				Title:        "Flashdrop",
				Description:  "An airdrop that is valid for one hour only",
				ZoneDrop:     nil,
				ClaimRecords: nil,
			},
			true,
		},
		{
			"invalid-01",
			fields{
				Title:        "Flashdrop",
				Description:  "An airdrop that is valid for one hour only",
				ZoneDrop:     &ZoneDrop{},
				ClaimRecords: []byte{},
			},
			true,
		},
		// HandleRegisterZoneDropProposal will deal with in depth validation,
		// as ClaimRecords is compressed data that needs to be decompressed.
		{
			"valid",
			fields{
				Title:       "Flashdrop",
				Description: "An airdrop that is valid for one hour only",
				ZoneDrop: &ZoneDrop{
					ChainId:    "test-1",
					StartTime:  time.Now().Add(-time.Hour),
					Duration:   time.Hour,
					Decay:      30 * time.Minute,
					Allocation: 16400,
					Actions: []sdk.Dec{
						sdk.MustNewDecFromStr("0.1"),
						sdk.MustNewDecFromStr("0.2"),
						sdk.MustNewDecFromStr("0.3"),
						sdk.MustNewDecFromStr("0.4"),
					},
					IsConcluded: false,
				},
				ClaimRecords: []byte{0},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := RegisterZoneDropProposal{
				Title:        tt.fields.Title,
				Description:  tt.fields.Description,
				ZoneDrop:     tt.fields.ZoneDrop,
				ClaimRecords: tt.fields.ClaimRecords,
			}

			err := m.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("Error:\n%v\n", err)
				return
			}
			require.NoError(t, err)
		})
	}
}
