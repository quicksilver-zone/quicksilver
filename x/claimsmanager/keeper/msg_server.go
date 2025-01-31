package keeper

import (
	"context"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
)

type msgServer struct {
	*Keeper
}

func (m msgServer) SubmitClaim(ctx context.Context, claim *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m msgServer) SubmitClaimableEvent(ctx context.Context, claim *types.MsgSubmitClaimableEventClaim) (*types.MsgSubmitClaimableEventClaimResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: &keeper}
}

var _ types.MsgServer = msgServer{}
