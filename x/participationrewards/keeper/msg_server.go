package keeper

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// SubmitClaim is used to verify, by proof, that the given user address has a
// claim against the given asset type for the given zone.
func (k msgServer) SubmitClaim(goCtx context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.GetClaimsEnabled(ctx) {
		return nil, errors.New("claims currently disabled")
	}
	// fetch zone
	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	if !ok {
		return nil, fmt.Errorf("invalid zone, chain id \"%s\" not found", msg.Zone)
	}
	var pd types.ProtocolData
	pd, ok = k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, msg.SrcZone)
	if !ok {
		return nil, fmt.Errorf("unable to obtain connection protocol data for %q", msg.SrcZone)
	}

	// protocol data
	iConnectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, pd.Data)
	if err != nil {
		k.Logger(ctx).Error("SubmitClaim: error unmarshalling protocol data")
	}
	connectionData, _ := iConnectionData.(*types.ConnectionProtocolData)

	for i, proof := range msg.Proofs {
		pl := fmt.Sprintf("Proof [%d]", i)

		if proof.Height != connectionData.LastEpoch {
			return nil, fmt.Errorf(
				"invalid claim for last epoch, %s expected height %d, got %d",
				pl,
				connectionData.LastEpoch,
				proof.Height,
			)
		}

		// if we are claiming against Quicksilver, use the SelfProofOpsFn.
		if msg.SrcZone == ctx.ChainID() {
			if err := k.ValidateSelfProofOps(
				ctx,
				k.icsKeeper.ClaimsManagerKeeper,
				"epoch",
				proof.ProofType,
				proof.Key,
				proof.Data,
				proof.ProofOps,
			); err != nil {
				return nil, fmt.Errorf("%s: %w", pl, err)
			}
		} else {
			if err := k.ValidateProofOps(
				ctx,
				&k.icsKeeper.IBCKeeper,
				connectionData.ConnectionID,
				connectionData.ChainID,
				proof.Height,
				proof.ProofType,
				proof.Key,
				proof.Data,
				proof.ProofOps,
			); err != nil {
				return nil, fmt.Errorf("%s: %w", pl, err)
			}
		}
	}

	// if we get here all data was validated; verifyClaim will write the claim to the correct store.
	if mod, ok := k.prSubmodules[msg.ClaimType]; ok {
		// vertifyClaim needs to return the amount!
		amount, err := mod.ValidateClaim(ctx, k.Keeper, msg)
		if err != nil {
			return nil, fmt.Errorf("claim validation failed: %w", err)
		}
		claim := k.icsKeeper.ClaimsManagerKeeper.NewClaim(msg.UserAddress, zone.ChainId, msg.ClaimType, msg.SrcZone, amount)
		k.icsKeeper.ClaimsManagerKeeper.SetClaim(ctx, &claim)
	}

	return &types.MsgSubmitClaimResponse{}, nil
}

// MsgGovRemoveProtocolData removes a protocoldata item.
func (k msgServer) GovRemoveProtocolData(goCtx context.Context, msg *types.MsgGovRemoveProtocolData) (*types.MsgGovRemoveProtocolDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checking msg authority is the gov module address
	if k.Keeper.GetGovAuthority(ctx) != msg.Authority {
		return &types.MsgGovRemoveProtocolDataResponse{},
			govtypes.ErrInvalidSigner.Wrapf(
				"invalid authority: expected %s, got %s",
				k.Keeper.GetGovAuthority(ctx), msg.Authority,
			)
	}

	k.Keeper.DeleteProtocolData(ctx, []byte(msg.Key))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeDeleteKeyProposal,
			sdk.NewAttribute(types.AttributeKeyProtocolDataKey, msg.Key),
		),
	})

	return &types.MsgGovRemoveProtocolDataResponse{}, nil
}
