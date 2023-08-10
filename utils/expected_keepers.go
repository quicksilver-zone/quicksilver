package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	claimsmanagertypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
)

type ClaimsManagerKeeper interface {
	IterateLastEpochUserClaims(ctx sdk.Context, chainID, address string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool)
}
