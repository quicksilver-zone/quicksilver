package types

import (
	"fmt"
	"time"

	osmosisgammtypes "github.com/ingenuity-build/quicksilver/osmosis-types/gamm"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
)

type OsmosisPoolProtocolData struct {
	PoolID      uint64
	PoolName    string
	LastUpdated time.Time
	PoolData    osmosisgammtypes.PoolI
	Zones       map[string]string // chainID: IBC/denom
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
