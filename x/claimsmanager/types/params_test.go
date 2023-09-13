package types_test

import (
	"testing"

	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/stretchr/testify/require"
)

func TestParams_Validate(t *testing.T) {
	type fields struct{}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"blank",
			fields{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := types.Params{}
			err := p.Validate()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestParams(t *testing.T) {
	// test default params
	testParams := types.Params{}
	defaultParams := types.DefaultParams()
	require.Equal(t, defaultParams, testParams)

	str := `{}
`
	require.Equal(t, str, testParams.String())
}
