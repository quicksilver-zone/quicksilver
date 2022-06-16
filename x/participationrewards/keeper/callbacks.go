package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

// Callbacks wrapper struct for interchainstaking keeper
type Callback func(Keeper, sdk.Context, []byte, types.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]Callback
}

var _ types.QueryCallbacks = Callbacks{}

func (k Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]Callback)}
}

//callback handler
func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query types.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) AddCallback(id string, fn interface{}) types.QueryCallbacks {
	c.callbacks[id] = fn.(Callback)
	return c
}

func (c Callbacks) RegisterCallbacks() types.QueryCallbacks {
	a := c.
		AddCallback("perfrewards", Callback(RewardsCallback))

	return a.(Callbacks)
}

// Callbacks

// zone callback >>>
func RewardsCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {

	delegatorRewards := distrtypes.QueryDelegationTotalRewardsResponse{}
	err := k.cdc.Unmarshal(response, &delegatorRewards)
	if err != nil {
		return err
	}

	zone, found := k.icsKeeper.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	zs, err := k.getZoneScores(ctx, zone, delegatorRewards)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info(
		"callback zone score",
		"zone", zs.ZoneId,
		"total voting power", zs.TotalVotingPower,
		"validator scores", zs.ValidatorScores,
	)

	zAllocation := zone.GetAllocation()

	userAllocations, err := k.calcUserAllocations(ctx, zone, *zs, zAllocation)
	if err != nil {
		return err
	}

	if err := k.distributeToUsers(ctx, userAllocations); err != nil {
		return err
	}

	// create snapshot of current intents for next epoch boundary
	for _, di := range k.icsKeeper.AllOrdinalizedIntents(ctx, zone, false) {
		k.icsKeeper.SetIntent(ctx, zone, di, true)
	}

	return nil
}
