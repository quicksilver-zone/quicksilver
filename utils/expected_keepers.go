package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

type ClaimsManagerKeeper interface {
	GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool)
}
