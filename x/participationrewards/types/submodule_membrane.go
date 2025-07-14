package types

import (
	"errors"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

// the following types are required to decode the membrane position from the claim.
// silly revive lint rule requires them to not be nested.
type MembranePosition struct {
	PositionID       math.Int                  `json:"position_id"`
	CollateralAssets []MembraneCollateralAsset `json:"collateral_assets"`
	CreditAmount     math.Int                  `json:"credit_amount"`
}

type MembraneCollateralAsset struct {
	Asset        MembraneCollateralAssetAsset `json:"asset"`
	MaxBorrowLTV sdk.Dec                      `json:"max_borrow_LTV"`
	MaxLTV       sdk.Dec                      `json:"max_LTV"`
	RateIndex    sdk.Dec                      `json:"rate_index"`
}

type MembraneCollateralAssetAsset struct {
	Info   MembraneCollateralAssetInfo `json:"info"`
	Amount math.Int                    `json:"amount"`
}

type MembraneCollateralAssetInfo struct {
	NativeToken MembraneNativeToken `json:"native_token"`
}

type MembraneNativeToken struct {
	Denom string `json:"denom"`
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
