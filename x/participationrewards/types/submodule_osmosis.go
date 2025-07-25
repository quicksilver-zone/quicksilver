package types

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/multierr"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/balancer"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/stableswap"
	gamm "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/types"
	"github.com/quicksilver-zone/quicksilver/utils"
)

const (
	PoolTypeBalancer   = "balancer"
	PoolTypeStableSwap = "stableswap"
)

// OsmosisPoolProtocolData defines protocol state to track qAssets locked in
// Osmosis pools.
type OsmosisPoolProtocolData struct {
	PoolID         uint64
	PoolName       string
	LastUpdated    time.Time
	PoolData       json.RawMessage
	PoolType       string
	Denoms         map[string]DenomWithZone
	IsIncentivized bool
}

type DenomWithZone struct {
	Denom   string
	ChainID string
}

func (opd *OsmosisPoolProtocolData) GetPool() (gamm.CFMMPoolI, error) {
	switch opd.PoolType {
	case PoolTypeBalancer:
		var poolData balancer.Pool
		if len(opd.PoolData) > 0 {
			err := json.Unmarshal(opd.PoolData, &poolData)
			if err != nil {
				return nil, fmt.Errorf("1: unable to unmarshal concrete PoolData: %w", err)
			}
		}
		return &poolData, nil

	case PoolTypeStableSwap:
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
func (opd *OsmosisPoolProtocolData) ValidateBasic() error {
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
		return multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return nil
}

func (opd *OsmosisPoolProtocolData) GenerateKey() []byte {
	return []byte(fmt.Sprintf("%d", opd.PoolID))
}

// -----------------------------------------------------

type OsmosisParamsProtocolData struct {
	ChainID   string
	BaseDenom string
	BaseChain string
}

// ValidateBasic satisfies ProtocolDataI and validates basic stateless data.
// LastUpdated and PoolData requires stateful access of keeper to validate.
func (oppd *OsmosisParamsProtocolData) ValidateBasic() error {
	errs := make(map[string]error)

	if oppd.ChainID == "" {
		errs["ChainID"] = ErrUndefinedAttribute
	}

	if oppd.BaseChain == "" {
		errs["BaseChain"] = ErrUndefinedAttribute
	}

	if oppd.BaseDenom == "" {
		errs["BaseDenom"] = ErrUndefinedAttribute
	}

	if len(errs) > 0 {
		return multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return nil
}

func (*OsmosisParamsProtocolData) GenerateKey() []byte {
	return []byte(OsmosisParamsKey)
}
