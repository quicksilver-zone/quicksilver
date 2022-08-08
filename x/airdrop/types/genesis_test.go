package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestValidateGenesis(t *testing.T) {
	type args struct {
		data GenesisState
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"null genesis",
			args{},
			false,
		},
		{
			"default genesis",
			args{
				*DefaultGenesisState(),
			},
			false,
		},
		{
			"duplicate zone drop",
			args{
				GenesisState{
					DefaultParams(),
					[]*ZoneDrop{
						{
							ChainId:    "test-1",
							StartTime:  time.Now().Add(1 * time.Minute),
							Duration:   time.Minute,
							Decay:      time.Minute,
							Allocation: 1000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
						{
							ChainId:    "test-1",
							StartTime:  time.Now().Add(1 * time.Hour),
							Duration:   time.Hour,
							Decay:      time.Hour,
							Allocation: 5000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
					},
					[]*ClaimRecord{},
				},
			},
			true,
		},
		{
			"invalid zone drop",
			args{
				GenesisState{
					DefaultParams(),
					[]*ZoneDrop{
						{
							ChainId:    "",
							StartTime:  time.Now().Add(1 * time.Minute),
							Duration:   -time.Minute,
							Decay:      -time.Hour,
							Allocation: 0,
							Actions:    []sdk.Dec{},
						},
					},
					[]*ClaimRecord{},
				},
			},
			true,
		},
		{
			"duplicate claim record",
			args{
				GenesisState{
					DefaultParams(),
					[]*ZoneDrop{
						{
							ChainId:    "test-1",
							StartTime:  time.Now().Add(1 * time.Minute),
							Duration:   time.Minute,
							Decay:      time.Minute,
							Allocation: 1000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
					},
					[]*ClaimRecord{
						{
							ChainId:          "test-1",
							Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
							ActionsCompleted: map[int32]*CompletedAction{},
							MaxAllocation:    500000,
						},
						{
							ChainId:          "test-1",
							Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
							ActionsCompleted: map[int32]*CompletedAction{},
							MaxAllocation:    500000,
						},
					},
				},
			},
			true,
		},
		{
			"invalid claim record",
			args{
				GenesisState{
					DefaultParams(),
					[]*ZoneDrop{
						{
							ChainId:    "test-1",
							StartTime:  time.Now().Add(1 * time.Minute),
							Duration:   time.Minute,
							Decay:      time.Minute,
							Allocation: 1000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
					},
					[]*ClaimRecord{
						{
							ChainId: "",
							Address: "",
							ActionsCompleted: map[int32]*CompletedAction{
								7: {
									CompleteTime: time.Now().Add(time.Hour),
									ClaimAmount:  1000000,
								},
							},
							MaxAllocation: 500000,
						},
					},
				},
			},
			true,
		},
		{
			"claim record no zone drop",
			args{
				GenesisState{
					DefaultParams(),
					[]*ZoneDrop{
						{
							ChainId:    "test-1",
							StartTime:  time.Now().Add(1 * time.Minute),
							Duration:   time.Minute,
							Decay:      time.Minute,
							Allocation: 1000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
					},
					[]*ClaimRecord{
						{
							ChainId: "test-2",
							Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
							ActionsCompleted: map[int32]*CompletedAction{
								0: {
									CompleteTime: time.Now().Add(-time.Hour),
									ClaimAmount:  100000,
								},
							},
							MaxAllocation: 500000,
						},
					},
				},
			},
			true,
		},
		{
			"claim record exceed zone drop",
			args{
				GenesisState{
					DefaultParams(),
					[]*ZoneDrop{
						{
							ChainId:    "test-1",
							StartTime:  time.Now().Add(1 * time.Minute),
							Duration:   time.Minute,
							Decay:      time.Minute,
							Allocation: 1000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
					},
					[]*ClaimRecord{
						{
							ChainId:          "test-1",
							Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
							ActionsCompleted: map[int32]*CompletedAction{},
							MaxAllocation:    600000,
						},
						{
							ChainId:          "test-1",
							Address:          "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
							ActionsCompleted: map[int32]*CompletedAction{},
							MaxAllocation:    600000,
						},
					},
				},
			},
			true,
		},
		{
			"no claim records",
			args{
				GenesisState{
					DefaultParams(),
					[]*ZoneDrop{
						{
							ChainId:    "test-1",
							StartTime:  time.Now().Add(1 * time.Minute),
							Duration:   time.Minute,
							Decay:      time.Minute,
							Allocation: 1000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
						{
							ChainId:    "test-2",
							StartTime:  time.Now().Add(1 * time.Hour),
							Duration:   time.Hour,
							Decay:      time.Hour,
							Allocation: 1000000,
							Actions:    []sdk.Dec{sdk.OneDec()},
						},
					},
					[]*ClaimRecord{
						{
							ChainId:          "test-1",
							Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
							ActionsCompleted: map[int32]*CompletedAction{},
							MaxAllocation:    500000,
						},
						{
							ChainId:          "test-1",
							Address:          "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
							ActionsCompleted: map[int32]*CompletedAction{},
							MaxAllocation:    500000,
						},
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGenesis(tt.args.data)
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
