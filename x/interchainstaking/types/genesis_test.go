package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Params           Params
		Zones            []Zone
		Receipts         []Receipt
		Delegations      []DelegationsForZone
		DelegatorIntents []DelegatorIntentsForZone
		PortConnections  []PortConnectionTuple
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
			gs := GenesisState{
				Params:           tt.fields.Params,
				Zones:            tt.fields.Zones,
				Receipts:         tt.fields.Receipts,
				Delegations:      tt.fields.Delegations,
				DelegatorIntents: tt.fields.DelegatorIntents,
				PortConnections:  tt.fields.PortConnections,
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
