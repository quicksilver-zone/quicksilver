package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

// Callbacks wrapper struct for participationrewards keeper
type Callback func(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]Callback
}

func (k Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]Callback)}
}

// // callback handler
// func (c Callbacks) Call(id string, ctx sdk.Context, args proto.Message) error {
// 	return c.callbacks[id](c.k, ctx, args)
// }
//callback handler
func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) AddCallback(id string, fn interface{}) {
	c.callbacks[id] = fn.(Callback)
}

func (c Callbacks) RemoveCallback(id string) {
	delete(c.callbacks, id)
}
