package simtypes

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"time"
	"unsafe"

	"golang.org/x/exp/constraints"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// shamelessly copied from
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang#31832326

// RandStringOfLength generates a random string of a particular length.
func RandStringOfLength(r *rand.Rand, n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// RandPositiveInt get a rand positive sdkmath.Int.
func RandPositiveInt(r *rand.Rand, maxPositiveInt sdkmath.Int) (sdkmath.Int, error) {
	if !maxPositiveInt.GTE(sdk.OneInt()) {
		return sdkmath.Int{}, errors.New("maxPositiveInt too small")
	}

	maxPositiveInt = maxPositiveInt.Sub(sdk.OneInt())

	return sdk.NewIntFromBigInt(new(big.Int).Rand(r, maxPositiveInt.BigInt())).Add(sdkmath.OneInt()), nil
}

// RandomAmount generates a random amount.
// Note: The range of RandomAmount includes max, and is, in fact, biased to return max as well as 0.
func RandomAmount(r *rand.Rand, maxAmount sdkmath.Int) sdkmath.Int {
	randInt := big.NewInt(0)

	switch r.Intn(10) {
	case 0:
		// do nothing
		// randInt = big.NewInt(0)
	case 1:
		randInt = maxAmount.BigInt()
	default: // NOTE: there are 10 total cases.
		randInt = big.NewInt(0).Rand(r, maxAmount.BigInt()) // up to maxAmount - 1
	}

	return sdk.NewIntFromBigInt(randInt)
}

// RandomDecAmount generates a random decimal amount
// Note: The range of RandomDecAmount includes max, and is, in fact, biased to return max as well as 0.
func RandomDecAmount(r *rand.Rand, maxDec sdk.Dec) sdk.Dec {
	randInt := big.NewInt(0)

	switch r.Intn(10) {
	case 0:
		return sdk.NewDecFromBigIntWithPrec(randInt, sdk.Precision)
	case 1:
		randInt = maxDec.BigInt() // the underlying big int with all precision bits.
	default: // NOTE: there are 10 total cases.
		randInt = big.NewInt(0).Rand(r, maxDec.BigInt())
	}

	return sdk.NewDecFromBigIntWithPrec(randInt, sdk.Precision)
}

// RandTimestamp generates a random timestamp.
func RandTimestamp(r *rand.Rand) time.Time {
	// json.Marshal breaks for timestamps greater with year greater than 9999
	unixTime := r.Int63n(253373529600)
	return time.Unix(unixTime, 0)
}

// RandIntBetween returns a random int between two numbers inclusively.
func RandIntBetween(r *rand.Rand, minInt, maxInt int) int {
	return r.Intn(maxInt-minInt) + minInt
}

// RandSubsetCoins returns random subset of the provided coins
// will return at least one coin unless coins argument is empty or malformed
// i.e. 0 amt in coins.
func RandSubsetCoins(r *rand.Rand, coins sdk.Coins) sdk.Coins {
	if len(coins) == 0 {
		return sdk.Coins{}
	}
	// make sure at least one coin added
	denomIdx := r.Intn(len(coins))
	coin := coins[denomIdx]
	amt, err := RandPositiveInt(r, coin.Amount)
	// malformed coin. 0 amt in coins
	if err != nil {
		return sdk.Coins{}
	}

	subset := sdk.Coins{sdk.NewCoin(coin.Denom, amt)}

	for i, c := range coins {
		// skip denom that we already chose earlier
		if i == denomIdx {
			continue
		}
		// coin flip if multiple coins
		// if there is single coin then return random amount of it
		if r.Intn(2) == 0 && len(coins) != 1 {
			continue
		}

		amt, err := RandPositiveInt(r, c.Amount)
		// ignore errors and try another denom
		if err != nil {
			continue
		}

		subset = append(subset, sdk.NewCoin(c.Denom, amt))
	}

	return subset.Sort()
}

func RandCoin(r *rand.Rand, coins sdk.Coins) sdk.Coins {
	if len(coins) == 0 {
		return sdk.Coins{}
	}
	// make sure at least one coin added
	denomIdx := r.Intn(len(coins))
	coin := coins[denomIdx]
	amt, err := RandPositiveInt(r, coin.Amount)
	// malformed coin. 0 amt in coins
	if err != nil {
		return sdk.Coins{}
	}

	return sdk.Coins{sdk.NewCoin(coin.Denom, amt)}
}

// RandExponentialCoin uniformly samples a denom from the addr's balances.
// Then it samples an Exponentially distributed amount of the addr's coins, with rate = 10.
// (Meaning that on average it samples 10% of the chosen balance)
// Pre-condition: Addr must have a spendable balance.
func RandExponentialCoin(r *rand.Rand, coin sdk.Coin) sdk.Coin {
	lambda := float64(10)
	sample := r.ExpFloat64() / lambda
	// truncate exp at 1, which will only be reached in .0045% of the time.
	// .000045 ~= (1 - CDF(1, Exp[\lambda=10])) = e^{-10}
	sample = math.Min(1, sample)
	// Do some hacky scaling to get this into an SDK decimal,
	// were going to treat it as an integer in the range [0, 10000]
	maxRange := int64(10000)
	intSample := int64(math.Round(sample * float64(maxRange)))
	newAmount := coin.Amount.MulRaw(intSample).QuoRaw(maxRange)
	return sdk.NewCoin(coin.Denom, newAmount)
}

func RandLTBound[T constraints.Integer](r *rand.Rand, upperbound T) T {
	return RandLTEBound(r, upperbound-T(1))
}

func RandLTEBound[T constraints.Integer](r *rand.Rand, upperbound T) T {
	maxUB := int(upperbound)
	i := RandIntBetween(r, 0, maxUB+1)
	t := T(i)
	return t
}

func RandSelect[T interface{}](r *rand.Rand, args ...T) T {
	choice := RandLTBound(r, len(args))
	return args[choice]
}
