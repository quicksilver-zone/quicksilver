package types

import (
	"errors"
	"fmt"
)

var ErrInvalidSubzoneForBasezone = errors.New("invalid subzone for basezone")

// IsSubzone returns true if this zone is a sub-zone.
func (z *Zone) IsSubzone() bool {
	return z.SubzoneInfo != nil
}

// ChainID returns the ID of the running chain for the given zone.
func (z *Zone) ChainID() string {
	if z.IsSubzone() {
		return z.SubzoneInfo.BaseChainID
	}

	return z.ChainId
}

// ID returns the unique identifier for the given zone.
func (z *Zone) ID() string {
	return z.ChainId
}

func ValidateSubzoneForBasezone(subZone, baseZone Zone) error {
	if baseZone.IsSubzone() {
		return fmt.Errorf("cannot make a subzone for a subzone: %w", ErrInvalidSubzoneForBasezone)
	}

	if !subZone.IsSubzone() ||
		subZone.SubzoneInfo.Authority == "" ||
		subZone.SubzoneInfo.BaseChainID == "" ||
		subZone.SubzoneInfo.ChainID == "" {
		return fmt.Errorf("all subzone info must be populated: %w", ErrInvalidSubzoneForBasezone)
	}

	if subZone.ConnectionId != baseZone.ConnectionId {
		return fmt.Errorf("connection IDs must be identical for subzone and basezone: %w", ErrInvalidSubzoneForBasezone)
	}

	if subZone.SubzoneInfo.BaseChainID != baseZone.ChainId {
		return fmt.Errorf("incorrect basechainID for subzone: %w", ErrInvalidSubzoneForBasezone)
	}

	if subZone.SubzoneInfo.ChainID == baseZone.ChainId {
		return fmt.Errorf("subzone chainID must be unique: %w", ErrInvalidSubzoneForBasezone)
	}

	if subZone.AccountPrefix != baseZone.AccountPrefix {
		return fmt.Errorf("account prefix mismatch: %w", ErrInvalidSubzoneForBasezone)
	}

	if subZone.BaseDenom != baseZone.BaseDenom {
		return fmt.Errorf("base denom mismatch: %w", ErrInvalidSubzoneForBasezone)
	}

	if subZone.LocalDenom == baseZone.LocalDenom {
		return fmt.Errorf("subzone local denom must be unique: %w", ErrInvalidSubzoneForBasezone)
	}

	if subZone.LiquidityModule != baseZone.LiquidityModule ||
		subZone.MultiSend != baseZone.MultiSend ||
		subZone.Decimals != baseZone.Decimals ||
		subZone.Is_118 != baseZone.Is_118 {
		return fmt.Errorf("chain capability mismatch: %w", ErrInvalidSubzoneForBasezone)
	}

	return nil
}
