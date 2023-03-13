package configurer

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	appconfig "github.com/ingenuity-build/quicksilver/cmd/config"
	"github.com/ingenuity-build/quicksilver/test/e2e/configurer/config"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/test/e2e/configurer/chain"
	"github.com/ingenuity-build/quicksilver/test/e2e/containers"
	"github.com/ingenuity-build/quicksilver/test/e2e/initialization"
)

type UpgradeSettings struct {
	IsEnabled  bool
	Version    string
	ForkHeight int64 // non-zero height implies that this is a fork upgrade.
}

type UpgradeConfigurer struct {
	baseConfigurer
	upgradeVersion string
	forkHeight     int64 // forkHeight > 0 implies that this is a fork upgrade. Otherwise, proposal upgrade.
}

var _ Configurer = (*UpgradeConfigurer)(nil)

func NewUpgradeConfigurer(t *testing.T, chainConfigs []*chain.Config, setupTests setupFn, containerManager *containers.Manager, upgradeVersion string, forkHeight int64) Configurer {
	return &UpgradeConfigurer{
		baseConfigurer: baseConfigurer{
			chainConfigs:     chainConfigs,
			containerManager: containerManager,
			setupTests:       setupTests,
			syncUntilHeight:  forkHeight + defaultSyncUntilHeight,
			t:                t,
		},
		forkHeight:     forkHeight,
		upgradeVersion: upgradeVersion,
	}
}

func (uc *UpgradeConfigurer) ConfigureChains() error {
	for _, chainConfig := range uc.chainConfigs {
		if err := uc.ConfigureChain(chainConfig); err != nil {
			return err
		}
	}
	return nil
}

func (uc *UpgradeConfigurer) ConfigureChain(chainConfig *chain.Config) error {
	uc.t.Logf("starting upgrade e2e infrastructure for chain-id: %s", chainConfig.ID)
	tmpDir, err := os.MkdirTemp("", "quicksilver-e2e-testnet-")
	if err != nil {
		return err
	}

	validatorConfigBytes, err := json.Marshal(chainConfig.ValidatorInitConfigs)
	if err != nil {
		return err
	}

	forkHeight := uc.forkHeight
	if forkHeight > 0 {
		forkHeight -= config.ForkHeightPreUpgradeOffset
	}

	chainInitResource, err := uc.containerManager.RunChainInitResource(chainConfig.ID, int(chainConfig.VotingPeriod), int(chainConfig.ExpeditedVotingPeriod), validatorConfigBytes, tmpDir, int(forkHeight))
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%v/%v-encode", tmpDir, chainConfig.ID)
	uc.t.Logf("serialized init file for chain-id %v: %v", chainConfig.ID, fileName)

	// loop through the reading and unmarshaling of the init file a total of maxRetries or until error is nil
	// without this, test attempts to unmarshal file before docker container is finished writing
	var initializedChain initialization.Chain
	for i := 0; i < config.MaxRetries; i++ {
		initializedChainBytes, _ := os.ReadFile(fileName)
		err = json.Unmarshal(initializedChainBytes, &initializedChain)
		if err == nil {
			break
		}

		if i == config.MaxRetries-1 {
			if err != nil {
				return err
			}
		}

		if i > 0 {
			time.Sleep(1 * time.Second)
		}
	}
	if err := uc.containerManager.PurgeResource(chainInitResource); err != nil {
		return err
	}
	uc.initializeChainConfigFromInitChain(&initializedChain, chainConfig)
	return nil
}

func (uc *UpgradeConfigurer) CreatePreUpgradeState() error {
	chainA := uc.chainConfigs[0]
	_, err := chainA.GetDefaultNode()
	if err != nil {
		return err
	}
	chainB := uc.chainConfigs[1]
	_, err = chainB.GetDefaultNode()
	if err != nil {
		return err
	}

	chainA.SendIBC(chainB, chainB.NodeConfigs[0].PublicAddress, initialization.QuickToken)
	chainB.SendIBC(chainA, chainA.NodeConfigs[0].PublicAddress, initialization.QuickToken)
	chainA.SendIBC(chainB, chainB.NodeConfigs[0].PublicAddress, initialization.StakeToken)
	chainB.SendIBC(chainA, chainA.NodeConfigs[0].PublicAddress, initialization.StakeToken)

	return nil
}

func (uc *UpgradeConfigurer) RunSetup() error {
	return uc.setupTests(uc)
}

func (uc *UpgradeConfigurer) RunUpgrade() error {
	var err error
	if uc.forkHeight > 0 {
		uc.runForkUpgrade()
	} else {
		err = uc.runProposalUpgrade()
	}
	if err != nil {
		return err
	}

	// Check if the nodes are running
	for chainIndex, chainConfig := range uc.chainConfigs {
		chain := uc.baseConfigurer.GetChainConfig(chainIndex)
		for validatorIdx := range chainConfig.NodeConfigs {
			node := chain.NodeConfigs[validatorIdx]
			// Check node status
			_, err = node.Status()
			if err != nil {
				uc.t.Errorf("node is not running after upgrade, chain-id %s, node %s", chainConfig.ID, node.Name)
				return err
			}
			uc.t.Logf("node %s upgraded successfully, address %s", node.Name, node.PublicAddress)
		}
	}
	return nil
}

func (uc *UpgradeConfigurer) runProposalUpgrade() error {
	// submit, deposit, and vote for upgrade proposal
	// prop height = current height + voting period + time it takes to submit proposal + small buffer
	for _, chainConfig := range uc.chainConfigs {
		for validatorIdx, node := range chainConfig.NodeConfigs {
			if validatorIdx == 0 {
				currentHeight, err := node.QueryCurrentHeight()
				if err != nil {
					return err
				}
				chainConfig.UpgradePropHeight = currentHeight + int64(chainConfig.VotingPeriod) + int64(config.PropSubmitBlocks) + int64(config.PropBufferBlocks)
				node.SubmitUpgradeProposal(uc.upgradeVersion, chainConfig.UpgradePropHeight, sdk.NewCoin(appconfig.BaseDenom, sdk.NewInt(config.InitialMinDeposit)))
				chainConfig.LatestProposalNumber++
				node.DepositProposal(chainConfig.LatestProposalNumber)
			}
			node.VoteYesProposal(initialization.ValidatorWalletName, chainConfig.LatestProposalNumber)
		}
	}

	// wait till all chains halt at upgrade height
	for _, chainConfig := range uc.chainConfigs {
		uc.t.Logf("waiting to reach upgrade height on chain %s", chainConfig.ID)
		chainConfig.WaitUntilHeight(chainConfig.UpgradePropHeight)
		uc.t.Logf("upgrade height reached on chain %s", chainConfig.ID)
	}

	// remove all containers so we can upgrade them to the new version
	for _, chainConfig := range uc.chainConfigs {
		for _, validatorConfig := range chainConfig.NodeConfigs {
			err := uc.containerManager.RemoveNodeResource(validatorConfig.Name)
			if err != nil {
				return err
			}
		}
	}

	// remove all containers so we can upgrade them to the new version
	for _, chainConfig := range uc.chainConfigs {
		if err := uc.upgradeContainers(chainConfig, chainConfig.UpgradePropHeight); err != nil {
			return err
		}
	}
	return nil
}

func (uc *UpgradeConfigurer) runForkUpgrade() {
	for _, chainConfig := range uc.chainConfigs {
		uc.t.Logf("waiting to reach fork height on chain %s", chainConfig.ID)
		chainConfig.WaitUntilHeight(uc.forkHeight)
		uc.t.Logf("fork height reached on chain %s", chainConfig.ID)
	}
}

func (uc *UpgradeConfigurer) upgradeContainers(chainConfig *chain.Config, propHeight int64) error {
	// upgrade containers to the locally compiled daemon
	uc.t.Logf("starting upgrade for chain-id: %s...", chainConfig.ID)
	uc.containerManager.QuicksilversRepository = containers.CurrentBranchQuickSilverRepository
	uc.containerManager.QuicksilverTag = containers.CurrentBranchQuicksilverTag

	for _, node := range chainConfig.NodeConfigs {
		if err := node.Run(); err != nil {
			return err
		}
	}

	uc.t.Logf("waiting to upgrade containers on chain %s", chainConfig.ID)
	chainConfig.WaitUntilHeight(propHeight)
	uc.t.Logf("upgrade successful on chain %s", chainConfig.ID)
	return nil
}
