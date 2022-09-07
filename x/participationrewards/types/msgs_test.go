package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgSubmitClaim_ValidateBasic(t *testing.T) {
	type fields struct {
		UserAddress string
		Zone        string
		ProofType   int64
		Proofs      []*Proof
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
			msg := MsgSubmitClaim{
				UserAddress: tt.fields.UserAddress,
				Zone:        tt.fields.Zone,
				ProofType:   tt.fields.ProofType,
				Proofs:      tt.fields.Proofs,
			}
			err := msg.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
