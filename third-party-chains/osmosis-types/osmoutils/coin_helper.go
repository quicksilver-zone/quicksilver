package osmoutils

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubDecCoinArrays subtracts the contents of the second param from the first (decCoinsArrayA - decCoinsArrayB)
// Note that this takes in two _arrays_ of DecCoins, meaning that each term itself is of type DecCoins (i.e. an array of DecCoin).
func SubDecCoinArrays(decCoinsArrayA []sdk.DecCoins, decCoinsArrayB []sdk.DecCoins) ([]sdk.DecCoins, error) {
	if len(decCoinsArrayA) != len(decCoinsArrayB) {
		return []sdk.DecCoins{}, errors.New("DecCoin arrays must be of equal length to be subtracted")
	}

	finalDecCoinArray := []sdk.DecCoins{}
	for i := range decCoinsArrayA {
		finalDecCoinArray = append(finalDecCoinArray, decCoinsArrayA[i].Sub(decCoinsArrayB[i]))
	}

	return finalDecCoinArray, nil
}

// SafeSubDecCoinArrays subtracts the contents of the second param from the first (decCoinsArrayA - decCoinsArrayB)
// Note that this takes in two _arrays_ of DecCoins, meaning that each term itself is of type DecCoins (i.e. an array of DecCoin).
// Contrary to SubDecCoinArrays, this subtractions allows for negative result values.
func SafeSubDecCoinArrays(decCoinsArrayA []sdk.DecCoins, decCoinsArrayB []sdk.DecCoins) ([]sdk.DecCoins, error) {
	if len(decCoinsArrayA) != len(decCoinsArrayB) {
		return []sdk.DecCoins{}, errors.New("DecCoin arrays must be of equal length to be subtracted")
	}

	finalDecCoinArray := []sdk.DecCoins{}
	for i := range decCoinsArrayA {
		subResult, _ := decCoinsArrayA[i].SafeSub(decCoinsArrayB[i])
		finalDecCoinArray = append(finalDecCoinArray, subResult)
	}

	return finalDecCoinArray, nil
}

// AddDecCoinArrays adds the contents of the second param from the first (decCoinsArrayA + decCoinsArrayB)
// Note that this takes in two _arrays_ of DecCoins, meaning that each term itself is of type DecCoins (i.e. an array of DecCoin).
func AddDecCoinArrays(decCoinsArrayA []sdk.DecCoins, decCoinsArrayB []sdk.DecCoins) ([]sdk.DecCoins, error) {
	if len(decCoinsArrayA) != len(decCoinsArrayB) {
		return []sdk.DecCoins{}, errors.New("DecCoin arrays must be of equal length to be added")
	}

	finalDecCoinArray := []sdk.DecCoins{}
	for i := range decCoinsArrayA {
		finalDecCoinArray = append(finalDecCoinArray, decCoinsArrayA[i].Add(decCoinsArrayB[i]...))
	}

	return finalDecCoinArray, nil
}

// CollapseDecCoinsArray takes an array of DecCoins and returns the sum of all the DecCoins in the array.
func CollapseDecCoinsArray(decCoinsArray []sdk.DecCoins) sdk.DecCoins {
	finalDecCoins := sdk.DecCoins{}
	for _, decCoins := range decCoinsArray {
		finalDecCoins = finalDecCoins.Add(decCoins...)
	}
	return finalDecCoins
}

// ConvertCoinsToDecCoins takes sdk.Coins and converts it to sdk.DecCoins
func ConvertCoinsToDecCoins(coins sdk.Coins) sdk.DecCoins {
	decCoins := sdk.DecCoins{}
	for _, coin := range coins {
		decCoins = append(decCoins, sdk.NewDecCoin(coin.Denom, coin.Amount))
	}
	return decCoins
}

// FilterDenoms returns the coins with only the passed in denoms
func FilterDenoms(coins sdk.Coins, denoms []string) sdk.Coins {
	filteredCoins := sdk.NewCoins()

	for _, denom := range denoms {
		filteredCoins = filteredCoins.Add(sdk.NewCoin(denom, coins.AmountOf(denom)))
	}

	return filteredCoins
}

// MergeCoinMaps takes two maps of type map[T]sdk.Coins and merges them together, adding the values of the second map to the first.
func MergeCoinMaps[T comparable](currentEpochExpectedDistributionsOne map[T]sdk.Coins, poolIDToExpectedDistributionMapOne map[T]sdk.Coins) map[T]sdk.Coins {
	newMap := map[T]sdk.Coins{}

	// Iterate over the first map and add all the values to the new map
	for poolID, expectedDistribution := range currentEpochExpectedDistributionsOne {
		newMap[poolID] = expectedDistribution
	}

	// Iterate over the second map and add all the values to the new map
	for poolID, expectedDistribution := range poolIDToExpectedDistributionMapOne {
		if _, ok := newMap[poolID]; ok {
			newMap[poolID] = newMap[poolID].Add(expectedDistribution...)
		} else {
			newMap[poolID] = expectedDistribution
		}
	}
	return newMap
}

func ConvertCoinArrayToCoins(coinArray []sdk.Coin) sdk.Coins {
	coins := sdk.Coins{}
	for _, coin := range coinArray {
		coins = append(coins, coin)
	}
	return coins
}
