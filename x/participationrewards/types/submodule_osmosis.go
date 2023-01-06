package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm/pool-models/balancer"
	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm/pool-models/stableswap"
)

// OsmosisPoolProtocolData defines protocol state to track qAssets locked in
// Osmosis pools.
type OsmosisPoolProtocolData struct {
	PoolID      uint64
	PoolName    string
	LastUpdated time.Time
	PoolData    json.RawMessage
	PoolType    string
	Zones       map[string]string // chainID: IBC/denom
}

func (opd *OsmosisPoolProtocolData) GetPool() (gamm.PoolI, error) {
	switch opd.PoolType {
	case "balancer":
		var poolData balancer.Pool
		if len(opd.PoolData) > 0 {
			err := json.Unmarshal(opd.PoolData, &poolData)
			if err != nil {
				return nil, fmt.Errorf("1: unable to unmarshal concrete PoolData: %w", err)
			}
		}
		return &poolData, nil

	case "stableswap":
		var poolData stableswap.Pool
		if len(opd.PoolData) > 0 {
			err := json.Unmarshal(opd.PoolData, &poolData)
			if err != nil {
				return nil, fmt.Errorf("2: unable to unmarshal concrete PoolData: %w", err)
			}
		}
		return &poolData, nil
	default:
		// this looks like an upgrade case fallback handler?
		// should probably be changed to a proper error case for unknown type
		// at some future point...
		var poolData balancer.Pool
		if len(opd.PoolData) > 0 {
			err := json.Unmarshal(opd.PoolData, &poolData)
			if err != nil {
				return nil, fmt.Errorf("3: unable to unmarshal concrete PoolData: %w", err)
			}
		}
		return &poolData, nil
	}
}

// ValidateBasic satisfies ProtocolDataI and validates basic stateless data.
// LastUpdated and PoolData requires stateful access of keeper to validate.
func (opd OsmosisPoolProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	if opd.PoolID == 0 {
		errors["PoolID"] = ErrUndefinedAttribute
	}

	if len(opd.PoolName) == 0 {
		errors["PoolName"] = ErrUndefinedAttribute
	}

	if len(opd.PoolType) == 0 {
		errors["PoolType"] = ErrUndefinedAttribute
	}

	i := 0
	for chainID, denom := range opd.Zones {
		el := fmt.Sprintf("Zones[%d]", i)

		if len(chainID) == 0 {
			errors[el+" key"] = fmt.Errorf("%w, chainID", ErrUndefinedAttribute)
		}

		if len(denom) == 0 {
			errors[el+" value"] = fmt.Errorf("%w, IBC/denom", ErrUndefinedAttribute)
		}

		i++
	}

	if i == 0 {
		errors["Zones"] = ErrUndefinedAttribute
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

// -----------------------------------------------------

type OsmosisParamsProtocolData struct {
	ChainID string
}

// ValidateBasic satisfies ProtocolDataI and validates basic stateless data.
// LastUpdated and PoolData requires stateful access of keeper to validate.
func (oppd OsmosisParamsProtocolData) ValidateBasic() error {
	errors := make(map[string]error)

	if len(oppd.ChainID) == 0 {
		errors["ChainID"] = ErrUndefinedAttribute
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
