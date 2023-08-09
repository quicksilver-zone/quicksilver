package keeper

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	interchainquerykeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc                 codec.Codec
	storeKey            storetypes.StoreKey
	scopedKeeper        *capabilitykeeper.ScopedKeeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICQKeeper           interchainquerykeeper.Keeper
	AccountKeeper       types.AccountKeeper
	BankKeeper          types.BankKeeper
	IBCKeeper           *ibckeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	ClaimsManagerKeeper types.ClaimsManagerKeeper
	Ir                  codectypes.InterfaceRegistry
	hooks               types.IcsHooks
	paramStore          paramtypes.Subspace
	msgRouter           types.MessageRouter
	authority           string
}

// NewKeeper returns a new instance of zones Keeper.
// This function will panic on failure.
func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	icaControllerKeeper icacontrollerkeeper.Keeper,
	scopedKeeper *capabilitykeeper.ScopedKeeper,
	icqKeeper interchainquerykeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper,
	transferKeeper ibctransferkeeper.Keeper,
	claimsManagerKeeper types.ClaimsManagerKeeper,
	ps paramtypes.Subspace,
	msgRouter types.MessageRouter,
	authority string,
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
		scopedKeeper:        scopedKeeper,
		ICAControllerKeeper: icaControllerKeeper,
		ICQKeeper:           icqKeeper,
		BankKeeper:          bankKeeper,
		AccountKeeper:       accountKeeper,
		IBCKeeper:           ibcKeeper,
		TransferKeeper:      transferKeeper,
		ClaimsManagerKeeper: claimsManagerKeeper,
		hooks:               nil,

		paramStore: ps,
		msgRouter:  msgRouter,
		authority:  authority,
	}
}

// SetHooks set the ics hooks.
func (k *Keeper) SetHooks(icsh types.IcsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set epochs hooks twice")
	}

	k.hooks = icsh

	return k
}

func (k *Keeper) GetGovAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) GetCodec() codec.Codec {
	return k.cdc
}

func (k *Keeper) ScopedKeeper() *capabilitykeeper.ScopedKeeper {
	return k.scopedKeeper
}

// ClaimCapability claims the channel capability passed via the OnOpenChanInit callback.
func (k *Keeper) ClaimCapability(ctx sdk.Context, capability *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, capability, name)
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

// ### Interval functions >>>
// * some of these functions (or portions thereof) may be changed to single
//   query type functions, dependent upon callback features / capabilities;

func (k *Keeper) SetValidatorsForZone(ctx sdk.Context, data []byte, icqQuery icqtypes.Query) error {
	zone, found := k.GetZone(ctx, icqQuery.ChainId)
	if !found {
		k.Logger(ctx).Error("unable find zone for query", "zone", icqQuery.ChainId)
		return fmt.Errorf("unable to find zone %s for query %s", icqQuery.ChainId, icqQuery.Id)
	}

	validatorsRes, err := k.UnmarshalValidatorsResponse(data)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validators info for zone", "zone", zone.ZoneID(), "err", err)
		return err
	}

	if validatorsRes.Pagination != nil && !bytes.Equal(validatorsRes.Pagination.NextKey, []byte{}) {
		validatorsReq, err := k.UnmarshalValidatorsRequest(icqQuery.Request)
		if err != nil {
			k.Logger(ctx).Error("unable to unmarshal request info for zone", "zone", zone.ZoneID(), "err", err)
			return err
		}

		if validatorsReq.Pagination == nil {
			validatorsReq.Pagination = new(query.PageRequest)
		}
		validatorsReq.Pagination.Key = validatorsRes.Pagination.NextKey
		k.Logger(ctx).Debug("Found pagination nextKey in valset; resubmitting...")
		err = k.EmitValSetQuery(ctx, icqQuery.ConnectionId, &zone, validatorsReq, sdkmath.NewInt(-1))
		if err != nil {
			return nil
		}
	}

	for _, validator := range validatorsRes.Validators {
		addr, err := addressutils.ValAddressFromBech32(validator.OperatorAddress, "")
		if err != nil {
			return err
		}
		val, found := k.GetValidator(ctx, &zone, addr)
		toQuery := false
		switch {
		case !found:
			k.Logger(ctx).Debug("Unable to find validator - fetching proof...", "valoper", validator.OperatorAddress)
			toQuery = true
		case !val.CommissionRate.Equal(validator.GetCommission()):
			k.Logger(ctx).Debug("Validator commission change; fetching proof", "valoper", validator.OperatorAddress, "from", val.CommissionRate, "to", validator.GetCommission())
			toQuery = true
		case !val.VotingPower.Equal(validator.Tokens):
			k.Logger(ctx).Debug("Validator voting power change; fetching proof", "valoper", validator.OperatorAddress, "from", val.VotingPower, "to", validator.Tokens)
			toQuery = true
		case !val.DelegatorShares.Equal(validator.DelegatorShares):
			k.Logger(ctx).Debug("Validator shares amount change; fetching proof", "valoper", validator.OperatorAddress, "from", val.DelegatorShares, "to", validator.DelegatorShares)
			toQuery = true
		}

		if toQuery {
			k.EmitValidatorQuery(ctx, icqQuery.ConnectionId, icqQuery.ChainId, validator)
		}
	}

	return nil
}

func (k *Keeper) SetValidatorForZone(ctx sdk.Context, zone *types.Zone, data []byte) error {
	if data == nil {
		k.Logger(ctx).Error("expected validator state, got nil")
		// return nil here, as if we receive nil we fail to unmarshal (as nil validators are invalid),
		// so we can never hope to resolve this query. Possibly received a valset update from a
		// different chain.
		return nil
	}
	validator, err := k.UnmarshalValidator(data)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validator info for zone", "zone", zone.BaseChainID(), "err", err)
		return err
	}

	valAddrBytes, err := addressutils.ValAddressFromBech32(validator.OperatorAddress, zone.GetValoperPrefix())
	if err != nil {
		return err
	}
	val, found := k.GetValidator(ctx, zone, valAddrBytes)
	if !found {
		k.Logger(ctx).Info("Unable to find validator - adding...", "valoper", validator.OperatorAddress)

		jailTime := time.Time{}
		if validator.IsJailed() {
			jailTime = ctx.BlockTime()
		}
		if err := k.SetValidator(ctx, zone, types.Validator{
			ValoperAddress:  validator.OperatorAddress,
			CommissionRate:  validator.GetCommission(),
			VotingPower:     validator.Tokens,
			DelegatorShares: validator.DelegatorShares,
			Score:           sdk.ZeroDec(),
			Status:          validator.Status.String(),
			Jailed:          validator.IsJailed(),
			JailedSince:     jailTime,
		}); err != nil {
			return err
		}

		if err := k.MakePerformanceDelegation(ctx, zone, validator.OperatorAddress); err != nil {
			return err
		}

	} else {
		if !val.Jailed && validator.IsJailed() {
			k.Logger(ctx).Info("Transitioning validator to jailed state", "valoper", validator.OperatorAddress, "old_vp", val.VotingPower, "new_vp", validator.Tokens, "new_shares", validator.DelegatorShares, "old_shares", val.DelegatorShares)

			val.Jailed = true
			val.JailedSince = ctx.BlockTime()
			if !val.VotingPower.IsPositive() {
				return fmt.Errorf("existing voting power must be greater than zero, received %s", val.VotingPower)
			}
			if !validator.Tokens.IsPositive() {
				return fmt.Errorf("incoming voting power must be greater than zero, received %s", validator.Tokens)
			}
			// determine difference between previous vp/shares ratio and new ratio.
			prevRatio := val.DelegatorShares.Quo(sdk.NewDecFromInt(val.VotingPower))
			newRatio := validator.DelegatorShares.Quo(sdk.NewDecFromInt(validator.Tokens))
			delta := newRatio.Quo(prevRatio)
			err = k.UpdateWithdrawalRecordsForSlash(ctx, zone, val.ValoperAddress, delta)
			if err != nil {
				return err
			}
		} else if val.Jailed && !validator.IsJailed() {
			k.Logger(ctx).Info("Transitioning validator to unjailed state", "valoper", validator.OperatorAddress)

			val.Jailed = false
			val.JailedSince = time.Time{}
		}

		if !val.CommissionRate.Equal(validator.GetCommission()) {
			k.Logger(ctx).Debug("Validator commission rate change; updating...", "valoper", validator.OperatorAddress, "oldRate", val.CommissionRate, "newRate", validator.GetCommission())
			val.CommissionRate = validator.GetCommission()
		}

		if !val.VotingPower.Equal(validator.Tokens) {
			k.Logger(ctx).Debug("Validator voting power change; updating", "valoper", validator.OperatorAddress, "oldPower", val.VotingPower, "newPower", validator.Tokens)
			val.VotingPower = validator.Tokens
		}

		if !val.DelegatorShares.Equal(validator.DelegatorShares) {
			k.Logger(ctx).Debug("Validator delegator shares change; updating", "valoper", validator.OperatorAddress, "oldShares", val.DelegatorShares, "newShares", validator.DelegatorShares)
			val.DelegatorShares = validator.DelegatorShares
		}

		if val.Status != validator.Status.String() {
			k.Logger(ctx).Debug("Transitioning validator status", "valoper", validator.OperatorAddress, "previous", val.Status, "current", validator.Status.String())

			val.Status = validator.Status.String()
		}

		if err := k.SetValidator(ctx, zone, val); err != nil {
			return err
		}

		if _, found := k.GetPerformanceDelegation(ctx, zone, validator.OperatorAddress); !found {
			if err := k.MakePerformanceDelegation(ctx, zone, validator.OperatorAddress); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k *Keeper) UpdateWithdrawalRecordsForSlash(ctx sdk.Context, zone *types.Zone, valoper string, delta sdk.Dec) error {
	var err error
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ZoneID(), types.WithdrawStatusUnbond, func(_ int64, record types.WithdrawalRecord) bool {
		recordSubAmount := sdkmath.ZeroInt()
		distr := record.Distribution
		for _, d := range distr {
			if d.Valoper != valoper {
				continue
			}
			newAmount := sdk.NewDec(int64(d.Amount)).Quo(delta).TruncateInt()
			thisSubAmount := sdkmath.NewInt(int64(d.Amount)).Sub(newAmount)
			recordSubAmount = recordSubAmount.Add(thisSubAmount)
			d.Amount = newAmount.Uint64()
			k.Logger(ctx).Info("Updated withdrawal record due to slashing", "valoper", valoper, "old_amount", d.Amount, "new_amount", newAmount.Int64(), "sub_amount", thisSubAmount.Int64())
		}
		record.Distribution = distr
		record.Amount = record.Amount.Sub(sdk.NewCoin(zone.BaseDenom, recordSubAmount))
		k.SetWithdrawalRecord(ctx, record)
		return false
	})
	return err
}

func (k *Keeper) depositInterval(ctx sdk.Context) zoneItrFn {
	return func(index int64, zone *types.Zone) (stop bool) {
		if zone.DepositAddress != nil {
			if !zone.DepositAddress.Balance.Empty() {
				k.Logger(ctx).Debug("balance is non zero", "balance", zone.DepositAddress.Balance)
				k.EmitDepositIntervalQuery(ctx, zone)

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

// MigrateParams fetches params, adds ClaimsEnabled field and re-sets params.
func (k *Keeper) MigrateParams(ctx sdk.Context) {
	params := types.Params{}
	params.DepositInterval = k.GetParam(ctx, types.KeyDepositInterval)
	params.CommissionRate = k.GetCommissionRate(ctx)
	params.ValidatorsetInterval = k.GetParam(ctx, types.KeyValidatorSetInterval)
	params.UnbondingEnabled = false

	k.paramStore.SetParamSet(ctx, &params)
}

func (k *Keeper) GetParams(clientCtx sdk.Context) (params types.Params) {
	k.paramStore.GetParamSet(clientCtx, &params)
	return params
}

// SetParams sets the distribution parameters to the param space.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}

func (k *Keeper) SetZoneIDForPortConnection(ctx sdk.Context, portID, connectionID, zoneID string) {
	key := fmt.Sprintf("%s-%s", portID, connectionID)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortConnectionZone)
	bz := []byte(zoneID)
	store.Set([]byte(key), bz)
}

func (k *Keeper) GetZoneIDFromPortConnection(ctx sdk.Context, portID, connectionID string) (string, error) {
	key := fmt.Sprintf("%s-%s", portID, connectionID)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortConnectionZone)
	bz := store.Get([]byte(key))
	if len(bz) == 0 {
		return "", fmt.Errorf("unable to find zone for port connection %s-%s", portID, connectionID)
	}

	return string(bz), nil
}

func (k *Keeper) GetChainIDFromConnection(ctx sdk.Context, connectionID string) (string, error) {
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

func (k *Keeper) GetChainIDFromContext(ctx sdk.Context) (string, error) {
	connectionID := ctx.Context().Value(utils.ContextKey("connectionID"))
	if connectionID == nil {
		return "", errors.New("connectionID not in context")
	}

	return k.GetChainIDFromConnection(ctx, connectionID.(string))
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
		zone.BaseChainID(),
		types.BankStoreKey,
		append(data, []byte(zone.BaseDenom)...),
		sdk.NewInt(-1),
		types.ModuleName,
		"perfbalance",
		100,
	)

	return nil
}

func (k *Keeper) EmitValSetQuery(ctx sdk.Context, connectionID string, zone *types.Zone, validatorsReq stakingtypes.QueryValidatorsRequest, period sdkmath.Int) error {
	bz, err := k.cdc.Marshal(&validatorsReq)
	if err != nil {
		return errors.New("failed to marshal valset pagination request")
	}

	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		zone.BaseChainID(),
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		period,
		types.ModuleName,
		"valset",
		0,
	)

	return nil
}

func (k *Keeper) EmitValidatorQuery(ctx sdk.Context, connectionID, zoneID string, validator stakingtypes.Validator) {
	_, addr, _ := bech32.DecodeAndConvert(validator.OperatorAddress)
	data := stakingtypes.GetValidatorKey(addr)
	k.ICQKeeper.MakeRequest(
		ctx,
		connectionID,
		zoneID,
		"store/staking/key",
		data,
		sdk.NewInt(-1),
		types.ModuleName,
		"validator",
		0,
	)
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
		zone.ZoneID(),
		"cosmos.tx.v1beta1.Service/GetTxsEvent",
		k.cdc.MustMarshal(&req),
		sdk.NewInt(-1),
		types.ModuleName,
		"depositinterval",
		0,
	)
}

func (k *Keeper) GetDelegationsInProcess(ctx sdk.Context, zone *types.Zone) sdkmath.Int {
	delegationsInProcess := sdkmath.ZeroInt()
	k.IterateZoneReceipts(ctx, zone, func(_ int64, receipt types.Receipt) (stop bool) {
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
	delegationsInProcess := k.GetDelegationsInProcess(ctx, zone)
	ratio, isZero := k.GetRatio(ctx, zone, epochRewards.Add(delegationsInProcess))
	k.Logger(ctx).Info("Epochly rewards", "coins", epochRewards)
	k.Logger(ctx).Info("Last redemption rate", "rate", zone.LastRedemptionRate)
	k.Logger(ctx).Info("Current redemption rate", "rate", zone.RedemptionRate)
	k.Logger(ctx).Info("New redemption rate", "rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, zone).Amount.Add(epochRewards).Add(delegationsInProcess))

	// soft cap redemption rate, instead of panicking.
	delta := ratio.Quo(zone.RedemptionRate)
	if delta.GT(sdk.NewDecWithPrec(102, 2)) {
		k.Logger(ctx).Error("ratio diverged by more than 2% upwards in the last epoch; capping at 1.02...")
		ratio = zone.RedemptionRate.Mul(sdk.NewDecWithPrec(102, 2))
	} else if delta.LT(sdk.NewDecWithPrec(95, 2)) && !isZero { // we allow a bigger downshift if all assets were withdrawn and we revert to zero.
		k.Logger(ctx).Error("ratio diverged by more than 5% downwards in the last epoch; 5% is the theoretical max if _all_ controlled tokens were tombstoned. capping at 0.95...")
		ratio = zone.RedemptionRate.Mul(sdk.NewDecWithPrec(95, 2))
	}

	zone.LastRedemptionRate = zone.RedemptionRate
	zone.RedemptionRate = ratio
	k.SetZone(ctx, zone)
}

func (k *Keeper) OverrideRedemptionRateNoCap(ctx sdk.Context, zone *types.Zone) {
	ratio, _ := k.GetRatio(ctx, zone, sdk.ZeroInt())
	k.Logger(ctx).Info("Last redemption rate", "rate", zone.LastRedemptionRate)
	k.Logger(ctx).Info("Current redemption rate", "rate", zone.RedemptionRate)
	k.Logger(ctx).Info("New redemption rate", "rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, zone).Amount)

	zone.RedemptionRate = ratio
	k.SetZone(ctx, zone)
}

func (k *Keeper) GetRatio(ctx sdk.Context, zone *types.Zone, epochRewards sdkmath.Int) (sdk.Dec, bool) {
	// native asset amount
	nativeAssetAmount := k.GetDelegatedAmount(ctx, zone).Amount
	nativeAssetUnbondingAmount := k.GetUnbondingAmount(ctx, zone).Amount

	// qAsset amount
	qAssetAmount := k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount

	// check if zone is fully withdrawn (no qAssets remain)
	if qAssetAmount.IsZero() {
		// ratio 1.0 (default 1:1 ratio between nativeAssets and qAssets)
		// native assets should not reach zero before qAssets (discount rate asymptote)
		return sdk.OneDec(), true
	}

	return sdk.NewDecFromInt(nativeAssetAmount.Add(epochRewards).Add(nativeAssetUnbondingAmount)).Quo(sdk.NewDecFromInt(qAssetAmount)), false
}

func (k *Keeper) GetAggregateIntentOrDefault(ctx sdk.Context, zone *types.Zone) (types.ValidatorIntents, error) {
	var intents types.ValidatorIntents
	var filteredIntents types.ValidatorIntents

	if len(zone.AggregateIntent) == 0 {
		intents = k.DefaultAggregateIntents(ctx, zone)
	} else {
		intents = zone.AggregateIntent
	}
	// filter intents here...
	// check validators for tombstoned
	for _, v := range intents {
		valAddrBytes, err := addressutils.ValAddressFromBech32(v.ValoperAddress, zone.GetValoperPrefix())
		if err != nil {
			return nil, err
		}
		val, found := k.GetValidator(ctx, zone, valAddrBytes)

		// this case should not happen as we check the validity of a validator entry when intent is set.
		if !found {
			continue
		}
		// we should never let tombstoned validators into the list, even if they are explicitly selected
		if val.Tombstoned {
			continue
		}

		// we should never let denylist validators into the list, even if they are explicitly selected
		// if in deny list {
		// continue
		// }
		filteredIntents = append(filteredIntents, v)
	}

	return filteredIntents, nil
}

func (k *Keeper) Rebalance(ctx sdk.Context, zone *types.Zone, epochNumber int64) error {
	currentAllocations, currentSum, currentLocked, lockedSum := k.GetDelegationMap(ctx, zone)
	targetAllocations, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return err
	}
	rebalances := types.DetermineAllocationsForRebalancing(currentAllocations, currentLocked, currentSum, lockedSum, targetAllocations, k.Logger(ctx))
	msgs := make([]proto.Message, 0)
	for _, rebalance := range rebalances {
		msgs = append(msgs, &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: rebalance.Source, ValidatorDstAddress: rebalance.Target, Amount: sdk.NewCoin(zone.BaseDenom, rebalance.Amount)})
		k.SetRedelegationRecord(ctx, types.RedelegationRecord{
			ChainId:     zone.ZoneID(),
			EpochNumber: epochNumber,
			Source:      rebalance.Source,
			Destination: rebalance.Target,
			Amount:      rebalance.Amount.Int64(),
		})
	}
	if len(msgs) == 0 {
		k.Logger(ctx).Info("No rebalancing required")
		return nil
	}
	k.Logger(ctx).Info("Send rebalancing messages", "msgs", msgs)
	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, types.EpochRebalanceMemo(epochNumber), zone.MessagesPerTx)
}

// UnmarshalValidatorsResponse attempts to umarshal  a byte slice into a QueryValidatorsResponse.
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

// UnmarshalValidatorsRequest attempts to umarshal a byte slice into a QueryValidatorsRequest.
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
