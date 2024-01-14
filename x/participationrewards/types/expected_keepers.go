package types

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint/types"

	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	interchainquerytypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	interchainstakingtypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862

	SetModuleAccount(sdk.Context, types.ModuleAccountI)
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

// BankKeeper defines the contract needed to be fulfilled for banking and supply
// dependencies.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, senderModule sdk.AccAddress, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
}

// StakingKeeper defines the contract for staking APIs.
type StakingKeeper interface {
	BondDenom(ctx sdk.Context) string
}

type ClaimsManagerKeeper interface {
	ArchiveAndGarbageCollectClaims(ctx sdk.Context, chainID string)
	IterateClaims(ctx sdk.Context, chainID string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	IterateLastEpochUserClaims(ctx sdk.Context, chainID, address string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool)
	SetClaim(ctx sdk.Context, claim *claimsmanagertypes.Claim)
}

type InterchainQueryKeeper interface {
	MakeRequest(
		ctx sdk.Context,
		connectionID,
		chainID,
		queryType string,
		request []byte,
		period sdkmath.Int,
		module string,
		callbackID string,
		ttl uint64,
	)
	GetQuery(ctx sdk.Context, id string) (interchainquerytypes.Query, bool)
}

type InterchainStakingKeeper interface {
	SubmitTx(ctx sdk.Context, msgs []sdk.Msg, account *interchainstakingtypes.ICAAccount, memo string, messagesPerTx int64) error
	IterateZones(ctx sdk.Context, fn func(index int64, zone *interchainstakingtypes.Zone) (stop bool))
	GetZone(ctx sdk.Context, chainID string) (interchainstakingtypes.Zone, bool)
	SetZone(ctx sdk.Context, zone *interchainstakingtypes.Zone)
	AllDelegatorIntents(ctx sdk.Context, zone *interchainstakingtypes.Zone, snapshot bool) []interchainstakingtypes.DelegatorIntent
	SetDelegatorIntent(ctx sdk.Context, zone *interchainstakingtypes.Zone, intent interchainstakingtypes.DelegatorIntent, snapshot bool)
	GetDelegatedAmount(ctx sdk.Context, zone *interchainstakingtypes.Zone) sdk.Coin
	GetDelegationsInProcess(ctx sdk.Context, chainID string) sdkmath.Int
	IterateDelegatorIntents(ctx sdk.Context, zone *interchainstakingtypes.Zone, snapshot bool, fn func(index int64, intent interchainstakingtypes.DelegatorIntent) (stop bool))
	GetValidators(ctx sdk.Context, chainID string) []interchainstakingtypes.Validator
	SetValidator(ctx sdk.Context, chainID string, val interchainstakingtypes.Validator) error
}
