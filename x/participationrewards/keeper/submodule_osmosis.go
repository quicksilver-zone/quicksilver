package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	osmosistypes "github.com/ingenuity-build/quicksilver/osmosis-types"
	osmolockup "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type OsmosisModule struct{}

var _ Submodule = &OsmosisModule{}

func (m *OsmosisModule) Hooks(ctx sdk.Context, k *Keeper) {
	// osmosis params
	params, found := k.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisParams, types.OsmosisParamsKey)
	if !found {
		k.Logger(ctx).Error("unable to query osmosisparams in OsmosisModule hook")
		return
	}

	paramsData := types.OsmosisParamsProtocolData{}
	if err := json.Unmarshal(params.Data, &paramsData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal osmosisparams in OsmosisModule hook", "error", err)
		return
	}

	data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, paramsData.ChainID)
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("unable to query connection/%s in OsmosisModule hook", paramsData.ChainID))
		return
	}

	connectionData := types.ConnectionProtocolData{}
	if err := json.Unmarshal(data.Data, &connectionData); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("unable to unmarshal connection/%s in OsmosisModule hook", paramsData.ChainID))
		return
	}

	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeOsmosisPool), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, data.Data)
		if err != nil {
			return false
		}
		pool, _ := ipool.(*types.OsmosisPoolProtocolData)

		// update pool datas
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/gamm/key",
			m.GetKeyPrefixPools(pool.PoolID),
			sdk.NewInt(-1),
			types.ModuleName,
			OsmosisPoolUpdateCallbackID,
			0,
		) // query pool data
		return false
	})
}

func (m *OsmosisModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (uint64, error) {
	var amount uint64
	for _, proof := range msg.Proofs {
		lock := osmolockup.PeriodLock{}
		err := k.cdc.Unmarshal(proof.Data, &lock)
		if err != nil {
			return 0, err
		}

		_, lockupOwner, err := bech32.DecodeAndConvert(lock.Owner)
		if err != nil {
			return 0, err
		}

		if sdk.AccAddress(lockupOwner).String() != msg.UserAddress {
			return 0, errors.New("not a valid proof for submitting user")
		}

		sdkAmount, err := osmosistypes.DetermineApplicableTokensInPool(ctx, k, lock, msg.Zone)
		if err != nil {
			return 0, err
		}

		if sdkAmount.IsNil() || sdkAmount.IsNegative() {
			return 0, errors.New("unexpected amount")
		}
		amount += sdkAmount.Uint64()
	}
	return amount, nil
}

func (m *OsmosisModule) GetKeyPrefixPools(poolID uint64) []byte {
	return append([]byte{0x02}, sdk.Uint64ToBigEndian(poolID)...)
}
