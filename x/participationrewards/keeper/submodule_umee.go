package keeper

import (
	"encoding/json"
	"fmt"

	umee "github.com/ingenuity-build/quicksilver/umee-types"

	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	umeetypes "github.com/ingenuity-build/quicksilver/umee-types/leverage/types"
	"github.com/ingenuity-build/quicksilver/utils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type UmeeModule struct{}

var _ Submodule = &UmeeModule{}

func (u UmeeModule) Hooks(ctx sdk.Context, k *Keeper) {
	// umee-types params
	params, found := k.GetProtocolData(ctx, types.ProtocolDataTypeUmeeParams, types.UmeeParamsKey)
	if !found {
		k.Logger(ctx).Error("unable to query umeeparams in UmeeModule hook")
		return
	}

	paramsData := types.UmeeParamsProtocolData{}
	if err := json.Unmarshal(params.Data, &paramsData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal umeeparams in UmeeModule hook", "error", err)
		return
	}

	data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, paramsData.ChainID)
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("unable to query connection/%s in UmeeModule hook", paramsData.ChainID))
		return
	}

	connectionData := types.ConnectionProtocolData{}
	if err := json.Unmarshal(data.Data, &connectionData); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("unable to unmarshal connection/%s in UmeeModule hook", paramsData.ChainID))
		return
	}

	// umee-types reserves update
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeReserves), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ireserves, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeReserves, data.Data)
		if err != nil {
			return false
		}
		reserves, _ := ireserves.(*types.UmeeReservesProtocolData)

		// update reserves
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
	// umee-types interest scalar update
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
	// umee-types utoken supply update
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
			UmeeUTokenSupplyUpdateCallbackID,
			0,
		) // query utoken supply

		return false
	})

	// TODO: check module spendable coins retrieval
	// assuming that module account is not a vesting account so there
	// will be no locked coins to subtract from the total balance
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeLeverageModuleBalance), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ibalance, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeLeverageModuleBalance, data.Data)
		if err != nil {
			return false
		}
		balance, _ := ibalance.(*types.UmeeLeverageModuleBalanceProtocolData)
		accountPrefix := banktypes.CreateAccountBalancesPrefix(authtypes.NewModuleAddress(umeetypes.LeverageModuleName))

		// update leverage module balance
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			icstypes.BankStoreKey,
			append(accountPrefix, []byte(balance.Denom)...),
			sdk.NewInt(-1),
			types.ModuleName,
			UmeeLeverageModuleBalanceUpdateCallbackID,
			0,
		) // query leverage module balance

		return false
	})
	// umee-types total borrowed update
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
			umeetypes.KeyAdjustedTotalBorrow(borrows.Denom),
			sdk.NewInt(-1),
			types.ModuleName,
			UmeeTotalBorrowsUpdateCallbackID,
			0,
		) // query leverage module balance

		return false
	})
}

func getDenomFromProof(proof *cmtypes.Proof, addr []byte) (string, error) {
	denom, err := utils.DenomFromRequestKey(proof.Key, addr)
	if err != nil {
		return "", err
	}
	if proof.ProofType == types.ProofTypeLeverage {
		denom = denom[:len(denom)-1]
	}
	return denom, err
}

func (u UmeeModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (uint64, error) {
	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	if !ok {
		return 0, fmt.Errorf("unable to find registered zone for chain id: %s", msg.Zone)
	}

	_, addr, err := bech32.DecodeAndConvert(msg.UserAddress)

	amount := uint64(0)
	for _, proof := range msg.Proofs {
		// determine denoms from keys
		if proof.Data == nil {
			continue
		}

		udenom, err := getDenomFromProof(proof, addr)
		if err != nil {
			return 0, err
		}

		denom := umeetypes.ToTokenDenom(udenom)

		data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, fmt.Sprintf("%s_%s", msg.SrcZone, denom))
		if !found {
			// we don't have a record for this denom, but this is okay, we don't want to submit records for every ibc denom.
			continue
		}
		denomData := types.LiquidAllowedDenomProtocolData{}
		err = json.Unmarshal(data.Data, &denomData)
		if err != nil {
			return 0, err
		}
		if denomData.QAssetDenom == zone.LocalDenom && denomData.IbcDenom == denom {
			uToken, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, proof.Data, udenom)
			if err != nil {
				return 0, err
			}
			token, err := umee.ExchangeUToken(ctx, uToken, k)
			if err != nil {
				return 0, err
			}
			amount += token.Amount.Uint64()
		}
	}

	return amount, err
}
