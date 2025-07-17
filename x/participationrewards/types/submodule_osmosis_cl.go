package types

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/multierr"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clpool "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/model"
	cl "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/types"
	"github.com/quicksilver-zone/quicksilver/utils"
)

const (
	PoolTypeCL = "concentrated-liquidity"
)

// OsmosisPoolProtocolData defines protocol state to track qAssets locked in
// Osmosis pools.
type OsmosisClPoolProtocolData struct {
	PoolID         uint64
	PoolName       string
	LastUpdated    time.Time
	PoolData       json.RawMessage
	PoolType       string
	Denoms         map[string]DenomWithZone
	IsIncentivized bool
}

func (opd *OsmosisClPoolProtocolData) GetPool() (cl.ConcentratedPoolExtension, error) {
	var poolData clpool.Pool
	if len(opd.PoolData) > 0 {
		err := json.Unmarshal(opd.PoolData, &poolData)
		if err != nil {
			return nil, fmt.Errorf("1: unable to unmarshal concrete PoolData: %w", err)
		}
	}
	return &poolData, nil
}

// ValidateBasic satisfies ProtocolDataI and validates basic stateless data.
// LastUpdated and PoolData requires stateful access of keeper to validate.
func (opd *OsmosisClPoolProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if opd.PoolID == 0 {
		errs["PoolID"] = ErrUndefinedAttribute
	}

	if opd.PoolName == "" {
		errs["PoolName"] = ErrUndefinedAttribute
	}

	if opd.PoolType == "" {
		errs["PoolType"] = ErrUndefinedAttribute
	}

	i := 0
	for _, ibcdenom := range utils.Keys(opd.Denoms) {
		el := fmt.Sprintf("Denoms[%s]", ibcdenom)

		if opd.Denoms[ibcdenom].ChainID == "" {
			errs[el+" key"] = fmt.Errorf("%w, chainID", ErrInvalidChainID)
		}

		if opd.Denoms[ibcdenom].Denom == "" || sdk.ValidateDenom(opd.Denoms[ibcdenom].Denom) != nil {
			errs[el+" value"] = fmt.Errorf("%w, IBC/denom", ErrInvalidDenom)
		}

		i++
	}

	if i == 0 {
		errs["Zones"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		var errList []error
		for _, err := range errs {
			errList = append(errList, err)
		}
		return multierr.Combine(errList...)
	}

	return nil
}

func (opd *OsmosisClPoolProtocolData) GenerateKey() []byte {
	return []byte(fmt.Sprintf("%d", opd.PoolID))
}
