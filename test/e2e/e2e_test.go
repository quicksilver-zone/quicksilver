package e2e

import (
	"fmt"
	"github.com/ingenuity-build/quicksilver/test/e2e/configurer/chain"
	"github.com/ingenuity-build/quicksilver/test/e2e/initialization"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"os"
	"path/filepath"
	"time"
)

// Reusable Checks

// CheckBalance Checks the balance of an address
func (s *IntegrationTestSuite) CheckBalance(node *chain.NodeConfig, addr, denom string, amount int64) {
	// check the balance of the contract
	s.Require().Eventually(func() bool {
		balance, err := node.QueryBalances(addr)
		s.Require().NoError(err)
		if len(balance) == 0 {
			return false
		}
		// check that the amount is in one of the balances inside the balance list
		for _, b := range balance {
			if b.Denom == denom && b.Amount.Int64() == amount {
				return true
			}
		}
		return false
	},
		1*time.Minute,
		10*time.Millisecond,
	)
}

// TestIBCTokenTransfer tests that IBC token transfers work as expected.
func (s *IntegrationTestSuite) TestIBCTokenTransfer() {
	if s.skipIBC {
		s.T().Skip("Skipping IBC tests")
	}
	chainA := s.configurer.GetChainConfig(0)
	chainB := s.configurer.GetChainConfig(1)
	chainA.SendIBC(chainB, chainB.NodeConfigs[0].PublicAddress, initialization.QuickToken)
	chainB.SendIBC(chainA, chainA.NodeConfigs[0].PublicAddress, initialization.QuickToken)
	chainA.SendIBC(chainB, chainB.NodeConfigs[0].PublicAddress, initialization.StakeToken)
	chainB.SendIBC(chainA, chainA.NodeConfigs[0].PublicAddress, initialization.StakeToken)

	_, err := chainA.GetDefaultNode()
	s.NoError(err)
}

func (s *IntegrationTestSuite) TestLargeWasmUpload() {
	chainA := s.configurer.GetChainConfig(0)
	node, err := chainA.GetDefaultNode()
	s.NoError(err)
	node.StoreWasmCode("bytecode/large.wasm", initialization.ValidatorWalletName)
}

func (s *IntegrationTestSuite) TestStateSync() {
	if s.skipStateSync {
		s.T().Skip()
	}

	chainA := s.configurer.GetChainConfig(0)
	runningNode, err := chainA.GetDefaultNode()
	s.Require().NoError(err)

	persistentPeers := chainA.GetPersistentPeers()

	stateSyncHostPort := fmt.Sprintf("%s:26657", runningNode.Name)
	stateSyncRPCServers := []string{stateSyncHostPort, stateSyncHostPort}

	// get trust height and trust hash.
	trustHeight, err := runningNode.QueryCurrentHeight()
	s.Require().NoError(err)

	trustHash, err := runningNode.QueryHashFromBlock(trustHeight)
	s.Require().NoError(err)

	stateSynchingNodeConfig := &initialization.NodeConfig{
		Name:               "state-sync",
		Pruning:            "default",
		PruningKeepRecent:  "0",
		PruningInterval:    "0",
		SnapshotInterval:   1500,
		SnapshotKeepRecent: 2,
	}

	tempDir, err := os.MkdirTemp("", "quicksilver-e2e-statesync-")
	s.Require().NoError(err)

	// configure genesis and config files for the state-syncing node.
	nodeInit, err := initialization.InitSingleNode(
		chainA.ID,
		tempDir,
		filepath.Join(runningNode.ConfigDir, "config", "genesis.json"),
		stateSynchingNodeConfig,
		trustHeight,
		trustHash,
		stateSyncRPCServers,
		persistentPeers,
	)
	s.Require().NoError(err)

	stateSynchingNode := chainA.CreateNode(nodeInit)

	// ensure that the running node has snapshots at a height > trustHeight.
	hasSnapshotsAvailable := func(syncInfo coretypes.SyncInfo) bool {
		snapshotHeight := runningNode.SnapshotInterval
		if uint64(syncInfo.LatestBlockHeight) < snapshotHeight {
			s.T().Logf("snapshot height is not reached yet, current (%d), need (%d)", syncInfo.LatestBlockHeight, snapshotHeight)
			return false
		}

		snapshots, err := runningNode.QueryListSnapshots()
		s.Require().NoError(err)

		for _, snapshot := range snapshots {
			if snapshot.Height > uint64(trustHeight) {
				s.T().Log("found state sync snapshot after trust height")
				return true
			}
		}
		s.T().Log("state sync snashot after trust height is not found")
		return false
	}
	runningNode.WaitUntil(hasSnapshotsAvailable)

	// start the state synchin node.
	err = stateSynchingNode.Run()
	s.Require().NoError(err)

	// ensure that the state synching node cathes up to the running node.
	s.Require().Eventually(func() bool {
		stateSyncNodeHeight, err := stateSynchingNode.QueryCurrentHeight()
		s.Require().NoError(err)
		runningNodeHeight, err := runningNode.QueryCurrentHeight()
		s.Require().NoError(err)
		return stateSyncNodeHeight == runningNodeHeight
	},
		3*time.Minute,
		500*time.Millisecond,
	)

	// stop the state synching node.
	err = chainA.RemoveNode(stateSynchingNode.Name)
	s.Require().NoError(err)
}
