package keeper

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v5/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"

	config "github.com/ingenuity-build/quicksilver/cmd/config"
	"github.com/ingenuity-build/quicksilver/utils"
	claimsmanagerkeeper "github.com/ingenuity-build/quicksilver/x/claimsmanager/keeper"
	interchainquerykeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc                 codec.Codec
	storeKey            storetypes.StoreKey
	scopedKeeper        *capabilitykeeper.ScopedKeeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICQKeeper           interchainquerykeeper.Keeper
	AccountKeeper       authKeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	IBCKeeper           ibckeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	ClaimsManagerKeeper claimsmanagerkeeper.Keeper
	paramStore          paramtypes.Subspace
}

// NewKeeper returns a new instance of zones Keeper.
// This function will panic on failure.
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, accountKeeper authKeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, icacontrollerkeeper icacontrollerkeeper.Keeper, scopedKeeper *capabilitykeeper.ScopedKeeper, icqKeeper interchainquerykeeper.Keeper, ibcKeeper ibckeeper.Keeper, transferKeeper ibctransferkeeper.Keeper, claimsManagerKeeper claimsmanagerkeeper.Keeper, ps paramtypes.Subspace) Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if addr := accountKeeper.GetModuleAddress(types.EscrowModuleAccount); addr == nil {
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
		TransferKeeper:      transferKeeper,
		ClaimsManagerKeeper: claimsManagerKeeper,

		paramStore: ps,
	}
}

func (k *Keeper) GetGovAuthority(_ sdk.Context) string {
	return sdk.MustBech32ifyAddressBytes(config.Bech32Prefix, k.AccountKeeper.GetModuleAddress(govtypes.ModuleName))
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetCodec() codec.Codec {
	return k.cdc
}

func (k *Keeper) ScopedKeeper() *capabilitykeeper.ScopedKeeper {
	return k.scopedKeeper
}

// ClaimCapability claims the channel capability passed via the OnOpenChanInit callback
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (k *Keeper) SetConnectionForPort(ctx sdk.Context, connectionID string, port string) {
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
func (k Keeper) IteratePortConnections(ctx sdk.Context, cb func(pc types.PortConnectionTuple) (stop bool)) {
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
func (k Keeper) AllPortConnections(ctx sdk.Context) (pcs []types.PortConnectionTuple) {
	k.IteratePortConnections(ctx, func(pc types.PortConnectionTuple) bool {
		pcs = append(pcs, pc)
		return false
	})

	return pcs
}

// ### Interval functions >>>
// * some of these functions (or portions thereof) may be changed to single
//   query type functions, dependent upon callback features / capabilities;

func SetValidatorsForZone(k *Keeper, ctx sdk.Context, zoneInfo types.Zone, data []byte, request []byte) error {
	validatorsRes := stakingTypes.QueryValidatorsResponse{}
	if len(data) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (8)")
	}
	err := k.cdc.Unmarshal(data, &validatorsRes)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validators info for zone", "zone", zoneInfo.ChainId, "err", err)
		return err
	}

	if validatorsRes.Pagination != nil && !bytes.Equal(validatorsRes.Pagination.NextKey, []byte{}) {
		validatorsReq := stakingTypes.QueryValidatorsRequest{Pagination: &query.PageRequest{}}
		err = k.cdc.Unmarshal(request, &validatorsReq)
		if err != nil {
			k.Logger(ctx).Error("unable to unmarshal request", "zone", zoneInfo.ChainId, "err", err)
			return err
		}
		if validatorsReq.Pagination == nil {
			k.Logger(ctx).Info("unmarshalled a QueryValidatorsRequest with a nil Pagination", "zone", zoneInfo.ChainId)
			validatorsReq.Pagination = new(query.PageRequest)
		}
		validatorsReq.Pagination.Key = validatorsRes.Pagination.NextKey
		bz, err := k.cdc.Marshal(&validatorsReq)
		if err != nil {
			return errors.New("failed to marshal valset pagination request")
		}
		k.Logger(ctx).Info("Found pagination nextKey in valset; resubmitting...")

		k.ICQKeeper.MakeRequest(
			ctx,
			zoneInfo.ConnectionId,
			zoneInfo.ChainId,
			"cosmos.staking.v1beta1.Query/Validators",
			bz,
			sdk.NewInt(-1),
			types.ModuleName,
			"valset",
			0,
		)
	}

	for _, validator := range validatorsRes.Validators {
		val, found := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
		toQuery := false
		switch {
		case !found:
			k.Logger(ctx).Info("Unable to find validator - fetching proof...", "valoper", validator.OperatorAddress)
			toQuery = true
		case !val.CommissionRate.Equal(validator.GetCommission()):
			k.Logger(ctx).Info("Validator commission change; fetching proof", "valoper", validator.OperatorAddress, "from", val.CommissionRate, "to", validator.GetCommission())
			toQuery = true
		case !val.VotingPower.Equal(validator.Tokens):
			k.Logger(ctx).Info("Validator voting power change; fetching proof", "valoper", validator.OperatorAddress, "from", val.VotingPower, "to", validator.Tokens)
			toQuery = true
		case !val.DelegatorShares.Equal(validator.DelegatorShares):
			k.Logger(ctx).Info("Validator shares amount change; fetching proof", "valoper", validator.OperatorAddress, "from", val.DelegatorShares, "to", validator.DelegatorShares)
			toQuery = true
		}

		if toQuery {
			_, addr, _ := bech32.DecodeAndConvert(validator.OperatorAddress)
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
	k.SetZone(ctx, &zoneInfo)
	return nil
}

func SetValidatorForZone(k *Keeper, ctx sdk.Context, zoneInfo types.Zone, data []byte) error {
	validator := stakingTypes.Validator{}
	if len(data) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (9)")
	}
	err := k.cdc.Unmarshal(data, &validator)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validator info for zone", "zone", zoneInfo.ChainId, "err", err)
		return err
	}

	val, found := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
	if !found {
		k.Logger(ctx).Info("Unable to find validator - adding...", "valoper", validator.OperatorAddress)

		jailTime := time.Time{}
		if validator.IsJailed() {
			jailTime = ctx.BlockTime()
		}
		zoneInfo.Validators = append(zoneInfo.Validators, &types.Validator{
			ValoperAddress:  validator.OperatorAddress,
			CommissionRate:  validator.GetCommission(),
			VotingPower:     validator.Tokens,
			DelegatorShares: validator.DelegatorShares,
			Score:           sdk.ZeroDec(),
			Status:          validator.Status.String(),
			Jailed:          validator.IsJailed(),
			JailedSince:     jailTime,
		})
		zoneInfo.Validators = zoneInfo.GetValidatorsSorted()

		if err := k.MakePerformanceDelegation(ctx, &zoneInfo, validator.OperatorAddress); err != nil {
			return err
		}

	} else {

		if !val.CommissionRate.Equal(validator.GetCommission()) {
			k.Logger(ctx).Info("Validator commission rate change; updating...", "valoper", validator.OperatorAddress, "oldRate", val.CommissionRate, "newRate", validator.GetCommission())
			val.CommissionRate = validator.GetCommission()
		}

		if !val.VotingPower.Equal(validator.Tokens) {
			k.Logger(ctx).Info("Validator voting power change; updating", "valoper", validator.OperatorAddress, "oldPower", val.VotingPower, "newPower", validator.Tokens)
			val.VotingPower = validator.Tokens
		}

		if !val.DelegatorShares.Equal(validator.DelegatorShares) {
			k.Logger(ctx).Info("Validator delegator shares change; updating", "valoper", validator.OperatorAddress, "oldShares", val.DelegatorShares, "newShares", validator.DelegatorShares)
			val.DelegatorShares = validator.DelegatorShares
		}

		if !val.Jailed && validator.IsJailed() {
			k.Logger(ctx).Info("Transitioning validator to jailed state", "valoper", validator.OperatorAddress)

			val.Jailed = true
			val.JailedSince = ctx.BlockTime()
		} else if val.Jailed && !validator.IsJailed() {
			k.Logger(ctx).Info("Transitioning validator to unjailed state", "valoper", validator.OperatorAddress)

			val.Jailed = false
			val.JailedSince = time.Time{}
		}

		if val.Status != validator.Status.String() {
			k.Logger(ctx).Info("Transitioning validator status", "valoper", validator.OperatorAddress, "previous", val.Status, "current", validator.Status.String())

			val.Status = validator.Status.String()
		}

		if _, found := k.GetPerformanceDelegation(ctx, &zoneInfo, validator.OperatorAddress); !found {
			if err := k.MakePerformanceDelegation(ctx, &zoneInfo, validator.OperatorAddress); err != nil {
				return err
			}
		}
	}

	k.SetZone(ctx, &zoneInfo)
	return nil
}

func (k Keeper) depositInterval(ctx sdk.Context) zoneItrFn {
	return func(index int64, zoneInfo types.Zone) (stop bool) {
		if zoneInfo.DepositAddress != nil {
			if !zoneInfo.DepositAddress.Balance.Empty() {
				k.Logger(ctx).Info("balance is non zero", "balance", zoneInfo.DepositAddress.Balance)

				req := tx.GetTxsEventRequest{Events: []string{"transfer.recipient='" + zoneInfo.DepositAddress.GetAddress() + "'"}, OrderBy: tx.OrderBy_ORDER_BY_DESC, Pagination: &query.PageRequest{Limit: types.TxRetrieveCount}}
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
func (k Keeper) MigrateParams(ctx sdk.Context) {
	params := types.Params{}
	params.DepositInterval = k.GetParam(ctx, types.KeyDepositInterval)
	params.CommissionRate = k.GetCommissionRate(ctx)
	params.ValidatorsetInterval = k.GetParam(ctx, types.KeyValidatorSetInterval)
	params.UnbondingEnabled = false

	k.paramStore.SetParamSet(ctx, &params)
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

func (k Keeper) GetChainIDFromContext(ctx sdk.Context) (string, error) {
	connectionID := ctx.Context().Value(utils.ContextKey("connectionID"))
	if connectionID == nil {
		return "", errors.New("connectionID not in context")
	}

	return k.GetChainID(ctx, connectionID.(string))
}

func (k Keeper) EmitPerformanceBalanceQuery(ctx sdk.Context, zone *types.Zone) error {
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

	return nil
}

// redemption rate

func (k *Keeper) UpdateRedemptionRate(ctx sdk.Context, zone types.Zone, epochRewards math.Int) {
	delegationsInProcess := sdk.ZeroInt()
	k.IterateZoneReceipts(ctx, &zone, func(_ int64, receipt types.Receipt) (stop bool) {
		if receipt.Completed == nil {
			for _, coin := range receipt.Amount {
				delegationsInProcess = delegationsInProcess.Add(coin.Amount) // we cannot simply choose
			}
		}
		return false
	})
	ratio, isZero := k.GetRatio(ctx, zone, epochRewards.Add(delegationsInProcess))
	k.Logger(ctx).Info("Epochly rewards", "coins", epochRewards)
	k.Logger(ctx).Info("Last redemption rate", "rate", zone.LastRedemptionRate)
	k.Logger(ctx).Info("Current redemption rate", "rate", zone.RedemptionRate)
	k.Logger(ctx).Info("New redemption rate", "rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, &zone).Amount.Add(epochRewards).Add(delegationsInProcess))

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
	k.SetZone(ctx, &zone)
}

func (k *Keeper) OverrideRedemptionRateNoCap(ctx sdk.Context, zone types.Zone) {
	delegationsInProcess := sdk.ZeroInt()
	k.IterateZoneReceipts(ctx, &zone, func(_ int64, receipt types.Receipt) (stop bool) {
		if receipt.Completed == nil {
			for _, coin := range receipt.Amount {
				delegationsInProcess = delegationsInProcess.Add(coin.Amount) // we cannot simply choose
			}
		}
		return false
	})
	ratio, _ := k.GetRatio(ctx, zone, delegationsInProcess)
	k.Logger(ctx).Info("Last redemption rate", "rate", zone.LastRedemptionRate)
	k.Logger(ctx).Info("Current redemption rate", "rate", zone.RedemptionRate)
	k.Logger(ctx).Info("New redemption rate", "rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, &zone).Amount.Add(delegationsInProcess))

	zone.RedemptionRate = ratio
	k.SetZone(ctx, &zone)
}

func (k *Keeper) GetRatio(ctx sdk.Context, zone types.Zone, epochRewards math.Int) (sdk.Dec, bool) {
	// native asset amount
	nativeAssetAmount := k.GetDelegatedAmount(ctx, &zone).Amount
	nativeAssetUnbondingAmount := k.GetUnbondingAmount(ctx, &zone).Amount

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

func (k *Keeper) Rebalance(ctx sdk.Context, zone types.Zone, epochNumber int64) error {
	currentAllocations, currentSum := k.GetDelegationMap(ctx, &zone)
	targetAllocations := zone.GetAggregateIntentOrDefault()
	rebalances := DetermineAllocationsForRebalancing(currentAllocations, currentSum, targetAllocations, k.ZoneRedelegationRecords(ctx, zone.ChainId))
	msgs := make([]sdk.Msg, 0)
	for _, rebalance := range rebalances {
		msgs = append(msgs, &stakingTypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: rebalance.Source, ValidatorDstAddress: rebalance.Target, Amount: sdk.NewCoin(zone.BaseDenom, rebalance.Amount)})
		k.SetRedelegationRecord(ctx, types.RedelegationRecord{
			ChainId:     zone.ChainId,
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
	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, fmt.Sprintf("rebalance/%d", epochNumber), zone.MessagesPerTx)
}

func (k *Keeper) CloseICAChannel(ctx sdk.Context, portID, channelID string) error {
	_, capability, err := k.IBCKeeper.ChannelKeeper.LookupModuleByChannel(ctx, portID, channelID)
	if err != nil {
		return err
	}

	if err := k.IBCKeeper.ChannelKeeper.ChanCloseInit(ctx, portID, channelID, capability); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeCloseICA,
			sdk.NewAttribute(types.AttributeKeyPortID, portID),
			sdk.NewAttribute(types.AttributeKeyChannelID, channelID),
		),
	})

	return nil
}

type RebalanceTarget struct {
	Amount math.Int
	Source string
	Target string
}

func DetermineAllocationsForRebalancing(currentAllocations map[string]math.Int, currentSum math.Int, targetAllocations types.ValidatorIntents, existingRedelegations []types.RedelegationRecord) []RebalanceTarget {
	out := make([]RebalanceTarget, 0)
	deltas := CalculateDeltas(currentAllocations, currentSum, targetAllocations)

	wantToRebalance := sdk.ZeroInt()
	canRebalanceFrom := sdk.ZeroInt()

	totalLocked := int64(0)
	lockedPerValidator := map[string]int64{}
	for _, redelegation := range existingRedelegations {
		totalLocked += redelegation.Amount
		thisLocked, found := lockedPerValidator[redelegation.Destination]
		if !found {
			thisLocked = 0
		}
		lockedPerValidator[redelegation.Destination] = thisLocked + redelegation.Amount
	}

	fmt.Println("Total locked (per-validator)", totalLocked, lockedPerValidator)

	// TODO: make this a param
	maxCanRebalance := currentSum.Sub(math.NewInt(totalLocked)).Quo(sdk.NewInt(2))
	fmt.Println("Can rebalance with current locking", maxCanRebalance)

	// TODO: make this a param
	maxCanRebalance = math.MinInt(maxCanRebalance, currentSum.Quo(sdk.NewInt(7)))
	fmt.Println("Can rebalance with locking this epoch", maxCanRebalance)

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].ValoperAddress < deltas[j].ValoperAddress
	})

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].Weight.GT(deltas[j].Weight)
	})

	for _, delta := range deltas {
		switch {
		case delta.Weight.IsZero():
			// do nothing
		case delta.Weight.IsPositive():
			// if delta > current value - locked value, truncate, as we cannot rebalance locked tokens.
			wantToRebalance = wantToRebalance.Add(delta.Weight.TruncateInt())
		case delta.Weight.IsNegative():

			if delta.Weight.Abs().GT(sdk.NewDecFromInt(currentAllocations[delta.ValoperAddress].Sub(math.NewInt(lockedPerValidator[delta.ValoperAddress])))) {
				delta.Weight = sdk.NewDecFromInt(currentAllocations[delta.ValoperAddress].Sub(math.NewInt(lockedPerValidator[delta.ValoperAddress]))).Neg()
				fmt.Printf("Truncated delta for %s to %d due to locked tokens\n", delta.ValoperAddress, delta.Weight.Abs())
			}
			canRebalanceFrom = canRebalanceFrom.Add(delta.Weight.Abs().TruncateInt())
		}
	}

	toRebalance := sdk.MinInt(sdk.MinInt(wantToRebalance, canRebalanceFrom), maxCanRebalance)

	if toRebalance.Equal(math.ZeroInt()) {
		fmt.Println("No rebalancing this epoch")
		return []RebalanceTarget{}
	}
	fmt.Println("Will rebalance this epoch", toRebalance)

	tgtIdx := 0
	srcIdx := len(deltas) - 1
	for i := 0; toRebalance.GT(sdk.ZeroInt()); {
		i++
		if i > 20 {
			break
		}
		src := deltas[srcIdx]
		tgt := deltas[tgtIdx]
		if src.ValoperAddress == tgt.ValoperAddress {
			break
		}
		var amount math.Int
		if src.Weight.Abs().TruncateInt().IsZero() { //nolint:gocritic
			srcIdx--
			continue
		} else if src.Weight.Abs().TruncateInt().GT(toRebalance) { // amount == rebalance
			amount = toRebalance
		} else {
			amount = src.Weight.Abs().TruncateInt()
		}

		if tgt.Weight.Abs().TruncateInt().IsZero() {
			tgtIdx++
			continue
		} else if tgt.Weight.Abs().TruncateInt().LTE(toRebalance) {
			amount = sdk.MinInt(amount, tgt.Weight.Abs().TruncateInt())
		}
		out = append(out, RebalanceTarget{Amount: amount, Target: tgt.ValoperAddress, Source: src.ValoperAddress})
		deltas[srcIdx].Weight = src.Weight.Add(sdk.NewDecFromInt(amount))
		deltas[tgtIdx].Weight = tgt.Weight.Sub(sdk.NewDecFromInt(amount))
		toRebalance = toRebalance.Sub(amount)

	}

	// sort keys by relative value of delta
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Source < out[j].Source
	})

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Target < out[j].Target
	})

	// sort keys by relative value of delta
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Amount.GT(out[j].Amount)
	})

	return out
}
