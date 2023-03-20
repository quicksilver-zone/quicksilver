package configurer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/test/e2e/configurer/chain"
	"github.com/ingenuity-build/quicksilver/test/e2e/containers"
	"github.com/ingenuity-build/quicksilver/test/e2e/initialization"
	"github.com/ingenuity-build/quicksilver/test/e2e/util"
)

// baseConfigurer is the base implementation for the
// other 2 types of configurers. It is not meant to be used
// on its own. Instead, it is meant to be embedded
// by composition into more concrete configurers.
type baseConfigurer struct {
	chainConfigs     []*chain.Config
	containerManager *containers.Manager
	setupTests       setupFn
	syncUntilHeight  int64 // the height until which to wait for validators to sync when first started.
	t                *testing.T
}

// defaultSyncUntilHeight arbitrary small height to make sure the chain is making progress.
const defaultSyncUntilHeight = 3

func (bc *baseConfigurer) ClearResources() error {
	bc.t.Log("tearing down e2e integration test suite...")

	if err := bc.containerManager.ClearResources(); err != nil {
		return err
	}

	for _, chainConfig := range bc.chainConfigs {
		os.RemoveAll(chainConfig.DataDir)
	}
	return nil
}

func (bc *baseConfigurer) GetChainConfig(chainIndex int) *chain.Config {
	return bc.chainConfigs[chainIndex]
}

func (bc *baseConfigurer) RunValidators() error {
	for _, chainConfig := range bc.chainConfigs {
		if err := bc.runValidators(chainConfig); err != nil {
			return err
		}
	}
	return nil
}

func (bc *baseConfigurer) runValidators(chainConfig *chain.Config) error {
	bc.t.Logf("starting %s validator containers...", chainConfig.ID)
	for _, node := range chainConfig.NodeConfigs {
		if err := node.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (bc *baseConfigurer) RunIBC() error {
	// Run a hermes relayer between every possible pair of chains.
	for i := 0; i < len(bc.chainConfigs); i++ {
		for j := i + 1; j < len(bc.chainConfigs); j++ {
			if err := bc.runIBCRelayer(bc.chainConfigs[i], bc.chainConfigs[j]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *baseConfigurer) RunICQ() error {
	// Run a icq relayer between every possible pair of chains.
	for i := 0; i < len(bc.chainConfigs); i++ {
		for j := i + 1; j < len(bc.chainConfigs); j++ {
			if err := bc.runICQRelayer(bc.chainConfigs[i], bc.chainConfigs[j]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *baseConfigurer) RunXCC() error {
	// Run a icq relayer between every possible pair of chains.
	for i := 0; i < len(bc.chainConfigs); i++ {
		for j := i + 1; j < len(bc.chainConfigs); j++ {
			if err := bc.runXCCLookup(bc.chainConfigs[i], bc.chainConfigs[j]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *baseConfigurer) runIBCRelayer(chainConfigA *chain.Config, chainConfigB *chain.Config) error {
	bc.t.Log("starting Hermes relayer container...")

	tmpDir, err := os.MkdirTemp("", "quicksilver-e2e-testnet-hermes-")
	if err != nil {
		return err
	}

	hermesCfgPath := path.Join(tmpDir, "hermes")
	if err := os.MkdirAll(hermesCfgPath, 0o755); err != nil {
		return err
	}
	_, err = util.CopyFile(
		filepath.Join("./scripts/", "hermes_bootstrap.sh"),
		filepath.Join(hermesCfgPath, "hermes_bootstrap.sh"),
	)
	if err != nil {
		return err
	}

	bc.t.Logf(filepath.Join(hermesCfgPath, "hermes_bootstrap.sh"))

	nodeConfigA := chainConfigA.NodeConfigs[0]
	nodeConfigB := chainConfigB.NodeConfigs[0]
	hermesResource, err := bc.containerManager.RunHermesResource(
		bc.t,
		chainConfigA.ID,
		nodeConfigA.Name,
		nodeConfigA.Mnemonic,
		chainConfigB.ID,
		nodeConfigB.Name,
		nodeConfigB.Mnemonic,
		hermesCfgPath,
	)
	if err != nil {
		return err
	}

	require.True(bc.t, hermesResource.Container.State.Running)
	endpoint := fmt.Sprintf("http://%s/state", hermesResource.GetHostPort("3031/tcp"))
	require.Eventually(bc.t, func() bool {
		resp, err := http.Get(endpoint) //nolint:gosec
		if err != nil {
			return false
		}

		defer resp.Body.Close()

		bz, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}

		var respBody map[string]interface{}
		if err := json.Unmarshal(bz, &respBody); err != nil {
			return false
		}

		status, ok := respBody["status"].(string)
		require.True(bc.t, ok)
		result, ok := respBody["result"].(map[string]interface{})
		require.True(bc.t, ok)

		chains, ok := result["chains"].([]interface{})
		require.True(bc.t, ok)

		return status == "success" && len(chains) == 2
	},
		5*time.Minute,
		time.Second,
		"hermes relayer not healthy")

	bc.t.Logf("started Hermes relayer container: %s", hermesResource.Container.ID)

	// XXX: Give time to both networks to start, otherwise we might see gRPC
	// transport errors.
	time.Sleep(10 * time.Second)

	// create the client, connection and channel between the two Quicksilver chains
	return bc.connectIBCChains(chainConfigA, chainConfigB)
}

func (bc *baseConfigurer) runICQRelayer(chainConfigA *chain.Config, chainConfigB *chain.Config) error {
	bc.t.Log("starting ICQ relayer container...")

	tmpDir, err := os.MkdirTemp("", "quicksilver-e2e-testnet-icq-")
	if err != nil {
		return err
	}

	icqCfgPath := path.Join(tmpDir, "icq")
	if err := os.MkdirAll(icqCfgPath, 0o755); err != nil {
		return err
	}
	_, err = util.CopyFile(
		filepath.Join("./scripts/", "icq_bootstrap.sh"),
		filepath.Join(icqCfgPath, "icq_bootstrap.sh"),
	)
	if err != nil {
		return err
	}

	bc.t.Logf(filepath.Join(icqCfgPath, "icq_bootstrap.sh"))

	nodeConfigA := chainConfigA.NodeConfigs[0]
	nodeConfigB := chainConfigB.NodeConfigs[0]
	icqResource, err := bc.containerManager.RunICQResource(
		bc.t,
		chainConfigA.ID,
		nodeConfigA.Name,
		chainConfigB.ID,
		nodeConfigB.Name,
		icqCfgPath,
	)
	if err != nil {
		return err
	}

	// endpoint = fmt.Sprintf("http://%s/metrics", icqResource.GetHostPort("2112/tcp"))
	require.True(bc.t, icqResource.Container.State.Running)

	bc.t.Logf("started icq relayer container: %s", icqResource.Container.ID)

	// XXX: Give time to both networks to start, otherwise we might see gRPC
	// transport errors.
	time.Sleep(10 * time.Second)

	return nil
}

func (bc *baseConfigurer) runXCCLookup(chainConfigA *chain.Config, chainConfigB *chain.Config) error {
	bc.t.Log("starting XCCLookup container...")

	tmpDir, err := os.MkdirTemp("", "quicksilver-e2e-testnet-xcc-")
	if err != nil {
		return err
	}

	xccCfgPath := path.Join(tmpDir, "icq")
	if err := os.MkdirAll(xccCfgPath, 0o755); err != nil {
		return err
	}
	_, err = util.CopyFile(
		filepath.Join("./scripts/", "xcc_bootstrap.sh"),
		filepath.Join(xccCfgPath, "xcc_bootstrap.sh"),
	)
	if err != nil {
		return err
	}

	bc.t.Logf(filepath.Join(xccCfgPath, "xcc_bootstrap.sh"))

	nodeConfigA := chainConfigA.NodeConfigs[0]
	nodeConfigB := chainConfigB.NodeConfigs[0]

	xccResource, err := bc.containerManager.RunXCCLookupResource(
		bc.t,
		chainConfigA.ID,
		nodeConfigA.Name,
		chainConfigB.ID,
		nodeConfigB.Name,
		xccCfgPath,
	)
	if err != nil {
		return err
	}

	// endpoint = fmt.Sprintf("http://%s/metrics", icqResource.GetHostPort("2112/tcp"))
	require.True(bc.t, xccResource.Container.State.Running)

	bc.t.Logf("started XCC-Lookup container: %s", xccResource.Container.ID)

	// XXX: Give time to both networks to start, otherwise we might see gRPC
	// transport errors.
	time.Sleep(10 * time.Second)

	return nil
}

func (bc *baseConfigurer) connectIBCChains(chainA *chain.Config, chainB *chain.Config) error {
	bc.t.Logf("connecting %s and %s chains via IBC", chainA.ChainMeta.ID, chainB.ChainMeta.ID)
	cmd := []string{"hermes", "create", "channel", "--a-chain", chainA.ChainMeta.ID, "--b-chain", chainB.ChainMeta.ID, "--a-port", "transfer", "--b-port", "transfer", "--new-client-connection", "--yes"}
	_, _, err := bc.containerManager.ExecHermesCmd(bc.t, cmd, "SUCCESS")
	if err != nil {
		return err
	}
	bc.t.Logf("connected %s and %s chains via IBC", chainA.ChainMeta.ID, chainB.ChainMeta.ID)
	return nil
}

func (bc *baseConfigurer) initializeChainConfigFromInitChain(initializedChain *initialization.Chain, chainConfig *chain.Config) {
	chainConfig.ChainMeta = initializedChain.ChainMeta
	chainConfig.NodeConfigs = make([]*chain.NodeConfig, 0, len(initializedChain.Nodes))
	setupTime := time.Now()
	for i, validator := range initializedChain.Nodes {
		conf := chain.NewNodeConfig(bc.t, validator, chainConfig.ValidatorInitConfigs[i], chainConfig.ID, bc.containerManager).WithSetupTime(setupTime)
		chainConfig.NodeConfigs = append(chainConfig.NodeConfigs, conf)
	}
}
