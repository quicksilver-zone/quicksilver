package keeper

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

type ClaimRecords []types.ClaimRecord

// HandleRegisterZoneDropProposal is a handler for executing a passed airdrop proposal.
func HandleRegisterZoneDropProposal(ctx sdk.Context, k Keeper, p *types.RegisterZoneDropProposal) error {
	// validate ZoneDrop
	if err := p.ZoneDrop.ValidateBasic(); err != nil {
		return err
	}

	// check for ClaimRecords
	if p.ClaimRecords == nil {
		return types.ErrUndefinedAttribute
	}
	if len(p.ClaimRecords) == 0 {
		return types.ErrUndefinedAttribute
	}

	// decompress claim records
	crsb, err := k.decompress(p.ClaimRecords)
	if err != nil {
		return err
	}

	// unmarshal json
	var crs ClaimRecords
	if err := json.Unmarshal(crsb, &crs); err != nil {
		return err
	}

	// validate ClaimRecords and process
	for i, cr := range crs {
		if err := cr.ValidateBasic(); err != nil {
			return fmt.Errorf("claim record %d, %w", i, err)
		}

		if len(cr.ActionsCompleted) != 0 {
			return fmt.Errorf("invalid zonedrop proposal claim record [%d]: contains completed actions", i)
		}

		if err := k.SetClaimRecord(ctx, cr); err != nil {
			return fmt.Errorf("invalid zonedrop proposal claim record [%d]: %w", i, err)
		}
	}

	// process ZoneDrop
	k.SetZoneDrop(ctx, *p.ZoneDrop)

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

func (k Keeper) decompress(data []byte) ([]byte, error) {
	// zip reader
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	return io.ReadAll(zr)
}
