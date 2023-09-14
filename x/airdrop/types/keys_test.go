package types_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetKeyZoneDrop(t *testing.T) {
	testID := "test-01"
	type args struct {
		chainID string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"valid",
			args{
				chainID: testID,
			},
			append([]byte{0x1}, []byte(testID)...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.GetKeyZoneDrop(tt.args.chainID)
			if !reflect.DeepEqual(got, tt.want) {
				err := fmt.Errorf("GetKeyZoneDrop() = %v, want %v", got, tt.want)
				require.NoError(t, err)
				return
			}
		})
	}
}

func TestGetKeyClaimRecord(t *testing.T) {
	testID := "test-01"
	testAddress := "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"
	testAcc, _ := sdk.AccAddressFromBech32(testAddress)
	type args struct {
		chainID string
		addr    sdk.AccAddress
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"valid",
			args{
				chainID: testID,
				addr:    testAcc,
			},
			append(append([]byte{0x2}, []byte(testID)...), testAcc...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.GetKeyClaimRecord(tt.args.chainID, tt.args.addr)
			if !reflect.DeepEqual(got, tt.want) {
				err := fmt.Errorf("GetKeyClaimRecord() = %v, want %v", got, tt.want)
				// t.Logf("Error:\n%v\n", err)
				require.NoError(t, err)
				return
			}
		})
	}
}

func TestGetPrefixClaimRecord(t *testing.T) {
	testID := "tester-01"
	type args struct {
		chainID string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"valid",
			args{
				chainID: testID,
			},
			append([]byte{0x2}, []byte(testID)...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.GetPrefixClaimRecord(tt.args.chainID)
			if !reflect.DeepEqual(got, tt.want) {
				err := fmt.Errorf("GetPrefixClaimRecord() = %v, want %v", got, tt.want)
				// t.Logf("Error:\n%v\n", err)
				require.NoError(t, err)

				return
			}
		})
	}
}
