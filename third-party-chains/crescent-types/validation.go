package crescenttypes

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	lpfarm "github.com/ingenuity-build/quicksilver/third-party-chains/crescent-types/lpfarm"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type ParticipationRewardsKeeper interface {
	GetProtocolData(ctx sdk.Context, pdType participationrewardstypes.ProtocolDataType, key string) (participationrewardstypes.ProtocolData, bool)
}

func DetermineApplicableTokensInPool(ctx sdk.Context, prKeeper ParticipationRewardsKeeper, position lpfarm.Position) (math.Int, error) {
	farmingAmount := position.FarmingAmount

	poolID := position.Denom[4:]
	pd, ok := prKeeper.GetProtocolData(ctx, participationrewardstypes.ProtocolDataTypeCrescentPool, poolID)
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to obtain crescent protocol data for poolID=%s", poolID)
	}

	ipool, err := participationrewardstypes.UnmarshalProtocolData(participationrewardstypes.ProtocolDataTypeCrescentPool, pd.Data)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	pool, _ := ipool.(*participationrewardstypes.CrescentPoolProtocolData)

	if pool.Denom == "" {
		return sdk.ZeroInt(), fmt.Errorf("invalid poolDenom")
	}

	poolData, err := pool.GetPool()
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if poolData.Disabled {
		return sdk.ZeroInt(), fmt.Errorf("pool%d is disabled", pool.PoolID)
	}

	reserveAddress := poolData.GetReserveAddress()

	pd, ok = prKeeper.GetProtocolData(ctx, participationrewardstypes.ProtocolDataTypeCrescentReserveAddressBalance, fmt.Sprintf("%s_%s", reserveAddress, pool.Denom))
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to obtain reserveaddressbalance protocoldata for address=%s, denom=%s", reserveAddress, pool.Denom)
	}
	ibalance, err := participationrewardstypes.UnmarshalProtocolData(participationrewardstypes.ProtocolDataTypeCrescentReserveAddressBalance, pd.Data)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	balance, _ := ibalance.(*participationrewardstypes.CrescentReserveAddressBalanceProtocolData)

	poolDenomBalance, err := balance.GetBalance()
	if err != nil {
		return sdk.ZeroInt(), err
	}

	pd, ok = prKeeper.GetProtocolData(ctx, participationrewardstypes.ProtocolDataTypeCrescentPoolCoinSupply, poolData.PoolCoinDenom)
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to obtain poolcoinsupply protocoldata for denom=%s", poolData.PoolCoinDenom)
	}
	isupply, err := participationrewardstypes.UnmarshalProtocolData(participationrewardstypes.ProtocolDataTypeCrescentPoolCoinSupply, pd.Data)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	supply, _ := isupply.(*participationrewardstypes.CrescentPoolCoinSupplyProtocolData)

	// calculate user PoolCoin ratio and LP asset amount
	poolSupply, err := supply.GetSupply() // total poolcoin supply
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if poolSupply.IsZero() {
		return sdk.ZeroInt(), fmt.Errorf("empty pool, %s", poolID)
	}
	uratio := sdk.NewDecFromInt(farmingAmount).QuoInt(poolSupply)

	uAmount := uratio.MulInt(poolDenomBalance).TruncateInt()

	return uAmount, nil
}
