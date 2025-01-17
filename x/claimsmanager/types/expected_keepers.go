package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	interchainquerytypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	interchainstakingtypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

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
	GetLocalAddressMap(ctx sdk.Context, remoteAddress sdk.AccAddress, chainID string) (sdk.AccAddress, bool)
}
