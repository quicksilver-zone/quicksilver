package keeper

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	umee "github.com/ingenuity-build/quicksilver/umee"
	umeetypes "github.com/ingenuity-build/quicksilver/umee/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type UmeeModule struct{}

var _ Submodule = &UmeeModule{}

func (u UmeeModule) Hooks(ctx sdk.Context, k *Keeper) {
	//update module balance, reserves, utoken supply
	// umee params
	params, found := k.GetProtocolData(ctx, types.ProtocolDataTypeUmeeParams, types.UmeeParamsKey)
	if !found {
		k.Logger(ctx).Error("unable to query osmosisparams in OsmosisModule hook")
		return
	}

	paramsData := types.UmeeParamsProtocolData{}
	if err := json.Unmarshal(params.Data, &paramsData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal umeeparams in UmeeModule hook", "error", err)
		return
	}

	data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, paramsData.ChainID)
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("unable to query connection/%s in OsmosisModule hook", paramsData.ChainID))
		return
	}

	connectionData := types.ConnectionProtocolData{}
	if err := json.Unmarshal(data.Data, &connectionData); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("unable to unmarshal connection/%s in UmeeModule hook", paramsData.ChainID))
		return
	}

	// umee reserves update
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeReserves), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ireserves, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeReserves, data.Data)
		if err != nil {
			return false
		}
		reserves, _ := ireserves.(*types.UmeeReservesProtocolData)

		//update reserves
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/leverage/key",
			umeetypes.KeyReserveAmount(reserves.Denom),
			sdk.NewInt(-1),
			types.ModuleName,
			UmeeReservesUpdateCallbackID,
			0,
		) // query reserve data
		return false
	})
	// umee interest scalar update
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeInterestScalar), func(idx int64, _ []byte, data types.ProtocolData) bool {
		iinterest, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeInterestScalar, data.Data)
		if err != nil {
			return false
		}
		interest, _ := iinterest.(*types.UmeeInterestScalarProtocolData)

		// update interest
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/leverage/key",
			umeetypes.KeyInterestScalar(interest.Denom),
			sdk.NewInt(-1),
			types.ModuleName,
			UmeeInterestScalarUpdateCallbackID,
			0,
		) // query interest data

		return false
	})
	// umee utoken supply update
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeUTokenSupply), func(idx int64, _ []byte, data types.ProtocolData) bool {
		isupply, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeUTokenSupply, data.Data)
		if err != nil {
			return false
		}
		supply, _ := isupply.(*types.UmeeUTokenSupplyProtocolData)

		// update utoken supply
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/leverage/key",
			umeetypes.KeyUTokenSupply(supply.Denom),
			sdk.NewInt(-1),
			types.ModuleName,
			"umeereservesupdate",
			0,
		) // query utoken supply

		return false
	})

	//TODO: check module balance retrieval
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeReserves), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ireserves, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeReserves, data.Data)
		if err != nil {
			return false
		}
		reserves, _ := ireserves.(*types.UmeeReservesProtocolData)

		// update leverage module balance
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/leverage/key",
			umeetypes.KeyUTokenSupply(reserves.Denom), //wrong key
			sdk.NewInt(-1),
			types.ModuleName,
			"umeereservesupdate",
			0,
		) // query leverage module balance

		return false
	})
	// umee total borrowed update
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeTotalBorrows), func(idx int64, _ []byte, data types.ProtocolData) bool {
		iborrows, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeTotalBorrows, data.Data)
		if err != nil {
			return false
		}
		borrows, _ := iborrows.(*types.UmeeTotalBorrowsProtocolData)

		// update total borrows for a denom
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			"store/leverage/key",
			umeetypes.KeyAdjustedTotalBorrow(borrows.Denom), //wrong key
			sdk.NewInt(-1),
			types.ModuleName,
			UmeeTotalBorrowsUpdateCallbackID,
			0,
		) // query leverage module balance

		return false
	})

}

func (u UmeeModule) IsActive() bool {
	return true
}

func (u UmeeModule) IsReady() bool {
	return true
}

func (u UmeeModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (uint64, error) {
	//the claim will have some u/assets in it, convert them to rewards, this will just contain the conversion logic
	uToken := sdk.Coin{}
	token, err := umee.ExchangeUToken(ctx, uToken, k)

	return token.Amount.Uint64(), err
}
