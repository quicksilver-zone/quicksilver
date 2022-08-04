package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type OsmosisModule struct{}

var _ Submodule = &OsmosisModule{}

func (m *OsmosisModule) Hooks(ctx sdk.Context, k Keeper) {
	// osmosis params
	data, found := k.GetProtocolData(ctx, "osmosis/connection")
	if !found {
		k.Logger(ctx).Error("unable to query osmosis/connection in OsmosisModule hook")
		return
	}
	connectionData := ConnectionProtocolData{}
	if err := json.Unmarshal(data.Data, &connectionData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal osmosis/connection in OsmosisModule hook", "error", err)
		return
	}

	k.IterateProtocolDatas(ctx, "osmosis/pools", func(idx int64, data types.ProtocolData) bool {
		ipool, err := UnmarshalProtocolData("osmosispool", data.Data)
		if err != nil {
			return false
		}
		pool, _ := ipool.(OsmosisPoolProtocolData)

		// update pool datas
		k.IcqKeeper.MakeRequest(ctx, connectionData.ConnectionId, connectionData.ChainId, "store/gamm/key", m.GetKeyPrefixPools(pool.PoolId), sdk.NewInt(-1), types.ModuleName, "osmosispoolupdate", 0) // query pool data
		return false
	})
}

func (m *OsmosisModule) IsActive() bool {
	return true
}

func (m *OsmosisModule) IsReady() bool {
	return true
}

func (m *OsmosisModule) VerifyClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) error {
	return nil
}

func (m *OsmosisModule) GetKeyPrefixPools(poolId uint64) []byte {
	return append([]byte{0x02}, sdk.Uint64ToBigEndian(poolId)...)
}
