package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	crescenttypes "github.com/quicksilver-zone/quicksilver/third-party-chains/crescent-types"
	liquiditytypes "github.com/quicksilver-zone/quicksilver/third-party-chains/crescent-types/liquidity/types"
	lpfarm "github.com/quicksilver-zone/quicksilver/third-party-chains/crescent-types/lpfarm"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type CrescentModule struct{}

var _ Submodule = &CrescentModule{}

func (CrescentModule) Hooks(ctx sdk.Context, k *Keeper) {
	// crescent-types params
	params, found := k.GetProtocolData(ctx, types.ProtocolDataTypeCrescentParams, types.CrescentParamsKey)
	if !found {
		k.Logger(ctx).Error("unable to query crescentparams in CrescentModule hook")
		return
	}

	paramsData := types.CrescentParamsProtocolData{}
	if err := json.Unmarshal(params.Data, &paramsData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal crescentparams in CrescentModule hook", "error", err)
		return
	}

	data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, paramsData.ChainID)
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("unable to query connection/%s in CrescentModule hook", paramsData.ChainID))
		return
	}

	connectionData := types.ConnectionProtocolData{}
	if err := json.Unmarshal(data.Data, &connectionData); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("unable to unmarshal connection/%s in CrescentModule hook", paramsData.ChainID))
		return
	}

	// update reserve address denom balance
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeCrescentReserveAddressBalance), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ibalance, err := types.UnmarshalProtocolData(types.ProtocolDataTypeCrescentReserveAddressBalance, data.Data)
		if err != nil {
			return false
		}
		balance, _ := ibalance.(*types.CrescentReserveAddressBalanceProtocolData)
		_, addrBytes, _ := bech32.DecodeAndConvert(balance.ReserveAddress)
		lookupKey := banktypes.CreateAccountBalancesPrefix(addrBytes)

		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			icstypes.BankStoreKey,
			append(lookupKey, []byte(balance.Denom)...),
			sdk.NewInt(-1),
			types.ModuleName,
			CrescentReserveBalanceUpdateCallbackID,
			0,
		)
		return false
	})

	// update pool data
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeCrescentPool), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeCrescentPool, data.Data)
		if err != nil {
			return false
		}
		pool, _ := ipool.(*types.CrescentPoolProtocolData)

		poolKey := liquiditytypes.GetPoolKey(pool.PoolID)

		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/liquidity/key",
			poolKey,
			sdk.NewInt(-1),
			types.ModuleName,
			CrescentPoolUpdateCallbackID,
			0,
		)
		return false
	})

	// update poolcoin supply
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeCrescentPoolCoinSupply), func(idx int64, _ []byte, data types.ProtocolData) bool {
		isupply, err := types.UnmarshalProtocolData(types.ProtocolDataTypeCrescentPoolCoinSupply, data.Data)
		if err != nil {
			return false
		}
		supply, _ := isupply.(*types.CrescentPoolCoinSupplyProtocolData)

		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			icstypes.BankStoreKey,
			append(banktypes.SupplyKey, []byte(supply.PoolCoinDenom)...),
			sdk.NewInt(-1),
			types.ModuleName,
			CrescentPoolCoinSupplyUpdateCallbackID,
			0,
		)
		return false
	})
}

func (CrescentModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (uint64, error) {
	var amount uint64
	for _, proof := range msg.Proofs {
		position := lpfarm.Position{}
		if proof.ProofType == types.ProofTypeBank {
			addr, poolDenom, err := banktypes.AddressAndDenomFromBalancesStore(proof.Key[1:])
			if err != nil {
				return 0, err
			}
			coin, err := keeper.UnmarshalBalanceCompat(k.cdc, proof.Data, poolDenom)
			if err != nil {
				return 0, err
			}
			address, err := addressutils.EncodeAddressToBech32("cre", addr)
			if err != nil {
				return 0, err
			}
			position = lpfarm.Position{
				Farmer:        address,
				Denom:         coin.GetDenom(),
				FarmingAmount: coin.Amount,
			}
		} else {
			err := k.cdc.Unmarshal(proof.Data, &position)
			if err != nil {
				return 0, err
			}
		}

		farmer, err := addressutils.AccAddressFromBech32(position.Farmer, "")
		if err != nil {
			return 0, err
		}

		user, _ := addressutils.AccAddressFromBech32(msg.UserAddress, "")
		if !farmer.Equals(user) {
			return 0, errors.New("not a valid proof for submitting user")
		}

		sdkAmount, err := crescenttypes.DetermineApplicableTokensInPool(ctx, k, position)
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
