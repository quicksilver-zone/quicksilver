package keeper

import (
	"sort"

	"github.com/tendermint/tendermint/libs/log"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/x/supply/types"
)

// Keeper of the mint store.
type Keeper struct {
	cdc             codec.BinaryCodec
	storeKey        storetypes.StoreKey
	accountKeeper   types.AccountKeeper
	bankKeeper      types.BankKeeper
	stakingKeeper   types.StakingKeeper
	moduleAccounts  []string
	endpointEnabled bool
}

// NewKeeper creates a new mint Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	moduleAccounts []string,
	endpointEnabled bool,
) Keeper {
	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		accountKeeper:   ak,
		bankKeeper:      bk,
		stakingKeeper:   sk,
		moduleAccounts:  moduleAccounts,
		endpointEnabled: endpointEnabled,
	}
}

// _____________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) CalculateCirculatingSupply(ctx sdk.Context, baseDenom string, excludeAddresses []string) math.Int {
	nonCirculating := math.ZeroInt()
	k.accountKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
		for _, addr := range excludeAddresses {
			if addr == account.GetAddress().String() {
				// matched excluded address
				nonCirculating = nonCirculating.Add(k.bankKeeper.GetBalance(ctx, account.GetAddress(), baseDenom).Amount)
				return false
			}
		}

		nonCirculating = nonCirculating.Add(k.bankKeeper.LockedCoins(ctx, account.GetAddress()).AmountOf(baseDenom))
		return false
	})

	for _, macc := range k.moduleAccounts {
		// exclude staking pools
		if macc != stakingtypes.BondedPoolName && macc != stakingtypes.NotBondedPoolName {
			addr := k.accountKeeper.GetModuleAddress(macc)
			maccBalance := k.bankKeeper.GetBalance(ctx, addr, baseDenom).Amount
			nonCirculating = nonCirculating.Add(maccBalance)
		}
	}

	return k.bankKeeper.GetSupply(ctx, baseDenom).Amount.Sub(nonCirculating)
}

func (k Keeper) TopN(ctx sdk.Context, baseDenom string, n uint64) []*types.Account {
	accountMap := map[string]math.Int{}

	modMap := map[string]bool{}

	for _, mod := range k.moduleAccounts {
		modMap[k.accountKeeper.GetModuleAddress(mod).String()] = true
	}

	k.accountKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
		if modMap[account.GetAddress().String()] {
			return false
		}
		balance := k.bankKeeper.GetBalance(ctx, account.GetAddress(), baseDenom).Amount
		accountMap[account.GetAddress().String()] = balance
		return false
	})

	k.stakingKeeper.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) (stop bool) {
		if modMap[delegation.GetDelegatorAddr().String()] {
			return false
		}
		balance := delegation.GetShares().TruncateInt()
		accountMap[delegation.GetDelegatorAddr().String()] = accountMap[delegation.GetDelegatorAddr().String()].Add(balance)
		return false
	})

	accountSlice := []*types.Account{}
	for addr, balance := range accountMap {
		accountSlice = append(accountSlice, &types.Account{Address: addr, Balance: balance})
	}

	sort.Slice(accountSlice, func(i, j int) bool {
		return accountSlice[i].Balance.GT(accountSlice[j].Balance)
	})

	return accountSlice[:n]
}
