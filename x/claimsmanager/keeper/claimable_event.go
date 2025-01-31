package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	participationrewardstypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func (k *Keeper) SetClaimableEvent(ctx sdk.Context, claimableEvent *types.ClaimableEvent) error {
	if claimableEvent.MaxClaimTime.Before(ctx.BlockHeader().Time) {
		return fmt.Errorf("max time cannot be less than current block time")
	}

	// todo : 1. call the store consensus state here after each call
	// this might include callbacks etc
	switch claimableEvent.EventName {
	case epochstypes.ModuleName:
		if err := k.StoreSelfConsensusState(ctx, "epoch"); err != nil {
			return err
		}
	case participationrewardstypes.ModuleName:
	// zones: can get it from protocol data
	default:
	}

	var listOfChains []string
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeConnection), func(index int64, _ []byte, data types.ProtocolData) (stop bool) {
		iConnectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, data.Data)
		if err != nil {
			k.Logger(ctx).Error("Error unmarshalling protocol data")
		}
		connectionData, _ := iConnectionData.(*types.ConnectionProtocolData)
		listOfChains = append(listOfChains, connectionData.ChainID)
		return
	})

	for _, chainID := range listOfChains {
		claimableEvent.Heights[chainID] = 0
	}

	//2.  store rest of the data as is for the claimable event
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	bz := k.cdc.MustMarshal(claimableEvent)
	store.Set(types.GetKeyClaimableEvent(claimableEvent.EventModule, claimableEvent.EventName), bz)
	return nil
}

func (k *Keeper) GetClaimableEvent(ctx sdk.Context, eventModule, eventName string) (types.ClaimableEvent, bool) {
	data := types.ClaimableEvent{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	key := types.GetGenericKeyClaimableEvent(eventModule, eventName)
	bz := store.Get(key)
	if len(bz) == 0 {
		return data, false
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

func (k *Keeper) IteratePrefixedClaimableEvent(ctx sdk.Context, fn func(index int64, key []byte, data types.ClaimableEvent) (stop bool)) {
	if fn == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixClaimableEvent)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		data := types.ClaimableEvent{}
		k.cdc.MustUnmarshal(iterator.Value(), &data)
		stop := fn(i, iterator.Key(), data)
		if stop {
			break
		}
		i++

	}
}

func (k *Keeper) DeleteClaimableEvent(ctx sdk.Context, eventModule, eventName string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)
	key := types.GetGenericKeyClaimableEvent(eventModule, eventName)

	store.Delete(key)
}
