package types_test

import (
	"fmt"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func TestMsgSubmitClaim_ValidateBasic(t *testing.T) {
	userAddress := addressutils.GenerateAccAddressForTest().String()

	type fields struct {
		UserAddress string
		Zone        string
		SrcZone     string
		ClaimType   cmtypes.ClaimType
		Proofs      []*cmtypes.Proof
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
				SrcZone:     "",
				ClaimType:   -1,
				Proofs:      []*cmtypes.Proof{},
			},
			true,
		},
		{
			"invalid_with_proof",
			fields{
				UserAddress: "cosmos1234567890abcde",
				Zone:        "",
				SrcZone:     "",
				ClaimType:   -1,
				Proofs: []*cmtypes.Proof{
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
				SrcZone:     "test-02",
				ClaimType:   cmtypes.ClaimTypeOsmosisPool,
				Proofs: []*cmtypes.Proof{
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
			msg := types.MsgSubmitClaim{
				UserAddress: tt.fields.UserAddress,
				Zone:        tt.fields.Zone,
				SrcZone:     tt.fields.SrcZone,
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
	validAddress := addressutils.GenerateAccAddressForTest().String()
	validAcc, _ := sdk.AccAddressFromBech32(validAddress)

	type fields struct {
		UserAddress string
		Zone        string
		ClaimType   cmtypes.ClaimType
		Proofs      []*cmtypes.Proof
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
			msg := types.MsgSubmitClaim{
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

func TestNewMsgSubmitClaim(t *testing.T) {
	userAddress := addressutils.GenerateAccAddressForTest()
	type args struct {
		userAddress sdk.Address
		srcZone     string
		zone        string
		claimType   cmtypes.ClaimType
		proofs      []*cmtypes.Proof
	}
	tests := []struct {
		name string
		args args
		want *types.MsgSubmitClaim
	}{
		{
			name: "test",
			args: args{
				userAddress,
				"osmosis-1",
				"juno",
				cmtypes.ClaimTypeOsmosisPool,
				[]*cmtypes.Proof{
					{
						Key:       []byte{1, 2, 3, 4, 5},
						Data:      []byte{0, 0, 1, 1, 2, 3, 4, 5},
						ProofOps:  &crypto.ProofOps{},
						Height:    123,
						ProofType: "lockup",
					},
				},
			},
			want: &types.MsgSubmitClaim{
				UserAddress: userAddress.String(),
				Zone:        "juno",
				SrcZone:     "osmosis-1",
				ClaimType:   cmtypes.ClaimTypeOsmosisPool,
				Proofs: []*cmtypes.Proof{
					{
						Key:       []byte{1, 2, 3, 4, 5},
						Data:      []byte{0, 0, 1, 1, 2, 3, 4, 5},
						ProofOps:  &crypto.ProofOps{},
						Height:    123,
						ProofType: "lockup",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.NewMsgSubmitClaim(tt.args.userAddress, tt.args.srcZone, tt.args.zone, tt.args.claimType, tt.args.proofs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMsgSubmitClaim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgGovRemoveProtocolData_ValidateBasic(t *testing.T) {
	userAddress := addressutils.GenerateAccAddressForTest().String()

	type fields struct {
		Title       string
		Description string
		Key         string
		Authority   string
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
			"invalid_empty_title",
			fields{
				Title:       "",
				Description: "Non-nil description",
				Key:         "123",
				Authority:   "quick123",
			},
			true,
		},
		{
			"invalid_empty_desc",
			fields{
				Title:       "Non-nil title",
				Description: "",
				Key:         "123",
				Authority:   "quick123",
			},
			true,
		},
		{
			"invalid_empty_key",
			fields{
				Title:       "Non-nil title",
				Description: "Non-nil description",
				Key:         "",
				Authority:   "quick123",
			},
			true,
		},
		{
			"invalid_empty_authority",
			fields{
				Title:       "Non-nil title",
				Description: "Non-nil description",
				Key:         "123",
				Authority:   "",
			},
			true,
		},
		{
			"invalid_invalid_authority",
			fields{
				Title:       "Non-nil title",
				Description: "Non-nil description",
				Key:         "123",
				Authority:   "quick123",
			},
			true,
		},
		{
			"valid",
			fields{
				Title:       "Non-nil title",
				Description: "Non-nil description",
				Key:         "123",
				Authority:   userAddress,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.MsgGovRemoveProtocolData{
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Key:         tt.fields.Key,
				Authority:   tt.fields.Authority,
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

func TestMsgGovRemoveProtocolData_GetSigners(t *testing.T) {
	validAddress := addressutils.GenerateAccAddressForTest().String()
	validAcc, _ := sdk.AccAddressFromBech32(validAddress)

	type fields struct {
		Title       string
		Description string
		Key         string
		Authority   string
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
				Authority: validAddress,
			},
			[]sdk.AccAddress{validAcc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.MsgGovRemoveProtocolData{
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Key:         tt.fields.Key,
				Authority:   tt.fields.Authority,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				err := fmt.Errorf("MsgGovRemoveProtocolData.GetSigners() = %v, want %v", got, tt.want)
				require.NoError(t, err)
			}
		})
	}
}
