package types

import (
	"time"

	osmosisgammtypes "github.com/osmosis-labs/osmosis/v9/x/gamm/types"
)

type OsmosisPoolProtocolData struct {
	PoolID      uint64
	PoolName    string
	LastUpdated time.Time
	PoolData    osmosisgammtypes.PoolI
	Zones       []string
}
