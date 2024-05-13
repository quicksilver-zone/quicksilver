package utils

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cdsTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ExecuteTokenizeShares(ctx context.Context, c *cosmos.CosmosChain, keyName, validatorAddress, delegatorAddress string, amount sdk.Coin) (string, error) {
	tn := c.Validators[0]
	if len(c.FullNodes) > 0 {
		tn = c.FullNodes[0]
	}

	command := []string{
		"staking", "tokenize-share",
		validatorAddress, amount.String(), delegatorAddress, "--gas", "auto", "--broadcast-mode", "sync",
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return txHash, fmt.Errorf("failed to delegate: %w", err)
	}
	return txHash, nil
}

func ExecuteValidatorBond(ctx context.Context, c *cosmos.CosmosChain, keyName, validatorAddress string) (string, error) {
	tn := c.Validators[0]
	if len(c.FullNodes) > 0 {
		tn = c.FullNodes[0]
	}

	command := []string{
		"staking", "validator-bond",
		validatorAddress, "--gas", "auto", "--broadcast-mode", "sync",
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return txHash, fmt.Errorf("failed to delegate: %w", err)
	}
	return txHash, nil
}

// StakingQueryValidators returns all validators.
func QueryStakingValidators(ctx context.Context, c *cosmos.CosmosChain, status string) ([]stakingtypes.Validator, error) {
	conn, err := grpc.Dial(c.GetHostGRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", c.GetHostGRPCAddress(), err)
	}
	queryClient := stakingtypes.NewQueryClient(conn)
	res, err := queryClient.Validators(ctx, &stakingtypes.QueryValidatorsRequest{
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query validators: %w", err)
	}
	return res.Validators, nil
}

func QueryStakingDelegationGRPC(ctx context.Context, c *cosmos.CosmosChain, delegatorAddress, validatorAddress string) (*cdsTypes.Delegation, error) {
	conn, err := grpc.Dial(c.GetHostGRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", c.GetHostGRPCAddress(), err)
	}
	queryClient := stakingtypes.NewQueryClient(conn)
	res, err := queryClient.Delegation(ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddress,
		ValidatorAddr: validatorAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query validators: %w", err)
	}
	return &res.DelegationResponse.Delegation, nil
}

func QueryStakingParamsGRPC(ctx context.Context, c *cosmos.CosmosChain) (*cdsTypes.Params, error) {
	conn, err := grpc.Dial(c.GetHostGRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", c.GetHostGRPCAddress(), err)
	}
	queryClient := stakingtypes.NewQueryClient(conn)
	res, err := queryClient.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to query validators: %w", err)
	}
	return &res.Params, nil
}
