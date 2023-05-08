package simtypes

import (
	"errors"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
)

func RandomSimAccount(r *rand.Rand, accs []simulation.Account) simulation.Account {
	acc, _ := simulation.RandomAcc(r, accs)
	return acc
}

type SimAccountConstraint = func(account simulation.Account) bool

// returns acc, accExists := sim.RandomSimAccountWithConstraint(f)
// where acc is a uniformly sampled account from all accounts satisfying the constraint f
// a constraint is satisfied for an account `acc` if f(acc) = true
// accExists is false, if there is no such account.
func RandomSimAccountWithConstraint(r *rand.Rand, f SimAccountConstraint, accs []simulation.Account) (simulation.Account, bool) {
	filteredAddrs := []simulation.Account{}
	for _, acc := range accs {
		if f(acc) {
			filteredAddrs = append(filteredAddrs, acc)
		}
	}

	if len(filteredAddrs) == 0 {
		return simulation.Account{}, false
	}
	return RandomSimAccount(r, filteredAddrs), true
}

func RandomSimAccountWithMinCoins(ctx sdk.Context, r *rand.Rand, accs []simulation.Account, coins sdk.Coins, bk BankKeeper) (simulation.Account, error) {
	accHasMinCoins := func(acc simulation.Account) bool {
		spendableCoins := bk.SpendableCoins(ctx, acc.Address)
		return spendableCoins.IsAllGTE(coins) && coins.DenomsSubsetOf(spendableCoins)
	}
	acc, found := RandomSimAccountWithConstraint(r, accHasMinCoins, accs)
	if !found {
		return simulation.Account{}, errors.New("no address with min balance found")
	}
	return acc, nil
}

func RandomExistingAddress(r *rand.Rand, accs []simulation.Account) sdk.AccAddress {
	acc := RandomSimAccount(r, accs)
	return acc.Address
}

func AddAccount(acc simulation.Account, accs []simulation.Account) []simulation.Account {
	if _, found := FindAccount(acc.Address, accs); !found {
		return append(accs, acc)
	}
	return accs
}

// FindAccount iterates over all the simulation accounts to find the one that matches
// the given address
// TODO: Benchmark time in here, we should probably just make a hashmap indexing this.
func FindAccount(address sdk.Address, accs []simulation.Account) (simulation.Account, bool) {
	for _, acc := range accs {
		if acc.Address.Equals(address) {
			return acc, true
		}
	}

	return simulation.Account{}, false
}

func RandomSimAccountWithBalance(ctx sdk.Context, r *rand.Rand, accs []simulation.Account, bk BankKeeper) (simulation.Account, error) {
	accHasBal := func(acc simulation.Account) bool {
		return len(bk.SpendableCoins(ctx, acc.Address)) != 0
	}
	acc, found := RandomSimAccountWithConstraint(r, accHasBal, accs)
	if !found {
		return simulation.Account{}, errors.New("no address with balance found. Check simulator configuration, this should be very rare")
	}
	return acc, nil
}

// Returns (account, randSubsetCoins, found), so if found = false, then no such address exists.
// randSubsetCoins is a random subset of the provided denoms, if the account is found.
// TODO: Write unit test.
func SelAddrWithDenoms(ctx sdk.Context, r *rand.Rand, accs []simulation.Account, denoms []string, bk BankKeeper) (simulation.Account, sdk.Coins, bool) {
	accHasDenoms := func(acc simulation.Account) bool {
		for _, denom := range denoms {
			if bk.GetBalance(ctx, acc.Address, denom).Amount.IsZero() {
				return false
			}
			// only return addr if it has spendable coins of requested denom
			coins := bk.SpendableCoins(ctx, acc.Address)
			for _, coin := range coins {
				if denom == coin.Denom {
					return true
				}
			}
		}
		return true
	}

	acc, accExists := RandomSimAccountWithConstraint(r, accHasDenoms, accs)
	if !accExists {
		return acc, sdk.Coins{}, false
	}
	balance, err := RandCoinSubsetFromBalance(ctx, r, acc.Address, denoms, bk)
	if err != nil {
		return acc, sdk.Coins{}, false
	}

	return acc, balance.Sort(), true
}

// SelAddrWithDenom attempts to find an address with the provided denom. This function
// returns (account, randSubsetCoins, found), so if found = false, then no such address exists.
// randSubsetCoins is a random subset of the provided denoms, if the account is found.
// TODO: Write unit test.
func SelAddrWithDenom(ctx sdk.Context, r *rand.Rand, accs []simulation.Account, denom string, bk BankKeeper) (simulation.Account, sdk.Coin, bool) {
	acc, subsetCoins, found := SelAddrWithDenoms(ctx, r, accs, []string{denom}, bk)
	if !found {
		return acc, sdk.Coin{}, found
	}
	return acc, subsetCoins[0], found
}

// GetRandSubsetOfKDenoms returns a random subset of coins of k unique denoms from the provided account
// TODO: Write unit test.
func GetRandSubsetOfKDenoms(ctx sdk.Context, r *rand.Rand, acc simulation.Account, k int, bk BankKeeper) (sdk.Coins, bool) {
	// get all spendable coins from provided account
	coins := bk.SpendableCoins(ctx, acc.Address)
	// ensure account coins are greater than or equal to the requested subset length
	if len(coins) < k {
		return sdk.Coins{}, false
	}

	for len(coins) != k {
		index := r.Intn(len(coins) - 1)
		coins = RemoveIndex(coins, index)
	}
	// append random amount less than or equal to existing amount to new subset array
	subset := sdk.Coins{}
	for _, c := range coins {
		amt, err := simulation.RandPositiveInt(r, c.Amount)
		if err != nil {
			return sdk.Coins{}, false
		}
		subset = append(subset, sdk.NewCoin(c.Denom, amt))
	}

	// return nothing if the coin struct length is less than requested (sanity check)
	if len(subset) < k {
		return sdk.Coins{}, false
	}

	return subset.Sort(), true
}

// RandomSimAccountWithKDenoms returns an account that possesses k unique denoms.
func RandomSimAccountWithKDenoms(ctx sdk.Context, r *rand.Rand, accs []simulation.Account, k int, bk BankKeeper) (simulation.Account, bool) {
	accHasBal := func(acc simulation.Account) bool {
		return len(bk.SpendableCoins(ctx, acc.Address)) >= k
	}
	return RandomSimAccountWithConstraint(r, accHasBal, accs)
}

// RandExponentialCoinFromBalance uniformly samples a denom from the addr's balances.
// Then it samples an Exponentially distributed amount of the addr's coins, with rate = 10.
// (Meaning that on average it samples 10% of the chosen balance)
// Pre-condition: Addr must have a spendable balance.
func RandExponentialCoinFromBalance(ctx sdk.Context, r *rand.Rand, addr sdk.AccAddress, bk BankKeeper) sdk.Coin {
	balances := bk.SpendableCoins(ctx, addr)
	if len(balances) == 0 {
		panic("precondition for RandExponentialCoin broken: Addr has 0 spendable balance")
	}

	coin := RandCoin(r, balances)
	// TODO: Reconsider if this becomes problematic in the future, but currently thinking it
	// should be fine for simulation.
	return RandExponentialCoin(r, coin[0])
}

func RandCoinSubsetFromBalance(ctx sdk.Context, r *rand.Rand, addr sdk.AccAddress, denoms []string, bk BankKeeper) (sdk.Coins, error) {
	subsetCoins := sdk.Coins{}
	for _, denom := range denoms {
		coins := bk.SpendableCoins(ctx, addr)
		for _, coin := range coins {
			if denom == coin.Denom {
				amt, err := RandPositiveInt(r, coin.Amount)
				if err != nil {
					return nil, err
				}
				subsetCoins = subsetCoins.Add(sdk.NewCoin(coin.Denom, amt))
			}
		}
	}
	return subsetCoins, nil
}

// RandomFees returns a random fee by selecting a random coin denomination and
// amount from the account's available balance. If the user doesn't have enough
// funds for paying fees, it returns empty coins.
func RandomFees(r *rand.Rand, spendableCoins sdk.Coins) (sdk.Coins, error) {
	if spendableCoins.Empty() {
		return nil, nil
	}

	// TODO: Revisit this
	perm := r.Perm(len(spendableCoins))
	var randCoin sdk.Coin
	for _, index := range perm {
		randCoin = spendableCoins[index]
		if !randCoin.Amount.IsZero() {
			break
		}
	}

	if randCoin.Amount.IsZero() {
		return nil, fmt.Errorf("no coins found for random fees")
	}

	amt, err := RandPositiveInt(r, randCoin.Amount)
	if err != nil {
		return nil, err
	}

	// Create a random fee and verify the fees are within the account's spendable
	// balance.
	fees := sdk.NewCoins(sdk.NewCoin(randCoin.Denom, amt))

	return fees, nil
}

func RemoveIndex(s sdk.Coins, index int) sdk.Coins {
	return append(s[:index], s[index+1:]...)
}
