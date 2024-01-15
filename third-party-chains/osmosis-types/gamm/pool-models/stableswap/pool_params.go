package stableswap

import (
	"github.com/quicksilver-zone/quicksilver/v7/third-party-chains/osmosis-types/gamm"
)

func (params PoolParams) Validate() error {
	if params.ExitFee.IsNegative() {
		return gamm.ErrNegativeExitFee
	}

	if params.ExitFee.GTE(sdkmath.LegacyOneDec()) {
		return gamm.ErrTooMuchExitFee
	}

	if params.SwapFee.IsNegative() {
		return gamm.ErrNegativeSwapFee
	}

	if params.SwapFee.GTE(sdkmath.LegacyOneDec()) {
		return gamm.ErrTooMuchSwapFee
	}
	return nil
}
