package types

import (
	fmt "fmt"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
)

func TestMsgSubmitClaim_ValidateBasic(t *testing.T) {
	userAddress := utils.GenerateAccAddressForTest().String()

	type fields struct {
		UserAddress string
		Zone        string
		ClaimType   ClaimType
		Proofs      []*Proof
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
			"invalid_empty_proof",
			fields{
				UserAddress: "cosmos1234567890abcde",
				Zone:        "",
				ClaimType:   -1,
				Proofs:      []*Proof{},
			},
			true,
		},
		{
			"invalid_with_proof",
			fields{
				UserAddress: "cosmos1234567890abcde",
				Zone:        "",
				ClaimType:   -1,
				Proofs: []*Proof{
					{}, // blank
					{
						Key:       []byte{1, 2, 3, 4, 5},
						Data:      []byte{0, 0, 1, 1, 2, 3, 4, 5},
						ProofOps:  nil,
						Height:    -1,
						ProofType: "bank",
					},
				},
			},
			true,
		},
		{
			"valid",
			fields{
				UserAddress: userAddress,
				Zone:        "test-01",
				ClaimType:   ClaimTypeOsmosisPool,
				Proofs: []*Proof{
					{
						Key:       []byte{1, 2, 3, 4, 5},
						Data:      []byte{0, 0, 1, 1, 2, 3, 4, 5},
						ProofOps:  &crypto.ProofOps{},
						Height:    123,
						ProofType: "lockup",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSubmitClaim{
				UserAddress: tt.fields.UserAddress,
				Zone:        tt.fields.Zone,
				ClaimType:   tt.fields.ClaimType,
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

func TestMsgSubmitClaim_GetSigners(t *testing.T) {
	validAddress := utils.GenerateAccAddressForTest().String()
	validAcc, _ := sdk.AccAddressFromBech32(validAddress)

	type fields struct {
		UserAddress string
		Zone        string
		ClaimType   ClaimType
		Proofs      []*Proof
	}
	tests := []struct {
		name   string
		fields fields
		want   []sdk.AccAddress
	}{
		{
			"blank",
			fields{},
			[]sdk.AccAddress{{}},
		},
		{
			"valid",
			fields{
				UserAddress: validAddress,
			},
			[]sdk.AccAddress{validAcc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSubmitClaim{
				UserAddress: tt.fields.UserAddress,
				Zone:        tt.fields.Zone,
				ClaimType:   tt.fields.ClaimType,
				Proofs:      tt.fields.Proofs,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				err := fmt.Errorf("MsgSubmitClaim.GetSigners() = %v, want %v", got, tt.want)
				require.NoError(t, err)
			}
		})
	}
}

// func TestNewMsgSubmitClaim(t *testing.T) {
// 	validAddress := utils.GenerateAccAddressForTest().String()
// 	testZone := "test-01"

// 	type args struct {
// 		userAddress sdk.Address
// 		zone        string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *MsgSubmitClaim
// 	}{
// 		{
// 			"valid",
// 			args{
// 				userAddress: sdk.MustAccAddressFromBech32(validAddress),
// 				zone:        testZone,
// 			},
// 			&MsgSubmitClaim{
// 				UserAddress: validAddress,
// 				Zone:        testZone,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewMsgSubmitClaim(tt.args.userAddress, tt.args.zone); !reflect.DeepEqual(got, tt.want) {
// 				err := fmt.Errorf("NewMsgSubmitClaim() = %v, want %v", got, tt.want)
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }
