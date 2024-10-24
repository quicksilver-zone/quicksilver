package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	umee "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types"
	leveragetypes "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type UmeeModule struct{}

var _ Submodule = &UmeeModule{}

func (UmeeModule) Hooks(ctx sdk.Context, k *Keeper) {
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
			leveragetypes.KeyReserveAmount(reserves.Denom),
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
			leveragetypes.KeyInterestScalar(interest.Denom),
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
			leveragetypes.KeyUTokenSupply(supply.Denom),
			sdk.NewInt(-1),
			types.ModuleName,
			UmeeUTokenSupplyUpdateCallbackID,
			0,
		) // query utoken supply

		return false
	})

	// umee-types leverage module balance update
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeUmeeLeverageModuleBalance), func(idx int64, _ []byte, data types.ProtocolData) bool {
		ibalance, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeLeverageModuleBalance, data.Data)
		if err != nil {
			return false
		}
		balance, _ := ibalance.(*types.UmeeLeverageModuleBalanceProtocolData)
		accountPrefix := banktypes.CreateAccountBalancesPrefix(authtypes.NewModuleAddress(leveragetypes.LeverageModuleName))

		// update leverage module balance
		k.IcqKeeper.MakeRequest(
			ctx,
			connectionData.ConnectionID,
			connectionData.ChainID,
			icstypes.BankStoreKey,
			append(accountPrefix, balance.Denom...),
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
			leveragetypes.KeyAdjustedTotalBorrow(borrows.Denom),
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

func (UmeeModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (math.Int, error) {
	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to find registered zone for chain id: %s", msg.Zone)
	}

	addr, err := addressutils.AccAddressFromBech32(msg.UserAddress, "")
	if err != nil {
		return sdk.ZeroInt(), err
	}

	amount := sdk.ZeroInt()
	keyCache := make(map[string]bool)

	for _, proof := range msg.Proofs {
		if _, found := keyCache[string(proof.Key)]; found {
			continue
		}
		keyCache[string(proof.Key)] = true

		if proof.Data == nil {
			continue
		}

		udenom, err := getDenomFromProof(proof, addr)
		if err != nil {
			mappedAddr, found := k.icsKeeper.GetLocalAddressMap(ctx, addr, msg.SrcZone)
			if found {
				udenom, err = getDenomFromProof(proof, mappedAddr)
				if err != nil {
					return sdk.ZeroInt(), errors.New("not a valid proof for submitting user or mapped account")
				}
			} else {
				return sdk.ZeroInt(), errors.New("not a valid proof for submitting user")
			}
		}

		denom := leveragetypes.ToTokenDenom(udenom)

		data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, fmt.Sprintf("%s_%s", msg.SrcZone, denom))
		if !found {
			// we don't have a record for this denom, but this is okay, we don't want to submit records for every ibc denom.
			continue
		}
		denomData := types.LiquidAllowedDenomProtocolData{}
		err = json.Unmarshal(data.Data, &denomData)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		if denomData.QAssetDenom == zone.LocalDenom && denomData.IbcDenom == denom {
			uToken, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, proof.Data, udenom)
			if err != nil {
				return sdk.ZeroInt(), err
			}
			token, err := umee.ExchangeUToken(ctx, uToken, k)
			if err != nil {
				return sdk.ZeroInt(), err
			}
			amount = amount.Add(token.Amount)
		}
	}

	return amount, err
}
