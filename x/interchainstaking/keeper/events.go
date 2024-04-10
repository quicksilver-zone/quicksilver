package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
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
	return c.
		AddCallback(ICQEmitDelegatorDelegations, EventCallback(EmitDelegatorDelegations)).
		AddCallback(TriggerCalculateRedemptionRate, EventCallback(CalculateRedemptionRate))
}

const (
	ICQEmitDelegatorDelegations    = "ICQEmitDelegatorDelegations"
	TriggerCalculateRedemptionRate = "CalculateRedemptionRate"
)

// -----------------------------------
// Callback Handlers
// -----------------------------------

type DelegatorDelegationsParams struct {
	ChainID      string
	ConnectionID string
	Request      []byte
}

func EmitDelegatorDelegations(k *Keeper, ctx sdk.Context, args []byte) error {

	var params DelegatorDelegationsParams
	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	k.ICQKeeper.MakeRequest(
		ctx,
		params.ConnectionID,
		params.ChainID,
		"cosmos.staking.v1beta1.Query/DelegatorDelegations",
		params.Request,
		sdk.NewInt(-1),
		types.ModuleName,
		"delegations_epoch",
		0,
	)
	return nil
}

func CalculateRedemptionRate(k *Keeper, ctx sdk.Context, args []byte) error {
	zone, found := k.GetZone(ctx, string(args))
	if !found {
		return fmt.Errorf("unable to find zone %s", args)
	}
	return k.TriggerRedemptionRate(ctx, &zone)
}
