package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	interchainquerykeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc                 codec.Codec
	storeKey            sdk.StoreKey
	scopedKeeper        capabilitykeeper.ScopedKeeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICQKeeper           interchainquerykeeper.Keeper
	AccountKeeper       authKeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	IBCKeeper           ibckeeper.Keeper
	paramStore          paramtypes.Subspace
}

// NewKeeper returns a new instance of zones Keeper
func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, accountKeeper authKeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, icacontrollerkeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper, icqKeeper interchainquerykeeper.Keeper, ibcKeeper ibckeeper.Keeper, ps paramtypes.Subspace) Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		scopedKeeper:        scopedKeeper,
		ICAControllerKeeper: icacontrollerkeeper,
		ICQKeeper:           icqKeeper,
		BankKeeper:          bankKeeper,
		AccountKeeper:       accountKeeper,
		IBCKeeper:           ibcKeeper,
		paramStore:          ps,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ClaimCapability claims the channel capability passed via the OnOpenChanInit callback
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (k *Keeper) SetConnectionForPort(ctx sdk.Context, connectionId string, port string) error {
	mapping := types.PortConnectionTuple{ConnectionId: connectionId, PortId: port}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortMapping)
	bz := k.cdc.MustMarshal(&mapping)
	store.Set([]byte(port), bz)
	return nil
}

func (k *Keeper) GetConnectionForPort(ctx sdk.Context, port string) (string, error) {
	mapping := types.PortConnectionTuple{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortMapping)
	bz := store.Get([]byte(port))
	if len(bz) == 0 {
		return "", fmt.Errorf("unable to find mapping for port %s", port)
	}

	k.cdc.MustUnmarshal(bz, &mapping)
	return mapping.ConnectionId, nil
}

// ### Interval functions >>>
// * some of these functions (or portions thereof) may be changed to single
//   query type functions, dependent upon callback features / capabilities;

func SetValidatorsForZone(k Keeper, ctx sdk.Context, zoneInfo types.RegisteredZone, data []byte) error {
	validatorsRes := stakingTypes.QueryValidatorsResponse{}
	err := k.cdc.UnmarshalJSON(data, &validatorsRes)
	if err != nil {
		k.Logger(ctx).Error("Unable to unmarshal validators info for zone", "zone", zoneInfo.ChainId, "err", err)
		return err
	}

	for _, validator := range validatorsRes.Validators {
		val, err := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
		if err != nil {
			k.Logger(ctx).Info("Unable to find validator - adding...", "valoper", validator.OperatorAddress)
			zoneInfo.Validators = append(zoneInfo.GetValidatorsSorted(), &types.Validator{
				ValoperAddress: validator.OperatorAddress,
				CommissionRate: validator.GetCommission(),
				VotingPower:    sdk.NewDecFromInt(validator.Tokens),
				Delegations:    []*types.Delegation{},
			})
			continue
		}

		if !val.CommissionRate.Equal(validator.GetCommission()) {
			val.CommissionRate = validator.GetCommission()
			k.Logger(ctx).Info("Validator commission rate change; updating...", "valoper", validator.OperatorAddress, "oldRate", val.CommissionRate, "newRate", validator.GetCommission())
		}

		if !val.VotingPower.Equal(sdk.NewDecFromInt(validator.Tokens)) {
			val.VotingPower = sdk.NewDecFromInt(validator.Tokens)
			k.Logger(ctx).Info("Validator voting power change; updating", "valoper", validator.OperatorAddress, "oldPower", val.VotingPower, "newPower", validator.Tokens.ToDec())
		}
	}

	// also do this for Unbonded and Unbonding
	k.SetRegisteredZone(ctx, zoneInfo)
	return nil
}

func (k Keeper) depositInterval(ctx sdk.Context) zoneItrFn {
	return func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
		if zoneInfo.DepositAddress != nil {
			if !zoneInfo.DepositAddress.Balance.Empty() {
				k.Logger(ctx).Info("Balance is non zero", "balance", zoneInfo.DepositAddress.Balance)

				var callback Callback = func(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
					txs := coretypes.ResultTxSearch{}
					err := json.Unmarshal(args, &txs)
					if err != nil {
						k.Logger(ctx).Error("Unable to unmarshal txs for deposit account", "deposit_address", zoneInfo.DepositAddress.GetAddress(), "err", err)
						return err
					}

					for _, tx := range txs.Txs {
						k.HandleReceiptTransaction(ctx, tx, zoneInfo)
					}
					return nil
				}

				k.ICQKeeper.MakeRequest(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.tx.v1beta1.Query/GetTxEvents", map[string]string{"transfer.recipient": zoneInfo.DepositAddress.GetAddress()}, sdk.NewInt(-1), types.ModuleName, callback)

			}
		} else {
			k.Logger(ctx).Error("Deposit account is nil")
		}
		return false
	}
}

// // temporary: this callback should be registered when the delegate account is created, in ibc_module.go but is currently here to
// // avoid a testnet restart (same logic, just called in a different way).

// func (k Keeper) delegateInterval(ctx sdk.Context) zoneItrFn {
// 	return func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
// 		for _, ica := range zoneInfo.DelegationAddresses {
// 			// emit a single balance query for each delegate account

// 			var cb Callback = func(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
// 				zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
// 				if !found {
// 					return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
// 				}
// 				return k.SetAccountBalance(ctx, zone, query.QueryParameters["address"], args)
// 			}

// 			k.ICQKeeper.MakeRequest(
// 				ctx,
// 				zoneInfo.ConnectionId,
// 				zoneInfo.ChainId,
// 				"cosmos.bank.v1beta1.Query/AllBalances",
// 				map[string]string{"address": ica.Address},
// 				sdk.NewInt(-1),
// 				types.ModuleName,
// 				cb,
// 			)
// 		}
// 		return false
// 	}
// }

func (k *Keeper) GetParam(ctx sdk.Context, key []byte) uint64 {
	var out uint64
	k.paramStore.Get(ctx, key, &out)
	return out
}

func (k *Keeper) GetCommissionRate(ctx sdk.Context) sdk.Dec {
	var out sdk.Dec
	k.paramStore.Get(ctx, types.KeyCommissionRate, &out)
	return out
}

func (k Keeper) GetParams(clientCtx sdk.Context) (params types.Params) {
	k.paramStore.GetParamSet(clientCtx, &params)
	return params
}

// SetParams sets the distribution parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}

func (k Keeper) GetChainID(ctx sdk.Context, connectionID string) (string, error) {
	conn, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return "", fmt.Errorf("invalid connection id, \"%s\" not found", connectionID)
	}
	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, conn.ClientId)
	if !found {
		return "", fmt.Errorf("client id \"%s\" not found for connection \"%s\"", conn.ClientId, connectionID)
	}
	client, ok := clientState.(*ibctmtypes.ClientState)
	if !ok {
		return "", fmt.Errorf("invalid client state for client \"%s\" on connection \"%s\"", conn.ClientId, connectionID)
	}

	return client.ChainId, nil
}
