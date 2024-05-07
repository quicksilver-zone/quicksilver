package cfmm_common

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"

	types "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/osmomath"
)

const errMsgFormatSharesLargerThanMax = "cannot exit all shares in a pool. Attempted to exit %s shares, max allowed is %s"

// CalcExitPool returns how many tokens should come out, when exiting k LP shares against a "standard" CFMM
func CalcExitPool(ctx sdk.Context, pool types.CFMMPoolI, exitingShares osmomath.Int, exitFee osmomath.Dec) (sdk.Coins, error) {
	totalShares := pool.GetTotalShares()
	if exitingShares.GTE(totalShares) {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrLimitMaxAmount, errMsgFormatSharesLargerThanMax, exitingShares, totalShares.Sub(osmomath.OneInt()))
	}

	// refundedShares = exitingShares * (1 - exit fee)
	// with 0 exit fee optimization
	var refundedShares osmomath.Dec
	if !exitFee.IsZero() {
		// exitingShares * (1 - exit fee)
		oneSubExitFee := osmomath.OneDec().Sub(exitFee)
		refundedShares = oneSubExitFee.MulIntMut(exitingShares)
	} else {
		refundedShares = exitingShares.ToLegacyDec()
	}

	shareOutRatio := refundedShares.QuoInt(totalShares)
	// exitedCoins = shareOutRatio * pool liquidity
	exitedCoins := sdk.Coins{}
	poolLiquidity := pool.GetTotalPoolLiquidity(ctx)

	for _, asset := range poolLiquidity {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.MulInt(asset.Amount).TruncateInt()
		if exitAmt.LTE(osmomath.ZeroInt()) {
			continue
		}
		if exitAmt.GTE(asset.Amount) {
			return sdk.Coins{}, errors.New("too many shares out")
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	}

	return exitedCoins, nil
}

// We binary search a number of LP shares, s.t. if we exited the pool with the updated liquidity,
// and swapped all the tokens back to the input denom, we'd get the same amount. (under 0 spread factor)
// Thanks to CFMM path-independence, we can estimate slippage with these swaps to be sure to get the right numbers here.
// (by path-independence, swap all of B -> A, and then swap all of C -> A will yield same amount of A, regardless
// of order and interleaving)
//
// This implementation requires each of pool.GetTotalPoolLiquidity, pool.ExitPool, and pool.SwapExactAmountIn
// to not update or read from state, and instead only do updates based upon the pool struct.
func BinarySearchSingleAssetJoin(
	pool types.CFMMPoolI,
	tokenIn sdk.Coin,
	poolWithAddedLiquidityAndShares func(newLiquidity sdk.Coin, newShares osmomath.Int) types.CFMMPoolI,
) (numLPShares osmomath.Int, err error) {
	// use dummy context
	ctx := sdk.Context{}
	// should be guaranteed to converge if above 256 since osmomath.Int has 256 bits
	maxIterations := 300
	// upperbound of number of LP shares = existingShares * tokenIn.Amount / pool.totalLiquidity.AmountOf(tokenIn.Denom)
	existingTokenLiquidity := pool.GetTotalPoolLiquidity(ctx).AmountOf(tokenIn.Denom)
	existingLPShares := pool.GetTotalShares()
	LPShareUpperBound := existingLPShares.Mul(tokenIn.Amount).ToLegacyDec().QuoInt(existingTokenLiquidity).Ceil().TruncateInt()
	LPShareLowerBound := osmomath.ZeroInt()

	// Creates a pool with tokenIn liquidity added, where it created `sharesIn` number of shares.
	// Returns how many tokens you'd get, if you then exited all of `sharesIn` for tokenIn.Denom
	estimateCoinOutGivenShares := func(sharesIn osmomath.Int) (tokenOut osmomath.Int, err error) {
		// new pool with added liquidity & LP shares, which we can mutate.
		poolWithUpdatedLiquidity := poolWithAddedLiquidityAndShares(tokenIn, sharesIn)
		swapToDenom := tokenIn.Denom
		// so now due to correctness of exitPool, we exitPool and swap all remaining assets to base asset
		exitFee := osmomath.ZeroDec()
		exitedCoins, err := poolWithUpdatedLiquidity.ExitPool(ctx, sharesIn, exitFee)
		if err != nil {
			return osmomath.Int{}, err
		}

		return SwapAllCoinsToSingleAsset(poolWithUpdatedLiquidity, ctx, exitedCoins, swapToDenom, osmomath.ZeroDec())
	}

	// We accept an additive tolerance of 1 LP share error and round down
	errTolerance := osmomath.ErrTolerance{AdditiveTolerance: osmomath.OneDec(), MultiplicativeTolerance: osmomath.Dec{}, RoundingDir: osmomath.RoundDown}

	numLPShares, err = osmomath.BinarySearch(
		estimateCoinOutGivenShares,
		LPShareLowerBound, LPShareUpperBound, tokenIn.Amount, errTolerance, maxIterations)
	if err != nil {
		return osmomath.Int{}, err
	}

	return numLPShares, nil
}

// SwapAllCoinsToSingleAsset iterates through each token in the input set and trades it against the same pool sequentially
func SwapAllCoinsToSingleAsset(pool types.CFMMPoolI, ctx sdk.Context, inTokens sdk.Coins, swapToDenom string,
	spreadFactor osmomath.Dec,
) (osmomath.Int, error) {
	tokenOutAmt := inTokens.AmountOfNoDenomValidation(swapToDenom)
	for _, coin := range inTokens {
		if coin.Denom == swapToDenom {
			continue
		}
		tokenOut, err := pool.SwapOutAmtGivenIn(ctx, sdk.NewCoins(coin), swapToDenom, spreadFactor)
		if err != nil {
			return osmomath.Int{}, err
		}
		tokenOutAmt = tokenOutAmt.Add(tokenOut.Amount)
	}
	return tokenOutAmt, nil
}
