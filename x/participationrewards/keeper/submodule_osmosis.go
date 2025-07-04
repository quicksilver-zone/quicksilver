package keeper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"

	osmosistypes "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types"
	osmolockup "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/lockup/types"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
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

func (*OsmosisModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (math.Int, error) {
	amount := sdk.ZeroInt()
	var lock osmolockup.PeriodLock

	addr, err := addressutils.AccAddressFromBech32(msg.UserAddress, "")
	if err != nil {
		return sdk.ZeroInt(), err
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

		if proof.ProofType == types.ProofTypeBank {
			poolDenom, err := utils.DenomFromRequestKey(proof.Key, addr)
			if err != nil {
				// check for mapped address for this user from SrcZone.
				mappedAddr, found := k.icsKeeper.GetRemoteAddressMap(ctx, addr, msg.SrcZone)
				if found {
					poolDenom, err = utils.DenomFromRequestKey(proof.Key, mappedAddr)
					if err != nil {
						return sdk.ZeroInt(), errors.New("not a valid proof for submitting user or mapped account")
					}
				} else {
					return sdk.ZeroInt(), errors.New("not a valid proof for submitting user")
				}
			}

			coin, err := keeper.UnmarshalBalanceCompat(k.cdc, proof.Data, poolDenom)
			if err != nil {
				return sdk.ZeroInt(), err
			}
			poolID, err := strconv.ParseUint(poolDenom[strings.LastIndex(poolDenom, "/")+1:], 10, 64)
			if err != nil {
				return sdk.ZeroInt(), err
			}
			lock = osmolockup.PeriodLock{
				ID:       poolID,
				Owner:    addr.String(),
				Duration: time.Hour,
				EndTime:  time.Time{},
				Coins:    sdk.NewCoins(coin),
			}
		} else {
			lock = osmolockup.PeriodLock{}
			err := k.cdc.Unmarshal(proof.Data, &lock)
			if err != nil {
				return sdk.ZeroInt(), err
			}

			_, lockupOwner, err := bech32.DecodeAndConvert(lock.Owner)
			if err != nil {
				return sdk.ZeroInt(), err
			}

			if !bytes.Equal(lockupOwner, addr) {
				mappedAddr, found := k.icsKeeper.GetLocalAddressMap(ctx, addr, msg.SrcZone)
				if !found || !bytes.Equal(lockupOwner, mappedAddr) {
					return sdk.ZeroInt(), errors.New("not a valid proof for submitting user or mapped account")
				}
			}
		}

		denom, found := k.ApplicableDenomForZone(ctx, msg.Zone)
		if !found {
			return math.ZeroInt(), errors.New("no applicable denom found for zone")
		}

		sdkAmount, err := osmosistypes.DetermineApplicableTokensInPool(ctx, k, lock, msg.Zone, denom)
		if err != nil {
			return sdk.ZeroInt(), err
		}

		if sdkAmount.IsNil() || sdkAmount.IsNegative() {
			return sdk.ZeroInt(), errors.New("unexpected amount")
		}
		amount = amount.Add(sdkAmount)
	}
	return amount, nil
}

func (*OsmosisModule) GetKeyPrefixPools(poolID uint64) []byte {
	return append([]byte{0x02}, sdk.Uint64ToBigEndian(poolID)...)
}
