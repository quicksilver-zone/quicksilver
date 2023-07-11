package crescenttypes

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/ingenuity-build/quicksilver/crescent-types/lpfarm"
	"github.com/ingenuity-build/quicksilver/utils"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type ParticipationRewardsKeeper interface {
	GetProtocolData(ctx sdk.Context, pdType participationrewardstypes.ProtocolDataType, key string) (participationrewardstypes.ProtocolData, bool)
}

func DetermineApplicableTokensInPool(ctx sdk.Context, prKeeper ParticipationRewardsKeeper, position types.Position, chainID string) (math.Int, error) {
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

	if poolData.Disabled {
		return sdk.ZeroInt(), errors.New(fmt.Sprintf("pool%d is disabled", pool.PoolID))
	}

	reserveAddress := poolData.GetReserveAddress()

	pd, ok = prKeeper.GetProtocolData(ctx, participationrewardstypes.ProtocolDataTypeCrescentReserveAddressBalance, reserveAddress.String()+poolDenom)
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to obtain reserveaddressbalance protocoldata for address=%s, denom=%s", reserveAddress.String(), poolDenom)
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
