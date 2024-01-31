package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/client/cli"
)

// ProposalHandler is the community spend proposal handler.

var AddProtocolDataProposalHandler = govclient.NewProposalHandler(cli.GetCmdAddProtocolDataProposal)
