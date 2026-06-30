package keeper

import (
	"context"
	"errors"

	"github.com/quicksilver-zone/quicksilver/x/supply/types"
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

func (k msgServer) IncentivePoolSpend(goCtx context.Context, msg *types.MsgIncentivePoolSpend) (*types.MsgIncentivePoolSpendResponse, error) {
	return nil, errors.New("incentive pool spend is currently disabled")
}
