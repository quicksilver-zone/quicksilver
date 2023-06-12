package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// Callback wrapper struct for interchainstaking keeper.
type Callback func(sdk.Context, *Keeper, []byte, icqtypes.Query) error

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
	return c.callbacks[id](ctx, c.k, args, query)
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
	a := c.
		AddCallback("validatorselectionrewards", Callback(ValidatorSelectionRewardsCallback)).
		AddCallback("osmosispoolupdate", Callback(OsmosisPoolUpdateCallback)).
		AddCallback("epochblock", Callback(SetEpochBlockCallback))

	return a.(Callbacks)
}

// Callbacks

func ValidatorSelectionRewardsCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
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
	userAllocations := k.CalcUserValidatorSelectionAllocations(ctx, &zone, *zs)

	if err := k.DistributeToUsers(ctx, userAllocations); err != nil {
		return err
	}

	// create snapshot of current intents for next epoch boundary
	for _, di := range k.icsKeeper.AllDelegatorIntents(ctx, &zone, false) {
		k.icsKeeper.SetDelegatorIntent(ctx, &zone, di, true)
	}

	// set zone ValidatorSelectionAllocation to zero
	zone.ValidatorSelectionAllocation = 0
	k.icsKeeper.SetZone(ctx, &zone)

	return nil
}

func OsmosisPoolUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	var pd gamm.PoolI
	if err := k.cdc.UnmarshalInterface(response, &pd); err != nil {
		return err
	}

	// check query.Request is at least 9 bytes in length. (0x02 + 8 bytes for uint64)
	if len(query.Request) < 9 {
		return errors.New("query request not sufficient length")
	}
	// assert first character is 0x02 as expected.
	if query.Request[0] != 0x02 {
		return errors.New("query request has unexpected prefix")
	}

	poolID := sdk.BigEndianToUint64(query.Request[1:])
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisPool, fmt.Sprintf("%d", poolID))
	if !ok {
		return fmt.Errorf("unable to find protocol data for osmosispools/%d", poolID)
	}
	ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, data.Data)
	if err != nil {
		return err
	}
	pool, ok := ipool.(*types.OsmosisPoolProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for osmosispools/%d", poolID)
	}
	pool.PoolData, err = json.Marshal(pd)
	if err != nil {
		return err
	}
	pool.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(pool)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, pool.GenerateKey(), &data)

	return nil
}

// SetEpochBlockCallback records the block height of the registered zone at the epoch boundary.
func SetEpochBlockCallback(ctx sdk.Context, k *Keeper, args []byte, query icqtypes.Query) error {
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, query.ChainId)
	if !ok {
		return fmt.Errorf("unable to find protocol data for connection/%s", query.ChainId)
	}
	k.Logger(ctx).Debug("epoch callback called")
	iConnectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, data.Data)
	connectionData, _ := iConnectionData.(*types.ConnectionProtocolData)

	if err != nil {
		return err
	}

	blockResponse := tmservice.GetLatestBlockResponse{}
	// block response is never expected to be nil
	if len(args) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (1)")
	}
	err = k.cdc.Unmarshal(args, &blockResponse)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug("got block response", "block", blockResponse)

	if blockResponse.SdkBlock == nil {
		// v0.45 and below
		//nolint:staticcheck // SA1019 ignore this!
		connectionData.LastEpoch = blockResponse.Block.Header.Height
	} else {
		// v0.46 and above
		connectionData.LastEpoch = blockResponse.SdkBlock.Header.Height
	}

	heightInBytes := sdk.Uint64ToBigEndian(uint64(connectionData.LastEpoch))
	// trigger a client update at the epoch boundary
	k.IcqKeeper.MakeRequest(
		ctx,
		query.ConnectionId,
		query.ChainId,
		"ibc.ClientUpdate",
		heightInBytes,
		sdk.NewInt(-1),
		types.ModuleName,
		"",
		0,
	)

	k.Logger(ctx).Debug("emitted client update", "height", connectionData.LastEpoch)

	data.Data, err = json.Marshal(connectionData)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, connectionData.GenerateKey(), &data)
	return nil
}
