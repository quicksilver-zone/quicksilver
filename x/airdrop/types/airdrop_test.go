package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestZoneDrop_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainID     string
		StartTime   time.Time
		Duration    time.Duration
		Decay       time.Duration
		Allocation  uint64
		Actions     []sdk.Dec
		IsConcluded bool
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
			"invalid-weights-exceeded",
			fields{
				ChainID:    "",
				StartTime:  time.Now().Add(time.Hour),
				Duration:   -time.Minute,
				Decay:      -time.Second,
				Allocation: 0,
				Actions: []sdk.Dec{
					sdk.MustNewDecFromStr("0.3"),
					sdk.MustNewDecFromStr("0.4"),
					sdk.MustNewDecFromStr("0.5"),
				},
				IsConcluded: false,
			},
			true,
		},
		{
			"invalid-weights-insufficient",
			fields{
				ChainID:    "",
				StartTime:  time.Now().Add(time.Hour),
				Duration:   0,
				Decay:      0,
				Allocation: 0,
				Actions: []sdk.Dec{
					sdk.MustNewDecFromStr("0.3"),
					sdk.MustNewDecFromStr("0.3"),
					sdk.MustNewDecFromStr("0.3"),
				},
				IsConcluded: false,
			},
			true,
		},
		{
			"invalid-actions-exceeded",
			fields{
				ChainID:    "test-1",
				StartTime:  time.Now().Add(-time.Hour),
				Duration:   time.Hour,
				Decay:      30 * time.Minute,
				Allocation: 16400,
				Actions: []sdk.Dec{
					sdk.MustNewDecFromStr("0.01"),
					sdk.MustNewDecFromStr("0.02"),
					sdk.MustNewDecFromStr("0.03"),
					sdk.MustNewDecFromStr("0.04"),
					sdk.MustNewDecFromStr("0.06"),
					sdk.MustNewDecFromStr("0.07"),
					sdk.MustNewDecFromStr("0.08"),
					sdk.MustNewDecFromStr("0.09"),
					sdk.MustNewDecFromStr("0.1"),
					sdk.MustNewDecFromStr("0.2"),
					sdk.MustNewDecFromStr("0.3"),
					sdk.MustNewDecFromStr("0.1"),
				},
				IsConcluded: false,
			},
			true,
		},
		{
			"valid",
			fields{
				ChainID:    "test-1",
				StartTime:  time.Now().Add(-time.Hour),
				Duration:   time.Hour,
				Decay:      30 * time.Minute,
				Allocation: 16400,
				Actions: []sdk.Dec{
					sdk.MustNewDecFromStr("0.1"),
					sdk.MustNewDecFromStr("0.2"),
					sdk.MustNewDecFromStr("0.3"),
					sdk.MustNewDecFromStr("0.4"),
				},
				IsConcluded: false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zd := ZoneDrop{
				ChainId:     tt.fields.ChainID,
				StartTime:   tt.fields.StartTime,
				Duration:    tt.fields.Duration,
				Decay:       tt.fields.Decay,
				Allocation:  tt.fields.Allocation,
				Actions:     tt.fields.Actions,
				IsConcluded: tt.fields.IsConcluded,
			}

			err := zd.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestClaimRecord_ValidateBasic(t *testing.T) {
	type fields struct {
		ChainID          string
		Address          string
		ActionsCompleted map[int32]*CompletedAction
		MaxAllocation    uint64
		BaseValue        uint64
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
			"invalid-00",
			fields{
				ChainID:       "",
				Address:       "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9lj",
				MaxAllocation: 0,
				ActionsCompleted: map[int32]*CompletedAction{
					0: {
						CompleteTime: time.Now().Add(-time.Hour),
						ClaimAmount:  0,
					},
				},
				BaseValue: 0,
			},
			true,
		},
		{
			"invalid-01",
			fields{
				ChainID:       "test-01",
				Address:       "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
				MaxAllocation: 144000,
				ActionsCompleted: map[int32]*CompletedAction{
					12: {
						CompleteTime: time.Now().Add(-time.Minute),
						ClaimAmount:  150000,
					},
				},
				BaseValue: 25000,
			},
			true,
		},
		{
			"invalid-02",
			fields{
				ChainID:       "test-01",
				Address:       "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
				MaxAllocation: 144000,
				ActionsCompleted: map[int32]*CompletedAction{
					1: {
						CompleteTime: time.Now().Add(-time.Hour),
						ClaimAmount:  50000,
					},
					2: {
						CompleteTime: time.Now().Add(-time.Minute),
						ClaimAmount:  100000,
					},
					3: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  50000,
					},
				},
				BaseValue: 25000,
			},
			true,
		},
		{
			"valid",
			fields{
				ChainID:       "test-01",
				Address:       "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
				MaxAllocation: 144000,
				ActionsCompleted: map[int32]*CompletedAction{
					1: {
						CompleteTime: time.Now().Add(-time.Hour),
						ClaimAmount:  12000,
					},
					2: {
						CompleteTime: time.Now().Add(-time.Minute),
						ClaimAmount:  12000,
					},
					3: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					4: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					5: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					6: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					7: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					8: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					9: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					10: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
					11: {
						CompleteTime: time.Now().Add(-time.Second),
						ClaimAmount:  12000,
					},
				},
				BaseValue: 25000,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := ClaimRecord{
				ChainId:          tt.fields.ChainID,
				Address:          tt.fields.Address,
				ActionsCompleted: tt.fields.ActionsCompleted,
				MaxAllocation:    tt.fields.MaxAllocation,
				BaseValue:        tt.fields.BaseValue,
			}

			err := cr.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAction_InBounds(t *testing.T) {
	tests := []struct {
		name string
		a    Action
		want bool
	}{
		{
			"exceed lower",
			ActionUndefined,
			false,
		},
		{
			"exceed upper",
			Action(len(Action_name)),
			false,
		},
		{
			"in bounds lower",
			1,
			true,
		},
		{
			"in bounds upper",
			Action(len(Action_name) - 1),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.InBounds()
			if got != tt.want {
				err := fmt.Errorf("Action.InBounds() = %v, want %v", got, tt.want)
				require.NoError(t, err)
			}
		})
	}
}
