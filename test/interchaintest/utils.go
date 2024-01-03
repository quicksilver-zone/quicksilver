package interchaintest

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	cdsTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	istypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/strangelove-ventures/interchaintest/v5/chain/cosmos"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"path/filepath"
	"strings"
	"time"
)

// TxProposalv1 contains chain proposal transaction detail for gov module v1 (sdk v0.46.0+)
type TxProposalv1 struct {
	Messages []json.RawMessage `json:"messages"`
	Metadata string            `json:"metadata"`
	Deposit  string            `json:"deposit"`
	Title    string            `json:"title"`
	Summary  string            `json:"summary"`

	// SDK v50 only
	Proposer  string `json:"proposer,omitempty"`
	Expedited bool   `json:"expedited,omitempty"`
}

func SubmitProposal(ctx context.Context, c *cosmos.CosmosChain, keyName string, prop TxProposalv1) (string, error) {
	tn := c.Validators[0]
	if len(c.FullNodes) > 0 {
		tn = c.FullNodes[0]
	}

	propJson, err := json.MarshalIndent(prop, "", " ")
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(propJson)
	proposalFilename := fmt.Sprintf("%x.json", hash)

	err = tn.WriteFile(ctx, propJson, proposalFilename)
	if err != nil {
		return "", fmt.Errorf("writing param change proposal: %w", err)
	}

	proposalPath := filepath.Join(tn.HomeDir(), proposalFilename)

	command := []string{
		"gov", "submit-proposal",
		proposalPath,
		"--gas", "auto",
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return txHash, fmt.Errorf("failed to submit gov v1 proposal: %w", err)
	}

	return TxProposal(tn, txHash)
}

func TxProposal(tn *cosmos.ChainNode, txHash string) (string, error) {
	var txResp *types.TxResponse
	err := retry.Do(func() error {
		var err error
		txResp, err = authTx.QueryTx(tn.CliContext(), txHash)
		fmt.Println("Tx proposal response: ", txResp)
		return err
	},
		// retry for total of 3 seconds
		retry.Attempts(15),
		retry.Delay(200*time.Millisecond),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction %s: %w", txHash, err)
	}
	events := txResp.Events
	evtSubmitProp := "submit_proposal"
	proposalID, _ := AttributeValue(events, evtSubmitProp, "proposal_id")

	return proposalID, nil
}

func AttributeValue(events []abcitypes.Event, eventType, attrKey string) (string, bool) {
	for _, event := range events {
		if event.Type != eventType {
			continue
		}
		for _, attr := range event.Attributes {
			if string(attr.Key) == attrKey {
				return string(attr.Value), true
			}
		}
	}
	return "", false
}

func RequestStakingDelegate(ctx context.Context, c *cosmos.CosmosChain, validatorAddress, keyName, amount string) (string, error) {
	tn := c.Validators[0]
	if len(c.FullNodes) > 0 {
		tn = c.FullNodes[0]
	}
	command := []string{
		"staking", "delegate",
		validatorAddress, amount,
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return txHash, fmt.Errorf("failed to delegate: %w", err)
	}
	return txHash, nil
}

func RequestStakingUnbond(ctx context.Context, c *cosmos.CosmosChain, validatorAddress, keyName, amount string) (string, error) {
	tn := c.Validators[0]
	if len(c.FullNodes) > 0 {
		tn = c.FullNodes[0]
	}
	command := []string{
		"staking", "unbond",
		validatorAddress, amount,
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return txHash, fmt.Errorf("failed to delegate: %w", err)
	}
	return txHash, nil
}

func RequestICSRedeem(ctx context.Context, c *cosmos.CosmosChain, keyName, amount string) (string, error) {
	tn := c.Validators[0]
	if len(c.FullNodes) > 0 {
		tn = c.FullNodes[0]
	}
	command := []string{
		"interchainstaking", "redeem",
		amount, keyName,
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return txHash, fmt.Errorf("failed to redeem: %w", err)
	}
	return txHash, nil
}

func QueryZoneICAAddress(ctx context.Context, c *cosmos.CosmosChain, address, connectionID string) (string, error) {
	queryICA := []string{
		c.Config().Bin, "query", "interchain-accounts", "controller", "interchain-accounts", address, connectionID,
		"--chain-id", c.Config().ChainID,
		"--home", c.HomeDir(),
		"--node", c.GetRPCAddress(),
	}
	stdout, _, err := c.Exec(ctx, queryICA, nil)
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(string(stdout), ":", 2)
	return strings.TrimSpace(parts[1]), err
}

func QueryZones(ctx context.Context, c *cosmos.CosmosChain) ([]istypes.Zone, error) {
	stdout, _, err := c.Validators[0].ExecQuery(ctx, "interchainstaking", "zones")
	if err != nil {
		return nil, fmt.Errorf("failed to query zones: %w", err)
	}
	var zones istypes.QueryZonesResponse
	err = codec.NewLegacyAmino().UnmarshalJSON(stdout, &zones)
	if err != nil {
		return nil, err
	}
	return zones.Zones, nil
}

func QueryStakingDelegation(ctx context.Context, c *cosmos.CosmosChain, validatorAddress, keyName string) (*cdsTypes.Delegation, error) {
	stdout, _, err := c.Validators[0].ExecQuery(ctx, "staking", "delegation", keyName, validatorAddress)
	var delegation cdsTypes.DelegationResponse
	err = c.Config().EncodingConfig.Codec.UnmarshalJSON(stdout, &delegation)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get staking validators: %w", err)
	}
	return &delegation.Delegation, nil
}

func QueryStakingValidators(ctx context.Context, c *cosmos.CosmosChain) (cdsTypes.Validators, error) {
	stdout, _, err := c.Validators[0].ExecQuery(ctx, "staking", "validators")
	if err != nil {
		return nil, fmt.Errorf("failed to query staking validators: %w", err)
	}
	var validators cdsTypes.QueryValidatorsResponse
	err = c.Config().EncodingConfig.Codec.UnmarshalJSON(stdout, &validators)
	if err != nil {
		return nil, err
	}
	return validators.Validators, nil
}
