package keeper

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cometbft/cometbft/libs/log"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/third-party-chains/gaia-types/liquid/types"
	epochskeeper "github.com/quicksilver-zone/quicksilver/x/epochs/keeper"
	interchainquerykeeper "github.com/quicksilver-zone/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

type TxSubmitFn func(ctx sdk.Context, k *Keeper, msgs []sdk.Msg, account *types.ICAAccount, memo string, messagesPerTx int64) error

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc                 codec.Codec
	storeKey            storetypes.StoreKey
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICQKeeper           interchainquerykeeper.Keeper
	AccountKeeper       types.AccountKeeper
	AuthzKeeper         types.AuthzKeeper
	BankKeeper          types.BankKeeper
	IBCKeeper           *ibckeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	ClaimsManagerKeeper types.ClaimsManagerKeeper
	EpochsKeeper        types.EpochsKeeper
	Ir                  codectypes.InterfaceRegistry
	hooks               types.IcsHooks
	paramStore          paramtypes.Subspace
	txSubmit            TxSubmitFn
	govAuthority        string
}

// NewKeeper returns a new instance of zones Keeper.
// This function will panic on failure.
func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	authzKeeper types.AuthzKeeper,
	bankKeeper types.BankKeeper,
	icaControllerKeeper icacontrollerkeeper.Keeper,
	icqKeeper interchainquerykeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper,
	transferKeeper ibctransferkeeper.Keeper,
	claimsManagerKeeper types.ClaimsManagerKeeper,
	ps paramtypes.Subspace,
	govAuthority string,
) *Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if addr := accountKeeper.GetModuleAddress(types.EscrowModuleAccount); addr == nil {
		panic(fmt.Sprintf("%s escrow account has not been set", types.EscrowModuleAccount))
	}

	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	if ibcKeeper == nil {
		panic("ibcKeeper is nil")
	}

	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		ICAControllerKeeper: icaControllerKeeper,
		ICQKeeper:           icqKeeper,
		BankKeeper:          bankKeeper,
		AccountKeeper:       accountKeeper,
		IBCKeeper:           ibcKeeper,
		TransferKeeper:      transferKeeper,
		ClaimsManagerKeeper: claimsManagerKeeper,
		hooks:               nil,
		txSubmit:            ProdSubmitTx,
		paramStore:          ps,
		AuthzKeeper:         authzKeeper,
		govAuthority:        govAuthority,
	}
}

func (k *Keeper) OverrideTxSubmit(fn TxSubmitFn) {
	k.txSubmit = fn
}

// SetHooks set the ics hooks.
func (k *Keeper) SetHooks(icsh types.IcsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set epochs hooks twice")
	}

	k.hooks = icsh

	return k
}

func (k *Keeper) GetGovAuthority(_ sdk.Context) string {
	return k.govAuthority
}

func (k *Keeper) SetEpochsKeeper(epochsKeeper epochskeeper.Keeper) {
	k.EpochsKeeper = &epochsKeeper
}

// Logger returns a module-specific logger.
func (*Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) GetCodec() codec.Codec {
	return k.cdc
}

func (k *Keeper) SetConnectionForPort(ctx sdk.Context, connectionID, port string) {
	mapping := types.PortConnectionTuple{ConnectionId: connectionID, PortId: port}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortMapping)
	bz := k.cdc.MustMarshal(&mapping)
	store.Set([]byte(port), bz)
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

func (k *Keeper) DeleteConnectionForPort(ctx sdk.Context, port string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortMapping)
	store.Delete([]byte(port))
}

// IteratePortConnections iterates through all of the delegations.
func (k *Keeper) IteratePortConnections(ctx sdk.Context, cb func(pc types.PortConnectionTuple) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixPortMapping)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pc := types.PortConnectionTuple{}
		k.cdc.MustUnmarshal(iterator.Value(), &pc)
		if cb(pc) {
			break
		}
	}
}

// AllPortConnections returns all delegations used during genesis dump.
func (k *Keeper) AllPortConnections(ctx sdk.Context) (pcs []types.PortConnectionTuple) {
	k.IteratePortConnections(ctx, func(pc types.PortConnectionTuple) bool {
		pcs = append(pcs, pc)
		return false
	})

	return pcs
}

// GetZoneByLocalDenom returns zone by denom.
func (k *Keeper) GetZoneByLocalDenom(ctx sdk.Context, denom string) *types.Zone {
	zone, existed := k.GetLocalDenomZoneMapping(ctx, denom)
	if existed {
		return zone
	}
	// If get zone from denom zone mapping not found, find it in zones list end set into the mapping
	k.IterateZones(ctx, func(_ int64, thisZone *types.Zone) bool {
		if thisZone.LocalDenom == denom {
			zone = thisZone
			k.SetLocalDenomZoneMapping(ctx, thisZone)
			return true
		}
		return false
	})
	return zone
}

// GetLocalDenomZoneMapping returns zone by denom.
func (k *Keeper) GetLocalDenomZoneMapping(ctx sdk.Context, denom string) (*types.Zone, bool) {
	zone := types.Zone{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLocalDenomZoneMapping)
	bz := store.Get([]byte(denom))
	if len(bz) == 0 {
		return nil, false
	}

	k.cdc.MustUnmarshal(bz, &zone)
	return &zone, true
}

// SetLocalDenomZoneMapping set denom <-> zone mapping.
func (k *Keeper) SetLocalDenomZoneMapping(ctx sdk.Context, zone *types.Zone) {
	if zone == nil {
		return
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLocalDenomZoneMapping)
	bz := k.cdc.MustMarshal(zone)
	store.Set([]byte(zone.LocalDenom), bz)
}

// DeleteDenomZoneMapping delete zone info in denom - zone mapping.
func (k *Keeper) DeleteDenomZoneMapping(ctx sdk.Context, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLocalDenomZoneMapping)
	store.Delete([]byte(denom))
}

// ### Interval functions >>>
// * some of these functions (or portions thereof) may be changed to single
//   query type functions, dependent upon callback features / capabilities;

func (k *Keeper) SetValidatorsForZone(ctx sdk.Context, data []byte, icqQuery icqtypes.Query) error {
	return nil
}

func (k *Keeper) SetValidatorForZone(ctx sdk.Context, zone *types.Zone, data []byte) error {
	return nil
}

func (k *Keeper) GetParam(ctx sdk.Context, key []byte) uint64 {
	var out uint64
	k.paramStore.Get(ctx, key, &out)
	return out
}

func (k *Keeper) GetUnbondingEnabled(ctx sdk.Context) bool {
	var out bool
	k.paramStore.Get(ctx, types.KeyUnbondingEnabled, &out)
	return out
}

func (k *Keeper) GetCommissionRate(ctx sdk.Context) sdk.Dec {
	var out sdk.Dec
	k.paramStore.Get(ctx, types.KeyCommissionRate, &out)
	return out
}

func (k *Keeper) GetAuthzAutoClaimAddress(ctx sdk.Context) string {
	var out string
	k.paramStore.Get(ctx, types.KeyAuthzAutoClaimAddress, &out)
	return out
}

// MigrateParams fetches params, adds ClaimsEnabled field and re-sets params.
func (k *Keeper) MigrateParams(ctx sdk.Context) {
	params := types.Params{}
	params.DepositInterval = k.GetParam(ctx, types.KeyDepositInterval)
	params.CommissionRate = k.GetCommissionRate(ctx)
	params.ValidatorsetInterval = k.GetParam(ctx, types.KeyValidatorSetInterval)
	params.UnbondingEnabled = false
	params.AuthzAutoClaimAddress = k.GetAuthzAutoClaimAddress(ctx)

	k.paramStore.SetParamSet(ctx, &params)
}

func (k *Keeper) GetParams(clientCtx sdk.Context) (params types.Params) {
	k.paramStore.GetParamSet(clientCtx, &params)
	return params
}

// SetParams sets the interchainstaking parameters to the param space.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}

func (k *Keeper) GetChainID(ctx sdk.Context, connectionID string) (string, error) {
	conn, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return "", fmt.Errorf("invalid connection id, %q not found", connectionID)
	}
	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, conn.ClientId)
	if !found {
		return "", fmt.Errorf("client id %q not found for connection %q", conn.ClientId, connectionID)
	}
	client, ok := clientState.(*ibctmtypes.ClientState)
	if !ok {
		return "", fmt.Errorf("invalid client state for client %q on connection %q", conn.ClientId, connectionID)
	}

	return client.ChainId, nil
}

func (k *Keeper) EmitPerformanceBalanceQuery(ctx sdk.Context, zone *types.Zone) error {
	_, addr, err := bech32.DecodeAndConvert(zone.PerformanceAddress.Address)
	if err != nil {
		return err
	}
	data := banktypes.CreateAccountBalancesPrefix(addr)

	// query performance account for baseDenom balance every 100 blocks.
	k.ICQKeeper.MakeRequest(
		ctx,
		zone.ConnectionId,
		zone.ChainId,
		types.BankStoreKey,
		append(data, zone.BaseDenom...),
		sdk.NewInt(-1),
		types.ModuleName,
		"perfbalance",
		100,
	)

	zone.PerformanceAddress.BalanceWaitgroup = 1
	k.SetZone(ctx, zone)

	return nil
}

func (k *Keeper) EmitValSetQuery(ctx sdk.Context, connectionID, chainID string, validatorsReq stakingtypes.QueryValidatorsRequest, period sdkmath.Int) error {
	bz, err := k.cdc.Marshal(&validatorsReq)
	if err != nil {
		return errors.New("failed to marshal valset pagination request")
	}

	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		period,
		types.ModuleName,
		"valset",
		0,
	)

	return nil
}

func (k *Keeper) EmitValidatorQuery(ctx sdk.Context, connectionID, chainID string, validator stakingtypes.Validator) error {
	_, addr, err := bech32.DecodeAndConvert(validator.OperatorAddress)
	if err != nil {
		return fmt.Errorf("EmitValidatorQuery failed to decode validator.OperatorAddress: %q got error: %w",
			validator.OperatorAddress, err)
	}
	data := stakingtypes.GetValidatorKey(addr)
	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"store/staking/key",
		data,
		sdk.NewInt(-1),
		types.ModuleName,
		"validator",
		0,
	)
	return nil
}

func (k *Keeper) EmitDepositIntervalQuery(ctx sdk.Context, zone *types.Zone) {
	req := tx.GetTxsEventRequest{
		Events: []string{
			"transfer.recipient='" + zone.DepositAddress.GetAddress() + "'",
		},
		OrderBy: tx.OrderBy_ORDER_BY_DESC,
		Pagination: &query.PageRequest{
			Limit: types.TxRetrieveCount,
		},
	}

	k.ICQKeeper.MakeRequest(
		ctx,
		zone.ConnectionId,
		zone.ChainId,
		"cosmos.tx.v1beta1.Service/GetTxsEvent",
		k.cdc.MustMarshal(&req),
		sdk.NewInt(-1),
		types.ModuleName,
		"depositinterval",
		0,
	)
}

func (k *Keeper) EmitLiquidValidatorQuery(ctx sdk.Context, connectionID, chainID string, addr sdk.ValAddress) {
	data := lsmstakingtypes.GetLiquidValidatorKey(addr)
	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"store/liquid/key",
		data,
		sdk.NewInt(-1),
		types.ModuleName,
		"lsminfo",
		0,
	)
}

func (k *Keeper) EmitSigningInfoQuery(ctx sdk.Context, connectionID, chainID string, validator stakingtypes.Validator) error {
	consAddress, err := validator.GetConsAddr()
	if err != nil {
		return err
	}

	data := slashingtypes.ValidatorSigningInfoKey(consAddress)
	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		chainID,
		"store/slashing/key",
		data,
		sdk.NewInt(-1),
		types.ModuleName,
		"signinginfo",
		0,
	)

	return nil
}

func (k *Keeper) GetDelegationsInProcess(ctx sdk.Context, chainID string) sdkmath.Int {
	delegationsInProcess := sdkmath.ZeroInt()
	k.IterateZoneReceipts(ctx, chainID, func(_ int64, receipt types.Receipt) (stop bool) {
		if receipt.Completed == nil {
			for _, coin := range receipt.Amount {
				delegationsInProcess = delegationsInProcess.Add(coin.Amount) // we cannot simply choose
			}
		}
		return false
	})
	return delegationsInProcess
}

// redemption rate

func (k *Keeper) UpdateRedemptionRate(ctx sdk.Context, zone *types.Zone, epochRewards sdkmath.Int) {
	return
	// delegationsInProcess := k.GetDelegationsInProcess(ctx, zone.ChainId)
	// ratio, isZero := k.GetRatio(ctx, zone, epochRewards.Add(delegationsInProcess))
	// k.Logger(ctx).Info("Redemption Rate Update", "chain", zone.ChainId, "epochly_rewards", epochRewards, "last_rate", zone.LastRedemptionRate, "current_rate", zone.RedemptionRate, "new_rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, zone).Amount.Add(epochRewards).Add(delegationsInProcess))

	// // TODO: make max deltas params.
	// // soft cap redemption rate, instead of panicking.
	// delta := ratio.Quo(zone.RedemptionRate)
	// if delta.GT(sdk.NewDecWithPrec(102, 2)) {
	// 	k.Logger(ctx).Error("ratio diverged by more than 2% upwards in the last epoch; capping at 1.02...")
	// 	ratio = zone.RedemptionRate.Mul(sdk.NewDecWithPrec(102, 2))
	// } else if delta.LT(sdk.NewDecWithPrec(95, 2)) && !isZero { // we allow a bigger downshift if all assets were withdrawn and we revert to zero.
	// 	k.Logger(ctx).Error("ratio diverged by more than 5% downwards in the last epoch; 5% is the theoretical max if _all_ controlled tokens were tombstoned. capping at 0.95...")
	// 	ratio = zone.RedemptionRate.Mul(sdk.NewDecWithPrec(95, 2))
	// }

	// zone.LastRedemptionRate = zone.RedemptionRate
	// zone.RedemptionRate = ratio
	// k.SetZone(ctx, zone)
}

// UnmarshalValidatorsResponse attempts to umarshal a byte slice into a QueryValidatorsResponse.
func (k *Keeper) UnmarshalValidatorsResponse(data []byte) (stakingtypes.QueryValidatorsResponse, error) {
	validatorsRes := stakingtypes.QueryValidatorsResponse{}
	if len(data) == 0 {
		return validatorsRes, errors.New("attempted to unmarshal zero length byte slice (8)")
	}
	err := k.cdc.Unmarshal(data, &validatorsRes)
	if err != nil {
		return validatorsRes, err
	}

	return validatorsRes, nil
}

// UnmarshalValidatorsRequest attempts to umarshal  a byte slice into a QueryValidatorsRequest.
func (k *Keeper) UnmarshalValidatorsRequest(data []byte) (stakingtypes.QueryValidatorsRequest, error) {
	validatorsReq := stakingtypes.QueryValidatorsRequest{}
	err := k.cdc.Unmarshal(data, &validatorsReq)
	if err != nil {
		return validatorsReq, err
	}

	return validatorsReq, nil
}

// UnmarshalValidator attempts to umarshal  a byte slice into a Validator.
func (k *Keeper) UnmarshalValidator(data []byte) (stakingtypes.Validator, error) {
	validator := stakingtypes.Validator{}
	if len(data) == 0 {
		return validator, errors.New("attempted to unmarshal zero length byte slice (9)")
	}
	err := k.cdc.Unmarshal(data, &validator)
	if err != nil {
		return validator, err
	}

	return validator, nil
}
