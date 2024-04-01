package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"

)

type ClaimsManagerKeeper interface {
	GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool)
}
