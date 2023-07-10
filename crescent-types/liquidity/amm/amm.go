package amm

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The minimum and maximum coin amount used in the amm package.
var (
	MinCoinAmount = sdk.NewInt(100)
	MaxCoinAmount = math.NewIntWithDecimal(1, 40)
)

var (
	MinPoolPrice               = sdk.NewDecWithPrec(1, 15)                        // 10^-15
	MaxPoolPrice               = sdk.NewDecFromInt(math.NewIntWithDecimal(1, 20)) // 10^20
	MinRangedPoolPriceGapRatio = sdk.NewDecWithPrec(1, 3)                         // 0.001, 0.1%
)
