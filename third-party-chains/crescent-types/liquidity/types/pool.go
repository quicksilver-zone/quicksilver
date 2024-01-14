package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/utils/addressutils"
)

func (pool Pool) GetReserveAddress() string {
	return pool.ReserveAddress
}

// Validate validates Pool for genesis.
func (pool Pool) Validate() error {
	if pool.Id == 0 {
		return fmt.Errorf("pool id must not be 0")
	}
	if pool.PairId == 0 {
		return fmt.Errorf("pair id must not be 0")
	}
	if _, err := addressutils.AccAddressFromBech32(pool.ReserveAddress, ""); err != nil {
		return fmt.Errorf("invalid reserve address %s: %w", pool.ReserveAddress, err)
	}
	if err := sdk.ValidateDenom(pool.PoolCoinDenom); err != nil {
		return fmt.Errorf("invalid pool coin denom: %w", err)
	}
	return nil
}
