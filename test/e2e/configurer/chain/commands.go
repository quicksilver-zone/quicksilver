package chain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ingenuity-build/quicksilver/test/e2e/configurer/config"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/p2p"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	quicksilver "github.com/ingenuity-build/quicksilver/app"
	appconfig "github.com/ingenuity-build/quicksilver/cmd/config"
	"github.com/ingenuity-build/quicksilver/test/e2e/initialization"
)

// The value is returned as a string, so we have to unmarshal twice
type params struct {
	Key      string `json:"key"`
	Subspace string `json:"subspace"`
	Value    string `json:"value"`
}

func (n *NodeConfig) QueryGRPCGateway(path string, parameters ...string) ([]byte, error) {
	if len(parameters)%2 != 0 {
		return nil, fmt.Errorf("invalid number of parameters, must follow the format of key + value")
	}

	// add the URL for the given validator ID, and pre-pend to to path.
	hostPort, err := n.containerManager.GetHostPort(n.Name, "1317/tcp")
	require.NoError(n.t, err)
	endpoint := fmt.Sprintf("http://%s", hostPort)
	fullQueryPath := fmt.Sprintf("%s/%s", endpoint, path)

	var resp *http.Response
	require.Eventually(n.t, func() bool {
		req, err := http.NewRequest("GET", fullQueryPath, nil)
		if err != nil {
			return false
		}

		if len(parameters) > 0 {
			q := req.URL.Query()
			for i := 0; i < len(parameters); i += 2 {
				q.Add(parameters[i], parameters[i+1])
			}
			req.URL.RawQuery = q.Encode()
		}

		resp, err = http.DefaultClient.Do(req) //nolint:bodyclose
		if err != nil {
			n.t.Logf("error while executing HTTP request: %s", err.Error())
			return false
		}

		return resp.StatusCode != http.StatusServiceUnavailable
	}, time.Minute, time.Millisecond*10, "failed to execute HTTP request")

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bz))
	}
	return bz, nil
}

func (n *NodeConfig) StoreWasmCode(wasmFile, from string) {
	n.LogActionF("storing wasm code from file %s", wasmFile)
	cmd := []string{"quicksilverd", "tx", "wasm", "store", wasmFile, fmt.Sprintf("--from=%s", from), "--gas=auto", "--gas-prices=0.1uqck", "--gas-adjustment=1.3"}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully stored")
}

func (n *NodeConfig) InstantiateWasmContract(codeID, initMsg, from string) {
	n.LogActionF("instantiating wasm contract %s with %s", codeID, initMsg)
	cmd := []string{"quicksilverd", "tx", "wasm", "instantiate", codeID, initMsg, fmt.Sprintf("--from=%s", from), "--no-admin", "--label=contract"}
	n.LogActionF(strings.Join(cmd, " "))
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully initialized")
}

func (n *NodeConfig) WasmExecute(contract, execMsg, from string) {
	n.LogActionF("executing %s on wasm contract %s from %s", execMsg, contract, from)
	cmd := []string{"quicksilverd", "tx", "wasm", "execute", contract, execMsg, fmt.Sprintf("--from=%s", from)}
	n.LogActionF(strings.Join(cmd, " "))
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully executed")
}

// QueryParams extracts the params for a given subspace and key. This is done generically via json to avoid having to
// specify the QueryParamResponse type (which may not exist for all params).
func (n *NodeConfig) QueryParams(subspace, key string) string {
	cmd := []string{"quicksilverd", "query", "params", "subspace", subspace, key, "--output=json"}

	out, _, err := n.containerManager.ExecCmd(n.t, n.Name, cmd, "")
	require.NoError(n.t, err)

	result := &params{}
	err = json.Unmarshal(out.Bytes(), &result)
	require.NoError(n.t, err)
	return result.Value
}

func (n *NodeConfig) QueryGovModuleAccount() string {
	cmd := []string{"quicksilverd", "query", "auth", "module-accounts", "--output=json"}

	out, _, err := n.containerManager.ExecCmd(n.t, n.Name, cmd, "")
	require.NoError(n.t, err)
	var result map[string][]interface{}
	err = json.Unmarshal(out.Bytes(), &result)
	require.NoError(n.t, err)
	for _, acc := range result["accounts"] {
		account, ok := acc.(map[string]interface{})
		require.True(n.t, ok)
		if account["name"] == "gov" {
			moduleAccount, ok := account["base_account"].(map[string]interface{})["address"].(string)
			require.True(n.t, ok)
			return moduleAccount
		}
	}
	require.True(n.t, false, "gov module account not found")
	return ""
}

func (n *NodeConfig) SubmitParamChangeProposal(proposalJSON, from string) {
	n.LogActionF("submitting param change proposal %s", proposalJSON)
	// ToDo: Is there a better way to do this?
	wd, err := os.Getwd()
	require.NoError(n.t, err)
	localProposalFile := wd + "/scripts/param_change_proposal.json"
	f, err := os.Create(localProposalFile)
	require.NoError(n.t, err)
	_, err = f.WriteString(proposalJSON)
	require.NoError(n.t, err)
	err = f.Close()
	require.NoError(n.t, err)

	cmd := []string{"quicksilverd", "tx", "gov", "submit-proposal", "param-change", "/quicksilver/param_change_proposal.json", fmt.Sprintf("--from=%s", from)}

	_, _, err = n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)

	err = os.Remove(localProposalFile)
	require.NoError(n.t, err)

	n.LogActionF("successfully submitted param change proposal")
}

func (n *NodeConfig) SendIBCTransfer(from, recipient, amount, memo string) {
	n.LogActionF("IBC sending %s from %s to %s. memo: %s", amount, from, recipient, memo)

	cmd := []string{"quicksilverd", "tx", "ibc-transfer", "transfer", "transfer", "channel-0", recipient, amount, fmt.Sprintf("--from=%s", from), "--memo", memo}

	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)

	n.LogActionF("successfully submitted sent IBC transfer")
}

func (n *NodeConfig) FailIBCTransfer(from, recipient, amount string) {
	n.LogActionF("IBC sending %s from %s to %s", amount, from, recipient)

	cmd := []string{"quicksilverd", "tx", "ibc-transfer", "transfer", "transfer", "channel-0", recipient, amount, fmt.Sprintf("--from=%s", from)}

	_, _, err := n.containerManager.ExecTxCmdWithSuccessString(n.t, n.chainID, n.Name, cmd, "rate limit exceeded")
	require.NoError(n.t, err)

	n.LogActionF("Failed to send IBC transfer (as expected)")
}

// SwapExactAmountIn swaps tokenInCoin to get at least tokenOutMinAmountInt of the other token's pool out.
// swapRoutePoolIds is the comma separated list of pool ids to swap through.
// swapRouteDenoms is the comma separated list of denoms to swap through.
// To reproduce locally:
// docker container exec <container id> quicksilverd tx gamm swap-exact-amount-in <tokeinInCoin> <tokenOutMinAmountInt> --swap-route-pool-ids <swapRoutePoolIds> --swap-route-denoms <swapRouteDenoms> --chain-id=<id>--from=<address> --keyring-backend=test -b=block --yes --log_format=json
func (n *NodeConfig) SwapExactAmountIn(tokenInCoin, tokenOutMinAmountInt string, swapRoutePoolIDs string, swapRouteDenoms string, from string) {
	n.LogActionF("swapping %s to get a minimum of %s with pool id routes (%s) and denom routes (%s)", tokenInCoin, tokenOutMinAmountInt, swapRoutePoolIDs, swapRouteDenoms)
	cmd := []string{"quicksilverd", "tx", "gamm", "swap-exact-amount-in", tokenInCoin, tokenOutMinAmountInt, fmt.Sprintf("--swap-route-pool-ids=%s", swapRoutePoolIDs), fmt.Sprintf("--swap-route-denoms=%s", swapRouteDenoms), fmt.Sprintf("--from=%s", from)}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully swapped")
}

func (n *NodeConfig) JoinPoolExactAmountIn(tokenIn string, poolID uint64, shareOutMinAmount string, from string) {
	n.LogActionF("join-swap-extern-amount-in (%s)  (%s) from (%s), pool id (%d)", tokenIn, shareOutMinAmount, from, poolID)
	cmd := []string{"quicksilverd", "tx", "gamm", "join-swap-extern-amount-in", tokenIn, shareOutMinAmount, fmt.Sprintf("--pool-id=%d", poolID), fmt.Sprintf("--from=%s", from)}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully joined pool")
}

func (n *NodeConfig) ExitPool(from, minAmountsOut string, poolID uint64, shareAmountIn string) {
	n.LogActionF("exiting gamm pool")
	cmd := []string{"quicksilverd", "tx", "gamm", "exit-pool", fmt.Sprintf("--min-amounts-out=%s", minAmountsOut), fmt.Sprintf("--share-amount-in=%s", shareAmountIn), fmt.Sprintf("--pool-id=%d", poolID), fmt.Sprintf("--from=%s", from)}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully exited pool %d, minAmountsOut %s, shareAmountIn %s", poolID, minAmountsOut, shareAmountIn)
}

func (n *NodeConfig) SubmitUpgradeProposal(upgradeVersion string, upgradeHeight int64, initialDeposit sdk.Coin) {
	n.LogActionF("submitting upgrade proposal %s for height %d", upgradeVersion, upgradeHeight)
	cmd := []string{"quicksilverd", "tx", "gov", "submit-proposal", "software-upgrade", upgradeVersion, fmt.Sprintf("--title=\"%s upgrade\"", upgradeVersion), "--description=\"upgrade proposal submission\"", fmt.Sprintf("--upgrade-height=%d", upgradeHeight), "--upgrade-info=\"\"", "--from=val", fmt.Sprintf("--deposit=%s", initialDeposit)}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully submitted upgrade proposal")
}

func (n *NodeConfig) SubmitSuperfluidProposal(asset string, initialDeposit sdk.Coin) {
	n.LogActionF("submitting superfluid proposal for asset %s", asset)
	cmd := []string{"quicksilverd", "tx", "gov", "submit-proposal", "set-superfluid-assets-proposal", fmt.Sprintf("--superfluid-assets=%s", asset), fmt.Sprintf("--title=\"%s superfluid asset\"", asset), fmt.Sprintf("--description=\"%s superfluid asset\"", asset), "--from=val", fmt.Sprintf("--deposit=%s", initialDeposit)}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully submitted superfluid proposal for asset %s", asset)
}

func (n *NodeConfig) SubmitTextProposal(text string, initialDeposit sdk.Coin, isExpedited bool) {
	n.LogActionF("submitting text gov proposal")
	cmd := []string{"quicksilverd", "tx", "gov", "submit-proposal", "--type=text", fmt.Sprintf("--title=\"%s\"", text), "--description=\"test text proposal\"", "--from=val", fmt.Sprintf("--deposit=%s", initialDeposit)}
	if isExpedited {
		cmd = append(cmd, "--is-expedited=true")
	}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully submitted text gov proposal")
}

func (n *NodeConfig) DepositProposal(proposalNumber int) {
	n.LogActionF("depositing on proposal: %d", proposalNumber)
	deposit := sdk.NewCoin(appconfig.BaseDenom, sdk.NewInt(config.MinDepositValue)).String()

	cmd := []string{"quicksilverd", "tx", "gov", "deposit", fmt.Sprintf("%d", proposalNumber), deposit, "--from=val"}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully deposited on proposal %d", proposalNumber)
}

func (n *NodeConfig) VoteYesProposal(from string, proposalNumber int) {
	n.LogActionF("voting yes on proposal: %d", proposalNumber)
	cmd := []string{"quicksilverd", "tx", "gov", "vote", fmt.Sprintf("%d", proposalNumber), "yes", fmt.Sprintf("--from=%s", from)}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully voted yes on proposal %d", proposalNumber)
}

func (n *NodeConfig) VoteNoProposal(from string, proposalNumber int) {
	n.LogActionF("voting no on proposal: %d", proposalNumber)
	cmd := []string{"quicksilverd", "tx", "gov", "vote", fmt.Sprintf("%d", proposalNumber), "no", fmt.Sprintf("--from=%s", from)}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully voted no on proposal: %d", proposalNumber)
}

func (n *NodeConfig) BankSend(amount string, sendAddress string, receiveAddress string) {
	n.LogActionF("bank sending %s from address %s to %s", amount, sendAddress, receiveAddress)
	cmd := []string{"quicksilverd", "tx", "bank", "send", sendAddress, receiveAddress, amount, "--from=val"}
	_, _, err := n.containerManager.ExecTxCmd(n.t, n.chainID, n.Name, cmd)
	require.NoError(n.t, err)
	n.LogActionF("successfully sent bank sent %s from address %s to %s", amount, sendAddress, receiveAddress)
}

// This method also funds fee tokens from the `initialization.ValidatorWalletName` account.
// TODO: Abstract this to be a fee token provider account.
func (n *NodeConfig) CreateWallet(walletName string) string {
	n.LogActionF("creating wallet %s", walletName)
	cmd := []string{"quicksilverd", "keys", "add", walletName, "--keyring-backend=test"}
	outBuf, _, err := n.containerManager.ExecCmd(n.t, n.Name, cmd, "")
	require.NoError(n.t, err)
	re := regexp.MustCompile("osmo1(.{38})")
	walletAddr := fmt.Sprintf("%s\n", re.FindString(outBuf.String()))
	walletAddr = strings.TrimSuffix(walletAddr, "\n")
	n.LogActionF("created wallet %s, wallet address - %s", walletName, walletAddr)
	n.BankSend(initialization.WalletFeeTokens.String(), initialization.ValidatorWalletName, walletAddr)
	n.LogActionF("Sent fee tokens from %s", initialization.ValidatorWalletName)
	return walletAddr
}

func (n *NodeConfig) CreateWalletAndFund(walletName string, tokensToFund []string) string {
	return n.CreateWalletAndFundFrom(walletName, initialization.ValidatorWalletName, tokensToFund)
}

func (n *NodeConfig) CreateWalletAndFundFrom(newWalletName string, fundingWalletName string, tokensToFund []string) string {
	n.LogActionF("Sending tokens to %s", newWalletName)

	walletAddr := n.CreateWallet(newWalletName)
	for _, tokenToFund := range tokensToFund {
		n.BankSend(tokenToFund, fundingWalletName, walletAddr)
	}

	n.LogActionF("Successfully sent tokens to %s", newWalletName)
	return walletAddr
}

func (n *NodeConfig) GetWallet(walletName string) string {
	n.LogActionF("retrieving wallet %s", walletName)
	cmd := []string{"quicksilverd", "keys", "show", walletName, "--keyring-backend=test"}
	outBuf, _, err := n.containerManager.ExecCmd(n.t, n.Name, cmd, "")
	require.NoError(n.t, err)
	re := regexp.MustCompile("osmo1(.{38})")
	walletAddr := fmt.Sprintf("%s\n", re.FindString(outBuf.String()))
	walletAddr = strings.TrimSuffix(walletAddr, "\n")
	n.LogActionF("wallet %s found, waller address - %s", walletName, walletAddr)
	return walletAddr
}

func (n *NodeConfig) QueryPropStatusTimed(proposalNumber int, desiredStatus string, totalTime chan time.Duration) {
	start := time.Now()
	require.Eventually(
		n.t,
		func() bool {
			status, err := n.QueryPropStatus(proposalNumber)
			if err != nil {
				return false
			}

			return status == desiredStatus
		},
		1*time.Minute,
		10*time.Millisecond,
		"Quicksilver node failed to retrieve prop tally",
	)
	elapsed := time.Since(start)
	totalTime <- elapsed
}

type validatorInfo struct {
	Address     bytes.HexBytes
	PubKey      cryptotypes.PubKey
	VotingPower int64
}

// ResultStatus is node's info, same as Tendermint, except that we use our own
// PubKey.
type ResultStatus struct {
	NodeInfo      p2p.DefaultNodeInfo
	SyncInfo      coretypes.SyncInfo
	ValidatorInfo validatorInfo
}

func (n *NodeConfig) Status() (ResultStatus, error) {
	cmd := []string{"quicksilverd", "status"}
	_, errBuf, err := n.containerManager.ExecCmd(n.t, n.Name, cmd, "")
	if err != nil {
		return ResultStatus{}, err
	}

	cfg := quicksilver.MakeEncodingConfig()
	legacyAmino := cfg.Amino
	var result ResultStatus
	err = legacyAmino.UnmarshalJSON(errBuf.Bytes(), &result)
	fmt.Println("result", result)

	if err != nil {
		return ResultStatus{}, err
	}
	return result, nil
}
