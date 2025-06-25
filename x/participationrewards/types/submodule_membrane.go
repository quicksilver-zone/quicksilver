package types

import (
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

type MembranePosition struct {
	PositionID       math.Int `json:"position_id"`
	CollateralAssets []struct {
		Asset struct {
			Info struct {
				NativeToken struct {
					Denom string `json:"denom"`
				} `json:"native_token"`
			} `json:"info"`
			Amount math.Int `json:"amount"`
		} `json:"asset"`
		MaxBorrowLTV sdk.Dec `json:"max_borrow_LTV"`
		MaxLTV       sdk.Dec `json:"max_LTV"`
		RateIndex    sdk.Dec `json:"rate_index"`
	} `json:"collateral_assets"`
	CreditAmount math.Int `json:"credit_amount"`
}

// MembraneProtocolData defines the contract address of the membrane contract.
type MembraneProtocolData struct {
	ContractAddress string
}

func (mpd *MembraneProtocolData) ValidateBasic() error {
	if mpd.ContractAddress == "" {
		return errors.New("contract address is required")
	}

	_, err := addressutils.AccAddressFromBech32(mpd.ContractAddress, "")
	return err
}

func (*MembraneProtocolData) GenerateKey() []byte {
	return []byte(MembraneParamsKey)
}
