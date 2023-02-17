package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ica "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts"
	"github.com/cosmos/ibc-go/v5/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v5/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v5/modules/core/02-client/client"

	"github.com/ingenuity-build/quicksilver/x/airdrop"
	"github.com/ingenuity-build/quicksilver/x/claimsmanager"
	"github.com/ingenuity-build/quicksilver/x/epochs"
	"github.com/ingenuity-build/quicksilver/x/interchainquery"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking"
	interchainstakingclient "github.com/ingenuity-build/quicksilver/x/interchainstaking/client"
	"github.com/ingenuity-build/quicksilver/x/mint"
	"github.com/ingenuity-build/quicksilver/x/participationrewards"
	participationrewardsclient "github.com/ingenuity-build/quicksilver/x/participationrewards/client"
	"github.com/ingenuity-build/quicksilver/x/tokenfactory"
)

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	distr.AppModuleBasic{},
	mint.AppModuleBasic{},
	gov.NewAppModuleBasic(
		[]govclient.ProposalHandler{
			paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.LegacyProposalHandler, upgradeclient.LegacyCancelProposalHandler,
			ibcclientclient.UpdateClientProposalHandler, ibcclientclient.UpgradeProposalHandler, interchainstakingclient.RegisterProposalHandler, interchainstakingclient.UpdateProposalHandler,
			participationrewardsclient.AddProtocolDataProposalHandler,
		},
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	ibc.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	ica.AppModuleBasic{},
	vesting.AppModuleBasic{},
	claimsmanager.AppModuleBasic{},
	epochs.AppModuleBasic{},
	interchainstaking.AppModuleBasic{},
	interchainquery.AppModuleBasic{},
	participationrewards.AppModuleBasic{},
	airdrop.AppModuleBasic{},
	tokenfactory.AppModuleBasic{},
	wasm.AppModuleBasic{},
)
