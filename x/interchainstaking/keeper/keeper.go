package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v4/modules/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"

	interchainquerykeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc                 codec.Codec
	storeKey            storetypes.StoreKey
	scopedKeeper        capabilitykeeper.ScopedKeeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICQKeeper           interchainquerykeeper.Keeper
	AccountKeeper       authKeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	IBCKeeper           ibckeeper.Keeper
	paramStore          paramtypes.Subspace
}

// NewKeeper returns a new instance of zones Keeper
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, accountKeeper authKeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, icacontrollerkeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper, icqKeeper interchainquerykeeper.Keeper, ibcKeeper ibckeeper.Keeper, ps paramtypes.Subspace) Keeper {
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

func (k Keeper) GetCodec() codec.Codec {
	return k.cdc
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
	err := k.cdc.Unmarshal(data, &validatorsRes)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validators info for zone", "zone", zoneInfo.ChainId, "err", err)
		return err
	}

	for _, validator := range validatorsRes.Validators {
		_, addr, _ := bech32.DecodeAndConvert(validator.OperatorAddress)
		val, err := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
		if err != nil {
			k.Logger(ctx).Info("Unable to find validator - fetching proof...", "valoper", validator.OperatorAddress)

			data := stakingTypes.GetValidatorKey(addr)
			k.ICQKeeper.MakeRequest(
				ctx,
				zoneInfo.ConnectionId,
				zoneInfo.ChainId,
				"store/staking/key",
				data,
				sdk.NewInt(-1),
				types.ModuleName,
				"validator",
				0,
			)
			continue
		}

		if !val.CommissionRate.Equal(validator.GetCommission()) || !val.VotingPower.Equal(validator.Tokens) || !val.DelegatorShares.Equal(validator.DelegatorShares) {
			k.Logger(ctx).Info("Validator state change; fetching proof", "valoper", validator.OperatorAddress)

			if err != nil {
				return err
			}

			data := stakingTypes.GetValidatorKey(addr)
			k.ICQKeeper.MakeRequest(
				ctx,
				zoneInfo.ConnectionId,
				zoneInfo.ChainId,
				"store/staking/key",
				data,
				sdk.NewInt(-1),
				types.ModuleName,
				"validator",
				0,
			)
		}
	}

	// also do this for Unbonded and Unbonding
	k.SetRegisteredZone(ctx, zoneInfo)
	return nil
}

func SetValidatorForZone(k Keeper, ctx sdk.Context, zoneInfo types.RegisteredZone, data []byte) error {
	validator := stakingTypes.Validator{}
	err := k.cdc.Unmarshal(data, &validator)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validator info for zone", "zone", zoneInfo.ChainId, "err", err)
		return err
	}

	val, err := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
	if err != nil {
		k.Logger(ctx).Info("Unable to find validator - adding...", "valoper", validator.OperatorAddress)

		zoneInfo.Validators = append(zoneInfo.Validators, &types.Validator{
			ValoperAddress:  validator.OperatorAddress,
			CommissionRate:  validator.GetCommission(),
			VotingPower:     validator.Tokens,
			DelegatorShares: validator.DelegatorShares,
			Score:           sdk.ZeroDec(),
		})
		zoneInfo.Validators = zoneInfo.GetValidatorsSorted()

	} else {

		if !val.CommissionRate.Equal(validator.GetCommission()) {
			val.CommissionRate = validator.GetCommission()
			k.Logger(ctx).Info("Validator commission rate change; updating...", "valoper", validator.OperatorAddress, "oldRate", val.CommissionRate, "newRate", validator.GetCommission())
		}

		if !val.VotingPower.Equal(validator.Tokens) {
			val.VotingPower = validator.Tokens
			k.Logger(ctx).Info("Validator voting power change; updating", "valoper", validator.OperatorAddress, "oldPower", val.VotingPower, "newPower", validator.Tokens)
		}

		if !val.DelegatorShares.Equal(validator.DelegatorShares) {
			val.DelegatorShares = validator.DelegatorShares
			k.Logger(ctx).Info("Validator delegator shares change; updating", "valoper", validator.OperatorAddress, "oldShares", val.DelegatorShares, "newShares", validator.DelegatorShares)
		}
	}

	k.SetRegisteredZone(ctx, zoneInfo)
	return nil
}

func (k Keeper) depositInterval(ctx sdk.Context) zoneItrFn {
	return func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
		if zoneInfo.DepositAddress != nil {
			if !zoneInfo.DepositAddress.Balance.Empty() {
				k.Logger(ctx).Info("balance is non zero", "balance", zoneInfo.DepositAddress.Balance)

				req := tx.GetTxsEventRequest{Events: []string{"transfer.recipient='" + zoneInfo.DepositAddress.GetAddress() + "'"}, Pagination: &query.PageRequest{Limit: types.TxRetrieveCount, Reverse: true}}
				k.ICQKeeper.MakeRequest(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.tx.v1beta1.Service/GetTxsEvent", k.cdc.MustMarshal(&req), sdk.NewInt(-1), types.ModuleName, "depositinterval", 0)

			}
		} else {
			k.Logger(ctx).Error("deposit account is nil")
		}
		return false
	}
}

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

func (k Keeper) GetChainIdFromContext(ctx sdk.Context) (string, error) {
	connectionId := ctx.Context().Value("connectionId")
	if connectionId == nil {
		return "", fmt.Errorf("connectionId not in context")
	}

	return k.GetChainID(ctx, connectionId.(string))
}

func (k Keeper) EmitPerformanceBalanceQuery(ctx sdk.Context, zone *types.RegisteredZone) error {
	balanceQuery := bankTypes.QueryAllBalancesRequest{Address: zone.PerformanceAddress.Address}
	bz, err := k.GetCodec().Marshal(&balanceQuery)
	if err != nil {
		return err
	}

	k.ICQKeeper.MakeRequest(
		ctx,
		zone.ConnectionId,
		zone.ChainId,
		"cosmos.bank.v1beta1.Query/AllBalances",
		bz,
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		"perfbalance",
		0,
	)

	return nil
}
