package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/eventmanager/keeper"
	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

var GLOBAL_VAR = 0

// ___________________________________________________________________________________________________

type EventCallback func(*keeper.Keeper, sdk.Context, []byte) error

// Callbacks wrapper struct for interchainstaking keeper.
type EventCallbacks struct {
	k         *keeper.Keeper
	callbacks map[string]EventCallback
}

var _ types.EventCallbacks = EventCallbacks{}

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

func (c EventCallbacks) AddCallback(id string, fn interface{}) types.EventCallbacks {
	c.callbacks[id], _ = fn.(EventCallback)
	return c
}

func (c EventCallbacks) RegisterCallbacks() types.EventCallbacks {
	a := c.
		AddCallback("testCallback", EventCallback(testCallback))

	return a.(EventCallbacks)
}

// -----------------------------------
// Callback Handlers
// -----------------------------------

func testCallback(k *keeper.Keeper, ctx sdk.Context, args []byte) error {
	GLOBAL_VAR = 12345
	return nil
}

// tests

func (suite *KeeperTestSuite) TestEventLifecycle() {
	app := suite.GetSimApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	callbackHandler := EventCallbacks{&app.EventManagerKeeper, make(map[string]EventCallback, 0)}

	app.EventManagerKeeper.SetCallbackHandler(types.ModuleName, callbackHandler)

	app.EventManagerKeeper.AddEvent(ctx, types.ModuleName, suite.chainB.ChainID, "test", "testCallback", types.EventTypeICADelegate, types.EventStatusPending, nil, nil)

	events := app.EventManagerKeeper.AllEvents(ctx)

	suite.Equal(1, len(events))

	GLOBAL_VAR = 0

	app.EventManagerKeeper.Trigger(ctx, types.ModuleName, suite.chainB.ChainID)

	event, found := app.EventManagerKeeper.GetEvent(ctx, types.ModuleName, suite.chainB.ChainID, "test")

	suite.True(found)
	suite.Equal(12345, GLOBAL_VAR)

	suite.Equal(event.EventStatus, types.EventStatusActive)

	app.EventManagerKeeper.DeleteEvent(ctx, types.ModuleName, suite.chainB.ChainID, "test")

	events = app.EventManagerKeeper.AllEvents(ctx)

	suite.Equal(0, len(events))
}

func (suite *KeeperTestSuite) TestEventLifecycleWithCondition() {
	app := suite.GetSimApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	callbackHandler := EventCallbacks{&app.EventManagerKeeper, make(map[string]EventCallback, 0)}

	app.EventManagerKeeper.SetCallbackHandler(types.ModuleName, callbackHandler)

	condition := &types.ConditionAll{Fields: []*types.FieldValue{{Field: types.FieldModule, Value: types.ModuleName, Negate: true}}, Negate: false}

	app.EventManagerKeeper.AddEvent(ctx, types.ModuleName, suite.chainB.ChainID, "test", "testCallback", types.EventTypeICADelegate, types.EventStatusPending, condition, nil)

	events := app.EventManagerKeeper.AllEvents(ctx)

	suite.Equal(1, len(events))

	GLOBAL_VAR = 0

	app.EventManagerKeeper.Trigger(ctx, types.ModuleName, suite.chainB.ChainID)

	event, found := app.EventManagerKeeper.GetEvent(ctx, types.ModuleName, suite.chainB.ChainID, "test")

	fmt.Println(event)
	suite.True(found)
	suite.Equal(12345, GLOBAL_VAR)

	suite.Equal(event.EventStatus, types.EventStatusActive)

	app.EventManagerKeeper.DeleteEvent(ctx, types.ModuleName, suite.chainB.ChainID, "test")

	events = app.EventManagerKeeper.AllEvents(ctx)

	suite.Equal(0, len(events))
}
