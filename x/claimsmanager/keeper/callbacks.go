package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
)

const (
	ValidatorSelectionRewardsCallbackID       = "validatorselectionrewards"
	OsmosisPoolUpdateCallbackID               = "osmosispoolupdate"
	OsmosisClPoolUpdateCallbackID             = "osmosisclpoolupdate"
	SetEpochBlockCallbackID                   = "epochblock"
	UmeeReservesUpdateCallbackID              = "umeereservesupdatecallback"
	UmeeTotalBorrowsUpdateCallbackID          = "umeetotalborrowsupdatecallback"
	UmeeInterestScalarUpdateCallbackID        = "umeeinterestscalarupdatecallback"
	UmeeUTokenSupplyUpdateCallbackID          = "umeeutokensupplyupdatecallback"
	UmeeLeverageModuleBalanceUpdateCallbackID = "umeeleveragemodulebalanceupdatecallback"
)

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
	if !c.Has(id) {
		return fmt.Errorf("callback %s not found", id)
	}
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
	a := c.AddCallback(SetEpochBlockCallbackID, Callback(SetEpochBlockCallback))

	return a.(Callbacks)
}

// SetEpochBlockCallback records the block height of the registered zone at the epoch boundary.
func SetEpochBlockCallback(ctx sdk.Context, k *Keeper, args []byte, query icqtypes.Query) error {
	k.Logger(ctx).Debug("epoch callback called")
	data, connectionData, err := GetAndUnmarshalProtocolData[*types.ConnectionProtocolData](ctx, k, query.ChainId, types.ProtocolDataTypeConnection)
	if err != nil {
		return err
	}

	// block response is never expected to be nil
	if len(args) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (1)")
	}

	blockResponse := tmservice.GetLatestBlockResponse{}
	err = k.cdc.Unmarshal(args, &blockResponse)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug("got block response", "block", blockResponse)

	if blockResponse.SdkBlock == nil {
		// v0.45 and below
		// nolint:staticcheck // SA1019 ignore this!
		connectionData.LastEpoch = blockResponse.Block.Header.Height
	} else {
		// v0.46 and above
		connectionData.LastEpoch = blockResponse.SdkBlock.Header.Height
	}

	// todo update claimable events with all the heights collected here
	k.IteratePrefixedClaimableEvent(ctx, func(index int64, key []byte, data types.ClaimableEvent) (stop bool) {
		if ctx.BlockHeader().Time.Before(data.MaxClaimTime) {
			// Update the heights for each claimable event
			data.Heights[query.ChainId] = connectionData.LastEpoch // Update the height for the current chain

			// Save the updated claimable event back to the store
			if err = k.SetClaimableEvent(ctx, &data); err != nil {
				k.Logger(ctx).Error("Failed to update claimable event", "event", data.EventName, "error", err)
			}
			return false // Continue iterating till all events have updated the chain height
		}
		return false
	})

	heightInBytes := sdk.Uint64ToBigEndian(uint64(connectionData.LastEpoch)) //nolint:gosec
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
