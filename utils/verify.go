package utils

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	tmmath "github.com/tendermint/tendermint/libs/math"
	"github.com/tendermint/tendermint/light"
	"github.com/tendermint/tendermint/types"
)

// VerifyNonAdjacent is identical to VerifyNonAdjacent in tendermint/tendermint/light/verifier.go, with the exception that
// it does not attempt to validate that the block is _newer_ than the current consensus state.
func VerifyNonAdjacent(
	trustedHeader *types.SignedHeader, // height=X
	trustedVals *types.ValidatorSet, // height=X or height=X+1
	untrustedHeader *types.SignedHeader, // height=Y
	untrustedVals *types.ValidatorSet, // height=Y
	trustingPeriod time.Duration,
	now time.Time,
	_ time.Duration,
	trustLevel tmmath.Fraction,
) error {
	if untrustedHeader.Height == trustedHeader.Height+1 {
		return errors.New("headers must be non adjacent in height")
	}

	if light.HeaderExpired(trustedHeader, trustingPeriod, now) {
		return light.ErrOldHeaderExpired{At: trustedHeader.Time.Add(trustingPeriod), Now: now}
	}

	// if err := verifyNewHeaderAndVals(
	// 	untrustedHeader, untrustedVals,
	// 	trustedHeader,
	// 	now, maxClockDrift); err != nil {
	// 	return ErrInvalidHeader{err}
	// }

	// Ensure that +`trustLevel` (default 1/3) or more of last trusted validators signed correctly.
	err := trustedVals.VerifyCommitLightTrusting(trustedHeader.ChainID, untrustedHeader.Commit, trustLevel)
	if err != nil {
		var typeError *types.ErrNotEnoughVotingPowerSigned
		if errors.As(err, &typeError) {
			return light.ErrNewValSetCantBeTrusted{Reason: *typeError}
		}

		return err
	}

	// Ensure that +2/3 of new validators signed correctly.
	//
	// NOTE: this should always be the last check because untrustedVals can be
	// intentionally made very large to DOS the light client. not the case for
	// VerifyAdjacent, where validator set is known in advance.
	if err := untrustedVals.VerifyCommitLight(trustedHeader.ChainID, untrustedHeader.Commit.BlockID,
		untrustedHeader.Height, untrustedHeader.Commit); err != nil {
		return light.ErrInvalidHeader{Reason: err}
	}

	return nil
}

func VerifyAdjacent(
	trustedHeader *types.SignedHeader, // height=X
	untrustedHeader *types.SignedHeader, // height=X+1
	untrustedVals *types.ValidatorSet, // height=X+1
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration, //nolint:revive
) error {
	if untrustedHeader.Height != trustedHeader.Height+1 {
		return errors.New("headers must be adjacent in height")
	}

	if light.HeaderExpired(trustedHeader, trustingPeriod, now) {
		return light.ErrOldHeaderExpired{At: trustedHeader.Time.Add(trustingPeriod), Now: now}
	}

	// if err := verifyNewHeaderAndVals(
	// 	untrustedHeader, untrustedVals,
	// 	trustedHeader,
	// 	now, maxClockDrift); err != nil {
	// 	return ErrInvalidHeader{err}
	// }

	// Check the validator hashes are the same
	if !bytes.Equal(untrustedHeader.ValidatorsHash, trustedHeader.NextValidatorsHash) {
		err := fmt.Errorf("expected old header next validators (%X) to match those from new header (%X)",
			trustedHeader.NextValidatorsHash,
			untrustedHeader.ValidatorsHash,
		)
		return err
	}

	// Ensure that +2/3 of new validators signed correctly.
	if err := untrustedVals.VerifyCommitLight(trustedHeader.ChainID, untrustedHeader.Commit.BlockID,
		untrustedHeader.Height, untrustedHeader.Commit); err != nil {
		return light.ErrInvalidHeader{Reason: err}
	}

	return nil
}

// Verify combines both VerifyAdjacent and VerifyNonAdjacent functions.
func Verify(
	trustedHeader *types.SignedHeader, // height=X
	trustedVals *types.ValidatorSet, // height=X or height=X+1
	untrustedHeader *types.SignedHeader, // height=Y
	untrustedVals *types.ValidatorSet, // height=Y
	trustingPeriod time.Duration,
	now time.Time,
	maxClockDrift time.Duration,
	trustLevel tmmath.Fraction,
) error {
	if untrustedHeader.Height != trustedHeader.Height+1 {
		return VerifyNonAdjacent(trustedHeader, trustedVals, untrustedHeader, untrustedVals,
			trustingPeriod, now, maxClockDrift, trustLevel)
	}

	return VerifyAdjacent(trustedHeader, untrustedHeader, untrustedVals, trustingPeriod, now, maxClockDrift)
}
