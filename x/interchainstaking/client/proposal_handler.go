package client

import (
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	RegisterProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRegisterProposal)
	UpdateProposalHandler   = govclient.NewProposalHandler(cli.GetCmdSubmitUpdateProposal)
)
