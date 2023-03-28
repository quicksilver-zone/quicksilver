package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClaim_ValidateBasic(t *testing.T) {
	type fields struct {
		UserAddress string
		ChainID     string
		Amount      uint64
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
			"invalid_address",
			fields{
				UserAddress: "cosmos1234567890",
				ChainID:     "testzone-1",
				Amount:      10000,
			},
			true,
		},
		{
			"invalid_chain_id",
			fields{
				UserAddress: GenerateAccAddressForTest().String(),
				ChainID:     "",
				Amount:      10000,
			},
			true,
		},
		{
			"invalid_chain_id",
			fields{
				UserAddress: GenerateAccAddressForTest().String(),
				ChainID:     "",
				Amount:      10000,
			},
			true,
		},
		{
			"invalid_amount",
			fields{
				UserAddress: GenerateAccAddressForTest().String(),
				ChainID:     "testzone-1",
				Amount:      0,
			},
			true,
		},
		{
			"valid",
			fields{
				UserAddress: GenerateAccAddressForTest().String(),
				ChainID:     "testzone-1",
				Amount:      1000000,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Claim{
				UserAddress: tt.fields.UserAddress,
				ChainId:     tt.fields.ChainID,
				Amount:      tt.fields.Amount,
			}
			err := c.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
