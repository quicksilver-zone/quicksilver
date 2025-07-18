package types

import (
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/logger"
)

type Config struct {
	BindPort    int               `yaml:"bind_port"`
	Timeout     int               `yaml:"timeout"`
	SourceChain string            `yaml:"source_chain"`
	SourceLcd   string            `yaml:"source_lcd"`
	Chains      map[string]string `yaml:"chains"`
	Ignore      Ignores           `yaml:"ignores"`
	Mocks       Mocks             `yaml:"mocks"`
	Logging     LoggingConfig     `yaml:"logging"`
}

// LoggingConfig defines the logging configuration
type LoggingConfig struct {
	Level string `yaml:"level" default:"info"`
}

// GetLogLevel returns the configured log level, defaulting to info if not set
func (c *LoggingConfig) GetLogLevel() logger.LogLevel {
	if c.Level == "" {
		return logger.InfoLevel
	}
	return logger.LogLevel(c.Level)
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
	OsmosisPools   []prewards.OsmosisPoolProtocolData `yaml:"osmosis_pools"`
	Connections    []prewards.ConnectionProtocolData  `yaml:"connections"`
	UmeeParams     []prewards.UmeeParamsProtocolData  `yaml:"umee_params"`
	MembraneParams []prewards.MembraneProtocolData    `yaml:"membrane_params"`
}
