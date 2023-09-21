package types

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/multierror"

	liquiditytypes "github.com/quicksilver-zone/quicksilver/third-party-chains/crescent-types/liquidity/types"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

type CrescentPoolProtocolData struct {
	PoolID      uint64
	Denom       string
	PoolData    json.RawMessage
	LastUpdated time.Time
}

func (cpd *CrescentPoolProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if cpd.PoolID == 0 {
		errs["PoolId"] = ErrUndefinedAttribute
	}

	if err := sdk.ValidateDenom(cpd.Denom); err != nil {
		errs["Denom"] = ErrInvalidDenom
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (cpd *CrescentPoolProtocolData) GenerateKey() []byte {
	return []byte(fmt.Sprintf("%d", cpd.PoolID))
}

func (cpd *CrescentPoolProtocolData) GetPool() (*liquiditytypes.Pool, error) {
	var poolData liquiditytypes.Pool
	if len(cpd.PoolData) > 0 {
		err := json.Unmarshal(cpd.PoolData, &poolData)
		if err != nil {
			return nil, fmt.Errorf("1: unable to unmarshal concrete PoolData: %w", err)
		}
	}
	return &poolData, nil
}

type CrescentReserveAddressBalanceProtocolData struct {
	ReserveAddress string
	Denom          string
	Balance        json.RawMessage
	LastUpdated    time.Time
}

func (crd CrescentReserveAddressBalanceProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if crd.ReserveAddress == "" {
		errs["ReserveAddress"] = ErrUndefinedAttribute
	}
	if _, err := addressutils.AccAddressFromBech32(crd.ReserveAddress, ""); err != nil {
		errs["ReserveAddress"] = ErrInvalidBech32
	}
	if err := sdk.ValidateDenom(crd.Denom); err != nil {
		errs["Denom"] = ErrInvalidDenom
	}
	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (crd CrescentReserveAddressBalanceProtocolData) GenerateKey() []byte {
	return []byte(fmt.Sprintf("%s_%s", crd.ReserveAddress, crd.Denom))
}

func (crd CrescentReserveAddressBalanceProtocolData) GetBalance() (math.Int, error) {
	var balanceData math.Int
	err := json.Unmarshal(crd.Balance, &balanceData)
	if err != nil {
		return balanceData, fmt.Errorf("1: unable to unmarshal concrete reservebalancedata: %w", err)
	}
	return balanceData, nil
}

type CrescentPoolCoinSupplyProtocolData struct {
	PoolCoinDenom string
	Supply        json.RawMessage
	LastUpdated   time.Time
}

func (cpd CrescentPoolCoinSupplyProtocolData) ValidateBasic() error {
	// poolcoindenom is always pool{poolid}
	if len(cpd.PoolCoinDenom) <= 4 || cpd.PoolCoinDenom[0:4] != "pool" {
		return ErrInvalidAssetName
	}

	return nil
}

func (cpd CrescentPoolCoinSupplyProtocolData) GenerateKey() []byte {
	return []byte(cpd.PoolCoinDenom)
}

func (cpd CrescentPoolCoinSupplyProtocolData) GetSupply() (math.Int, error) {
	var supplyData math.Int
	err := json.Unmarshal(cpd.Supply, &supplyData)
	if err != nil {
		return supplyData, fmt.Errorf("1: unable to unmarshal concrete supplydata: %w", err)
	}
	return supplyData, nil
}

type CrescentParamsProtocolData struct {
	ChainID string
}

func (uppd CrescentParamsProtocolData) ValidateBasic() error {
	if uppd.ChainID == "" {
		return ErrUndefinedAttribute
	}

	return nil
}

func (uppd CrescentParamsProtocolData) GenerateKey() []byte {
	return []byte(CrescentParamsKey)
}
