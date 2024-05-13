package utils

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/cosmos/cosmos-sdk/types"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
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
