package keeper

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	config "github.com/ingenuity-build/quicksilver/cmd/config"
	osmosistypes "github.com/ingenuity-build/quicksilver/osmosis-types"
	"github.com/ingenuity-build/quicksilver/utils"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	epochskeeper "github.com/ingenuity-build/quicksilver/x/epochs/keeper"
	icqkeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// UserAllocation is an internal keeper struct to track transient state for
// rewards distribution. It contains the user address and the coins that are
// allocated to it.
type UserAllocation struct {
	Address string
	Amount  sdkmath.Int
}

var _ osmosistypes.ParticipationRewardsKeeper = &Keeper{}

type Keeper struct {
	cdc                  codec.BinaryCodec
	storeKey             storetypes.StoreKey
	paramSpace           paramtypes.Subspace
	accountKeeper        types.AccountKeeper
	bankKeeper           types.BankKeeper
	stakingKeeper        types.StakingKeeper
	IcqKeeper            icqkeeper.Keeper
	icsKeeper            *icskeeper.Keeper
	epochsKeeper         epochskeeper.Keeper
	feeCollectorName     string
	prSubmodules         map[cmtypes.ClaimType]Submodule
	ValidateProofOps     utils.ProofOpsFn
	ValidateSelfProofOps utils.SelfProofOpsFn
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
	icqk icqkeeper.Keeper,
	icsk *icskeeper.Keeper,
	feeCollectorName string,
	proofValidationFn utils.ProofOpsFn,
	selfProofValidationFn utils.SelfProofOpsFn,
) *Keeper {
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:                  cdc,
		storeKey:             key,
		paramSpace:           ps,
		accountKeeper:        ak,
		bankKeeper:           bk,
		stakingKeeper:        sk,
		IcqKeeper:            icqk,
		icsKeeper:            icsk,
		feeCollectorName:     feeCollectorName,
		prSubmodules:         LoadSubmodules(),
		ValidateProofOps:     proofValidationFn,
		ValidateSelfProofOps: selfProofValidationFn,
	}
}

func (k *Keeper) GetGovAuthority(_ sdk.Context) string {
	return sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), k.accountKeeper.GetModuleAddress(govtypes.ModuleName))
}

func (k *Keeper) SetEpochsKeeper(epochsKeeper epochskeeper.Keeper) {
	k.epochsKeeper = epochsKeeper
}

// GetParams returns the total set of participationrewards parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of participationrewards parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k *Keeper) GetClaimsEnabled(ctx sdk.Context) bool {
	var out bool
	k.paramSpace.Get(ctx, types.KeyClaimsEnabled, &out)
	return out
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) GetCodec() codec.BinaryCodec {
	return k.cdc
}

func (k *Keeper) UpdateSelfConnectionData(ctx sdk.Context) error {
	selfConnectionData := types.ConnectionProtocolData{
		ConnectionID: types.SelfConnection,
		ChainID:      ctx.ChainID(),
		LastEpoch:    ctx.BlockHeight() - 2, // reason why -2 works here.
		Prefix:       config.Bech32Prefix,
	}
	selfConnectionDataBlob, err := json.Marshal(selfConnectionData)
	if err != nil {
		k.Logger(ctx).Info("Error Marshalling self connection Data")
		return err
	}

	data := types.ProtocolData{
		Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
		Data: selfConnectionDataBlob,
	}
	k.Logger(ctx).Info("Setting self protocol data", "data", data)
	k.SetProtocolData(ctx, selfConnectionData.GenerateKey(), &data)

	return nil
}

func (k *Keeper) GetModuleBalance(ctx sdk.Context) sdkmath.Int {
	denom := k.stakingKeeper.BondDenom(ctx)
	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := k.bankKeeper.GetBalance(ctx, moduleAddress, denom)

	k.Logger(ctx).Info("module account", "address", moduleAddress, "balance", moduleBalance)

	return moduleBalance.Amount
}

func LoadSubmodules() map[cmtypes.ClaimType]Submodule {
	out := make(map[cmtypes.ClaimType]Submodule, 0)
	out[cmtypes.ClaimTypeLiquidToken] = &LiquidTokensModule{}
	out[cmtypes.ClaimTypeOsmosisPool] = &OsmosisModule{}
	out[cmtypes.ClaimTypeUmeeToken] = &UmeeModule{}
	return out
}
