package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"

	"github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc       codec.Codec
	storeKey  storetypes.StoreKey
	callbacks map[string]types.QueryCallbacks
	IBCKeeper *ibckeeper.Keeper
}

// NewKeeper returns a new instance of zones Keeper.
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, ibcKeeper *ibckeeper.Keeper) Keeper {
	if ibcKeeper == nil {
		panic("ibcKeeper is nil")
	}

	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		callbacks: make(map[string]types.QueryCallbacks),
		IBCKeeper: ibcKeeper,
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

func (k *Keeper) MakeRequest(
	ctx sdk.Context,
	connectionID,
	chainID,
	queryType string,
	request []byte,
	period math.Int,
	module string,
	callbackID string,
	ttl uint64,
) {
	k.Logger(ctx).Debug(
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
		newQuery := k.NewQuery(module, connectionID, chainID, queryType, request, period, callbackID, ttl)
		k.SetQuery(ctx, *newQuery)
	} else {
		// a re-request of an existing query triggers resetting of height to trigger immediately.
		k.Logger(ctx).Debug("re-request", "LastHeight", existingQuery.LastHeight)
		existingQuery.LastHeight = sdk.ZeroInt()
		k.SetQuery(ctx, existingQuery)
	}
}

// Heights

func (k *Keeper) SetLatestHeight(ctx sdk.Context, chainID string, height uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLatestHeight)
	store.Set([]byte(chainID), sdk.Uint64ToBigEndian(height))
}

func (k *Keeper) GetLatestHeight(ctx sdk.Context, chainID string) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLatestHeight)
	return sdk.BigEndianToUint64(store.Get([]byte(chainID)))
}
