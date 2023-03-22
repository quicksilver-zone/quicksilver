package chain

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsutils "github.com/cosmos/cosmos-sdk/x/params/client/utils"
	"github.com/stretchr/testify/require"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ingenuity-build/quicksilver/test/e2e/configurer/config"
	"github.com/ingenuity-build/quicksilver/test/e2e/containers"
	"github.com/ingenuity-build/quicksilver/test/e2e/initialization"
)

type Config struct {
	initialization.ChainMeta

	ValidatorInitConfigs []*initialization.NodeConfig
	// voting period is number of blocks it takes to deposit, 1.2 seconds per validator to vote on the prop, and a buffer.
	VotingPeriod float32
	// upgrade proposal height for chain.
	UpgradePropHeight    int64
	LatestProposalNumber int
	LatestLockNumber     int
	NodeConfigs          []*NodeConfig

	LatestCodeID int

	t                *testing.T
	containerManager *containers.Manager
}

const (
	// defaultNodeIndex to use for querying and executing transactions.
	// It is used when we are indifferent about the node we are working with.
	defaultNodeIndex = 0
	// waitUntilRepeatPauseTime is the time to wait between each check of the node status.
	waitUntilRepeatPauseTime = 2 * time.Second
	// waitUntilrepeatMax is the maximum number of times to repeat the wait until condition.
	waitUntilrepeatMax = 60

	proposalStatusPassed = "PROPOSAL_STATUS_PASSED"
)

func New(t *testing.T, containerManager *containers.Manager, id string, initValidatorConfigs []*initialization.NodeConfig) *Config {
	numVal := float32(len(initValidatorConfigs))
	return &Config{
		ChainMeta: initialization.ChainMeta{
			ID: id,
		},
		ValidatorInitConfigs: initValidatorConfigs,
		VotingPeriod:         config.PropDepositBlocks + numVal*config.PropVoteBlocks + config.PropBufferBlocks,
		t:                    t,
		containerManager:     containerManager,
	}
}

// CreateNode returns new initialized NodeConfig.
func (c *Config) CreateNode(initNode *initialization.Node) *NodeConfig {
	nodeConfig := &NodeConfig{
		Node:             *initNode,
		chainID:          c.ID,
		containerManager: c.containerManager,
		t:                c.t,
	}
	c.NodeConfigs = append(c.NodeConfigs, nodeConfig)
	return nodeConfig
}

// RemoveNode removes node and stops it from running.
func (c *Config) RemoveNode(nodeName string) error {
	for i, node := range c.NodeConfigs {
		if node.Name == nodeName {
			c.NodeConfigs = append(c.NodeConfigs[:i], c.NodeConfigs[i+1:]...)
			return node.Stop()
		}
	}
	return fmt.Errorf("node %s not found", nodeName)
}

// WaitUntilEpoch waits for the chain to reach the specified epoch.
func (c *Config) WaitUntilEpoch(epoch int64, epochIdentifier string) {
	node, err := c.GetDefaultNode()
	require.NoError(c.t, err)
	node.WaitUntil(func(_ coretypes.SyncInfo) bool {
		newEpochNumber := node.QueryCurrentEpoch(epochIdentifier)
		c.t.Logf("current epoch number is (%d), waiting to reach (%d)", newEpochNumber, epoch)
		return newEpochNumber >= epoch
	})
}

// WaitForNumEpochs waits for the chain to to go through a given number of epochs.
func (c *Config) WaitForNumEpochs(epochsToWait int64, epochIdentifier string) {
	node, err := c.GetDefaultNode()
	require.NoError(c.t, err)
	oldEpochNumber := node.QueryCurrentEpoch(epochIdentifier)
	c.WaitUntilEpoch(oldEpochNumber+epochsToWait, epochIdentifier)
}

// WaitUntilHeight waits for all validators to reach the specified height at the minimum.
// returns error, if any.
func (c *Config) WaitUntilHeight(height int64) {
	// Ensure the nodes are making progress.
	doneCondition := func(syncInfo coretypes.SyncInfo) bool {
		curHeight := syncInfo.LatestBlockHeight

		if curHeight < height {
			c.t.Logf("current block height is %d, waiting to reach: %d", curHeight, height)
			return false
		}

		return !syncInfo.CatchingUp
	}

	for _, node := range c.NodeConfigs {
		c.t.Logf("node container: %s, waiting to reach height %d", node.Name, height)
		node.WaitUntil(doneCondition)
	}
}

// WaitForNumHeights waits for all nodes to go through a given number of heights.
func (c *Config) WaitForNumHeights(heightsToWait int64) {
	node, err := c.GetDefaultNode()
	require.NoError(c.t, err)
	currentHeight, err := node.QueryCurrentHeight()
	require.NoError(c.t, err)
	c.WaitUntilHeight(currentHeight + heightsToWait)
}

func (c *Config) SendIBC(dstChain *Config, recipient string, token sdk.Coin) {
	c.t.Logf("IBC sending %s from %s to %s (%s)", token, c.ID, dstChain.ID, recipient)

	dstNode, err := dstChain.GetDefaultNode()
	require.NoError(c.t, err)

	// removes the fee token from balances for calculating the difference in other tokens
	// before and after the IBC send.
	removeFeeTokenFromBalance := func(balance sdk.Coins) sdk.Coins {
		_, feeTokenBalance := balance.Find(initialization.E2EFeeToken)
		return balance.Sub(feeTokenBalance)
	}

	balancesDstPreWithTxFeeBalance, err := dstNode.QueryBalances(recipient)
	require.NoError(c.t, err)
	balancesDstPre := removeFeeTokenFromBalance(balancesDstPreWithTxFeeBalance)
	cmd := []string{"hermes", "tx", "ft-transfer", "--dst-chain", dstChain.ID, "--src-chain", c.ID, "--src-port", "transfer", "--src-channel", "channel-0", "--amount", token.Amount.String(), fmt.Sprintf("--denom=%s", token.Denom), fmt.Sprintf("--receiver=%s", recipient), "--timeout-height-offset=1000"}
	outBuf, _, err := c.containerManager.ExecHermesCmd(c.t, cmd, "SUCCESS")
	require.NoError(c.t, err)
	if c.containerManager.IsDebugLogEnabled {
		c.t.Logf("Hermes SendIBC: %s\n", outBuf.String())
	}

	require.Eventually(
		c.t,
		func() bool {
			balancesDstPostWithTxFeeBalance, err := dstNode.QueryBalances(recipient)
			require.NoError(c.t, err)

			balancesDstPost := removeFeeTokenFromBalance(balancesDstPostWithTxFeeBalance)

			ibcCoin := balancesDstPost.Sub(balancesDstPre...)
			if ibcCoin.Len() == 1 {
				tokenPre := balancesDstPre.AmountOfNoDenomValidation(ibcCoin[0].Denom)
				tokenPost := balancesDstPost.AmountOfNoDenomValidation(ibcCoin[0].Denom)
				resPre := token.Amount
				resPost := tokenPost.Sub(tokenPre)
				return resPost.Uint64() == resPre.Uint64()
			}
			return false
		},
		5*time.Minute,
		time.Second,
		"tx not received on destination chain",
	)

	c.t.Log("successfully sent IBC tokens")
}

// GetDefaultNode returns the default node of the chain.
// The default node is the first one created. Returns error if no
// ndoes created.
func (c *Config) GetDefaultNode() (*NodeConfig, error) {
	return c.getNodeAtIndex(defaultNodeIndex)
}

// GetPersistentPeers returns persistent peers from every node
// associated with a chain.
func (c *Config) GetPersistentPeers() []string {
	peers := make([]string, len(c.NodeConfigs))
	for i, node := range c.NodeConfigs {
		peers[i] = node.PeerID
	}
	return peers
}

// GetNodeAtIndex returns the nodeIndex'th node on the chain
func (c *Config) GetNodeAtIndex(nodeIndex int) (*NodeConfig, error) {
	return c.getNodeAtIndex(nodeIndex)
}

func (c *Config) getNodeAtIndex(nodeIndex int) (*NodeConfig, error) {
	if nodeIndex > len(c.NodeConfigs) {
		return nil, fmt.Errorf("node index (%d) is greater than the number of nodes available (%d)", nodeIndex, len(c.NodeConfigs))
	}
	return c.NodeConfigs[nodeIndex], nil
}

func (c *Config) SubmitParamChangeProposal(subspace, key string, value []byte) error {
	node, err := c.GetDefaultNode()
	if err != nil {
		return err
	}

	proposal := paramsutils.ParamChangeProposalJSON{
		Title:       "Param Change",
		Description: fmt.Sprintf("Changing the %s param", key),
		Changes: paramsutils.ParamChangesJSON{
			paramsutils.ParamChangeJSON{
				Subspace: subspace,
				Key:      key,
				Value:    value,
			},
		},
		Deposit: "625000000uqck",
	}
	proposalJSON, err := json.Marshal(proposal)
	if err != nil {
		return err
	}

	node.SubmitParamChangeProposal(string(proposalJSON), initialization.ValidatorWalletName)
	c.LatestProposalNumber++

	for _, n := range c.NodeConfigs {
		n.VoteYesProposal(initialization.ValidatorWalletName, c.LatestProposalNumber)
	}

	require.Eventually(c.t, func() bool {
		status, err := node.QueryPropStatus(c.LatestProposalNumber)
		if err != nil {
			return false
		}
		return status == proposalStatusPassed
	}, time.Second*30, time.Millisecond*500)
	return nil
}
