package types

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	liquiditytypes "github.com/ingenuity-build/quicksilver/crescent-types/liquidity/types"
	"time"
)

type PoolType int32

const (
	// POOL_TYPE_UNSPECIFIED specifies unknown pool type
	PoolTypeUnspecified PoolType = 0
	// POOL_TYPE_BASIC specifies the basic pool type
	PoolTypeBasic PoolType = 1
	// POOL_TYPE_RANGED specifies the ranged pool type
	PoolTypeRanged PoolType = 2
)

type CrescentPoolProtocolData struct {
	Type           PoolType
	PoolId         uint64
	PairId         uint64
	ReserveAddress string
	PoolCoinDenom  string
	Disabled       bool
	Denoms         map[string]DenomWithZone
	PoolData       json.RawMessage
	LastUpdated    time.Time
}

func (cpd *CrescentPoolProtocolData) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
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

type CrescentPairProtocolData struct {
	PairId      uint64
	PairData    json.RawMessage
	LastUpdated time.Time
}

func (cpd CrescentPairProtocolData) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (cpd CrescentPairProtocolData) GenerateKey() []byte {
	return []byte(fmt.Sprintf("%d", cpd.PairId))
}

func (cpd CrescentPairProtocolData) GetPair() (*liquiditytypes.Pair, error) {
	var pairData liquiditytypes.Pair
	if len(cpd.PairData) > 0 {
		err := json.Unmarshal(cpd.PairData, &pairData)
		if err != nil {
			return nil, fmt.Errorf("1: unable to unmarshal concrete PairData: %w", err)
		}
	}
	return &pairData, nil
}

type CrescentReserveAddressBalanceProtocolData struct {
	ReserveAddress string
	Denom          string
	Balance        json.RawMessage
}

func (crd CrescentReserveAddressBalanceProtocolData) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
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
}

func (cpd CrescentPoolCoinSupplyProtocolData) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
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
