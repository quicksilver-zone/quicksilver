package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
)

type ClaimsManagerKeeper interface {
	GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool)
}
