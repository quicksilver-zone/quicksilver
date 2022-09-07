package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams_Validate(t *testing.T) {
	type fields struct {
		DistributionProportions DistributionProportions
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
			p := Params{
				DistributionProportions: tt.fields.DistributionProportions,
			}
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
