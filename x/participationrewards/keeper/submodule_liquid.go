package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type LiquidTokensModule struct{}

var _ Submodule = &LiquidTokensModule{}

func (m *LiquidTokensModule) Hooks(ctx sdk.Context, k Keeper) {
}

func (m *LiquidTokensModule) IsActive() bool {
	return true
}

func (m *LiquidTokensModule) IsReady() bool {
	return true
}

func (m *LiquidTokensModule) ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (uint64, error) {
	// message
	// check denom is valid vs allowed

	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	if !ok {
		return 0, fmt.Errorf("unable to find registered zone for chain id: %s", msg.Zone)
	}

	_, addr, err := bech32.DecodeAndConvert(msg.UserAddress)
	if err != nil {
		return 0, err
	}

	amount := uint64(0)
	for _, proof := range msg.Proofs {
		// determine denoms from key
		if proof.Data == nil {
			continue
		}

		// DenomFromRequestKey will error if the user address does not match the address in the key
		// or if the denom found is not valid.
		denom, err := utils.DenomFromRequestKey(proof.Key, addr)
		if err != nil {
			return 0, err
		}

		data, found := k.GetProtocolData(ctx, fmt.Sprintf("liquid/%s/%s", msg.SrcZone, denom))
		if !found {
			// we don't have a record for this denom, but this is okay, we don't want to submit records for every ibc denom.
			continue
		}
		denomData := types.LiquidAllowedDenomProtocolData{}
		err = json.Unmarshal(data.Data, &denomData)
		if err != nil {
			return 0, err
		}
		if denomData.Denom == zone.BaseDenom && denomData.LocalDenom == denom {
			coin := sdk.Coin{}
			err = k.cdc.Unmarshal(proof.Data, &coin)
			if err != nil {
				return 0, err
			}
			amount += coin.Amount.Uint64()
		}

	}
	return amount, nil
}
