package containers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

const (
	hermesContainerName    = "hermes-relayer"
	icqContainerName       = "icq-relayer"
	xccLookupContainerName = "xcc-lookup"

	// The maximum number of times debug logs are printed to console
	// per CLI command.
	maxDebugLogsPerCommand = 3

	GasLimit = 400000
)

var (
	// We set consensus min fee = .0025 uqck / gas * 400000 gas = 1000.
	defaultErrRegex = regexp.MustCompile(`([Ee])rror`)

	txArgs = []string{"-b=block", "--yes", "--keyring-backend=test", "--log_format=json"}
)

// Manager is a wrapper around all Docker instances, and the Docker API.
// It provides utilities to run and interact with all Docker containers used within e2e testing.
type Manager struct {
	ImageConfig
	pool              *dockertest.Pool
	network           *dockertest.Network
	resources         map[string]*dockertest.Resource
	IsDebugLogEnabled bool
}

// NewManager creates a new Manager instance and initializes
// all Docker specific utilizes. Returns an error if initialization fails.
func NewManager(isUpgrade, isFork, isDebugLogEnabled bool) (manager *Manager, err error) {
	manager = &Manager{
		ImageConfig:       NewImageConfig(isUpgrade, isFork),
		resources:         make(map[string]*dockertest.Resource),
		IsDebugLogEnabled: isDebugLogEnabled,
	}
	manager.pool, err = dockertest.NewPool("")
	if err != nil {
		return nil, err
	}
	manager.network, err = manager.pool.CreateNetwork("quicksilver-testnet")
	if err != nil {
		return nil, err
	}
	return manager, nil
}

// ExecTxCmd Runs ExecTxCmdWithSuccessString searching for `code: 0`.
func (m *Manager) ExecTxCmd(t *testing.T, chainID, containerName string, command []string) (outBuf, errBuf bytes.Buffer, err error) {
	t.Helper()
	return m.ExecTxCmdWithSuccessString(t, chainID, containerName, command, "code: 0")
}

// ExecTxCmdWithSuccessString Runs ExecCmd, with flags for txs added.
// namely adding flags `--chain-id={chain-id} -b=block --yes --keyring-backend=test "--log_format=json" --gas=400000`,
// and searching for `successStr`.
func (m *Manager) ExecTxCmdWithSuccessString(t *testing.T, chainID, containerName string, command []string, successStr string) (outBuf, errBuf bytes.Buffer, err error) {
	t.Helper()

	allTxArgs := []string{fmt.Sprintf("--chain-id=%s", chainID)}
	allTxArgs = append(allTxArgs, txArgs...)

	txCommand := append(command, allTxArgs...) //nolint:gocritic
	return m.ExecCmd(t, containerName, txCommand, successStr)
}

// ExecHermesCmd executes command on the hermes relayer container.
func (m *Manager) ExecHermesCmd(t *testing.T, command []string, success string) (outBuf, errBuf bytes.Buffer, err error) {
	t.Helper()
	return m.ExecCmd(t, hermesContainerName, command, success)
}

// ExecICQCmd executes command on the ICQ relayer container.
func (m *Manager) ExecICQCmd(t *testing.T, command []string, success string) (outBuf, errBuf bytes.Buffer, err error) {
	t.Helper()
	return m.ExecCmd(t, icqContainerName, command, success)
}

// ExecCmd executes command by running it on the node container (specified by containerName)
// success is the output of the command that needs to be observed for the command to be deemed successful.
// It is found by checking if stdout or stderr contains the success string anywhere within it.
// returns container std out, container std err, and error if any.
// An error is returned if the command fails to execute or if the success string is not found in the output.
func (m *Manager) ExecCmd(t *testing.T, containerName string, command []string, success string) (outBuf, errBuf bytes.Buffer, err error) {
	t.Helper()
	if _, ok := m.resources[containerName]; !ok {
		return bytes.Buffer{}, bytes.Buffer{}, fmt.Errorf("no resource %s found", containerName)
	}
	containerID := m.resources[containerName].Container.ID

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if m.IsDebugLogEnabled {
		t.Logf("\n\nRunning: \"%s\", success condition is \"%s\"", command, success)
	}
	maxDebugLogTriesLeft := maxDebugLogsPerCommand

	// We use the `require.Eventually` function because it is only allowed to do one transaction per block without
	// sequence numbers. For simplicity, we avoid keeping track of the sequence number and just use the `require.Eventually`.
	require.Eventually(
		t,
		func() bool {
			exec, err := m.pool.Client.CreateExec(docker.CreateExecOptions{
				Context:      ctx,
				AttachStdout: true,
				AttachStderr: true,
				Container:    containerID,
				User:         "root",
				Cmd:          command,
			})
			require.NoError(t, err)

			err = m.pool.Client.StartExec(exec.ID, docker.StartExecOptions{
				Context:      ctx,
				Detach:       false,
				OutputStream: &outBuf,
				ErrorStream:  &errBuf,
			})
			if err != nil {
				return false
			}

			errBufString := errBuf.String()
			// Note that this does not match all errors.
			// This only works if CLI outpurs "Error" or "error"
			// to stderr.
			if (defaultErrRegex.MatchString(errBufString) || m.IsDebugLogEnabled) && maxDebugLogTriesLeft > 0 {
				t.Log("\nstderr:")
				t.Log(errBufString)

				t.Log("\nstdout:")
				t.Log(outBuf.String())
				// N.B: We should not be returning false here
				// because some applications such as Hermes might log
				// "error" to stderr when they function correctly,
				// causing test flakiness. This log is needed only for
				// debugging purposes.
				maxDebugLogTriesLeft--
			}

			if success != "" {
				return strings.Contains(outBuf.String(), success) || strings.Contains(errBufString, success)
			}

			return true
		},
		time.Minute,
		50*time.Millisecond,
		fmt.Sprintf("success condition (%s) was not met.\nstdout:\n %s\nstderr:\n %s\n",
			success, outBuf.String(), errBuf.String()),
	)

	return outBuf, errBuf, nil
}

// RunHermesResource runs a Hermes container. Returns the container resource and error if any.
// the name of the hermes container is "<chain A id>-<chain B id>-relayer".
func (m *Manager) RunHermesResource(
	t *testing.T,
	chainAID,
	quickARelayerNodeName,
	quickAValMnemonic,
	chainBID,
	quickBRelayerNodeName,
	quickBValMnemonic,
	hermesCfgPath string,
) (*dockertest.Resource, error) {
	t.Helper()

	hermesResource, err := m.pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       hermesContainerName,
			Repository: m.HermesRepository,
			Tag:        m.HermesTag,
			NetworkID:  m.network.Network.ID,
			User:       "root:root",
			Mounts: []string{
				fmt.Sprintf("%s/:/root/hermes", hermesCfgPath),
			},
			ExposedPorts: []string{
				"3031",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"3031/tcp": {{HostIP: "", HostPort: "3031"}},
			},
			Env: []string{
				fmt.Sprintf("QUICK_A_E2E_CHAIN_ID=%s", chainAID),
				fmt.Sprintf("QUICK_B_E2E_CHAIN_ID=%s", chainBID),
				fmt.Sprintf("QUICK_A_E2E_VAL_MNEMONIC=%s", quickAValMnemonic),
				fmt.Sprintf("QUICK_B_E2E_VAL_MNEMONIC=%s", quickBValMnemonic),
				fmt.Sprintf("QUICK_A_E2E_VAL_HOST=%s", quickARelayerNodeName),
				fmt.Sprintf("QUICK_B_E2E_VAL_HOST=%s", quickBRelayerNodeName),
			},
			Entrypoint: []string{
				"sh",
				"-c",
				"chmod +x /root/hermes/hermes_bootstrap.sh && /root/hermes/hermes_bootstrap.sh",
			},
		},
		noRestart,
	)
	if err != nil {
		return nil, err
	}
	m.resources[hermesContainerName] = hermesResource

	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	err = m.pool.Client.Logs(docker.LogsOptions{
		Context:           ctx,
		Container:         hermesResource.Container.ID,
		OutputStream:      &outBuf,
		ErrorStream:       &errBuf,
		InactivityTimeout: 0,
		Stderr:            true,
		Stdout:            true,
	})
	if err != nil {
		return nil, err
	}

	if m.IsDebugLogEnabled {
		t.Logf(outBuf.String())
		t.Logf(errBuf.String())
	}

	return hermesResource, nil
}

// RunICQResource runs an ICQ container. Returns the container resource and error if any.
// the name of the ICQ container is "<chain A id>-<chain B id>-relayer".
func (m *Manager) RunICQResource(t *testing.T, chainAID, quickANodeName, chainBID, quickBNodeName, icqCfgPath string) (*dockertest.Resource, error) {
	t.Helper()

	icqResource, err := m.pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       icqContainerName,
			Repository: m.ICQRepository,
			Tag:        m.ICQTag,
			NetworkID:  m.network.Network.ID,
			User:       "root:root",
			Mounts: []string{
				fmt.Sprintf("%s/:/root/icq", icqCfgPath),
			},
			ExposedPorts: []string{
				"2112",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"2112/tcp": {{HostIP: "", HostPort: "2112"}},
			},
			Env: []string{
				fmt.Sprintf("QUICK_A_E2E_CHAIN_ID=%s", chainAID),
				fmt.Sprintf("QUICK_B_E2E_CHAIN_ID=%s", chainBID),
				fmt.Sprintf("QUICK_A_E2E_VAL_HOST=%s", quickANodeName),
				fmt.Sprintf("QUICK_B_E2E_VAL_HOST=%s", quickBNodeName),
			},
			Entrypoint: []string{
				"sh",
				"-c",
				"chmod +x /root/icq/icq_bootstrap.sh && /root/icq/icq_bootstrap.sh",
			},
		},
		noRestart,
	)
	if err != nil {
		return nil, err
	}
	m.resources[icqContainerName] = icqResource

	_, err = m.pool.Client.InspectContainer(icqResource.Container.ID)
	if err != nil {
		return nil, err
	}

	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	err = m.pool.Client.Logs(docker.LogsOptions{
		Context:           ctx,
		Container:         icqResource.Container.ID,
		OutputStream:      &outBuf,
		ErrorStream:       &errBuf,
		InactivityTimeout: 0,
		Stderr:            true,
		Stdout:            true,
	})
	if err != nil {
		return nil, err
	}

	if m.IsDebugLogEnabled {
		t.Logf(outBuf.String())
		t.Logf(errBuf.String())
	}

	return icqResource, nil
}

// RunXCCLookupResource runs an XCCLookup container. Returns the container resource and error if any.
// the name of the XCCLookup container is "<chain A id>-<chain B id>-relayer".
func (m *Manager) RunXCCLookupResource(t *testing.T, chainAID, quickANodeName, chainBID, quickBNodeName, xccLookupCfgPath string) (*dockertest.Resource, error) {
	t.Helper()

	xccLookupResource, err := m.pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       xccLookupContainerName,
			Repository: m.XCCLookupRepository,
			Tag:        m.XCCLookupTag,
			NetworkID:  m.network.Network.ID,
			User:       "root:root",
			Mounts: []string{
				fmt.Sprintf("%s/:/root/xcclookup", xccLookupCfgPath),
			},
			ExposedPorts: []string{
				"3033",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"3033/tcp": {{HostIP: "", HostPort: "3033"}},
			},
			Env: []string{
				fmt.Sprintf("QUICK_A_E2E_CHAIN_ID=%s", chainAID),
				fmt.Sprintf("QUICK_B_E2E_CHAIN_ID=%s", chainBID),
				fmt.Sprintf("QUICK_A_E2E_VAL_HOST=%s", quickANodeName),
				fmt.Sprintf("QUICK_B_E2E_VAL_HOST=%s", quickBNodeName),
			},
			Entrypoint: []string{
				"sh",
				"-c",
				"chmod +x /root/xcclookup/xcc_bootstrap.sh && /root/xcclookup/xcc_bootstrap.sh",
			},
		},
		noRestart,
	)
	if err != nil {
		return nil, err
	}
	m.resources[xccLookupContainerName] = xccLookupResource

	_, err = m.pool.Client.InspectContainer(xccLookupResource.Container.ID)
	if err != nil {
		return nil, err
	}

	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	err = m.pool.Client.Logs(docker.LogsOptions{
		Context:           ctx,
		Container:         xccLookupResource.Container.ID,
		OutputStream:      &outBuf,
		ErrorStream:       &errBuf,
		InactivityTimeout: 0,
		Stderr:            true,
		Stdout:            true,
	})
	if err != nil {
		return nil, err
	}

	if m.IsDebugLogEnabled {
		t.Logf(outBuf.String())
		t.Logf(errBuf.String())
	}

	return xccLookupResource, nil
}

// RunNodeResource runs a node container. Assigns containerName to the container.
// Mounts the container on valConfigDir volume on the running host. Returns the container resource and error if any.
func (m *Manager) RunNodeResource(containerName, valCondigDir string) (*dockertest.Resource, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	runOpts := &dockertest.RunOptions{
		Name:       containerName,
		Repository: m.QuicksilversRepository,
		Tag:        m.QuicksilverTag,
		NetworkID:  m.network.Network.ID,
		User:       "root:root",
		Cmd: []string{
			"start",
		},
		Mounts: []string{
			fmt.Sprintf("%s/:/quicksilver/.quicksilverd", valCondigDir),
			fmt.Sprintf("%s/scripts:/quicksilver", pwd),
		},
	}

	resource, err := m.pool.RunWithOptions(runOpts, noRestart)
	if err != nil {
		return nil, err
	}

	m.resources[containerName] = resource

	return resource, nil
}

// RunChainInitResource runs a chain init container to initialize genesis and configs for a chain with chainId.
// The chain is to be configured with chainVotingPeriod and validators deserialized from validatorConfigBytes.
// The genesis and configs are to be mounted on the init container as volume on mountDir path.
// Returns the container resource and error if any. This method does not Purge the container. The caller
// must deal with removing the resource.
func (m *Manager) RunChainInitResource(chainID string, chainVotingPeriod int, validatorConfigBytes []byte, mountDir string, forkHeight int) (*dockertest.Resource, error) {
	votingPeriodDuration := time.Duration(chainVotingPeriod * 1000000000)
	fmt.Printf("initializing chain resource...\nRepository: %s\nTag: %s\n", m.ImageConfig.InitRepository, m.ImageConfig.InitTag)

	initResource, err := m.pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       chainID,
			Repository: m.ImageConfig.InitRepository,
			Tag:        m.ImageConfig.InitTag,
			NetworkID:  m.network.Network.ID,
			Cmd: []string{
				fmt.Sprintf("--data-dir=%s", mountDir),
				fmt.Sprintf("--chain-id=%s", chainID),
				fmt.Sprintf("--config=%s", validatorConfigBytes),
				fmt.Sprintf("--voting-period=%v", votingPeriodDuration),
				fmt.Sprintf("--fork-height=%v", forkHeight),
			},
			User: "root:root",
			Mounts: []string{
				fmt.Sprintf("%s:%s", mountDir, mountDir),
			},
		},
		noRestart,
	)
	if err != nil {
		return nil, err
	}
	return initResource, nil
}

// PurgeResource purges the container resource and returns an error if any.
func (m *Manager) PurgeResource(resource *dockertest.Resource) error {
	return m.pool.Purge(resource)
}

// GetNodeResource returns the node resource for containerName.
func (m *Manager) GetNodeResource(containerName string) (*dockertest.Resource, error) {
	resource, exists := m.resources[containerName]
	if !exists {
		return nil, fmt.Errorf("node resource not found: container name: %s", containerName)
	}
	return resource, nil
}

// GetHostPort returns the port-forwarding address of the running host
// necessary to connect to the portId exposed inside the container.
// The container is determined by containerName.
// Returns the host-port or error if any.
func (m *Manager) GetHostPort(containerName, portID string) (string, error) {
	resource, err := m.GetNodeResource(containerName)
	if err != nil {
		return "", err
	}
	return resource.GetHostPort(portID), nil
}

// RemoveNodeResource removes a node container specified by containerName.
// Returns error if any.
func (m *Manager) RemoveNodeResource(containerName string) error {
	resource, err := m.GetNodeResource(containerName)
	if err != nil {
		return err
	}
	var opts docker.RemoveContainerOptions
	opts.ID = resource.Container.ID
	opts.Force = true
	if err := m.pool.Client.RemoveContainer(opts); err != nil {
		return err
	}
	delete(m.resources, containerName)
	return nil
}

// ClearResources removes all outstanding Docker resources created by the Manager.
func (m *Manager) ClearResources() error {
	for _, resource := range m.resources {
		if err := m.pool.Purge(resource); err != nil {
			return err
		}
	}

	return m.pool.RemoveNetwork(m.network)
}

func noRestart(config *docker.HostConfig) {
	// in this case we don't want the nodes to restart on failure
	config.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
}
