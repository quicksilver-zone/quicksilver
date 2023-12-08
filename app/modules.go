package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	packetforward "github.com/strangelove-ventures/packet-forward-middleware/v5/router"
	packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v5/router/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v5/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v5/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v5/modules/core/02-client/client"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"

	"github.com/quicksilver-zone/quicksilver/x/airdrop"
	airdroptypes "github.com/quicksilver-zone/quicksilver/x/airdrop/types"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager"
	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/x/epochs"
	epochstypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainquery"
	interchainquerytypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking"
	interchainstakingclient "github.com/quicksilver-zone/quicksilver/x/interchainstaking/client"
	interchainstakingtypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/mint"
	minttypes "github.com/quicksilver-zone/quicksilver/x/mint/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards"
	participationrewardsclient "github.com/quicksilver-zone/quicksilver/x/participationrewards/client"
	participationrewardstypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	"github.com/quicksilver-zone/quicksilver/x/supply"
	supplytypes "github.com/quicksilver-zone/quicksilver/x/supply/types"
	"github.com/quicksilver-zone/quicksilver/x/tokenfactory"
	tokenfactorytypes "github.com/quicksilver-zone/quicksilver/x/tokenfactory/types"
)

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
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
		packetforward.AppModuleBasic{},
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
		supply.AppModuleBasic{},
	)

	// module account permissions.
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:                 nil,
		distrtypes.ModuleName:                      nil,
		minttypes.ModuleName:                       {authtypes.Minter},
		stakingtypes.BondedPoolName:                {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:             {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:                        {authtypes.Burner},
		ibctransfertypes.ModuleName:                {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:                        nil,
		claimsmanagertypes.ModuleName:              nil,
		interchainstakingtypes.ModuleName:          {authtypes.Minter},
		interchainstakingtypes.EscrowModuleAccount: {authtypes.Burner},
		interchainquerytypes.ModuleName:            nil,
		participationrewardstypes.ModuleName:       nil,
		airdroptypes.ModuleName:                    nil,
		packetforwardtypes.ModuleName:              nil,
		wasm.ModuleName:                            {authtypes.Burner},
		tokenfactorytypes.ModuleName:               {authtypes.Minter, authtypes.Burner},
	}
)

func appModules(
	app *Quicksilver,
	encodingConfig EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Marshaler

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	return []module.AppModule{
		// SDK app modules
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		staking.NewAppModule(appCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		// ibc modules
		ibc.NewAppModule(app.IBCKeeper),
		app.TransferModule,
		app.PacketForwardModule,
		app.ICAModule,
		// Quicksilver app modules
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, app.BankKeeper),
		claimsmanager.NewAppModule(appCodec, app.ClaimsManagerKeeper),
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		interchainstaking.NewAppModule(appCodec, app.InterchainstakingKeeper),
		interchainquery.NewAppModule(appCodec, app.InterchainQueryKeeper),
		participationrewards.NewAppModule(appCodec, app.ParticipationRewardsKeeper),
		airdrop.NewAppModule(appCodec, app.AirdropKeeper),
		tokenfactory.NewAppModule(app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		supply.NewAppModule(appCodec, app.SupplyKeeper),
	}
}

// simulationModules returns modules for simulation manager
// define the order of the modules for deterministic simulations.
func simulationModules(
	app *Quicksilver,
	encodingConfig EncodingConfig,
) []module.AppModuleSimulation {
	appCodec := encodingConfig.Marshaler
	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	return []module.AppModuleSimulation{
		// SDK app modules
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		staking.NewAppModule(appCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		// ibc modules
		ibc.NewAppModule(app.IBCKeeper),
		app.TransferModule,
		app.ICAModule,
		// Quicksilver app modules
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, app.BankKeeper),
		claimsmanager.NewAppModule(appCodec, app.ClaimsManagerKeeper),
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		interchainstaking.NewAppModule(appCodec, app.InterchainstakingKeeper),
		interchainquery.NewAppModule(appCodec, app.InterchainQueryKeeper),
		participationrewards.NewAppModule(appCodec, app.ParticipationRewardsKeeper),
		airdrop.NewAppModule(appCodec, app.AirdropKeeper),
		tokenfactory.NewAppModule(app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper),
		//supply.NewAppModule(appCodec, app.SupplyKeeper),
		// wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
	}
}

/*
orderBeginBlockers tells the app's module manager how to set the order of
BeginBlockers, which are run at the beginning of every block.
Interchain Security Requirements:
During begin block slashing happens after distr.BeginBlocker so that
there is nothing left over in the validator fee pool, so as to keep the
CanWithdrawInvariant invariant.
NOTE: staking module is required if HistoricalEntries param > 0
NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
*/

func orderBeginBlockers() []string {
	return []string{
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		// Note: epochs' begin should be "real" start of epochs, we keep epochs beginblock at the beginning
		epochstypes.ModuleName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,
		interchainstakingtypes.ModuleName,
		interchainquerytypes.ModuleName, // check ordering here.
		// no-op modules
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		packetforwardtypes.ModuleName,
		claimsmanagertypes.ModuleName,
		participationrewardstypes.ModuleName,
		airdroptypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		tokenfactorytypes.ModuleName,
		wasm.ModuleName,
		supplytypes.ModuleName,
	}
}

/*
Interchain Security Requirements:
- provider.EndBlock gets validator updates from the staking module;
thus, staking.EndBlock must be executed before provider.EndBlock;
- creating a new consumer chain requires the following order,
CreateChildClient(), staking.EndBlock, provider.EndBlock;
thus, gov.EndBlock must be executed before staking.EndBlock.
*/
func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		// Note: epochs' endblock should be "real" end of epochs, we keep epochs endblock at the end
		interchainquerytypes.ModuleName,
		epochstypes.ModuleName,
		// no-op modules
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		packetforwardtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		claimsmanagertypes.ModuleName,
		interchainstakingtypes.ModuleName,
		participationrewardstypes.ModuleName,
		airdroptypes.ModuleName,
		tokenfactorytypes.ModuleName,
		wasm.ModuleName,
		supplytypes.ModuleName,
		// currently no-op.
	}
}

/*
NOTE: The genutils module must occur after staking so that pools are
properly initialized with tokens from genesis accounts.
NOTE: The genutils module must also occur after auth so that it can access the params from auth.
NOTE: Capability module must occur first so that it can initialize any capabilities
so that other modules that want to create or claim capabilities afterwards in InitChain
can do so safely.
*/
func orderInitBlockers() []string {
	return []string{
		// SDK modules
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		packetforwardtypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		icatypes.ModuleName,
		// Quicksilver modules
		epochstypes.ModuleName,
		claimsmanagertypes.ModuleName,
		interchainstakingtypes.ModuleName,
		interchainquerytypes.ModuleName,
		participationrewardstypes.ModuleName,
		airdroptypes.ModuleName,
		tokenfactorytypes.ModuleName,
		supplytypes.ModuleName,
		// wasmd
		wasm.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	}
}
