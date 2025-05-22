package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type MembraneModule struct{}

var _ Submodule = &MembraneModule{}

func (*MembraneModule) Hooks(_ sdk.Context, _ *Keeper) {
}

func (*MembraneModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (math.Int, error) {
	// message
	// check denom is valid vs allowed

	params, found := k.GetProtocolData(ctx, types.ProtocolDataTypeMembraneParams, types.MembraneParamsKey)
	if !found {
		k.Logger(ctx).Error("unable to query membraneparams in MembraneModule")
		return sdk.ZeroInt(), errors.New("unable to query membraneparams in MembraneModule")
	}

	paramsData := types.MembraneProtocolData{}
	if err := json.Unmarshal(params.Data, &paramsData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal membraneparams in MembraneModule", "error", err)
		return sdk.ZeroInt(), err
	}

	osmosisParams, found := k.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisParams, types.OsmosisParamsKey)
	if !found {
		k.Logger(ctx).Error("unable to query osmosisparams in MembraneModule")
		return sdk.ZeroInt(), errors.New("unable to query osmosisparams in MembraneModule")
	}

	osmosisParamsData := types.OsmosisParamsProtocolData{}
	if err := json.Unmarshal(osmosisParams.Data, &osmosisParamsData); err != nil {
		k.Logger(ctx).Error("unable to unmarshal osmosisparams in MembraneModule", "error", err)
		return sdk.ZeroInt(), err
	}

	if msg.SrcZone != osmosisParamsData.ChainID {
		return sdk.ZeroInt(), errors.New("src zone does not match osmosis chain id")
	}

	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	if !ok {
		return sdk.ZeroInt(), fmt.Errorf("unable to find registered zone for chain id: %s", msg.Zone)
	}

	_, addr, err := bech32.DecodeAndConvert(msg.UserAddress)
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

		// Validate the address is valid from the key.
		contractAddr, parts, err := utils.DecodeCwNamespacedKey(proof.Key, 2)
		if err != nil {
			return sdk.ZeroInt(), err
		}

		if !contractAddr.Equals(sdk.AccAddress(paramsData.ContractAddress)) {
			return sdk.ZeroInt(), errors.New("not a valid membrane contract address")
		}

		if len(parts) != 2 {
			return sdk.ZeroInt(), errors.New("not a key length for membrane claims")
		}

		if string(parts[0]) != "positions" {
			return sdk.ZeroInt(), errors.New("not a valid key for membrane claims")
		}

		if sdk.AccAddress(parts[1]).Equals(sdk.AccAddress(addr)) {
			mappedAddr, found := k.icsKeeper.GetLocalAddressMap(ctx, addr, msg.SrcZone)
			if found {
				if sdk.AccAddress(parts[1]).Equals(mappedAddr) {
					return sdk.ZeroInt(), errors.New("not a valid key for submitting user")
				}
			} else {
				return sdk.ZeroInt(), errors.New("not a valid key for submitting user (mapped address not found)")
			}
		}

		denom, err := utils.DenomFromRequestKey(proof.Key, addr)
		if err != nil {
			// check for mapped address for this user from SrcZone.
			mappedAddr, found := k.icsKeeper.GetLocalAddressMap(ctx, addr, msg.SrcZone)
			if found {
				denom, err = utils.DenomFromRequestKey(proof.Key, mappedAddr)
				if err != nil {
					return sdk.ZeroInt(), errors.New("not a valid proof for submitting user or mapped account")
				}
			} else {
				return sdk.ZeroInt(), errors.New("not a valid proof for submitting user")
			}
		}

		var positions []types.MembranePosition
		err = json.Unmarshal(proof.Data, &positions)
		if err != nil {
			return sdk.ZeroInt(), err
		}

		for _, position := range positions {
			for _, collateralAsset := range position.CollateralAssets {
				data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeLiquidToken, fmt.Sprintf("%s_%s", msg.SrcZone, collateralAsset.Asset.Info.NativeToken.Denom))

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
					amount = amount.Add(collateralAsset.Asset.Amount)
				}

			}
		}

	}
	return amount, nil
}
