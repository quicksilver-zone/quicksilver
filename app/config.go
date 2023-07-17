package app

// DONTCOVER

import (
	"fmt"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func DefaultConfig() network.Config {
	encCfg := MakeEncodingConfig()
	return network.Config{
		Codec:             encCfg.Marshaler,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor:    NewAppConstructor(encCfg, "quicktest-1"),
		GenesisState:      ModuleBasics.DefaultGenesis(encCfg.Marshaler),
		TimeoutCommit:     1 * time.Second / 2,
		ChainID:           "quicktest-1",
		NumValidators:     1,
		BondDenom:         sdk.DefaultBondDenom,
		MinGasPrices:      fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:     sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:     sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:      sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		CleanupDir:        true,
		SigningAlgo:       string(hd.Secp256k1Type),
		KeyringOptions:    []keyring.Option{},
	}
}

func NewAppConstructor(_ EncodingConfig, chainID string) network.AppConstructor {
	return func(val network.ValidatorI) servertypes.Application {
		return NewQuicksilver(
			val.GetCtx().Logger,
			dbm.NewMemDB(),
			nil,
			true,
			map[int64]bool{},
			DefaultNodeHome,
			0,
			wasm.EnableAllProposals,
			EmptyAppOptions{},
			GetWasmOpts(EmptyAppOptions{}),
			false,
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetChainID(chainID),
			// baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}
