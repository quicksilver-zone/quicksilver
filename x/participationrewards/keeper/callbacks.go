package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	osmosisgammtypes "github.com/osmosis-labs/osmosis/v9/x/gamm/types"

	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// Callbacks wrapper struct for interchainstaking keeper
type Callback func(Keeper, sdk.Context, []byte, icqtypes.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]Callback
}

var _ icqtypes.QueryCallbacks = Callbacks{}

func (k Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]Callback)}
}

// callback handler
func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id] = fn.(Callback)
	return c
}

func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	a := c.
		AddCallback("validatorselectionrewards", Callback(ValidatorSelectionRewardsCallback)).
		AddCallback("osmosispoolupdate", Callback(OsmosisPoolUpdateCallback))

	return a.(Callbacks)
}

// Callbacks

func ValidatorSelectionRewardsCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	delegatorRewards := distrtypes.QueryDelegationTotalRewardsResponse{}
	err := k.cdc.Unmarshal(response, &delegatorRewards)
	if err != nil {
		return err
	}

	zone, found := k.icsKeeper.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	zs, err := k.getZoneScores(ctx, zone, delegatorRewards)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info(
		"callback zone score",
		"zone", zs.ZoneID,
		"total voting power", zs.TotalVotingPower,
		"validator scores", zs.ValidatorScores,
	)

	// snapshot obtained and used here
	userAllocations := k.calcUserValidatorSelectionAllocations(ctx, zone, *zs)

	if err := k.distributeToUsers(ctx, userAllocations); err != nil {
		return err
	}

	// create snapshot of current intents for next epoch boundary
	for _, di := range k.icsKeeper.AllIntents(ctx, zone, false) {
		k.icsKeeper.SetIntent(ctx, zone, di, true)
	}

	// set zone ValidatorSelectionAllocation to zero
	zone.ValidatorSelectionAllocation = sdk.NewCoins(
		sdk.NewCoin(
			k.stakingKeeper.BondDenom(ctx),
			sdk.ZeroInt(),
		),
	)
	k.icsKeeper.SetZone(ctx, &zone)

	return nil
}

func OsmosisPoolUpdateCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	var acc osmosisgammtypes.PoolI
	err := k.cdc.UnmarshalInterface(response, &acc)
	if err != nil {
		return err
	}
	// check query.Request is at least 9 bytes in length. (0x02 + 8 bytes for uint64)
	if len(query.Request) < 9 {
		return fmt.Errorf("query request not sufficient length")
	}
	// assert first character is 0x02 as expected.
	if query.Request[0] != 0x02 {
		return fmt.Errorf("query request has unexpected prefix")
	}

	poolID := sdk.BigEndianToUint64(query.Request[1:])
	data, ok := k.GetProtocolData(ctx, fmt.Sprintf("pools/%d", poolID))
	if !ok {
		return fmt.Errorf("unable to find protocol data for osmosis/pools/%d", poolID)
	}
	ipool, err := UnmarshalProtocolData("osmosispool", data.Data)
	if err != nil {
		return err
	}
	pool, ok := ipool.(types.OsmosisPoolProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for osmosis/pools/%d", poolID)
	}
	pool.PoolData = acc
	data.Data, err = json.Marshal(pool)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, fmt.Sprintf("osmosis/pools/%d", poolID), &data)

	return nil
}
