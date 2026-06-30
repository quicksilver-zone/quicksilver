package keeper

import (
	"context"
	"errors"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
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
	return nil, errors.New("unsupported action")}

// MsgGovRemoveProtocolData removes a protocoldata item.
func (k msgServer) GovRemoveProtocolData(goCtx context.Context, msg *types.MsgGovRemoveProtocolData) (*types.MsgGovRemoveProtocolDataResponse, error) {
	return &types.MsgGovRemoveProtocolDataResponse{}, nil
}
