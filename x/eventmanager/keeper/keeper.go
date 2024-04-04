package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc       codec.Codec
	storeKey  storetypes.StoreKey
	callbacks map[string]types.EventCallbacks
}

// NewKeeper returns a new instance of zones Keeper.
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey) Keeper {

	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		callbacks: make(map[string]types.EventCallbacks),
	}
}

func (k *Keeper) SetCallbackHandler(module string, handler types.EventCallbacks) error {
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

func (k Keeper) Call(ctx sdk.Context, moduleName string, callbackID string, payload []byte) error {

	module, found := k.callbacks[moduleName]
	if !found {
		return fmt.Errorf("bad module %s", moduleName)
	}
	if module.Has(callbackID) {
		// we have executed a callback; only a single callback is expected per request, so break here.
		return module.Call(ctx, callbackID, payload)
	}

	return fmt.Errorf("callback %s not found for module %s", callbackID, moduleName)
}
