package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
<<<<<<< HEAD
=======
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
>>>>>>> origin/develop
)

func TestGenesisState(t *testing.T) {
	// test default genesis state
	testGenesisState := types.GenesisState{}
	defaultGenesisState := types.DefaultGenesisState()
	require.Equal(t, *defaultGenesisState, testGenesisState)
	// test new genesis state
	newGenesisState := types.NewGenesisState(
		types.Params{},
	)
	testGenesisState = types.GenesisState{}
	require.Equal(t, *newGenesisState, testGenesisState)
}

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Params types.Params
		Claims []*types.Claim
	}
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
		{
			"invalid_claim",
			fields{
				Params: types.DefaultParams(),
				Claims: []*types.Claim{
					{
						UserAddress: utils.GenerateAccAddressForTest().String(),
						ChainId:     "testzone-1",
						Amount:      0,
					},
				},
			},
			true,
		},
		{
			"valid_claim",
			fields{
				Params: types.DefaultParams(),
				Claims: []*types.Claim{
					{
						UserAddress: utils.GenerateAccAddressForTest().String(),
						ChainId:     "testzone-1",
						Amount:      1000000,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := types.GenesisState{
				Params: tt.fields.Params,
				Claims: tt.fields.Claims,
			}

			err := gs.Validate()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
