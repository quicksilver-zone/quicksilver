package keeper

import (
	"fmt"
	"github.com/quicksilver-zone/quicksilver/utils"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	osmosistypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types"
	umeetypes "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
)

var (
	_ osmosistypes.ClaimsManagerKeeper = &Keeper{}
	_ umeetypes.ClaimsManagerKeeper    = &Keeper{}
)

type Keeper struct {
	cdc                  codec.BinaryCodec
	storeKey             storetypes.StoreKey
	IBCKeeper            *ibckeeper.Keeper
	paramSpace           paramtypes.Subspace
	IcqKeeper            types.InterchainQueryKeeper
	icsKeeper            types.InterchainStakingKeeper
	PrSubmodules         map[types.ClaimType]Submodule
	ValidateProofOps     utils.ProofOpsFn
	ValidateSelfProofOps utils.SelfProofOpsFn
}

// NewKeeper returns a new instance of participationrewards Keeper.
// This function will panic on failure.
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	ibcKeeper *ibckeeper.Keeper,
	ps paramtypes.Subspace,
	icsk types.InterchainStakingKeeper,
	icqk types.InterchainQueryKeeper,
	proofValidationFn utils.ProofOpsFn,
	selfProofValidationFn utils.SelfProofOpsFn,
) Keeper {
	if ibcKeeper == nil {
		panic("ibcKeeper is nil")
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:                  cdc,
		storeKey:             key,
		IBCKeeper:            ibcKeeper,
		paramSpace:           ps,
		icsKeeper:            icsk,
		IcqKeeper:            icqk,
		PrSubmodules:         LoadSubmodules(),
		ValidateProofOps:     proofValidationFn,
		ValidateSelfProofOps: selfProofValidationFn,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreSelfConsensusState(ctx sdk.Context, key string) error {
	var height ibcclienttypes.Height

	blockHeight := ctx.BlockHeight() - 1
	if blockHeight < 0 {
		return fmt.Errorf("block height is negative: %d", blockHeight)
	}

	if strings.Contains(ctx.ChainID(), "-") {
		chainParts := strings.Split(ctx.ChainID(), "-")
		revisionNum, err := strconv.ParseUint(chainParts[len(chainParts)-1], 10, 64)
		if err != nil {
			k.Logger(ctx).Error("Error getting revision number for client ", "chainID", ctx.ChainID())
			return err
		}

		height = ibcclienttypes.Height{
			RevisionNumber: revisionNum,
			RevisionHeight: uint64(blockHeight),
		}
	} else {
		// ONLY FOR TESTING - ibctesting module chains donot follow standard [chainname]-[num] structure
		height = ibcclienttypes.Height{
			RevisionNumber: 0, // revision number for testchain1 is 0 (because parseChainId splits on '-')
			RevisionHeight: uint64(blockHeight),
		}
	}

	selfConsState, err := k.IBCKeeper.ClientKeeper.GetSelfConsensusState(ctx, height)
	if err != nil {
		k.Logger(ctx).Error("Error getting self consensus state of previous height")
		return err
	}

	state, _ := selfConsState.(*ibctmtypes.ConsensusState)
	k.SetSelfConsensusState(ctx, key, state)

	return nil
}

func (k Keeper) GetClaimsEnabled(ctx sdk.Context) bool {
	var out bool
	k.paramSpace.Get(ctx, types.KeyClaimsEnabled, &out)
	return out
}

func LoadSubmodules() map[types.ClaimType]Submodule {
	out := make(map[types.ClaimType]Submodule, 0)
	out[types.ClaimTypeLiquidToken] = &LiquidTokensModule{}
	out[types.ClaimTypeOsmosisPool] = &OsmosisModule{}
	out[types.ClaimTypeOsmosisCLPool] = &OsmosisClModule{}
	out[types.ClaimTypeUmeeToken] = &UmeeModule{}
	return out
}
