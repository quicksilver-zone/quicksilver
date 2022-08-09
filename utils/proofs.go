package utils

import (
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v3/modules/core/23-commitment/types"
	ibcKeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	tmclienttypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
)

func ValidateProofOps(ctx sdk.Context, ibcKeeper *ibcKeeper.Keeper, connectionID string, chainID string, height int64, module string, key []byte, data []byte, proofOps *crypto.ProofOps) error {
	if proofOps == nil {
		return fmt.Errorf("unable to validate proof. No proof submitted")
	}
	connection, _ := ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)

	csHeight := clienttypes.NewHeight(clienttypes.ParseChainID(chainID), uint64(height)+1)
	consensusState, found := ibcKeeper.ClientKeeper.GetClientConsensusState(ctx, connection.ClientId, csHeight)

	if !found {
		return fmt.Errorf("unable to fetch consensus state")
	}

	clientState, found := ibcKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return fmt.Errorf("unable to fetch client state")
	}

	path := commitmenttypes.NewMerklePath([]string{module, url.PathEscape(string(key))}...)

	merkleProof, err := commitmenttypes.ConvertProofs(proofOps)
	if err != nil {
		return fmt.Errorf("error converting proofs")
	}

	tmClientState, ok := clientState.(*tmclienttypes.ClientState)
	if !ok {
		return fmt.Errorf("error unmarshaling client state")
	}

	if len(data) != 0 {
		// if we got a non-nil response, verify inclusion proof.
		if err := merkleProof.VerifyMembership(tmClientState.ProofSpecs, consensusState.GetRoot(), path, data); err != nil {
			return fmt.Errorf("unable to verify inclusion proof: %s", err)
		}
		return nil

	}
	// if we got a nil response, verify non inclusion proof.
	if err := merkleProof.VerifyNonMembership(tmClientState.ProofSpecs, consensusState.GetRoot(), path); err != nil {
		return fmt.Errorf("unable to verify non-inclusion proof: %s", err)
	}
	return nil
}
