package types

import (
	"log"
	"maps"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	gokitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Config represents the config file for the relayer
type Config struct {
	BindPort       int
	MaxMsgsPerTx   int
	MaxTxsPerQuery uint64
	AllowedQueries []string
	SkipEpoch      bool
	HA             HAConfig
	DefaultChain   *ChainConfig
	Chains         map[string]*ReadOnlyChainConfig
	ProtoCodec     *codec.ProtoCodec `toml:"-"`
	ClientContext  *client.Context   `toml:"-"`
	HomePath       string            `toml:"-"`
}

type HAConfig struct {
	NodeCount  int
	NodeIndex  int
	Redundancy int
}

var (
	DefaultHomePath, _ = os.UserHomeDir()
	DefaultConfigPath  = filepath.Join(DefaultHomePath, ".icq-relayer")
)

func InitializeConfigFromToml(homepath string, logger gokitlog.Logger) Config {
	config := Config{}
	_, err := toml.DecodeFile(filepath.Join(homepath, "config.toml"), &config)
	if err != nil {
		level.Warn(logger).Log("msg", "Error decoding config", "err", err)
	}

	if config.DefaultChain == nil {
		config = NewConfig()
		file, err := os.Create(filepath.Join(homepath, "config.toml"))
		if err != nil {
			level.Error(logger).Log("msg", "Error creating config file", "err", err)
			log.Fatalf("Error creating config file: %v", err)
		}
		if err := toml.NewEncoder(file).Encode(config); err != nil {
			file.Close()
			level.Error(logger).Log("msg", "Error encoding config", "err", err)
		}
		file.Close()
	}
	config.HomePath = homepath
	if config.MaxTxsPerQuery == 0 {
		config.MaxTxsPerQuery = 50
	}

	if config.MaxMsgsPerTx == 0 {
		config.MaxMsgsPerTx = 40
	}
	return config
}

func NewConfig() Config {
	return Config{
		BindPort:       2112,
		MaxMsgsPerTx:   40,
		MaxTxsPerQuery: 50,
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
			"archway-1":      DefaultReadOnlyChainConfig("archway-1", "https://archway-1.rpc.quicksilver.zone:443"),
			"injective-1":    DefaultReadOnlyChainConfig("injective-1", "https://injective-1.rpc.quicksilver.zone:443"),
			"centauri-1":     DefaultReadOnlyChainConfig("centauri-1", "https://centauri-1.rpc.quicksilver.zone:443"),
			"phoenix-1":      DefaultReadOnlyChainConfig("phoenix-1", "https://phoenix-1.rpc.quicksilver.zone:443"),
			"omniflixhub-1":  DefaultReadOnlyChainConfig("omniflixhub-1", "https://omniflixhub-1.rpc.quicksilver.zone:443"),
			"Oraichain":      DefaultReadOnlyChainConfig("Oraichain", "https://oraichain.rpc.quicksilver.zone:443"),
		},
		HA: HAConfig{
			NodeIndex:  0,
			NodeCount:  1,
			Redundancy: 1,
		},
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

func (c *Config) FilterHA() {
	seed := 0

	var result map[int][]string
	chains := slices.Collect(maps.Keys(c.Chains))
	sort.Strings(chains)
	for true {
		seed += 1
		result = DistributeConfigs(c.HA.NodeCount, c.HA.Redundancy, seed, chains)
		if CheckDistributions(result) {
			break
		}
	}

	filterMap(c.Chains, result[c.HA.NodeIndex])
}

func filterMap[V any, K comparable](A map[K]*V, B []K) {
	// Step 1: Convert slice B into a set for fast lookup
	keySet := make(map[K]bool)
	for _, key := range B {
		keySet[key] = true
	}
	for key := range A {
		if _, exists := keySet[key]; !exists {
			delete(A, key)
		}
	}
}

func DistributeConfigs(N, R, seed int, configs []string) map[int][]string {
	if N == 0 {
		N = 1
	}
	if R == 0 {
		R = 1
	}
	nodeConfigs := make(map[int][]string)

	// Step 1: Create a list with `redundancy` copies of each config
	configPool := make([]string, 0, len(configs)*R)
	for i := 0; i < len(configs); i++ {
		for j := 0; j < R; j++ {
			configPool = append(configPool, configs[i])
		}
	}

	// Step 2: Shuffle the config pool with a fixed seed for deterministic results

	r := rand.New(rand.NewSource(int64(seed)))
	r.Shuffle(len(configPool), func(i, j int) { configPool[i], configPool[j] = configPool[j], configPool[i] })
	for i, config := range configPool {
		node := i % N // Ensures an even spread
		nodeConfigs[node] = append(nodeConfigs[node], config)
	}

	return nodeConfigs
}

func CheckDistributions(dist map[int][]string) bool {
	for i := range maps.Values(dist) {
		if !IsUniqueSliceElements(i) {
			return false
		}
	}
	return true
}

func IsUniqueSliceElements[T comparable](inputSlice []T) bool {
	seen := make(map[T]bool, len(inputSlice))
	for _, element := range inputSlice {
		if seen[element] {
			return false
		}
		seen[element] = true
	}
	return true
}
