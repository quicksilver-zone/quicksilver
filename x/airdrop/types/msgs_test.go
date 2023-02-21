package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

func TestMsgClaim_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainId string
		Action  int64
		Address string
		Proofs  []*cmtypes.Proof
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
			"invalid_no_zone",
			fields{
				ChainId: "",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_action_out_of_bounds_low",
			fields{
				ChainId: "cosmoshub-4",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_action_out_of_bounds",
			fields{
				ChainId: "cosmoshub-4",
				Action:  999,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_address_empty",
			fields{
				ChainId: "cosmoshub-4",
				Action:  0,
				Address: "",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_address",
			fields{
				ChainId: "cosmoshub-4",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lkq437x2w",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		// TODO: add more address checks
		//   - currently it fails using quick address (no sdk setup done)
		{
			"invalid_ActionUndefined",
			fields{
				ChainId: "cosmoshub-4",
				Action:  int64(ActionUndefined),
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs: []*cmtypes.Proof{
					{
						Key:       []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						ProofOps:  &crypto.ProofOps{},
						Height:    10,
						ProofType: "lockup",
					},
				},
			},
			true,
		},
		{
			"valid",
			fields{
				ChainId: "cosmoshub-4",
				Action:  int64(ActionInitialClaim),
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs: []*cmtypes.Proof{
					{
						Key:       []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						ProofOps:  &crypto.ProofOps{},
						Height:    10,
						ProofType: "lockup",
					},
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
