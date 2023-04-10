package interchaintest

import (
	"github.com/strangelove-ventures/interchaintest/v5/ibc"
)

var (
	QuickSilverE2ERepo  = "ghcr.io/ingenuity-build/quicksilver-e2e"
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
		Version:    "v0.8.8-beta1",
		UidGid:     "1027:1027",
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
		SidecarConfigs: []ibc.SidecarConfig{
			{
				ProcessName:      "icq",
				Image:            ICQImage,
				Ports:            []string{"2112"},
				StartCmd:         nil,
				PreStart:         false,
				ValidatorProcess: true,
			},
			{
				ProcessName:      "xcc",
				Image:            XccLookupImage,
				Ports:            []string{"3033"},
				StartCmd:         nil,
				PreStart:         false,
				ValidatorProcess: true,
			},
		},
	}

	pathQuicksilverOsmosis = "quicksilver-osmosis"
	genesisWalletAmount    = int64(10_000_000)
)
