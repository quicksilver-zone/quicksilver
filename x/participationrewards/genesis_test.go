package participationrewards_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/participationrewards"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/keeper"
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
		t.Fatal(fmt.Errorf("unable to marshal protocol data"))
	}
	protocolData := keeper.NewProtocolData("osmosispool", "osmosis", bz)

	app.ParticipationRewardsKeeper.SetProtocolData(ctx, fmt.Sprintf("osmosis/pools/%d", pool.PoolID), protocolData)

	genesis := participationrewards.ExportGenesis(ctx, app.ParticipationRewardsKeeper)

	require.Equal(t, "osmosis/pools/1", genesis.ProtocolData[0].Key)
	require.Equal(t, "osmosis", genesis.ProtocolData[0].ProtocolData.Protocol)
	require.Equal(t, "osmosispool", genesis.ProtocolData[0].ProtocolData.Type)
}

func TestParticipationRewardsInitGenesis(t *testing.T) {
	// setup params
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	now := time.Now()
	ctx = ctx.WithBlockHeight(1)
	ctx = ctx.WithBlockTime(now)

	validOsmosisData := `{
		"poolname": "osmosis/pools/1",
		"zones": {
			"zone_id": "IBC/zone_denom"
		}
	}`

	kpd := &types.KeyedProtocolData{
		Key: "pools/6",
		ProtocolData: &types.ProtocolData{
			Protocol: "osmosis",
			Type:     "osmosispool",
			Data:     []byte(validOsmosisData),
		},
	}

	claim := &types.Claim{UserAddress: "cosmos1e9adutp4mvamq7m8eqarz57u8ymh7mhqxqfxpr", ChainId: "cosmoshub-1", Amount: 100, Module: types.ClaimTypeLiquidToken, SourceChainId: "osmosis-1"}

	// test genesisState validation
	genesisState := types.GenesisState{
		Params: types.Params{
			DistributionProportions: types.DistributionProportions{
				ValidatorSelectionAllocation: sdk.NewDecWithPrec(5, 1),
				HoldingsAllocation:           sdk.NewDecWithPrec(5, 1),
				LockupAllocation:             sdk.ZeroDec(),
			},
		},
		Claims:       []*types.Claim{claim},
		ProtocolData: []*types.KeyedProtocolData{kpd},
	}
	require.NoError(t, genesisState.Validate(), "genesis validation failed")

	participationrewards.InitGenesis(ctx, app.ParticipationRewardsKeeper, genesisState)

	require.Equal(t, app.ParticipationRewardsKeeper.GetParams(ctx).DistributionProportions.ValidatorSelectionAllocation, sdk.NewDecWithPrec(5, 1))
	require.Equal(t, app.ParticipationRewardsKeeper.GetParams(ctx).DistributionProportions.HoldingsAllocation, sdk.NewDecWithPrec(5, 1))
	require.Equal(t, app.ParticipationRewardsKeeper.GetParams(ctx).DistributionProportions.LockupAllocation, sdk.ZeroDec())

	pd, found := app.ParticipationRewardsKeeper.GetProtocolData(ctx, "pools/6")
	require.True(t, found)
	require.Equal(t, "osmosis", pd.Protocol)
	require.Equal(t, "osmosispool", pd.Type)

	clm, found := app.ParticipationRewardsKeeper.GetClaim(ctx, "cosmoshub-1", "cosmos1e9adutp4mvamq7m8eqarz57u8ymh7mhqxqfxpr", types.ClaimTypeLiquidToken, "osmosis-1")
	require.True(t, found)
	require.Equal(t, uint64(100), clm.Amount)
}
