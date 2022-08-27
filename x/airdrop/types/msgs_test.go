package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgClaim_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainId string
		Action  int32
		Address string
		Proofs  []*Proof
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"empty MsgClaim",
			fields{},
			true,
		},
		{
			"no zone",
			fields{
				ChainId: "",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*Proof{},
			},
			true,
		},
		{
			"action out of bounds",
			fields{
				ChainId: "cosmoshub-4",
				Action:  -1,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*Proof{},
			},
			true,
		},
		{
			"action out of bounds",
			fields{
				ChainId: "cosmoshub-4",
				Action:  999,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*Proof{},
			},
			true,
		},
		{
			"invalid address",
			fields{
				ChainId: "cosmoshub-4",
				Action:  0,
				Address: "",
				Proofs:  []*Proof{},
			},
			true,
		},
		{
			"invalid address",
			fields{
				ChainId: "cosmoshub-4",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lkq437x2w",
				Proofs:  []*Proof{},
			},
			true,
		},
		// TODO: add more address checks
		//   - currently it fails using quick address (no sdk setup done)
		{
			"valid",
			fields{
				ChainId: "cosmoshub-4",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs: []*Proof{
					{},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgClaim{
				ChainId: tt.fields.ChainId,
				Action:  tt.fields.Action,
				Address: tt.fields.Address,
				Proofs:  tt.fields.Proofs,
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
