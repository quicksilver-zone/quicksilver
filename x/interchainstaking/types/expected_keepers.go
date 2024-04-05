package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v5/modules/core/exported"

	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

// ChannelKeeper defines the expected IBC channel keeper.
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	GetConnection(ctx sdk.Context, connectionID string) (ibcexported.ConnectionI, error)
}

// PortKeeper defines the expected IBC port keeper.
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
	IsBound(ctx sdk.Context, portID string) bool
}

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, senderAddr sdk.AccAddress, recipientADdr sdk.AccAddress, amt sdk.Coins) error
}

type IcsHooks interface {
	AfterZoneCreated(ctx sdk.Context, zone *Zone) error
}

type ClaimsManagerKeeper interface {
	IterateLastEpochClaims(ctx sdk.Context, chainID string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	IterateLastEpochUserClaims(ctx sdk.Context, chainID, address string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	IterateUserClaims(ctx sdk.Context, chainID, address string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	SetClaim(ctx sdk.Context, claim *claimsmanagertypes.Claim)
}

type EpochsKeeper interface {
	GetEpochInfo(ctx sdk.Context, identifier string) epochstypes.EpochInfo
}

type EventManagerKeeper interface {
	AddEvent(ctx sdk.Context, module, chainID, identifier, callback string, eventType, status int32, condtion emtypes.ConditionI, payload []byte)
}
