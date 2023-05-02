package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func TestOsmosisParamsProtocolData_ValidateBasic(t *testing.T) {
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
				"test-01",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oppd := types.OsmosisParamsProtocolData{
				ChainID: tt.fields.ChainID,
			}
			err := oppd.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
