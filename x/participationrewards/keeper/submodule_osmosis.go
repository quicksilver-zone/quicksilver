package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	osmosistypes "github.com/ingenuity-build/quicksilver/osmosis-types"
	osmolockup "github.com/ingenuity-build/quicksilver/osmosis-types/lockup"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type OsmosisModule struct{}

var _ Submodule = &OsmosisModule{}

func (m *OsmosisModule) Hooks(ctx sdk.Context, k Keeper) {
	// osmosis params
	data, found := k.GetProtocolData(ctx, "connection/osmosis")
	if !found {
		k.Logger(ctx).Error("unable to query connection/osmosis in OsmosisModule hook")
		return
	}
	connectionData := types.ConnectionProtocolData{}
	if err := json.Unmarshal(data.Data, &connectionData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal osmosis/connection in OsmosisModule hook", "error", err)
		return
	}

	k.IteratePrefixedProtocolDatas(ctx, "osmosispools", func(idx int64, data types.ProtocolData) bool {
		ipool, err := types.UnmarshalProtocolData(types.ProtocolDataOsmosisPool, data.Data)
		if err != nil {
			return false
		}
		pool, _ := ipool.(types.OsmosisPoolProtocolData)

		// update pool datas
		k.IcqKeeper.MakeRequest(ctx, connectionData.ConnectionID, connectionData.ChainID, "store/gamm/key", m.GetKeyPrefixPools(pool.PoolID), sdk.NewInt(-1), types.ModuleName, "osmosispoolupdate", 0) // query pool data
		return false
	})
}

func (m *OsmosisModule) IsActive() bool {
	return true
}

func (m *OsmosisModule) IsReady() bool {
	return true
}

func (m *OsmosisModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (uint64, error) {
	var amount uint64
	for _, proof := range msg.Proofs {
		lockupResponse := osmolockup.LockedResponse{}
		err := k.cdc.Unmarshal(proof.Data, &lockupResponse)
		if err != nil {
			return 0, err
		}

		sdkAmount, err := osmosistypes.DetermineApplicableTokensInPool(ctx, k, lockupResponse, msg.Zone)
		if err != nil {
			return 0, err
		}
		amount += sdkAmount.Uint64()
	}
	return amount, nil
}

func (m *OsmosisModule) GetKeyPrefixPools(poolID uint64) []byte {
	return append([]byte{0x02}, sdk.Uint64ToBigEndian(poolID)...)
}
