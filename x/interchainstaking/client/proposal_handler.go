package client

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/client/cli"
)

// ProposalHandler is the community spend proposal handler.
var (
	RegisterProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRegisterProposal, emptyRestHandler)
	UpdateProposalHandler   = govclient.NewProposalHandler(cli.GetCmdSubmitUpdateProposal, emptyRestHandler)
)

func emptyRestHandler(client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "unsupported-ibc-client",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for ICS proposals")
		},
	}
}
