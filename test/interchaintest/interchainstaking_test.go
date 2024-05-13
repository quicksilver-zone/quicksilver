package interchaintest

import (
	"context"
	"fmt"
	"os"
	"testing"

	module "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"gopkg.in/yaml.v2"
)

// TestInterchainStaking TODO
func TestInterchainStaking(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Quicksilver and gaia
	numVals := 3
	numFullNodes := 3

	config, err := createConfig()
	require.NoError(t, err)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "quicksilver",
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "gaia",
			Version:       "v14.1.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				GasPrices:      "0.0uatom",
				EncodingConfig: gaiaEncoding(),
			},
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	quicksilver, gaia := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Create relayer factory to utilize the go-relayer
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(t, client, network)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain().
		AddChain(quicksilver).
		AddChain(gaia).
		AddRelayer(r, "rly").
		AddLink(interchaintest.InterchainLink{
			Chain1:  quicksilver,
			Chain2:  gaia,
			Relayer: r,
			Path:    pathQuicksilverGaia,
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

	err = ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,

		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Start the relayer
	require.NoError(t, r.StartRelayer(ctx, eRep, pathQuicksilverGaia))
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}
		},
	)

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, quicksilver, gaia)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, quicksilver, gaia)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	quickUser, gaiaUser := users[0], users[1]

	quickUserAddr := quickUser.FormattedAddress()
	gaiaUserAddr := gaiaUser.FormattedAddress()
	_ = quickUserAddr
	_ = gaiaUserAddr

	runSidecars(ctx, t, quicksilver, gaia)
}

func runSidecars(ctx context.Context, t *testing.T, quicksilver, gaia *cosmos.CosmosChain) {
	t.Helper()

	runICQ(ctx, t, quicksilver, gaia)
	// runXCC(t, ctx, quicksilver, gaia)
}

func runICQ(ctx context.Context, t *testing.T, quicksilver, gaia *cosmos.CosmosChain) *cosmos.SidecarProcess {
	t.Helper()

	var icq *cosmos.SidecarProcess
	for _, sidecar := range quicksilver.Sidecars {
		if sidecar.ProcessName == "icq" {
			icq = sidecar
		}
	}
	require.NotNil(t, icq)

	containerCfg := "config.yaml"

	cfg := Config{
		BindPort:     2112,
		MaxMsgsPerTx: 40,
		DefaultChain: "quicksilver-2",
		AllowedQueries: []string{
			"deposittx",
			"depositinterval",
		},
		Chains: map[string]ChainClientConfig{
			"quicksilver-2": {
				Key:            "default",
				ChainID:        quicksilver.Config().ChainID,
				RPCAddr:        quicksilver.GetRPCAddress(),
				GRPCAddr:       quicksilver.GetGRPCAddress(),
				AccountPrefix:  quicksilver.Config().Bech32Prefix,
				KeyringBackend: "test",
				GasAdjustment:  1.3,
				GasPrices:      "0.0001uqck",
				KeyDirectory:   fmt.Sprintf("%s/keys", icq.HomeDir()),
				Debug:          false,
				Timeout:        "20s",
				BlockTimeout:   "10s",
				OutputFormat:   "json",
				SignModeStr:    "direct",
			},
			"gaia-1": {
				Key:            "default",
				ChainID:        gaia.Config().ChainID,
				RPCAddr:        gaia.GetRPCAddress(),
				GRPCAddr:       gaia.GetGRPCAddress(),
				AccountPrefix:  gaia.Config().Bech32Prefix,
				KeyringBackend: "test",
				KeyDirectory:   fmt.Sprintf("%s/keys", icq.HomeDir()),
				GasAdjustment:  1.2,
				GasPrices:      "0.0uatom",
				MinGasAmount:   0,
				Debug:          false,
				Timeout:        "20s",
				BlockTimeout:   "10s",
				OutputFormat:   "json",
				SignModeStr:    "direct",
			},
		},
	}
	file := cfg.MustYAML()

	err := icq.WriteFile(ctx, file, containerCfg)
	require.NoError(t, err)
	_, err = icq.ReadFile(ctx, containerCfg)
	require.NoError(t, err)
	// Copy all file from test/interchaintest/relayer/icq/keys
	// to /icq/keys

	files, err := os.ReadDir("relayer/icq/keys/quicksilver-2/keyring-test")
	require.NoError(t, err)

	for _, file := range files {
		err := icq.CopyFile(ctx, fmt.Sprintf("relayer/icq/keys/quicksilver-2/keyring-test/%s", file.Name()), fmt.Sprintf("keys/quicksilver-2/keyring-test/%s", file.Name()))
		require.NoError(t, err)
	}
	err = icq.CreateContainer(ctx)
	require.NoError(t, err)
	err = icq.StartContainer(ctx)
	require.NoError(t, err)

	return icq
}

func runXCC(t *testing.T, ctx context.Context, quicksilver, gaia *cosmos.CosmosChain) {
	// 	t.Helper()

	// 	var xcc *cosmos.SidecarProcess
	// 	for _, sidecar := range quicksilver.Sidecars {
	// 		if sidecar.ProcessName == "xcc" {
	// 			xcc = sidecar
	// 		}
	// 	}
	// 	require.NotNil(t, xcc)

	// 	containerCfg := "config.yaml"

	// 	file := fmt.Sprintf(`source_chain: '%s'
	// chains:
	//   quick-1: '%s'
	//   gaia-1: '%s'
	// `,
	// 		quicksilver.Config().ChainID,
	// 		quicksilver.GetRPCAddress(),
	// 		gaia.GetRPCAddress(),
	// 	)

	// 	err := xcc.WriteFile(ctx, []byte(file), containerCfg)
	// 	require.NoError(t, err)
	// 	_, err = xcc.ReadFile(ctx, containerCfg)
	// 	require.NoError(t, err)

	// 	err = xcc.StartContainer(ctx)
	// 	require.NoError(t, err)

	// err = xcc.Running(ctx)
	// require.NoError(t, err)
}

// Config represents the config file for the relayer
type Config struct {
	BindPort       int                          `yaml:"bind_port" json:"bind_port"`
	MaxMsgsPerTx   int                          `yaml:"max_msgs_per_tx" json:"max_msgs_per_tx"`
	DefaultChain   string                       `yaml:"default_chain" json:"default_chain"`
	AllowedQueries []string                     `yaml:"allowed_queries" json:"allowed_queries"`
	SkipEpoch      bool                         `yaml:"skip_epoch" json:"skip_epoch"`
	Chains         map[string]ChainClientConfig `yaml:"chains" json:"chains"`
	// Cl             map[string]ChainClient       `yaml:",omitempty" json:",omitempty"`
}

// MustYAML returns the yaml string representation of the Paths
func (c Config) MustYAML() []byte {
	out, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return out
}

type ChainClientConfig struct {
	Key            string                  `json:"key" yaml:"key"`
	ChainID        string                  `json:"chain-id" yaml:"chain-id"`
	RPCAddr        string                  `json:"rpc-addr" yaml:"rpc-addr"`
	GRPCAddr       string                  `json:"grpc-addr" yaml:"grpc-addr"`
	AccountPrefix  string                  `json:"account-prefix" yaml:"account-prefix"`
	KeyringBackend string                  `json:"keyring-backend" yaml:"keyring-backend"`
	GasAdjustment  float64                 `json:"gas-adjustment" yaml:"gas-adjustment"`
	GasPrices      string                  `json:"gas-prices" yaml:"gas-prices"`
	MinGasAmount   uint64                  `json:"min-gas-amount" yaml:"min-gas-amount"`
	KeyDirectory   string                  `json:"key-directory" yaml:"key-directory"`
	Debug          bool                    `json:"debug" yaml:"debug"`
	Timeout        string                  `json:"timeout" yaml:"timeout"`
	BlockTimeout   string                  `json:"block-timeout" yaml:"block-timeout"`
	OutputFormat   string                  `json:"output-format" yaml:"output-format"`
	SignModeStr    string                  `json:"sign-mode" yaml:"sign-mode"`
	ExtraCodecs    []string                `json:"extra-codecs" yaml:"extra-codecs"`
	Modules        []module.AppModuleBasic `json:"-" yaml:"-"`
}

// type ChainClient struct {
// 	log *zap.Logger

// 	Config         *ChainClientConfig
// 	Keybase        keyring.Keyring
// 	KeyringOptions []keyring.Option
// 	RPCClient      rpcclient.Client
// 	LightProvider  provtypes.Provider
// 	Input          io.Reader
// 	Output         io.Writer
// 	// TODO: GRPC Client type?

// 	Codec Codec
// }
