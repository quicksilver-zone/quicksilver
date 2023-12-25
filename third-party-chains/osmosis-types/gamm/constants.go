package gamm

import (
	"cosmossdk.io/math"

	sdkmath "cosmossdk.io/math"
)

const (
	MinPoolAssets = 2
	MaxPoolAssets = 8

	OneShareExponent = 18
	// Raise 10 to the power of SigFigsExponent to determine number of significant figures.
	// i.e. SigFigExponent = 8 is 10^8 which is 100000000. This gives 8 significant figures.
	SigFigsExponent = 8
)

var (
	// OneShare represents the amount of subshares in a single pool share.
	OneShare = math.NewIntWithDecimal(1, OneShareExponent)

	// InitPoolSharesSupply is the amount of new shares to initialize a pool with.
	InitPoolSharesSupply = OneShare.MulRaw(100)

	// SpotPriceSigFigs is the amount of significant figures used in return value of calculate SpotPrice
	SpotPriceSigFigs = sdkmath.LegacyNewDec(10).Power(SigFigsExponent).TruncateInt()
	// MaxSpotPrice is the maximum supported spot price. Anything greater than this will error.
	MaxSpotPrice = sdkmath.LegacyNewDec(2).Power(128).Sub(sdkmath.LegacyOneDec())
)
