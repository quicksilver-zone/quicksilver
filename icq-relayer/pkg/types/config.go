package types

import (
	"fmt"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/BurntSushi/toml"
)

// Config represents the config file for the relayer
type Config struct {
	BindPort       int
	MaxMsgsPerTx   int
	AllowedQueries []string
	SkipEpoch      bool
	DefaultChain   *ChainConfig
	Chains         map[string]*ReadOnlyChainConfig
	ProtoCodec     *codec.ProtoCodec `toml:"-"`
	ClientContext  *client.Context   `toml:"-"`
	HomePath       string            `toml:"-"`
}

var DefaultHomePath = "~/.icq-relayer"

func InitializeConfigFromToml(homepath string) Config {
	config := NewConfig()
	_, err := toml.DecodeFile(filepath.Join(homepath, "config.toml"), &config)
	if err != nil {
		//log.Fatal().Msg(fmt.Sprintf("Error Decoding config: %v\n", err.Error()))
		fmt.Printf("Error Decoding config: %v\n", err.Error())
	}
	//zerolog.SetGlobalLevel(config.LogLevel)
	config.HomePath = homepath
	return config
}

func NewConfig() Config {
	return Config{
		BindPort:       2112,
		MaxMsgsPerTx:   40,
		AllowedQueries: []string{},
		SkipEpoch:      false,
		DefaultChain: &ChainConfig{
			ReadOnlyChainConfig: DefaultReadOnlyChainConfig("quicksilver-2", "https://rpc.quicksilver.zone:443"),
			Prefix:              "quick",
			MnemonicPath:        "./seed",
			GasLimit:            150000000,
			GasPrice:            "0.00025uqck",
			GasMultiplier:       1.25,
		},
		Chains: map[string]*ReadOnlyChainConfig{
			"cosmoshub-4":    DefaultReadOnlyChainConfig("cosmoshub-4", "https://rpc.cosmoshub-4.quicksilver.zone:443"),
			"osmosis-1":      DefaultReadOnlyChainConfig("osmosis-1", "https://rpc.osmosis-1.quicksilver.zone:443"),
			"regen-1":        DefaultReadOnlyChainConfig("regen-1", "https://rpc.regen-1.quicksilver.zone:443"),
			"stargaze-1":     DefaultReadOnlyChainConfig("stargaze-1", "https://rpc.stargaze-1.quicksilver.zone:443"),
			"juno-1":         DefaultReadOnlyChainConfig("juno-1", "https://rpc.juno-1.quicksilver.zone:443"),
			"sommelier-3":    DefaultReadOnlyChainConfig("sommelier-3", "https://rpc.sommelier-3.quicksilver.zone:443"),
			"ssc-1":          DefaultReadOnlyChainConfig("ssc-1", "https://rpc.ssc-1.quicksilver.zone:443"),
			"dydx-mainnet-1": DefaultReadOnlyChainConfig("dydx-mainnet-1", "https://rpc.dydx-mainnet-1.quicksilver.zone:443"),
			"agoric-3":       DefaultReadOnlyChainConfig("agoric-3", "https://rpc.agoric-3.quicksilver.zone:443"),
		},
		ProtoCodec:    nil,
		ClientContext: &client.Context{},
	}
}

func DefaultReadOnlyChainConfig(chainID string, rpcUrl string) *ReadOnlyChainConfig {
	return &ReadOnlyChainConfig{
		ChainID:                     chainID,
		RpcUrl:                      rpcUrl,
		ConnectTimeoutSeconds:       10,
		QueryTimeoutSeconds:         5,
		QueryRetries:                5,
		QueryRetryDelayMilliseconds: 400,
	}
}
