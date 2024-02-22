package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
)

func GenerateQueryHash(connectionID, chainID, queryType string, request []byte, module string, callbackID string) string {
	return fmt.Sprintf("%x", crypto.Sha256(append([]byte(module+connectionID+chainID+queryType+callbackID), request...)))
}

// ----------------------------------------------------------------

func (k Keeper) NewQuery(
	module,
	connectionID,
	chainID,
	queryType string,
	request []byte,
	period math.Int,
	callbackID string,
	ttl uint64,
) *types.Query {
	return &types.Query{
		Id:           GenerateQueryHash(connectionID, chainID, queryType, request, module, callbackID),
		ConnectionId: connectionID,
		ChainId:      chainID,
		QueryType:    queryType,
		Request:      request,
		Period:       period,
		LastHeight:   sdk.ZeroInt(),
		CallbackId:   callbackID,
		Ttl:          ttl,
	}
}

// GetQuery returns query.
func (k Keeper) GetQuery(ctx sdk.Context, id string) (types.Query, bool) {
	query := types.Query{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return query, false
	}
	k.cdc.MustUnmarshal(bz, &query)
	return query, true
}

// SetQuery set query info.
func (k Keeper) SetQuery(ctx sdk.Context, query types.Query) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	bz := k.cdc.MustMarshal(&query)
	store.Set([]byte(query.Id), bz)
}

// DeleteQuery delete query info.
func (k Keeper) DeleteQuery(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	store.Delete([]byte(id))
}

// IterateQueries iterate through queries.
func (k Keeper) IterateQueries(ctx sdk.Context, fn func(index int64, queryInfo types.Query) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		query := types.Query{}
		k.cdc.MustUnmarshal(iterator.Value(), &query)
		stop := fn(i, query)

		if stop {
			break
		}
		i++
	}
}

// AllQueries returns every queryInfo in the store.
func (k Keeper) AllQueries(ctx sdk.Context) []types.Query {
	queries := []types.Query{}
	k.IterateQueries(ctx, func(_ int64, queryInfo types.Query) (stop bool) {
		queries = append(queries, queryInfo)
		return false
	})
	return queries
}
