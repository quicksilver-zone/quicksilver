package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/v7/x/claimsmanager/types"
	epochstypes "github.com/quicksilver-zone/quicksilver/v7/x/epochs/types"
)

// ChannelKeeper defines the expected IBC channel keeper.
type ChannelKeeper interface {
	GetChannel(ctx context.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx context.Context, portID, channelID string) (uint64, bool)
	GetConnection(ctx context.Context, connectionID string) (ibcexported.ConnectionI, error)
}

// PortKeeper defines the expected IBC port keeper.
type PortKeeper interface {
	BindPort(ctx context.Context, portID string) *capabilitytypes.Capability
	IsBound(ctx context.Context, portID string) bool
}

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx context.Context, denom string) sdk.Coin
	HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

type IcsHooks interface {
	AfterZoneCreated(ctx context.Context, connectionID, chainID, accountPrefix string) error
}

type ClaimsManagerKeeper interface {
	IterateLastEpochUserClaims(ctx sdk.Context, chainID, address string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	SetClaim(ctx sdk.Context, claim *claimsmanagertypes.Claim)
}

type EpochsKeeper interface {
	GetEpochInfo(ctx sdk.Context, identifier string) epochstypes.EpochInfo
}
