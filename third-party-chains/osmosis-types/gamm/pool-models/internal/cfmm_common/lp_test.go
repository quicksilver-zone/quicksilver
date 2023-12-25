package cfmm_common_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/balancer"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/internal/cfmm_common"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/stableswap"
)

// a helper function used to multiply coins
func mulCoins(coins sdk.Coins, multiplier sdkmath.LegacyDec) sdk.Coins {
	outCoins := sdk.Coins{}
	for _, coin := range coins {
		outCoin := sdk.NewCoin(coin.Denom, multiplier.MulInt(coin.Amount).TruncateInt())
		if !outCoin.Amount.IsZero() {
			outCoins = append(outCoins, outCoin)
		}
	}
	return outCoins
}

func TestCalcExitPool(t *testing.T) {
	emptyContext := sdk.Context{}

	twoStablePoolAssets := sdk.NewCoins(
		sdk.NewInt64Coin("foo", 1000000000),
		sdk.NewInt64Coin("bar", 1000000000),
	)

	threeBalancerPoolAssets := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin("foo", 2000000000), Weight: sdkmath.NewIntFromUint64(5)},
		{Token: sdk.NewInt64Coin("bar", 3000000000), Weight: sdkmath.NewIntFromUint64(5)},
		{Token: sdk.NewInt64Coin("baz", 4000000000), Weight: sdkmath.NewIntFromUint64(5)},
	}

	// create these pools used for testing
	twoAssetPool, err := stableswap.NewStableswapPool(
		1,
		stableswap.PoolParams{ExitFee: sdkmath.LegacyZeroDec()},
		twoStablePoolAssets,
		[]uint64{1, 1},
		"",
	)
	require.NoError(t, err)

	threeAssetPool, err := balancer.NewBalancerPool(
		1,
		balancer.PoolParams{SwapFee: sdkmath.LegacyZeroDec(), ExitFee: sdkmath.LegacyZeroDec()},
		threeBalancerPoolAssets,
		"",
		time.Now(),
	)
	require.NoError(t, err)

	twoAssetPoolWithExitFee, err := stableswap.NewStableswapPool(
		1,
		stableswap.PoolParams{ExitFee: sdkmath.LegacyMustNewDecFromStr("0.0001")},
		twoStablePoolAssets,
		[]uint64{1, 1},
		"",
	)
	require.NoError(t, err)

	threeAssetPoolWithExitFee, err := balancer.NewBalancerPool(
		1,
		balancer.PoolParams{SwapFee: sdkmath.LegacyZeroDec(), ExitFee: sdkmath.LegacyMustNewDecFromStr("0.0002")},
		threeBalancerPoolAssets,
		"",
		time.Now(),
	)
	require.NoError(t, err)

	tests := []struct {
		name          string
		pool          gamm.PoolI
		exitingShares sdkmath.Int
		expError      bool
	}{
		{
			name:          "two-asset pool, exiting shares grater than total shares",
			pool:          &twoAssetPool,
			exitingShares: twoAssetPool.GetTotalShares().AddRaw(1),
			expError:      true,
		},
		{
			name:          "three-asset pool, exiting shares grater than total shares",
			pool:          &threeAssetPool,
			exitingShares: threeAssetPool.GetTotalShares().AddRaw(1),
			expError:      true,
		},
		{
			name:          "two-asset pool, valid exiting shares",
			pool:          &twoAssetPool,
			exitingShares: twoAssetPool.GetTotalShares().QuoRaw(2),
			expError:      false,
		},
		{
			name:          "three-asset pool, valid exiting shares",
			pool:          &threeAssetPool,
			exitingShares: sdkmath.NewIntFromUint64(3000000000000),
			expError:      false,
		},
		{
			name:          "two-asset pool with exit fee, valid exiting shares",
			pool:          &twoAssetPoolWithExitFee,
			exitingShares: twoAssetPoolWithExitFee.GetTotalShares().QuoRaw(2),
			expError:      false,
		},
		{
			name:          "three-asset pool with exit fee, valid exiting shares",
			pool:          &threeAssetPoolWithExitFee,
			exitingShares: sdkmath.NewIntFromUint64(7000000000000),
			expError:      false,
		},
	}

	for _, test := range tests {
		// using empty context since, currently, the context is not used anyway. This might be changed in the future
		exitFee := test.pool.GetExitFee(emptyContext)
		exitCoins, err := cfmm_common.CalcExitPool(emptyContext, test.pool, test.exitingShares, exitFee)
		if test.expError {
			require.Error(t, err, "test: %v", test.name)
		} else {
			require.NoError(t, err, "test: %v", test.name)

			// exitCoins = ( (1 - exitFee) * exitingShares / poolTotalShares ) * poolTotalLiquidity
			expExitCoins := mulCoins(test.pool.GetTotalPoolLiquidity(emptyContext), (sdkmath.LegacyOneDec().Sub(exitFee)).MulInt(test.exitingShares).QuoInt(test.pool.GetTotalShares()))
			require.Equal(t, expExitCoins.Sort().String(), exitCoins.Sort().String(), "test: %v", test.name)
		}
	}
}

func TestMaximalExactRatioJoin(t *testing.T) {
	emptyContext := sdk.Context{}

	balancerPoolAsset := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin("foo", 100), Weight: sdkmath.NewIntFromUint64(5)},
		{Token: sdk.NewInt64Coin("bar", 100), Weight: sdkmath.NewIntFromUint64(5)},
	}

	tests := []struct {
		name        string
		pool        func() gamm.PoolI
		tokensIn    sdk.Coins
		expNumShare sdkmath.Int
		expRemCoin  sdk.Coins
	}{
		{
			name: "two asset pool, same tokenIn ratio",
			pool: func() gamm.PoolI {
				balancerPool, err := balancer.NewBalancerPool(
					1,
					balancer.PoolParams{SwapFee: sdkmath.LegacyZeroDec(), ExitFee: sdkmath.LegacyZeroDec()},
					balancerPoolAsset,
					"",
					time.Now(),
				)
				require.NoError(t, err)
				return &balancerPool
			},
			tokensIn:    sdk.NewCoins(sdk.NewCoin("foo", sdkmath.NewInt(10)), sdk.NewCoin("bar", sdkmath.NewInt(10))),
			expNumShare: sdkmath.NewIntFromUint64(10000000000000000000),
			expRemCoin:  sdk.Coins{},
		},
		{
			name: "two asset pool, different tokenIn ratio with pool",
			pool: func() gamm.PoolI {
				balancerPool, err := balancer.NewBalancerPool(
					1,
					balancer.PoolParams{SwapFee: sdkmath.LegacyZeroDec(), ExitFee: sdkmath.LegacyZeroDec()},
					balancerPoolAsset,
					"",
					time.Now(),
				)
				require.NoError(t, err)
				return &balancerPool
			},
			tokensIn:    sdk.NewCoins(sdk.NewCoin("foo", sdkmath.NewInt(10)), sdk.NewCoin("bar", sdkmath.NewInt(11))),
			expNumShare: sdkmath.NewIntFromUint64(10000000000000000000),
			expRemCoin:  sdk.NewCoins(sdk.NewCoin("bar", sdkmath.NewIntFromUint64(1))),
		},
	}

	for _, test := range tests {
		balancerPool, err := balancer.NewBalancerPool(
			1,
			balancer.PoolParams{SwapFee: sdkmath.LegacyZeroDec(), ExitFee: sdkmath.LegacyZeroDec()},
			balancerPoolAsset,
			"",
			time.Now(),
		)
		require.NoError(t, err)

		numShare, remCoins, err := cfmm_common.MaximalExactRatioJoin(&balancerPool, emptyContext, test.tokensIn)

		require.NoError(t, err)
		require.Equal(t, test.expNumShare, numShare)
		require.Equal(t, test.expRemCoin, remCoins)
	}
}
