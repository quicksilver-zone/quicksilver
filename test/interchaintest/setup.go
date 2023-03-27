package interchaintest

import (
	"github.com/strangelove-ventures/ibctest/v5/ibc"
)

var (
	QuickSilverE2ERepo  = "quicksilverzone/quicksilver@latest"
	QuicksilverMainRepo = "quicksilverzone/quicksilver@v1.2.7"

	// repo, version = GetDockerImageInfo()

	QuicksilverImage = ibc.DockerImage{
		Repository: "quicksilverzone/quicksilver",
		Version:    "latest",
		UidGid:     "1025:1025",
	}

	config = ibc.ChainConfig{
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
		EncodingConfig:      nil,
	}

	pathQuicksilverGaia = "quicksilver-gaia"
	genesisWalletAmount = int64(10_000_000)
)
