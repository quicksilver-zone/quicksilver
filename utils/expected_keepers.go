package utils

import (
	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
)

type ClaimsManagerKeeper interface {
	IterateLastEpochUserClaims(ctx sdk.Context, chainID, address string, fn func(index int64, data claimsmanagertypes.Claim) (stop bool))
	GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool)
}
