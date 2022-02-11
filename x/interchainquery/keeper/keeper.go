package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc      codec.Codec
	storeKey sdk.StoreKey
}

// NewKeeper returns a new instance of zones Keeper
func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// func (k *Keeper) SetConnectionForPort(ctx sdk.Context, connectionId string, port string) error {
// 	mapping := types.PortConnectionTuple{ConnectionId: connectionId, PortId: port}
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortMapping)
// 	bz := k.cdc.MustMarshal(&mapping)
// 	store.Set([]byte(port), bz)
// 	return nil
// }

// func (k *Keeper) GetConnectionForPort(ctx sdk.Context, port string) (string, error) {
// 	mapping := types.PortConnectionTuple{}
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPortMapping)
// 	bz := store.Get([]byte(port))
// 	if len(bz) == 0 {
// 		return "", fmt.Errorf("unable to find mapping for port %s", port)
// 	}

// 	k.cdc.MustUnmarshal(bz, &mapping)
// 	return mapping.ConnectionId, nil
// }
