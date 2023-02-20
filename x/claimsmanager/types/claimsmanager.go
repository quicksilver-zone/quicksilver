package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

func (c Claim) ValidateBasic() error {
	errors := make(map[string]error)

	_, err := sdk.AccAddressFromBech32(c.UserAddress)
	if err != nil {
		errors["UserAddress"] = err
	}

	if len(c.ChainId) == 0 {
		errors["ChainId"] = ErrUndefinedAttribute
	}

	if c.Amount <= 0 {
		errors["Amount"] = ErrNotPositive
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
