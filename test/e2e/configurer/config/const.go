package config

import (
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

const (
	// ForkHeightPreUpgradeOffset if not skipping upgrade, how many blocks we allow for fork to run pre upgrade state creation.
	ForkHeightPreUpgradeOffset int64 = 60
	// PropSubmitBlocks is estimated number of blocks it takes to submit for a proposal.
	PropSubmitBlocks float32 = 10
	// PropDepositBlocks is estimated number of blocks it takes to deposit for a proposal.
	PropDepositBlocks float32 = 10
	// PropVoteBlocks is number of blocks it takes to vote for a single validator to vote for a proposal.
	PropVoteBlocks float32 = 1.2
	// PropBufferBlocks is number of blocks used as a calculation buffer.
	PropBufferBlocks float32 = 6
	// MaxRetries is max retries for json unmarshalling.
	MaxRetries = 60
)

var (
	// MinDepositValue is minimum deposit value for a proposal to enter a voting period.
	MinDepositValue = govv1types.DefaultMinDepositTokens.Int64()
	// InitialMinDeposit is minimum deposit value for proposal to be submitted.
	InitialMinDeposit = MinDepositValue / 4
)
