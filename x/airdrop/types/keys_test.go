package types

import (
	"fmt"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetKeyZoneDrop(t *testing.T) {
	testId := "test-01"
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
				chainID: testId,
			},
			append([]byte{0x1}, []byte(testId)...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetKeyZoneDrop(tt.args.chainID)
			if !reflect.DeepEqual(got, tt.want) {
				err := fmt.Errorf("GetKeyZoneDrop() = %v, want %v", got, tt.want)
				// t.Logf("Error:\n%v\n", err)
				require.NoError(t, err)
				return
			}
		})
	}
}

func TestGetKeyClaimRecord(t *testing.T) {
	testId := "test-01"
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
				chainID: testId,
				addr:    testAcc,
			},
			append(append([]byte{0x2}, []byte(testId)...), testAcc...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetKeyClaimRecord(tt.args.chainID, tt.args.addr)
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
	testId := "tester-01"
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
				chainID: testId,
			},
			append([]byte{0x2}, []byte(testId)...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPrefixClaimRecord(tt.args.chainID)
			if !reflect.DeepEqual(got, tt.want) {
				err := fmt.Errorf("GetPrefixClaimRecord() = %v, want %v", got, tt.want)
				// t.Logf("Error:\n%v\n", err)
				require.NoError(t, err)

				return
			}
		})
	}
}
