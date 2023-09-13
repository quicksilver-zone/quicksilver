package stableswap

import (
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
