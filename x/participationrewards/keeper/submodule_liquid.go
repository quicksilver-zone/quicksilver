package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type LiquidTokensModule struct{}

var _ Submodule = &LiquidTokensModule{}

func (*LiquidTokensModule) Hooks(_ sdk.Context, _ *Keeper) {
}

func (*LiquidTokensModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (math.Int, error) {
	// message
	// check denom is valid vs allowed

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

		// DenomFromRequestKey will error if the user address does not match the address in the key
		// or if the denom found is not valid.
		denom, err := utils.DenomFromRequestKey(proof.Key, addr)
		if err != nil {
			// check for mapped address for this user from SrcZone.
			mappedAddr, found := k.icsKeeper.GetRemoteAddressMap(ctx, addr, msg.SrcZone)
			if found {
				denom, err = utils.DenomFromRequestKey(proof.Key, mappedAddr)
				if err != nil {
					return sdk.ZeroInt(), errors.New("not a valid proof for submitting user or mapped account")
				}
			} else {
				return sdk.ZeroInt(), errors.New("not a valid proof for submitting user")
			}
		}

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
			coin, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, proof.Data, denomData.IbcDenom)
			if err != nil {
				return sdk.ZeroInt(), err
			}
			amount = amount.Add(coin.Amount)
		}
	}
	return amount, nil
}
