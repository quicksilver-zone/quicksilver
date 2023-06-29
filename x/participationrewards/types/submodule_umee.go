package types

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/internal/multierror"
	"time"
)

type UmeeProtocolData struct {
	Denom       string
	LastUpdated time.Time
	Data        json.RawMessage
}

func (upd *UmeeProtocolData) ValidateBasic() error {
	if upd.Denom == "" {
		return ErrUndefinedAttribute
	}
	return nil
}

func (upd *UmeeProtocolData) GenerateKey() []byte {
	return []byte(upd.Denom)
}

func GetUnderlyingData[V math.Int | sdk.Dec](upd *UmeeProtocolData) (V, error) {
	var data V
	err := json.Unmarshal(upd.Data, &data)
	if err != nil {
		return data, fmt.Errorf("1: unable to unmarshal concrete reservedata: %w", err)
	}
	return data, nil
}

// UmeeReservesProtocolData defines protocol state to track qAssets in
// Umee reserves.
type UmeeReservesProtocolData struct {
	UmeeProtocolData
}

func (upd *UmeeReservesProtocolData) GetReserveAmount() (math.Int, error) {
	return GetUnderlyingData[math.Int](&upd.UmeeProtocolData)
}

type UmeeTotalBorrowsProtocolData struct {
	UmeeProtocolData
}

func (upd *UmeeTotalBorrowsProtocolData) GetTotalBorrows() (sdk.Dec, error) {
	return GetUnderlyingData[sdk.Dec](&upd.UmeeProtocolData)
}

type UmeeInterestScalarProtocolData struct {
	UmeeProtocolData
}

func (upd *UmeeInterestScalarProtocolData) GetInterestScalar() (sdk.Dec, error) {
	return GetUnderlyingData[sdk.Dec](&upd.UmeeProtocolData)
}

type UmeeUTokenSupplyProtocolData struct {
	UmeeProtocolData
}

func (upd *UmeeUTokenSupplyProtocolData) GetUTokenSupply() (math.Int, error) {
	return GetUnderlyingData[math.Int](&upd.UmeeProtocolData)
}

type UmeeLeverageModuleBalanceProtocolData struct {
	UmeeProtocolData
}

func (upd *UmeeLeverageModuleBalanceProtocolData) GetModuleBalance() (math.Int, error) {
	return GetUnderlyingData[math.Int](&upd.UmeeProtocolData)
}

// -----------------------------------------------------

type UmeeParamsProtocolData struct {
	ChainID   string
	BaseDenom string
	BaseChain string
}

func (uppd UmeeParamsProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if uppd.ChainID == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	if uppd.BaseChain == "" {
		errs["BaseChain"] = ErrUndefinedAttribute
	}

	if uppd.BaseDenom == "" {
		errs["BaseDenom"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		return multierror.New(errs)
	}

	return nil
}

func (uppd UmeeParamsProtocolData) GenerateKey() []byte {
	return []byte(UmeeParamsKey)
}
