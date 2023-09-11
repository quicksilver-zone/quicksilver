package types // noalias

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	participationrewardstypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the contract needed to be fulfilled for banking and supply
// dependencies.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error
	BlockedAddr(addr sdk.AccAddress) bool
}

// StakingKeeper defines the contract for staking APIs.
type StakingKeeper interface {
	BondDenom(ctx sdk.Context) string
	GetDelegatorBonded(ctx sdk.Context, delegator sdk.AccAddress) sdkmath.Int
}

type GovKeeper interface {
	IterateProposals(ctx sdk.Context, cb func(proposal v1.Proposal) (stop bool))
	GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote v1.Vote, found bool)
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
