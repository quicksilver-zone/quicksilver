package utils

import (
	"errors"
	"fmt"
	"net/url"

	"cosmossdk.io/api/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	tmclienttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
)

type ProofOpsFn func(ctx sdk.Context, ibcKeeper *ibckeeper.Keeper, connectionID, chainID string, height int64, module string, key []byte, data []byte, proofOps *crypto.ProofOps) error

type SelfProofOpsFn func(ctx sdk.Context, claimsKeeper claimsmanagerkeeper.Keeper, consensusStateKey, module string, key []byte, data []byte, proofOps *crypto.ProofOps) error

func ValidateProofOps(
	ctx sdk.Context,
	ibcKeeper *ibckeeper.Keeper,
	connectionID,
	chainID string,
	height int64,
	module string,
	key,
	data []byte,
	proofOps *crypto.ProofOps,
) error {
	if proofOps == nil {
		return errors.New("unable to validate proof. No proof submitted")
	}
	connection, _ := ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)

	csHeight := clienttypes.NewHeight(clienttypes.ParseChainID(chainID), uint64(height)+1)
	consensusState, found := ibcKeeper.ClientKeeper.GetClientConsensusState(ctx, connection.ClientId, csHeight)

	if !found {
		return errors.New("unable to fetch consensus state")
	}

	clientState, found := ibcKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return errors.New("unable to fetch client state")
	}

	path := commitmenttypes.NewMerklePath([]string{module, url.PathEscape(string(key))}...)

	merkleProof, err := commitmenttypes.ConvertProofs(proofOps)
	if err != nil {
		return errors.New("error converting proofs")
	}

	tmClientState, ok := clientState.(*tmclienttypes.ClientState)
	if !ok {
		return errors.New("error unmarshaling client state")
	}

	if len(data) != 0 {
		// if we got a non-nil response, verify inclusion proof.
		if err := merkleProof.VerifyMembership(tmClientState.ProofSpecs, consensusState.GetRoot(), path, data); err != nil {
			return fmt.Errorf("unable to verify inclusion proof: %w", err)
		}
		return nil

	}
	// if we got a nil response, verify non inclusion proof.
	if err := merkleProof.VerifyNonMembership(tmClientState.ProofSpecs, consensusState.GetRoot(), path); err != nil {
		return fmt.Errorf("unable to verify non-inclusion proof: %w", err)
	}
	return nil
}

func ValidateSelfProofOps(ctx sdk.Context, claimsKeeper claimsmanagerkeeper.Keeper, consensusStateKey, module string, key, data []byte, proofOps *crypto.ProofOps) error {
	if proofOps == nil {
		return errors.New("unable to validate proof. No proof submitted")
	}

	consensusState, found := claimsKeeper.GetSelfConsensusState(ctx, consensusStateKey)
	if !found {
		return errors.New("unable to lookup self-consensus state")
	}

	proofSpecs := commitmenttypes.GetSDKSpecs()

	path := commitmenttypes.NewMerklePath([]string{module, url.PathEscape(string(key))}...)

	merkleProof, err := commitmenttypes.ConvertProofs(proofOps)
	if err != nil {
		return errors.New("error converting proofs")
	}

	if len(data) != 0 {
		// if we got a non-nil response, verify inclusion proof.
		if err := merkleProof.VerifyMembership(proofSpecs, consensusState.GetRoot(), path, data); err != nil {
			return fmt.Errorf("unable to verify inclusion proof: %w", err)
		}
		return nil

	}
	// if we got a nil response, verify non inclusion proof.
	if err := merkleProof.VerifyNonMembership(proofSpecs, consensusState.GetRoot(), path); err != nil {
		return fmt.Errorf("unable to verify non-inclusion proof: %w", err)
	}
	return nil
}

func MockSelfProofOps(_ sdk.Context, _ claimsmanagerkeeper.Keeper, _, _ string, _, _ []byte, _ *crypto.ProofOps) error {
	return nil
}

func MockProofOps(_ sdk.Context, _ *ibckeeper.Keeper, _, _ string, _ int64, _ string, _, _ []byte, _ *crypto.ProofOps) error {
	return nil
}
