package stableswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
)

func (params PoolParams) Validate() error {
	if params.ExitFee.IsNegative() {
		return gamm.ErrNegativeExitFee
	}

	if params.ExitFee.GTE(sdk.OneDec()) {
		return gamm.ErrTooMuchExitFee
	}

	if params.SwapFee.IsNegative() {
		return gamm.ErrNegativeSwapFee
	}

	if params.SwapFee.GTE(sdk.OneDec()) {
		return gamm.ErrTooMuchSwapFee
	}
	return nil
}
