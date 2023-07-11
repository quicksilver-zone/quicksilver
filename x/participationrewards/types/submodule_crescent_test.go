package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrescentParamsProtocolData_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainID string
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
			"valid",
			fields{
				ChainID: "test-01",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uppd := CrescentParamsProtocolData{
				ChainID: tt.fields.ChainID,
			}
			err := uppd.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
