package interchaintest

import (
	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
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

	// XccLookupImage = ibc.DockerImage{
	// 	Repository: "quicksilverzone/xcclookup",
	// 	Version:    "v0.4.3",
	// 	UidGid:     "1026:1026",
	// }

	ICQImage = ibc.DockerImage{
		Repository: "quicksilverzone/interchain-queries",
		Version:    "v1.0.0-beta.2",
		UidGid:     "1000:1000",
	}

	pathQuicksilverJuno = "quicksilver-juno"
	genesisWalletAmount = math.NewInt(10_000_000)
)

func createConfig() (ibc.ChainConfig, error) {
	genesis := []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.voting_params.voting_period", "15s"),
		cosmos.NewGenesisKV("app_state.gov.deposit_params.max_deposit_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.deposit_params.min_deposit.0.denom", "uqck"),
		cosmos.NewGenesisKV("app_state.gov.deposit_params.min_deposit.0.amount", "1"),

		cosmos.NewGenesisKV("app_state.epochs.epochs.0.duration", "60s"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.0.start_time", "0001-01-01T00:00:00Z"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.0.current_epoch_start_time", "0001-01-01T00:00:00Z"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.0.current_epoch", "0"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.0.identifier", "epoch"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.0.epoch_counting_started", false),

		cosmos.NewGenesisKV("app_state.epochs.epochs.1.duration", "30s"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.1.start_time", "0001-01-01T00:00:00Z"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.1.current_epoch_start_time", "0001-01-01T00:00:00Z"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.1.current_epoch", "0"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.1.identifier", "day"),
		cosmos.NewGenesisKV("app_state.epochs.epochs.1.epoch_counting_started", false),

		cosmos.NewGenesisKV("app_state.interchainstaking.params.unbonding_enabled", true),
		cosmos.NewGenesisKV("app_state.interchainstaking.params.deposit_interval", "10"),
		cosmos.NewGenesisKV("app_state.mint.params.epoch_identifier", "epoch"),
	}

	return ibc.ChainConfig{
			Type:                "cosmos",
			Name:                "quicksilver",
			ChainID:             "quicksilver-1",
			Images:              []ibc.DockerImage{QuicksilverImage},
			Bin:                 "quicksilverd",
			Bech32Prefix:        "quick",
			Denom:               "uqck",
			GasPrices:           "0.0uqck",
			GasAdjustment:       1.1,
			TrustingPeriod:      "112h",
			NoHostMount:         false,
			ModifyGenesis:       cosmos.ModifyGenesis(genesis),
			ConfigFileOverrides: map[string]any{"config/config.toml": testutil.Toml{"consensus": testutil.Toml{"timeout_commit": "1s", "timeout_propose": "500ms", "timeout_prevote": "500ms", "timeout_precommit": "500ms"}}},
			EncodingConfig:      nil,
			SidecarConfigs: []ibc.SidecarConfig{
				{
					ProcessName:      "icq",
					Image:            ICQImage,
					Ports:            []string{"2112"},
					StartCmd:         []string{"icq-relayer", "start", "--home", "/icq/.icq-relayer"},
					PreStart:         false,
					ValidatorProcess: false,
					HomeDir:          "/icq/.icq-relayer",
				},
				// {
				// 	ProcessName:      "xcc",
				// 	Image:            XccLookupImage,
				// 	Ports:            []string{"3033"},
				// 	StartCmd:         []string{"/xcc", "-a", "serve", "-f", "/var/sidecar/processes/xcc/config.yaml"},
				// 	PreStart:         true,
				// 	ValidatorProcess: false,
				// },
			},
		},
		nil
}
