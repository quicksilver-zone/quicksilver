package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		{
			"blank",
			fields{},
			true,
		},
		// currently params sturct only contains DistributionProportions,
		// thus its testcases are to provide sufficient coverage.
		{
			"valid",
			fields{
				DistributionProportions: DistributionProportions{
					ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.34"),
					HoldingsAllocation:           sdk.MustNewDecFromStr("0.33"),
					LockupAllocation:             sdk.MustNewDecFromStr("0.33"),
				},
			},
			false,
		},
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
