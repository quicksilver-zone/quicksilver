package types

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	errNilEpochProvisions      = errors.New("epoch provisions was nil in genesis")
	errNegativeEpochProvisions = errors.New("epoch provisions should be non-negative")
)

// NewMinter returns a new Minter object with the given epoch
// provisions values.
func NewMinter(epochProvisions sdkmath.LegacyDec) Minter {
	return Minter{
		EpochProvisions: epochProvisions,
	}
}

// InitialMinter returns an initial Minter object.
func InitialMinter() Minter {
	return NewMinter(sdkmath.LegacyNewDec(0))
}

// DefaultInitialMinter returns a default initial Minter object for a new chain.
func DefaultInitialMinter() Minter {
	return InitialMinter()
}

// Validate checks if minter epoch provisions are valid.
func (m Minter) Validate() error {
	if m.EpochProvisions.IsNil() {
		return errNilEpochProvisions
	}

	if m.EpochProvisions.IsNegative() {
		return errNegativeEpochProvisions
	}
	return nil
}

// NextEpochProvisions returns the epoch provisions.
func (m Minter) NextEpochProvisions(params Params) sdkmath.LegacyDec {
	return m.EpochProvisions.Mul(params.ReductionFactor)
}

// EpochProvision returns the provisions for a block based on the epoch
// provisions rate.
func (m Minter) EpochProvision(params Params) sdk.Coin {
	provisionAmt := m.EpochProvisions
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
