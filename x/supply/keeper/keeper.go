package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/x/supply/types"
)

// Keeper of the mint store.
type Keeper struct {
	cdc            codec.BinaryCodec
	storeKey       storetypes.StoreKey
	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	stakingKeeper  types.StakingKeeper
	moduleAccounts []string
	baseDenom      string
}

// NewKeeper creates a new mint Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	moduleAccounts []string,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		accountKeeper:  ak,
		bankKeeper:     bk,
		stakingKeeper:  sk,
		moduleAccounts: moduleAccounts,
	}
}

// _____________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) CalculateCirculatingSupply(ctx sdk.Context, excludeAddresses []string) math.Int {
	baseDenom := k.stakingKeeper.BondDenom(ctx)
	// Creates context with current height and checks txs for ctx to be usable by start of next block
	nonCirculating := math.ZeroInt()
	k.accountKeeper.IterateAccounts(ctx, func(account sdk.AccountI) (stop bool) {
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
			fmt.Println(macc, maccBalance)
			nonCirculating = nonCirculating.Add(maccBalance)
		}
	}

	return k.bankKeeper.GetSupply(ctx, k.baseDenom).Amount.Sub(nonCirculating)
}
