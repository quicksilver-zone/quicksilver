package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	paramSpace    paramtypes.Subspace
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	govKeeper     types.GovKeeper
	ibcKeeper     *ibckeeper.Keeper
	icsKeeper     types.InterchainStakingKeeper
	prKeeper      types.ParticipationRewardsKeeper

	ValidateProofOps utils.ProofOpsFn

	// the address capable of executing authority-scoped messages (ex. params, props). Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper returns a new instance of participationrewards Keeper.
// This function will panic on failure.
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	gk types.GovKeeper,
	ibcKeeper *ibckeeper.Keeper,
	icsk types.InterchainStakingKeeper,
	prk types.ParticipationRewardsKeeper,
	pofn utils.ProofOpsFn,
	authority string,
) *Keeper {
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	if ibcKeeper == nil {
		panic("ibcKeeper is nil")
	}

	return &Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       ps,
		accountKeeper:    ak,
		bankKeeper:       bk,
		stakingKeeper:    sk,
		govKeeper:        gk,
		ibcKeeper:        ibcKeeper,
		icsKeeper:        icsk,
		prKeeper:         prk,
		ValidateProofOps: pofn,
		authority:        authority,
	}
}

// GetAuthority returns the x/airdrop module's authority.
func (k *Keeper) GetAuthority() string {
	return k.authority
}

// GetParams returns the total set of airdrop parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of airdrop parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetModuleAccountAddress gets the airdrop module account address.
func (k *Keeper) GetModuleAccountAddress(_ sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// GetModuleAccountBalance gets the airdrop module account coin balance.
func (k *Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coin {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	return k.bankKeeper.GetBalance(ctx, moduleAccAddr, k.stakingKeeper.BondDenom(ctx))
}

func (k *Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amount)
}

func (k *Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAccount sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAccount, amount)
}

func (k *Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAccount sdk.AccAddress, recipientModule string, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAccount, recipientModule, amount)
}

func (k *Keeper) BondDenom(ctx sdk.Context) string {
	return k.stakingKeeper.BondDenom(ctx)
}
