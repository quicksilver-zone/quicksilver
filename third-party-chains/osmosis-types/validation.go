package osmosistypes

import (
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	osmosislockuptypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/lockup"
	"github.com/quicksilver-zone/quicksilver/utils"
	participationrewardstypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	cl "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity"
	clmodel "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/model"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/osmomath"
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
	pool, _ := ipool.(*participationrewardstypes.OsmosisPoolProtocolData)

	poolDenom := ""
	for _, zk := range utils.Keys(pool.Denoms) {
		if pool.Denoms[zk].ChainID == chainID {
			poolDenom = zk
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

func CalculateUnderlyingAssetsFromPosition(ctx sdk.Context, position clmodel.Position, pool cl.ConcentratedPoolExtension) (sdk.Coin, sdk.Coin, error) {
	token0 := pool.GetToken0()
	token1 := pool.GetToken1()

	if position.Liquidity.IsZero() {
		return sdk.NewCoin(token0, osmomath.ZeroInt()), sdk.NewCoin(token1, osmomath.ZeroInt()), nil
	}

	// Calculate the amount of underlying assets in the position
	asset0, asset1, err := pool.CalcActualAmounts(ctx, position.LowerTick, position.UpperTick, position.Liquidity)
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	// Create coin objects from the underlying assets.
	coin0 := sdk.NewCoin(token0, asset0.TruncateInt())
	coin1 := sdk.NewCoin(token1, asset1.TruncateInt())

	return coin0, coin1, nil
}

func DetermineApplicableTokensInClPool(ctx sdk.Context, prKeeper ParticipationRewardsKeeper, position clmodel.Position, chainID string) (math.Int, error) {
	poolID := position.PoolId
	pd, ok := prKeeper.GetProtocolData(ctx, participationrewardstypes.ProtocolDataTypeOsmosisPool, fmt.Sprintf("%d", poolID))
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to obtain protocol data for poolID=%d", poolID)
	}

	ipool, err := participationrewardstypes.UnmarshalProtocolData(participationrewardstypes.ProtocolDataTypeOsmosisCLPool, pd.Data)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	pool, _ := ipool.(*participationrewardstypes.OsmosisClPoolProtocolData)

	poolDenom := ""
	for _, zk := range utils.Keys(pool.Denoms) {
		if pool.Denoms[zk].ChainID == chainID {
			poolDenom = zk
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

	asset := sdk.Coin{}
	asset0, asset1, err := CalculateUnderlyingAssetsFromPosition(ctx, position, poolData)
	if err != nil {
		return sdk.ZeroInt(), errors.New("unable to determine underlying assets for position")
	}
	switch true {
	case asset0.Denom == poolDenom:
		asset = asset0
	case asset1.Denom == poolDenom:
		asset = asset1
	default:
		return sdk.ZeroInt(), fmt.Errorf("position does not match local denom for %s", chainID)
	}

	return asset.Amount, nil
}
