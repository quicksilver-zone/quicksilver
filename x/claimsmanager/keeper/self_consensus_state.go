package keeper

import (
	"errors"
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/cosmos/ibc-go/v5/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
)

// GetSelfConsensusState returns consensus state stored every epoch
func (k Keeper) GetSelfConsensusState(ctx sdk.Context, key string) (ibctmtypes.ConsensusState, bool) {
	store := ctx.KVStore(k.storeKey)

	var selfConsensusState ibctmtypes.ConsensusState
	k.cdc.MustUnmarshal(store.Get(append(types.KeySelfConsensusState, []byte(key)...)), &selfConsensusState)

	return selfConsensusState, true
}

// SetSelfConsensusState sets the self consensus state
func (k Keeper) SetSelfConsensusState(ctx sdk.Context, key string, consState ibctmtypes.ConsensusState) {
	store := ctx.KVStore(k.storeKey)
	store.Set(store.Get(append(types.KeySelfConsensusState, []byte(key)...)), k.cdc.MustMarshal(&consState))
}

// ValidateSelfProofOps Validate Proof Ops against a consensus state stored in the claimsmanager Keeper. ConsensusStateKey is the key to lookup the state.
func (k Keeper) ValidateSelfProofOps(ctx sdk.Context, consensusStateKey string, chainID string, height int64, module string, key []byte, data []byte, proofOps *crypto.ProofOps) error {
	if proofOps == nil {
		return errors.New("unable to validate proof. No proof submitted")
	}

	consensusState, found := k.GetSelfConsensusState(ctx, consensusStateKey)
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
