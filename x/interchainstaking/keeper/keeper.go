package keeper

import (
	"bytes"
	"fmt"
	"sort"

	"cosmossdk.io/math"
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
	icacontrollerkeeper "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ingenuity-build/quicksilver/utils"
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
	paramStore          paramtypes.Subspace
}

// NewKeeper returns a new instance of zones Keeper.
// This function will panic on failure.
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, accountKeeper authKeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, icacontrollerkeeper icacontrollerkeeper.Keeper, scopedKeeper *capabilitykeeper.ScopedKeeper, icqKeeper interchainquerykeeper.Keeper, ibcKeeper ibckeeper.Keeper, ps paramtypes.Subspace) Keeper {
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

func SetValidatorsForZone(k *Keeper, ctx sdk.Context, zoneInfo types.Zone, data []byte) error {
	validatorsRes := stakingTypes.QueryValidatorsResponse{}
	if bytes.Equal(data, []byte("")) {
		return fmt.Errorf("attempted to unmarshal zero length byte slice (8)")
	}
	err := k.cdc.Unmarshal(data, &validatorsRes)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validators info for zone", "zone", zoneInfo.ChainId, "err", err)
		return err
	}

	for _, validator := range validatorsRes.Validators {
		_, addr, _ := bech32.DecodeAndConvert(validator.OperatorAddress)
		val, found := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
		if !found {
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
	if bytes.Equal(data, []byte("")) {
		return fmt.Errorf("attempted to unmarshal zero length byte slice (9)")
	}
	err := k.cdc.Unmarshal(data, &validator)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal validator info for zone", "zone", zoneInfo.ChainId, "err", err)
		return err
	}

	val, found := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
	if !found {
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

		if validator.GetCommission().IsNil() || !val.CommissionRate.Equal(validator.GetCommission()) {
			val.CommissionRate = validator.GetCommission()
			k.Logger(ctx).Info("Validator commission rate change; updating...", "valoper", validator.OperatorAddress, "oldRate", val.CommissionRate, "newRate", validator.GetCommission())
		}

		if validator.Tokens.IsNil() || !val.VotingPower.Equal(validator.Tokens) {
			val.VotingPower = validator.Tokens
			k.Logger(ctx).Info("Validator voting power change; updating", "valoper", validator.OperatorAddress, "oldPower", val.VotingPower, "newPower", validator.Tokens)
		}

		if validator.DelegatorShares.IsNil() || !val.DelegatorShares.Equal(validator.DelegatorShares) {
			val.DelegatorShares = validator.DelegatorShares
			k.Logger(ctx).Info("Validator delegator shares change; updating", "valoper", validator.OperatorAddress, "oldShares", val.DelegatorShares, "newShares", validator.DelegatorShares)
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

func (k Keeper) GetChainIDFromContext(ctx sdk.Context) (string, error) {
	connectionID := ctx.Context().Value(utils.ContextKey("connectionID"))
	if connectionID == nil {
		return "", fmt.Errorf("connectionID not in context")
	}

	return k.GetChainID(ctx, connectionID.(string))
}

func (k Keeper) EmitPerformanceBalanceQuery(ctx sdk.Context, zone *types.Zone) error {
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

// redemption rate

func (k *Keeper) assertRedemptionRateWithinBounds(ctx sdk.Context, previousRate sdk.Dec, newRate sdk.Dec) error {
	// TODO: what is an acceptable deviation?
	return nil
}

func (k *Keeper) updateRedemptionRate(ctx sdk.Context, zone types.Zone, epochRewards math.Int) {
	ratio := k.getRatio(ctx, zone, epochRewards)
	k.Logger(ctx).Info("Epochly rewards", "coins", epochRewards)
	k.Logger(ctx).Info("Last redemption rate", "rate", zone.LastRedemptionRate)
	k.Logger(ctx).Info("Current redemption rate", "rate", zone.RedemptionRate)
	k.Logger(ctx).Info("New redemption rate", "rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount, "lv", k.GetDelegatedAmount(ctx, &zone).Amount.Add(epochRewards))

	if err := k.assertRedemptionRateWithinBounds(ctx, zone.RedemptionRate, ratio); err != nil {
		panic("Redemption rate out of bounds")
	}
	zone.LastRedemptionRate = zone.RedemptionRate
	zone.RedemptionRate = ratio
	k.SetZone(ctx, &zone)
}

func (k *Keeper) getRatio(ctx sdk.Context, zone types.Zone, epochRewards math.Int) sdk.Dec {
	// native asset amount
	nativeAssetAmount := k.GetDelegatedAmount(ctx, &zone).Amount
	// qAsset amount
	qAssetAmount := k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount

	// check if zone is fully withdrawn (no qAssets remain)
	if qAssetAmount.IsZero() {
		// ratio 1.0 (default 1:1 ratio between nativeAssets and qAssets)
		// native assets should not reach zero before qAssets (discount rate asymptote)
		return sdk.OneDec()
	}

	return sdk.NewDecFromInt(nativeAssetAmount.Add(epochRewards)).Quo(sdk.NewDecFromInt(qAssetAmount))
}

func (k *Keeper) Rebalance(ctx sdk.Context, zone types.Zone) error {
	currentAllocations, currentSum := k.GetDelegationMap(ctx, &zone)
	targetAllocations := zone.GetAggregateIntentOrDefault()
	rebalances := DetermineAllocationsForRebalancing(currentAllocations, currentSum, targetAllocations)
	msgs := make([]sdk.Msg, 0)
	for _, rebalance := range rebalances {
		msgs = append(msgs, &stakingTypes.MsgBeginRedelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorSrcAddress: rebalance.Source, ValidatorDstAddress: rebalance.Target, Amount: sdk.NewCoin(zone.BaseDenom, rebalance.Amount)})
	}
	if len(msgs) == 0 {
		k.Logger(ctx).Info("No rebalancing required")
		return nil
	}
	k.Logger(ctx).Info("Send rebalancing messages", "msgs", msgs)
	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, "epoch %d rebalancing")
}

type RebalanceTarget struct {
	Amount math.Int
	Source string
	Target string
}

func DetermineAllocationsForRebalancing(currentAllocations map[string]math.Int, currentSum math.Int, targetAllocations map[string]*types.ValidatorIntent) []RebalanceTarget {
	out := make([]RebalanceTarget, 0)
	deltas := calculateDeltas(currentAllocations, currentSum, targetAllocations)

	wantToRebalance := sdk.ZeroInt()
	maxCanRebalance := currentSum.Quo(sdk.NewInt(2))

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].ValoperAddress < deltas[j].ValoperAddress
	})

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].Weight.GT(deltas[j].Weight)
	})

	for _, delta := range deltas {
		if delta.Weight.IsPositive() {
			wantToRebalance = wantToRebalance.Add(delta.Weight.TruncateInt())
		}
	}

	toRebalance := sdk.MinInt(wantToRebalance, maxCanRebalance)

	fmt.Println("deltas", deltas, wantToRebalance, maxCanRebalance, toRebalance)

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

		if tgt.Weight.Abs().TruncateInt().IsZero() { //nolint:gocritic
			tgtIdx++
			continue
		} else if tgt.Weight.Abs().TruncateInt().GT(toRebalance) {
			// amount == amount!
		} else {
			amount = sdk.MinInt(amount, tgt.Weight.Abs().TruncateInt())
		}
		out = append(out, RebalanceTarget{Amount: amount, Target: tgt.ValoperAddress, Source: src.ValoperAddress})
		deltas[srcIdx].Weight = src.Weight.Add(sdk.NewDecFromInt(amount))
		deltas[tgtIdx].Weight = tgt.Weight.Sub(sdk.NewDecFromInt(amount))
		toRebalance = toRebalance.Sub(amount)
		fmt.Printf("source: %s [%d], target : %s [%d], amount: %d, toRebalance: %d\n", src.ValoperAddress, src.Weight.TruncateInt().Int64(), tgt.ValoperAddress, tgt.Weight.TruncateInt().Int64(), amount.Int64(), toRebalance.Int64())

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
