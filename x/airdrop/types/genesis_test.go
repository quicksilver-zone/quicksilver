package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
)

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Params       types.Params
		ZoneDrops    []*types.ZoneDrop
		ClaimRecords []*types.ClaimRecord
	}

	defaultGenesis := types.DefaultGenesisState()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"null genesis",
			fields{},
			false,
		},
		{
			"default genesis",
			fields{
				defaultGenesis.Params,
				defaultGenesis.ZoneDrops,
				defaultGenesis.ClaimRecords,
			},
			false,
		},
		{
			"duplicate zone drop",
			fields{
				types.DefaultParams(),
				[]*types.ZoneDrop{
					{
						ChainId:    "test-1",
						StartTime:  time.Now().Add(1 * time.Minute),
						Duration:   time.Minute,
						Decay:      time.Minute,
						Allocation: 1000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
					{
						ChainId:    "test-1",
						StartTime:  time.Now().Add(1 * time.Hour),
						Duration:   time.Hour,
						Decay:      time.Hour,
						Allocation: 5000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
				},
				[]*types.ClaimRecord{},
			},
			true,
		},
		{
			"invalid zone drop",
			fields{
				types.DefaultParams(),
				[]*types.ZoneDrop{
					{
						ChainId:    "",
						StartTime:  time.Now().Add(1 * time.Minute),
						Duration:   -time.Minute,
						Decay:      -time.Hour,
						Allocation: 0,
						Actions:    []sdkmath.LegacyDec{},
					},
				},
				[]*types.ClaimRecord{},
			},
			true,
		},
		{
			"duplicate claim record",
			fields{
				types.DefaultParams(),
				[]*types.ZoneDrop{
					{
						ChainId:    "test-1",
						StartTime:  time.Now().Add(1 * time.Minute),
						Duration:   time.Minute,
						Decay:      time.Minute,
						Allocation: 1000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
				},
				[]*types.ClaimRecord{
					{
						ChainId:          "test-1",
						Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
						ActionsCompleted: map[int32]*types.CompletedAction{},
						MaxAllocation:    500000,
					},
					{
						ChainId:          "test-1",
						Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
						ActionsCompleted: map[int32]*types.CompletedAction{},
						MaxAllocation:    500000,
					},
				},
			},
			true,
		},
		{
			"invalid claim record",
			fields{
				types.DefaultParams(),
				[]*types.ZoneDrop{
					{
						ChainId:    "test-1",
						StartTime:  time.Now().Add(1 * time.Minute),
						Duration:   time.Minute,
						Decay:      time.Minute,
						Allocation: 1000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
				},
				[]*types.ClaimRecord{
					{
						ChainId: "",
						Address: "",
						ActionsCompleted: map[int32]*types.CompletedAction{
							999: {
								CompleteTime: time.Now().Add(time.Hour),
								ClaimAmount:  1000000,
							},
						},
						MaxAllocation: 500000,
					},
				},
			},
			true,
		},
		{
			"claim record no zone drop",
			fields{
				types.DefaultParams(),
				[]*types.ZoneDrop{
					{
						ChainId:    "test-1",
						StartTime:  time.Now().Add(1 * time.Minute),
						Duration:   time.Minute,
						Decay:      time.Minute,
						Allocation: 1000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
				},
				[]*types.ClaimRecord{
					{
						ChainId: "test-2",
						Address: "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
						ActionsCompleted: map[int32]*types.CompletedAction{
							0: {
								CompleteTime: time.Now().Add(-time.Hour),
								ClaimAmount:  100000,
							},
						},
						MaxAllocation: 500000,
					},
				},
			},
			true,
		},
		{
			"claim record exceed zone drop",
			fields{
				types.DefaultParams(),
				[]*types.ZoneDrop{
					{
						ChainId:    "test-1",
						StartTime:  time.Now().Add(1 * time.Minute),
						Duration:   time.Minute,
						Decay:      time.Minute,
						Allocation: 1000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
				},
				[]*types.ClaimRecord{
					{
						ChainId:          "test-1",
						Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
						ActionsCompleted: map[int32]*types.CompletedAction{},
						MaxAllocation:    600000,
					},
					{
						ChainId:          "test-1",
						Address:          "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
						ActionsCompleted: map[int32]*types.CompletedAction{},
						MaxAllocation:    600000,
					},
				},
			},
			true,
		},
		{
			"no claim records",
			fields{
				types.DefaultParams(),
				[]*types.ZoneDrop{
					{
						ChainId:    "test-1",
						StartTime:  time.Now().Add(1 * time.Minute),
						Duration:   time.Minute,
						Decay:      time.Minute,
						Allocation: 1000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
					{
						ChainId:    "test-2",
						StartTime:  time.Now().Add(1 * time.Hour),
						Duration:   time.Hour,
						Decay:      time.Hour,
						Allocation: 1000000,
						Actions:    []sdkmath.LegacyDec{sdkmath.LegacyOneDec()},
					},
				},
				[]*types.ClaimRecord{
					{
						ChainId:          "test-1",
						Address:          "cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w",
						ActionsCompleted: map[int32]*types.CompletedAction{},
						MaxAllocation:    500000,
					},
					{
						ChainId:          "test-1",
						Address:          "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
						ActionsCompleted: map[int32]*types.CompletedAction{},
						MaxAllocation:    500000,
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := types.GenesisState{
				Params:       tt.fields.Params,
				ZoneDrops:    tt.fields.ZoneDrops,
				ClaimRecords: tt.fields.ClaimRecords,
			}

			err := gs.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
