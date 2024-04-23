package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
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
		AddCallback(CalculateValues, EventCallback(CalculateTokenValues)).
		AddCallback(Submodules, EventCallback(SubmoduleHooks)).
		AddCallback(DistributeRewards, EventCallback(DistributeParticipationRewards))
}

const (
	CalculateValues   = "CalculateValues"
	Submodules        = "Submodules"
	DistributeRewards = "DistributeRewards"
)

// -----------------------------------
// Callback Handlers
// -----------------------------------

func CalculateTokenValues(k *Keeper, ctx sdk.Context, args []byte) error {
	defer k.EventManagerKeeper.MarkCompleted(ctx, types.ModuleName, "", "calc_tokens")

	tvs, err := k.CalcTokenValues(ctx)
	if err != nil {
		return err
	}

	err = k.SetZoneAllocations(ctx, tvs)
	if err != nil {
		return err
	}

	k.QueryValidatorDelegationPerformance(ctx)

	return nil
}

func SubmoduleHooks(k *Keeper, ctx sdk.Context, args []byte) error {
	for _, sub := range k.PrSubmodules {
		sub.Hooks(ctx, k)

	}
	return nil
}

func DistributeParticipationRewards(k *Keeper, ctx sdk.Context, args []byte) error {
	// calculate, based on latest token values
	// allocation based on calculations
	return nil
}
