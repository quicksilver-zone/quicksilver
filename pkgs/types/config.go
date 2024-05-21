package types

import (
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type Config struct {
	BindPort    int               `yaml:"bind_port"`
	SourceChain string            `yaml:"source_chain"`
	SourceLcd   string            `yaml:"source_lcd"`
	Chains      map[string]string `yaml:"chains"`
	Ignore      Ignores           `yaml:"ignores"`
	Mocks       Mocks             `yaml:"mocks"`
}

type Ignores []Ignore

type Ignore struct {
	Type string
	Key  string
}

func (i Ignores) GetIgnoresForType(ignoreType string) Ignores {
	out := make(Ignores, 0)
	for _, ignore := range i {
		if ignore.Type == ignoreType {
			out = append(out, ignore)
		}
	}
	return out
}

func (i Ignores) Contains(key string) bool {
	for _, ignore := range i {
		if ignore.Key == key {
			return true
		}
	}
	return false
}

const (
	IgnoreTypeLiquid        = "liquid"
	IgnoreTypeOsmosisPool   = "osmosispool"
	IgnoreTypeOsmosisCLPool = "osmosisclpool"
)

type Mocks struct {
	OsmosisPools []prewards.OsmosisPoolProtocolData `yaml:"osmosis_pools"`
	Connections  []prewards.ConnectionProtocolData  `yaml:"connections"`
	UmeeParams   []prewards.UmeeParamsProtocolData  `yaml:"umee_params"`
}
