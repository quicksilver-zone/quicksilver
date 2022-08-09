package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (m *LiquidTokensModule) VerifyClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) error {
	// message
	// check denom is valid vs allowed

	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	if !ok {
		return fmt.Errorf("unable to find registered zone for chain id: %s", msg.Zone)
	}

	_, found := k.GetProtocolData(ctx, fmt.Sprintf("liquid/%s/%s", msg.Zone, zone.BaseDenom))
	if !found {
		return fmt.Errorf("unable to query liquid/%s/%s", msg.Zone, zone.BaseDenom)
	}

	return nil
}
