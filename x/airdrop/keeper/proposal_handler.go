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

var chunkSize = int64(4096)

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

		if cr.ActionsCompleted != nil {
			for j, ca := range cr.ActionsCompleted {
				if ctx.BlockTime().Before(ca.CompleteTime) {
					return fmt.Errorf("invalid zonedrop proposal claim record [%d]: completed action [%d] is in future", i, j)
				}
			}
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
	// input buffer
	var ibuf bytes.Buffer
	ibuf.Write(data)

	// zip reader
	zr, err := zlib.NewReader(&ibuf)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	// copy reader data to output buffer (writer)
	// - prevents data going out of scope;
	var obuf bytes.Buffer
	for {
		_, err := io.CopyN(&obuf, zr, chunkSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return obuf.Bytes(), nil
}
