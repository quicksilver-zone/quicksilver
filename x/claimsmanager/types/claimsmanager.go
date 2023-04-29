package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

// ValidateBasic performs stateless validation of a Claim.
func (c *Claim) ValidateBasic() error {
	errs := make(map[string]error)

	_, err := sdk.AccAddressFromBech32(c.UserAddress)
	if err != nil {
		errs["UserAddress"] = err
	}

	if c.ChainId == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	if c.Amount <= 0 {
		errs["Amount"] = ErrNotPositive
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

// UserAllocation is an internal keeper struct to track transient state for
// rewards distribution. It contains the user address and the coins that are
// allocated to it.
type UserAllocation struct {
	Address string
	Amount  math.Int
}

type CustomeZone struct {
	ChainId            string
	HoldingsAllocation uint64
	LocalDenom         string
}
