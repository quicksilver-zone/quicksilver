package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	interchainquerytypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
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
