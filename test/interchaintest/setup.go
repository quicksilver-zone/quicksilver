package interchaintest

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/icza/dyno"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"

	istypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

var (
	QuickSilverE2ERepo  = "quicksilverzone/quicksilver-e2e"
	QuicksilverMainRepo = "quicksilverzone/quicksilver"

	repo, version = GetDockerImageInfo()

	QuicksilverImage = ibc.DockerImage{
		Repository: repo,
		Version:    version,
		UidGid:     "1025:1025",
	}

	XccLookupImage = ibc.DockerImage{
		Repository: "quicksilverzone/xcclookup",
		Version:    "v0.4.3",
		UidGid:     "1026:1026",
	}

	ICQImage = ibc.DockerImage{
		Repository: "quicksilverzone/interchain-queries",
		Version:    "e2e",
		UidGid:     "1027:1027",
	}

	pathQuicksilverJuno = "quicksilver-juno"
	genesisWalletAmount = math.NewInt(10_000_000_000)
	votingPeriod     = "30s"
	maxDepositPeriod = "10s"
)

func createConfig() (ibc.ChainConfig, error) {
	return ibc.ChainConfig{
			Type:                "cosmos",
			Name:                "quicksilver",
			ChainID:             "quicksilver-2",
			Images:              []ibc.DockerImage{QuicksilverImage},
			Bin:                 "quicksilverd",
			Bech32Prefix:        "quick",
			Denom:               "uqck",
			GasPrices:           "0.0uqck",
			GasAdjustment:       1.1,
			TrustingPeriod:      "112h",
			NoHostMount:         false,
			ModifyGenesis:        ModifyGenesisShortProposals(votingPeriod, maxDepositPeriod),
			ConfigFileOverrides: nil,
			EncodingConfig:      quicksilverEncoding(),
			SidecarConfigs: []ibc.SidecarConfig{
				{
					ProcessName:      "icq",
					Image:            ICQImage,
					Ports:            []string{"2112"},
					StartCmd:         []string{"interchain-queries", "run", "--home", "/var/sidecar-processes/icq"},
					PreStart:         true,
					ValidatorProcess: false,
				},
				{
					ProcessName:      "xcc",
					Image:            XccLookupImage,
					Ports:            []string{"3033"},
					StartCmd:         []string{"/xcc", "-a", "serve", "-f", "/var/sidecar/processes/xcc/config.yaml"},
					PreStart:         true,
					ValidatorProcess: false,
				},
			},
		},
		nil
}

// quicksilverEncoding registers the Quicksilver specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func quicksilverEncoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	istypes.RegisterInterfaces(cfg.InterfaceRegistry)
	govv1.RegisterInterfaces(cfg.InterfaceRegistry)
	return &cfg
}


func ModifyGenesisShortProposals(votingPeriod, maxDepositPeriod string) func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(chainConfig ibc.ChainConfig, genbz []byte) ([]byte, error) {
		g := make(map[string]interface{})
		if err := json.Unmarshal(genbz, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}
		if err := dyno.Set(g, votingPeriod, "app_state", "gov", "params", "voting_period"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		if err := dyno.Set(g, maxDepositPeriod, "app_state", "gov", "params", "max_deposit_period"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		if err := dyno.Set(g, chainConfig.Denom, "app_state", "gov", "params", "min_deposit", 0, "denom"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		out, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}
		return out, nil
	}
}