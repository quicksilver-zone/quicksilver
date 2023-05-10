package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

func TestMsgClaim_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainID string
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
				ChainID: "",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},

		{
			"invalid_action_out_of_bounds_low",
			fields{
				ChainID: "cosmoshub-4",
				Action:  0,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_action_out_of_bounds",
			fields{
				ChainID: "cosmoshub-4",
				Action:  999,
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_address_empty",
			fields{
				ChainID: "cosmoshub-4",
				Action:  0,
				Address: "",
				Proofs:  []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_address",
			fields{
				ChainID: "cosmoshub-4",
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
				ChainID: "cosmoshub-4",
				Action:  int64(types.ActionUndefined),
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
			"invalid proof no key",
			fields{
				ChainID: "cosmoshub-4",
				Action:  int64(types.ActionInitialClaim),
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs: []*cmtypes.Proof{
					{
						Key:       nil,
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
			"invalid proof no data",
			fields{
				ChainID: "cosmoshub-4",
				Action:  int64(types.ActionInitialClaim),
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs: []*cmtypes.Proof{
					{
						Key:       []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						Data:      nil,
						ProofOps:  &crypto.ProofOps{},
						Height:    10,
						ProofType: "lockup",
					},
				},
			},
			true,
		},
		{
			"invalid proof no proof ops",
			fields{
				ChainID: "cosmoshub-4",
				Action:  int64(types.ActionInitialClaim),
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs: []*cmtypes.Proof{
					{
						Key:       []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						ProofOps:  nil,
						Height:    10,
						ProofType: "lockup",
					},
				},
			},
			true,
		},
		{
			"invalid proof no height",
			fields{
				ChainID: "cosmoshub-4",
				Action:  int64(types.ActionInitialClaim),
				Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
				Proofs: []*cmtypes.Proof{
					{
						Key:       []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
						ProofOps:  &crypto.ProofOps{},
						Height:    -1,
						ProofType: "lockup",
					},
				},
			},
			true,
		},
		{
			"valid",
			fields{
				ChainID: "cosmoshub-4",
				Action:  int64(types.ActionInitialClaim),
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
			msg := types.MsgClaim{
				ChainId: tt.fields.ChainID,
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

func TestMsgIncentivePoolSpendValidateBasic(t *testing.T) {
	type fields struct {
		Authority string
		ToAddress string
		Amount    sdk.Coins
	}

	validTestCoins := sdk.NewCoins(sdk.NewCoin("test", sdk.NewIntFromUint64(10000)))
	addr1 := utils.GenerateAccAddressForTest().String()
	addr2 := utils.GenerateAccAddressForTest().String()

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
			"invalid authority",
			fields{
				Authority: "invalid",
				ToAddress: addr2,
				Amount:    validTestCoins,
			},
			true,
		},
		{
			"invalid to address",
			fields{
				Authority: addr1,
				ToAddress: "invalid",
				Amount:    validTestCoins,
			},
			true,
		},
		{
			"invalid equal addresses",
			fields{
				Authority: addr1,
				ToAddress: addr1,
				Amount:    validTestCoins,
			},
			true,
		},
		{
			"non positive amount",
			fields{
				Authority: addr1,
				ToAddress: addr2,
				Amount:    sdk.Coins{},
			},
			true,
		}, {
			"invalid amount",
			fields{
				Authority: addr1,
				ToAddress: addr2,
				Amount: sdk.Coins{sdk.Coin{
					Denom:  "",
					Amount: sdkmath.NewInt(1000),
				}},
			},
			true,
		},
		{
			"valid",
			fields{
				Authority: addr1,
				ToAddress: addr2,
				Amount:    validTestCoins,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.MsgIncentivePoolSpend{
				Authority: tt.fields.Authority,
				ToAddress: tt.fields.ToAddress,
				Amount:    tt.fields.Amount,
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
