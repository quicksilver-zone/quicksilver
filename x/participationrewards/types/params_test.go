package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParams_Validate(t *testing.T) {
	type fields struct {
		DistributionProportions DistributionProportions
		ClaimsEnabled           bool
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
		// currently params struct only contains DistributionProportions,
		// thus its testcases are to provide sufficient coverage.
		{
			"valid",
			fields{
				DistributionProportions: DistributionProportions{
					ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.34"),
					HoldingsAllocation:           sdk.MustNewDecFromStr("0.33"),
					LockupAllocation:             sdk.MustNewDecFromStr("0.33"),
				},
				ClaimsEnabled: false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				DistributionProportions: tt.fields.DistributionProportions,
				ClaimsEnabled:           tt.fields.ClaimsEnabled,
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

func TestParams(t *testing.T) {
	// test default params
	testParams := Params{
		DistributionProportions: DistributionProportions{
			ValidatorSelectionAllocation: sdk.MustNewDecFromStr("0.34"),
			HoldingsAllocation:           sdk.MustNewDecFromStr("0.33"),
			LockupAllocation:             sdk.MustNewDecFromStr("0.33"),
		},
		ClaimsEnabled: false,
	}
	defaultParams := DefaultParams()
	require.Equal(t, defaultParams, testParams)

	str := `distributionproportions:
  validatorselectionallocation: "0.340000000000000000"
  holdingsallocation: "0.330000000000000000"
  lockupallocation: "0.330000000000000000"
claimsenabled: false
`
	require.Equal(t, str, testParams.String())
}
