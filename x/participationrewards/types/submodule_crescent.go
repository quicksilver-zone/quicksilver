package types

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquiditytypes "github.com/ingenuity-build/quicksilver/crescent-types/liquidity/types"
	"github.com/ingenuity-build/quicksilver/internal/multierror"
	"github.com/ingenuity-build/quicksilver/utils"
	"strings"
	"time"
)

type CrescentPoolProtocolData struct {
	PoolId         uint64
	ReserveAddress string
	PoolCoinDenom  string
	Disabled       bool
	Denoms         map[string]DenomWithZone
	PoolData       json.RawMessage
	LastUpdated    time.Time
}

func (cpd *CrescentPoolProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if cpd.PoolId == 0 {
		errs["PoolId"] = ErrUndefinedAttribute
	}

	if cpd.ReserveAddress == "" {
		errs["ReserveAddress"] = ErrUndefinedAttribute
	}

	if cpd.PoolCoinDenom == "" {
		errs["PoolCoinDenom"] = ErrUndefinedAttribute
	}

	i := 0
	for _, ibcdenom := range utils.Keys(cpd.Denoms) {
		el := fmt.Sprintf("Denoms[%s]", ibcdenom)

		if cpd.Denoms[ibcdenom].ChainID == "" || len(strings.Split(cpd.Denoms[ibcdenom].ChainID, "-")) < 2 {
			errs[el+" key"] = fmt.Errorf("%w, chainID", ErrInvalidChainID)
		}

		if cpd.Denoms[ibcdenom].Denom == "" || sdk.ValidateDenom(cpd.Denoms[ibcdenom].Denom) != nil {
			errs[el+" value"] = fmt.Errorf("%w, IBC/denom", ErrInvalidDenom)
		}

		i++
	}

	if i == 0 {
		errs["Zones"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (cpd *CrescentPoolProtocolData) GenerateKey() []byte {
	return []byte(fmt.Sprintf("%d", cpd.PoolId))
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
	if crd.Denom == "" {
		errs["ReserveAddress"] = ErrUndefinedAttribute
	}
	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (crd CrescentReserveAddressBalanceProtocolData) GenerateKey() []byte {
	return []byte(crd.ReserveAddress + crd.Denom)
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
	if cpd.PoolCoinDenom == "" {
		return ErrUndefinedAttribute
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
