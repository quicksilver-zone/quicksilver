package interchaintest

import (
	quicksilver "github.com/ingenuity-build/quicksilver/app"
	"github.com/strangelove-ventures/ibctest/v5/ibc"
)

var (
	QuickSilverE2ERepo  = "ghcr.io/cosmoscontracts/juno-e2e"
	QuicksilverMainRepo = "ghcr.io/cosmoscontracts/juno"

	repo, version = GetDockerImageInfo()

	QuicksilverImage = ibc.DockerImage{
		Repository: repo,
		Version:    version,
		UidGid:     "1025:1025",
	}

	junoConfig = ibc.ChainConfig{
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
		ModifyGenesis:       nil,
		ConfigFileOverrides: nil,
		EncodingConfig:      quicksilver.MakeEncodingConfig(),
	}

	pathQuicksilverGaia = "quicksilver-gaia"
	genesisWalletAmount = int64(10_000_000)
)
