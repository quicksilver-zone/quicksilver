package keeper

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v5/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	epochskeeper "github.com/quicksilver-zone/quicksilver/x/epochs/keeper"
	interchainquerykeeper "github.com/quicksilver-zone/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/x/lsmtypes"
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
	EpochsKeeper        types.EpochsKeeper
	Ir                  codectypes.InterfaceRegistry
	hooks               types.IcsHooks
	paramStore          paramtypes.Subspace
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

func (k *Keeper) GetGovAuthority(_ sdk.Context) string {
	return sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), k.AccountKeeper.GetModuleAddress(govtypes.ModuleName))
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
	validatorsRes, err := k.UnmarshalValidatorsResponse(data)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validators info for zone", "zone", icqQuery.ChainId, "err", err)
		return err
	}

	if validatorsRes.Pagination != nil && !bytes.Equal(validatorsRes.Pagination.NextKey, []byte{}) {
		validatorsReq, err := k.UnmarshalValidatorsRequest(icqQuery.Request)
		if err != nil {
			k.Logger(ctx).Error("unable to unmarshal request info for zone", "zone", icqQuery.ChainId, "err", err)
			return err
		}

		if validatorsReq.Pagination == nil {
			k.Logger(ctx).Debug("unmarshalled a QueryValidatorsRequest with a nil Pagination", "zone", icqQuery.ChainId)
			validatorsReq.Pagination = new(query.PageRequest)
		}
		validatorsReq.Pagination.Key = validatorsRes.Pagination.NextKey
		k.Logger(ctx).Debug("Found pagination nextKey in valset; resubmitting...")
		err = k.EmitValSetQuery(ctx, icqQuery.ConnectionId, icqQuery.ChainId, validatorsReq, sdkmath.NewInt(-1))
		if err != nil {
			return nil
		}
	}

	for _, validator := range validatorsRes.Validators {
		addr, err := addressutils.ValAddressFromBech32(validator.OperatorAddress, "")
		if err != nil {
			return err
		}
		val, found := k.GetValidator(ctx, icqQuery.ChainId, addr)
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
		case val.Jailed != validator.Jailed:
			k.Logger(ctx).Debug("jail status change; fetching proof", "valoper", validator.OperatorAddress, "from", val.Jailed, "to", validator.Jailed)
			toQuery = true
		case val.Status != validator.Status.String():
			k.Logger(ctx).Debug("bond status change; fetching proof", "valoper", validator.OperatorAddress, "from", val.Status, "to", validator.Status.String())
			toQuery = true
		case !validator.LiquidShares.IsNil() && !val.LiquidShares.Equal(validator.LiquidShares):
			k.Logger(ctx).Debug("liquid shares amount change; fetching proof", "valoper", validator.OperatorAddress, "from", val.LiquidShares, "to", validator.LiquidShares)
			toQuery = true
		case !validator.ValidatorBondShares.IsNil() && !val.ValidatorBondShares.Equal(validator.ValidatorBondShares):
			k.Logger(ctx).Debug("Validator bond shares amount change; fetching proof", "valoper", validator.OperatorAddress, "from", val.ValidatorBondShares, "to", validator.ValidatorBondShares)
			toQuery = true
		}

		if toQuery {
			if err := k.EmitValidatorQuery(ctx, icqQuery.ConnectionId, icqQuery.ChainId, validator); err != nil {
				k.Logger(ctx).Error("EmitValidatorQuery error", "valoper", validator.OperatorAddress, "err", err)
				return err
			}
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
		k.Logger(ctx).Error("unable to unmarshal validator info for zone", "zone", zone.ChainId, "err", err)
		return err
	}

	valAddrBytes, err := addressutils.ValAddressFromBech32(validator.OperatorAddress, zone.GetValoperPrefix())
	if err != nil {
		return err
	}
	val, found := k.GetValidator(ctx, zone.ChainId, valAddrBytes)
	if !found {
		k.Logger(ctx).Debug("Unable to find validator - adding...", "valoper", validator.OperatorAddress)

		jailTime := time.Time{}
		if validator.IsJailed() {
			var pk cryptotypes.PubKey
			err := k.cdc.UnpackAny(validator.ConsensusPubkey, &pk)
			if err != nil {
				return err
			}
			consAddr := sdk.ConsAddress(pk.Address().Bytes())
			k.SetValidatorAddrByConsAddr(ctx, zone.ChainId, validator.OperatorAddress, consAddr)
			jailTime = ctx.BlockTime()

			err = k.EmitSigningInfoQuery(ctx, zone.ConnectionId, zone.ChainId, validator)
			if err != nil {
				return err
			}
		}

		if err := k.SetValidator(ctx, zone.ChainId, types.Validator{
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
		if val.Tombstoned {
			k.Logger(ctx).Debug(fmt.Sprintf("%q on chainID: %q was found to already have been tombstoned; not updating state.", validator.OperatorAddress, zone.ChainId))
			return nil
		}

		if !val.Jailed && validator.IsJailed() {
			k.Logger(ctx).Info("Transitioning validator to jailed state", "valoper", validator.OperatorAddress, "old_vp", val.VotingPower, "new_vp", validator.Tokens, "new_shares", validator.DelegatorShares, "old_shares", val.DelegatorShares)

			var pk cryptotypes.PubKey
			err := k.cdc.UnpackAny(validator.ConsensusPubkey, &pk)
			if err != nil {
				return err
			}
			consAddr := sdk.ConsAddress(pk.Address().Bytes())
			k.SetValidatorAddrByConsAddr(ctx, zone.ChainId, validator.OperatorAddress, consAddr)

			err = k.EmitSigningInfoQuery(ctx, zone.ConnectionId, zone.ChainId, validator)
			if err != nil {
				return err
			}

			val.Jailed = true
			val.JailedSince = ctx.BlockTime()

			// be defensive, so we don't get divison weirdness!
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
			k.Logger(ctx).Debug("Transitioning validator to unjailed state", "valoper", validator.OperatorAddress)

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

		if !validator.ValidatorBondShares.IsNil() && !val.ValidatorBondShares.Equal(validator.ValidatorBondShares) {
			k.Logger(ctx).Info("Validator bonded shares change; updating", "valoper", validator.OperatorAddress, "oldShares", val.ValidatorBondShares, "newShares", validator.ValidatorBondShares)
			val.ValidatorBondShares = validator.ValidatorBondShares
		}

		if !validator.LiquidShares.IsNil() && !val.LiquidShares.Equal(validator.LiquidShares) {
			k.Logger(ctx).Info("Validator liquid shares change; updating", "valoper", validator.OperatorAddress, "oldShares", val.LiquidShares, "newShares", validator.LiquidShares)
			val.LiquidShares = validator.LiquidShares
		}

		if err := k.SetValidator(ctx, zone.ChainId, val); err != nil {
			return err
		}

		if _, found := k.GetPerformanceDelegation(ctx, zone.ChainId, zone.PerformanceAddress, validator.OperatorAddress); !found {
			if err := k.MakePerformanceDelegation(ctx, zone, validator.OperatorAddress); err != nil {
				return err
			}
		}
	}

	return nil
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

func (k *Keeper) GetChainIDFromContext(ctx sdk.Context) (string, error) {
	connectionID := ctx.Context().Value(utils.ContextKey("connectionID"))
	if connectionID == nil {
		return "", errors.New("connectionID not in context")
	}

	return k.GetChainID(ctx, connectionID.(string))
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
		append(data, []byte(zone.BaseDenom)...),
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

func (k *Keeper) EmitValidatorQuery(ctx sdk.Context, connectionID, chainID string, validator lsmstakingtypes.Validator) error {
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

func (k *Keeper) EmitSigningInfoQuery(ctx sdk.Context, connectionID, chainID string, validator lsmstakingtypes.Validator) error {
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
	delegationsInProcess := k.GetDelegationsInProcess(ctx, zone.ChainId)
	ratio, isZero := k.GetRatio(ctx, zone, epochRewards.Add(delegationsInProcess))
	k.Logger(ctx).Info("Redemption Rate Update", "chain", zone.ChainId, "epochly_rewards", epochRewards, "last_rate", zone.LastRedemptionRate, "current_rate", zone.RedemptionRate, "new_rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, zone).Amount.Add(epochRewards).Add(delegationsInProcess))

	// TODO: make max deltas params.
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
	delegationsInProcess := k.GetDelegationsInProcess(ctx, zone.ChainId)
	ratio, _ := k.GetRatio(ctx, zone, delegationsInProcess)
	k.Logger(ctx).Info("Forced Redemption Rate Update", "chain", zone.ChainId, "last_rate", zone.LastRedemptionRate, "current_rate", zone.RedemptionRate, "new_rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, zone).Amount.Add(delegationsInProcess))

	zone.LastRedemptionRate = zone.RedemptionRate
	zone.RedemptionRate = ratio
	k.SetZone(ctx, zone)
}

func (k *Keeper) GetRatio(ctx sdk.Context, zone *types.Zone, epochRewards sdkmath.Int) (sdk.Dec, bool) {
	// native asset amount
	nativeAssetAmount := k.GetDelegatedAmount(ctx, zone).Amount
	nativeAssetUnbondingAmount := k.GetUnbondingAmount(ctx, zone).Amount
	nativeAssetUnbonded := zone.DelegationAddress.Balance.AmountOf(zone.BaseDenom)

	// qAsset amount
	qAssetAmount := k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount

	// check if zone is fully withdrawn (no qAssets remain)
	if qAssetAmount.IsZero() {
		// ratio 1.0 (default 1:1 ratio between nativeAssets and qAssets)
		// native assets should not reach zero before qAssets (discount rate asymptote)
		return sdk.OneDec(), true
	}

	return sdk.NewDecFromInt(nativeAssetAmount.Add(epochRewards).Add(nativeAssetUnbondingAmount).Add(nativeAssetUnbonded)).Quo(sdk.NewDecFromInt(qAssetAmount)), false
}

func (k *Keeper) GetAggregateIntentOrDefault(ctx sdk.Context, zone *types.Zone) (types.ValidatorIntents, error) {
	var intents types.ValidatorIntents
	var filteredIntents types.ValidatorIntents

	if len(zone.AggregateIntent) == 0 {
		intents = k.DefaultAggregateIntents(ctx, zone.ChainId)
	} else {
		intents = zone.AggregateIntent
	}

	jailedThreshold := k.EpochsKeeper.GetEpochInfo(ctx, "epoch").Duration * 2

	// filter intents here...
	// check validators for tombstoned
	for _, validatorIntent := range intents {
		valAddrBytes, err := addressutils.ValAddressFromBech32(validatorIntent.ValoperAddress, zone.GetValoperPrefix())
		if err != nil {
			return nil, err
		}
		validator, found := k.GetValidator(ctx, zone.ChainId, valAddrBytes)

		// this case should not happen as we check the validity of a validator entry when intent is set.
		if !found {
			continue
		}
		// we should never let tombstoned validators into the list, even if they are explicitly selected
		if validator.Tombstoned {
			continue
		}

		// if the validator has been jailed for > two epochs, remove them.
		if validator.Jailed && validator.JailedSince.Add(jailedThreshold).Before(ctx.BlockTime()) {
			continue
		}

		// we should never let denylist validators into the list, even if they are explicitly selected
		// if in deny list {
		// continue
		// }
		filteredIntents = append(filteredIntents, validatorIntent)
	}

	return filteredIntents, nil
}

func (k *Keeper) Rebalance(ctx sdk.Context, zone *types.Zone, epochNumber int64) error {
	currentAllocations, currentSum, currentLocked, lockedSum := k.GetDelegationMap(ctx, zone.ChainId)
	targetAllocations, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return err
	}
	maxCanAllocate := k.DetermineMaximumValidatorAllocations(ctx, zone)
	rebalances := types.DetermineAllocationsForRebalancing(currentAllocations, currentLocked, currentSum, lockedSum, targetAllocations, maxCanAllocate, k.Logger(ctx)).RemoveDuplicates()
	msgs := make([]sdk.Msg, 0)
	for _, rebalance := range rebalances {
		msgs = append(msgs, &stakingtypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: rebalance.Source, ValidatorDstAddress: rebalance.Target, Amount: sdk.NewCoin(zone.BaseDenom, rebalance.Amount)})
		k.SetRedelegationRecord(ctx, types.RedelegationRecord{
			ChainId:     zone.ChainId,
			EpochNumber: epochNumber,
			Source:      rebalance.Source,
			Destination: rebalance.Target,
			Amount:      rebalance.Amount.Int64(),
		})
	}
	if len(msgs) == 0 {
		k.Logger(ctx).Debug("No rebalancing required")
		return nil
	}
	k.Logger(ctx).Info("Send rebalancing messages", "msgs", msgs)
	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, types.EpochRebalanceMemo(epochNumber), zone.MessagesPerTx)
}

// UnmarshalValidatorsResponse attempts to umarshal  a byte slice into a QueryValidatorsResponse.
func (k *Keeper) UnmarshalValidatorsResponse(data []byte) (lsmstakingtypes.QueryValidatorsResponse, error) {
	validatorsRes := lsmstakingtypes.QueryValidatorsResponse{}
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
func (k *Keeper) UnmarshalValidator(data []byte) (lsmstakingtypes.Validator, error) {
	validator := lsmstakingtypes.Validator{}
	if len(data) == 0 {
		return validator, errors.New("attempted to unmarshal zero length byte slice (9)")
	}
	err := k.cdc.Unmarshal(data, &validator)
	if err != nil {
		return validator, err
	}

	return validator, nil
}
