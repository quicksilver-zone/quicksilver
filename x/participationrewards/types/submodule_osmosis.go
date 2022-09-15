package types

import (
	"time"

	osmosisgammtypes "github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
)

type OsmosisPoolProtocolData struct {
	PoolID      uint64
	PoolName    string
	LastUpdated time.Time
	PoolData    osmosisgammtypes.PoolI
	Zones       map[string]string // chainID: IBC/denom
}
