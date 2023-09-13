package client

import (
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.

var AddProtocolDataProposalHandler = govclient.NewProposalHandler(cli.GetCmdAddProtocolDataProposal)
