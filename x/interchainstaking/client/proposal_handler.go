package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/client/cli"
)

// ProposalHandler is the community spend proposal handler.
var (
	RegisterProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRegisterProposal)
	UpdateProposalHandler   = govclient.NewProposalHandler(cli.GetCmdSubmitUpdateProposal)
)
