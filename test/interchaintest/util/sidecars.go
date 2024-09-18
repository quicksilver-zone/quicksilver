package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/stretchr/testify/require"
)

func RunICQ(t *testing.T, ctx context.Context, quicksilver, juno *cosmos.CosmosChain, icqUser ibc.Wallet) {
	t.Helper()

	var icq *cosmos.SidecarProcess
	for _, sidecar := range quicksilver.Sidecars {
		if sidecar.ProcessName == "icq" {
			icq = sidecar
			break
		}
	}
	require.NotNil(t, icq)

	icq.WriteFile(ctx, []byte(icqUser.Mnemonic()), "seed")

	file := fmt.Sprintf(`BindPort = 2112
MaxMsgsPerTx = 50
AllowedQueries = []
SkipEpoch = false

[DefaultChain]
ChainID = "%s"
RpcUrl = "%s"
ConnectTimeoutSeconds = 10
QueryTimeoutSeconds = 5
QueryRetries = 5
QueryRetryDelayMilliseconds = 400
MnemonicPath = "/icq/.icq-relayer/seed"
Prefix = "quick"
TxSubmitTimeoutSeconds = 0
GasLimit = 150000000
GasPrice = "0.003uqck"
GasMultiplier = 1.3

[Chains]
[Chains.%s]
ChainID = "%s"
RpcUrl = "%s"
ConnectTimeoutSeconds = 10
QueryTimeoutSeconds = 5
QueryRetries = 5
QueryRetryDelayMilliseconds = 400
`,
		quicksilver.Config().ChainID,
		quicksilver.GetRPCAddress(),
		juno.Config().ChainID,
		juno.Config().ChainID,
		juno.GetRPCAddress(),
	)

	err := icq.WriteFile(ctx, []byte(file), "config.toml")
	require.NoError(t, err)

	err = icq.CreateContainer(ctx)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := icq.StopContainer(ctx)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}

			err = icq.RemoveContainer(ctx)
			if err != nil {
				panic(fmt.Errorf("an error occurred while removing the container: %s", err))
			}
		},
	)
	err = icq.StartContainer(ctx)
	require.NoError(t, err)

}

func RunXCC(t *testing.T, ctx context.Context, quicksilver, juno *cosmos.CosmosChain) {
	t.Helper()

	var xcc *cosmos.SidecarProcess
	for _, sidecar := range quicksilver.Sidecars {
		if sidecar.ProcessName == "xcc" {
			xcc = sidecar
		}
	}
	require.NotNil(t, xcc)

	containerCfg := "config.yaml"

	file := fmt.Sprintf(`source_chain: '%s'
chains:
  quick-1: '%s'
  juno-1: '%s'
`,
		quicksilver.Config().ChainID,
		quicksilver.GetRPCAddress(),
		juno.GetRPCAddress(),
	)

	err := xcc.WriteFile(ctx, []byte(file), containerCfg)
	require.NoError(t, err)
	_, err = xcc.ReadFile(ctx, containerCfg)
	require.NoError(t, err)

	err = xcc.StartContainer(ctx)
	require.NoError(t, err)

	// err = xcc.Running(ctx)
	// require.NoError(t, err)
}
