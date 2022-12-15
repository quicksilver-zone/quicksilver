package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc       codec.Codec
	storeKey  storetypes.StoreKey
	callbacks map[string]types.QueryCallbacks
	IBCKeeper *ibckeeper.Keeper
}

// NewKeeper returns a new instance of zones Keeper
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, ibckeeper *ibckeeper.Keeper) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		callbacks: make(map[string]types.QueryCallbacks),
		IBCKeeper: ibckeeper,
	}
}

func (k *Keeper) SetCallbackHandler(module string, handler types.QueryCallbacks) error {
	_, found := k.callbacks[module]
	if found {
		return fmt.Errorf("callback handler already set for %s", module)
	}
	k.callbacks[module] = handler.RegisterCallbacks()
	return nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetDatapointForID(ctx sdk.Context, id string, result []byte, height math.Int) error {
	mapping := types.DataPoint{Id: id, RemoteHeight: height, LocalHeight: sdk.NewInt(ctx.BlockHeight()), Value: result}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixData)
	bz := k.cdc.MustMarshal(&mapping)
	store.Set([]byte(id), bz)
	return nil
}

func (k *Keeper) GetDatapointForID(ctx sdk.Context, id string) (types.DataPoint, error) {
	mapping := types.DataPoint{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixData)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return types.DataPoint{}, fmt.Errorf("unable to find data for id %s", id)
	}

	k.cdc.MustUnmarshal(bz, &mapping)
	return mapping, nil
}

// IterateDatapoints iterate through datapoints
func (k Keeper) IterateDatapoints(ctx sdk.Context, fn func(index int64, dp types.DataPoint) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixData)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		datapoint := types.DataPoint{}
		k.cdc.MustUnmarshal(iterator.Value(), &datapoint)
		stop := fn(i, datapoint)

		if stop {
			break
		}
		i++
	}
}

// DeleteQuery delete datapoint
func (k Keeper) DeleteDatapoint(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixData)
	store.Delete([]byte(id))
}

func (k *Keeper) GetDatapoint(ctx sdk.Context, module string, connectionID string, chainID string, queryType string, request []byte) (types.DataPoint, error) {
	id := GenerateQueryHash(connectionID, chainID, queryType, request, module)
	return k.GetDatapointForID(ctx, id)
}

func (k *Keeper) GetDatapointOrRequest(ctx sdk.Context, module string, connectionID string, chainID string, queryType string, request []byte, maxAge uint64) (types.DataPoint, error) {
	val, err := k.GetDatapoint(ctx, module, connectionID, chainID, queryType, request)
	if err != nil {
		// no datapoint
		k.MakeRequest(ctx, connectionID, chainID, queryType, request, sdk.NewInt(-1), "", "", maxAge)
		return types.DataPoint{}, errors.New("no data; query submitted")
	}

	if val.LocalHeight.LT(sdk.NewInt(ctx.BlockHeight() - int64(maxAge))) { // this is somewhat arbitrary; TODO: make this better
		k.MakeRequest(ctx, connectionID, chainID, queryType, request, sdk.NewInt(-1), "", "", maxAge)
		return types.DataPoint{}, errors.New("stale data; query submitted")
	}
	// check ttl
	return val, nil
}

func (k *Keeper) MakeRequest(ctx sdk.Context, connectionID string, chainID string, queryType string, request []byte, period math.Int, module string, callbackID string, ttl uint64) {
	k.Logger(ctx).Info(
		"MakeRequest",
		"connection_id", connectionID,
		"chain_id", chainID,
		"query_type", queryType,
		"request", request,
		"period", period,
		"module", module,
		"callback", callbackID,
		"ttl", ttl,
	)
	key := GenerateQueryHash(connectionID, chainID, queryType, request, module)
	existingQuery, found := k.GetQuery(ctx, key)
	if !found {
		if module != "" && callbackID != "" {
			if _, exists := k.callbacks[module]; !exists {
				err := fmt.Errorf("no callback handler registered for module %s", module)
				k.Logger(ctx).Error(err.Error())
				panic(err)
			}
			if exists := k.callbacks[module].Has(callbackID); !exists {
				err := fmt.Errorf("no callback %s registered for module %s", callbackID, module)
				k.Logger(ctx).Error(err.Error())
				panic(err)
			}
		}
		newQuery := k.NewQuery(ctx, module, connectionID, chainID, queryType, request, period, callbackID, ttl)
		k.SetQuery(ctx, *newQuery)
	} else {
		// a re-request of an existing query triggers resetting of height to trigger immediately.
		k.Logger(ctx).Info("re-request", "LastHeight", existingQuery.LastHeight)
		existingQuery.LastHeight = sdk.ZeroInt()
		k.SetQuery(ctx, existingQuery)
	}
}
