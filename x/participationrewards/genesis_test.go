package participationrewards_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	simapp "github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/participationrewards"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func TestParticipationRewardsExportGenesis(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	chainStartTime := ctx.BlockTime()

	pool := types.OsmosisPoolProtocolData{
		PoolID:      1,
		PoolName:    "atom/osmo",
		LastUpdated: chainStartTime,
	}

	bz, err := json.Marshal(pool)
	if err != nil {
		t.Fatalf("unable to marshal protocol data: %v", err)
	}
	protocolData := types.NewProtocolData(types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)], bz)

	app.ParticipationRewardsKeeper.SetProtocolData(ctx, fmt.Sprintf("%d", pool.PoolID), protocolData)

	genesis := participationrewards.ExportGenesis(ctx, app.ParticipationRewardsKeeper)

	// 0,0,0,4 (binary encoded types.ProtocolDataTypeOsmosisPool)
	// 49 (ASCII value of '1')
	require.Equal(t, string([]byte{0, 0, 0, 0, 0, 0, 0, 4, 49}), genesis.ProtocolData[0].Key)
	require.Equal(t, types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)], genesis.ProtocolData[0].ProtocolData.Type)
}

func TestParticipationRewardsInitGenesis(t *testing.T) {
	// setup params
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	now := time.Now()
	ctx = ctx.WithBlockHeight(1)
	ctx = ctx.WithBlockTime(now)

	validOsmosisData := `{
	"poolid": 1,
	"poolname": "atom/osmo",
	"pooltype": "balancer",
	"zones": {
		"zone_id": "IBC/zone_denom"
	}
}`

	kpd := &types.KeyedProtocolData{
		Key: "6",
		ProtocolData: &types.ProtocolData{
			Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)],
			Data: []byte(validOsmosisData),
		},
	}

	// test genesisState validation
	genesisState := types.GenesisState{
		Params: types.Params{
			DistributionProportions: types.DistributionProportions{
				ValidatorSelectionAllocation: sdk.NewDecWithPrec(5, 1),
				HoldingsAllocation:           sdk.NewDecWithPrec(5, 1),
				LockupAllocation:             sdk.ZeroDec(),
			},
		},
		ProtocolData: []*types.KeyedProtocolData{kpd},
	}
	require.NoError(t, genesisState.Validate(), "genesis validation failed")

	participationrewards.InitGenesis(ctx, app.ParticipationRewardsKeeper, genesisState)

	require.Equal(t, app.ParticipationRewardsKeeper.GetParams(ctx).DistributionProportions.ValidatorSelectionAllocation, sdk.NewDecWithPrec(5, 1))
	require.Equal(t, app.ParticipationRewardsKeeper.GetParams(ctx).DistributionProportions.HoldingsAllocation, sdk.NewDecWithPrec(5, 1))
	require.Equal(t, app.ParticipationRewardsKeeper.GetParams(ctx).DistributionProportions.LockupAllocation, sdk.ZeroDec())

	pd, found := app.ParticipationRewardsKeeper.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisPool, "6")
	require.True(t, found)
	require.Equal(t, types.ProtocolDataType_name[int32(types.ProtocolDataTypeOsmosisPool)], pd.Type)
}
