package epochs_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	simapp "github.com/quicksilver-zone/quicksilver/v7/app"
	"github.com/quicksilver-zone/quicksilver/v7/x/epochs"
	"github.com/quicksilver-zone/quicksilver/v7/x/epochs/types"
)

func TestEpochsExportGenesis(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	chainStartTime := ctx.BlockTime()
	chainStartHeight := ctx.BlockHeight()

	genesis := epochs.ExportGenesis(ctx, app.EpochsKeeper)
	require.Len(t, genesis.Epochs, 3)

	require.Equal(t, genesis.Epochs[0].Identifier, "day")
	require.Equal(t, genesis.Epochs[0].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[0].Duration, time.Hour*24)
	require.Equal(t, genesis.Epochs[0].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[0].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[0].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[0].EpochCountingStarted, false)
	require.Equal(t, genesis.Epochs[1].Identifier, "epoch")
	require.Equal(t, genesis.Epochs[1].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[1].Duration, time.Second*240)
	require.Equal(t, genesis.Epochs[1].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[1].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[1].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[1].EpochCountingStarted, false)
	require.Equal(t, genesis.Epochs[2].Identifier, "week")
	require.Equal(t, genesis.Epochs[2].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[2].Duration, time.Hour*24*7)
	require.Equal(t, genesis.Epochs[2].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[2].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[2].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[2].EpochCountingStarted, false)
}

func TestEpochsInitGenesis(t *testing.T) {
	// setup feemarketGenesis params
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// On init genesis, default epochs information is set
	// To check init genesis again, should make it fresh status
	epochInfos := app.EpochsKeeper.AllEpochInfos(ctx)
	for _, epochInfo := range epochInfos {
		app.EpochsKeeper.DeleteEpochInfo(ctx, epochInfo.Identifier)
	}

	now := time.Now()
	ctx = ctx.WithBlockHeight(1)
	ctx = ctx.WithBlockTime(now)

	// test genesisState validation
	genesisState := types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}
	require.EqualError(t, genesisState.Validate(), "value #2: epoch identifier should be unique, got duplicate \"monthly\"")

	genesisState = types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}

	epochs.InitGenesis(ctx, app.EpochsKeeper, genesisState)

	epochInfo := app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
	require.Equal(t, epochInfo.Identifier, "monthly")
	require.Equal(t, epochInfo.StartTime.UTC().String(), now.UTC().String())
	require.Equal(t, epochInfo.Duration, time.Hour*24)
	require.Equal(t, epochInfo.CurrentEpoch, int64(0))
	require.Equal(t, epochInfo.CurrentEpochStartHeight, ctx.BlockHeight())
	require.Equal(t, epochInfo.CurrentEpochStartTime.UTC().String(), time.Time{}.String())
	require.Equal(t, epochInfo.EpochCountingStarted, true)
}
