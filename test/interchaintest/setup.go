package interchaintest

import (
	"encoding/json"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/strangelove-ventures/interchaintest/v5/ibc"
	"time"
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
	genesisWalletAmount = int64(10_000_000_000_000)
)

func createConfig() (ibc.ChainConfig, error) {
	encodingConfig := app.MakeEncodingConfig()
	return ibc.ChainConfig{
			Type:           "cosmos",
			Name:           "quicksilver",
			ChainID:        "quicksilver-2",
			Images:         []ibc.DockerImage{QuicksilverImage},
			Bin:            "quicksilverd",
			Bech32Prefix:   "quick",
			Denom:          "uqck",
			GasPrices:      "0.0uqck",
			GasAdjustment:  1.1,
			TrustingPeriod: "112h",
			NoHostMount:    false,
			ModifyGenesis: func(chainCfg ibc.ChainConfig, bytes []byte) ([]byte, error) {
				genesis := app.GenesisState{}
				err := json.Unmarshal(bytes, &genesis)
				if err != nil {
					return nil, err
				}

				gov := govtypes.DefaultGenesisState()
				gov.VotingParams.VotingPeriod = time.Minute
				genesis["gov"] = encodingConfig.Marshaler.MustMarshalJSON(gov)

				return json.MarshalIndent(genesis, "", "")
			},
			ConfigFileOverrides: nil,
			EncodingConfig: &simappparams.EncodingConfig{
				InterfaceRegistry: encodingConfig.InterfaceRegistry,
				Codec:             encodingConfig.Marshaler,
				TxConfig:          encodingConfig.TxConfig,
				Amino:             encodingConfig.Amino,
			},
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
