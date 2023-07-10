package amm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Pool is the interface of a pool.
type Pool interface {
	Balances() (rx, ry sdk.Int)
	SetBalances(rx, ry sdk.Int, derive bool)
	PoolCoinSupply() sdk.Int
	Price() sdk.Dec
	IsDepleted() bool

	HighestBuyPrice() (sdk.Dec, bool)
	LowestSellPrice() (sdk.Dec, bool)
	BuyAmountOver(price sdk.Dec, inclusive bool) sdk.Int
	SellAmountUnder(price sdk.Dec, inclusive bool) sdk.Int
	BuyAmountTo(price sdk.Dec) sdk.Int
	SellAmountTo(price sdk.Dec) sdk.Int

	Clone() Pool
}
