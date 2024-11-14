package types

import (
	"fmt"
	"log"
	"os"
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

var (
	DefaultHomePath, _ = os.UserHomeDir()
	DefaultConfigPath  = filepath.Join(DefaultHomePath, ".icq-relayer")
)

func InitializeConfigFromToml(homepath string) Config {
	config := Config{}
	_, err := toml.DecodeFile(filepath.Join(homepath, "config.toml"), &config)
	if err != nil {
		log.Printf("Error decoding config: %v\n", err)
	}

	if config.DefaultChain == nil {
		config = NewConfig()
		file, err := os.Create(filepath.Join(homepath, "config.toml"))
		if err != nil {
			log.Fatalf("Error creating config file: %v", err)
		}
		if err := toml.NewEncoder(file).Encode(config); err != nil {
			file.Close()
			log.Fatalf("Error encoding config: %v", err)
		}
		file.Close()
	}
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
			ReadOnlyChainConfig: DefaultReadOnlyChainConfig("quicksilver-2", "https://quicksilver-2.rpc.quicksilver.zone:443"),
			Prefix:              "quick",
			MnemonicPath:        "./seed",
			GasLimit:            150000000,
			GasPrice:            "0.00025uqck",
			GasMultiplier:       1.25,
		},
		Chains: map[string]*ReadOnlyChainConfig{
			"cosmoshub-4":    DefaultReadOnlyChainConfig("cosmoshub-4", "https://cosmoshub-4.rpc.quicksilver.zone:443"),
			"osmosis-1":      DefaultReadOnlyChainConfig("osmosis-1", "https://osmosis-1.rpc.quicksilver.zone:443"),
			"regen-1":        DefaultReadOnlyChainConfig("regen-1", "https://regen-1.rpc.quicksilver.zone:443"),
			"stargaze-1":     DefaultReadOnlyChainConfig("stargaze-1", "https://stargaze-1.rpc.quicksilver.zone:443"),
			"juno-1":         DefaultReadOnlyChainConfig("juno-1", "https://juno-1.rpc.quicksilver.zone:443"),
			"sommelier-3":    DefaultReadOnlyChainConfig("sommelier-3", "https://sommelier-3.rpc.quicksilver.zone:443"),
			"ssc-1":          DefaultReadOnlyChainConfig("ssc-1", "https://ssc-1.rpc.quicksilver.zone:443"),
			"dydx-mainnet-1": DefaultReadOnlyChainConfig("dydx-mainnet-1", "https://dydx-mainnet-1.rpc.quicksilver.zone:443"),
			"agoric-3":       DefaultReadOnlyChainConfig("agoric-3", "https://agoric-3.rpc.quicksilver.zone:443"),
			"secret-4":       DefaultReadOnlyChainConfig("secret-4", "https://secret-4.rpc.quicksilver.zone:443"),
			"celestia":       DefaultReadOnlyChainConfig("celestia", "https://celestia.rpc.quicksilver.zone:443"),
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
