package keeper

import (
	"fmt"


	sdk "github.com/cosmos/cosmos-sdk/types"

	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
)

// ___________________________________________________________________________________________________

type Callback func(*Keeper, sdk.Context, []byte, icqtypes.Query) error

// Callbacks wrapper struct for interchainstaking keeper.
type Callbacks struct {
	k         *Keeper
	callbacks map[string]Callback
}

var _ icqtypes.QueryCallbacks = Callbacks{}

func (k *Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]Callback)}
}

// Call calls callback handler.
func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	if !c.Has(id) {
		return fmt.Errorf("callback %s not found", id)
	}
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id], _ = fn.(Callback)
	return c
}

func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	return c
}
