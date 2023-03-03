package osmosistypes

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	osmosislockuptypes "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type ParticipationRewardsKeeper interface {
	GetProtocolData(ctx sdk.Context, pdType participationrewardstypes.ProtocolDataType, key string) (participationrewardstypes.ProtocolData, bool)
}

func DetermineApplicableTokensInPool(ctx sdk.Context, prKeeper ParticipationRewardsKeeper, lock osmosislockuptypes.PeriodLock, chainID string) (math.Int, error) {
	gammtoken, err := lock.SingleCoin()
	if err != nil {
		return sdk.ZeroInt(), err
	}

	poolID := gammtoken.Denom[strings.LastIndex(gammtoken.Denom, "/")+1:]
	pd, ok := prKeeper.GetProtocolData(ctx, participationrewardstypes.ProtocolDataTypeOsmosisPool, poolID)
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to obtain protocol data for poolID=%s", poolID)
	}

	ipool, err := participationrewardstypes.UnmarshalProtocolData(participationrewardstypes.ProtocolDataTypeOsmosisPool, pd.Data)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	pool, _ := ipool.(participationrewardstypes.OsmosisPoolProtocolData)

	poolDenom := ""
	for zk, zd := range pool.Zones {
		if zk == chainID {
			poolDenom = zd
			break
		}
	}

	if poolDenom == "" {
		return sdk.ZeroInt(), fmt.Errorf("invalid zone, pool zone must match %s", chainID)
	}

	poolData, err := pool.GetPool()
	if err != nil {
		return sdk.ZeroInt(), err
	}
	// calculate user gamm ratio and LP asset amount
	ugamm := gammtoken.Amount          // user's gamm amount
	pgamm := poolData.GetTotalShares() // total pool gamm amount
	if pgamm.IsZero() {
		return sdk.ZeroInt(), fmt.Errorf("empty pool, %s", poolID)
	}
	uratio := sdk.NewDecFromInt(ugamm).QuoInt(pgamm)

	zasset := poolData.GetTotalPoolLiquidity(ctx).AmountOf(poolDenom) // pool zone asset amount
	uAmount := uratio.MulInt(zasset).TruncateInt()

	return uAmount, nil
}
