package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	cdsTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	istypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
)

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
	return &delegation.Delegation, nil
}
