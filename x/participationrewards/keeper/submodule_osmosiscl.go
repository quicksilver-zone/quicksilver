package keeper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	osmosistypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types"
	osmoclmodel "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/model"
	osmocl "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type OsmosisClModule struct{}

var _ Submodule = &OsmosisClModule{}

func (m *OsmosisClModule) Hooks(ctx sdk.Context, k *Keeper) {
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

	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeOsmosisCLPool), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisCLPool, data.Data)
		if err != nil {
			return false
		}
		pool, _ := ipool.(*types.OsmosisClPoolProtocolData)

		// update pool datas
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/concentratedliquidity/key",
			m.KeyPool(pool.PoolID),
			sdk.NewInt(-1),
			types.ModuleName,
			OsmosisClPoolUpdateCallbackID,
			0,
		) // query pool data
		return false
	})
}

func (*OsmosisClModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (math.Int, error) {
	claimAmount := math.ZeroInt()
	var position osmoclmodel.Position

	addr, err := addressutils.AccAddressFromBech32(msg.UserAddress, "")
	if err != nil {
		return math.ZeroInt(), err
	}

	keyCache := make(map[string]bool)

	for _, proof := range msg.Proofs {
		if _, found := keyCache[string(proof.Key)]; found {
			continue
		}
		keyCache[string(proof.Key)] = true

		if proof.Data == nil {
			continue
		}

		position = osmoclmodel.Position{}
		err := k.cdc.Unmarshal(proof.Data, &position)
		if err != nil {
			return math.ZeroInt(), err
		}

		_, lockupOwner, err := bech32.DecodeAndConvert(position.Address)
		if err != nil {
			return math.ZeroInt(), err
		}

		if !bytes.Equal(lockupOwner, addr) {
			mappedAddr, found := k.icsKeeper.GetLocalAddressMap(ctx, addr, msg.SrcZone)
			if !found || !bytes.Equal(lockupOwner, mappedAddr) {
				return math.ZeroInt(), errors.New("not a valid proof for submitting user or mapped account")
			}
		}

		denom, found := k.ApplicableDenomForZone(ctx, msg.Zone)
		if !found {
			return math.ZeroInt(), errors.New("no applicable denom found for zone")
		}

		sdkAmount, err := osmosistypes.DetermineApplicableTokensInClPool(ctx, k, position, msg.Zone, denom)
		if err != nil {
			return math.ZeroInt(), err
		}

		if sdkAmount.IsNil() || sdkAmount.IsNegative() {
			return math.ZeroInt(), errors.New("unexpected amount")
		}
		claimAmount = claimAmount.Add(sdkAmount)
	}
	return claimAmount, nil
}

func (*OsmosisClModule) KeyPool(poolID uint64) []byte {
	return osmocl.KeyPool(poolID)
}

func (k *Keeper) ApplicableDenomForZone(ctx sdk.Context, chainID string) (denom string, found bool) {
	zone, found := k.icsKeeper.GetZone(ctx, chainID)
	if !found {
		return "", false
	}

	params, found := k.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisParams, types.OsmosisParamsKey)
	if !found {
		return "", false
	}

	paramsData := types.OsmosisParamsProtocolData{}
	if err := json.Unmarshal(params.Data, &paramsData); err != nil {
		return "", false
	}

	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeLiquidToken), func(idx int64, key []byte, data types.ProtocolData) bool {
		liquidToken, _ := types.UnmarshalProtocolData(types.ProtocolDataTypeLiquidToken, data.Data)
		liquidTokenData := liquidToken.(*types.LiquidAllowedDenomProtocolData)
		if liquidTokenData.ChainID == paramsData.ChainID && liquidTokenData.QAssetDenom == zone.LocalDenom {
			found = true
			denom = liquidTokenData.IbcDenom
			return true
		}
		return false
	})
	return denom, found
}
