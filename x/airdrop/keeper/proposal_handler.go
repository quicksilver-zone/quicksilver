package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

// HandleRegisterZoneDropProposal is a handler for executing a passed airdrop proposal.
func HandleRegisterZoneDropProposal(ctx sdk.Context, k Keeper, p *types.RegisterZoneDropProposal) error {
	k.SetZoneDrop(ctx, *p.ZoneDrop)

	for i, cr := range p.ClaimRecords {
		if err := k.SetClaimRecord(ctx, *cr); err != nil {
			return fmt.Errorf("invalid zonedrop proposal claim record: [%d] %w", i, err)
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
		sdk.NewEvent(
			types.EventTypeRegisterZoneDrop,
			sdk.NewAttribute(types.AttributeKeyZoneID, p.ZoneDrop.ChainId),
		),
	})

	return nil
}
