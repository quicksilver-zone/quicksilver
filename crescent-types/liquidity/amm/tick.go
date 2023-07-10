package amm

import (
	"math"
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TickPrecision represents a tick precision.
type TickPrecision int

func (prec TickPrecision) PriceToDownTick(price sdk.Dec) sdk.Dec {
	return PriceToDownTick(price, int(prec))
}

func (prec TickPrecision) PriceToUpTick(price sdk.Dec) sdk.Dec {
	return PriceToUpTick(price, int(prec))
}

func (prec TickPrecision) UpTick(price sdk.Dec) sdk.Dec {
	return UpTick(price, int(prec))
}

func (prec TickPrecision) DownTick(price sdk.Dec) sdk.Dec {
	return DownTick(price, int(prec))
}

func (prec TickPrecision) HighestTick() sdk.Dec {
	return HighestTick(int(prec))
}

func (prec TickPrecision) LowestTick() sdk.Dec {
	return LowestTick(int(prec))
}

func (prec TickPrecision) TickToIndex(tick sdk.Dec) int {
	return TickToIndex(tick, int(prec))
}

func (prec TickPrecision) TickFromIndex(i int) sdk.Dec {
	return TickFromIndex(i, int(prec))
}

func (prec TickPrecision) RoundPrice(price sdk.Dec) sdk.Dec {
	return RoundPrice(price, int(prec))
}

func (prec TickPrecision) TickGap(price sdk.Dec) sdk.Dec {
	return TickGap(price, int(prec))
}

func (prec TickPrecision) RandomTick(r *rand.Rand, minPrice, maxPrice sdk.Dec) sdk.Dec {
	return RandomTick(r, minPrice, maxPrice, int(prec))
}

// char returns the characteristic(integral part) of
// log10(x * pow(10, sdk.Precision)).
func char(x sdk.Dec) int {
	if x.IsZero() {
		panic("cannot calculate log10 for 0")
	}
	return len(x.BigInt().Text(10)) - 1
}

// pow10 returns pow(10, n - sdk.Precision).
func pow10(n int) sdk.Dec {
	x := big.NewInt(10)
	x.Exp(x, big.NewInt(int64(n)), nil)
	return sdk.NewDecFromBigIntWithPrec(x, sdk.Precision)
}

// isPow10 returns whether x is a power of 10 or not.
func isPow10(x sdk.Dec) bool {
	b := x.BigInt()
	if b.Sign() <= 0 {
		return false
	}
	ten := big.NewInt(10)
	if b.Cmp(ten) == 0 {
		return true
	}
	zero := big.NewInt(0)
	m := new(big.Int)
	for b.Cmp(ten) >= 0 {
		b.DivMod(b, ten, m)
		if m.Cmp(zero) != 0 {
			return false
		}
	}
	return b.Cmp(big.NewInt(1)) == 0
}

// PriceToDownTick returns the highest price tick under(or equal to) the price.
func PriceToDownTick(price sdk.Dec, prec int) sdk.Dec {
	b := price.BigInt()
	l := char(price)
	d := int64(l - prec)
	if d > 0 {
		p := big.NewInt(10)
		p.Exp(p, big.NewInt(d), nil)
		b.Quo(b, p).Mul(b, p)
	}
	return sdk.NewDecFromBigIntWithPrec(b, sdk.Precision)
}

// PriceToUpTick returns the lowest price tick greater or equal than
// the price.
func PriceToUpTick(price sdk.Dec, prec int) sdk.Dec {
	tick := PriceToDownTick(price, prec)
	if !tick.Equal(price) {
		return UpTick(tick, prec)
	}
	return tick
}

// UpTick returns the next lowest price tick above the price.
func UpTick(price sdk.Dec, prec int) sdk.Dec {
	tick := PriceToDownTick(price, prec)
	if tick.Equal(price) {
		l := char(price)
		return price.Add(pow10(l - prec))
	}
	l := char(tick)
	return tick.Add(pow10(l - prec))
}

// DownTick returns the next highest price tick under the price.
// DownTick doesn't check if the price is the lowest price tick.
func DownTick(price sdk.Dec, prec int) sdk.Dec {
	tick := PriceToDownTick(price, prec)
	if tick.Equal(price) {
		l := char(price)
		var d sdk.Dec
		if isPow10(price) {
			d = pow10(l - prec - 1)
		} else {
			d = pow10(l - prec)
		}
		return price.Sub(d)
	}
	return tick
}

// HighestTick returns the highest possible price tick.
func HighestTick(prec int) sdk.Dec {
	i := big.NewInt(2)
	// Maximum 315 bits possible, but take slightly less value for safety.
	i.Exp(i, big.NewInt(300), nil).Sub(i, big.NewInt(1))
	return PriceToDownTick(sdk.NewDecFromBigIntWithPrec(i, sdk.Precision), prec)
}

// LowestTick returns the lowest possible price tick.
func LowestTick(prec int) sdk.Dec {
	return sdk.NewDecWithPrec(1, int64(sdk.Precision-prec))
}

// TickToIndex returns a tick index for given price.
// Tick index 0 means the lowest possible price fit in ticks.
func TickToIndex(price sdk.Dec, prec int) int {
	b := price.BigInt()
	l := len(b.Text(10)) - 1
	d := int64(l - prec)
	if d > 0 {
		q := big.NewInt(10)
		q.Exp(q, big.NewInt(d), nil)
		b.Quo(b, q)
	}
	p := int(math.Pow10(prec))
	b.Sub(b, big.NewInt(int64(p)))
	return (l-prec)*9*p + int(b.Int64())
}

// TickFromIndex returns a price for given tick index.
// See TickToIndex for more details about tick indices.
func TickFromIndex(i, prec int) sdk.Dec {
	p := int(math.Pow10(prec))
	l := i/(9*p) + prec
	t := big.NewInt(int64(p + i%(p*9)))
	if l > prec {
		m := big.NewInt(10)
		m.Exp(m, big.NewInt(int64(l-prec)), nil)
		t.Mul(t, m)
	}
	return sdk.NewDecFromBigIntWithPrec(t, sdk.Precision)
}

// RoundTickIndex returns rounded tick index using banker's rounding.
func RoundTickIndex(i int) int {
	return (i + 1) / 2 * 2
}

// RoundPrice returns rounded price using banker's rounding.
func RoundPrice(price sdk.Dec, prec int) sdk.Dec {
	tick := PriceToDownTick(price, prec)
	if price.Equal(tick) {
		return price
	}
	return TickFromIndex(RoundTickIndex(TickToIndex(tick, prec)), prec)
}

// TickGap returns tick gap at given price.
func TickGap(price sdk.Dec, prec int) sdk.Dec {
	tick := PriceToDownTick(price, prec)
	l := char(tick)
	return pow10(l - prec)
}

// RandomTick returns a random tick within range [minPrice, maxPrice].
// If prices are not on ticks, then prices are adjusted to the nearest
// ticks.
func RandomTick(r *rand.Rand, minPrice, maxPrice sdk.Dec, prec int) sdk.Dec {
	minPrice = PriceToUpTick(minPrice, prec)
	maxPrice = PriceToDownTick(maxPrice, prec)
	minPriceIdx := TickToIndex(minPrice, prec)
	maxPriceIdx := TickToIndex(maxPrice, prec)
	return TickFromIndex(minPriceIdx+r.Intn(maxPriceIdx-minPriceIdx), prec)
}
