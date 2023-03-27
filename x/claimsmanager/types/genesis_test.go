package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
)

func TestGenesisState(t *testing.T) {
	// test default genesis state
	testGenesisState := GenesisState{
		Params{},
		nil,
	}
	defaultGenesisState := DefaultGenesisState()
	require.Equal(t, *defaultGenesisState, testGenesisState)
	// test new genesis state
	newGenesisState := NewGenesisState(
		Params{},
	)
	testGenesisState = GenesisState{
		Params{},
		nil,
	}
	require.Equal(t, *newGenesisState, testGenesisState)
}

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Params Params
		Claims []*Claim
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
				Params: DefaultParams(),
				Claims: []*Claim{
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
				Params: DefaultParams(),
				Claims: []*Claim{
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
			gs := GenesisState{
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
