package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	prtypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"

	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v5/modules/core/exported"
)

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	GetConnection(ctx sdk.Context, connectionID string) (ibcexported.ConnectionI, error)
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
	IsBound(ctx sdk.Context, portID string) bool
}

type ParticipationRewardsKeeper interface {
	IterateUserClaims(ctx sdk.Context, chainID string, address string, fn func(index int64, data prtypes.Claim) (stop bool))
	IterateLastEpochUserClaims(ctx sdk.Context, chainID string, address string, fn func(index int64, data prtypes.Claim) (stop bool))
}
