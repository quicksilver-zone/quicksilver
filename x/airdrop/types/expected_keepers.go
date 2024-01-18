package types // noalias

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icstypes "github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
	participationrewardstypes "github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/types"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the contract needed to be fulfilled for banking and supply
// dependencies.
type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	IsSendEnabledCoins(ctx context.Context, coins ...sdk.Coin) error
	BlockedAddr(addr sdk.AccAddress) bool
}

// StakingKeeper defines the contract for staking APIs.
type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
	GetDelegatorBonded(ctx context.Context, delegator sdk.AccAddress) (sdkmath.Int, error)
}

type GovKeeper interface {
	// Proposals(ctx context.Context, req *v1.QueryProposalsRequest) (*v1.QueryProposalsResponse, error)
	// Vote(ctx context.Context, req *v1.QueryVoteRequest) (*v1.QueryVoteResponse, error)
}

type InterchainStakingKeeper interface {
	GetZone(ctx sdk.Context, chainID string) (icstypes.Zone, bool)
	GetDelegatorIntent(ctx sdk.Context, zone *icstypes.Zone, delegator string, snapshot bool) (icstypes.DelegatorIntent, bool)
	IterateZones(ctx sdk.Context, fn func(index int64, zone *icstypes.Zone) (stop bool))
	UserZoneReceipts(ctx sdk.Context, zone *icstypes.Zone, addr sdk.AccAddress) ([]icstypes.Receipt, error)
}

type ParticipationRewardsKeeper interface {
	GetProtocolData(ctx sdk.Context, pdType participationrewardstypes.ProtocolDataType, key string) (participationrewardstypes.ProtocolData, bool)
}
