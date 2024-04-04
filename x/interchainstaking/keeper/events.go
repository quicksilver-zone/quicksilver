package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

// ___________________________________________________________________________________________________

type EventCallback func(*Keeper, sdk.Context, []byte) error

// Callbacks wrapper struct for interchainstaking keeper.
type EventCallbacks struct {
	k         *Keeper
	callbacks map[string]EventCallback
}

var _ emtypes.EventCallbacks = EventCallbacks{}

func (k *Keeper) EventCallbackHandler() EventCallbacks {
	return EventCallbacks{k, make(map[string]EventCallback)}
}

// Call calls callback handler.
func (c EventCallbacks) Call(ctx sdk.Context, id string, args []byte) error {
	if !c.Has(id) {
		return fmt.Errorf("callback %s not found", id)
	}
	return c.callbacks[id](c.k, ctx, args)
}

func (c EventCallbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c EventCallbacks) AddCallback(id string, fn interface{}) emtypes.EventCallbacks {
	c.callbacks[id], _ = fn.(EventCallback)
	return c
}

func (c EventCallbacks) RegisterCallbacks() emtypes.EventCallbacks {
	a := c.
		AddCallback("valset", EventCallback(TestCallback))

	return a.(EventCallbacks)
}

// -----------------------------------
// Callback Handlers
// -----------------------------------

func TestCallback(k *Keeper, ctx sdk.Context, args []byte) error {
	k.Logger(ctx).Error("TEST CALLBACK")
	return nil
}
