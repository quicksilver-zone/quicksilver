package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

func (k *Keeper) HandleAcknowledgement(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, connectionID string) error {
	return nil
}

func (*Keeper) HandleTimeout(_ sdk.Context, _ channeltypes.Packet) error {
	return nil
}

// ----------------------------------------------------------------

func (k *Keeper) HandleMsgTransfer(ctx sdk.Context, msg ibctransfertypes.FungibleTokenPacketData, ibcDenom string) error {
	return nil
}
