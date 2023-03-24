package initialization

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	tmjson "github.com/tendermint/tendermint/libs/json"

	"github.com/ingenuity-build/quicksilver/test/e2e/util"
	epochtypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
)

// NodeConfig is a configuration for the node supplied from the test runner
// to initialization scripts. It should be backwards compatible with earlier
// versions. If this struct is updated, the change must be backported to earlier
// branches that might be used for upgrade testing.
type NodeConfig struct {
	Name               string // name of the config that will also be assigned to Docke container.
	Pruning            string // default, nothing, everything, or custom
	PruningKeepRecent  string // keep all of the last N states (only used with custom pruning)
	PruningInterval    string // delete old states from every Nth block (only used with custom pruning)
	SnapshotInterval   uint64 // statesync snapshot every Nth block (0 to disable)
	SnapshotKeepRecent uint32 // number of recent snapshots to keep and serve (0 to keep all)
	IsValidator        bool   // flag indicating whether a node should be a validator
}

const (
	// common
	QuickDenom          = "uqck"
	IonDenom            = "uion"
	StakeDenom          = "stake"
	AtomDenom           = "uatom"
	OsmoIBCDenom        = "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518"
	StakeIBCDenom       = "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B7787"
	E2EFeeToken         = "e2e-default-feetoken" //nolint:gosec
	UstIBCDenom         = "ibc/BE1BB42D4BE3C30D50B68D7C41DB4DFCE9678E8EF8C539F6E6A9345048894FCC"
	LuncIBCDenom        = "ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0"
	MinGasPrice         = "0.000"
	IbcSendAmount       = 3300000000
	ValidatorWalletName = "val"
	// chainA
	ChainAID      = "quick-test-a"
	QuickBalanceA = 20000000000000
	IonBalanceA   = 100000000000
	StakeBalanceA = 110000000000
	StakeAmountA  = 100000000000
	UstBalanceA   = 500000000000000
	LuncBalanceA  = 500000000000000
	// chainB
	ChainBID          = "quick-test-b"
	QuickalanceB      = 500000000000
	IonBalanceB       = 100000000000
	StakeBalanceB     = 440000000000
	StakeAmountB      = 400000000000
	GenesisFeeBalance = 100000000000
	WalletFeeBalance  = 100000000

	EpochDayDuration      = time.Second * 60
	EpochWeekDuration     = time.Second * 120
	TWAPPruningKeepPeriod = EpochDayDuration / 4

	// Denoms for testing Stride migration in v15.
	// Can be removed after v15 upgrade.
	StOsmoDenom               = "stOsmo"
	JunoDenom                 = "juno"
	StJunoDenom               = "stJuno"
	StarsDenom                = "stars"
	StStarsDenom              = "stStars"
	DefaultStrideDenomBalance = QuickBalanceA
)

var (
	StakeAmountIntA  = sdk.NewInt(StakeAmountA)
	StakeAmountCoinA = sdk.NewCoin(QuickDenom, StakeAmountIntA)
	StakeAmountIntB  = sdk.NewInt(StakeAmountB)
	StakeAmountCoinB = sdk.NewCoin(QuickDenom, StakeAmountIntB)

	// Pool balances for testing Stride migration in v15.
	// Can be removed after v15 upgrade.
	StridePoolBalances = fmt.Sprintf("%d%s,%d%s,%d%s,%d%s,%d%s", DefaultStrideDenomBalance, StOsmoDenom, DefaultStrideDenomBalance, JunoDenom, DefaultStrideDenomBalance, StJunoDenom, DefaultStrideDenomBalance, StarsDenom, DefaultStrideDenomBalance, StStarsDenom)

	InitBalanceStrA = fmt.Sprintf("%d%s,%d%s,%d%s,%d%s,%d%s", QuickBalanceA, QuickDenom, StakeBalanceA, StakeDenom, IonBalanceA, IonDenom, UstBalanceA, UstIBCDenom, LuncBalanceA, LuncIBCDenom)
	InitBalanceStrB = fmt.Sprintf("%d%s,%d%s,%d%s", QuickalanceB, QuickDenom, StakeBalanceB, StakeDenom, IonBalanceB, IonDenom)
	QuickToken      = sdk.NewInt64Coin(QuickDenom, IbcSendAmount) // 3,300quick
	StakeToken      = sdk.NewInt64Coin(StakeDenom, IbcSendAmount) // 3,300ustake
	tenQuick        = sdk.Coins{sdk.NewInt64Coin(QuickDenom, 10_000_000)}
	fiftyQuick      = sdk.Coins{sdk.NewInt64Coin(QuickDenom, 50_000_000)} //nolint:unused
	WalletFeeTokens = sdk.NewCoin(E2EFeeToken, sdk.NewInt(WalletFeeBalance))
)

func addAccount(path, moniker, amountStr string, accAddr sdk.AccAddress, forkHeight int) error {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config

	config.SetRoot(path)
	config.Moniker = moniker

	coins, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}
	coins = coins.Add(sdk.NewCoin(E2EFeeToken, sdk.NewInt(GenesisFeeBalance)))

	balances := banktypes.Balance{Address: accAddr.String(), Coins: coins.Sort()}
	genAccount := authtypes.NewBaseAccount(accAddr, nil, 0, 0)

	// TODO: Make the SDK make it far cleaner to add an account to GenesisState
	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	genDoc.InitialHeight = int64(forkHeight)

	authGenState := authtypes.GetGenesisStateFromAppState(util.Cdc, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return fmt.Errorf("failed to get accounts from any: %w", err)
	}

	if accs.Contains(accAddr) {
		return fmt.Errorf("failed to add account to genesis state; account already exists: %s", accAddr)
	}

	// Add the newInternal account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	accs = append(accs, genAccount)
	accs = authtypes.SanitizeGenesisAccounts(accs)

	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return fmt.Errorf("failed to convert accounts into any's: %w", err)
	}

	authGenState.Accounts = genAccs

	authGenStateBz, err := util.Cdc.MarshalJSON(&authGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}

	appState[authtypes.ModuleName] = authGenStateBz

	bankGenState := banktypes.GetGenesisStateFromAppState(util.Cdc, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenStateBz, err := util.Cdc.MarshalJSON(bankGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}

	appState[banktypes.ModuleName] = bankGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

func updateModuleGenesis[V proto.Message](appGenState map[string]json.RawMessage, moduleName string, protoVal V, updateGenesis func(V)) error {
	if err := util.Cdc.UnmarshalJSON(appGenState[moduleName], protoVal); err != nil {
		return err
	}
	updateGenesis(protoVal)
	newGenState := protoVal

	bz, err := util.Cdc.MarshalJSON(newGenState)
	if err != nil {
		return err
	}
	appGenState[moduleName] = bz
	return nil
}

func initGenesis(chain *internalChain, votingPeriod time.Duration, forkHeight int) error {
	// initialize a genesis file
	configDir := chain.nodes[0].configDir()
	for _, val := range chain.nodes {
		addr, err := val.keyRecord.GetAddress()
		if err != nil {
			return err
		}

		if chain.chainMeta.ID == ChainAID {
			if err := addAccount(configDir, "", InitBalanceStrA+","+StridePoolBalances, addr, forkHeight); err != nil {
				return err
			}
		} else if chain.chainMeta.ID == ChainBID {
			if err := addAccount(configDir, "", InitBalanceStrB+","+StridePoolBalances, addr, forkHeight); err != nil {
				return err
			}
		}
	}

	// copy the genesis file to the remaining validators
	for _, val := range chain.nodes[1:] {
		_, err := util.CopyFile(
			filepath.Join(configDir, "config", "genesis.json"),
			filepath.Join(val.configDir(), "config", "genesis.json"),
		)
		if err != nil {
			return err
		}
	}

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config

	config.SetRoot(chain.nodes[0].configDir())
	config.Moniker = chain.nodes[0].moniker

	genFilePath := config.GenesisFile()
	appGenState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFilePath)
	if err != nil {
		return err
	}

	err = updateModuleGenesis(appGenState, staketypes.ModuleName, &staketypes.GenesisState{}, updateStakeGenesis)
	if err != nil {
		return err
	}

	err = updateModuleGenesis(appGenState, minttypes.ModuleName, &minttypes.GenesisState{}, updateMintGenesis)
	if err != nil {
		return err
	}

	err = updateModuleGenesis(appGenState, banktypes.ModuleName, &banktypes.GenesisState{}, updateBankGenesis())
	if err != nil {
		return err
	}

	err = updateModuleGenesis(appGenState, epochtypes.ModuleName, &epochtypes.GenesisState{}, updateEpochGenesis)
	if err != nil {
		return err
	}

	err = updateModuleGenesis(appGenState, crisistypes.ModuleName, &crisistypes.GenesisState{}, updateCrisisGenesis)
	if err != nil {
		return err
	}

	err = updateModuleGenesis(appGenState, govtypes.ModuleName, &govv1types.GenesisState{}, updateGovGenesis(votingPeriod))
	if err != nil {
		return err
	}

	err = updateModuleGenesis(appGenState, genutiltypes.ModuleName, &genutiltypes.GenesisState{}, updateGenUtilGenesis(chain))
	if err != nil {
		return err
	}

	bz, err := json.MarshalIndent(appGenState, "", "  ")
	if err != nil {
		return err
	}

	genDoc.AppState = bz

	genesisJSON, err := tmjson.MarshalIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}

	// write the updated genesis file to each validator
	for _, val := range chain.nodes {
		if err := util.WriteFile(filepath.Join(val.configDir(), "config", "genesis.json"), genesisJSON); err != nil {
			return err
		}
	}
	return nil
}

func updateBankGenesis() func(s *banktypes.GenesisState) {
	return func(bankGenState *banktypes.GenesisState) {
		denomsToRegister := []string{StakeDenom, IonDenom, QuickDenom, AtomDenom, LuncIBCDenom, UstIBCDenom}
		for _, denom := range denomsToRegister {
			setDenomMetadata(bankGenState, denom)
		}
	}
}

func updateStakeGenesis(stakeGenState *staketypes.GenesisState) {
	stakeGenState.Params = staketypes.Params{
		BondDenom:         QuickDenom,
		MaxValidators:     100,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		UnbondingTime:     240000000000,
		MinCommissionRate: sdk.ZeroDec(),
	}
}

func updateMintGenesis(mintGenState *minttypes.GenesisState) {
	mintGenState.Params.MintDenom = QuickDenom
	mintGenState.Params.EpochIdentifier = "day"
}

func updateEpochGenesis(epochGenState *epochtypes.GenesisState) {
	epochGenState.Epochs = []epochtypes.EpochInfo{
		// override week epochs which are in default integrations, to be 2min
		epochtypes.NewGenesisEpochInfo("week", time.Second*120),
		// override day epochs which are in default integrations, to be 1min
		epochtypes.NewGenesisEpochInfo("day", time.Second*60),
	}
}

func updateCrisisGenesis(crisisGenState *crisistypes.GenesisState) {
	crisisGenState.ConstantFee.Denom = QuickDenom
}

func updateGovGenesis(votingPeriod time.Duration) func(*govv1types.GenesisState) {
	return func(govGenState *govv1types.GenesisState) {
		govGenState.VotingParams.VotingPeriod = &votingPeriod
		govGenState.DepositParams.MinDeposit = tenQuick
	}
}

func updateGenUtilGenesis(c *internalChain) func(*genutiltypes.GenesisState) {
	return func(genUtilGenState *genutiltypes.GenesisState) {
		// generate genesis txs
		genTxs := make([]json.RawMessage, 0, len(c.nodes))
		for _, node := range c.nodes {
			if !node.isValidator {
				continue
			}

			stakeAmountCoin := StakeAmountCoinA
			if c.chainMeta.ID != ChainAID {
				stakeAmountCoin = StakeAmountCoinB
			}
			createValmsg, err := node.buildCreateValidatorMsg(stakeAmountCoin)
			if err != nil {
				panic("genutil genesis setup failed: " + err.Error())
			}

			signedTx, err := node.signMsg(createValmsg)
			if err != nil {
				panic("genutil genesis setup failed: " + err.Error())
			}

			txRaw, err := util.Cdc.MarshalJSON(signedTx)
			if err != nil {
				panic("genutil genesis setup failed: " + err.Error())
			}
			genTxs = append(genTxs, txRaw)
		}
		genUtilGenState.GenTxs = genTxs
	}
}

func setDenomMetadata(genState *banktypes.GenesisState, denom string) {
	genState.DenomMetadata = append(genState.DenomMetadata, banktypes.Metadata{
		Description: fmt.Sprintf("Registered denom %s for e2e testing", denom),
		Display:     denom,
		Base:        denom,
		Symbol:      denom,
		Name:        denom,
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 0,
			},
		},
	})
}
