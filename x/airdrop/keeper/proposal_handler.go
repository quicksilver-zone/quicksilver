package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/airdrop/types"
)

type ClaimRecords []types.ClaimRecord

// HandleRegisterZoneDropProposal is a handler for executing a passed airdrop proposal.
func HandleRegisterZoneDropProposal(ctx sdk.Context, k *Keeper, p *types.RegisterZoneDropProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	_, found := k.icsKeeper.GetZone(ctx, p.ZoneDrop.ChainId)
	if !found {
		return fmt.Errorf("zone not found, %q", p.ZoneDrop.ChainId)
	}

	if p.ZoneDrop.StartTime.Before(ctx.BlockTime()) {
		return errors.New("zone airdrop already started")
	}

	// decompress claim records
	crsb, err := types.Decompress(p.ClaimRecords)
	if err != nil {
		return err
	}

	// unmarshal json
	var crs ClaimRecords
	if err := json.Unmarshal(crsb, &crs); err != nil {
		return err
	}

	sumMax := uint64(0)
	// validate ClaimRecords and process
	for i, cr := range crs {
		if err := cr.ValidateBasic(); err != nil {
			return fmt.Errorf("claim record %d, %w", i, err)
		}

		if len(cr.ActionsCompleted) != 0 {
			return fmt.Errorf("invalid zonedrop proposal claim record [%d]: contains completed actions", i)
		}

		if cr.ChainId != p.ZoneDrop.ChainId {
			return fmt.Errorf("invalid zonedrop proposal claim record [%d]: chainID missmatch, expected %q got %q", i, p.ZoneDrop.ChainId, cr.ChainId)
		}

		sumMax += cr.MaxAllocation
	}

	// check allocations
	if sumMax > p.ZoneDrop.Allocation {
		return fmt.Errorf("sum of claim records max allocations (%v) exceed zone airdrop allocation (%v)", sumMax, p.ZoneDrop.Allocation)
	}

	// process ZoneDrop
	k.SetZoneDrop(ctx, *p.ZoneDrop)
	for i, cr := range crs {
		if err := k.SetClaimRecord(ctx, cr); err != nil {
			return fmt.Errorf("invalid zonedrop proposal claim record [%d]: %w", i, err)
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
