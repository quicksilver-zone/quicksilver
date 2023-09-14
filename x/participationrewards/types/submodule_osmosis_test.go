package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func TestOsmosisParamsProtocolData_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainID   string
		BaseChain string
		BaseDenom string
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
			"missing-fields",
			fields{
				ChainID: "test-01",
			},
			true,
		},
		{
			"value",
			fields{
				ChainID:   "test-01",
				BaseDenom: "uosmo",
				BaseChain: "test-01",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oppd := types.OsmosisParamsProtocolData{
				ChainID:   tt.fields.ChainID,
				BaseDenom: tt.fields.BaseDenom,
				BaseChain: tt.fields.BaseChain,
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
