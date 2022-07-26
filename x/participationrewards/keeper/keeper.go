package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/tendermint/tendermint/libs/log"

	icqkeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         storetypes.StoreKey
	paramSpace       paramtypes.Subspace
	accountKeeper    authkeeper.AccountKeeper
	bankKeeper       bankkeeper.Keeper
	stakingKeeper    stakingkeeper.Keeper
	icqKeeper        icqkeeper.Keeper
	icsKeeper        icskeeper.Keeper
	feeCollectorName string
}

// NewKeeper returns a new instance of participationrewards Keeper
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	sk stakingkeeper.Keeper,
	icqk icqkeeper.Keeper,
	icsk icskeeper.Keeper,
	feeCollectorName string,
) Keeper {
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       ps,
		accountKeeper:    ak,
		bankKeeper:       bk,
		stakingKeeper:    sk,
		icqKeeper:        icqk,
		icsKeeper:        icsk,
		feeCollectorName: feeCollectorName,
	}
}

// GetParams returns the total set of participationrewards parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of participationrewards parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAllocation(ctx sdk.Context, balance sdk.Coin, portion sdk.Dec) sdk.Coin {
	return sdk.NewCoin(balance.Denom, sdk.NewDecFromInt(balance.Amount).Mul(portion).TruncateInt())
}
